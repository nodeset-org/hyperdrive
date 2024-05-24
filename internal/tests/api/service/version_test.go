package service

import (
	"testing"

	"github.com/nodeset-org/hyperdrive-daemon/shared"
	"github.com/stretchr/testify/require"
)

// Test getting the server version
func TestServerVersion(t *testing.T) {
	version := shared.HyperdriveVersion

	// Run the round-trip test
	response, err := apiClient.Service.Version()
	require.NoError(t, err)
	require.Equal(t, version, response.Data.Version)
	t.Logf("Received correct version: %s", version)
}
