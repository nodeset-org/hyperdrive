package config

import (
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"time"

	"github.com/nodeset-org/hyperdrive/shared/types"
	"github.com/pbnjay/memory"
)

const localBnStepID string = "step-local-bn"

func createLocalCcStep(wiz *wizard, currentStep int, totalSteps int) *choiceWizardStep {
	// Create the button names and descriptions from the config
	clientNames := []string{"Random (Recommended)"}
	clientDescriptions := []string{"Select a client randomly to help promote the diversity of the Beacon Chain. We recommend you do this unless you have a strong reason to pick a specific client. To learn more about why client diversity is important, please visit https://clientdiversity.org for an explanation."}
	clients := []*types.ParameterOption[types.BeaconNode]{}
	for _, client := range wiz.md.Config.Hyperdrive.LocalBeaconConfig.BeaconNode.Options {
		clientNames = append(clientNames, client.Name)
		clientDescriptions = append(clientDescriptions, getAugmentedBnDescription(client.Value, client.Description))
		clients = append(clients, client)
	}

	helperText := "Please select the Beacon Node you would like to use.\n\nHighlight each one to see a brief description of it, or go to https://clientdiversity.org/ to learn more about them."

	show := func(modal *choiceModalLayout) {
		wiz.md.setPage(modal.page)
		modal.focus(0) // Catch-all for safety

		if !wiz.md.isNew {
			var bnName string
			for _, option := range wiz.md.Config.Hyperdrive.LocalBeaconConfig.BeaconNode.Options {
				if option.Value == wiz.md.Config.Hyperdrive.LocalBeaconConfig.BeaconNode.Value {
					bnName = option.Name
					break
				}
			}
			for i, clientName := range clientNames {
				if bnName == clientName {
					modal.focus(i)
					break
				}
			}
		}
	}

	done := func(buttonIndex int, buttonLabel string) {
		if buttonIndex == 0 {
			wiz.md.pages.RemovePage(randomBnPrysmID)
			wiz.md.pages.RemovePage(randomBnID)
			selectRandomBn(clients, true, wiz, currentStep, totalSteps)
		} else {
			buttonLabel = strings.TrimSpace(buttonLabel)
			selectedClient := types.BeaconNode_Unknown
			for _, client := range wiz.md.Config.Hyperdrive.LocalBeaconConfig.BeaconNode.Options {
				if client.Name == buttonLabel {
					selectedClient = client.Value
					break
				}
			}
			if selectedClient == types.BeaconNode_Unknown {
				panic(fmt.Sprintf("Local BN selection buttons didn't match any known clients, buttonLabel = %s\n", buttonLabel))
			}
			wiz.md.Config.Hyperdrive.LocalBeaconConfig.BeaconNode.Value = selectedClient
			switch selectedClient {
			//case config.ConsensusClient_Prysm:
			//	wiz.consensusLocalPrysmWarning.show()
			case types.BeaconNode_Teku:
				totalMemoryGB := memory.TotalMemory() / 1024 / 1024 / 1024
				if runtime.GOARCH == "arm64" || totalMemoryGB < 15 {
					wiz.bnLocalTekuWarning.show()
				} else {
					wiz.checkpointSyncProviderModal.show()
				}
			default:
				wiz.checkpointSyncProviderModal.show()
			}
		}
	}

	back := func() {
		wiz.modeModal.show()
	}

	return newChoiceStep(
		wiz,
		currentStep,
		totalSteps,
		helperText,
		clientNames,
		clientDescriptions,
		100,
		"Beacon Node > Selection",
		DirectionalModalVertical,
		show,
		done,
		back,
		localBnStepID,
	)
}

// Get a random client compatible with the user's hardware and EC choices.
func selectRandomBn(goodOptions []*types.ParameterOption[types.BeaconNode], includeSupermajority bool, wiz *wizard, currentStep int, totalSteps int) {
	// Get system specs
	totalMemoryGB := memory.TotalMemory() / 1024 / 1024 / 1024
	isLowPower := (totalMemoryGB < 15 || runtime.GOARCH == "arm64")

	// Filter out the clients based on system specs
	filteredClients := []types.BeaconNode{}
	for _, clientOption := range goodOptions {
		client := clientOption.Value
		switch client {
		case types.BeaconNode_Teku:
			if !isLowPower {
				filteredClients = append(filteredClients, client)
			}
		/*
			case types.BeaconNode_Prysm:
				if includeSupermajority {
					filteredClients = append(filteredClients, client)
				}
		*/
		default:
			filteredClients = append(filteredClients, client)
		}
	}

	// Select a random client
	rand.Seed(time.Now().UnixNano())
	selectedClient := filteredClients[rand.Intn(len(filteredClients))]
	wiz.md.Config.Hyperdrive.LocalBeaconConfig.BeaconNode.Value = selectedClient

	// Show the selection page
	/*
		if selectedClient == types.BeaconNode_Prysm {
			wiz.consensusLocalRandomPrysmModal = createRandomPrysmStep(wiz, currentStep, totalSteps, goodOptions)
			wiz.consensusLocalRandomPrysmModal.show()
		} else {
			wiz.consensusLocalRandomModal = createRandomStep(wiz, currentStep, totalSteps, goodOptions)
			wiz.consensusLocalRandomModal.show()
		}
	*/
	wiz.bnLocalRandomModal = createRandomBnStep(wiz, currentStep, totalSteps, goodOptions)
	wiz.bnLocalRandomModal.show()
}

// Get a more verbose client description, including warnings
func getAugmentedBnDescription(client types.BeaconNode, originalDescription string) string {
	switch client {
	/*
		case types.BeaconNode_Prysm:
			return fmt.Sprintf("%s\n\n[orange]NOTE: Prysm currently has a very high representation of the Beacon Chain. For the health of the network and the overall safety of your funds, please consider choosing a client with a lower representation. Please visit https://clientdiversity.org to learn more.", originalDescription)
	*/
	case types.BeaconNode_Teku:
		totalMemoryGB := memory.TotalMemory() / 1024 / 1024 / 1024
		if runtime.GOARCH == "arm64" || totalMemoryGB < 15 {
			return fmt.Sprintf("%s\n\n[orange]WARNING: Teku is a resource-heavy client and will likely not perform well on your system given your CPU power or amount of available RAM. We recommend you pick a lighter client instead.", originalDescription)
		}
	}

	return originalDescription
}
