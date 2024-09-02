package receipt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Receipt struct {
	Id           int     `json:"id" db:"id"`
	Time         int     `json:"time" db:"time"`
	Date         string  `json:"date" db:"date"`
	NumOfAzsNode int     `json:"num_azs_node" db:"num_azs_node"`
	GasType      string  `json:"gas_type" db:"gas_type"`
	CountLitres  float32 `json:"count_litres" db:"count_litres"`
	Cash         float32 `json:"cash" db:"cash"`
	Cashless     float32 `json:"cashless" db:"cashless"`
	Online       float32 `json:"online" db:"online"`
	Sum          float32 `json:"sum" db:"sum"`
}

type FilterParams struct {
	StartTime   int64  // Start time for filtering
	EndTime     int64  // End time
	PaymentType string // Type of payment: "cash", "cashless", "online" or an empty string for all
}

type ReceiptRepo struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *ReceiptRepo {
	return &ReceiptRepo{pool: pool}
}

func getTableName(id_azs int) (table string) {
	return fmt.Sprintf("azs_id_%d_receipts_v2", id_azs)
}

func (r *ReceiptRepo) fetchReceipts(ctx context.Context, query string, args ...interface{}) ([]Receipt, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var receipts []Receipt
	for rows.Next() {
		var receipt Receipt
		if err := rows.Scan(&receipt.Id, &receipt.Time, &receipt.Date, &receipt.NumOfAzsNode, &receipt.GasType, &receipt.CountLitres, &receipt.Cash, &receipt.Cashless, &receipt.Online, &receipt.Sum); err != nil {
			return nil, err
		}
		receipts = append(receipts, receipt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return receipts, nil
}

func (r *ReceiptRepo) CreateReceipt(ctx context.Context, id_azs int) error {
	table := getTableName(id_azs)
	query := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    time BIGINT,
    date VARCHAR(20) NOT NULL,
    num_azs_node INT,
    gas_type VARCHAR(10) NOT NULL,
    count_litres NUMERIC(10, 2),
    cash NUMERIC(10, 2),
    cashless NUMERIC(10, 2),
	online NUMERIC(10, 2),  
    sum NUMERIC(10, 2)
);`, table)
	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to create table %s: %v", table, err)
		return err
	}
	return nil
}

func (r *ReceiptRepo) Add(ctx context.Context, id_azs int, receipt Receipt) error {
	table := getTableName(id_azs)
	query := fmt.Sprintf("INSERT INTO %s (time, date, num_azs_node, gas_type, count_litres, cash, cashless, online, sum) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", table)
	_, err := r.pool.Exec(ctx, query, receipt.Time, receipt.Date, receipt.NumOfAzsNode, receipt.GasType, receipt.CountLitres, receipt.Cash, receipt.Cashless, receipt.Online, receipt.Sum)
	if err != nil {
		log.Printf("Failed to add receipt to table %s: %v", table, err)
		return err
	}
	return nil
}

func (r *ReceiptRepo) DeleteAll(ctx context.Context, id_azs int) error {
	table := getTableName(id_azs)

	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", table)
	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		log.Printf("Error deleting all receipts from table %s: %v", table, err)
		return err
	}
	return nil
}

func (r *ReceiptRepo) GetFilteredReceipts(ctx context.Context, id_azs int, filter FilterParams) ([]Receipt, error) {
	table := getTableName(id_azs)
	baseQuery := fmt.Sprintf("SELECT id, time, date, num_azs_node, gas_type, count_litres, cash, cashless, online, sum FROM %s", table)

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
	} else if filter.PaymentType == "online" {
		whereClauses = append(whereClauses, "online > 0")
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
