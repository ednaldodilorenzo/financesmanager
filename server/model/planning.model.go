package model

type Planning struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Total       int64  `json:"total"`
	Planned     int64  `json:"planned"`
	Accumulated int64  `json:"accumulated"`
}
