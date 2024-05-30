package service

import (
	"runtime/debug"
	"testing"

	"github.com/nodeset-org/hyperdrive-daemon/shared"
	"github.com/stretchr/testify/require"
)

// Test getting the server version
func TestServerVersion(t *testing.T) {
	// Set up panic handling and cleanup
	defer func() {
		r := recover()
		if r != nil {
			debug.PrintStack()
			fail("Recovered from panic: %v", r)
		} else {
			testMgr.RevertToBaseline()
		}
	}()

	version := shared.HyperdriveVersion

	// Run the round-trip test
	response, err := apiClient.Service.Version()
	require.NoError(t, err)
	require.Equal(t, version, response.Data.Version)
	t.Logf("Received correct version: %s", version)
}
