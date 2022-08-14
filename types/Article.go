package types

type Article struct {
	Aid            int64            `json:"aid,omitempty"` //Article ID 文章ID, 作为主键
	ArticleType    string           `json:"article_type"`  // announcement进公告区，galgame进galgame讨论区，discussion进综合讨论区
	ArticleContent []ArticleContent `json:"article_content,omitempty"`
	ArticleTags    []string         `json:"article_tags"`
	RelatedGalGame *GalGame         `json:"related_galgame,omitempty"`
}

type ArticleContent struct {
	Acid           int64  `json:"acid,omitempty"`
	ContentFloor   int64  `json:"content_floor"` // 内容的楼数，1楼即正文
	ArticleContent string `json:"article_content"`
	Author         User   `json:"author"` // 后端处理时按照Uid实时生成
	Likes          int64  `json:"likes"`  // 点赞数，之后可能会用到
}
