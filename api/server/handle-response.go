package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/nodeset-org/hyperdrive/api"
)

// The request couldn't complete because of a custom error
func HandleError(logger *slog.Logger, w http.ResponseWriter, code int, err error) error {
	msg := err.Error()
	return writeResponse(w, logger, code, formatError(msg))
}

// The request completed successfully
func HandleSuccess(logger *slog.Logger, w http.ResponseWriter, response any) error {
	// Serialize the response
	bytes, err := json.Marshal(response)
	if err != nil {
		err := fmt.Errorf("error serializing response: %w", err)
		msg := err.Error()
		return writeResponse(w, logger, http.StatusInternalServerError, formatError(msg))
	}

	// Write it
	logger.Debug(
		"Response body",
		"body", string(bytes),
	)
	return writeResponse(w, logger, http.StatusOK, bytes)
}

// Writes a response to an HTTP request back to the client and logs it
func writeResponse(w http.ResponseWriter, logger *slog.Logger, statusCode int, message []byte) error {
	// Prep the log attributes
	codeMsg := fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode))
	attrs := []any{
		"code", codeMsg,
		"message", string(message),
	}

	// Log the response
	logMsg := "Responded with:"
	switch statusCode {
	case http.StatusOK:
		logger.Info(logMsg, attrs...)
	case http.StatusInternalServerError:
		logger.Error(logMsg, attrs...)
	default:
		logger.Warn(logMsg, attrs...)
	}

	// Write it to the client
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, writeErr := w.Write(message)
	return writeErr
}

// JSONifies an error for responding to requests
func formatError(message string) []byte {
	msg := api.ApiResponse[any]{
		Error: message,
	}

	bytes, _ := json.Marshal(msg)
	return bytes
}
