package repository

import (
	"context"
	"fmt"
)

// create table if not exists azs_receipts
// (
// id      bigint primary key generated always as identity,
// id_azs  bigint,
// time   	bigint,
// info 	varchar(500) not null
// );

type AzsReceiptData struct {
	Id   int    `json:"id" db:"id"`
	Time string `json:"time" db:"time"`
	Info string `json:"info" db:"info"`
}

func (r *Repository) CreateAzsReceipt(ctx context.Context, id_azs int) (err error) {
	_, err = r.pool.Exec(ctx,
		`create table if not exists azs_id_$1_receipts(
		id      bigint primary key generated always as identity,
		time   	bigint,
		info 	varchar(500) not null);`,
		id_azs)
	if err != nil {
		err = fmt.Errorf("failed to exec data: %w", err)
		return
	}
	return
}

func (r *Repository) AddAzsReceipt(ctx context.Context, id_azs int, time int, info string) (err error) {
	_, err = r.pool.Exec(ctx,
		`insert into azs_id_$1_receipts (time, info) values ($2, $3)`,
		id_azs, time, info)
	if err != nil {
		err = fmt.Errorf("failed to exec data: %w", err)
		return
	}
	return
}

func (r *Repository) DeleteAzsReceiptAll(ctx context.Context, id_azs int) (err error) {
	_, err = r.pool.Exec(ctx, `DROP TABLE azs_id_$1_receipts`, id_azs)
	if err != nil {
		err = fmt.Errorf("failed to exec data: %w", err)
		return
	}
	return
}

func (r *Repository) DeleteLaterAzsReceipt(ctx context.Context, id_azs int, time int) (err error) {
	_, err = r.pool.Exec(ctx, `DELETE FROM azs_id_$1_receipts WHERE time < $2`, id_azs, time)
	if err != nil {
		err = fmt.Errorf("failed to exec data: %w", err)
		return
	}
	return
}

func (r *Repository) GetAzsReceiptInRange(ctx context.Context, id_azs int, time1, time2 int) (receipts []AzsReceiptData, err error) {
	rows, err := r.pool.Query(ctx, `SELECT * FROM azs_id_$1_receipts WHERE time > $2 and time < $3`, id_azs, time1, time2)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var receipt AzsReceiptData
		if err = rows.Scan(&receipt.Id, &receipt.Time, &receipt.Info); err != nil {
			err = fmt.Errorf("failed to query data: %w", err)
			return
		}
		receipts = append(receipts, receipt)
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}
	return
}

func (r *Repository) GetAzsReceiptAll(ctx context.Context, id_azs int) (receipts []AzsReceiptData, err error) {
	rows, err := r.pool.Query(ctx, `SELECT * FROM azs_id_$1_receipts`, id_azs)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var receipt AzsReceiptData
		if err = rows.Scan(&receipt.Id, &receipt.Time, &receipt.Info); err != nil {
			err = fmt.Errorf("failed to query data: %w", err)
			return
		}
		receipts = append(receipts, receipt)
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}
	return
}
