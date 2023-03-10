package application

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Vadosss63/t-azs/internal/repository"
	"github.com/julienschmidt/httprouter"
)

func (a app) ShowUsersPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	lp := filepath.Join("public", "html", "users_page.html")
	navi := filepath.Join("public", "html", "admin_navi.html")
	tmpl := template.Must(template.ParseFiles(lp, navi))

	users, err := a.repo.GetUserAll(a.ctx)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = tmpl.ExecuteTemplate(rw, "User", users)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a app) ShowHistoryReceiptsPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fromSearchDate := r.FormValue("formSearch")
	toSearchDate := r.FormValue("toSearch")

	// TODO: add checking fo date from < to
	// Parse the date string
	fromSearchTime, err := time.Parse("2006-01-02", fromSearchDate)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse the date string
	toSearchTime, err := time.Parse("2006-01-02", toSearchDate)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	a.HistoryReceiptsPage(rw, r, p, fromSearchTime, toSearchTime)
}

func (a app) HistoryReceiptsPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params, fromSearchTime, toSearchTime time.Time) {
	// user := r.Context().Value("user").(*repository.User)

	id_azs, ok := getIntVal(r.FormValue("id_azs"))

	if !ok {
		http.Error(rw, "Ошибка id_azs"+r.FormValue("id_azs"), http.StatusBadRequest)
		return
	}

	receipts, err := a.repo.GetAzsReceiptInRange(a.ctx, id_azs, fromSearchTime.Unix(), toSearchTime.Unix())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	lp := filepath.Join("public", "html", "azs_receipt.html")
	navi := filepath.Join("public", "html", "user_navi.html")
	tmpl := template.Must(template.ParseFiles(lp, navi))

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	azs, err := a.repo.GetAzs(a.ctx, id_azs)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	type AzsReceiptDatas struct {
		Azs           repository.AzsStatsData
		FormSearchVal string
		ToSearchVal   string
		Receipts      []repository.AzsReceiptData
		Count         int
	}

	azsReceiptDatas := AzsReceiptDatas{
		Azs:           azs,
		FormSearchVal: fromSearchTime.Format("2006-01-02"),
		ToSearchVal:   toSearchTime.Format("2006-01-02"),
		Receipts:      receipts,
		Count:         len(receipts),
	}

	err = tmpl.ExecuteTemplate(rw, "AzsReceiptDatas", azsReceiptDatas)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a app) UserPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params, u repository.User) {

	azs_statses, err := a.repo.GetAzsAllForUser(a.ctx, u.Id)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	lp := filepath.Join("public", "html", "azs_stats.html")
	navi := filepath.Join("public", "html", "user_navi.html")
	tmpl := template.Must(template.ParseFiles(lp, navi))

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	azses := []repository.AzsStatsDataFull{}

	for _, azs_stats := range azs_statses {

		azsStatsDataFull, err := repository.ParseStats(azs_stats)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		azses = append(azses, azsStatsDataFull)
	}

	type AzsStatsTemplate struct {
		User  repository.User
		Azses []repository.AzsStatsDataFull
	}

	azsStatsTemplate := AzsStatsTemplate{
		User:  u,
		Azses: azses,
	}

	err = tmpl.ExecuteTemplate(rw, "AzsStatsTemplate", azsStatsTemplate)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}
