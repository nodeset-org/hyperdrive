package adapter

// KeyedRequest is a request that contains a key used for authentication
type KeyedRequest struct {
	Key string `json:"key"`
}
