package repository

import (
	"context"
)

type LogButton struct {
	IdAzs    int `json:"id_azs" db:"id_azs"`
	Download int `json:"download" db:"download"`
}

func (r *Repository) CreateLogButtonTable(ctx context.Context) (err error) {
	_, err = r.pool.Query(ctx,
		"create table if not exists log_button"+
			"(id_azs  bigint,"+
			"download  int);")
	return
}

func (r *Repository) DeleteLogButtonTable(ctx context.Context) (err error) {
	_, err = r.pool.Exec(ctx, "DROP TABLE log_button")
	return
}

func (r *Repository) AddLogButton(ctx context.Context, id_azs int) (err error) {
	_, err = r.pool.Exec(ctx, `insert into log_button (id_azs, download) values ($1, 0)`, id_azs)
	return
}

func (r *Repository) UpdateLogButton(ctx context.Context, id_azs, download int) (err error) {
	_, err = r.pool.Exec(ctx, `UPDATE log_button SET download = $2 WHERE id_azs = $1`, id_azs, download)
	return
}

func (r *Repository) DeleteLogButton(ctx context.Context, id_azs int) (err error) {
	_, err = r.pool.Exec(ctx, `DELETE FROM log_button WHERE id_azs = $1`, id_azs)
	return
}

func (r *Repository) GetLogButton(ctx context.Context, id_azs int) (LogButton LogButton, err error) {
	row := r.pool.QueryRow(ctx, `SELECT * FROM log_button where id_azs = $1`, id_azs)
	if err != nil {
		return
	}

	err = row.Scan(&LogButton.IdAzs, &LogButton.Download)
	return
}

func (r *Repository) GetLogButtonAll(ctx context.Context) (LogButtons []LogButton, err error) {
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
