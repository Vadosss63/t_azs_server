package trbl_button

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type TrblButtonRepository interface {
	Add(ctx context.Context, id_azs int) (err error)
	Update(ctx context.Context, id_azs, download int) (err error)
	Delete(ctx context.Context, id_azs int) (err error)
	Get(ctx context.Context, id_azs int) (LogButton LogButton, err error)
	GetAll(ctx context.Context) (LogButtons []LogButton, err error)
	CreateTable(ctx context.Context) (err error)
	DeleteTable(ctx context.Context) (err error)
}

type TrblButtonRepo struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *TrblButtonRepo {
	return &TrblButtonRepo{pool: pool}
}

type LogButton struct {
	IdAzs    int `json:"id_azs" db:"id_azs"`
	Download int `json:"download" db:"download"`
}

func (r *TrblButtonRepo) CreateTable(ctx context.Context) (err error) {
	_, err = r.pool.Query(ctx,
		"create table if not exists log_button"+
			"(id_azs  bigint,"+
			"download  int);")
	return
}

func (r *TrblButtonRepo) DeleteTable(ctx context.Context) (err error) {
	_, err = r.pool.Exec(ctx, "DROP TABLE log_button")
	return
}

func (r *TrblButtonRepo) Add(ctx context.Context, id_azs int) (err error) {
	_, err = r.pool.Exec(ctx, `insert into log_button (id_azs, download) values ($1, 0)`, id_azs)
	return
}

func (r *TrblButtonRepo) Update(ctx context.Context, id_azs, download int) (err error) {
	_, err = r.pool.Exec(ctx, `UPDATE log_button SET download = $2 WHERE id_azs = $1`, id_azs, download)
	return
}

func (r *TrblButtonRepo) Delete(ctx context.Context, id_azs int) (err error) {
	_, err = r.pool.Exec(ctx, `DELETE FROM log_button WHERE id_azs = $1`, id_azs)
	return
}

func (r *TrblButtonRepo) Get(ctx context.Context, id_azs int) (LogButton LogButton, err error) {
	row := r.pool.QueryRow(ctx, `SELECT * FROM log_button where id_azs = $1`, id_azs)
	if err != nil {
		return
	}

	err = row.Scan(&LogButton.IdAzs, &LogButton.Download)
	return
}

func (r *TrblButtonRepo) GetAll(ctx context.Context) (LogButtons []LogButton, err error) {
	rows, err := r.pool.Query(ctx, `SELECT * FROM log_button`)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var LogButton LogButton
		if err = rows.Scan(&LogButton.IdAzs, &LogButton.Download); err != nil {
			return
		}
		LogButtons = append(LogButtons, LogButton)
	}
	return
}
