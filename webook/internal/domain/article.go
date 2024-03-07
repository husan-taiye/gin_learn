package domain

import "time"

type Article struct {
	Id      int64
	Title   string
	Content string
	Author  Author
	Status  ArticleStatus
	Ctime   time.Time
	Utime   time.Time
}

func (a Article) Abstract() string {
	// 摘要取前几句
	// 考虑中文问题
	cs := []rune(a.Content)
	if len(cs) < 100 {
		return a.Content
	}
	// 英文？ 完整单词？ 别纠结！
	return string(cs[:100])
}

type ArticleStatus int8

const (
	// ArticleStatusUnknown 为了避免零值之类的问题
	ArticleStatusUnknown ArticleStatus = iota
	ArticleStatusUnpublished
	ArticleStatusPublished
	ArticleStatusPrivate
)

func (arts ArticleStatus) ToUint8() uint8 {
	return uint8(arts)
}
func (arts ArticleStatus) Valid() bool {
	return arts.ToUint8() > 0
}

func (arts ArticleStatus) NonPublished() bool {
	return arts != ArticleStatusPublished
}

func (arts ArticleStatus) String() string {
	switch arts {
	case ArticleStatusPrivate:
		return "private"
	case ArticleStatusUnpublished:
		return "unpublished"
	case ArticleStatusPublished:
		return "published"
	case ArticleStatusUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}

// ArticleStatusV1 如果状态很复杂，有很多行为，状态里面需要一些额外的字段
type ArticleStatusV1 struct {
	Val  uint8
	Name string
}

var (
	ArticleStatusV1Unknown = ArticleStatusV1{Val: 0, Name: "unknown"}
)

type Author struct {
	Id   int64
	Name string
}
