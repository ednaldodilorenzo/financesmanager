package model

type Planning struct {
	Name    string `json:"name"`
	Total   int64  `json:"total"`
	Planned int64  `json:"planned"`
}
