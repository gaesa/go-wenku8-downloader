package scraper

import (
	"log"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func GetChapterArray(volume *Volume) ([]*Chapter, error) {
	doc, err := Get(volume.CatalogueUrl)
	if err != nil {
		return nil, err
	}
	chapterArray := make([]*Chapter, 0)
	rows := doc.Find("tbody").Children()
	insertMap(rows, &chapterArray, volume.RowNumber, volume.EndRow, volume.Name, volume.CatalogueUrl)
	return chapterArray, nil
}

func insertMap(rows *goquery.Selection, chapterArray *[]*Chapter, start int, end int, volumeName string, catalogueUrl string) {
	baseUrl, err := url.Parse(catalogueUrl)
	if err != nil {
		log.Printf("Error: Invalid catalogue URL %s: %v", catalogueUrl, err)
		return
	}

	rows.Slice(start, end).Find("a").Each(func(i int, s *goquery.Selection) {
		chapterIndex := i + 1
		chapterTitle := s.Text()

		chapterUrl, exists := s.Attr("href")
		if !exists {
			log.Printf("Warning: Missing chapter URL in volume %s, skipping entry %d", volumeName, chapterIndex)
			return
		}

		parsedChaperUrl, err := url.Parse(chapterUrl)
		if err != nil {
			log.Printf("Error: Invalid chapter URL %s in volume %s: %v", chapterUrl, volumeName, err)
			return
		}
		fullUrl := baseUrl.ResolveReference(parsedChaperUrl)

		chapter := &Chapter{
			Index: chapterIndex,
			Title: chapterTitle,
			Url:   fullUrl.String(),
		}
		*chapterArray = append(*chapterArray, chapter)
	})
}
