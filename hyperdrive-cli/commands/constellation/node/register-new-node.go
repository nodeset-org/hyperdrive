package node

import (
	"fmt"
	"math/big"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/urfave/cli/v2"
)

const (
	inputTimezoneLocation string = "timezone-location"
	inputBondAmount       string = "bond-amount"
	inputMinimumNodeFee   string = "node-fee"
	inputPrefix           string = "prefix"
	inputThreadsFlag      string = "threads"
	inputNodeAddressFlag  string = "node-address"
	inputSaltFlag         string = "salt"
)

type Result struct {
	Salt    *big.Int
	Address common.Address
}

func registerNewNode(c *cli.Context) error {
	// cs := client.NewConstellationClientFromCtx(c)
	// hd, err := client.NewHyperdriveClientFromCtx(c)
	// timezone := c.String(inputTimezoneLocation)
	// if timezone == "" {
	// 	timezone = clituils.Prompt("Please enter the timezone and location of the node (e.g. 'Europe/London'):", `^[a-zA-Z]+\/[a-zA-Z]+$`, "Invalid string format")
	// }

	// bondAmount := c.String(inputBondAmount)
	// if bondAmount == "" {
	// 	bondAmount = clituils.Prompt("Please enter the bond amount as an integer in gwei(e.g. 1000):", `^[0-9]+$`, "Invalid amount format")
	// }

	// minimumNodeFee := c.String(inputMinimumNodeFee)
	// if minimumNodeFee == "" {
	// 	minimumNodeFee = clituils.Prompt("Please enter the minimum node fee as an integer in gwei(e.g. 1000):", `^[0-9]+$`, "Invalid amount format")
	// }

	// prefix := c.String(inputPrefix)
	// if prefix == "" {
	// 	prefix = clituils.Prompt("Please enter the prefix for the vanity address (e.g. '0xb0b0'):", `^(0x[a-fA-F0-9]*|)$`, "Invalid prefix format")
	// }
	// seed := 1
	// fmt.Printf("prefix: %s", prefix)

	// // Generate vanity address from prefix
	// if !strings.HasPrefix(prefix, "0x") {
	// 	return fmt.Errorf("prefix must start with 0x")
	// }
	// targetPrefix, success := big.NewInt(0).SetString(prefix, 0)
	// if !success {
	// 	return fmt.Errorf("invalid prefix: %s", prefix)
	// }

	// // Get the starting salt
	// saltString := c.String(inputSaltFlag)
	// var salt *big.Int
	// if saltString == "" {
	// 	salt = big.NewInt(0)
	// } else {
	// 	salt, success = big.NewInt(0).SetString(saltString, 0)
	// 	if !success {
	// 		return fmt.Errorf("invalid starting salt: %s", salt)
	// 	}
	// }

	// // Get the core count
	// threads := c.Int(inputThreadsFlag)
	// if threads == 0 {
	// 	threads = runtime.GOMAXPROCS(0)
	// } else if threads < 0 {
	// 	threads = 1
	// } else if threads > runtime.GOMAXPROCS(0) {
	// 	threads = runtime.GOMAXPROCS(0)
	// }
	// // Get the node address
	// nodeAddressStr := c.String(inputNodeAddressFlag)
	// if nodeAddressStr == "" {
	// 	return fmt.Errorf("node address is required")
	// }

	// // Get the vanity generation artifacts
	// // Move to constellation
	// vanityArtifacts, err := cs.Api.Node.GetVanityArtifacts(nodeAddressStr)
	// if err != nil {
	// 	return err
	// }

	// // Set up some variables
	// nodeAddress := vanityArtifacts.Data.NodeAddress.Bytes()
	// minipoolFactoryAddress := vanityArtifacts.Data.MinipoolFactoryAddress
	// initHash := vanityArtifacts.Data.InitHash.Bytes()
	// shiftAmount := uint(42 - len(prefix))

	// // Run the search
	// fmt.Printf("Running with %d threads.\n", threads)

	// wg := new(sync.WaitGroup)
	// wg.Add(threads)
	// stop := false
	// stopPtr := &stop

	// // Create a channel to receive results from goroutines
	// resultChan := make(chan Result, 1) // Buffer of 1 to allow non-blocking send on find

	// // Start a goroutine to listen for the first result
	// var result Result
	// done := make(chan bool)
	// go func() {
	// 	result = <-resultChan
	// 	done <- true
	// }()

	// for i := 0; i < threads; i++ {
	// 	saltOffset := big.NewInt(int64(i))
	// 	workerSalt := big.NewInt(0).Add(salt, saltOffset)

	// 	go func(i int, workerSalt *big.Int) {
	// 		defer wg.Done()
	// 		localSalt, localAddress := runWorker(i == 0, stopPtr, targetPrefix, nodeAddress, minipoolFactoryAddress, initHash, workerSalt, int64(threads), shiftAmount)
	// 		if localSalt != nil {
	// 			select {
	// 			case resultChan <- Result{Salt: localSalt, Address: localAddress}:
	// 				fmt.Printf("Worker %d found salt and address: %s, %s\n", i, localSalt, localAddress.Hex())
	// 			default:
	// 			}
	// 		}
	// 	}(i, workerSalt)
	// }

	// wg.Wait()
	// close(resultChan) // Ensure no more writes to the channel

	// select {
	// case result := <-resultChan:
	// 	fmt.Printf("Found salt and address: %s, %s\n", result.Salt, result.Address.Hex())
	// 	// Continue processing with result...
	// default:
	// 	return nil
	// }

	// // Fetch from  HD Daemon isntead of hardcode:
	// expectedMinipoolResponse, err := cs.Api.Node.GetExpectedMinipoolAddress(result.Address.Hex(), result.Salt)
	// if err != nil {
	// 	fmt.Printf("Failed to get expected minipool address: %s", err.Error())
	// 	return err
	// }
	// depositData, err := hd.Api.Wallet.GenerateDepositData(expectedMinipoolResponse.Data.ExpectedMinipoolAddress)
	// if err != nil {
	// 	fmt.Printf("Failed to generate deposit data: %s", err.Error())
	// 	return err
	// }

	// validatorPubkey := depositData.Data.PublicKey
	// validatorSignature := depositData.Data.Signature
	// depositDataRoot := depositData.Data.DepositDataRoot

	// cc := client.NewConstellationClientFromCtx(c)
	// bondAmountInt, err := strconv.Atoi(bondAmount)
	// if err != nil {
	// 	fmt.Printf("Failed to convert bond amount to int: %s", err.Error())
	// 	return err
	// }
	// minimumNodeFeeInt, err := strconv.Atoi(minimumNodeFee)
	// if err != nil {
	// 	fmt.Printf("Failed to convert minimum node fee to int: %s", err.Error())
	// 	return err
	// }

	// // Build the TX
	// resp, err := cc.Api.Node.RegisterNewNode(timezone, bondAmountInt, minimumNodeFeeInt, seed, validatorPubkey, validatorSignature, depositDataRoot)
	// if err != nil {
	// 	fmt.Printf("Failed to register new node: %s", err.Error())
	// 	return err
	// }

	// // Run the TX
	// registered, err := tx.HandleTx(c, hd, resp.Data.TxInfo,
	// 	"Are you sure you want to register a new node?",
	// 	"registering new node",
	// 	"Registering the new node...",
	// )
	// if err != nil {
	// 	return err
	// }
	// if !registered {
	// 	return nil
	// }

	// // Log & return
	// fmt.Println("New node registered successfully.")
	return nil
}

