package repository

import (
	"context"
	"fmt"
)

type User struct {
	Id             int    `json:"id" db:"id"`
	Login          string `json:"login" db:"login"`
	HashedPassword string `json:"hashed_password" db:"hashed_password"`
	Name           string `json:"name" db:"name"`
	Surname        string `json:"surname" db:"surname"`
}

func (r *Repository) AddNewUser(ctx context.Context, name, surname, login, hashedPassword string) (err error) {
	_, err = r.pool.Exec(ctx, `insert into users (name, surname, login, hashed_password) values ($1, $2,$3, $4)`, name, surname, login, hashedPassword)
	if err != nil {
		err = fmt.Errorf("failed to exec data: %w", err)
		return
	}
	return
}

func (r *Repository) DeleteUser(ctx context.Context, id int) (err error) {
	_, err = r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		err = fmt.Errorf("failed to exec data: %w", err)
		return
	}
	return
}

func (r *Repository) FindUser(ctx context.Context, login string) (u User, err error) {
	row := r.pool.QueryRow(ctx, `select id, login, name, surname from users where login = $1`, login)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}
	err = row.Scan(&u.Id, &u.Login, &u.Name, &u.Surname)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}
	return
}

func (r *Repository) UpdateUser(ctx context.Context, user User) (err error) {
	_, err = r.pool.Exec(ctx,
		`UPDATE users SET login = '$2', name = '$3', surname = '$4' WHERE id = $1`,
		user.Id, user.Login, user.Name, user.Surname)
	if err != nil {
		err = fmt.Errorf("failed to exec data: %w", err)
		return
	}
	return
}

func (r *Repository) UpdateUserPassword(ctx context.Context, id int, hashedPassword string) (err error) {
	_, err = r.pool.Exec(ctx,
		`UPDATE users SET hashed_password = '$2' WHERE id = $1`,
		id, hashedPassword)
	if err != nil {
		err = fmt.Errorf("failed to exec data: %w", err)
		return
	}
	return
}

func (r *Repository) GetUserAll(ctx context.Context) (users []User, err error) {
	rows, err := r.pool.Query(ctx, `SELECT * FROM users`)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var u User
		if err = rows.Scan(&u.Id, &u.Login, &u.Name, &u.Surname); err != nil {
			err = fmt.Errorf("failed to query data: %w", err)
			return
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}
	return
}

func (r *Repository) Login(ctx context.Context, login, hashedPassword string) (u User, err error) {
	row := r.pool.QueryRow(ctx, `select id, login, name, surname from users where login = $1 AND hashed_password = $2`, login, hashedPassword)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}
	err = row.Scan(&u.Id, &u.Login, &u.Name, &u.Surname)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}
	return
}
