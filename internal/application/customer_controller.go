package application

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
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
	FromTimeVal   string
	ToTimeVal     string
	Receipts      []repository.Receipt
	Count         int
	TotalSum      string
}

func addSpaces(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}

	var result string
	for i, c := range s {
		if i > 0 && (n-i)%3 == 0 {
			result += " "
		}
		result += string(c)
	}

	return result
}

func formatNumber(num float64) string {
	numStr := fmt.Sprintf("%0.2f", num)
	parts := strings.Split(numStr, ".")
	formattedInteger := addSpaces(parts[0])
	formattedNumber := formattedInteger + "." + parts[1]

	return formattedNumber
}

func (a app) showHistoryReceiptsPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fromSearchDate := r.FormValue("fromSearch")
	toSearchDate := r.FormValue("toSearch")
	fromTimeStr := r.FormValue("fromTime")
	toTimeStr := r.FormValue("toTime")

	// Парсинг даты
	fromSearchDateTime, fromErr := time.Parse("2006-01-02 15:04", fromSearchDate+" "+fromTimeStr)
	toSearchDateTime, toErr := time.Parse("2006-01-02 15:04", toSearchDate+" "+toTimeStr)

	if fromErr != nil || toErr != nil {
		http.Error(rw, "Error parsing dates or times", http.StatusBadRequest)
		return
	}

	a.historyReceiptsPage(rw, r, p, fromSearchDateTime, toSearchDateTime)
}

func (a app) historyReceiptsPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params, fromSearchTime, toSearchTime time.Time) {

	id_azs, ok := getIntVal(r.FormValue("id_azs"))

	if !ok {
		http.Error(rw, "Ошибка id_azs"+r.FormValue("id_azs"), http.StatusBadRequest)
		return
	}

	loc := time.Now().Location()

	fromTime := time.Date(fromSearchTime.Year(), fromSearchTime.Month(), fromSearchTime.Day(), fromSearchTime.Hour(), fromSearchTime.Minute(), fromSearchTime.Second(), 0, loc)
	toTime := time.Date(toSearchTime.Year(), toSearchTime.Month(), toSearchTime.Day(), toSearchTime.Hour(), toSearchTime.Minute(), toSearchTime.Second(), 0, loc)

	receipts, err := a.repo.GetReceiptInRange(a.ctx, id_azs, fromTime.Unix(), toTime.Unix())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	azs, err := a.repo.GetAzs(a.ctx, id_azs)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	totalSum := 0.0
	for _, receipt := range receipts {
		// Convert Sum from string to float64
		sum, err := strconv.ParseFloat(receipt.Sum, 64)
		if err != nil {
			fmt.Println("Error parsing Sum:", err)
			continue
		}
		totalSum += sum
	}

	azsReceiptDatas := AzsReceiptTemplate{
		Azs:           azs,
		FormSearchVal: fromSearchTime.Format("2006-01-02"),
		ToSearchVal:   toSearchTime.Format("2006-01-02"),
		FromTimeVal:   fromSearchTime.Format("15:04"),
		ToTimeVal:     toSearchTime.Format("15:04"),
		Receipts:      receipts,
		Count:         len(receipts),
		TotalSum:      formatNumber(totalSum),
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

func (a app) userPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params, u repository.User) {

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
