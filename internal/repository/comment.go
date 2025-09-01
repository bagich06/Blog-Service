package repository

import (
	"blog/internal/models"
	"context"
)

func (repo *PGRepo) CreateComment(content string, userID int, postID int) (commentID int, err error) {
	err = repo.pool.QueryRow(context.Background(), `INSERT INTO comments(content, user_id, post_id) VALUES($1, $2, $3) RETURNING id`, content, userID, postID).Scan(&commentID)
	if err != nil {
		return 0, err
	}
	return commentID, nil
}

func (repo *PGRepo) GetAllComments(postID int) (comments []models.Comment, err error) {
	rows, err := repo.pool.Query(context.Background(), `SELECT * FROM comments WHERE post_id = $1`, postID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.Content, &comment.UserId, &comment.PostId)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (repo *PGRepo) DeleteComment(commentID int) error {
	_, err := repo.pool.Exec(context.Background(), `DELETE FROM comments WHERE id = $1`, commentID)
	return err
}

func (repo *PGRepo) GetCommentByID(commentID int) (comment models.Comment, err error) {
	err = repo.pool.QueryRow(context.Background(), `SELECT id, content, user_id, post_id FROM comments WHERE id = $1`, commentID).Scan(&comment.ID, &comment.Content, &comment.UserId, &comment.PostId)
	return comment, err
}

func (repo *PGRepo) GetCommentsByUserID(userID int) (comments []models.Comment, err error) {
	rows, err := repo.pool.Query(context.Background(), `SELECT * FROM comments WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.Content, &comment.UserId, &comment.PostId)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (repo *PGRepo) GetCommentIDByUserID(userID int) (int, error) {
	var commentID int
	err := repo.pool.QueryRow(context.Background(), `SELECT comment_id FROM comments WHERE user_id = $1 LIMIT 1`, userID).Scan(&commentID)
	if err != nil {
		return 0, err
	}
	return commentID, nil
}
