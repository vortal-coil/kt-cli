package internal

import (
	"fmt"
	"github.com/kt-soft-dev/kt-cli/pkg"
	"log"
	"os"
)

const (
	// ModeLog is logging with timestamp
	ModeLog = iota
	// ModePlain is plain printing without a timestamp
	ModePlain
	// ModeNoNewline is printing like plain but without a newline
	ModeNoNewline
)

// printMode is the singleton represents current way of printing messages
var printMode = ModeLog

// SetPrintMode sets the way of printing messages. See constants like Mode* for available modes
// This mode is ignored in some cases in interactive mode
func SetPrintMode(mode int) {
	printMode = mode
	pkg.SetLogger(Print)
}

// Print prints the content with optional parameters in the way defined by printMode
func Print(content string, params ...interface{}) {
	text := fmt.Sprintf(content, params...)

	switch printMode {
	case ModePlain:
		fmt.Println(text)
	case ModeNoNewline:
		fmt.Print(text)
	default:
		log.Println(text)
	}
}

// errorLogger is the singleton for logging errors
var errorLogger = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lmsgprefix)

// PrintError prints the error message with optional parameters in the way defined by printMode
// It's a wrapper around fmt.Print that respects printMode
func PrintError(err string, params ...interface{}) {
	text := fmt.Sprintf(err, params...)

	switch printMode {
	case ModePlain:
		_, _ = fmt.Fprintln(os.Stderr, text)
	case ModeNoNewline:
		_, _ = fmt.Fprint(os.Stderr, text)
	default:
		errorLogger.Println(text)
	}
}
