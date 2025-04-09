package analysis

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
	"log/slog"

	"web-analyzer/internal/models"
	"web-analyzer/internal/utils"
	"web-analyzer/pkg/metrics"
)

func AnalyzePage(ctx context.Context, targetURL string) (*models.PageAnalysis, error) {

	parsedURL, err := url.Parse(targetURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") || parsedURL.Host == "" {
		return nil, fmt.Errorf("invalid URL format: %s", targetURL)
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctxWithTimeout, http.MethodGet, targetURL, nil)
	if err != nil {
		metrics.Requests.WithLabelValues("failed").Inc()
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "WebAnalyzer/1.0")

	resp, err := client.Do(req)
	if err != nil {
		metrics.Requests.WithLabelValues("failed").Inc()
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		metrics.Requests.WithLabelValues("error").Inc()
		return nil, fmt.Errorf("received non-OK status code: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		metrics.Requests.WithLabelValues("parse_error").Inc()
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	analysis := &models.PageAnalysis{
		Headings:    make(map[string]int),
		LinksStatus: make(map[string]string),
		PageSize:    resp.ContentLength,
		LoadTime:    time.Since(time.Now().Add(-10 * time.Second)).Milliseconds(), // Approximation of load Time
	}

	versionChan := make(chan string, 1)
	go func() {
		versionChan <- utils.DetectHTMLVersion(resp)
	}()

	linksChan := make(chan models.LinkInfo, 100)
	resultChan := make(chan error, 1)

	go func() {
		utils.TraverseHTML(doc, analysis, targetURL, linksChan)
		close(linksChan)
	}()

	linkWg := &sync.WaitGroup{}
	linksProcessed := 0

	maxWorkers := runtime.NumCPU() * 2 // Use CPU count = 6  determine concurrency level

	// Create a worker pool for checking links
	for i := 0; i < maxWorkers; i++ {
		linkWg.Add(1)
		go func() {
			defer linkWg.Done()
			for link := range linksChan {
				if ctx.Err() != nil {
					return
				}
				utils.CheckLink(link, analysis)
				linksProcessed++
			}
		}()
	}

	go func() {
		linkWg.Wait()
		close(resultChan)
	}()

	select {
	case <-resultChan: // successfully processed links
	case <-ctx.Done():
		return nil, fmt.Errorf("analysis cancelled or timed out: %w", ctx.Err())
	}

	analysis.HTMLVersion = <-versionChan

	metrics.Requests.WithLabelValues("success").Inc()
	metrics.LinksProcessed.Add(float64(linksProcessed))
	metrics.AnalysisTime.Observe(float64(analysis.LoadTime) / 1000.0)

	return analysis, nil
}

func HandleAnalyze(c *gin.Context) {
	startTime := time.Now()
	logger := slog.With("handler", "analyze", "requestID", c.GetString("requestID"))

	targetURL := c.Query("url")
	if targetURL == "" {
		logger.Warn("missing URL parameter")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "URL parameter is required",
			Details: "Please provide a valid URL to analyze",
		})
		return
	}

	logger.Info("starting analysis..", "url", targetURL)

	resultChan := make(chan struct {
		result *models.PageAnalysis
		err    error
	}, 1)

	go func() {
		result, err := AnalyzePage(c.Request.Context(), targetURL)
		resultChan <- struct {
			result *models.PageAnalysis
			err    error
		}{result, err}
	}()

	select {
	case res := <-resultChan:
		if res.err != nil {
			logger.Error("analysis failed..", "error", res.err)
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "Analysis failed",
				Details: res.err.Error(),
			})
			return
		}
		res.result.AnalysisDuration = time.Since(startTime).String()

		logger.Info("analysis completed..", "duration", res.result.AnalysisDuration,
			"htmlVersion", res.result.HTMLVersion, "internalLinks", res.result.InternalLinks,
			"externalLinks", res.result.ExternalLinks)

		c.JSON(http.StatusOK, res.result)

	case <-c.Request.Context().Done():
		logger.Error("request cancelled by client")
		c.JSON(http.StatusRequestTimeout, models.ErrorResponse{
			Error:   "Request cancelled",
			Details: "The analysis request was cancelled / timed out",
		})
	}
}
