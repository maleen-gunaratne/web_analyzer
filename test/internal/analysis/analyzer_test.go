package analysis_test

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	"web-analyzer/internal/analysis"
)

func TestHandleAnalyze(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Valid URL Parameter", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("GET", "/?url=https://example.com", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		analysis.HandleAnalyze(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "html_version")
	})

	t.Run("Missing URL Parameter", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		analysis.HandleAnalyze(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "URL parameter is required")
	})

	t.Run("Invalid URL", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("GET", "/?url=not-a-valid-url", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		analysis.HandleAnalyze(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid URL format")
	})
}
