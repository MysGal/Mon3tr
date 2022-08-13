package types

type User struct {
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
	Uid       int    `json:"uid"`
	Email     string `json:"email"`
	UserGroup string `json:"user_group"`
}
