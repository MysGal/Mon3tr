package types

type Article struct {
	Aid            int     `json:"aid"`
	ArticleType    string  `json:"article_type"`
	ArticleContent string  `json:"article_content"`
	Author         User    `json:"author"`
	RelatedGalGame GalGame `json:"related_galgame,omitempty"`
}
