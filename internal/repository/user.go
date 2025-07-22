package repository

import (
	"blog/internal/models"
	"context"
)

func (repo *PGRepo) CreateUser(user models.User) (taskID int, err error) {
	err = repo.pool.QueryRow(context.Background(), `INSERT INTO users (name, phone, password) VALUES ($1, $2, $3) RETURNING id`, user.Name, user.Phone, user.Password).Scan(&taskID)
	if err != nil {
		return 0, err
	}
	return taskID, nil
}

func (repo *PGRepo) GetUserByID(id int) (user models.User, err error) {
	err = repo.pool.QueryRow(context.Background(), `SELECT id, name, phone, password FROM users WHERE id=$1`, id).Scan(&user.ID, &user.Name, &user.Phone, &user.Password)
	return user, nil
}

func (repo *PGRepo) GetAllUsers() (users []models.User, err error) {
	row, err := repo.pool.Query(context.Background(), `SELECT id, name, phone, password FROM users`)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		var user models.User
		err = row.Scan(&user.ID, &user.Name, &user.Phone, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (repo *PGRepo) Delete(id int) error {
	_, err := repo.pool.Exec(context.Background(), `DELETE FROM users WHERE id=$1`, id)
	return err
}

func (repo *PGRepo) GetUserByPhone(phone string) (user models.User, err error) {
	err = repo.pool.QueryRow(context.Background(), `SELECT id, name, phone, password FROM users WHERE phone=$1`, phone).Scan(&user.ID, &user.Name, &user.Phone, &user.Password)
	return user, nil
}
