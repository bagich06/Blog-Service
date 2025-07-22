package models

type Comment struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	UserId  int    `json:"user_id"`
	PostId  int    `json:"post_id"`
}
