package main

import (
	"flag"
	"github.com/kt-soft-dev/kt-cli/internal"
	"github.com/kt-soft-dev/kt-cli/pkg"
	"os"
)

func main() {
	flag.Parse()
	internal.SetPrintMode(*internal.PrintModeFlag)
	pkg.SetInteractiveMode(!*internal.NotInteractive)
	internal.ScanEnv()
	isStdIn := internal.IsStdin()

	// When not in debug mode, catch panics and print them in more user-friendly way like error messages
	if !*internal.Debug {
		defer func() {
			if err := recover(); err != nil {
				internal.PrintError("%v", err)
				os.Exit(1)
			}
		}()
	}

	// globalContext, cancel := context.WithCancel(context.Background())
	config, err := internal.LoadConfig(*internal.ConfigFilename)
	if err != nil {
		internal.PrintError("Failed to load config file and/or create a new one. Exiting...")
		os.Exit(1)
	}

	if !*internal.NoConfigSave {
		// Save the config file on exit. It could change during the program execution in some cases
		defer func() {
			err = internal.SaveConfig(config, *internal.ConfigFilename)
			if err != nil {
				internal.PrintError("Failed to save config file")
			}
		}()
	}

	// Set the token from the command line flag to config
	if *internal.Auth != "" {
		config.Token = *internal.Auth
	}

	// If the token is not set, and we are not in non-interactive mode, ask for it now
	if config.Token == "" && !*internal.NotInteractive {
		internal.ActionAskForToken(config)
	}

	switch {
	case *internal.Method != "":
		internal.ActionApiCall(config)

	case *internal.Ping:
		internal.ActionPing()

	case *internal.Upload != "" || isStdIn:
		internal.ActionUpload(config, isStdIn)

	case *internal.Download != "":
		internal.ActionDownload(config)

	case *internal.GetKeys != "":
		internal.ActionGetKeys(config)

	case *internal.FilesList != "":
		internal.ActionFilesList(config)

	default:
		internal.ActionDefault(config)
	}
}
