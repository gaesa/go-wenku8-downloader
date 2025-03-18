package downloader

import (
	"encoding/json"
	"os"
	"path"

	"github.com/gaesa/go-wenku-downloader/scraper"
	"github.com/gaesa/go-wenku-downloader/util"
)

func DownloadNovelMetadata(novel *scraper.Novel, dirPath string) error {
	filePath := path.Join(dirPath, "metadata.json")

	if err := util.CheckDir(dirPath); err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(novel); err != nil {
		return err
	}

	return nil
}

func DownloadVolumeMetadata(volume *scraper.Volume, dirPath string) error {
	filePath := path.Join(dirPath, "metadata.json")

	if err := util.CheckDir(dirPath); err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(volume); err != nil {
		return err
	}

	return nil
}
