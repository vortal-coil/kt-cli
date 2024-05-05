package internal

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/kt-soft-dev/kt-cli/pkg"
	"github.com/rodaine/table"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"strings"
)

// Actions represent the available CLI commands. Each action is a function that can be called from the CLI
// and perform some operations. Each action can have its own flags and parameters.
// The actions and parameters are defined in the flags.go file.
// The actions are called from the main.go file and use global state without returning any values.

// ActionPing checks if the API is alive and responds to requests
func ActionPing() {
	if pkg.CheckApiAlive() {
		Print("API is alive")
	} else {
		PrintError("API is not alive")
	}
}

func ActionDefault(config *Config) {
	// Usually, in case of empty method and non-empty token,
	// we should take this as a request to validate and store the token
	if *Auth != "" {
		_ = CheckTokenAndAssign(config.Token, config)
		Print("Token is validated and saved")
		// Config will be saved because of the deferring above (if no -no-save flag is set)
		return
	}

	flag.PrintDefaults()
}

// ActionDownload downloads a file by its ID and saves it to the specified path
func ActionDownload(config *Config) {
	savePath := strings.TrimSpace(*DownloadPath)
	if savePath == "" {
		PrintError("Save path is required")
		return
	} else if savePath == "." {
		Print("Save path is set to current directory. You can change it by -act.download.path flag")
	}

	// @todo streaming download for big files
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	name, _, err := pkg.DownloadFile(config.Token, *Download, writer, &pkg.CryptoInfo{Password: *Passwd})
	if err != nil {
		PrintError(err.Error())
		return
	}

	pathInfo, err := os.Stat(savePath)
	if err == nil && pathInfo.IsDir() {
		savePath = savePath + string(os.PathSeparator) + name
	}

	out, err := os.Create(savePath)
	if err != nil {
		PrintError("Failed to create file %s", savePath)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, &buffer)
	if err != nil {
		PrintError("Failed to save file %s", savePath)
	}
}

// ActionUpload uploads a file to the cloud. The file can be provided by path or by stdin.
func ActionUpload(config *Config, isStdIn bool) {
	var reader io.Reader
	var name string

	if isStdIn {
		name = *UploadName
		if name == "" {
			PrintError("File name is required for stdin upload. Use -act.upload.name flag")
			return
		}
		reader = os.Stdin
	} else {
		path := *Upload
		if path == "" {
			path = pkg.ScanOrDefault("Enter file path: ", "")
			if path == "" {
				PrintError("File path is required")
				return
			}
		}

		fileInfo, err := os.Stat(path)
		if err != nil {
			PrintError("Failed to access file")
			return
		}
		if fileInfo.IsDir() {
			// @todo directory uploading
			PrintError("Directory uploading is not supported yet")
			return
		}

		file, err := os.Open(path)
		if err != nil {
			PrintError("Failed to open file")
			return
		}

		if *UploadName != "" {
			name = *UploadName
		} else {
			name = file.Name()
		}

		reader = file
	}

	_, err := pkg.UploadFile(config.Token, name, "", *UploadDisk, *UploadFolder, &pkg.CryptoInfo{Password: *Passwd}, reader)
	if err != nil {
		PrintError(err.Error())
		return
	}
}

func ActionFilesList(config *Config) {
	// @todo offsets for big lists
	filesList, err := pkg.ApiRequest(config.Token, "files.get", map[string]interface{}{"disk": *FilesList, "offset": 0})
	if err != nil {
		PrintError(err.Error())
		return
	}
	if filesList.Error.Code != 0 {
		PrintError(filesList.Error.Message)
		return
	}

	// At the moment we cast the list to the interface{} and then to the []interface{} to avoid the type assertion
	// In the future, we will create a struct for the response and use it directly
	rawList, ok := filesList.Result["list"]
	if !ok {
		PrintError("Bad file get response")
		return
	}

	list, ok := rawList.([]interface{})
	if !ok {
		PrintError("Files list parameter is not a list itself")
		return
	}
	if len(list) == 0 {
		PrintError("File list is empty")
		return
	}

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("ID", "Name", "Type", "Size")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, fileInfo := range list {
		fileMap := fileInfo.(map[string]interface{})
		tbl.AddRow(fileMap["id"], fileMap["name"], fileMap["type_desc"], ByteCount(int64(fileMap["size"].(float64))))
	}

	tbl.Print()
}

func ActionApiCall(config *Config) {
	paramsMap := ParseKeyValues(*Params)
	resp, err := pkg.ApiRequest(config.Token, *Method, paramsMap)
	err = GetActualError(resp, err)
	if err != nil {
		PrintError(err.Error())
		return
	}

	Print(JsonToString(resp.Result, *Pretty))
}

// ActionAskForToken asks the user to enter the access token. The token is not displayed on the screen.
func ActionAskForToken(config *Config) {
	if config.Token != "" && *NotInteractive {
		return
	}

	// @todo prompt for email and password to get the token or use web auth
	Print("Enter your access token to use most functions or leave it blank to proceed with anonymous requests." +
		"\n When you enter your password, the characters will not be displayed." +
		"\n This is a security measure to prevent it from being stored in SSH logs.\n")
	fmt.Print("Access token: ")
	password, err := terminal.ReadPassword(0)
	if err == nil && len(password) > 0 {
		if CheckTokenAndAssign(string(password), config) != nil {
			config.Token = string(password)
		}
	} else {
		PrintError(err.Error())
	}
}
