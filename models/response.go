package models

import (
	"net/http"
)

// Response модель http.Response, импортируемая в лог
type Response struct {
	Status     string         `json:"status,omitempty"`
	StatusCode int            `json:"status_code,omitempty"`
	Proto      string         `json:"proto,omitempty"`
	Location   string         `json:"location,omitempty"`
	Body       string         `json:"body,omitempty"`
	Header     http.Header    `json:"headers,omitempty"`
	Cookies    []*http.Cookie `json:"cookies,omitempty"`
}
