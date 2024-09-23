package minipool

import (
	"fmt"
	"math/big"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli/v2"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
)

var (
	vanityPrefixFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "prefix",
		Aliases: []string{"p"},
		Usage:   "The prefix of the minipool address to search for (must start with 0x)",
	}

	vanitySaltFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "salt",
		Aliases: []string{"s"},
		Usage:   "The starting salt to search from (must start with 0x)",
	}

	vanityThreadsFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "threads",
		Aliases: []string{"t"},
		Usage:   "The number of threads to use for searching (defaults to your CPU thread count)",
	}

	vanityAddressFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "node-address",
		Aliases: []string{"n"},
		Usage:   "The node address to search for (leave blank to use the local node)",
	}
)

func findVanitySalt(c *cli.Context) error {
	// Get RP client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}
	cs, err := client.NewConstellationClientFromCtx(c, hd)
	if err != nil {
		return err
	}

	// Get the target prefix
	prefix := c.String(vanityPrefixFlag.Name)
	if prefix == "" {
		prefix = utils.Prompt("Please specify the minipool address prefix you would like to search for (must start with 0x):", "^0x[0-9a-fA-F]+$", "Invalid hex string")
	}
	if !strings.HasPrefix(prefix, "0x") {
		return fmt.Errorf("Prefix must start with 0x.")
	}
	targetPrefix, success := big.NewInt(0).SetString(prefix, 0)
	if !success {
		return fmt.Errorf("Invalid prefix: %s", prefix)
	}

	// Get the starting salt
	saltString := c.String(vanitySaltFlag.Name)
	var salt *big.Int
	if saltString == "" {
		salt = big.NewInt(0)
	} else {
		salt, success = big.NewInt(0).SetString(saltString, 0)
		if !success {
			return fmt.Errorf("Invalid starting salt: %s", salt)
		}
	}

	// Get the core count
	threads := c.Int(vanityThreadsFlag.Name)
	if threads == 0 {
		threads = runtime.GOMAXPROCS(0)
	} else if threads < 0 {
		threads = 1
	} else if threads > runtime.GOMAXPROCS(0) {
		threads = runtime.GOMAXPROCS(0)
	}

	// Get the node address
	nodeAddressStr := c.String(vanityAddressFlag.Name)
	if nodeAddressStr == "" {
		nodeAddressStr = "0"
	}

	// Get the vanity generation artifacts
	vanityArtifacts, err := cs.Api.Minipool.GetVanityArtifacts(nodeAddressStr)
	if err != nil {
		return err
	}

	// Set up some variables
	subNodeAddress := vanityArtifacts.Data.SubNodeAddress.Bytes()
	superNodeAddress := vanityArtifacts.Data.SuperNodeAddress.Bytes()
	minipoolFactoryAddress := vanityArtifacts.Data.MinipoolFactoryAddress
	initHash := vanityArtifacts.Data.InitHash.Bytes()
	shiftAmount := uint(42 - len(prefix))

	// Run the search
	fmt.Printf("Running with %d threads.\n", threads)

	wg := new(sync.WaitGroup)
	wg.Add(threads)
	stop := false
	stopPtr := &stop

	// Spawn worker threads
	start := time.Now()
	for i := 0; i < threads; i++ {
		saltOffset := big.NewInt(int64(i))
		workerSalt := big.NewInt(0).Add(salt, saltOffset)

		go func(i int) {
			foundSalt, foundAddress := runWorker(i == 0, stopPtr, targetPrefix, subNodeAddress, superNodeAddress, minipoolFactoryAddress, initHash, workerSalt, int64(threads), shiftAmount)
			if foundSalt != nil {
				fmt.Printf("Found on thread %d: salt 0x%x = %s\n", i, foundSalt, foundAddress.Hex())
				*stopPtr = true
			}
			wg.Done()
		}(i)
	}

	// Wait for the workers to finish and print the elapsed time
	wg.Wait()
	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Printf("Finished in %s\n", elapsed)

	// Return
	return nil

}

func runWorker(report bool, stop *bool, targetPrefix *big.Int, subNodeAddress []byte, superNodeAddress []byte, minipoolManagerAddress common.Address, initHash []byte, salt *big.Int, increment int64, shiftAmount uint) (*big.Int, common.Address) {
	saltBytes := [32]byte{}
	hashInt := big.NewInt(0)
	incrementInt := big.NewInt(increment)
	hasher := crypto.NewKeccakState()
	nodeSalt := common.Hash{}
	internalSaltHash := common.Hash{}
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

		// Prep the internal salt by combining with the subnode address
		salt.FillBytes(saltBytes[:])
		hasher.Write(saltBytes[:])
		hasher.Write(subNodeAddress)
		_, err := hasher.Read(internalSaltHash[:])
		if err != nil {
			panic(err)
		}
		hasher.Reset()

		// Some speed optimizations -
		// This block is the fast way to do `nodeSalt := crypto.Keccak256Hash(superNodeAddress, internalSaltHash)`
		hasher.Write(superNodeAddress)
		hasher.Write(internalSaltHash[:])
		_, err = hasher.Read(nodeSalt[:])
		if err != nil {
			panic(err)
		}
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
		_, err = hasher.Read(addressResult[:])
		if err != nil {
			panic(err)
		}
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
