package status

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/urfave/cli/v2"
)

func camelToSnake(input string) string {
	var result strings.Builder
	for i, r := range input {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func getNodeStatus(c *cli.Context) error {
	sw := client.NewStakewiseClientFromCtx(c)
	response, err := sw.Api.Status.GetValidatorStatuses()
	if err != nil {
		fmt.Printf("error fetching validator statuses: %v\n", err)
		return err
	}

	fmt.Printf("Beacon Statuses:\n")
	for pubKey, status := range response.Data.BeaconStatus {
		fmt.Printf("%v: %v\n", pubKey, status)
	}

	fmt.Printf("\n\nNodeset Statuses:\n")
	for pubKey, status := range response.Data.NodesetStatus {
		fmt.Printf("%v: %v\n", pubKey, camelToSnake(string(status)))
	}

	return nil
}
