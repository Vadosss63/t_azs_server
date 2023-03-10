package repository

import (
	"context"
	"encoding/json"
)

// insert into users (login, hashed_password, name, surname)
// values ('alextonkonogov', '827ccb0eea8a706c4c34a16891f84e7b', 'Alex', 'Tonkonogov');

type AzsStatsData struct {
	Id           int    `json:"id" db:"id"`
	IdAzs        int    `json:"id_azs" db:"id_azs"`
	IdUser       int    `json:"id_user" db:"id_user"`
	IsAuthorized int    `json:"is_authorized" db:"is_authorized"`
	CountColum   int    `json:"count_colum" db:"count_colum"`
	Time         string `json:"time" db:"time"`
	Name         string `json:"name" db:"name"`
	Address      string `json:"address" db:"address"`
	Stats        string `json:"stats" db:"stats"`
}

type Info struct {
	CommonOnlineSum   int `json:"commonOnlineSum"`
	CommonSumCash     int `json:"commonSumCash"`
	CommonSumCashless int `json:"commonSumCashless"`
	DailyOnlineSum    int `json:"dailyOnlineSum"`
	DailySumCash      int `json:"dailySumCash"`
	DailySumCashless  int `json:"dailySumCashless"`
}

type Columns struct {
	CommonLiters        string `json:"commonLiters"`
	DailyLiters         string `json:"dailyLiters"`
	RemainingFuelLiters string `json:"remainingFuelLiters"`
}

type AzsStatsDataFull struct {
	Id           int       `json:"id" db:"id"`
	IdAzs        int       `json:"id_azs" db:"id_azs"`
	IsAuthorized int       `json:"is_authorized" db:"is_authorized"`
	Time         string    `json:"time" db:"time"`
	Name         string    `json:"name" db:"name"`
	Address      string    `json:"address" db:"address"`
	CountColum   int       `json:"count_colum" db:"count_colum"`
	Info         Info      `json:"info"`
	Columns      []Columns `json:"columns"`
}

func (r *Repository) AddAzs(ctx context.Context, id_azs int, is_authorized, count_colum int, time, name, address, stats string) (err error) {
	_, err = r.pool.Exec(ctx,
		`insert into azses (id_azs, id_user, is_authorized, time, name, address, count_colum, stats) values ($1, $2, $3, $4, $5, $6, $7, $8)`,
		id_azs, -1, is_authorized, time, name, address, count_colum, stats)
	return
}

func (r *Repository) UpdateAzs(ctx context.Context, azs AzsStatsData) (err error) {
	_, err = r.pool.Exec(ctx,
		`UPDATE azses SET is_authorized = $2, count_colum = $3, time = $4, name = $5, address = $6, stats = $7 WHERE id_azs = $1`,
		azs.IdAzs, azs.IsAuthorized, azs.CountColum, azs.Time, azs.Name, azs.Address, azs.Stats)
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

	err = row.Scan(&azs.Id, &azs.IdAzs, &azs.IdUser, &azs.IsAuthorized, &azs.CountColum, &azs.Time, &azs.Name, &azs.Address, &azs.Stats)
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
		if err = rows.Scan(&azs.Id, &azs.IdAzs, &azs.IdUser, &azs.IsAuthorized, &azs.CountColum, &azs.Time, &azs.Name, &azs.Address,
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

func (r *Repository) GetAzsAllForUser(ctx context.Context, id_user int) (azses []AzsStatsData, err error) {
	rows, err := r.pool.Query(ctx, `SELECT * FROM azses where id_user = $1`, id_user)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var azs AzsStatsData
		if err = rows.Scan(&azs.Id, &azs.IdAzs, &azs.IdUser, &azs.IsAuthorized, &azs.CountColum, &azs.Time, &azs.Name, &azs.Address,
			&azs.Stats); err != nil {
			return
		}
		azses = append(azses, azs)
	}
	return
}

func ParseStats(azsStatsData AzsStatsData) (azsStatsDataFull AzsStatsDataFull, err error) {

	type Values struct {
		Info    Info      `json:"info"`
		Columns []Columns `json:"columns"`
	}
	var stats Values
	err = json.Unmarshal([]byte(azsStatsData.Stats), &stats)

	if err != nil {
		return
	}

	azsStatsDataFull = AzsStatsDataFull{
		Id:           azsStatsData.Id,
		IdAzs:        azsStatsData.IdAzs,
		IsAuthorized: azsStatsData.IsAuthorized,
		Time:         azsStatsData.Time,
		Name:         azsStatsData.Name,
		Address:      azsStatsData.Address,
		CountColum:   azsStatsData.CountColum,
		Info:         stats.Info,
		Columns:      stats.Columns,
	}
	return
}
