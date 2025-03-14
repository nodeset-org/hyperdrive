package api

type ApiResponse[Data any] struct {
	Data  *Data  `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

type SuccessData struct{}
