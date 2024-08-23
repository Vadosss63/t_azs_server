package repository

import (
	"context"
	"fmt"
	"time"
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
	Id                string    `json:"Id"`
	DateCreate        time.Time `json:"DateCreate"`
	OrderType         string    `json:"OrderType"`
	OrderVolume       float64   `json:"OrderVolume"`
	StationId         string    `json:"StationId"`
	StationExtendedId string    `json:"StationExtendedId"`
	ColumnId          int       `json:"ColumnId"`
	FuelId            string    `json:"FuelId"`
	FuelMarka         string    `json:"FuelMarka"`
	PriceId           string    `json:"PriceId"`
	FuelExtendedId    string    `json:"FuelExtendedId"`
	PriceFuel         float64   `json:"PriceFuel"`
	Sum               float64   `json:"Sum"`
	Litre             float64   `json:"Litre"`
	SumPaid           float64   `json:"SumPaid"`
	Status            string    `json:"Status"`
	DateEnd           time.Time `json:"DateEnd"`
	ReasonId          string    `json:"ReasonId"`
	Reason            string    `json:"Reason"`
	LitreCompleted    float64   `json:"LitreCompleted"`
	SumPaidCompleted  float64   `json:"SumPaidCompleted"`
	ContractId        string    `json:"ContractId"`
}

type PriceEntry struct {
	StationId string  `json:"StationId"`
	ProductId string  `json:"ProductId"`
	Price     float64 `json:"Price"`
}

type StationStatus struct {
	ID      string
	Active  bool
	Columns map[int]bool // Ключ - ID колонки, значение - активность колонки
}

const (
	yaAzsInfoName = "ya_azs_info"
	columnAzsId   = "id_azs"
	columnLat     = "lat"
	columnLon     = "lon"
	enable        = "enable"
)

func (r *Repository) CreateYaAzsInfoTable(ctx context.Context) error {
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

func (r *Repository) DeleteYaAzsInfoTable(ctx context.Context) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s`, yaAzsInfoName)
	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to drop %s table: %w", yaAzsInfoName, err)
	}
	return nil
}

func (r *Repository) AddYaAzsInfo(ctx context.Context, idAzs int) error {
	query := fmt.Sprintf(`INSERT INTO %s (%s, %s, %s, %s ) VALUES ($1, 0, 0, FALSE)`, yaAzsInfoName, columnAzsId, columnLat, columnLon, enable)
	_, err := r.pool.Exec(ctx, query, idAzs)
	if err != nil {
		return fmt.Errorf("failed to add to %s: %w", yaAzsInfoName, err)
	}
	return nil
}

func (r *Repository) UpdateYaAzsInfoLocation(ctx context.Context, idAzs int, location Location) error {
	query := fmt.Sprintf(`UPDATE %s SET %s = $2, %s = $3 WHERE %s = $1`, yaAzsInfoName, columnLat, columnLon, columnAzsId)
	_, err := r.pool.Exec(ctx, query, idAzs, location.Lat, location.Lon)
	if err != nil {
		return fmt.Errorf("failed to update %s: %w", yaAzsInfoName, err)
	}
	return nil
}

func (r *Repository) UpdateYaAzsInfoEnable(ctx context.Context, idAzs int, isEnable bool) error {
	query := fmt.Sprintf(`UPDATE %s SET %s = $2 WHERE %s = $1`, yaAzsInfoName, enable, columnAzsId)
	_, err := r.pool.Exec(ctx, query, idAzs, isEnable)
	if err != nil {
		return fmt.Errorf("failed to update %s: %w", yaAzsInfoName, err)
	}
	return nil
}

func (r *Repository) DeleteYaAzsInfo(ctx context.Context, idAzs int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE %s = $1`, yaAzsInfoName, columnAzsId)
	_, err := r.pool.Exec(ctx, query, idAzs)
	if err != nil {
		return fmt.Errorf("failed to delete from %s: %w", yaAzsInfoName, err)
	}
	return nil
}

func (r *Repository) GetYaAzsInfoLocation(ctx context.Context, idAzs int) (Location, error) {
	query := fmt.Sprintf(`SELECT %s, %s FROM %s WHERE %s = $1`, columnLat, columnLon, yaAzsInfoName, columnAzsId)
	row := r.pool.QueryRow(ctx, query, idAzs)

	var location Location
	err := row.Scan(&location.Lat, &location.Lon)
	if err != nil {
		return location, fmt.Errorf("failed to get from %s: %w", yaAzsInfoName, err)
	}

	return location, nil
}

func (r *Repository) GetYaAzsInfoEnable(ctx context.Context, idAzs int) (bool, error) {
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE %s = $1`, enable, yaAzsInfoName, columnAzsId)
	row := r.pool.QueryRow(ctx, query, idAzs)

	var isEnable bool
	err := row.Scan(&isEnable)
	if err != nil {
		return isEnable, fmt.Errorf("failed to get from %s: %w", yaAzsInfoName, err)
	}

	return isEnable, nil
}
