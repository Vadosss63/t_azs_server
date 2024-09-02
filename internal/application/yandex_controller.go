package application

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Vadosss63/t-azs/internal/repository/azs"
	"github.com/Vadosss63/t-azs/internal/repository/ya_azs"

	"github.com/julienschmidt/httprouter"
)

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

func (a app) getPriceListHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// http://127.0.0.1:8086/tanker/price?apikey=expected_api_key
	if !checkAPIKey(rw, r) {
		return
	}

	azsIds, err := a.repo.YaAzsRepo.GetYaAzsInfoEnableList(a.ctx)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	var samplePrices = []ya_azs.PriceEntry{}
	for i := 0; i < len(azsIds); i++ {
		azsStats, err := a.repo.AzsRepo.GetAzs(a.ctx, azsIds[i])
		if err != nil {
			sendError(rw, "Server error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		azsStatsDataFull, err := azs.ParseStats(azsStats)
		if err != nil {
			sendError(rw, "Server error: "+err.Error(), http.StatusInternalServerError)
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

func (a app) getStationsHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// http://127.0.0.1:8086/tanker/station?apikey=expected_api_key
	if !checkAPIKey(rw, r) {
		return
	}

	stations, err := a.repo.YaAzsRepo.GetYaAzsInfoAllEnable(a.ctx)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	for i := 0; i < len(stations); i++ {

		idInt, _ := getIntVal(stations[i].Id)
		azsStats, err := a.repo.AzsRepo.GetAzs(a.ctx, idInt)
		if err != nil {
			sendError(rw, "Server error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		azsStatsDataFull, err := azs.ParseStats(azsStats)
		if err != nil {
			sendError(rw, "Server error: "+err.Error(), http.StatusInternalServerError)
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

func (a app) pingHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

	//a.repo.CreateYaPayTable(a.ctx)

	// columnIDInt, err := strconv.Atoi(columnId)
	// if err != nil {
	// 	http.Error(w, "Bad Request: Invalid column ID", http.StatusBadRequest)
	// 	return
	// }

	enable, err := a.repo.YaAzsRepo.GetYaAzsInfoEnable(a.ctx, idInt)
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

// Функция для обновления статуса заказа
func (a app) updateOrderStatusHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var order ya_azs.Order
	err := decoder.Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stationID, _ := strconv.Atoi(order.StationId)
	enable, err := a.repo.YaAzsRepo.GetYaAzsInfoEnable(a.ctx, stationID)
	if err != nil || !enable {
		http.Error(w, "Not Found: Station not found", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Order with ID %s updated to status %s", order.Id, order.Status)
}

func (a app) updateYandexPayStatusHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

	err = a.repo.YaAzsRepo.UpdateYaAzsInfoEnable(a.ctx, requestData.IdAzs, requestData.IsEnabled)

	if err != nil {
		http.Error(w, "Ошибка обновления", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	sendJsonData(w, map[string]string{"status": "success"})
}

const baseURL = "https://app.tanker.yandex.net" // Константа базового URL

// Функция отправки статуса заказа с использованием данных формы
func sendOrderStatus(endpoint string, params url.Values) error {
	fullURL := baseURL + endpoint // Полный URL для запроса
	resp, err := http.Post(fullURL, "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Проверка ответа от Яндекс.Заправки
	if resp.StatusCode != http.StatusOK {
		bodyText, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("received non-OK response: %v, body: %s", resp.Status, string(bodyText))
	}

	return nil
}

// Обработчик для статуса Accept
func handleAccept(apiKey, orderID string) error {
	params := url.Values{}
	params.Add("apikey", apiKey)
	params.Add("orderId", orderID)

	return sendOrderStatus("/api/order/accept", params)
}

// Обработчик для статуса Fueling
func handleFueling(apiKey, orderID string) error {
	params := url.Values{}
	params.Add("apikey", apiKey)
	params.Add("orderId", orderID)

	return sendOrderStatus("/api/order/fueling", params)
}

// Обработчик для статуса Canceled
func handleCanceled(apiKey, orderID, reason string) error {
	params := url.Values{}
	params.Add("apikey", apiKey)
	params.Add("orderId", orderID)
	params.Add("reason", reason)

	return sendOrderStatus("/api/order/canceled", params)
}

// Обработчик для статуса Completed
func handleCompleted(apiKey, orderID string, litre float64, extendedOrderID, extendedDate string) error {
	params := url.Values{}
	params.Add("apikey", apiKey)
	params.Add("orderId", orderID)
	params.Add("litre", fmt.Sprintf("%.2f", litre)) // Конвертация float в string
	params.Add("extendedOrderId", extendedOrderID)
	params.Add("extendedDate", extendedDate)

	return sendOrderStatus("/api/order/completed", params)
}
