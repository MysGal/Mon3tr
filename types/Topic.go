package types

type Topic struct {
	Topic string   `json:"topic"` // Url中主题标识符，给前端进行route使用，此即为主题的ID
	Name  string   `json:"name"`
	Tags  []string `json:"tags"`
}
