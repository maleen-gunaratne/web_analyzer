package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"web-analyzer/internal/models"
)

func TestPageAnalysis(t *testing.T) {
	t.Run("JSON Serialization", func(t *testing.T) {
		// Create a sample PageAnalysis object
		analysis := models.PageAnalysis{
			HTMLVersion:   "HTML5",
			Title:         "Test Page",
			InternalLinks: 5,
			ExternalLinks: 10,
			BrokenLinks:   2,
			HasLoginForm:  true,
			PageSize:      12345,
			LoadTime:      500,
			Headings: map[string]int{
				"h1": 1,
				"h2": 3,
			},
			LinksStatus: map[string]string{
				"https://example.com": "OK",
				"https://broken.com":  "Error: connection refused",
			},
			AnalysisDuration: "1.5s",
			MetaTags: map[string]string{
				"description": "A test page",
				"keywords":    "test, page",
			},
		}

		jsonData, err := json.Marshal(analysis)
		assert.NoError(t, err)

		jsonStr := string(jsonData)
		assert.Contains(t, jsonStr, `"html_version":"HTML5"`)
		assert.Contains(t, jsonStr, `"title":"Test Page"`)
		assert.Contains(t, jsonStr, `"internal_links":5`)
		assert.Contains(t, jsonStr, `"external_links":10`)
		assert.Contains(t, jsonStr, `"broken_links":2`)
		assert.Contains(t, jsonStr, `"has_login_form":true`)

		var parsedAnalysis models.PageAnalysis
		err = json.Unmarshal(jsonData, &parsedAnalysis)
		assert.NoError(t, err)

		assert.Equal(t, analysis.HTMLVersion, parsedAnalysis.HTMLVersion)
		assert.Equal(t, analysis.Title, parsedAnalysis.Title)
		assert.Equal(t, analysis.InternalLinks, parsedAnalysis.InternalLinks)
		assert.Equal(t, analysis.ExternalLinks, parsedAnalysis.ExternalLinks)
		assert.Equal(t, analysis.HasLoginForm, parsedAnalysis.HasLoginForm)
	})
}

func TestErrorResponse(t *testing.T) {
	t.Run("JSON Serialization", func(t *testing.T) {

		errResp := models.ErrorResponse{
			Error:   "Analysis failed",
			Details: "Invalid URL format",
		}

		jsonData, err := json.Marshal(errResp)
		assert.NoError(t, err)

		jsonStr := string(jsonData)
		assert.Contains(t, jsonStr, `"error":"Analysis failed"`)
		assert.Contains(t, jsonStr, `"details":"Invalid URL format"`)

		var parsedError models.ErrorResponse
		err = json.Unmarshal(jsonData, &parsedError)
		assert.NoError(t, err)

		assert.Equal(t, errResp.Error, parsedError.Error)
		assert.Equal(t, errResp.Details, parsedError.Details)
	})
}
