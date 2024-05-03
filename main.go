package main

import (
	"bufio"
	"bytes"
	"flag"
	"github.com/kt-soft-dev/kt-cli/internal"
	"github.com/kt-soft-dev/kt-cli/pkg"
	"io"
	"os"
	"strings"
)

func main() {
	// Do not catch panics and prettify them, show debug info
	debug := flag.Bool("debug", false, "Enable debug mode")

	// Global flags
	configFilename := flag.String("config", "config.yaml", "Set config file path")
	printMode := flag.Int("output", internal.ModeLog, "Output mode (0 - log with timestamp, 1 - plain log, 2 - no newline)")
	noConfigSave := flag.Bool("no-save", false, "Do not save the config file on exit (including token)")
	auth := flag.String("token", "", "Set auth token for future requests (will be saved in config file)")
	pretty := flag.Bool("pretty", false, "Pretty-print JSON responses")
	passwd := flag.String("passwd", "", "Set password for encryption. Also you can use environment variable KT_CLI_PASSWD")

	// Actions to perform
	method := flag.String("act.method", "", "Call API method")
	ping := flag.Bool("act.ping", false, "Check if API is alive")
	download := flag.String("act.download", "", "Download file by file ID")
	downloadPath := flag.String("act.download.path", ".", "Set path to save downloaded file")

	params := flag.String("params", "", "Set API method key=value parameters separated by space (format: k=v k=v k=v...)")
	flag.Parse()

	internal.SetPrintMode(*printMode)
	if *passwd == "" {
		*passwd = os.Getenv("KT_CLI_PASSWD")
	}

	// When not in debug mode, catch panics and print them in more user-friendly way like error messages
	if !*debug {
		defer func() {
			if err := recover(); err != nil {
				internal.Print("%v", err)
				os.Exit(1)
			}
		}()
	}

	// globalContext, cancel := context.WithCancel(context.Background())
	config, err := internal.LoadConfig(*configFilename)
	if err != nil {
		internal.Print("Failed to load config file and/or create a new one. Exiting...")
		os.Exit(1)
	}

	if !*noConfigSave {
		// Save the config file on exit. It could change during the program execution in some cases
		defer func() {
			err = internal.SaveConfig(config, *configFilename)
			if err != nil {
				internal.Print("Failed to save config file")
			}
		}()
	}

	if *auth != "" {
		config.Token = *auth
	}

	switch {
	case *method != "":
		paramsMap := pkg.ParseKeyValues(*params)
		resp, err := pkg.ApiRequest(config.Token, *method, paramsMap)
		err = pkg.GetActualError(resp, err)
		if err != nil {
			internal.Print(err.Error())
			return
		}

		internal.Print(pkg.JsonToString(resp.Result, *pretty))

	case *ping:
		if pkg.CheckApiAlive() {
			internal.Print("API is alive")
		} else {
			internal.Print("API is not alive")
		}

	case *download != "":
		savePath := strings.TrimSpace(*downloadPath)
		if savePath == "" {
			internal.Print("Save path is required")
			return
		} else if savePath == "." {
			internal.Print("Save path is set to current directory. You can change it by -act.download.path flag")
		}

		// @todo streaming download for big files
		var buffer bytes.Buffer
		writer := bufio.NewWriter(&buffer)
		name, _, err := pkg.DownloadFile(config.Token, *download, writer, &pkg.CryptoInfo{Password: *passwd})
		if err != nil {
			internal.Print(err.Error())
			return
		}

		pathInfo, err := os.Stat(savePath)
		if err == nil && pathInfo.IsDir() {
			savePath = savePath + string(os.PathSeparator) + name
		}

		out, err := os.Create(savePath)
		if err != nil {
			internal.Print("Failed to create file %s", savePath)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, &buffer)
		if err != nil {
			internal.Print("Failed to save file %s", savePath)
		}

	default:
		// Usually, in case of empty method and non-empty token,
		// we should take this as a request to validate and store the token
		if *auth != "" {
			id, err := pkg.CheckToken(*auth)
			if err != nil {
				internal.Print("Failed to check token")
				return
			}

			internal.Print("Logged in as user id %s", id)
			config.Token = *auth
			config.UserID = id
			// Config will be saved because of the deferring above (if no -no-save flag is set)
			return
		}

		flag.PrintDefaults()
		os.Exit(0)
	}
}
