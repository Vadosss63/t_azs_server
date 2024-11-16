package azs_statistics_controller

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/Vadosss63/t-azs/internal/application"
	"github.com/Vadosss63/t-azs/internal/repository/azs_statistics"
	"github.com/julienschmidt/httprouter"
)

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

type AzsStatisticsTemplate struct {
	IdAzs                string
	FormSearchVal        string
	ToSearchVal          string
	FromTimeVal          string
	ToTimeVal            string
	Statistics           []azs_statistics.Statistics
	TotalCash            string
	TotalCashless        string
	TotalOnline          string
	TotalLitersCol1      string
	TotalLitersCol2      string
	TotalFuelArrivalCol1 string
	TotalFuelArrivalCol2 string
}

type AzsStatisticsController struct {
	app *application.App
}

func NewController(app *application.App) *AzsStatisticsController {
	return &AzsStatisticsController{app: app}
}

func (c AzsStatisticsController) CheckDB() error {

	azsList, err := c.app.Repo.AzsRepo.GetAll(c.app.Ctx)

	if err != nil {
		log.Fatalf("Failed to get azs list: %v", err)
		return err
	}

	for i := 0; i < len(azsList); i++ {
		err = c.app.Repo.AzsStatRepo.CreateStatisticsTable(c.app.Ctx, azsList[i].IdAzs)
		if err != nil {
			log.Fatalf("Failed to create table for azs %d: %v", azsList[i].IdAzs, err)

		}
	}

	return nil
}

func (c AzsStatisticsController) Routes(router *httprouter.Router) {

	router.GET("/azs/statistics", c.app.Authorized(func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		now := time.Now()
		loc := now.Location()
		paymentType := ""

		fromSearchDateTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		toSearchDateTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)

		c.azsStatisticsPage(rw, r, p, fromSearchDateTime, toSearchDateTime, paymentType)
	}))

	router.POST("/azs/statistics", c.app.Authorized(c.showAzsStatisticsPage))

	router.POST("/azs/statistics/add", c.app.Authorized(c.azsStatistics))

}

func (c AzsStatisticsController) azsStatistics(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, ok := application.GetIntVal(strings.TrimSpace(r.FormValue("id")))
	statisticsJson := strings.TrimSpace(r.FormValue("statistics"))

	if !ok || statisticsJson == "" {
		application.SendJsonResponse(rw, http.StatusBadRequest, "All fields must be filled!", "Error")
		return
	}

	statistics, err := azs_statistics.ParseStatisticsFromJson(statisticsJson)
	if err != nil {
		application.SendJsonResponse(rw, http.StatusBadRequest, err.Error(), "Error")
		return
	}

	err = c.app.Repo.AzsStatRepo.AddStatistics(c.app.Ctx, id, statistics)
	if err != nil {
		application.SendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")
		return
	}

	application.SendJsonResponse(rw, http.StatusOK, "Ok", "Ok")
}

func (c AzsStatisticsController) showAzsStatisticsPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fromSearchDate := r.FormValue("fromSearch")
	toSearchDate := r.FormValue("toSearch")
	fromTimeStr := r.FormValue("fromTime")
	toTimeStr := r.FormValue("toTime")
	paymentType := r.FormValue("paymentType")

	fromSearchDateTime, fromErr := time.Parse("2006-01-02 15:04", fromSearchDate+" "+fromTimeStr)
	toSearchDateTime, toErr := time.Parse("2006-01-02 15:04", toSearchDate+" "+toTimeStr)

	if fromErr != nil || toErr != nil {
		http.Error(rw, "Error parsing dates or times", http.StatusBadRequest)
		return
	}

	c.azsStatisticsPage(rw, r, p, fromSearchDateTime, toSearchDateTime, paymentType)
}

func (c AzsStatisticsController) azsStatisticsPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params, fromSearchTime, toSearchTime time.Time, paymentType string) {
	idAzs, ok := application.GetIntVal(r.FormValue("id_azs"))
	if !ok {
		http.Error(rw, "Invalid id_azs value", http.StatusBadRequest)
		return
	}

	loc := time.Now().Location()
	fromTime := time.Date(fromSearchTime.Year(), fromSearchTime.Month(), fromSearchTime.Day(), fromSearchTime.Hour(), fromSearchTime.Minute(), 0, 0, loc)
	toTime := time.Date(toSearchTime.Year(), toSearchTime.Month(), toSearchTime.Day(), toSearchTime.Hour(), toSearchTime.Minute(), 0, 0, loc)

	filterParams := azs_statistics.StatisticsFilterParams{
		StartTime: fromTime.Unix(),
		EndTime:   toTime.Unix(),
	}

	statistics, err := c.app.Repo.AzsStatRepo.GetFilteredStatistics(c.app.Ctx, idAzs, filterParams)
	if err != nil {
		http.Error(rw, "Failed to retrieve filtered statistics: "+err.Error(), http.StatusInternalServerError)
		return
	}

	totalCash := 0.0
	totalCashless := 0.0
	totalOnline := 0.0
	totalLitersCol1 := 0.0
	totalLitersCol2 := 0.0
	totalFuelArrivalCol1 := 0.0
	totalFuelArrivalCol2 := 0.0

	for _, stat := range statistics {
		totalCash += float64(stat.DailyCash)
		totalCashless += float64(stat.DailyCashless)
		totalOnline += float64(stat.DailyOnline)
		totalLitersCol1 += float64(stat.DailyLitersCol1)
		totalLitersCol2 += float64(stat.DailyLitersCol2)
		totalFuelArrivalCol1 += float64(stat.FuelArrivalCol1)
		totalFuelArrivalCol2 += float64(stat.FuelArrivalCol2)
	}

	azsStatisticsData := AzsStatisticsTemplate{
		IdAzs:                fmt.Sprintf("%d", idAzs),
		FormSearchVal:        fromSearchTime.Format("2006-01-02"),
		ToSearchVal:          toSearchTime.Format("2006-01-02"),
		FromTimeVal:          fromSearchTime.Format("15:04"),
		ToTimeVal:            toSearchTime.Format("15:04"),
		Statistics:           statistics,
		TotalCash:            formatNumber(totalCash),
		TotalCashless:        formatNumber(totalCashless),
		TotalOnline:          formatNumber(totalOnline),
		TotalLitersCol1:      formatNumber(totalLitersCol1),
		TotalLitersCol2:      formatNumber(totalLitersCol2),
		TotalFuelArrivalCol1: formatNumber(totalFuelArrivalCol1),
		TotalFuelArrivalCol2: formatNumber(totalFuelArrivalCol2),
	}

	lp := filepath.Join("public", "html", "azs_statistics.html")
	navi := filepath.Join("public", "html", "user_navi.html")
	tmpl, err := template.ParseFiles(lp, navi)
	if err != nil {
		http.Error(rw, "Failed to parse template files: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tmpl.ExecuteTemplate(rw, "AzsStatisticsTemplate", azsStatisticsData); err != nil {
		http.Error(rw, "Failed to execute template: "+err.Error(), http.StatusInternalServerError)
	}
}
