package utils

import (
	"context"
	"golang.org/x/net/html"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"web-analyzer/internal/models"
)

var (
	checkedLinks sync.Map // Maintain already checked links
)

func TraverseHTML(n *html.Node, analysis *models.PageAnalysis, baseURL string, linksChan chan<- models.LinkInfo) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "title":
			if n.FirstChild != nil {
				analysis.Title = n.FirstChild.Data
			}
		case "h1", "h2", "h3", "h4", "h5", "h6":
			analysis.Headings[n.Data]++
		case "a":
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					linkURL := attr.Val
					if strings.TrimSpace(linkURL) == "" || strings.HasPrefix(linkURL, "#") {
						continue
					}

					isExternal := IsExternalLink(linkURL, baseURL)
					if isExternal {
						analysis.ExternalLinks++
					} else {
						analysis.InternalLinks++
					}

					linksChan <- models.LinkInfo{
						URL:        NormalizeURL(linkURL, baseURL),
						IsExternal: isExternal,
						BaseURL:    baseURL,
					}
				}
			}
		case "form":
			isLoginForm := false
			for _, attr := range n.Attr {
				if attr.Key == "action" && strings.Contains(strings.ToLower(attr.Val), "login") {
					isLoginForm = true
					break
				}
			}

			if !isLoginForm {
				var formNode func(*html.Node) bool
				formNode = func(node *html.Node) bool {
					if node.Type == html.ElementNode && node.Data == "input" {
						for _, attr := range node.Attr {
							if attr.Key == "type" && attr.Val == "password" {
								return true
							}
						}
					}

					for c := node.FirstChild; c != nil; c = c.NextSibling {
						if formNode(c) {
							return true
						}
					}
					return false
				}
				isLoginForm = formNode(n)
			}

			if isLoginForm {
				analysis.HasLoginForm = true
			}
		case "meta":
			if analysis.MetaTags == nil {
				analysis.MetaTags = make(map[string]string)
			}

			var name, content string
			for _, attr := range n.Attr {
				if attr.Key == "name" || attr.Key == "property" {
					name = attr.Val
				} else if attr.Key == "content" {
					content = attr.Val
				}
			}

			if name != "" && content != "" {
				analysis.MetaTags[name] = content
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		TraverseHTML(c, analysis, baseURL, linksChan)
	}
}

func CheckLink(link models.LinkInfo, analysis *models.PageAnalysis) {
	if _, exists := checkedLinks.Load(link.URL); exists {
		return
	}
	checkedLinks.Store(link.URL, true)

	if !strings.HasPrefix(link.URL, "http://") && !strings.HasPrefix(link.URL, "https://") {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // keep timeout for each request
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, link.URL, nil)
	if err != nil {
		slog.Debug("error creating request", "url", link.URL, "error", err)
		analysis.Mutex.Lock()
		analysis.BrokenLinks++
		if analysis.LinksStatus != nil {
			analysis.LinksStatus[link.URL] = "Error: " + err.Error()
			analysis.Mutex.Unlock()
		}
		return
	}

	req.Header.Set("User-Agent", "WebAnalyzer/1.0")

	client := &http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow up to 10 redirects
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		slog.Debug("error checking link", "url", link.URL, "error", err)

		analysis.Mutex.Lock()
		analysis.BrokenLinks++
		if analysis.LinksStatus != nil {
			analysis.LinksStatus[link.URL] = "Error: " + err.Error()
			analysis.Mutex.Unlock()
		}
		return
	}
	defer resp.Body.Close()

	analysis.Mutex.Lock()
	if resp.StatusCode >= 400 {
		analysis.BrokenLinks++
		if analysis.LinksStatus != nil {
			analysis.LinksStatus[link.URL] = "Status: " + resp.Status
		}
	} else if analysis.LinksStatus != nil {
		analysis.LinksStatus[link.URL] = "OK"
	}
	analysis.Mutex.Unlock()
}

func IsExternalLink(linkURL, baseURL string) bool {
	if strings.HasPrefix(linkURL, "http://") || strings.HasPrefix(linkURL, "https://") {
		baseParsed, baseErr := url.Parse(baseURL)
		linkParsed, linkErr := url.Parse(linkURL)

		if baseErr == nil && linkErr == nil {
			return baseParsed.Hostname() != linkParsed.Hostname()
		}
	}
	return false
}

func NormalizeURL(linkURL, baseURL string) string {

	if strings.HasPrefix(linkURL, "http://") || strings.HasPrefix(linkURL, "https://") {
		return linkURL
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return linkURL
	}

	if strings.HasPrefix(linkURL, "/") {
		return base.Scheme + "://" + base.Host + linkURL
	} else {
		basePath := base.Path
		if !strings.HasSuffix(basePath, "/") {
			lastSlash := strings.LastIndex(basePath, "/")
			if lastSlash != -1 {
				basePath = basePath[:lastSlash+1]
			} else {
				basePath = "/"
			}
		}
		return base.Scheme + "://" + base.Host + basePath + linkURL
	}
}
