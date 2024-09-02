package updater_button

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type UpdateCommand struct {
	IdAzs int    `json:"id_azs" db:"id_azs"`
	Url   string `json:"url" db:"url"`
}

type UpdaterButtonRepo struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *UpdaterButtonRepo {
	return &UpdaterButtonRepo{pool: pool}
}

func (r *UpdaterButtonRepo) CreateTable(ctx context.Context) (err error) {
	_, err = r.pool.Query(ctx,
		"create table if not exists update_command"+
			"(id_azs  bigint,"+
			"url varchar(100) not null);")
	return
}

func (r *UpdaterButtonRepo) DeleteTable(ctx context.Context) (err error) {
	_, err = r.pool.Exec(ctx, "DROP TABLE update_command")
	return
}

func (r *UpdaterButtonRepo) Add(ctx context.Context, id_azs int) (err error) {
	_, err = r.pool.Exec(ctx, `insert into update_command (id_azs, url) values ($1, '')`, id_azs)
	return
}

func (r *UpdaterButtonRepo) Update(ctx context.Context, id_azs int, url string) (err error) {
	_, err = r.pool.Exec(ctx, `UPDATE update_command SET url = $2 WHERE id_azs = $1`, id_azs, url)
	return
}

func (r *UpdaterButtonRepo) Delete(ctx context.Context, id_azs int) (err error) {
	_, err = r.pool.Exec(ctx, `DELETE FROM update_command WHERE id_azs = $1`, id_azs)
	return
}

func (r *UpdaterButtonRepo) Get(ctx context.Context, id_azs int) (UpdateCommand UpdateCommand, err error) {
	row := r.pool.QueryRow(ctx, `SELECT * FROM update_command where id_azs = $1`, id_azs)
	if err != nil {
		return
	}

	err = row.Scan(&UpdateCommand.IdAzs, &UpdateCommand.Url)
	return
}

func (r *UpdaterButtonRepo) GetAll(ctx context.Context) (UpdateCommands []UpdateCommand, err error) {
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
