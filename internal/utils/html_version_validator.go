package utils

import (
	"bufio"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func DetectHTMLVersion(resp *http.Response) string {

	reader := bufio.NewReader(resp.Body)

	documentTypeBuffer, err := reader.Peek(4096) // Limit to initial 4KB
	if err != nil && err != io.EOF && err != bufio.ErrBufferFull {
		return "Unknown (detection error)"
	}

	snippet := string(documentTypeBuffer)

	contentType := resp.Header.Get("Content-Type")
	isXHTML := strings.Contains(contentType, "application/xhtml+xml")

	switch {
	case html5Regex.MatchString(snippet):
		return "HTML5"
	case html401StrictRegex.MatchString(snippet):
		return "HTML 4.01 Strict"
	case html401TransRegex.MatchString(snippet):
		return "HTML 4.01 Transitional"
	case html401FrameRegex.MatchString(snippet):
		return "HTML 4.01 Frameset"
	case xhtml10StrictRegex.MatchString(snippet):
		return "XHTML 1.0 Strict"
	case xhtml10TransRegex.MatchString(snippet):
		return "XHTML 1.0 Transitional"
	case xhtml10FrameRegex.MatchString(snippet):
		return "XHTML 1.0 Frameset"
	case xhtml11Regex.MatchString(snippet):
		return "XHTML 1.1"
	case isXHTML:
		return "XHTML (Content-Type based detection)"
	case strings.Contains(snippet, "<!DOCTYPE html"):
		return "HTML (Non-standard DOCTYPE)"
	default:
		if strings.Contains(strings.ToLower(snippet), "<html") {
			return "HTML (No DOCTYPE)"
		}
		return "Unknown"
	}
}

var (
	html5Regex         = regexp.MustCompile(`(?i)<!DOCTYPE\s+html>`)
	html401StrictRegex = regexp.MustCompile(`(?i)<!DOCTYPE\s+HTML\s+PUBLIC\s+"-//W3C//DTD HTML 4\.01//EN"`)
	html401TransRegex  = regexp.MustCompile(`(?i)<!DOCTYPE\s+HTML\s+PUBLIC\s+"-//W3C//DTD HTML 4\.01 Transitional//EN"`)
	html401FrameRegex  = regexp.MustCompile(`(?i)<!DOCTYPE\s+HTML\s+PUBLIC\s+"-//W3C//DTD HTML 4\.01 Frameset//EN"`)
	xhtml10StrictRegex = regexp.MustCompile(`(?i)<!DOCTYPE\s+html\s+PUBLIC\s+"-//W3C//DTD XHTML 1\.0 Strict//EN"`)
	xhtml10TransRegex  = regexp.MustCompile(`(?i)<!DOCTYPE\s+html\s+PUBLIC\s+"-//W3C//DTD XHTML 1\.0 Transitional//EN"`)
	xhtml10FrameRegex  = regexp.MustCompile(`(?i)<!DOCTYPE\s+html\s+PUBLIC\s+"-//W3C//DTD XHTML 1\.0 Frameset//EN"`)
	xhtml11Regex       = regexp.MustCompile(`(?i)<!DOCTYPE\s+html\s+PUBLIC\s+"-//W3C//DTD XHTML 1\.1//EN"`)
)
