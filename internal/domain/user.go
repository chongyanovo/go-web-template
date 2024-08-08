package domain

type User struct {
	ID       int64  `json:"id"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}
