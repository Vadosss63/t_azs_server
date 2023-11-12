package repository

import (
	"context"
	"encoding/json"
	"fmt"
)

//insert into azs_id_10111991_receipts (time, info) values (1677101899, 'Дата: 22.23.2023 23:38\nКолонка: 2\nБензин: АИ 95\nСумма: 1500.00₽ руб')

type Receipt struct {
	Id           int    `json:"id" db:"id"`
	Time         int    `json:"time" db:"time"`
	Data         string `json:"data" db:"data"`
	NumOfAzsNode int    `json:"num_azs_node" db:"num_azs_node"`
	GasType      string `json:"gas_type" db:"gas_type"`
	CountLitres  string `json:"count_litres" db:"count_litres"`
	Sum          string `json:"sum" db:"sum"`
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
			"data varchar(20) not null,"+
			"num_azs_node int,"+
			"gas_type varchar(10) not null,"+
			"count_litres varchar(20) not null,"+
			"sum varchar(20) not null);")
	return
}

func (r *Repository) AddReceipt(ctx context.Context, id_azs int, receipt Receipt) (err error) {
	table := getTableName(id_azs)
	_, err = r.pool.Exec(ctx,
		"insert into "+table+" (time, data, num_azs_node, gas_type, count_litres, sum) values ($1, $2, $3, $4, $5, $6)",
		receipt.Time, receipt.Data, receipt.NumOfAzsNode, receipt.GasType, receipt.CountLitres, receipt.Sum)
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
	rows, err := r.pool.Query(ctx, "SELECT * FROM "+table+" WHERE (time > $1 and time < $2 AND sum != '0.00') ORDER BY id DESC", time1, time2)

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var receipt Receipt
		if err = rows.Scan(&receipt.Id, &receipt.Time, &receipt.Data, &receipt.NumOfAzsNode, &receipt.GasType, &receipt.CountLitres, &receipt.Sum); err != nil {
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
		if err = rows.Scan(&receipt.Id, &receipt.Time, &receipt.Data, &receipt.NumOfAzsNode, &receipt.GasType, &receipt.CountLitres, &receipt.Sum); err != nil {
			return
		}
		receipts = append(receipts, receipt)
	}
	return
}

func ParseReceiptFromJson(receiptJson string) (receipt Receipt, err error) {

	err = json.Unmarshal([]byte(receiptJson), &receipt)

	if err != nil {
		return
	}
	return
}
