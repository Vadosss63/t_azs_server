package application

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

// trunk-ignore(gitleaks/generic-api-key)
var token = "ef4cfcf144999ed560e9f9ad2be18101"

type answer struct {
	Status string `json:"status"`
	Msg    string `json:"Msg"`
}

func (a app) azsStats(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Header().Set("Content-Type", "application/json")
	tokenReq := strings.TrimSpace(r.FormValue("token"))
	if token != tokenReq {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(answer{"error", "invalid token"})
		return
	}

	id := strings.TrimSpace(r.FormValue("id"))
	idInt, ok := getIntVal(id)
	t := time.Now()
	name := strings.TrimSpace(r.FormValue("name"))
	address := strings.TrimSpace(r.FormValue("address"))
	count_colum, ok_count_colum := getIntVal(strings.TrimSpace(r.FormValue("count_colum")))
	stats := strings.TrimSpace(r.FormValue("stats"))
	// fmt.Println(stats)
	answerStat := answer{Msg: "Ok"}
	if !ok || !ok_count_colum || id == "" || name == "" || address == "" || stats == "" {
		answerStat = answer{Msg: "error", Status: "Все поля должны быть заполнены!"}
	} else {
		azs, err := a.repo.GetAzs(a.ctx, idInt)

		if azs.Id == -1 {
			err = a.repo.AddAzs(a.ctx, idInt, 0, count_colum, t.Format(time.RFC822), name, address, stats)

			if err == nil {
				err = a.repo.CreateReceipt(a.ctx, idInt)
			}

		} else if err == nil {
			azs.Time = t.Format(time.RFC822)
			azs.CountColum = count_colum
			azs.Name = name
			azs.Address = address
			azs.Stats = stats
			err = a.repo.UpdateAzs(a.ctx, azs)
		}

		if err != nil {
			answerStat.Status = "error"
			answerStat.Msg = err.Error()
		} else {
			answerStat = answer{"Ok", "Ok"}
		}

	}
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(answerStat)
}

func (a app) azsReceipt(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	rw.Header().Set("Content-Type", "application/json")
	tokenReq := strings.TrimSpace(r.FormValue("token"))
	if token != tokenReq {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(answer{"error", "invalid token"})
		return
	}

	id, ok_id := getIntVal(strings.TrimSpace(r.FormValue("id")))
	time, ok_time := getIntVal(strings.TrimSpace(r.FormValue("time")))
	receipt := strings.TrimSpace(r.FormValue("receipt"))

	if !ok_time || !ok_id || receipt == "" {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(answer{"error", "Все поля должны быть заполнены!"})
		return
	}

	answerStat := answer{Msg: "Ok"}
	err := a.repo.AddReceipt(a.ctx, id, time, receipt)
	if err != nil {
		answerStat.Status = "error"
		answerStat.Msg = err.Error()
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(answerStat)
}
