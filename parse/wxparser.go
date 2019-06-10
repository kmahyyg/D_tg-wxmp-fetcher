package parse

import (
	"io"
)

type WxArticle struct {
	// Identifier
	AccountID    int32
	MessageID    int32
	ArticleIndex int32
	Signature    string
	// Article
	Title           string
	Date            string
	AuthorName      string
	AccountName     string
	AccountImageURL string
	ArticleImageURL string
	Brief           string
	ContentHTML     string
}

func Consume(stream io.Reader) (*WxArticle, error) {
	return &WxArticle{}, nil
}
