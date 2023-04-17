package repository

import (
	"context"
)

type AzsButton struct {
	IdAzs  int `json:"id_azs" db:"id_azs"`
	Price  int `json:"price" db:"price"`
	Button int `json:"button" db:"button"`
}

func (r *Repository) CreateAzsButtonTable(ctx context.Context) (err error) {
	_, err = r.pool.Query(ctx,
		"create table if not exists azs_button"+
			"(id_azs  bigint,"+
			"price   int,"+
			"button  int);")
	return
}

func (r *Repository) DeleteAzsButtonTable(ctx context.Context) (err error) {
	_, err = r.pool.Exec(ctx, "DROP TABLE azs_button")
	return
}

func (r *Repository) AddAzsButton(ctx context.Context, id_azs int) (err error) {
	_, err = r.pool.Exec(ctx, `insert into azs_button (id_azs, price, button) values ($1, 0, 0)`, id_azs)
	return
}

func (r *Repository) UpdateAzsButton(ctx context.Context, id_azs, price, button int) (err error) {
	_, err = r.pool.Exec(ctx, `UPDATE azs_button SET price = $2, button = $3 WHERE id_azs = $1`, id_azs, price, button)
	return
}

func (r *Repository) DeleteAzsButton(ctx context.Context, id_azs int) (err error) {
	_, err = r.pool.Exec(ctx, `DELETE FROM azs_button WHERE id_azs = $1`, id_azs)
	return
}

func (r *Repository) GetAzsButton(ctx context.Context, id_azs int) (azsButton AzsButton, err error) {
	row := r.pool.QueryRow(ctx, `SELECT * FROM azs_button where id_azs = $1`, id_azs)
	if err != nil {
		return
	}

	err = row.Scan(&azsButton.IdAzs, &azsButton.Price, &azsButton.Button)
	return
}

func (r *Repository) GetAzsButtonAll(ctx context.Context) (azsButtons []AzsButton, err error) {
	rows, err := r.pool.Query(ctx, `SELECT * FROM azs_button`)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var azsButton AzsButton
		if err = rows.Scan(&azsButton.IdAzs, &azsButton.Price, &azsButton.Button); err != nil {
			return
		}
		azsButtons = append(azsButtons, azsButton)
	}
	return
}
