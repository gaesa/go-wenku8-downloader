package downloader

import (
	"os"
	"path"
	"time"
)

var Root = path.Join(os.TempDir(), "wenku8-download")

const (
	ErrorJsonName   = "error.json"
	ImageFolderName = "images"
	DownloadTimer   = time.Second
	RetryTimes      = 6
	RetryTimer      = 6 * time.Second
)
