package ya_controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Vadosss63/t-azs/internal/application"
	"github.com/Vadosss63/t-azs/internal/repository/azs"
	"github.com/Vadosss63/t-azs/internal/repository/ya_azs"

	"github.com/julienschmidt/httprouter"
)

func getYaPayApiKey() string {
	return "expected_api_key"
}

func checkAPIKey(rw http.ResponseWriter, r *http.Request) bool {
	apiKey := r.URL.Query().Get("apikey")

	if apiKey != "expected_api_key" {
		http.Error(rw, "Invalid API key", http.StatusUnauthorized)
		return false
	}
	return true
}

func sendJsonData(rw http.ResponseWriter, data any) bool {
	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(data); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}

func convertTypeFuelToYaFormat(s string) string {
	switch s {
	case "ДТ":
		return "diesel"
	case "АИ-92":
		return "a92"
	case "АИ-95":
		return "a95"
	case "АИ-98":
		return "a98"
	case "АИ-100":
		return "a100"
	case "Метан":
		return "metan"
	case "Пропан":
		return "propan"
	default:
		return ""
	}
}

type YaController struct {
	app *application.App
}

func NewController(app *application.App) *YaController {
	return &YaController{app: app}
}

func (c YaController) Routes(router *httprouter.Router) {
	router.POST("/update_yandexpay_status", c.app.Authorized(c.UpdateYandexPayStatusHandler))

	router.GET("/tanker/station", c.GetStationsHandler)

	router.GET("/tanker/price", c.GetPriceListHandler)

	router.GET("/tanker/ping", c.PingHandler)

	router.POST("/tanker/order", c.UpdateOrderStatusHandler)
	router.POST("/api/azs_order", c.GetOrderHandler)
	router.POST("/api/azs_cancel_order", c.CanceledHandler)
	router.POST("/api/azs_accept_order", c.AcceptOrderHandler)
	router.POST("/api/azs_fueling_order", c.FuelingHandler)
	router.POST("/api/azs_completed_order", c.CompletedHandler)

}

func (c YaController) GetPriceListHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// http://127.0.0.1:8086/tanker/price?apikey=expected_api_key
	if !checkAPIKey(rw, r) {
		return
	}

	azsIds, err := c.app.Repo.YaAzsRepo.GetEnableList(c.app.Ctx)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	var samplePrices = []ya_azs.PriceEntry{}
	for i := 0; i < len(azsIds); i++ {
		azsStats, err := c.app.Repo.AzsRepo.Get(c.app.Ctx, azsIds[i])
		if err != nil {
			application.SendError(rw, "Server error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		azsStatsDataFull, err := azs.ParseStats(azsStats)
		if err != nil {
			application.SendError(rw, "Server error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		for j := 0; j < len(azsStatsDataFull.AzsNodes); j++ {
			typeFuel := convertTypeFuelToYaFormat(azsStatsDataFull.AzsNodes[j].TypeFuel)
			if typeFuel == "" {
				continue
			}

			price := fmt.Sprintf("%.2f", azsStatsDataFull.AzsNodes[j].Price)

			priceFloat, _ := strconv.ParseFloat(price, 64)

			var samplePrice = ya_azs.PriceEntry{
				StationId: strconv.Itoa(azsIds[i]),
				ProductId: typeFuel,
				Price:     priceFloat,
			}

			samplePrices = append(samplePrices, samplePrice)
		}

	}

	sendJsonData(rw, samplePrices)
}

func (c YaController) GetStationsHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// http://127.0.0.1:8086/tanker/station?apikey=expected_api_key
	if !checkAPIKey(rw, r) {
		return
	}

	stations, err := c.app.Repo.YaAzsRepo.GetEnableAll(c.app.Ctx)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	for i := 0; i < len(stations); i++ {

		idInt, _ := application.GetIntVal(stations[i].Id)
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

		stations[i].Address = azsStatsDataFull.Address
		stations[i].Name = azsStatsDataFull.Name

		if stations[i].Columns == nil {
			stations[i].Columns = make(map[int32]ya_azs.Column)
		}

		for j := 0; j < len(azsStatsDataFull.AzsNodes); j++ {
			typeFuel := convertTypeFuelToYaFormat(azsStatsDataFull.AzsNodes[j].TypeFuel)
			if typeFuel == "" {
				continue
			}
			column := stations[i].Columns[int32(j)]

			if column.Fuels == nil {
				column.Fuels = []string{}
			}
			column.Fuels = append(column.Fuels, typeFuel)
			stations[i].Columns[int32(j)] = column
		}
	}

	sendJsonData(rw, stations)
}

func (c YaController) PingHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// http://127.0.0.1:8086/tanker/ping?apikey=expected_api_key&stationId=11111111&columnId=0
	if !checkAPIKey(w, r) {
		return
	}

	stationId := r.URL.Query().Get("stationId")
	// columnId := r.URL.Query().Get("columnId")

	idInt, err := strconv.Atoi(stationId)
	if err != nil {
		http.Error(w, "Bad Request: Invalid station ID", http.StatusBadRequest)
		return
	}

	//c.app.Repo.CreateYaPayTable(c.app.Ctx)

	// columnIDInt, err := strconv.Atoi(columnId)
	// if err != nil {
	// 	http.Error(w, "Bad Request: Invalid column ID", http.StatusBadRequest)
	// 	return
	// }

	enable, err := c.app.Repo.YaAzsRepo.GetEnable(c.app.Ctx, idInt)
	if err != nil || !enable {
		http.Error(w, "Not Found: Station not found", http.StatusNotFound)
		return
	}

	// // Проверка доступности станции
	// if !station.Active {
	// 	http.Error(w, "Service Unavailable: Station is not active", http.StatusServiceUnavailable)
	// 	return
	// }

	// if active, exists := station.Columns[columnIDInt]; !exists || !active {
	// 	http.Error(w, "Conflict: Column not found or not ready", http.StatusConflict)
	// 	return
	// }

	w.WriteHeader(http.StatusOK)
}

