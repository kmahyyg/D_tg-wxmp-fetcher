package article

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

// Metadata stores all necessary information to generate a link preview
type Metadata struct {
	Link   string
	Title  string
	Author string
	Image  string
	Brief  string
}
