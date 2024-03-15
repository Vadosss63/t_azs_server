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

const (
	SetPrice1         = 0x01
	SetPrice2         = 0x02
	SetPrice1Cashless = 0x03
	SetPrice2Cashless = 0x04
	BlockAzs          = 0x11
	UnblockAzs        = 0x12
	TakeZReport       = 0x21
	CancelN           = 0x22
	EnableN           = 0x23
	SetFuelArrival1   = 0x31
	SetFuelArrival2   = 0x32
	SetLockFuelValue1 = 0x33
	SetLockFuelValue2 = 0x34
	Encashment        = 0xFF
)

var azsPageTemplate = template.Must(template.ParseFiles(
	filepath.Join("public", "html", "azs_page.html"),
	filepath.Join("public", "html", "user_navi.html"),
))

type answer struct {
	Msg    string `json:"msg"`
	Status string `json:"status"`
}

type responseMessage struct {
	Msg    string `json:"msg"`
	Status string `json:"status"`
}

func sendJsonResponse(rw http.ResponseWriter, statusCode int, msg, status string) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	json.NewEncoder(rw).Encode(responseMessage{Msg: msg, Status: status})
}

func sendError(rw http.ResponseWriter, message string, statusCode int) {
	rw.WriteHeader(statusCode)
	http.Error(rw, message, statusCode)
}

func (a app) azsStats(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	tokenReq := strings.TrimSpace(r.FormValue("token"))
	if a.token != tokenReq {
		sendJsonResponse(rw, http.StatusBadRequest, "Invalid token", "Error")
		return
	}

	id := strings.TrimSpace(r.FormValue("id"))
	idInt, ok := getIntVal(id)
	t := time.Now()
	name := strings.TrimSpace(r.FormValue("name"))
	address := strings.TrimSpace(r.FormValue("address"))
	count_colum, ok_count_colum := getIntVal(strings.TrimSpace(r.FormValue("count_colum")))
	is_second_price, ok_is_second_price := getIntVal(strings.TrimSpace(r.FormValue("is_second_price")))
	stats := strings.TrimSpace(r.FormValue("stats"))

	if !ok || !ok_count_colum || !ok_is_second_price || id == "" || name == "" || address == "" || stats == "" {
		sendJsonResponse(rw, http.StatusOK, "Все поля должны быть заполнены!", "error")

		return
	}

	azs, err := a.repo.GetAzs(a.ctx, idInt)

	if azs.Id == -1 {
		err = a.repo.AddAzs(a.ctx, idInt, 0, count_colum, is_second_price, t.Format(time.RFC822), name, address, stats)
		if err == nil {
			err = a.repo.AddAzsButton(a.ctx, idInt)
			err = a.repo.CreateReceipt(a.ctx, idInt)
		}

	} else if err == nil {
		azs.Time = t.Format(time.RFC822)
		azs.CountColum = count_colum
		azs.Name = name
		azs.Address = address
		azs.Stats = stats
		azs.IsSecondPriceEnable = is_second_price
		err = a.repo.UpdateAzs(a.ctx, azs)
	}

	if err != nil {
		sendJsonResponse(rw, http.StatusOK, err.Error(), "Error")
		return
	}

	sendJsonResponse(rw, http.StatusOK, "Ok", "Ok")
}

func (a app) azsReceipt(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	tokenReq := strings.TrimSpace(r.FormValue("token"))
	if a.token != tokenReq {
		sendJsonResponse(rw, http.StatusBadRequest, "Invalid token", "Error")
		return
	}

	id, ok_id := getIntVal(strings.TrimSpace(r.FormValue("id")))
	receiptJson := strings.TrimSpace(r.FormValue("receipt"))

	if !ok_id || receiptJson == "" {
		sendJsonResponse(rw, http.StatusBadRequest, "Все поля должны быть заполнены!", "Error")
		return
	}

	receipt, err := repository.ParseReceiptFromJson(receiptJson)

	if err != nil {
		sendJsonResponse(rw, http.StatusBadRequest, err.Error(), "Error")
		return
	}
	err = a.repo.AddReceipt(a.ctx, id, receipt)
	if err != nil {
		sendJsonResponse(rw, http.StatusBadRequest, err.Error(), "Error")
		return
	}

	sendJsonResponse(rw, http.StatusOK, "Ok", "Ok")
}

