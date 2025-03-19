package downloader

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/gaesa/go-wenku8-downloader/scraper"
	"github.com/gaesa/go-wenku8-downloader/util"
)

func DownloadChapter(chapter *scraper.Chapter, dirPath string) error {
	filePath := path.Join(dirPath, fmt.Sprintf("%v.json", chapter.Index))

	// folder exist if not then create
	if err := util.CheckDir(dirPath); err != nil {
		return fmt.Errorf("章节，创建目录失败 %v", err)
	}

	if strings.TrimSpace(chapter.Content.Article) == "" {
		chapter.Content.Article = ""
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(chapter); err != nil {
		return err
	}

	return nil
}
