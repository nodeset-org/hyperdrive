# hyperdrive-daemon
Base daemon for the Hyperdrive service

This README is still a work in progress.


## Testing

Running the integration tests requires an externally-managed instance of [Hardhat](https://hardhat.org/hardhat-runner/docs/getting-started#overview) installed and running. Among other things, Hardhat is a special standalone Execution client that simulates an Execution layer blockchain. It provides testing functions such as the ability to take snapshots of the EVM, mine new blocks, fast forward time, and so on.


### Setting up Hardhat

Hardhat requires NodeJS to run. If you don't already have a NodeJS environment, we recommend [installing nvm](https://github.com/nvm-sh/nvm) (the Node version manager); the README there has installation instructions. Once installed, Hardhat recommends using Node v20 which can be set up via these commands:
``` 
nvm install 20
```
```
nvm use 20
```

Next, head to the `internal/hardhat` dir and run the following to set up Hardhat in your local copy of the repository:

```
npm ci
```

Finally, run the following to start a local instance of Hardhat's EVM runner:

```
npx hardhat node --port 8545
```


### Setting Environment Variables

`hyperdrive-daemon` uses some environment variables to configure its test manager:

- `HARDHAT_URL` specifies the URL to use for your Hardhat instance (e.g., `http://localhost:8545`)

If you are using VS Code as your editor, we recommend setting this in `.vscode/settings.json`:
```
{
    "go.testEnvVars": {
        "HARDHAT_URL": "http://localhost:8545",
    }
}
```

This is required for the editor's built-in test explorer to be able to run the tests properly.