package application

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/Vadosss63/t-azs/internal/repository"
	"github.com/julienschmidt/httprouter"
)

type answer struct {
	Status string `json:"status"`
	Msg    string `json:"Msg"`
}

func (a app) azsStats(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Header().Set("Content-Type", "application/json")
	tokenReq := strings.TrimSpace(r.FormValue("token"))
	if a.token != tokenReq {
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
	if a.token != tokenReq {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(answer{"error", "invalid token"})
		return
	}

	id, ok_id := getIntVal(strings.TrimSpace(r.FormValue("id")))
	receiptJson := strings.TrimSpace(r.FormValue("receipt"))

	if !ok_id || receiptJson == "" {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(answer{"error", "Все поля должны быть заполнены!"})
		return
	}

	answerStat := answer{Msg: "Ok"}

	receipt, err := repository.ParseReceiptFromJson(receiptJson)

	if err != nil {
		answerStat.Status = "error"
		answerStat.Msg = err.Error()
	} else {
		err := a.repo.AddReceipt(a.ctx, id, receipt)
		if err != nil {
			answerStat.Status = "error"
			answerStat.Msg = err.Error()
		}
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(answerStat)
}

func (a app) getAzsButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Header().Set("Content-Type", "application/json")
	tokenReq := strings.TrimSpace(r.FormValue("token"))
	if a.token != tokenReq {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(answer{"error", "invalid token"})
		return
	}

	id := strings.TrimSpace(r.FormValue("id"))
	idInt, ok := getIntVal(id)

	if ok {
		azsButton, err := a.repo.GetAzsButton(a.ctx, idInt)
		if err == nil {
			rw.WriteHeader(http.StatusOK)
			json.NewEncoder(rw).Encode(azsButton)
			return
		}
	}
	rw.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(rw).Encode(answer{Msg: "error", Status: "error id or GetAzsButton"})
}

func (a app) resetAzsButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Header().Set("Content-Type", "application/json")
	tokenReq := strings.TrimSpace(r.FormValue("token"))
	if a.token != tokenReq {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(answer{"error", "invalid token"})
		return
	}

	a.resetAzs(rw, r, p)
}

func (a app) resetAzs(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := strings.TrimSpace(r.FormValue("id"))
	idInt, ok := getIntVal(id)

	if ok {
		err := a.repo.UpdateAzsButton(a.ctx, idInt, 0, 0, 0)
		if err == nil {
			rw.WriteHeader(http.StatusOK)
			json.NewEncoder(rw).Encode(answer{Msg: "Ok", Status: "Ok"})
			return
		}
	}
	rw.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(rw).Encode(answer{Msg: "error", Status: "error"})
}

func (a app) pushAzsButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	id_azs, ok := getIntVal(r.FormValue("id_azs"))

	if !ok {
		rw.WriteHeader(http.StatusBadRequest)
		http.Error(rw, "Ошибка id_azs"+r.FormValue("id_azs"), http.StatusBadRequest)
		return
	}
	pushedBtn, ok := getIntVal(r.FormValue("pushedBtn"))
	if !ok {
		rw.WriteHeader(http.StatusBadRequest)
		http.Error(rw, "Ошибка pushedBtn"+r.FormValue("pushedBtn"), http.StatusBadRequest)
		return
	}
	price, ok := getIntVal(r.FormValue("price"))
	if !ok {
		rw.WriteHeader(http.StatusBadRequest)
		http.Error(rw, "Ошибка price"+r.FormValue("price"), http.StatusBadRequest)
		return
	}
	// pushedBtn = 1 - set price1
	// pushedBtn = 2 - set price2
	//0x11 – сброс колонки 1,
	//0x12 – сброс колонки 2,
	//0x21 – сервисная кнопка 1 (зарезервировано на будущее)
	//0x22 – сервисная кнопка 2 (зарезервировано на будущее)
	//0x23 – сервисная кнопка 3 (зарезервировано на будущее)
	//0xFF – сброс суточных счетчиков
	err := error(nil)

	switch pushedBtn {
	case 0x01:
		err = a.repo.UpdateAzsButton(a.ctx, id_azs, price, 0, 0)
	case 0x02:
		err = a.repo.UpdateAzsButton(a.ctx, id_azs, 0, price, 0)
	case 0x11, 0x12, 0x21, 0x22, 0x23, 0xFF:
		err = a.repo.UpdateAzsButton(a.ctx, id_azs, 0, 0, pushedBtn)
	default:
		rw.WriteHeader(http.StatusBadRequest)
		http.Error(rw, "Ошибка pushedBtn"+r.FormValue("pushedBtn"), http.StatusBadRequest)
		return
	}

	if err == nil {
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(answer{Msg: "Ok", Status: "Ok"})
	}
}

func (a app) azsButtonReady(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	id_azs, ok := getIntVal(r.FormValue("id_azs"))

	if !ok {
		rw.WriteHeader(http.StatusBadRequest)
		http.Error(rw, "Ошибка id_azs"+r.FormValue("id_azs"), http.StatusBadRequest)
		return
	}
	azsButton, err := a.repo.GetAzsButton(a.ctx, id_azs)
	if err == nil && azsButton.Button == 0 && azsButton.Price1 == 0 && azsButton.Price2 == 0 {
		rw.WriteHeader(http.StatusOK)
		return
	}
	rw.WriteHeader(http.StatusBadRequest)
}

func (a app) azsPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	id := strings.TrimSpace(r.FormValue("id_azs"))
	idInt, ok := getIntVal(id)

	if !ok {
		http.Error(rw, "Ошибка id_azs:"+id, http.StatusBadRequest)
		return
	}

	azs_stats, err := a.repo.GetAzs(a.ctx, idInt)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	azsStatsDataFull, err := repository.ParseStats(azs_stats)

	lp := filepath.Join("public", "html", "azs_page.html")
	navi := filepath.Join("public", "html", "user_navi.html")
	tmpl := template.Must(template.ParseFiles(lp, navi))

	err = tmpl.ExecuteTemplate(rw, "azsStatsDataFull", azsStatsDataFull)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}
