package pkg

// @todo auto-generate this

// ApiResponse is the basic mapstructure-RPC response structure
type ApiResponse struct {
	ID    uint `mapstructure:"id"`
	Error struct {
		Code    uint   `mapstructure:"code"`
		Message string `mapstructure:"message"`
	} `mapstructure:"error,omitempty"`
	Result map[string]interface{} `mapstructure:"result,omitempty"`
}

type File struct {
	Date       int    `mapstructure:"date"`
	Disk       string `mapstructure:"disk"`
	Encrypted  bool   `mapstructure:"encrypted"`
	Folder     string `mapstructure:"folder"`
	ID         string `mapstructure:"id"`
	Mime       string `mapstructure:"mime"`
	Name       string `mapstructure:"name"`
	NameCrypto string `mapstructure:"name_crypto"`
	Rating     int    `mapstructure:"rating"`
	Size       int    `mapstructure:"size"`
	Text       string `mapstructure:"text"`
	Type       string `mapstructure:"type"`
	TypeDesc   string `mapstructure:"type_desc"`
	URLSecret  string `mapstructure:"url_secret"`
	URLShared  bool   `mapstructure:"url_shared"`
}

type Folder struct {
	Disk   string `mapstructure:"disk"`
	ID     string `mapstructure:"id"`
	Name   string `mapstructure:"name"`
	Parent string `mapstructure:"parent"`
}

type FilesGetResponse struct {
	Folders    []*Folder `mapstructure:"folders"`
	HasFiles   bool      `mapstructure:"has_files"`
	HasFolders bool      `mapstructure:"has_folders"`
	List       []*File   `mapstructure:"list"`
	Offset     int       `mapstructure:"offset"`
}

type UserInfo struct {
	AvailableSpace int64  `mapstructure:"available_space"`
	Avatar         string `mapstructure:"avatar"`
	Email          string `mapstructure:"email"`
	ID             string `mapstructure:"id"`
	PrepaidSpace   int64  `mapstructure:"prepaid_space"`
}

type FileGetByIdResponse struct {
	Count int     `mapstructure:"count"`
	List  []*File `mapstructure:"list"`
}

type DownloadResponse struct {
	Crypto bool   `mapstructure:"crypto"`
	Name   string `mapstructure:"name"`
	URL    string `mapstructure:"url"`
}

type UploadResult struct {
	FileID string `mapstructure:"file_id"`
	Ok     bool   `mapstructure:"ok"`
}

type Disk struct {
	CryptoKey string `mapstructure:"crypto_key"`
	ID        string `mapstructure:"id"`
	PublicKey string `mapstructure:"public_key"`
	Title     string `mapstructure:"title"`
}

type DisksInfo struct {
	Count int     `mapstructure:"count"`
	List  []*Disk `mapstructure:"list"`
}
