package types

type Floor struct {
	//Pdid    string `json:"pdid"`  // Parent Discussion ID 标识所属于的帖子
	Floor   int64  `json:"floor,omitempty"` // 所属楼层，似乎没啥必要
	Author  int64  `json:"author"`
	Content string `json:"content"`
}
