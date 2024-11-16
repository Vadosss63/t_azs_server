package azs_statistics

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Statistics struct {
	Id              int     `json:"id" db:"id"`
	Time            int     `json:"time" db:"time"`
	Date            string  `json:"date" db:"date"`
	DailyCash       float32 `json:"daily_cash" db:"daily_cash"`
	DailyCashless   float32 `json:"daily_cashless" db:"daily_cashless"`
	DailyOnline     float32 `json:"daily_online" db:"daily_online"`
	DailyLitersCol1 float32 `json:"daily_liters_col1" db:"daily_liters_col1"`
	DailyLitersCol2 float32 `json:"daily_liters_col2" db:"daily_liters_col2"`
	FuelArrivalCol1 float32 `json:"fuel_arrival_col1" db:"fuel_arrival_col1"`
	FuelArrivalCol2 float32 `json:"fuel_arrival_col2" db:"fuel_arrival_col2"`
}

type StatisticsFilterParams struct {
	StartTime int64 // Start time for filtering
	EndTime   int64 // End time
}

type StatisticsRepository interface {
	CreateStatisticsTable(ctx context.Context, id_azs int) error
	AddStatistics(ctx context.Context, id_azs int, stats Statistics) error
	DeleteAllStatistics(ctx context.Context, id_azs int) error
	GetFilteredStatistics(ctx context.Context, id_azs int, filter StatisticsFilterParams) ([]Statistics, error)
}

type StatisticsRepo struct {
	pool *pgxpool.Pool
}

func NewStatisticsRepository(pool *pgxpool.Pool) *StatisticsRepo {
	return &StatisticsRepo{pool: pool}
}

func getStatisticsTableName(id_azs int) string {
	return fmt.Sprintf("azs_id_%d_statistics", id_azs)
}

func (r *StatisticsRepo) fetchStatistics(ctx context.Context, query string, args ...interface{}) ([]Statistics, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statistics []Statistics
	for rows.Next() {
		var stat Statistics
		if err := rows.Scan(&stat.Id, &stat.Time, &stat.Date, &stat.DailyCash, &stat.DailyCashless, &stat.DailyOnline,
			&stat.DailyLitersCol1, &stat.DailyLitersCol2, &stat.FuelArrivalCol1, &stat.FuelArrivalCol2); err != nil {
			return nil, err
		}
		statistics = append(statistics, stat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return statistics, nil
}

func (r *StatisticsRepo) CreateStatisticsTable(ctx context.Context, id_azs int) error {
	table := getStatisticsTableName(id_azs)
	query := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    time BIGINT,
    date VARCHAR(20) NOT NULL,
    daily_cash NUMERIC(15, 2),
    daily_cashless NUMERIC(15, 2),
    daily_online NUMERIC(15, 2),
    daily_liters_col1 NUMERIC(10, 2),
    daily_liters_col2 NUMERIC(10, 2),
    fuel_arrival_col1 NUMERIC(10, 2),
    fuel_arrival_col2 NUMERIC(10, 2)
);`, table)

	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to create table %s: %v", table, err)
		return err
	}
	return nil
}

func (r *StatisticsRepo) AddStatistics(ctx context.Context, id_azs int, stats Statistics) error {
	table := getStatisticsTableName(id_azs)
	query := fmt.Sprintf("INSERT INTO %s (time, date, daily_cash, daily_cashless, daily_online, daily_liters_col1, daily_liters_col2, fuel_arrival_col1, fuel_arrival_col2) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", table)
	_, err := r.pool.Exec(ctx, query, stats.Time, stats.Date, stats.DailyCash, stats.DailyCashless, stats.DailyOnline, stats.DailyLitersCol1, stats.DailyLitersCol2, stats.FuelArrivalCol1, stats.FuelArrivalCol2)
	if err != nil {
		log.Printf("Failed to add statistics to table %s: %v", table, err)
		return err
	}
	return nil
}

func (r *StatisticsRepo) DeleteAllStatistics(ctx context.Context, id_azs int) error {
	table := getStatisticsTableName(id_azs)
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", table)
	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		log.Printf("Error deleting all statistics from table %s: %v", table, err)
		return err
	}
	return nil
}

func (r *StatisticsRepo) GetFilteredStatistics(ctx context.Context, id_azs int, filter StatisticsFilterParams) ([]Statistics, error) {
	table := getStatisticsTableName(id_azs)
	baseQuery := fmt.Sprintf("SELECT id, time, date, daily_cash, daily_cashless, daily_online, daily_liters_col1, daily_liters_col2, fuel_arrival_col1, fuel_arrival_col2 FROM %s", table)

	whereClauses := []string{"1=1"}
	args := []interface{}{}
	argCount := 1

	if filter.StartTime != 0 && filter.EndTime != 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("time BETWEEN $%d AND $%d", argCount, argCount+1))
		args = append(args, filter.StartTime, filter.EndTime)
		argCount += 2
	}

	query := fmt.Sprintf("%s WHERE %s ORDER BY id DESC", baseQuery, strings.Join(whereClauses, " AND "))

	return r.fetchStatistics(ctx, query, args...)
}

func ParseStatisticsFromJson(statisticsJson string) (stats Statistics, err error) {
	err = json.Unmarshal([]byte(statisticsJson), &stats)
	if err != nil {
		return
	}
	return
}
