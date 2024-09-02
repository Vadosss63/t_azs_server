package user

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	Id             int    `json:"id" db:"id"`
	Login          string `json:"login" db:"login"`
	HashedPassword string `json:"hashed_password" db:"hashed_password"`
	Name           string `json:"name" db:"name"`
	Surname        string `json:"surname" db:"surname"`
}

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Add(ctx context.Context, name, surname, login, hashedPassword string) (err error) {
	_, err = r.pool.Exec(ctx, `insert into users (name, surname, login, hashed_password) values ($1, $2,$3, $4)`, name, surname, login, hashedPassword)
	return
}

func (r *UserRepo) Delete(ctx context.Context, id int) (err error) {
	_, err = r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	return
}

func (r *UserRepo) Get(ctx context.Context, id int) (u User, err error) {
	row := r.pool.QueryRow(ctx, `select id, login, name, surname from users where id = $1`, id)

	if err != nil {
		return
	}
	err = row.Scan(&u.Id, &u.Login, &u.Name, &u.Surname)
	return
}

func (r *UserRepo) Find(ctx context.Context, login string) (u User, err error) {
	row := r.pool.QueryRow(ctx, `select id, login, name, surname from users where login = $1`, login)
	if err != nil {
		return
	}
	err = row.Scan(&u.Id, &u.Login, &u.Name, &u.Surname)
	return
}

func (r *UserRepo) Update(ctx context.Context, user User) (err error) {
	_, err = r.pool.Exec(ctx,
		`UPDATE users SET login = '$2', name = '$3', surname = '$4' WHERE id = $1`,
		user.Id, user.Login, user.Name, user.Surname)
	return
}

func (r *UserRepo) UpdateUserPassword(ctx context.Context, id int, hashedPassword string) (err error) {
	_, err = r.pool.Exec(ctx,
		`UPDATE users SET hashed_password = $2 WHERE id = $1`,
		id, hashedPassword)
	return
}

func (r *UserRepo) GetAll(ctx context.Context) (users []User, err error) {
	rows, err := r.pool.Query(ctx, `SELECT * FROM users`)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var u User
		if err = rows.Scan(&u.Id, &u.Login, &u.HashedPassword, &u.Name, &u.Surname); err != nil {
			err = fmt.Errorf("failed to query data: %w", err)
			return
		}
		u.HashedPassword = ""
		users = append(users, u)
	}
	return
}

func (r *UserRepo) Login(ctx context.Context, login, hashedPassword string) (u User, err error) {
	row := r.pool.QueryRow(ctx, `select id, login, name, surname from users where login = $1 AND hashed_password = $2`, login, hashedPassword)
	if err != nil {
		return
	}
	err = row.Scan(&u.Id, &u.Login, &u.Name, &u.Surname)
	return
}