func (c YaController) UpdateOrderStatusHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var order ya_azs.Order
	err := decoder.Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stationID, _ := strconv.Atoi(order.StationId)
	enable, err := c.app.Repo.YaAzsRepo.GetEnable(c.app.Ctx, stationID)
	if err != nil || !enable {
		http.Error(w, "Not Found: Station not found", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Order with ID %s updated to status %s", order.Id, order.Status)
}

func (c YaController) UpdateYandexPayStatusHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		IdAzs     int  `json:"idAzs"`
		IsEnabled bool `json:"isEnabled"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}

	err = c.app.Repo.YaAzsRepo.UpdateEnable(c.app.Ctx, requestData.IdAzs, requestData.IsEnabled)

	if err != nil {
		http.Error(w, "Ошибка обновления", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	sendJsonData(w, map[string]string{"status": "success"})
}

const baseURL = "https://app.tanker.yandex.net"

func sendOrderStatus(endpoint string, params url.Values) error {
	fullURL := baseURL + endpoint
	resp, err := http.Post(fullURL, "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyText, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("received non-OK response: %v, body: %s", resp.Status, string(bodyText))
	}

	return nil
}

func handleAccept(apiKey, orderID string) error {
	params := url.Values{}
	params.Add("apikey", apiKey)
	params.Add("orderId", orderID)

	return sendOrderStatus("/api/order/accept", params)
}

func handleFueling(apiKey, orderID string) error {
	params := url.Values{}
	params.Add("apikey", apiKey)
	params.Add("orderId", orderID)

	return sendOrderStatus("/api/order/fueling", params)
}

func handleCanceled(apiKey, orderID, reason string) error {
	params := url.Values{}
	params.Add("apikey", apiKey)
	params.Add("orderId", orderID)
	params.Add("reason", reason)

	return sendOrderStatus("/api/order/canceled", params)
}

func handleCompleted(apiKey, orderID string, litre float64, extendedOrderID, extendedDate string) error {
	params := url.Values{}
	params.Add("apikey", apiKey)
	params.Add("orderId", orderID)
	params.Add("litre", fmt.Sprintf("%.2f", litre)) // Конвертация float в string
	params.Add("extendedOrderId", extendedOrderID)
	params.Add("extendedDate", extendedDate)

	return sendOrderStatus("/api/order/completed", params)
}

func (c YaController) GetOrderHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var requestData struct {
		IdAzs int    `json:"Id"`
		Token string `json:"Token"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		application.SendJsonResponse(rw, http.StatusBadRequest, "Incorrect data", "Error")
		return
	}

	if !c.app.ValidateToken(requestData.Token) {
		application.SendJsonResponse(rw, http.StatusUnauthorized, "Invalid token", "Error")
		return
	}

	yaPayData, err := c.app.Repo.YaPayRepo.Get(c.app.Ctx, requestData.IdAzs)

	if err != nil {
		application.SendJsonResponse(rw, http.StatusInternalServerError, "Error", "Error")
		return
	}

	application.SendJsonResponse(rw, http.StatusOK, yaPayData.Data, "Success")
}

