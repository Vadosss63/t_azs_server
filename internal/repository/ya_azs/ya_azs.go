package ya_azs

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Location struct {
	Lat float64 `json:"Lat" db:"lat"`
	Lon float64 `json:"Lon" db:"lon"`
}

type Column struct {
	Fuels []string `json:"Fuels" db:"fuels"`
}

type Station struct {
	Id       string           `json:"Id"`
	Enable   bool             `json:"Enable"`
	Name     string           `json:"Name"`
	Address  string           `json:"Address"`
	Location Location         `json:"Location"`
	Columns  map[int32]Column `json:"Columns"`
}

type Stations struct {
	Stations []Station `json:"Stations"`
}

type Order struct {
	Id                string  `json:"Id"`
	DateCreate        string  `json:"DateCreate"`
	OrderType         string  `json:"OrderType"`
	OrderVolume       float64 `json:"OrderVolume"`
	StationExtendedId string  `json:"StationExtendedId"`
	ColumnId          int     `json:"ColumnId"`
	FuelExtendedId    string  `json:"FuelExtendedId"`
	PriceFuel         float64 `json:"PriceFuel"`
	Status            string  `json:"Status"`
}

type PriceEntry struct {
	StationId string  `json:"StationId"`
	ProductId string  `json:"ProductId"`
	Price     float64 `json:"Price"`
}

type StationStatus struct {
	ID      string
	Active  bool
	Columns map[int]bool
}

type YaAzsRepository interface {
	CreateTable(ctx context.Context) error
	DeleteTable(ctx context.Context) error
	Add(ctx context.Context, idAzs int) error
	UpdateLocation(ctx context.Context, idAzs int, location Location) error
	UpdateEnable(ctx context.Context, idAzs int, isEnable bool) error
	Delete(ctx context.Context, idAzs int) error
	GetLocation(ctx context.Context, idAzs int) (Location, error)
	GetEnable(ctx context.Context, idAzs int) (bool, error)
	GetEnableAll(ctx context.Context) ([]Station, error)
	GetEnableList(ctx context.Context) ([]int, error)
}

type YaAzsRepo struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *YaAzsRepo {
	return &YaAzsRepo{pool: pool}
}

func (r *YaAzsRepo) CreateTable(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS ya_azs_info (
    			id_azs  BIGINT,
    			lat  DOUBLE PRECISION,
    			lon  DOUBLE PRECISION,
				enable  BOOLEAN);`

	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create ya_azs_info table: %w", err)
	}
	return nil
}

func (r *YaAzsRepo) DeleteTable(ctx context.Context) error {
	query := `DROP TABLE IF EXISTS ya_azs_info`
	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to drop ya_azs_info table: %w", err)
	}
	return nil
}

func (r *YaAzsRepo) Add(ctx context.Context, idAzs int) error {
	var exists bool
	queryCheck := `SELECT EXISTS(SELECT 1 FROM ya_azs_info WHERE id_azs = $1)`
	err := r.pool.QueryRow(ctx, queryCheck, idAzs).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check existence: %w", err)
	}

	if exists {
		return nil
	}

	queryInsert := `INSERT INTO ya_azs_info (id_azs, lat, lon, enable) VALUES ($1, 0, 0, FALSE)`
	_, err = r.pool.Exec(ctx, queryInsert, idAzs)
	if err != nil {
		return fmt.Errorf("failed to add to ya_azs_info: %w", err)
	}
	return nil
}

func (r *YaAzsRepo) UpdateLocation(ctx context.Context, idAzs int, location Location) error {
	query := `UPDATE ya_azs_info SET lat = $2, lon = $3 WHERE id_azs = $1`
	_, err := r.pool.Exec(ctx, query, idAzs, location.Lat, location.Lon)
	if err != nil {
		return fmt.Errorf("failed to update ya_azs_info: %w", err)
	}
	return nil
}

func (r *YaAzsRepo) UpdateEnable(ctx context.Context, idAzs int, isEnable bool) error {
	query := `UPDATE ya_azs_info SET enable = $2 WHERE id_azs = $1`
	_, err := r.pool.Exec(ctx, query, idAzs, isEnable)
	if err != nil {
		return fmt.Errorf("failed to update ya_azs_info: %w", err)
	}
	return nil
}

func (r *YaAzsRepo) Delete(ctx context.Context, idAzs int) error {
	query := `DELETE FROM ya_azs_info WHERE id_azs = $1`
	_, err := r.pool.Exec(ctx, query, idAzs)
	if err != nil {
		return fmt.Errorf("failed to delete from ya_azs_info: %w", err)
	}
	return nil
}

func (r *YaAzsRepo) GetLocation(ctx context.Context, idAzs int) (Location, error) {
	query := `SELECT lat, lon FROM ya_azs_info WHERE id_azs = $1`
	row := r.pool.QueryRow(ctx, query, idAzs)

	var location Location
	err := row.Scan(&location.Lat, &location.Lon)
	if err != nil {
		return location, fmt.Errorf("failed to get from ya_azs_info: %w", err)
	}

	return location, nil
}

func (r *YaAzsRepo) GetEnable(ctx context.Context, idAzs int) (bool, error) {
	query := `SELECT enable FROM ya_azs_info WHERE id_azs = $1`
	row := r.pool.QueryRow(ctx, query, idAzs)

	var isEnable bool
	err := row.Scan(&isEnable)
	if err != nil {
		return isEnable, fmt.Errorf("failed to get from ya_azs_info: %w", err)
	}

	return isEnable, nil
}

func (r *YaAzsRepo) GetEnableAll(ctx context.Context) ([]Station, error) {

	query := `SELECT id_azs, lat, lon FROM ya_azs_info WHERE enable = true`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query ya_azs_info: %w", err)
	}
	defer rows.Close()

	var stations []Station
	for rows.Next() {
		var station Station
		var id int
		if err := rows.Scan(&id, &station.Location.Lat, &station.Location.Lon); err != nil {
			return nil, fmt.Errorf("failed to scan from ya_azs_info: %w", err)
		}
		station.Enable = true
		station.Id = strconv.Itoa(id)
		stations = append(stations, station)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating over ya_azs_info: %w", err)
	}

	return stations, nil
}

func (r *YaAzsRepo) GetEnableList(ctx context.Context) ([]int, error) {

	query := `SELECT id_azs FROM ya_azs_info WHERE enable = true`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query ya_azs_info: %w", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan from ya_azs_info: %w", err)
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating over ya_azs_info: %w", err)
	}

	return ids, nil
}
