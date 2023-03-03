package repository

import (
	"context"
	"fmt"
)

//insert into azs_id_10111991_receipts (time, info) values (1677101899, 'Дата: 22.23.2023 23:38\nКолонка: 2\nБензин: АИ 95\nСумма: 1500.00₽ руб')

type AzsReceiptData struct {
	Id   int    `json:"id" db:"id"`
	Time int    `json:"time" db:"time"`
	Info string `json:"info" db:"info"`
}

func (r *Repository) CreateAzsReceipt(ctx context.Context, id_azs int) (err error) {
	query := fmt.Sprintf("create table if not exists azs_id_%d_receipts"+
		"(id bigint primary key generated always as identity,"+
		"time bigint,"+
		"info varchar(500) not null);", id_azs)
	
	_, err = r.pool.Query(ctx, query)

	if err != nil {
		err = fmt.Errorf("failed to exec data: %w", err)
		return
	}
	return
}

func (r *Repository) AddAzsReceipt(ctx context.Context, id_azs int, time int, info string) (err error) {
	table := fmt.Sprintf("azs_id_%d_receipts", id_azs)
	_, err = r.pool.Exec(ctx,
		"insert into "+table+" (time, info) values ($1, $2)", time, info)
	if err != nil {
		err = fmt.Errorf("failed to exec data: %w", err)
		return
	}
	return
}

func (r *Repository) DeleteAzsReceiptAll(ctx context.Context, id_azs int) (err error) {
	table := fmt.Sprintf("azs_id_%d_receipts", id_azs)
	_, err = r.pool.Exec(ctx, "DROP TABLE "+table, id_azs)
	if err != nil {
		err = fmt.Errorf("failed to exec data: %w", err)
		return
	}
	return
}

func (r *Repository) DeleteLaterAzsReceipt(ctx context.Context, id_azs int, time int) (err error) {
	table := fmt.Sprintf("azs_id_%d_receipts", id_azs)
	_, err = r.pool.Exec(ctx, "DELETE FROM "+table+"WHERE time < $2", id_azs, time)
	if err != nil {
		err = fmt.Errorf("failed to exec data: %w", err)
		return
	}
	return
}

func (r *Repository) GetAzsReceiptInRange(ctx context.Context, id_azs int, time1, time2 int64) (receipts []AzsReceiptData, err error) {
	table := fmt.Sprintf("azs_id_%d_receipts", id_azs)
	rows, err := r.pool.Query(ctx, "SELECT * FROM "+table+" WHERE time > $1 and time < $2", time1, time2)
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

	table := fmt.Sprintf("azs_id_%d_receipts", id_azs)
	rows, err := r.pool.Query(ctx, "SELECT * FROM "+table+" ORDER BY id DESC")

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
