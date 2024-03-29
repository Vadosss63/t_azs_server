package repository

import (
	"context"
)

type UpdateCommand struct {
	IdAzs int    `json:"id_azs" db:"id_azs"`
	Url   string `json:"url" db:"url"`
}

func (r *Repository) CreateUpdateCommandTable(ctx context.Context) (err error) {
	_, err = r.pool.Query(ctx,
		"create table if not exists update_command"+
			"(id_azs  bigint,"+
			"url varchar(100) not null);")
	return
}

func (r *Repository) DeleteUpdateCommandTable(ctx context.Context) (err error) {
	_, err = r.pool.Exec(ctx, "DROP TABLE update_command")
	return
}

func (r *Repository) AddUpdateCommand(ctx context.Context, id_azs int) (err error) {
	_, err = r.pool.Exec(ctx, `insert into update_command (id_azs, url) values ($1, '')`, id_azs)
	return
}

func (r *Repository) UpdateUpdateCommand(ctx context.Context, id_azs int, url string) (err error) {
	_, err = r.pool.Exec(ctx, `UPDATE update_command SET url = $2 WHERE id_azs = $1`, id_azs, url)
	return
}

func (r *Repository) DeleteUpdateCommand(ctx context.Context, id_azs int) (err error) {
	_, err = r.pool.Exec(ctx, `DELETE FROM update_command WHERE id_azs = $1`, id_azs)
	return
}

func (r *Repository) GetUpdateCommand(ctx context.Context, id_azs int) (UpdateCommand UpdateCommand, err error) {
	row := r.pool.QueryRow(ctx, `SELECT * FROM update_command where id_azs = $1`, id_azs)
	if err != nil {
		return
	}

	err = row.Scan(&UpdateCommand.IdAzs, &UpdateCommand.Url)
	return
}

func (r *Repository) GetUpdateCommandAll(ctx context.Context) (UpdateCommands []UpdateCommand, err error) {
	rows, err := r.pool.Query(ctx, `SELECT * FROM update_command`)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var UpdateCommand UpdateCommand
		if err = rows.Scan(&UpdateCommand.IdAzs, &UpdateCommand.Url); err != nil {
			return
		}
		UpdateCommands = append(UpdateCommands, UpdateCommand)
	}
	return
}
