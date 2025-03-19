package prompt

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/adrg/xdg"
	"github.com/gaesa/go-wenku-downloader/downloader"
	"github.com/gaesa/go-wenku-downloader/scraper"
	"github.com/gaesa/go-wenku-downloader/util"
	"github.com/go-shiori/go-epub"
)

func download(novelId int) error {
	downloadPath := path.Join(downloader.Root, strconv.Itoa(novelId))
	if err := util.CheckDir(downloadPath); err != nil {
		return fmt.Errorf("创建目录失败: %e", err)
	}

	// get novel metadata
	novel, err := promptNovelDetails(int(novelId))
	if err != nil {
		return fmt.Errorf("获取小说信息失败: %e", err)
	}

	// download novel metadata and coverImg
	downloader.DownloadNovelMetadata(novel, downloadPath)

	// download cover image
	success := false
	for range [3]int{} {
		err = downloader.DownloadImage(novel.Cover, downloadPath)
		if err == nil {
			success = true
			break
		}
	}
	if !success {
		return fmt.Errorf("download cover image failed %v", err)
	}

	// get selected volume list
	volumeArray, err := promptVolumeSelect(novel.CatalogueUrl)
	if err != nil {
		return fmt.Errorf("下载小说卷信息失败: %e", err)
	}

	// get coverIndex from input
	converIndex, err := inputCoverIndex()
	if err != nil {
		return fmt.Errorf("prompt failed %v", err)
	}

	// get onlyWenku8Img
	onlyWenku8Img, err := getInputBool("是否只下载wenku8的图片(推荐使用默认数值, 非文库图片资源大多数情况已失效), 默认:y(y/n)", true)
	if err != nil {
		return fmt.Errorf("prompt failed %v", err)
	}

	// download volume
	for _, volume := range volumeArray {
		volumePath := path.Join(downloadPath, formatFilename(volume.Name))
		err = downloader.DownloadVolume(volume, volumePath, onlyWenku8Img)
		if err != nil {
			log.Printf("download volume error %v", err)
			continue
		}
		err = createEpub(novel, volume.Name, volume.ChapterCount, converIndex, downloadPath)
		if err != nil {
			log.Printf("create epub failed: %v", err)
			continue
		}
	}

	return nil
}

func createEpub(novel *scraper.Novel, volumeName string, chapterCount int, coverIndex int, downloadPath string) error {
	formatedVolumeName := formatFilename(volumeName)
	formatedNovelName := formatFilename(novel.NovelName + "231231")
	var imagePathList []string
	// output epub path
	epubFilePath := path.Join(xdg.UserDirs.Download, fmt.Sprintf("%s %s.epub", formatedNovelName, formatedVolumeName))

	if util.CheckFileExist(epubFilePath) {
		fmt.Println()
		fmt.Printf("文件 %s 已经存在。是否覆盖？(y/N)：", epubFilePath)
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			return nil
		}
	}

	// volume path
	volumePath := path.Join(downloadPath, formatedVolumeName)
	imagePath := path.Join(volumePath, downloader.ImageFolderName)
	coverPath := path.Join(downloadPath, util.GetUrlLastString(novel.Cover))

	// create epub
	epub, err := epub.NewEpub(novel.NovelName + " " + volumeName)
	if err != nil {
		return err
	}
	epub.SetAuthor(novel.Author)
	epub.SetLang("zh-CN")

	// add coverImage to epub
	if util.CheckFileExist(coverPath) {
		_, err := util.AddImage(epub, coverPath)
		if err != nil {
			return fmt.Errorf("add image to epub failed")
		}
		imagePathList = append(imagePathList, coverPath)
	}

	for i := 1; i <= chapterCount; i++ {
		file, err := os.ReadFile(path.Join(volumePath, fmt.Sprintf("%d.json", i)))
		if err != nil {
			return err
		}
		chapter := &scraper.Chapter{}
		err = json.Unmarshal(file, &chapter)
		if err != nil {
			return err
		}
		xhtml := util.CreateSectionXhtml(chapter.Title, chapter.Content.Article)
		if len(chapter.Content.Images) != 0 {
			for _, img := range chapter.Content.Images {
				imgFile := path.Join(imagePath, util.GetUrlLastString(img))
				if !util.CheckFileExist(imgFile) {
					continue
				}
				internalPath, _ := util.AddImage(epub, imgFile)
				xhtml = util.AddImageToXhtml(internalPath, xhtml)
				imagePathList = append(imagePathList, imgFile)
			}
		}
		err = util.AddSectionXhtml(epub, chapter.Title, xhtml)
		if err != nil {
			return err
		}
	}
	tempCoverPath := imagePathList[0]
	if coverIndex < len(imagePathList) {
		tempCoverPath = imagePathList[coverIndex]
	}
	internalCoverPath, _ := util.AddCover(epub, tempCoverPath)
	epub.SetCover(internalCoverPath, "")

	err = epub.Write(epubFilePath)
	if err != nil {
		return err
	}
	return nil
}

func formatFilename(str string) string {
	newFilename := strings.ReplaceAll(str, "/", "-")
	re := regexp.MustCompile(`\p{P}|[0-9|=]`)
	newFilename = re.ReplaceAllString(newFilename, "")
	return newFilename
}
