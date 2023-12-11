package api

type ApiResponse[Data any] struct {
	Data *Data `json:"data"`
}

type SuccessData struct {
	Success bool `json:"success"`
}
