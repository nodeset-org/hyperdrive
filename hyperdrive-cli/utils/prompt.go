package utils

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"golang.org/x/term"
)

// Prompt for user input
func Prompt(initialPrompt string, expectedFormat string, incorrectFormatPrompt string) string {

	// Print initial prompt
	fmt.Println(initialPrompt)

	// Get valid user input
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan(); !regexp.MustCompile(expectedFormat).MatchString(scanner.Text()); scanner.Scan() {
		fmt.Println("")
		fmt.Println(incorrectFormatPrompt)
	}
	fmt.Println("")

	// Return user input
	return scanner.Text()

}

// Prompt for confirmation
func Confirm(initialPrompt string) bool {
	response := Prompt(fmt.Sprintf("%s [y/n]", initialPrompt), "(?i)^(y|yes|n|no)$", "Please answer 'y' or 'n'")
	return (strings.ToLower(response[:1]) == "y")
}

// Prompt for 'I agree' confirmation (used on important questions to avoid a quick 'y' response from the user)
func ConfirmWithIAgree(initialPrompt string) bool {
	response := Prompt(fmt.Sprintf("%s [Type 'I agree' or 'n']", initialPrompt), "(?i)^(i agree|n|no)$", "Please answer 'I agree' or 'n'")
	return (len(response) == 7 && strings.ToLower(response[:7]) == "i agree")
}

// Prompt for user selection
func Select(initialPrompt string, options []string) (int, string) {

	// Get prompt
	prompt := initialPrompt
	for i, option := range options {
		prompt += fmt.Sprintf("\n%d: %s", (i + 1), option)
	}

	// Get expected response format
	optionNumbers := []string{}
	for i := range options {
		optionNumbers = append(optionNumbers, strconv.Itoa(i+1))
	}
	expectedFormat := fmt.Sprintf("^(%s)$", strings.Join(optionNumbers, "|"))

	// Prompt user
	response := Prompt(prompt, expectedFormat, "Please enter a number corresponding to an option")

	// Get selected option
	index, _ := strconv.Atoi(response)
	selectedIndex := index - 1
	selectedOption := options[selectedIndex]

	// Return
	return selectedIndex, selectedOption

}

// Prompts the user to verify that there is nobody looking over their shoulder before printing sensitive information.
func ConfirmSecureSession(warning string) bool {
	if !Confirm(fmt.Sprintf("%s%s%s\nAre you sure you want to continue?", terminal.ColorYellow, warning, terminal.ColorReset)) {
		fmt.Println("Cancelled.")
		return false
	}

	return true
}

// Prompt for password input
func PromptPassword(initialPrompt string, expectedFormat string, incorrectFormatPrompt string) string {
	// Print initial prompt
	fmt.Println(initialPrompt)

	// Get valid user input
	var input string
	var init bool
	for !init || !regexp.MustCompile(expectedFormat).MatchString(input) {

		// Incorrect format
		if init {
			fmt.Println("")
			fmt.Println(incorrectFormatPrompt)
		} else {
			init = true
		}

		// Read password
		if bytes, err := term.ReadPassword(syscall.Stdin); err != nil {
			fmt.Println(fmt.Errorf("error reading password: %w", err))
		} else {
			input = string(bytes)
		}

	}
	fmt.Println("")
	return input
}
