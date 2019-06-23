package article

import (
	"encoding/base64"
	"fmt"
	"strconv"

	"bitbucket.org/mutze5/wxfetcher/proto"
)

const (
	fmtWxArticleLink = "https://mp.weixin.qq.com/s?__biz=%s&mid=%d&idx=%d&sn=%s"
)

// WxArticle is a struct to store full article informations
type WxArticle struct {
	// Identifier
	AccountID  int64  `jsvar:"biz" encoding:"base64"`
	MessageID  int64  `jsvar:"mid"`
	ArticleIdx int64  `jsvar:"idx"`
	Signature  string `jsvar:"sn"`
	// Article
	Title       string `jsvar:"msg_title"`
	AccountName string `jsvar:"nickname"`
	AuthorName  string // Will be filled during HTML parse
	Brief       string `jsvar:"msg_desc"`
	Timestamp   int64  `jsvar:"ct"`
	// Image
	AccountImageURL string `jsvar:"hd_head_img"`
	ArticleImageURL string `jsvar:"msg_cdn_url"`
	ContentHTML     string // will be filled during HTML parse

	// Internal
	jsVarUnfilled map[string]int // jsvar to field id
}

// WxArticleLink generate a link to WeChat article
func WxArticleLink(accountID, messageID, articleIdx int64, signature string) string {
	encodedBiz := base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(accountID, 10)))
	return fmt.Sprintf(fmtWxArticleLink, encodedBiz, messageID, articleIdx, signature)
}

// Meta generates a ArticleMeta for current WeChat article
func (a *WxArticle) Meta() (meta *proto.ArticleMeta) {
	meta = &proto.ArticleMeta{
		Link:      WxArticleLink(a.AccountID, a.MessageID, a.ArticleIdx, a.Signature),
		Title:     a.Title,
		Timestamp: a.Timestamp,
		Image:     a.ArticleImageURL,
		Brief:     a.Brief,
	}
	if a.AuthorName != "" {
		meta.Author = fmt.Sprintf("%s | %s", a.AccountName, a.AuthorName)
	} else {
		meta.Author = a.AccountName
	}
	return
}
