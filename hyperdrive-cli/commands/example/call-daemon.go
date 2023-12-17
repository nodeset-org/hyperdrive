package example

import (
	"encoding/json"
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
)

func callDaemon(client *client.Client, cmd string) error {
	response, err := client.Api.Example.CallDaemon(cmd)
	if err != nil {
		return fmt.Errorf("error running command: %w", err)
	}

	// Check the error message
	if response.Data.Error != "" {
		return fmt.Errorf("error running command: %s", response.Data.Error)
	}

	// Parse the response
	var responseData map[string]string
	err = json.Unmarshal([]byte(response.Data.Response), &responseData)
	if err != nil {
		return fmt.Errorf("error deserializing response data: %w", err)
	}

	// Print it
	fmt.Println(responseData)
	return nil
}
