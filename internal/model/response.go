package model

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ErrorResp struct {
	Message string `json:"message"`
}
