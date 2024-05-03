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
	notInteractive := flag.Bool("no-interactive", false, "Do not ask for any input, use default values")
	noConfigSave := flag.Bool("no-save", false, "Do not save the config file on exit (including token)")
	auth := flag.String("token", "", "Set auth token for future requests (will be saved in config file)")
	pretty := flag.Bool("pretty", false, "Pretty-print JSON responses")
	passwd := flag.String("passwd", "", "Set password for encryption/decryption. Also you can use environment variable KT_CLI_PASSWD")

	// Actions to perform
	method := flag.String("act.method", "", "Call API method")
	ping := flag.Bool("act.ping", false, "Check if API is alive")
	download := flag.String("act.download", "", "Download file by file ID")
	downloadPath := flag.String("act.download.path", ".", "Set path to save downloaded file")
	upload := flag.String("act.upload", "", "Upload file by path; stdin is also supported")
	uploadName := flag.String("act.upload.name", "", "Set file name for upload (required for stdin)")
	uploadDisk := flag.String("act.upload.disk", "", "Set disk for upload")
	uploadFolder := flag.String("act.upload.folder", "", "Set folder for upload")
	// @todo replace file method

	params := flag.String("params", "", "Set API method key=value parameters separated by space (format: k=v k=v k=v...)")
	flag.Parse()

	internal.SetPrintMode(*printMode)
	pkg.SetInteractiveMode(!*notInteractive)
	if *passwd == "" {
		*passwd = os.Getenv("KT_CLI_PASSWD")
	}
	isStdIn := internal.IsStdin()

	// When not in debug mode, catch panics and print them in more user-friendly way like error messages
	if !*debug {
		defer func() {
			if err := recover(); err != nil {
				internal.PrintError("%v", err)
				os.Exit(1)
			}
		}()
	}

	// globalContext, cancel := context.WithCancel(context.Background())
	config, err := internal.LoadConfig(*configFilename)
	if err != nil {
		internal.PrintError("Failed to load config file and/or create a new one. Exiting...")
		os.Exit(1)
	}

	if !*noConfigSave {
		// Save the config file on exit. It could change during the program execution in some cases
		defer func() {
			err = internal.SaveConfig(config, *configFilename)
			if err != nil {
				internal.PrintError("Failed to save config file")
			}
		}()
	}

	if *auth != "" {
		config.Token = *auth
	}

	switch {
	case *method != "":
		paramsMap := internal.ParseKeyValues(*params)
		resp, err := pkg.ApiRequest(config.Token, *method, paramsMap)
		err = internal.GetActualError(resp, err)
		if err != nil {
			internal.PrintError(err.Error())
			return
		}

		internal.Print(internal.JsonToString(resp.Result, *pretty))

	case *ping:
		if pkg.CheckApiAlive() {
			internal.Print("API is alive")
		} else {
			internal.PrintError("API is not alive")
		}

	case *upload != "" || isStdIn:
		var reader io.Reader
		var name string

		if isStdIn {
			name = *uploadName
			if name == "" {
				internal.PrintError("File name is required for stdin upload. Use -act.upload.name flag")
				return
			}
			reader = os.Stdin
		} else {
			path := *upload
			if path == "" {
				path = pkg.ScanOrDefault("Enter file path: ", "")
				if path == "" {
					internal.PrintError("File path is required")
					return
				}
			}

			fileInfo, err := os.Stat(path)
			if err != nil {
				internal.PrintError("Failed to access file")
				return
			}
			if fileInfo.IsDir() {
				// @todo directory uploading
				internal.PrintError("Directory uploading is not supported yet")
				return
			}

			file, err := os.Open(path)
			if err != nil {
				internal.PrintError("Failed to open file")
				return
			}

			if *uploadName != "" {
				name = *uploadName
			} else {
				name = file.Name()
			}

			reader = file
		}

		_, err := pkg.UploadFile(config.Token, name, "", *uploadDisk, *uploadFolder, &pkg.CryptoInfo{Password: *passwd}, reader)
		if err != nil {
			internal.PrintError(err.Error())
			return
		}

	case *download != "":
		savePath := strings.TrimSpace(*downloadPath)
		if savePath == "" {
			internal.PrintError("Save path is required")
			return
		} else if savePath == "." {
			internal.Print("Save path is set to current directory. You can change it by -act.download.path flag")
		}

		// @todo streaming download for big files
		var buffer bytes.Buffer
		writer := bufio.NewWriter(&buffer)
		name, _, err := pkg.DownloadFile(config.Token, *download, writer, &pkg.CryptoInfo{Password: *passwd})
		if err != nil {
			internal.PrintError(err.Error())
			return
		}

		pathInfo, err := os.Stat(savePath)
		if err == nil && pathInfo.IsDir() {
			savePath = savePath + string(os.PathSeparator) + name
		}

		out, err := os.Create(savePath)
		if err != nil {
			internal.PrintError("Failed to create file %s", savePath)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, &buffer)
		if err != nil {
			internal.PrintError("Failed to save file %s", savePath)
		}

	default:
		// Usually, in case of empty method and non-empty token,
		// we should take this as a request to validate and store the token
		if *auth != "" {
			id, err := pkg.CheckToken(*auth)
			if err != nil {
				internal.PrintError("Failed to check token")
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
