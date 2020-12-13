package models

import (
	"net/http"
)

// Request модель http.Request, импортируемая в лог
type Request struct {
	Method     string         `json:"method,omitempty"`
	Proto      string         `json:"proto,omitempty"`
	URL        string         `json:"url,omitempty"`
	RequestURI string         `json:"request_uri,omitempty"`
	Body       string         `json:"body,omitempty"`
	Header     http.Header    `json:"headers,omitempty"`
	Cookies    []*http.Cookie `json:"cookies,omitempty"`
	RemoteAddr string         `json:"remote_address,omitempty"`
	UserAgent  string         `json:"user_agent,omitempty"`
	Referer    string         `json:"referer,omitempty"`
}
