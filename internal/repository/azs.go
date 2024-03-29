package repository

import (
	"context"
	"encoding/json"
)

type AzsStatsData struct {
	Id                  int    `json:"id" db:"id"`
	IdAzs               int    `json:"id_azs" db:"id_azs"`
	IdUser              int    `json:"id_user" db:"id_user"`
	IsAuthorized        int    `json:"is_authorized" db:"is_authorized"`
	CountColum          int    `json:"count_colum" db:"count_colum"`
	IsSecondPriceEnable int    `json:"is_second_price" db:"is_second_price"`
	Time                string `json:"time" db:"time"`
	Name                string `json:"name" db:"name"`
	Address             string `json:"address" db:"address"`
	Stats               string `json:"stats" db:"stats"`
}

type Info struct {
	CommonOnlineSum   int    `json:"commonOnline"`
	CommonSumCash     int    `json:"commonCash"`
	CommonSumCashless int    `json:"commonCashless"`
	DailyOnlineSum    int    `json:"dailyOnline"`
	DailySumCash      int    `json:"dailyCash"`
	DailySumCashless  int    `json:"dailyCashless"`
	Version           string `json:"version"`
	IsBlock           bool   `json:"isBlock"`
}

type AzsNode struct {
	CommonLiters       string  `json:"commonLiters"`
	DailyLiters        string  `json:"dailyLiters"`
	FuelVolume         string  `json:"fuelVolume"`
	TypeFuel           string  `json:"typeFuel"`
	Price              float32 `json:"price"`
	PriceCashless      float32 `json:"priceCashless"`
	FuelVolumePerc     string  `json:"fuelVolumePerc"`
	Density            string  `json:"density"`
	AverageTemperature string  `json:"averageTemperature"`
	LockFuelValue      int     `json:"lockFuelValue"`
	FuelArrival        int     `json:"fuelArrival"`
}

type AzsStatsDataFull struct {
	Id                  int       `json:"id"`
	IdAzs               int       `json:"id_azs"`
	IsAuthorized        int       `json:"is_authorized"`
	Time                string    `json:"time"`
	Name                string    `json:"name" `
	Address             string    `json:"address" `
	CountColum          int       `json:"count_colum"`
	IsSecondPriceEnable int       `json:"is_second_price"`
	Info                Info      `json:"info"`
	AzsNodes            []AzsNode `json:"azs_nodes"`
}

func (r *Repository) AddAzs(ctx context.Context, id_azs int, is_authorized, count_colum, is_second_price int, time, name, address, stats string) (err error) {
	_, err = r.pool.Exec(ctx,
		`insert into azses (id_azs, id_user, is_authorized, time, name, address, count_colum, stats, is_second_price) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		id_azs, -1, is_authorized, time, name, address, count_colum, stats, is_second_price)
	return
}

func (r *Repository) UpdateAzs(ctx context.Context, azs AzsStatsData) (err error) {
	_, err = r.pool.Exec(ctx,
		`UPDATE azses SET is_authorized = $2, count_colum = $3, time = $4, name = $5, address = $6, stats = $7, is_second_price = $8 WHERE id_azs = $1`,
		azs.IdAzs, azs.IsAuthorized, azs.CountColum, azs.Time, azs.Name, azs.Address, azs.Stats, azs.IsSecondPriceEnable)
	return
}

func (r *Repository) DeleteAzs(ctx context.Context, id_azs int) (err error) {
	_, err = r.pool.Exec(ctx, `DELETE FROM azses WHERE id_azs = $1`, id_azs)
	return
}

func (r *Repository) GetAzs(ctx context.Context, id_azs int) (azs AzsStatsData, err error) {
	row := r.pool.QueryRow(ctx, `SELECT * FROM azses where id_azs = $1`, id_azs)
	if err != nil {
		azs.Id = -1
		return
	}

	err = row.Scan(&azs.Id, &azs.IdAzs, &azs.IdUser, &azs.IsAuthorized, &azs.CountColum, &azs.IsSecondPriceEnable, &azs.Time, &azs.Name, &azs.Address, &azs.Stats)
	if err != nil {
		azs.Id = -1
	}
	return
}

func (r *Repository) GetAzsAll(ctx context.Context) (azses []AzsStatsData, err error) {
	rows, err := r.pool.Query(ctx, `SELECT * FROM azses`)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var azs AzsStatsData
		if err = rows.Scan(&azs.Id, &azs.IdAzs, &azs.IdUser, &azs.IsAuthorized, &azs.CountColum, &azs.IsSecondPriceEnable, &azs.Time, &azs.Name, &azs.Address,
			&azs.Stats); err != nil {
			return
		}
		azses = append(azses, azs)
	}
	return
}

func (r *Repository) AddAzsToUser(ctx context.Context, id_user, id_azs int) (err error) {
	_, err = r.pool.Exec(ctx,
		`UPDATE azses SET id_user = $2 WHERE id_azs = $1`, id_azs, id_user)
	return
}

func (r *Repository) RemoveUserFromAzsAll(ctx context.Context, id_user int) (err error) {
	_, err = r.pool.Exec(ctx,
		`UPDATE azses SET id_user = -1 WHERE id_user = $1`, id_user)
	return
}

func (r *Repository) GetAzsAllForUser(ctx context.Context, id_user int) (azses []AzsStatsData, err error) {
	rows, err := r.pool.Query(ctx, `SELECT * FROM azses where id_user = $1`, id_user)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var azs AzsStatsData
		if err = rows.Scan(&azs.Id, &azs.IdAzs, &azs.IdUser, &azs.IsAuthorized, &azs.CountColum, &azs.IsSecondPriceEnable, &azs.Time, &azs.Name, &azs.Address,
			&azs.Stats); err != nil {
			return
		}
		azses = append(azses, azs)
	}
	return
}

func ParseStats(azsStatsData AzsStatsData) (azsStatsDataFull AzsStatsDataFull, err error) {

	type Values struct {
		Info     Info      `json:"main_info"`
		AzsNodes []AzsNode `json:"azs_nodes"`
	}
	var stats Values
	err = json.Unmarshal([]byte(azsStatsData.Stats), &stats)

	if err != nil {
		return
	}

	azsStatsDataFull = AzsStatsDataFull{
		Id:                  azsStatsData.Id,
		IdAzs:               azsStatsData.IdAzs,
		IsAuthorized:        azsStatsData.IsAuthorized,
		Time:                azsStatsData.Time,
		Name:                azsStatsData.Name,
		Address:             azsStatsData.Address,
		CountColum:          azsStatsData.CountColum,
		IsSecondPriceEnable: azsStatsData.IsSecondPriceEnable,
		Info:                stats.Info,
		AzsNodes:            stats.AzsNodes,
	}

	for i := 0; i < azsStatsDataFull.CountColum; i++ {
		azsStatsDataFull.AzsNodes[i].Price = azsStatsDataFull.AzsNodes[i].Price / 100
		azsStatsDataFull.AzsNodes[i].PriceCashless = azsStatsDataFull.AzsNodes[i].PriceCashless / 100
	}

	return
}
