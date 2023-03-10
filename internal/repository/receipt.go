package repository

import (
	"context"
	"fmt"
)

//insert into azs_id_10111991_receipts (time, info) values (1677101899, 'Дата: 22.23.2023 23:38\nКолонка: 2\nБензин: АИ 95\nСумма: 1500.00₽ руб')

type Receipt struct {
	Id   int    `json:"id" db:"id"`
	Time int    `json:"time" db:"time"`
	Info string `json:"info" db:"info"`
}

func getTableName(id_azs int) (table string) {
	return fmt.Sprintf("azs_id_%d_receipts", id_azs)
}

func (r *Repository) CreateReceipt(ctx context.Context, id_azs int) (err error) {
	table := getTableName(id_azs)
	_, err = r.pool.Query(ctx,
		"create table if not exists "+table+
			" (id bigint primary key generated always as identity,"+
			"time bigint,"+
			"info varchar(500) not null);")
	return
}

func (r *Repository) AddReceipt(ctx context.Context, id_azs int, time int, info string) (err error) {
	table := getTableName(id_azs)
	_, err = r.pool.Exec(ctx,
		"insert into "+table+" (time, info) values ($1, $2)", time, info)
	return
}

func (r *Repository) DeleteReceiptAll(ctx context.Context, id_azs int) (err error) {
	table := getTableName(id_azs)
	_, err = r.pool.Exec(ctx, "DROP TABLE "+table)
	return
}

func (r *Repository) DeleteLaterReceipt(ctx context.Context, id_azs int, time int) (err error) {
	table := getTableName(id_azs)
	_, err = r.pool.Exec(ctx, "DELETE FROM "+table+"WHERE time < $1", time)
	return
}

func (r *Repository) GetReceiptInRange(ctx context.Context, id_azs int, time1, time2 int64) (receipts []Receipt, err error) {
	table := getTableName(id_azs)
	rows, err := r.pool.Query(ctx, "SELECT * FROM "+table+" WHERE time > $1 and time < $2", time1, time2)

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var receipt Receipt
		if err = rows.Scan(&receipt.Id, &receipt.Time, &receipt.Info); err != nil {
			return
		}
		receipts = append(receipts, receipt)
	}
	return
}

func (r *Repository) GetReceiptAll(ctx context.Context, id_azs int) (receipts []Receipt, err error) {

	table := getTableName(id_azs)

	rows, err := r.pool.Query(ctx, "SELECT * FROM "+table+" ORDER BY id DESC")

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var receipt Receipt
		if err = rows.Scan(&receipt.Id, &receipt.Time, &receipt.Info); err != nil {
			return
		}
		receipts = append(receipts, receipt)
	}
	return
}
