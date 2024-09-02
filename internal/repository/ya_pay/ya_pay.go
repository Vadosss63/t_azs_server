package ya_pay

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	yaPayName        = "ya_pay"
	yaPayColumnID    = "id_azs"
	yaPayColumnValue = "value"
	yaPayColumnData  = "data"
)

type YaPay struct {
	IdAzs int    `json:"id_azs" db:"id_azs"`
	Value int    `json:"value" db:"value"`
	Data  string `json:"data" db:"data"`
}

type YaPayRepo struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *YaPayRepo {
	return &YaPayRepo{pool: pool}
}

func (r *YaPayRepo) CreateTable(ctx context.Context) error {
	query := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
    %s  BIGINT,
    %s  INT,
    %s  VARCHAR(500)
);`, yaPayName, yaPayColumnID, yaPayColumnValue, yaPayColumnData)
	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create %s table: %w", yaPayName, err)
	}
	return nil
}

func (r *YaPayRepo) DeleteTable(ctx context.Context) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s`, yaPayName)
	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to drop %s table: %w", yaPayName, err)
	}
	return nil
}

func (r *YaPayRepo) Add(ctx context.Context, idAzs int) error {
	query := fmt.Sprintf(`INSERT INTO %s (%s, %s, %s) VALUES ($1, 0, 0)`, yaPayName, yaPayColumnID, yaPayColumnValue, yaPayColumnData)
	_, err := r.pool.Exec(ctx, query, idAzs)
	if err != nil {
		return fmt.Errorf("failed to add to %s: %w", yaPayName, err)
	}
	return nil
}

func (r *YaPayRepo) Update(ctx context.Context, idAzs, value int, data string) error {
	query := fmt.Sprintf(`UPDATE %s SET %s = $2, %s = $3 WHERE %s = $1`, yaPayName, yaPayColumnValue, yaPayColumnData, yaPayColumnID)
	_, err := r.pool.Exec(ctx, query, idAzs, value, data)
	if err != nil {
		return fmt.Errorf("failed to update %s: %w", yaPayName, err)
	}
	return nil
}

func (r *YaPayRepo) Delete(ctx context.Context, idAzs int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE %s = $1`, yaPayName, yaPayColumnID)
	_, err := r.pool.Exec(ctx, query, idAzs)
	if err != nil {
		return fmt.Errorf("failed to delete from %s: %w", yaPayName, err)
	}
	return nil
}

func (r *YaPayRepo) Get(ctx context.Context, idAzs int) (YaPay, error) {
	query := fmt.Sprintf(`SELECT %s, %s, %s FROM %s WHERE %s = $1`, yaPayColumnID, yaPayColumnValue, yaPayColumnData, yaPayName, yaPayColumnID)
	row := r.pool.QueryRow(ctx, query, idAzs)

	var yaPay YaPay
	err := row.Scan(&yaPay.IdAzs, &yaPay.Value, &yaPay.Data)
	if err != nil {
		return yaPay, fmt.Errorf("failed to get from %s: %w", yaPayName, err)
	}

	return yaPay, nil
}

func (r *YaPayRepo) GetAll(ctx context.Context) ([]YaPay, error) {
	query := fmt.Sprintf(`SELECT %s, %s, %s FROM %s`, yaPayColumnID, yaPayColumnValue, yaPayColumnData, yaPayName)
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query %s: %w", yaPayName, err)
	}
	defer rows.Close()

	var yaPays []YaPay
	for rows.Next() {
		var yaPay YaPay
		if err := rows.Scan(&yaPay.IdAzs, &yaPay.Value, &yaPay.Data); err != nil {
			return nil, fmt.Errorf("failed to scan from %s: %w", yaPayName, err)
		}
		yaPays = append(yaPays, yaPay)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating over %s: %w", yaPayName, err)
	}

	return yaPays, nil
}
