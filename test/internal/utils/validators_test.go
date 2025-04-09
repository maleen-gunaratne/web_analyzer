package utils_test

import (
	"bytes"
	"golang.org/x/net/html"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"web-analyzer/internal/models"
	"web-analyzer/internal/utils"
)

func TestTraverseHTML(t *testing.T) {
	t.Run("HTML Element Parsing", func(t *testing.T) {

		htmlContent := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Test Page</title>
			<meta name="description" content="Test description">
		</head>
		<body>
			<h1>Heading 1</h1>
			<h2>Heading 2</h2>
			<h2>Another H2</h2>
			<a href="https://external.com">External Link</a>
			<a href="/internal">Internal Link</a>
			<form action="/login">
				<input type="text" name="username">
				<input type="password" name="password">
			</form>
		</body>
		</html>
		`
		doc, err := html.Parse(bytes.NewReader([]byte(htmlContent)))
		require.NoError(t, err)

		analysis := &models.PageAnalysis{
			Headings:    make(map[string]int),
			LinksStatus: make(map[string]string),
			MetaTags:    make(map[string]string),
		}
		linksChan := make(chan models.LinkInfo, 10)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range linksChan {
			}
		}()

		baseURL := "https://example.com"
		utils.TraverseHTML(doc, analysis, baseURL, linksChan)
		close(linksChan)
		wg.Wait()

		assert.Equal(t, "Test Page", analysis.Title)
		assert.Equal(t, 1, analysis.Headings["h1"])
		assert.Equal(t, 2, analysis.Headings["h2"])
		assert.Equal(t, 1, analysis.ExternalLinks)
		assert.Equal(t, 1, analysis.InternalLinks)
		assert.True(t, analysis.HasLoginForm)
		assert.Equal(t, "Test description", analysis.MetaTags["description"])
	})
}

func TestIsExternalLink(t *testing.T) {
	baseURL := "https://example.com"

	testCases := []struct {
		name     string
		link     string
		expected bool
	}{
		{"Absolute External URL", "https://external.com", true},
		{"Absolute Internal URL", "https://example.com/page", false},
		{"Absolute Internal URL with Path", "https://example.com/page/subpage", false},
		{"Relative URL", "/page", false},
		{"Root Relative URL", "/", false},
		{"Relative URL with Path", "page/subpage", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.IsExternalLink(tc.link, baseURL)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestNormalizeURL(t *testing.T) {
	baseURL := "https://example.com/base/path"

	testCases := []struct {
		name     string
		link     string
		expected string
	}{
		{"Absolute URL", "https://external.com", "https://external.com"},
		{"Root Relative URL", "/page", "https://example.com/page"},
		{"Relative URL", "subpage", "https://example.com/base/subpage"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.NormalizeURL(tc.link, baseURL)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCheckLink(t *testing.T) {
	t.Run("Link Checking", func(t *testing.T) {

		analysis := &models.PageAnalysis{
			LinksStatus: make(map[string]string),
		}

		link := models.LinkInfo{
			URL:        "mailto:test@example.com",
			IsExternal: true,
			BaseURL:    "https://example.com",
		}

		initialBrokenLinks := analysis.BrokenLinks
		utils.CheckLink(link, analysis)

		assert.Equal(t, initialBrokenLinks, analysis.BrokenLinks)

		link = models.LinkInfo{
			URL:        "http://invalid-url-that-will-fail",
			IsExternal: true,
			BaseURL:    "https://example.com",
		}

		utils.CheckLink(link, analysis)

		assert.Equal(t, initialBrokenLinks+1, analysis.BrokenLinks)
		assert.Contains(t, analysis.LinksStatus[link.URL], "Error")
	})
}
