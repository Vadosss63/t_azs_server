package azs_controller

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/Vadosss63/t-azs/internal/application"
	"github.com/Vadosss63/t-azs/internal/repository/azs"
	"github.com/Vadosss63/t-azs/internal/repository/receipt"
	"github.com/julienschmidt/httprouter"
)

type AzsController struct {
	app *application.App
}

func NewController(app *application.App) *AzsController {
	return &AzsController{app: app}
}

func (c AzsController) Routes(router *httprouter.Router) {
	router.GET("/azs/control", c.app.Authorized(c.azsPage))

	router.POST("/azs_stats", c.app.Authorized(c.azsStats))
	router.DELETE("/azs_stats", c.app.Authorized(c.deleteAsz))

	router.POST("/azs_receipt", c.app.Authorized(c.azsReceipt))

	router.POST("/get_azs_button", c.app.Authorized(c.getAzsButton))

	router.POST("/reset_azs_button", c.app.Authorized(c.resetAzsButton))

	router.GET("/reset_azs_button", c.app.Authorized(c.resetAzs))
	router.POST("/push_azs_button", c.app.Authorized(c.pushAzsButton))
	router.GET("/azs_button_ready", c.app.Authorized(c.azsButtonReady))

}

func (c AzsController) azsStats(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	idInt, ok := application.GetIntVal(strings.TrimSpace(r.FormValue("id")))
	if !ok {
		application.SendJsonResponse(rw, http.StatusBadRequest, "Invalid ID format", "Error")
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	address := strings.TrimSpace(r.FormValue("address"))
	countColum, okCountColum := application.GetIntVal(strings.TrimSpace(r.FormValue("count_colum")))
	isSecondPrice, okIsSecondPrice := application.GetIntVal(strings.TrimSpace(r.FormValue("is_second_price")))
	stats := strings.TrimSpace(r.FormValue("stats"))

	if name == "" || address == "" || stats == "" || !okCountColum || !okIsSecondPrice {
		application.SendJsonResponse(rw, http.StatusBadRequest, "All fields must be filled!", "Error")
		return
	}

	if err := c.manageAzs(idInt, countColum, isSecondPrice, name, address, stats); err != nil {
		application.SendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")
		return
	}

	application.SendJsonResponse(rw, http.StatusOK, "Operation successful", "Ok")
}

func (c AzsController) manageAzs(idInt, countColum, isSecondPrice int, name, address, stats string) error {
	t := time.Now().Format(time.RFC822)

	azs, err := c.app.Repo.AzsRepo.Get(c.app.Ctx, idInt)
	if azs.Id == -1 {
		return c.createAzs(idInt, countColum, isSecondPrice, name, address, stats, t)
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
	return c.app.Repo.AzsRepo.Update(c.app.Ctx, azs)
}

func (c AzsController) createAzs(idInt, countColum, isSecondPrice int, name, address, stats, time string) error {
	if err := c.app.Repo.AzsRepo.Add(c.app.Ctx, idInt, 0, countColum, isSecondPrice, time, name, address, stats); err != nil {
		return err
	}
	if err := c.app.Repo.AzsButtonRepo.Add(c.app.Ctx, idInt); err != nil {
		return err
	}
	if err := c.app.Repo.UpdaterButtonRepo.Add(c.app.Ctx, idInt); err != nil {
		return err
	}
	if err := c.app.Repo.TrblButtonRepo.Add(c.app.Ctx, idInt); err != nil {
		return err
	}
	if err := c.app.Repo.YaAzsRepo.Add(c.app.Ctx, idInt); err != nil {
		return err
	}
	if err := c.app.Repo.YaPayRepo.Add(c.app.Ctx, idInt); err != nil {
		return err
	}
	if err := c.app.Repo.AzsStatRepo.CreateStatisticsTable(c.app.Ctx, idInt); err != nil {
		return err
	}
	return c.app.Repo.ReceiptRepo.CreateReceipt(c.app.Ctx, idInt)
}

func (c AzsController) deleteAsz(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idAzs, ok := application.GetIntVal(r.FormValue("id_azs"))
	if !ok {
		application.SendError(rw, "Invalid id_azs", http.StatusBadRequest)
		return
	}

	if err := c.app.Repo.AzsRepo.Delete(c.app.Ctx, idAzs); err != nil {
		application.SendError(rw, "Failed to delete AZS: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := c.app.Repo.AzsButtonRepo.Delete(c.app.Ctx, idAzs); err != nil {
		application.SendError(rw, "Failed to delete AZS button: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := c.app.Repo.TrblButtonRepo.Delete(c.app.Ctx, idAzs); err != nil {
		application.SendError(rw, "Failed to delete AZS Log button: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := c.app.Repo.UpdaterButtonRepo.Delete(c.app.Ctx, idAzs); err != nil {
		application.SendError(rw, "Failed to delete Update Command: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := c.app.Repo.YaAzsRepo.Delete(c.app.Ctx, idAzs); err != nil {
		application.SendError(rw, "Failed to delete Ya Azs Info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := c.app.Repo.YaPayRepo.Delete(c.app.Ctx, idAzs); err != nil {
		application.SendError(rw, "Failed to delete YaPay: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := c.app.Repo.ReceiptRepo.DeleteAll(c.app.Ctx, idAzs); err != nil {
		application.SendError(rw, "Failed to delete all receipts for AZS: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := c.app.Repo.AzsStatRepo.DeleteAllStatistics(c.app.Ctx, idAzs); err != nil {
		application.SendError(rw, "Failed to delete all statistics for AZS: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/", http.StatusSeeOther)
}

func (c AzsController) azsReceipt(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, ok := application.GetIntVal(strings.TrimSpace(r.FormValue("id")))
	receiptJson := strings.TrimSpace(r.FormValue("receipt"))

	if !ok || receiptJson == "" {
		application.SendJsonResponse(rw, http.StatusBadRequest, "All fields must be filled!", "Error")
		return
	}

	receipt, err := receipt.ParseReceiptFromJson(receiptJson)

	if err != nil {
		application.SendJsonResponse(rw, http.StatusBadRequest, err.Error(), "Error")
		return
	}
	err = c.app.Repo.ReceiptRepo.Add(c.app.Ctx, id, receipt)
	if err != nil {
		application.SendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")
		return
	}

	application.SendJsonResponse(rw, http.StatusOK, "Ok", "Ok")
}

func (c AzsController) getAzsButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	idInt, ok := application.GetIntVal(strings.TrimSpace(r.FormValue("id")))

	if !ok {
		application.SendJsonResponse(rw, http.StatusBadRequest, "Error id or GetAzsButton", "Error")
		return
	}

	azsButton, err := c.app.Repo.AzsButtonRepo.Get(c.app.Ctx, idInt)
	if err != nil {
		application.SendJsonResponse(rw, http.StatusBadRequest, err.Error(), "Error")
		return
	}
	application.SendJson(rw, http.StatusOK, azsButton)
}

func (c AzsController) resetAzsButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	c.resetAzs(rw, r, p)
}

func (c AzsController) resetAzs(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := strings.TrimSpace(r.FormValue("id"))
	idInt, ok := application.GetIntVal(id)

	if ok {
		err := c.app.Repo.AzsButtonRepo.Update(c.app.Ctx, idInt, 0, 0)
		if err == nil {
			application.SendJsonResponse(rw, http.StatusOK, "Ok", "Ok")
			return
		}
	}
	application.SendJsonResponse(rw, http.StatusBadRequest, "Error", "Error")
}

func (c AzsController) pushAzsButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

	id_azs, ok := application.GetIntVal(r.FormValue("id_azs"))
	if !ok {
		application.SendError(rw, "Invalid id_azs value: "+r.FormValue("id_azs"), http.StatusBadRequest)
		return
	}

	pushedBtn := validBtns[r.FormValue("pushedBtn")]
	if pushedBtn == 0 {
		application.SendError(rw, "Invalid pushedBtn value: "+r.FormValue("pushedBtn"), http.StatusBadRequest)
		return
	}

	value, ok := application.GetIntVal(r.FormValue("value"))
	if !ok {
		application.SendError(rw, "Invalid value value: "+r.FormValue("value"), http.StatusBadRequest)
		return
	}

	err := c.app.Repo.AzsButtonRepo.Update(c.app.Ctx, id_azs, value, pushedBtn)
	if err != nil {
		application.SendError(rw, "Failed to update button: "+err.Error(), http.StatusInternalServerError)
		return
	}

	application.SendJsonResponse(rw, http.StatusOK, "Ok", "Success")
}

func (c AzsController) azsButtonReady(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idAzs, ok := application.GetIntVal(r.FormValue("id_azs"))
	if !ok {
		application.SendError(rw, "Invalid id_azs: "+r.FormValue("id_azs"), http.StatusBadRequest)
		return
	}

	azsButton, err := c.app.Repo.AzsButtonRepo.Get(c.app.Ctx, idAzs)
	if err != nil {
		application.SendError(rw, "Error fetching AZS button: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if azsButton.Button == 0 && azsButton.Value == 0 {
		application.SendJsonResponse(rw, http.StatusOK, "Ok", "ready")
	} else {
		application.SendJsonResponse(rw, http.StatusOK, "Ok", "not_ready")
	}
}

func (c AzsController) azsPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := strings.TrimSpace(r.FormValue("id_azs"))
	idInt, ok := application.GetIntVal(id)
	if !ok {
		application.SendError(rw, "Invalid id_azs: "+id, http.StatusBadRequest)
		return
	}

	azsStats, err := c.app.Repo.AzsRepo.Get(c.app.Ctx, idInt)
	if err != nil {
		application.SendError(rw, "Server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	azsStatsDataFull, err := azs.ParseStats(azsStats)
	if err != nil {
		application.SendError(rw, "Server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var azsPageTemplate = template.Must(template.ParseFiles(
		filepath.Join("public", "html", "azs_page.html"),
		filepath.Join("public", "html", "user_navi.html"),
	))

	if err := azsPageTemplate.ExecuteTemplate(rw, "azsStatsDataFull", azsStatsDataFull); err != nil {
		application.SendError(rw, "Server error: "+err.Error(), http.StatusInternalServerError)
	}
}
