package application

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/Vadosss63/t-azs/internal/repository/azs"
	"github.com/Vadosss63/t-azs/internal/repository/receipt"
	"github.com/julienschmidt/httprouter"
)

func (a App) azsStats(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	idInt, ok := GetIntVal(strings.TrimSpace(r.FormValue("id")))
	if !ok {
		SendJsonResponse(rw, http.StatusBadRequest, "Invalid ID format", "Error")
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	address := strings.TrimSpace(r.FormValue("address"))
	countColum, okCountColum := GetIntVal(strings.TrimSpace(r.FormValue("count_colum")))
	isSecondPrice, okIsSecondPrice := GetIntVal(strings.TrimSpace(r.FormValue("is_second_price")))
	stats := strings.TrimSpace(r.FormValue("stats"))

	if name == "" || address == "" || stats == "" || !okCountColum || !okIsSecondPrice {
		SendJsonResponse(rw, http.StatusBadRequest, "All fields must be filled!", "Error")
		return
	}

	if err := a.manageAzs(idInt, countColum, isSecondPrice, name, address, stats); err != nil {
		SendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")
		return
	}

	SendJsonResponse(rw, http.StatusOK, "Operation successful", "Ok")
}

func (a App) manageAzs(idInt, countColum, isSecondPrice int, name, address, stats string) error {
	t := time.Now().Format(time.RFC822)

	azs, err := a.Repo.AzsRepo.Get(a.Ctx, idInt)
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
	return a.Repo.AzsRepo.Update(a.Ctx, azs)
}

func (a App) createAzs(idInt, countColum, isSecondPrice int, name, address, stats, time string) error {
	if err := a.Repo.AzsRepo.Add(a.Ctx, idInt, 0, countColum, isSecondPrice, time, name, address, stats); err != nil {
		return err
	}
	if err := a.Repo.AzsButtonRepo.Add(a.Ctx, idInt); err != nil {
		return err
	}
	if err := a.Repo.UpdaterButtonRepo.Add(a.Ctx, idInt); err != nil {
		return err
	}
	if err := a.Repo.TrblButtonRepo.Add(a.Ctx, idInt); err != nil {
		return err
	}
	if err := a.Repo.YaAzsRepo.Add(a.Ctx, idInt); err != nil {
		return err
	}
	if err := a.Repo.YaPayRepo.Add(a.Ctx, idInt); err != nil {
		return err
	}
	return a.Repo.ReceiptRepo.CreateReceipt(a.Ctx, idInt)
}

func (a App) deleteAsz(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idAzs, ok := GetIntVal(r.FormValue("id_azs"))
	if !ok {
		SendError(rw, "Invalid id_azs", http.StatusBadRequest)
		return
	}

	if err := a.Repo.AzsRepo.Delete(a.Ctx, idAzs); err != nil {
		SendError(rw, "Failed to delete AZS: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.Repo.AzsButtonRepo.Delete(a.Ctx, idAzs); err != nil {
		SendError(rw, "Failed to delete AZS button: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.Repo.TrblButtonRepo.Delete(a.Ctx, idAzs); err != nil {
		SendError(rw, "Failed to delete AZS Log button: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.Repo.UpdaterButtonRepo.Delete(a.Ctx, idAzs); err != nil {
		SendError(rw, "Failed to delete Update Command: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.Repo.YaAzsRepo.Delete(a.Ctx, idAzs); err != nil {
		SendError(rw, "Failed to delete Ya Azs Info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.Repo.YaPayRepo.Delete(a.Ctx, idAzs); err != nil {
		SendError(rw, "Failed to delete YaPay: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.Repo.ReceiptRepo.DeleteAll(a.Ctx, idAzs); err != nil {
		SendError(rw, "Failed to delete all receipts for AZS: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/", http.StatusSeeOther)
}

func (a App) azsReceipt(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, ok := GetIntVal(strings.TrimSpace(r.FormValue("id")))
	receiptJson := strings.TrimSpace(r.FormValue("receipt"))

	if !ok || receiptJson == "" {
		SendJsonResponse(rw, http.StatusBadRequest, "All fields must be filled!", "Error")
		return
	}

	receipt, err := receipt.ParseReceiptFromJson(receiptJson)

	if err != nil {
		SendJsonResponse(rw, http.StatusBadRequest, err.Error(), "Error")
		return
	}
	err = a.Repo.ReceiptRepo.Add(a.Ctx, id, receipt)
	if err != nil {
		SendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")
		return
	}

	SendJsonResponse(rw, http.StatusOK, "Ok", "Ok")
}

func (a App) getAzsButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	idInt, ok := GetIntVal(strings.TrimSpace(r.FormValue("id")))

	if !ok {
		SendJsonResponse(rw, http.StatusBadRequest, "Error id or GetAzsButton", "Error")
		return
	}

	azsButton, err := a.Repo.AzsButtonRepo.Get(a.Ctx, idInt)
	if err != nil {
		SendJsonResponse(rw, http.StatusBadRequest, err.Error(), "Error")
		return
	}
	SendJson(rw, http.StatusOK, azsButton)
}

func (a App) resetAzsButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	a.resetAzs(rw, r, p)
}

func (a App) resetAzs(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := strings.TrimSpace(r.FormValue("id"))
	idInt, ok := GetIntVal(id)

	if ok {
		err := a.Repo.AzsButtonRepo.Update(a.Ctx, idInt, 0, 0)
		if err == nil {
			SendJsonResponse(rw, http.StatusOK, "Ok", "Ok")
			return
		}
	}
	SendJsonResponse(rw, http.StatusBadRequest, "Error", "Error")
}

func (a App) pushAzsButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

	id_azs, ok := GetIntVal(r.FormValue("id_azs"))
	if !ok {
		SendError(rw, "Invalid id_azs value: "+r.FormValue("id_azs"), http.StatusBadRequest)
		return
	}

	pushedBtn := validBtns[r.FormValue("pushedBtn")]
	if pushedBtn == 0 {
		SendError(rw, "Invalid pushedBtn value: "+r.FormValue("pushedBtn"), http.StatusBadRequest)
		return
	}

	value, ok := GetIntVal(r.FormValue("value"))
	if !ok {
		SendError(rw, "Invalid value value: "+r.FormValue("value"), http.StatusBadRequest)
		return
	}

	err := a.Repo.AzsButtonRepo.Update(a.Ctx, id_azs, value, pushedBtn)
	if err != nil {
		SendError(rw, "Failed to update button: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SendJsonResponse(rw, http.StatusOK, "Ok", "Success")
}

func (a App) azsButtonReady(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idAzs, ok := GetIntVal(r.FormValue("id_azs"))
	if !ok {
		SendError(rw, "Invalid id_azs: "+r.FormValue("id_azs"), http.StatusBadRequest)
		return
	}

	azsButton, err := a.Repo.AzsButtonRepo.Get(a.Ctx, idAzs)
	if err != nil {
		SendError(rw, "Error fetching AZS button: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if azsButton.Button == 0 && azsButton.Value == 0 {
		SendJsonResponse(rw, http.StatusOK, "Ok", "ready")
	} else {
		SendJsonResponse(rw, http.StatusOK, "Ok", "not_ready")
	}
}

func (a App) azsPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := strings.TrimSpace(r.FormValue("id_azs"))
	idInt, ok := GetIntVal(id)
	if !ok {
		SendError(rw, "Invalid id_azs: "+id, http.StatusBadRequest)
		return
	}

	azsStats, err := a.Repo.AzsRepo.Get(a.Ctx, idInt)
	if err != nil {
		SendError(rw, "Server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	azsStatsDataFull, err := azs.ParseStats(azsStats)
	if err != nil {
		SendError(rw, "Server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var azsPageTemplate = template.Must(template.ParseFiles(
		filepath.Join("public", "html", "azs_page.html"),
		filepath.Join("public", "html", "user_navi.html"),
	))

	if err := azsPageTemplate.ExecuteTemplate(rw, "azsStatsDataFull", azsStatsDataFull); err != nil {
		SendError(rw, "Server error: "+err.Error(), http.StatusInternalServerError)
	}
}