func (c YaController) CanceledHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	var requestData struct {
		IdAzs   int    `json:"Id"`
		Token   string `json:"Token"`
		OrderId string `json:"OrderId"`
		Reason  string `json:"Reason"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		application.SendJsonResponse(rw, http.StatusBadRequest, "Incorrect data", "Error")
		return
	}

	if !c.app.ValidateToken(requestData.Token) {
		application.SendJsonResponse(rw, http.StatusUnauthorized, "Invalid token", "Error")
		return
	}

	apiKey := getYaPayApiKey()

	err = handleCanceled(apiKey, requestData.OrderId, requestData.Reason)
	if err != nil {
		application.SendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")
		return
	}

	application.SendJsonResponse(rw, http.StatusOK, "Status updated", "Success")
}

func (c YaController) AcceptOrderHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var requestData struct {
		IdAzs   int    `json:"Id"`
		Token   string `json:"Token"`
		OrderId string `json:"OrderId"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		application.SendJsonResponse(rw, http.StatusBadRequest, "Incorrect data", "Error")
		return
	}

	if !c.app.ValidateToken(requestData.Token) {
		application.SendJsonResponse(rw, http.StatusUnauthorized, "Invalid token", "Error")
		return
	}

	apiKey := getYaPayApiKey()

	err = handleAccept(apiKey, requestData.OrderId)
	if err != nil {
		application.SendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")
		return
	}

	application.SendJsonResponse(rw, http.StatusOK, "Status updated", "Success")
}

func (c YaController) FuelingHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var requestData struct {
		IdAzs   int    `json:"Id"`
		Token   string `json:"Token"`
		OrderId string `json:"OrderId"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		application.SendJsonResponse(rw, http.StatusBadRequest, "Incorrect data", "Error")
		return
	}

	if !c.app.ValidateToken(requestData.Token) {
		application.SendJsonResponse(rw, http.StatusUnauthorized, "Invalid token", "Error")
		return
	}

	apiKey := getYaPayApiKey()

	err = handleFueling(apiKey, requestData.OrderId)

	if err != nil {
		application.SendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")
		return
	}
	application.SendJsonResponse(rw, http.StatusOK, "Status updated", "Success")
}

func (c YaController) CompletedHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var requestData struct {
		IdAzs           int     `json:"Id"`
		Token           string  `json:"Token"`
		OrderId         string  `json:"OrderId"`
		Litre           float64 `json:"Litre"`
		ExtendedOrderId string  `json:"ExtendedOrderId"`
		ExtendedDate    string  `json:"ExtendedDate"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		application.SendJsonResponse(rw, http.StatusBadRequest, "Incorrect data", "Error")
		return
	}

	if !c.app.ValidateToken(requestData.Token) {
		application.SendJsonResponse(rw, http.StatusUnauthorized, "Invalid token", "Error")
		return
	}

	apiKey := getYaPayApiKey()

	err = handleCompleted(apiKey, requestData.OrderId, requestData.Litre, requestData.ExtendedOrderId, requestData.ExtendedDate)
	if err != nil {
		application.SendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")
		return
	}

	application.SendJsonResponse(rw, http.StatusOK, "Status updated", "Success")
}
