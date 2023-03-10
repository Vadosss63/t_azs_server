package application

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Vadosss63/t-azs/internal/repository"
	"github.com/julienschmidt/httprouter"
)

type AzsStatsTemplate struct {
	User  repository.User
	Azses []repository.AzsStatsDataFull
}

type AzsReceiptTemplate struct {
	Azs           repository.AzsStatsData
	FormSearchVal string
	ToSearchVal   string
	Receipts      []repository.AzsReceiptData
	Count         int
}

func (a app) ShowHistoryReceiptsPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fromSearchDate := r.FormValue("formSearch")
	toSearchDate := r.FormValue("toSearch")

	// TODO: add checking fo date from < to
	fromSearchTime, fromErr := time.Parse("2006-01-02", fromSearchDate)
	toSearchTime, toErr := time.Parse("2006-01-02", toSearchDate)

	if fromErr != nil || toErr != nil {
		http.Error(rw, "Error parsing: SearchDate", http.StatusBadRequest)
		return
	}

	a.HistoryReceiptsPage(rw, r, p, fromSearchTime, toSearchTime)
}

func (a app) HistoryReceiptsPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params, fromSearchTime, toSearchTime time.Time) {

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

	azs, err := a.repo.GetAzs(a.ctx, id_azs)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	azsReceiptDatas := AzsReceiptTemplate{
		Azs:           azs,
		FormSearchVal: fromSearchTime.Format("2006-01-02"),
		ToSearchVal:   toSearchTime.Format("2006-01-02"),
		Receipts:      receipts,
		Count:         len(receipts),
	}

	lp := filepath.Join("public", "html", "azs_receipt.html")
	navi := filepath.Join("public", "html", "user_navi.html")
	tmpl := template.Must(template.ParseFiles(lp, navi))
	err = tmpl.ExecuteTemplate(rw, "AzsReceiptTemplate", azsReceiptDatas)
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

	azsStatsTemplate := AzsStatsTemplate{
		User:  u,
		Azses: []repository.AzsStatsDataFull{},
	}

	for _, azs_stats := range azs_statses {

		azsStatsDataFull, err := repository.ParseStats(azs_stats)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		azsStatsTemplate.Azses = append(azsStatsTemplate.Azses, azsStatsDataFull)
	}

	lp := filepath.Join("public", "html", "azs_stats.html")
	navi := filepath.Join("public", "html", "user_navi.html")
	tmpl := template.Must(template.ParseFiles(lp, navi))

	err = tmpl.ExecuteTemplate(rw, "AzsStatsTemplate", azsStatsTemplate)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}
