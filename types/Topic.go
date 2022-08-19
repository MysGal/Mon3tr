package types

type Topic struct {
	Topic       string   `json:"topic"` // Url中主题标识符，给前端进行route使用，此即为主题的ID
	Type        string   `json:"type"`
	Name        string   `json:"name"`
	Tags        []string `json:"tags"`
	RelatedData struct {
		GalGamePublisher []string `json:"galgame_publisher,omitempty"`
		GalGameAuthor    []string `json:"galgame_author,omitempty"`
	} `json:"related_data"`
}
