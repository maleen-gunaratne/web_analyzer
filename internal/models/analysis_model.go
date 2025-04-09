package models

type PageAnalysis struct {
	HTMLVersion      string            `json:"html_version"`
	Title            string            `json:"title"`
	Headings         map[string]int    `json:"headings"`
	InternalLinks    int               `json:"internal_links"`
	ExternalLinks    int               `json:"external_links"`
	BrokenLinks      int               `json:"broken_links"`
	HasLoginForm     bool              `json:"has_login_form"`
	PageSize         int64             `json:"page_size_bytes,omitempty"`
	LoadTime         int64             `json:"load_time_ms,omitempty"`
	LinksStatus      map[string]string `json:"links_status,omitempty"`
	AnalysisDuration string            `json:"analysis_duration,omitempty"`
	MetaTags         map[string]string `json:"meta_tags,omitempty"`
}

type LinkInfo struct {
	URL        string
	IsExternal bool
	BaseURL    string
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
