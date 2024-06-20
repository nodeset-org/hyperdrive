# hyperdrive-daemon
Base daemon for the Hyperdrive service

This README is still a work in progress.


## Testing

Running the integration tests requires an externally-managed instance of [Hardhat](https://hardhat.org/hardhat-runner/docs/getting-started#overview) installed and running. Among other things, Hardhat is a special standalone Execution client that simulates an Execution layer blockchain. It provides testing functions such as the ability to take snapshots of the EVM, mine new blocks, fast forward time, and so on.

We recommend using the [Hardhat start script](https://github.com/nodeset-org/osha/blob/main/hardhat/start.sh) provided in the [OSHA repository](https://github.com/nodeset-org/osha) if you don't already have experience running Hardhat.


### Configuration

When running the tests, make sure you follow two important rules:
1. Set the `HARDHAT_URL` environment variable to the URL of your instance prior to running tests
2. Set package parallelization to 1 (so multiple tests don't run concurrently, which will break Hardhat's snapshotting system)

For example:

`HARDHAT_URL="http://localhost:8545" go test -p 1 ./...`

During testing you'll notice Hardhat's node will print many event messages such as `evm_snapshot`. This is expected and part of the test suite, as each test that modifies the EL state will snapshot Hardhat prior to execution and revert to the snapshot at the end.