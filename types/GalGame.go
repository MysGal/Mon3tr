package types

type GalGame struct {
	Gid       int      `json:"gid"` // gid即Game ID,对应每一个Gal 与数据库上自增主键相对应
	Name      string   `json:"name"`
	Author    []string `json:"author"`
	Publisher []string `json:"publisher"`
	Tags      []string `json:"tags"`
	//detail    string   `json:"detail"` // 包含Markdown原文本，应以base64或其他编码传输
}
