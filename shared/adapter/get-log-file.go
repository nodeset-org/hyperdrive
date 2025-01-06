package adapter

import "context"

const (
	GetLogFileCommandString string = HyperdriveModuleCommand + " get-log-file"
)

// Request format for `get-log-file`
type GetLogFileRequest struct {
	KeyedRequest

	// The log file source to retrieve
	Source string `json:"source"`
}

// Response format for `get-log-file`
type GetLogFileResponse struct {
	// The path to the log file
	Path string `json:"path"`
}

// Get a log file path from the adapter
func (c *AdapterClient) GetLogFile(ctx context.Context, source string) (*GetLogFileResponse, error) {
	request := &GetLogFileRequest{
		KeyedRequest: KeyedRequest{
			Key: c.key,
		},
		Source: source,
	}
	response := &GetLogFileResponse{}
	err := runCommand(c, ctx, GetLogFileCommandString, request, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
