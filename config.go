package epub_gen

import "time"

// 书籍
type Book struct {
	Title       string
	Description string
	Publisher   string
	Author      []string
	Date        time.Time
	Lang        string
	Fonts       []string
	TocTitle    string
	Cover       string
	Sections    []*Section
	Version     int
}

// 章节
type Section struct {
	Title  string
	Author []string
	Date   string
	Link   string
	Data   string

	Filename       string
	Href           string
	ExcludeFromToc bool
	BeforeToc      bool

	id  string
	dir string
}

// 下载生成控制
type Control struct {
	Looking struct {
		AppendChapterTitles bool
	}
	HttpDL struct {
		Headers    map[string]string
		Timeout    time.Time
		RetryTimes int
	}
	Output struct {
		Path string
	}
	Debug struct {

	}
}