package constellation

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/nodeset-org/hyperdrive/clients"
)


var ConstellationCmd = &cobra.Command{
    Use:   "constellation initialize",
    Short: "todo",
    Run: func(cmd *cobra.Command, args []string) {
        // TODO: Fetch URL from user-settings.yml
        status, error := eth2_client.CheckStatus("https://eth.llamarpc.com")

        if !status {
            fmt.Printf("Could not connect to ETH2 Client: %s\n", error)
            return
        }

        fmt.Println("Hyperdrive Constellation successfully initialized.")
  },
}
