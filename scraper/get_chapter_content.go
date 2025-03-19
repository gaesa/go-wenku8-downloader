package scraper

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getChapterContentAndroid(chapter *Chapter) error {
	doc, err := AndroidGet(chapter.Url)
	if err != nil {
		return err
	}

	content := doc.Find("body").Text()

	content = strings.ReplaceAll(content, "&nbsp;", "")
	content = strings.ReplaceAll(content, "更多精彩热门日本轻小说、动漫小说，轻小说文库(http://www.wenku8.com) 为你一网打尽！", "")

	picReg := regexp.MustCompile(`http:\/\/pic\.wenku8\.com\/pictures\/[\/0-9]+.jpg`)
	picRegL := regexp.MustCompile(`http:\/\/pic\.wenku8\.com\/pictures\/[\/0-9]+.jpg\([0-9]+K\)`)
	images := picReg.FindAllString(content, -1)
	content = picRegL.ReplaceAllString(content, "")
	content = picReg.ReplaceAllString(content, "")

	chapterContent := &ChapterContent{
		Images:  images,
		Article: content,
	}
	chapter.Content = chapterContent
	return nil
}

func isContentFetched(doc *goquery.Document) bool {
	return strings.TrimSpace(doc.Find("#contentmain span").First().Text()) != "null"
}

func GetChapterContent(chapter *Chapter) error {
	doc, err := Get(chapter.Url)
	if err != nil {
		return err
	}
	if !isContentFetched(doc) {
		return getChapterContentAndroid(chapter)
	}

	content := doc.Find("#content").Text()
	content = strings.ReplaceAll(content, "本文来自 轻小说文库(http://www.wenku8.com)", "")
	content = strings.ReplaceAll(content, "台版 转自 轻之国度", "")
	content = strings.ReplaceAll(content, "最新最全的日本动漫轻小说 轻小说文库(http://www.wenku8.com) 为你一网打尽！", "")

	images := []string{}
	doc.Find("img").Each(func(i int, imgEle *goquery.Selection) {
		src, _ := imgEle.Attr("src")
		images = append(images, src)
	})

	chapterContent := &ChapterContent{
		Images:  images,
		Article: content,
	}
	chapter.Content = chapterContent
	return nil
}
