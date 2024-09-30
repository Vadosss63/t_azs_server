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

const (
	yaAzsInfoName = "ya_azs_info"
	columnAzsId   = "id_azs"
	columnLat     = "lat"
	columnLon     = "lon"
	enable        = "enable"
)

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
	query := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
    %s  BIGINT,
    %s  DOUBLE PRECISION,
    %s  DOUBLE PRECISION,
	%s  BOOLEAN
);`, yaAzsInfoName, columnAzsId, columnLat, columnLon, enable)
	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create %s table: %w", yaAzsInfoName, err)
	}
	return nil
}

func (r *YaAzsRepo) DeleteTable(ctx context.Context) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s`, yaAzsInfoName)
	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to drop %s table: %w", yaAzsInfoName, err)
	}
	return nil
}

func (r *YaAzsRepo) Add(ctx context.Context, idAzs int) error {
	query := fmt.Sprintf(`INSERT INTO %s (%s, %s, %s, %s ) VALUES ($1, 0, 0, FALSE)`, yaAzsInfoName, columnAzsId, columnLat, columnLon, enable)
	_, err := r.pool.Exec(ctx, query, idAzs)
	if err != nil {
		return fmt.Errorf("failed to add to %s: %w", yaAzsInfoName, err)
	}
	return nil
}

func (r *YaAzsRepo) UpdateLocation(ctx context.Context, idAzs int, location Location) error {
	query := fmt.Sprintf(`UPDATE %s SET %s = $2, %s = $3 WHERE %s = $1`, yaAzsInfoName, columnLat, columnLon, columnAzsId)
	_, err := r.pool.Exec(ctx, query, idAzs, location.Lat, location.Lon)
	if err != nil {
		return fmt.Errorf("failed to update %s: %w", yaAzsInfoName, err)
	}
	return nil
}

func (r *YaAzsRepo) UpdateEnable(ctx context.Context, idAzs int, isEnable bool) error {
	query := fmt.Sprintf(`UPDATE %s SET %s = $2 WHERE %s = $1`, yaAzsInfoName, enable, columnAzsId)
	_, err := r.pool.Exec(ctx, query, idAzs, isEnable)
	if err != nil {
		return fmt.Errorf("failed to update %s: %w", yaAzsInfoName, err)
	}
	return nil
}

func (r *YaAzsRepo) Delete(ctx context.Context, idAzs int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE %s = $1`, yaAzsInfoName, columnAzsId)
	_, err := r.pool.Exec(ctx, query, idAzs)
	if err != nil {
		return fmt.Errorf("failed to delete from %s: %w", yaAzsInfoName, err)
	}
	return nil
}

func (r *YaAzsRepo) GetLocation(ctx context.Context, idAzs int) (Location, error) {
	query := fmt.Sprintf(`SELECT %s, %s FROM %s WHERE %s = $1`, columnLat, columnLon, yaAzsInfoName, columnAzsId)
	row := r.pool.QueryRow(ctx, query, idAzs)

	var location Location
	err := row.Scan(&location.Lat, &location.Lon)
	if err != nil {
		return location, fmt.Errorf("failed to get from %s: %w", yaAzsInfoName, err)
	}

	return location, nil
}

func (r *YaAzsRepo) GetEnable(ctx context.Context, idAzs int) (bool, error) {
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE %s = $1`, enable, yaAzsInfoName, columnAzsId)
	row := r.pool.QueryRow(ctx, query, idAzs)

	var isEnable bool
	err := row.Scan(&isEnable)
	if err != nil {
		return isEnable, fmt.Errorf("failed to get from %s: %w", yaAzsInfoName, err)
	}

	return isEnable, nil
}

func (r *YaAzsRepo) GetEnableAll(ctx context.Context) ([]Station, error) {

	query := fmt.Sprintf(`SELECT %s, %s, %s FROM %s WHERE %s = true`, columnAzsId, columnLat, columnLon, yaAzsInfoName, enable)

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query %s: %w", yaAzsInfoName, err)
	}
	defer rows.Close()

	var stations []Station
	for rows.Next() {
		var station Station
		var id int
		if err := rows.Scan(&id, &station.Location.Lat, &station.Location.Lon); err != nil {
			return nil, fmt.Errorf("failed to scan from %s: %w", yaAzsInfoName, err)
		}
		station.Enable = true
		station.Id = strconv.Itoa(id)
		stations = append(stations, station)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating over %s: %w", yaAzsInfoName, err)
	}

	return stations, nil
}

func (r *YaAzsRepo) GetEnableList(ctx context.Context) ([]int, error) {

	query := fmt.Sprintf(`SELECT %s FROM %s WHERE %s = true`, columnAzsId, yaAzsInfoName, enable)

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query %s: %w", yaAzsInfoName, err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan from %s: %w", yaAzsInfoName, err)
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating over %s: %w", yaAzsInfoName, err)
	}

	return ids, nil
}
