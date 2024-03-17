package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type Receipt struct {
	Id           int     `json:"id" db:"id"`
	Time         int     `json:"time" db:"time"`
	Data         string  `json:"data" db:"data"`
	NumOfAzsNode int     `json:"num_azs_node" db:"num_azs_node"`
	GasType      string  `json:"gas_type" db:"gas_type"`
	CountLitres  string  `json:"count_litres" db:"count_litres"`
	Cash         float32 `json:"cash" db:"cash"`
	Cashless     float32 `json:"cashless" db:"cashless"`
	Sum          string  `json:"sum" db:"sum"`
}

type FilterParams struct {
	StartTime   int64  // Start time for filtering
	EndTime     int64  // End time
	PaymentType string // Type of payment: "cash", "cashless", or an empty string for all
}

func getTableName(id_azs int) (table string) {
	return fmt.Sprintf("azs_id_%d_receipts_v2", id_azs)
}

func (r *Repository) fetchReceipts(ctx context.Context, query string, args ...interface{}) ([]Receipt, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var receipts []Receipt
	for rows.Next() {
		var receipt Receipt
		if err := rows.Scan(&receipt.Id, &receipt.Time, &receipt.Data, &receipt.NumOfAzsNode, &receipt.GasType, &receipt.CountLitres, &receipt.Cash, &receipt.Cashless, &receipt.Sum); err != nil {
			return nil, err
		}
		receipts = append(receipts, receipt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return receipts, nil
}

func (r *Repository) CreateReceipt(ctx context.Context, id_azs int) error {
	table := getTableName(id_azs)
	query := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    time BIGINT,
    data VARCHAR(20) NOT NULL,
    num_azs_node INT,
    gas_type VARCHAR(10) NOT NULL,
    count_litres VARCHAR(20) NOT NULL,
    cash NUMERIC(10, 2),
    cashless NUMERIC(10, 2), 
    sum VARCHAR(20) NOT NULL
);`, table)
	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to create table %s: %v", table, err)
		return err
	}
	return nil
}

func (r *Repository) AddReceipt(ctx context.Context, id_azs int, receipt Receipt) error {
	table := getTableName(id_azs)
	query := fmt.Sprintf("INSERT INTO %s (time, data, num_azs_node, gas_type, count_litres, cash, cashless, sum) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", table)
	_, err := r.pool.Exec(ctx, query, receipt.Time, receipt.Data, receipt.NumOfAzsNode, receipt.GasType, receipt.CountLitres, receipt.Cash, receipt.Cashless, receipt.Sum)
	if err != nil {
		log.Printf("Failed to add receipt to table %s: %v", table, err)
		return err
	}
	return nil
}

func (r *Repository) DeleteReceiptAll(ctx context.Context, id_azs int) error {
	table := getTableName(id_azs)

	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", table)
	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		log.Printf("Error deleting all receipts from table %s: %v", table, err)
		return err
	}
	return nil
}

func (r *Repository) DeleteLaterReceipt(ctx context.Context, id_azs int, time int) error {
	table := getTableName(id_azs)
	query := fmt.Sprintf("DELETE FROM %s WHERE time < $1", table)
	_, err := r.pool.Exec(ctx, query, time)
	if err != nil {
		return fmt.Errorf("failed to delete from %s: %v", table, err)
	}
	return nil
}

func (r *Repository) GetReceiptInRange(ctx context.Context, id_azs int, time1, time2 int64) ([]Receipt, error) {
	table := getTableName(id_azs)
	query := fmt.Sprintf("SELECT id, time, data, num_azs_node, gas_type, count_litres, cash, cashless, sum FROM %s WHERE time > $1 AND time < $2 AND sum != '0.00' ORDER BY id DESC", table)
	return r.fetchReceipts(ctx, query, time1, time2)
}

func (r *Repository) GetReceiptAll(ctx context.Context, id_azs int) ([]Receipt, error) {
	table := getTableName(id_azs)
	query := fmt.Sprintf("SELECT id, time, data, num_azs_node, gas_type, count_litres, cash, cashless, sum FROM %s ORDER BY id DESC", table)
	return r.fetchReceipts(ctx, query)
}

func (r *Repository) GetReceiptsFiltered(ctx context.Context, id_azs int, filter FilterParams) ([]Receipt, error) {
	table := getTableName(id_azs)
	baseQuery := fmt.Sprintf("SELECT id, time, data, num_azs_node, gas_type, count_litres, cash, cashless, sum FROM %s", table)

	whereClauses := []string{"1=1"}
	args := []interface{}{}
	argCount := 1

	if filter.StartTime != 0 && filter.EndTime != 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("time BETWEEN $%d AND $%d", argCount, argCount+1))
		args = append(args, filter.StartTime, filter.EndTime)
		argCount += 2
	}

	if filter.PaymentType == "cash" {
		whereClauses = append(whereClauses, "cash > 0")
	} else if filter.PaymentType == "cashless" {
		whereClauses = append(whereClauses, "cashless > 0")
	}

	query := fmt.Sprintf("%s WHERE %s ORDER BY id DESC", baseQuery, strings.Join(whereClauses, " AND "))

	return r.fetchReceipts(ctx, query, args...)
}

func ParseReceiptFromJson(receiptJson string) (receipt Receipt, err error) {

	err = json.Unmarshal([]byte(receiptJson), &receipt)

	if err != nil {
		return
	}
	return
}
