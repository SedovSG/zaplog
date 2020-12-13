package models

// Zaplog модель
type Zaplog struct {
	Level    string      `json:"level"`
	Caller   string      `json:"caller"`
	Function string      `json:"function"`
	Message  string      `json:"message"`
	Method   string      `json:"method"`
	URL      string      `json:"url"`
	Code     int         `json:"code"`
	Request  Request     `json:"request"`
	Response Response    `json:"response"`
	Extend   interface{} `json:"extend"`
}
