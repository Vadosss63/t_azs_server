package azs_button

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	tableName    = "azs_button_v2"
	columnID     = "id_azs"
	columnValue  = "value"
	columnButton = "button"
)

type AzsButtonRepository interface {
	CreateTable(ctx context.Context) error
	DeleteTable(ctx context.Context) error
	Add(ctx context.Context, idAzs int) error
	Update(ctx context.Context, idAzs, price, button int) error
	Delete(ctx context.Context, idAzs int) error
	Get(ctx context.Context, idAzs int) (AzsButton, error)
	GetAll(ctx context.Context) ([]AzsButton, error)
}

type AzsButtonRepo struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *AzsButtonRepo {
	return &AzsButtonRepo{pool: pool}
}

type AzsButton struct {
	IdAzs  int `json:"id_azs" db:"id_azs"`
	Value  int `json:"value" db:"value"`
	Button int `json:"button" db:"button"`
}

func (r *AzsButtonRepo) CreateTable(ctx context.Context) error {
	query := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
    %s  BIGINT,
    %s  INT,
    %s  INT
);`, tableName, columnID, columnValue, columnButton)
	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create %s table: %w", tableName, err)
	}
	return nil
}

func (r *AzsButtonRepo) DeleteTable(ctx context.Context) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s`, tableName)
	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to drop %s table: %w", tableName, err)
	}
	return nil
}

func (r *AzsButtonRepo) Add(ctx context.Context, idAzs int) error {
	query := fmt.Sprintf(`INSERT INTO %s (%s, %s, %s) VALUES ($1, 0, 0)`, tableName, columnID, columnValue, columnButton)
	_, err := r.pool.Exec(ctx, query, idAzs)
	if err != nil {
		return fmt.Errorf("failed to add to %s: %w", tableName, err)
	}
	return nil
}

func (r *AzsButtonRepo) Update(ctx context.Context, idAzs, price, button int) error {
	query := fmt.Sprintf(`UPDATE %s SET %s = $2, %s = $3 WHERE %s = $1`, tableName, columnValue, columnButton, columnID)
	_, err := r.pool.Exec(ctx, query, idAzs, price, button)
	if err != nil {
		return fmt.Errorf("failed to update %s: %w", tableName, err)
	}
	return nil
}

func (r *AzsButtonRepo) Delete(ctx context.Context, idAzs int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE %s = $1`, tableName, columnID)
	_, err := r.pool.Exec(ctx, query, idAzs)
	if err != nil {
		return fmt.Errorf("failed to delete from %s: %w", tableName, err)
	}
	return nil
}

func (r *AzsButtonRepo) Get(ctx context.Context, idAzs int) (AzsButton, error) {
	query := fmt.Sprintf(`SELECT %s, %s, %s FROM %s WHERE %s = $1`, columnID, columnValue, columnButton, tableName, columnID)
	row := r.pool.QueryRow(ctx, query, idAzs)

	var azsButton AzsButton
	err := row.Scan(&azsButton.IdAzs, &azsButton.Value, &azsButton.Button)
	if err != nil {
		return azsButton, fmt.Errorf("failed to get from %s: %w", tableName, err)
	}

	return azsButton, nil
}

func (r *AzsButtonRepo) GetAll(ctx context.Context) ([]AzsButton, error) {
	query := fmt.Sprintf(`SELECT %s, %s, %s FROM %s`, columnID, columnValue, columnButton, tableName)
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query %s: %w", tableName, err)
	}
	defer rows.Close()

	var azsButtons []AzsButton
	for rows.Next() {
		var azsButton AzsButton
		if err := rows.Scan(&azsButton.IdAzs, &azsButton.Value, &azsButton.Button); err != nil {
			return nil, fmt.Errorf("failed to scan from %s: %w", tableName, err)
		}
		azsButtons = append(azsButtons, azsButton)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating over %s: %w", tableName, err)
	}

	return azsButtons, nil
}
