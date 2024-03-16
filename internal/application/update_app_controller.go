package application

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (a app) getAppUpdateButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !a.validateToken(rw, r.FormValue("token")) {
		return
	}

	id := strings.TrimSpace(r.FormValue("id"))
	idInt, ok := getIntVal(id)

	if ok {
		updateCommand, err := a.repo.GetUpdateCommand(a.ctx, idInt)
		if err == nil {
			sendJson(rw, http.StatusOK, updateCommand)
			return
		}
	}
	sendJsonResponse(rw, http.StatusBadRequest, "Error", "Error")
}

func (a app) setAppUpdateCmd(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !a.validateToken(rw, r.FormValue("token")) {
		return
	}
	id := strings.TrimSpace(r.FormValue("id"))
	idInt, ok := getIntVal(id)

	if ok {
		err := a.repo.UpdateUpdateCommand(a.ctx, idInt, "http://t-azs.ru:8085/public/update/GasStationPro.tar.gz")
		if err == nil {
			sendJsonResponse(rw, http.StatusOK, "Ok", "Ok")
			return
		}
	}
	sendJsonResponse(rw, http.StatusBadRequest, "Error", "Error")
}

func (a app) resetAppUpdateButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !a.validateToken(rw, r.FormValue("token")) {
		return
	}

	a.resetAppUpdateAzs(rw, r, p)
}

func (a app) resetAppUpdateAzs(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idInt, ok := getIntVal(strings.TrimSpace(r.FormValue("id")))

	if !ok {
		sendJsonResponse(rw, http.StatusBadRequest, "Error id", "Error")
		return
	}

	err := a.repo.UpdateUpdateCommand(a.ctx, idInt, "")
	if err != nil {
		sendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")

		return
	}
	sendJsonResponse(rw, http.StatusOK, "Ok", "Ok")

}
