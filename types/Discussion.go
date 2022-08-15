package types

type Discussion struct {
	Topic      string `json:"topic"`         // 所属主题
	Did        string `json:"did,omitempty"` // Discussion ID 由雪花算法生成，是一个帖子的全局唯一ID
	Type       string `json:"type"`          // release/discussion
	Title      string `json:"title"`
	Creator    int64  `json:"creator"`
	FirstFloor *Floor `json:"first_floor,omitempty"` // 首层
}
