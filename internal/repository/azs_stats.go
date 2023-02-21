package repository

import (
	"context"
	"fmt"
)

// create table if not exists azses
// (
//  id      bigint primary key generated always as identity,
// 	id_azs  bigint,
// 	is_authorized  int,
//  time   varchar(100) not null,
//  name    varchar(100) not null,
//  address varchar(100) not null,
// 	stats varchar(500) not null
// );
// insert into users (login, hashed_password, name, surname)
// values ('alextonkonogov', '827ccb0eea8a706c4c34a16891f84e7b', 'Alex', 'Tonkonogov');

type AzsStatsData struct {
	Id           int    `json:"id" db:"id"`
	IdAzs        int    `json:"id_azs" db:"id_azs"`
	IsAuthorized int    `json:"is_authorized" db:"is_authorized"`
	Time         string `json:"time" db:"time"`
	Name         string `json:"name" db:"name"`
	Address      string `json:"address" db:"address"`
	Stats        string `json:"stats" db:"stats"`
}

func (r *Repository) AddAzs(ctx context.Context, id_azs int, is_authorized int, time, name, address, stats string) (err error) {
	_, err = r.pool.Exec(ctx,
		`insert into azses (id_azs, is_authorized, time, name, address, stats) values ($1, $2, $3, $4, $5, $6)`,
		id_azs, is_authorized, time, name, address, stats)
	if err != nil {
		err = fmt.Errorf("failed to exec data: %w", err)
		return
	}
	return
}

func (r *Repository) UpdateAzsStats(ctx context.Context, azs AzsStatsData) (err error) {
	_, err = r.pool.Exec(ctx,
		`UPDATE azses SET is_authorized = '$2', time = '$3', name = '$4', address = '$5', stats = '$6' WHERE id_azs = $1`,
		azs.IdAzs, azs.IsAuthorized, azs.Time, azs.Name, azs.Address, azs.Stats)
	if err != nil {
		err = fmt.Errorf("failed to exec data: %w", err)
		return
	}
	return
}

func (r *Repository) DeleteAzsStats(ctx context.Context, id_azs int) (err error) {
	_, err = r.pool.Exec(ctx, `DELETE FROM azses WHERE id_azs = $1`, id_azs)
	if err != nil {
		err = fmt.Errorf("failed to exec data: %w", err)
		return
	}
	return
}

func (r *Repository) GetAzs(ctx context.Context, id_azs int) (azs AzsStatsData, err error) {
	row := r.pool.QueryRow(ctx, `SELECT * FROM azses where id_azs = $1`, id_azs)
	if err != nil {
		azs.Id = -1
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}

	err = row.Scan(&azs.Id, &azs.IdAzs, &azs.IsAuthorized, &azs.Time, &azs.Name, &azs.Address, &azs.Stats)
	if err != nil {
		azs.Id = -1
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}
	return
}

func (r *Repository) GetAzsAll(ctx context.Context) (azses []AzsStatsData, err error) {
	rows, err := r.pool.Query(ctx, `SELECT * FROM azses`)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var azs AzsStatsData
		if err = rows.Scan(&azs.Id, &azs.IdAzs, &azs.IsAuthorized, &azs.Time, &azs.Name,
			&azs.Address, &azs.Stats); err != nil {
			err = fmt.Errorf("failed to query data: %w", err)
			return
		}
		azses = append(azses, azs)
	}
	if err = rows.Err(); err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}
	return
}