func (a app) getAzsButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	tokenReq := strings.TrimSpace(r.FormValue("token"))
	if a.token != tokenReq {
		sendJsonResponse(rw, http.StatusBadRequest, "Invalid token", "Error")
		return
	}

	id := strings.TrimSpace(r.FormValue("id"))
	idInt, ok := getIntVal(id)

	if !ok {
		sendJsonResponse(rw, http.StatusBadRequest, "Error id or GetAzsButton", "Error")
		return
	}

	azsButton, err := a.repo.GetAzsButton(a.ctx, idInt)
	if err != nil {
		sendJsonResponse(rw, http.StatusBadRequest, err.Error(), "Error")
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(azsButton)
}

func (a app) resetAzsButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	tokenReq := strings.TrimSpace(r.FormValue("token"))
	if a.token != tokenReq {
		sendJsonResponse(rw, http.StatusBadRequest, "error", "invalid token")
		return
	}

	a.resetAzs(rw, r, p)
}

func (a app) resetAzs(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := strings.TrimSpace(r.FormValue("id"))
	idInt, ok := getIntVal(id)

	if ok {
		err := a.repo.UpdateAzsButton(a.ctx, idInt, 0, 0)
		if err == nil {
			sendJsonResponse(rw, http.StatusOK, "Ok", "Ok")
			return
		}
	}
	sendJsonResponse(rw, http.StatusBadRequest, "Error", "Error")
}

func (a app) pushAzsButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	validBtns := map[int]bool{
		SetPrice1: true, SetPrice2: true, SetPrice1Cashless: true, SetPrice2Cashless: true,
		BlockAzs: true, UnblockAzs: true, TakeZReport: true, CancelN: true, EnableN: true,
		SetFuelArrival1: true, SetFuelArrival2: true, SetLockFuelValue1: true, SetLockFuelValue2: true,
		Encashment: true,
	}

	id_azs, ok := getIntVal(r.FormValue("id_azs"))
	if !ok {
		sendError(rw, "Invalid id_azs value: "+r.FormValue("id_azs"), http.StatusBadRequest)
		return
	}

	pushedBtn, ok := getIntVal(r.FormValue("pushedBtn"))
	if !ok || !validBtns[pushedBtn] {
		sendError(rw, "Invalid pushedBtn value: "+r.FormValue("pushedBtn"), http.StatusBadRequest)
		return
	}

	value, ok := getIntVal(r.FormValue("value"))
	if !ok {
		sendError(rw, "Invalid value value: "+r.FormValue("value"), http.StatusBadRequest)
		return
	}

	err := a.repo.UpdateAzsButton(a.ctx, id_azs, value, pushedBtn)
	if err != nil {
		sendError(rw, "Failed to update button: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sendJsonResponse(rw, http.StatusOK, "Ok", "Success")
}

func (a app) azsButtonReady(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idAzs, ok := getIntVal(r.FormValue("id_azs"))
	if !ok {
		sendError(rw, "Invalid id_azs: "+r.FormValue("id_azs"), http.StatusBadRequest)
		return
	}

	azsButton, err := a.repo.GetAzsButton(a.ctx, idAzs)
	if err != nil {
		sendError(rw, "Error fetching AZS button: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if azsButton.Button == 0 && azsButton.Value == 0 {
		sendJsonResponse(rw, http.StatusOK, "Ok", "ready")
	} else {
		sendJsonResponse(rw, http.StatusOK, "Ok", "noready")
	}
}

func (a app) deleteAsz(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idAzs, ok := getIntVal(r.FormValue("id_azs"))
	if !ok {
		sendError(rw, "Invalid id_azs", http.StatusBadRequest)
		return
	}

	if err := a.repo.DeleteAzs(a.ctx, idAzs); err != nil {
		sendError(rw, "Failed to delete AZS: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.repo.DeleteReceiptAll(a.ctx, idAzs); err != nil {
		sendError(rw, "Failed to delete all receipts for AZS: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.repo.DeleteAzsButton(a.ctx, idAzs); err != nil {
		sendError(rw, "Failed to delete AZS button: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/", http.StatusSeeOther)
}

func (a app) azsPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := strings.TrimSpace(r.FormValue("id_azs"))
	idInt, ok := getIntVal(id)
	if !ok {
		sendError(rw, "Invalid id_azs: "+id, http.StatusBadRequest)
		return
	}

	azsStats, err := a.repo.GetAzs(a.ctx, idInt)
	if err != nil {
		sendError(rw, "Server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	azsStatsDataFull, err := repository.ParseStats(azsStats)
	if err != nil {
		sendError(rw, "Server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := azsPageTemplate.ExecuteTemplate(rw, "azsStatsDataFull", azsStatsDataFull); err != nil {
		sendError(rw, "Server error: "+err.Error(), http.StatusInternalServerError)
	}
}
