package repository

import (
	"blog/internal/models"
	"context"
)

func (repo *PGRepo) CreatePost(post models.Post, userID int) (id int, err error) {
	err = repo.pool.QueryRow(context.Background(), `INSERT INTO posts(content, user_id) VALUES($1, $2) RETURNING id`, post.Content, userID).Scan(&post.ID)
	if err != nil {
		return 0, err
	}
	return post.ID, nil
}

func (repo *PGRepo) GetAllPosts(userID int) (posts []models.Post, err error) {
	rows, err := repo.pool.Query(context.Background(), `SELECT id, content, user_id FROM posts WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Content, &post.UserId)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (repo *PGRepo) GetPostByID(postID int, userID int) (post models.Post, err error) {
	err = repo.pool.QueryRow(context.Background(), `SELECT id, content, user_id FROM posts WHERE id = $1 and user_id=$2`, postID, userID).Scan(&post.ID, &post.Content, &post.UserId)
	return post, err
}

func (repo *PGRepo) DeletePost(postID int, userID int) error {
	_, err := repo.pool.Exec(context.Background(), `DELETE FROM posts WHERE id = $1 AND user_id = $2`, postID, userID)
	return err
}
