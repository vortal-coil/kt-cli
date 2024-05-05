package internal

import (
	"flag"
	"os"
)

var (
	// Debug - Do not catch panics and prettify them, show debug info
	Debug = flag.Bool("Debug", false, "Enable Debug mode")
	// Params - Set parameters for API method if called
	Params = flag.String("params", "", "Set API method key=value parameters separated by space (format: k=v k=v k=v...)")

	// Global flags

	ConfigFilename = flag.String("config", "config.yaml", "Set config file path")
	PrintModeFlag  = flag.Int("output", ModeLog, "Output mode (0 - log with timestamp, 1 - plain log, 2 - no newline)")
	NotInteractive = flag.Bool("no-interactive", false, "Do not ask for any input, use default values")
	NoConfigSave   = flag.Bool("no-save", false, "Do not save the config file on exit (including token)")
	Auth           = flag.String("token", "", "Set auth token for future requests (will be saved in config file; also you can use environment variable KT_CLI_TOKEN)")
	Pretty         = flag.Bool("pretty", false, "Pretty-print JSON responses")
	Passwd         = flag.String("passwd", "", "Set password for encryption/decryption. Also you can use environment variable KT_CLI_PASSWD")
	PublicKeyFile  = flag.String("public", "public_key.pub", "Set public key file path for encryption/decryption (will be downloaded from the server if empty)")
	PrivateKeyFile = flag.String("private", "private_key.asc", "Set private key file path for encryption/decryption (will be downloaded and decrypted from the server if empty)")

	// Actions to perform

	Method = flag.String("act.method", "", "Call API method")
	Ping   = flag.Bool("act.ping", false, "Check if API is alive")

	GetKeys            = flag.String("act.keys", "", "Download keys for the provided disk (\".\" for default disk)")
	GetKeysPublicName  = flag.String("act.keys.public", "public_key.pub", "Set public key name for download")
	GetKeysPrivateName = flag.String("act.keys.private", "private_key.asc", "Set private key name for download")

	Download     = flag.String("act.download", "", "Download file by file ID")
	DownloadPath = flag.String("act.download.path", ".", "Set path to save downloaded file")

	Upload       = flag.String("act.upload", "", "Upload file by path; stdin is also supported")
	UploadName   = flag.String("act.upload.name", "", "Set file name for upload (required for stdin)")
	UploadDisk   = flag.String("act.upload.disk", "", "Set disk for upload")
	UploadFolder = flag.String("act.upload.folder", "", "Set folder for upload")

	FilesList = flag.String("act.files", "", "List files in provided disk")
	// @todo method to replace files contents
)

// ScanEnv scans environment variables as replacement for the flags that are not set
func ScanEnv() {
	if *Auth == "" {
		*Auth = os.Getenv("KT_CLI_TOKEN")
	}
	if *Passwd == "" {
		*Passwd = os.Getenv("KT_CLI_PASSWD")
	}
}
