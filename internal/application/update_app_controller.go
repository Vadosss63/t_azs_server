package application

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (a app) getAppUpdateButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
		updateCommand, err := a.repo.GetUpdateCommand(a.ctx, idInt)
		if err == nil {
			rw.WriteHeader(http.StatusOK)
			json.NewEncoder(rw).Encode(updateCommand)
			return
		}
	}
	rw.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(rw).Encode(answer{Msg: "error", Status: "error id or GetAppUpdateButton"})
}

func (a app) setAppUpdateCmd(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
		err := a.repo.UpdateUpdateCommand(a.ctx, idInt, "http://t-azs.ru:8085/public/update/GasStationPro.tar.gz")
		if err == nil {
			rw.WriteHeader(http.StatusOK)
			json.NewEncoder(rw).Encode(answer{Msg: "Ok", Status: "Ok"})
			return
		}
	}
	rw.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(rw).Encode(answer{Msg: "error", Status: "error"})
}

func (a app) resetAppUpdateButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Header().Set("Content-Type", "application/json")
	tokenReq := strings.TrimSpace(r.FormValue("token"))
	if a.token != tokenReq {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(answer{"error", "invalid token"})
		return
	}

	a.resetAppUpdateAzs(rw, r, p)
}

func (a app) resetAppUpdateAzs(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := strings.TrimSpace(r.FormValue("id"))
	idInt, ok := getIntVal(id)

	if ok {
		err := a.repo.UpdateUpdateCommand(a.ctx, idInt, "")
		if err == nil {
			rw.WriteHeader(http.StatusOK)
			json.NewEncoder(rw).Encode(answer{Msg: "Ok", Status: "Ok"})
			return
		}
	}
	rw.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(rw).Encode(answer{Msg: "error", Status: "error"})
}
