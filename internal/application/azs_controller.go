package application

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/Vadosss63/t-azs/internal/repository"
	"github.com/julienschmidt/httprouter"
)

func (a app) azsStats(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	idInt, ok := getIntVal(strings.TrimSpace(r.FormValue("id")))
	if !ok {
		sendJsonResponse(rw, http.StatusBadRequest, "Invalid ID format", "Error")
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	address := strings.TrimSpace(r.FormValue("address"))
	countColum, okCountColum := getIntVal(strings.TrimSpace(r.FormValue("count_colum")))
	isSecondPrice, okIsSecondPrice := getIntVal(strings.TrimSpace(r.FormValue("is_second_price")))
	stats := strings.TrimSpace(r.FormValue("stats"))

	if name == "" || address == "" || stats == "" || !okCountColum || !okIsSecondPrice {
		sendJsonResponse(rw, http.StatusBadRequest, "All fields must be filled!", "Error")
		return
	}

	if err := a.manageAzs(idInt, countColum, isSecondPrice, name, address, stats); err != nil {
		sendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")
		return
	}

	sendJsonResponse(rw, http.StatusOK, "Operation successful", "Ok")
}

func (a app) manageAzs(idInt, countColum, isSecondPrice int, name, address, stats string) error {
	t := time.Now().Format(time.RFC822)

	azs, err := a.repo.GetAzs(a.ctx, idInt)
	if azs.Id == -1 {
		return a.createAzs(idInt, countColum, isSecondPrice, name, address, stats, t)
	}

	if err != nil {
		return err
	}

	azs.Time = t
	azs.CountColum = countColum
	azs.Name = name
	azs.Address = address
	azs.Stats = stats
	azs.IsSecondPriceEnable = isSecondPrice
	return a.repo.UpdateAzs(a.ctx, azs)
}

func (a app) createAzs(idInt, countColum, isSecondPrice int, name, address, stats, time string) error {
	if err := a.repo.AddAzs(a.ctx, idInt, 0, countColum, isSecondPrice, time, name, address, stats); err != nil {
		return err
	}
	if err := a.repo.AzsButtonRepo.AddAzsButton(a.ctx, idInt); err != nil {
		return err
	}
	if err := a.repo.UpdaterButtonRepo.AddUpdateCommand(a.ctx, idInt); err != nil {
		return err
	}
	if err := a.repo.TrblButtonRepo.AddLogButton(a.ctx, idInt); err != nil {
		return err
	}
	if err := a.repo.AddYaAzsInfo(a.ctx, idInt); err != nil {
		return err
	}
	if err := a.repo.YaPayRepo.AddYaPay(a.ctx, idInt); err != nil {
		return err
	}
	return a.repo.CreateReceipt(a.ctx, idInt)
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

	if err := a.repo.AzsButtonRepo.DeleteAzsButton(a.ctx, idAzs); err != nil {
		sendError(rw, "Failed to delete AZS button: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.repo.TrblButtonRepo.DeleteLogButton(a.ctx, idAzs); err != nil {
		sendError(rw, "Failed to delete AZS Log button: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.repo.UpdaterButtonRepo.DeleteUpdateCommand(a.ctx, idAzs); err != nil {
		sendError(rw, "Failed to delete Update Command: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.repo.DeleteYaAzsInfo(a.ctx, idAzs); err != nil {
		sendError(rw, "Failed to delete Ya Azs Info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.repo.YaPayRepo.DeleteYaPay(a.ctx, idAzs); err != nil {
		sendError(rw, "Failed to delete YaPay: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.repo.DeleteReceiptAll(a.ctx, idAzs); err != nil {
		sendError(rw, "Failed to delete all receipts for AZS: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/", http.StatusSeeOther)
}

func (a app) azsReceipt(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, ok := getIntVal(strings.TrimSpace(r.FormValue("id")))
	receiptJson := strings.TrimSpace(r.FormValue("receipt"))

	if !ok || receiptJson == "" {
		sendJsonResponse(rw, http.StatusBadRequest, "All fields must be filled!", "Error")
		return
	}

	receipt, err := repository.ParseReceiptFromJson(receiptJson)

	if err != nil {
		sendJsonResponse(rw, http.StatusBadRequest, err.Error(), "Error")
		return
	}
	err = a.repo.AddReceipt(a.ctx, id, receipt)
	if err != nil {
		sendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")
		return
	}

	sendJsonResponse(rw, http.StatusOK, "Ok", "Ok")
}

func (a app) getAzsButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	idInt, ok := getIntVal(strings.TrimSpace(r.FormValue("id")))

	if !ok {
		sendJsonResponse(rw, http.StatusBadRequest, "Error id or GetAzsButton", "Error")
		return
	}

	azsButton, err := a.repo.AzsButtonRepo.GetAzsButton(a.ctx, idInt)
	if err != nil {
		sendJsonResponse(rw, http.StatusBadRequest, err.Error(), "Error")
		return
	}
	sendJson(rw, http.StatusOK, azsButton)
}

func (a app) resetAzsButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	a.resetAzs(rw, r, p)
}

func (a app) resetAzs(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := strings.TrimSpace(r.FormValue("id"))
	idInt, ok := getIntVal(id)

	if ok {
		err := a.repo.AzsButtonRepo.UpdateAzsButton(a.ctx, idInt, 0, 0)
		if err == nil {
			sendJsonResponse(rw, http.StatusOK, "Ok", "Ok")
			return
		}
	}
	sendJsonResponse(rw, http.StatusBadRequest, "Error", "Error")
}

func (a app) pushAzsButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	validBtns := map[string]int{
		"serviceBtn1":       0x01,
		"serviceBtn2":       0x02,
		"serviceBtn3":       0x03,
		"resetCounters":     0x10,
		"blockAzsNode":      0x11,
		"unblockAzsNode":    0x12,
		"setPriceCash1":     0x30,
		"setPriceCashless1": 0x38,
		"setPriceCash2":     0x31,
		"setPriceCashless2": 0x39,
		"setFuelArrival1":   0x48,
		"setLockFuelValue1": 0x40,
		"setFuelArrival2":   0x49,
		"setLockFuelValue2": 0x41,
	}

	id_azs, ok := getIntVal(r.FormValue("id_azs"))
	if !ok {
		sendError(rw, "Invalid id_azs value: "+r.FormValue("id_azs"), http.StatusBadRequest)
		return
	}

	pushedBtn := validBtns[r.FormValue("pushedBtn")]
	if pushedBtn == 0 {
		sendError(rw, "Invalid pushedBtn value: "+r.FormValue("pushedBtn"), http.StatusBadRequest)
		return
	}

	value, ok := getIntVal(r.FormValue("value"))
	if !ok {
		sendError(rw, "Invalid value value: "+r.FormValue("value"), http.StatusBadRequest)
		return
	}

	err := a.repo.AzsButtonRepo.UpdateAzsButton(a.ctx, id_azs, value, pushedBtn)
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

	azsButton, err := a.repo.AzsButtonRepo.GetAzsButton(a.ctx, idAzs)
	if err != nil {
		sendError(rw, "Error fetching AZS button: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if azsButton.Button == 0 && azsButton.Value == 0 {
		sendJsonResponse(rw, http.StatusOK, "Ok", "ready")
	} else {
		sendJsonResponse(rw, http.StatusOK, "Ok", "not_ready")
	}
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

	var azsPageTemplate = template.Must(template.ParseFiles(
		filepath.Join("public", "html", "azs_page.html"),
		filepath.Join("public", "html", "user_navi.html"),
	))

	if err := azsPageTemplate.ExecuteTemplate(rw, "azsStatsDataFull", azsStatsDataFull); err != nil {
		sendError(rw, "Server error: "+err.Error(), http.StatusInternalServerError)
	}
}
