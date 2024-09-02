package application

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/Vadosss63/t-azs/internal/repository/azs"
	"github.com/Vadosss63/t-azs/internal/repository/receipt"
	"github.com/Vadosss63/t-azs/internal/repository/user"
	"github.com/julienschmidt/httprouter"
)

type AzsStatsTemplate struct {
	User  user.User
	Azses []azs.AzsStatsDataFull
}

type AzsReceiptTemplate struct {
	Azs             azs.AzsStatsData
	FormSearchVal   string
	ToSearchVal     string
	FromTimeVal     string
	ToTimeVal       string
	Receipts        []receipt.Receipt
	Count           int
	TotalSum        string
	TotalLiters     string
	FormPaymentType string
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
	paymentType := r.FormValue("paymentType")

	// Парсинг даты
	fromSearchDateTime, fromErr := time.Parse("2006-01-02 15:04", fromSearchDate+" "+fromTimeStr)
	toSearchDateTime, toErr := time.Parse("2006-01-02 15:04", toSearchDate+" "+toTimeStr)

	if fromErr != nil || toErr != nil {
		http.Error(rw, "Error parsing dates or times", http.StatusBadRequest)
		return
	}

	a.historyReceiptsPage(rw, r, p, fromSearchDateTime, toSearchDateTime, paymentType)
}

func (a app) historyReceiptsPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params, fromSearchTime, toSearchTime time.Time, paymentType string) {
	id_azs, ok := getIntVal(r.FormValue("id_azs"))
	if !ok {
		http.Error(rw, "Invalid id_azs value", http.StatusBadRequest)
		return
	}

	loc := time.Now().Location()
	fromTime := time.Date(fromSearchTime.Year(), fromSearchTime.Month(), fromSearchTime.Day(), fromSearchTime.Hour(), fromSearchTime.Minute(), 0, 0, loc)
	toTime := time.Date(toSearchTime.Year(), toSearchTime.Month(), toSearchTime.Day(), toSearchTime.Hour(), toSearchTime.Minute(), 0, 0, loc)

	filterParams := receipt.FilterParams{
		StartTime:   fromTime.Unix(),
		EndTime:     toTime.Unix(),
		PaymentType: paymentType,
	}

	receipts, err := a.repo.ReceiptRepo.GetFilteredReceipts(a.ctx, id_azs, filterParams)
	if err != nil {
		http.Error(rw, "Failed to retrieve filtered receipts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	azs, err := a.repo.AzsRepo.Get(a.ctx, id_azs)
	if err != nil {
		http.Error(rw, "Failed to retrieve AZS data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	totalSum := 0.0
	totalLiters := 0.0
	for _, receipt := range receipts {
		switch paymentType {
		case "cash":
			totalSum += float64(receipt.Cash)
		case "cashless":
			totalSum += float64(receipt.Cashless)
		case "online":
			totalSum += float64(receipt.Online)
		default:
			totalSum += float64(receipt.Sum)
		}

		totalLiters += float64(receipt.CountLitres)
	}

	azsReceiptDatas := AzsReceiptTemplate{
		Azs:             azs,
		FormSearchVal:   fromSearchTime.Format("2006-01-02"),
		ToSearchVal:     toSearchTime.Format("2006-01-02"),
		FromTimeVal:     fromSearchTime.Format("15:04"),
		ToTimeVal:       toSearchTime.Format("15:04"),
		Receipts:        receipts,
		Count:           len(receipts),
		TotalSum:        formatNumber(totalSum),
		TotalLiters:     formatNumber(totalLiters),
		FormPaymentType: paymentType, // Установка выбранного типа оплаты
	}

	lp := filepath.Join("public", "html", "azs_receipt.html")
	navi := filepath.Join("public", "html", "user_navi.html")
	tmpl, err := template.ParseFiles(lp, navi)
	if err != nil {
		http.Error(rw, "Failed to parse template files: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tmpl.ExecuteTemplate(rw, "AzsReceiptTemplate", azsReceiptDatas); err != nil {
		http.Error(rw, "Failed to execute template: "+err.Error(), http.StatusInternalServerError)
	}
}

func (a app) userPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params, u user.User) {

	azs_statses, err := a.repo.AzsRepo.GetAzsAllForUser(a.ctx, u.Id)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	azsStatsTemplate := AzsStatsTemplate{
		User:  u,
		Azses: []azs.AzsStatsDataFull{},
	}

	for _, azs_stats := range azs_statses {

		azsStatsDataFull, err := azs.ParseStats(azs_stats)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		azsStatsTemplate.Azses = append(azsStatsTemplate.Azses, azsStatsDataFull)
	}

	lp := filepath.Join("public", "html", "azs_stats.html")
	navi := filepath.Join("public", "html", "user_navi.html")
	azs := filepath.Join("public", "html", "azs_container.html")
	tmpl := template.Must(template.ParseFiles(lp, navi, azs))

	err = tmpl.ExecuteTemplate(rw, "AzsStatsTemplate", azsStatsTemplate)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}
