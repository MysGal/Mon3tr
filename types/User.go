package types

type User struct {
	Uid       int64  `json:"uid"`
	UserName  string `json:"user_name,omitempty"`
	Password  string `json:"password,omitempty"`
	Email     string `json:"email,omitempty"`
	UserGroup string `json:"user_group,omitempty"`
}
