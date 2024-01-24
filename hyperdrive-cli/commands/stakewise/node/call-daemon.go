package node

import (
	"encoding/json"
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
)

func uploadDepositData(client *client.Client, cmd string) error {
	response, err := client.Api.Node.UploadDepositData(cmd)
	if err != nil {
		return fmt.Errorf("error running command: %w", err)
	}

	// Parse the response
	var responseData map[string]any
	err = json.Unmarshal([]byte(response.Data.Response), &responseData)
	if err != nil {
		return fmt.Errorf("error deserializing response data: %w", err)
	}

	// Print it
	fmt.Println(responseData)
	return nil
}