func runWorker(report bool, stop *bool, targetPrefix *big.Int, nodeAddress []byte, minipoolManagerAddress common.Address, initHash []byte, salt *big.Int, increment int64, shiftAmount uint) (*big.Int, common.Address) {
	saltBytes := [32]byte{}
	hashInt := big.NewInt(0)
	incrementInt := big.NewInt(increment)
	hasher := crypto.NewKeccakState()
	nodeSalt := common.Hash{}
	addressResult := common.Hash{}

	// Set up the reporting ticker if requested
	var ticker *time.Ticker
	var tickerChan chan struct{}
	lastSalt := big.NewInt(0).Set(salt)
	if report {
		start := time.Now()
		reportInterval := 5 * time.Second
		ticker = time.NewTicker(reportInterval)
		tickerChan = make(chan struct{})
		go func() {
			for {
				select {
				case <-ticker.C:
					delta := big.NewInt(0).Sub(salt, lastSalt)
					deltaFloat, suffix := humanize.ComputeSI(float64(delta.Uint64()) / 5.0)
					deltaString := humanize.FtoaWithDigits(deltaFloat, 2) + suffix
					fmt.Printf("At salt 0x%x... %s (%s salts/sec)\n", salt, time.Since(start), deltaString)
					lastSalt.Set(salt)
				case <-tickerChan:
					ticker.Stop()
					return
				}
			}
		}()
	}

	// Run the main salt finder loop
	for {
		if *stop {
			return nil, common.Address{}
		}

		// Some speed optimizations -
		// This block is the fast way to do `nodeSalt := crypto.Keccak256Hash(nodeAddress, saltBytes)`
		salt.FillBytes(saltBytes[:])
		hasher.Write(nodeAddress)
		hasher.Write(saltBytes[:])
		hasher.Read(nodeSalt[:])
		hasher.Reset()

		// This block is the fast way to do `crypto.CreateAddress2(minipoolManagerAddress, nodeSalt, initHash)`
		// except instead of capturing the returned value as an address, we keep it as bytes. The first 12 bytes
		// are ignored, since they are not part of the resulting address.
		//
		// Because we didn't call CreateAddress2 here, we have to call common.BytesToAddress below, but we can
		// postpone that until we find the correct salt.
		hasher.Write([]byte{0xff})
		hasher.Write(minipoolManagerAddress.Bytes())
		hasher.Write(nodeSalt[:])
		hasher.Write(initHash)
		hasher.Read(addressResult[:])
		hasher.Reset()

		hashInt.SetBytes(addressResult[12:])
		hashInt.Rsh(hashInt, shiftAmount*4)
		if hashInt.Cmp(targetPrefix) == 0 {
			if report {
				close(tickerChan)
			}
			address := common.BytesToAddress(addressResult[12:])
			return salt, address
		}
		salt.Add(salt, incrementInt)
	}
}
