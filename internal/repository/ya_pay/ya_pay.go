package ya_pay

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type YaPay struct {
	IdAzs    int    `json:"id_azs" db:"id_azs"`
	ColumnId int    `json:"columnId" db:"columnId"`
	Status   int    `json:"status" db:"status"`
	Data     string `json:"data" db:"data"`
}

type YaPayRepository interface {
	CreateTable(ctx context.Context) error
	DeleteTable(ctx context.Context) error
	Add(ctx context.Context, idAzs int) error
	Update(ctx context.Context, idAzs, columnId, status int, data string) error
	UpdateStatus(ctx context.Context, idAzs, columnId, status int) error
	ClearData(ctx context.Context, idAzs, columnId int) error
	Delete(ctx context.Context, idAzs int) error
	Get(ctx context.Context, idAzs, columnId int) (YaPay, error)
	GetAll(ctx context.Context) ([]YaPay, error)
}

type YaPayRepo struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *YaPayRepo {
	return &YaPayRepo{pool: pool}
}

func (r *YaPayRepo) CreateTable(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS ya_pay (id_azs  BIGINT, columnId INT, status INT, data VARCHAR(500));`)

	if err != nil {
		return fmt.Errorf("failed to create ya_pay table: %w", err)
	}
	return nil
}

func (r *YaPayRepo) DeleteTable(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, `DROP TABLE IF EXISTS ya_pay`)
	if err != nil {
		return fmt.Errorf("failed to drop ya_pay table: %w", err)
	}
	return nil
}

func (r *YaPayRepo) Add(ctx context.Context, idAzs int) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO ya_pay (id_azs, columnId, status, data) VALUES ($1, 0, 0, "")`, idAzs)
	if err != nil {
		return fmt.Errorf("failed to add to ya_pay: %w", err)
	}
	_, err = r.pool.Exec(ctx, `INSERT INTO ya_pay (id_azs, columnId, status, data) VALUES ($1, 1, 0, "")`, idAzs)
	if err != nil {
		return fmt.Errorf("failed to add to ya_pay: %w", err)
	}
	return nil
}

func (r *YaPayRepo) Update(ctx context.Context, idAzs, columnId, status int, data string) error {
	_, err := r.pool.Exec(ctx, `UPDATE ya_pay SET status = $1, data = $2 WHERE id_azs = $3 and columnId = $4`, status, data, idAzs, columnId)
	if err != nil {
		return fmt.Errorf("failed to update ya_pay: %w", err)
	}
	return nil
}

func (r *YaPayRepo) UpdateStatus(ctx context.Context, idAzs, columnId, status int) error {
	_, err := r.pool.Exec(ctx, `UPDATE ya_pay SET status = $1, WHERE id_azs = $3 and columnId = $4`, status, idAzs, columnId)
	if err != nil {
		return fmt.Errorf("failed to update ya_pay: %w", err)
	}
	return nil
}

func (r *YaPayRepo) ClearData(ctx context.Context, idAzs, columnId int) error {
	_, err := r.pool.Exec(ctx, `UPDATE ya_pay SET data = "" WHERE id_azs = $1 and columnId = $2`, idAzs, columnId)
	if err != nil {
		return fmt.Errorf("failed to clear data in ya_pay: %w", err)
	}
	return nil
}

func (r *YaPayRepo) Delete(ctx context.Context, idAzs int) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM ya_pay WHERE id_azs = $1`, idAzs)
	if err != nil {
		return fmt.Errorf("failed to delete from ya_pay: %w", err)
	}
	return nil
}

func (r *YaPayRepo) Get(ctx context.Context, idAzs, columnId int) (YaPay, error) {
	row := r.pool.QueryRow(ctx, `SELECT id_azs, columnId, status, data FROM ya_pay WHERE  id_azs = $1 and columnId = $2`, idAzs, columnId)

	var yaPay YaPay
	err := row.Scan(&yaPay.IdAzs, &yaPay.ColumnId, &yaPay.Status, &yaPay.Data)
	if err != nil {
		return yaPay, fmt.Errorf("failed to get from ya_pay: %w", err)
	}

	return yaPay, nil
}

func (r *YaPayRepo) GetAll(ctx context.Context) ([]YaPay, error) {
	rows, err := r.pool.Query(ctx, `SELECT id_azs, status, data FROM ya_pay`)
	if err != nil {
		return nil, fmt.Errorf("failed to query ya_pay: %w", err)
	}
	defer rows.Close()

	var yaPays []YaPay
	for rows.Next() {
		var yaPay YaPay
		if err := rows.Scan(&yaPay.IdAzs, &yaPay.ColumnId, &yaPay.Status, &yaPay.Data); err != nil {
			return nil, fmt.Errorf("failed to scan from ya_pay: %w", err)
		}
		yaPays = append(yaPays, yaPay)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating over ya_pay: %w", err)
	}

	return yaPays, nil
}
