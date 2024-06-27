package api_v2

// All responses from the NodeSet API will have this format
// `message` may or may not be populated (but should always be populated if `ok` is false)
// `data` should be populated if `ok` is true, and will be omitted if `ok` is false
type NodeSetResponse[DataType any] struct {
	OK      bool     `json:"ok"`
	Message string   `json:"message,omitempty"`
	Data    DataType `json:"data,omitempty"`
	Error   string   `json:"error,omitempty"`
}
