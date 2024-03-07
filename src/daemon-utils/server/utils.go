package server

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/goccy/go-json"
	"github.com/nodeset-org/hyperdrive/shared/utils/input"
	"github.com/rocket-pool/node-manager-core/utils/log"
)

// Function for validating an argument (wraps the old CLI validators)
type ArgValidator[ArgType any] func(string, string) (ArgType, error)

// Validates an argument, ensuring it exists and can be converted to the required type
func ValidateArg[ArgType any](name string, args url.Values, impl ArgValidator[ArgType], result_Out *ArgType) error {
	// Make sure it exists
	arg, exists := args[name]
	if !exists {
		return fmt.Errorf("missing argument '%s'", name)
	}

	// Run the parser
	result, err := impl(name, arg[0])
	if err != nil {
		return err
	}

	// Set the result
	*result_Out = result
	return nil
}

// Validates an argument representing a batch of inputs, ensuring it exists and the inputs can be converted to the required type
func ValidateArgBatch[ArgType any](name string, args url.Values, batchLimit int, impl ArgValidator[ArgType], result_Out *[]ArgType) error {
	// Make sure it exists
	arg, exists := args[name]
	if !exists {
		return fmt.Errorf("missing argument '%s'", name)
	}

	// Run the parser
	result, err := input.ValidateBatch[ArgType](name, arg[0], impl)
	if err != nil {
		return err
	}

	// Make sure there aren't too many entries
	if len(result) > batchLimit {
		return fmt.Errorf("too many inputs in arg %s (provided %d, max = %d)", name, len(result), batchLimit)
	}

	// Set the result
	*result_Out = result
	return nil
}

// Validates an optional argument, converting to the required type if it exists
func ValidateOptionalArg[ArgType any](name string, args url.Values, impl ArgValidator[ArgType], result_Out *ArgType, exists_Out *bool) error {
	// Make sure it exists
	arg, exists := args[name]
	if !exists {
		if exists_Out != nil {
			*exists_Out = false
		}
		return nil
	}

	// Run the parser
	result, err := impl(name, arg[0])
	if err != nil {
		return err
	}

	// Set the result
	*result_Out = result
	if exists_Out != nil {
		*exists_Out = true
	}
	return nil
}

// Gets a string argument, ensuring that it exists in the provided vars list
func GetStringFromVars(name string, args url.Values, result_Out *string) error {
	// Make sure it exists
	arg, exists := args[name]
	if !exists {
		return fmt.Errorf("missing argument '%s'", name)
	}

	// Set the result
	*result_Out = arg[0]
	return nil
}

// Gets an optional string argument from the provided vars list
func GetOptionalStringFromVars(name string, args url.Values, result_Out *string) bool {
	// Make sure it exists
	arg, exists := args[name]
	if !exists {
		return false
	}

	// Set the result
	*result_Out = arg[0]
	return true
}

// Handles an error related to parsing the input parameters of a request
func HandleInputError(log *log.ColorLogger, w http.ResponseWriter, err error) {
	// Write out any errors
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Printlnf("[%d BAD_REQUEST] <= %s", http.StatusBadRequest, err.Error())
	}
}

// Handle routes called with an invalid method
func HandleInvalidMethod(log *log.ColorLogger, w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte{})
	log.Printlnf("[%d METHOD_NOT_ALLOWED]", http.StatusMethodNotAllowed)
}

// Handles a Node daemon response
func HandleResponse(log *log.ColorLogger, w http.ResponseWriter, response any, err error, debug bool) {
	// Write out any errors
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Printlnf("[%d INTERNAL_SERVER_ERROR] <= %s", http.StatusInternalServerError, err.Error())
	}

	// Write the serialized response
	bytes, err := json.Marshal(response)
	if err != nil {
		err = fmt.Errorf("error serializing response: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Printlnf("[%d INTERNAL_SERVER_ERROR] <= %s", http.StatusInternalServerError, err.Error())
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.Write(bytes)
		if debug {
			log.Printlnf("[%d OK] <= %s", http.StatusOK, string(bytes))
		} else {
			log.Printlnf("[%d OK]", http.StatusOK)
		}
	}
}
