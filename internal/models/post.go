package models

type Post struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	UserId    int    `json:"user_id"`
	CommentID int    `json:"comment_id"`
}
