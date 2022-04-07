package utils

type AppErrors struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status"`
	Code       string `json:"code"`
}
