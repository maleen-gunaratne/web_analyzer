package utils_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"web-analyzer/internal/utils"
)

func TestDetectHTMLVersion(t *testing.T) {
	testCases := []struct {
		name           string
		doctype        string
		contentType    string
		expectedResult string
	}{
		{"HTML5", "<!DOCTYPE html><html>", "text/html", "HTML5"},
		{"HTML 4.01 Strict", `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN">`, "text/html", "HTML 4.01 Strict"},
		{"HTML 4.01 Transitional", `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">`, "text/html", "HTML 4.01 Transitional"},
		{"HTML 4.01 Frameset", `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Frameset//EN">`, "text/html", "HTML 4.01 Frameset"},
		{"XHTML 1.0 Strict", `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN">`, "text/html", "XHTML 1.0 Strict"},
		{"XHTML 1.0 Transitional", `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN">`, "text/html", "XHTML 1.0 Transitional"},
		{"XHTML 1.0 Frameset", `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Frameset//EN">`, "text/html", "XHTML 1.0 Frameset"},
		{"XHTML 1.1", `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN">`, "text/html", "XHTML 1.1"},
		{"XHTML Content-Type", `<html>`, "application/xhtml+xml", "XHTML (Content-Type based detection)"},
		{"Non-standard DOCTYPE", `<!DOCTYPE html SYSTEM "about:legacy-compat">`, "text/html", "HTML (Non-standard DOCTYPE)"},
		{"No DOCTYPE", `<html>`, "text/html", "HTML (No DOCTYPE)"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(tc.doctype + "<html><body>Test</body></html>")),
				Header:     make(http.Header),
			}
			resp.Header.Set("Content-Type", tc.contentType)

			// Detect HTML version
			version := utils.DetectHTMLVersion(resp)
			assert.Equal(t, tc.expectedResult, version)
		})
	}
}
