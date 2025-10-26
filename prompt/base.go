package prompt

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/gaesa/go-wenku8-downloader/scraper/enums"
)

type Questions int

const (
	ViewPopularNovels Questions = iota
	SearchNovels
	DownloadNovel
	DoNothing
)

var QuestionsText = []string{
	"查看(今日更新/热门轻小说/总推荐榜/...)",
	"搜索小说",
	"下载小说",
	"什么也不做",
}

func InitPrompt() {
	for {
		selectedIndex, err := getSelectedIndex("你打算做什么", QuestionsText)
		if err != nil {
			if !errors.Is(err, terminal.InterruptErr) {
				log.Print(err)
			}
			break
		} else {
			questionTwo(Questions(selectedIndex))
		}
	}
}

func questionTwo(question Questions) {
	handleError := func(err error, s string) {
		if err != nil {
			if !errors.Is(err, terminal.InterruptErr) {
				log.Printf("%v failed %v\n", s, err)
				return
			} else {
				os.Exit(0)
			}
		}
	}

	switch question {
	case ViewPopularNovels:
		selectedIndex, err := getSelectedIndex("请选择分类", enums.TopSoftText)
		handleError(err, "Search")
		promptTopList(enums.TopSortType(selectedIndex))

	case SearchNovels:
		selectedIndex, err := getSelectedIndex("请选择搜索类型", enums.SearchTypeText)
		handleError(err, "Search")

		str, err := getInputString(fmt.Sprintf("请输入要搜索的%s", enums.SearchTypeText[selectedIndex]))
		handleError(err, "Search")

		err = searchNovels(str, enums.SearchType(selectedIndex))
		handleError(err, "Search")

	case DownloadNovel:
		novelId, err := inputNovelId()
		handleError(err, "Prompt")
		err = download(novelId)
		handleError(err, "Download")

	case DoNothing:
		os.Exit(0)
	default:
		fmt.Println()
	}
}
