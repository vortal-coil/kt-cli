package pkg

import (
	"fmt"
)

var isInteractive bool

func init() {
	// By default, the library is not interactive. Interactive mode is used for CLI applications when you need to ask for user input
	isInteractive = false
}

func SetInteractiveMode(interactive bool) {
	isInteractive = interactive
}

// ScanOrDefault scans user input and returns it. If the input is empty, it returns the default value.
// Input is not scanned in non-interactive mode and the default value is returned
func ScanOrDefault(prompt, defaultValue string) (input string) {
	if !isInteractive {
		return defaultValue
	}

	// We don't use current logger here, because we want to print the prompt without newlines
	fmt.Print(prompt)

	_, err := fmt.Scanln(&input)
	if err != nil {
		return defaultValue
	}

	return
}
