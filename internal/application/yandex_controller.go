package application

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Station struct {
	Id       string `json:"Id"`
	Enable   bool   `json:"Enable"`
	Name     string `json:"Name"`
	Address  string `json:"Address"`
	Location struct {
		Lat float64 `json:"Lat"`
		Lon float64 `json:"Lon"`
	} `json:"Location"`
	Columns map[int32][]string `json:"Columns"`
}

type Order struct {
	Id                string    `json:"Id"`
	DateCreate        time.Time `json:"DateCreate"`
	OrderType         string    `json:"OrderType"`
	OrderVolume       float64   `json:"OrderVolume"`
	StationId         string    `json:"StationId"`
	StationExtendedId string    `json:"StationExtendedId"`
	ColumnId          int       `json:"ColumnId"`
	FuelId            string    `json:"FuelId"`
	FuelMarka         string    `json:"FuelMarka"`
	PriceId           string    `json:"PriceId"`
	FuelExtendedId    string    `json:"FuelExtendedId"`
	PriceFuel         float64   `json:"PriceFuel"`
	Sum               float64   `json:"Sum"`
	Litre             float64   `json:"Litre"`
	SumPaid           float64   `json:"SumPaid"`
	Status            string    `json:"Status"`
	DateEnd           time.Time `json:"DateEnd"`
	ReasonId          string    `json:"ReasonId"`
	Reason            string    `json:"Reason"`
	LitreCompleted    float64   `json:"LitreCompleted"`
	SumPaidCompleted  float64   `json:"SumPaidCompleted"`
	ContractId        string    `json:"ContractId"`
}

type PriceEntry struct {
	StationId string  `json:"StationId"`
	ProductId string  `json:"ProductId"`
	Price     float64 `json:"Price"`
}

// Заглушка данных для демонстрации ответа
var samplePrices = []PriceEntry{
	{StationId: "0001", ProductId: "a92", Price: 38.66},
	{StationId: "0001", ProductId: "a95_premium", Price: 45.21},
	{StationId: "0002", ProductId: "a92", Price: 38.98},
}

// Функция обработчика для получения прайс-листа
func getPriceListHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	apiKey := r.URL.Query().Get("apikey")

	// Проверка API ключа
	if apiKey != "expected_api_key" {
		http.Error(w, "Invalid API key", http.StatusUnauthorized)
		return
	}

	// Заголовок ответа
	w.Header().Set("Content-Type", "application/json")

	// Отправка данных в JSON формате
	if err := json.NewEncoder(w).Encode(samplePrices); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Функция для получения списка АЗС
func getStationsHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// В реальном приложении здесь будет логика получения данных
	stations := []Station{
		{
			Id:      "1",
			Enable:  true,
			Name:    "Station 1",
			Address: "123 Main St",
			Location: struct {
				Lat float64 `json:"Lat"`
				Lon float64 `json:"Lon"`
			}{Lat: 55.7558, Lon: 37.6173},
			Columns: map[int32][]string{
				1: {"a92", "a95"},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stations)
}

// Функция для обновления статуса заказа
func updateOrderStatusHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var order Order
	err := decoder.Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Здесь должна быть логика обработки заказа

	fmt.Fprintf(w, "Order with ID %s updated to status %s", order.Id, order.Status)
}

//PING handler

// Модель данных станции и ТРК для демонстрации
type StationStatus struct {
	ID      string
	Active  bool
	Columns map[int]bool // Ключ - ID колонки, значение - активность колонки
}

// Заглушка данных о станциях
var stations = map[string]StationStatus{
	"station1": {
		ID:     "station1",
		Active: true,
		Columns: map[int]bool{
			1: true,
			2: false, // Эта колонка не активна
		},
	},
	"station2": {
		ID:     "station2",
		Active: false, // Станция не активна
		Columns: map[int]bool{
			1: true,
			2: true,
		},
	},
}

// Функция обработчика для проверки станции и колонки
func pingHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	apiKey := r.URL.Query().Get("apikey")
	stationId := r.URL.Query().Get("stationId")
	columnId := r.URL.Query().Get("columnId")

	// Проверка API ключа
	if apiKey != "expected_api_key" {
		http.Error(w, "Unauthorized: Invalid API key", http.StatusUnauthorized)
		return
	}

	// Поиск станции по ID
	station, ok := stations[stationId]
	if !ok {
		http.Error(w, "Not Found: Station not found", http.StatusNotFound)
		return
	}

	// Проверка доступности станции
	if !station.Active {
		http.Error(w, "Service Unavailable: Station is not active", http.StatusServiceUnavailable)
		return
	}

	// Проверка наличия и состояния колонки
	columnIDInt, err := strconv.Atoi(columnId)
	if err != nil {
		http.Error(w, "Bad Request: Invalid column ID", http.StatusBadRequest)
		return
	}

	if active, exists := station.Columns[columnIDInt]; !exists || !active {
		http.Error(w, "Conflict: Column not found or not ready", http.StatusConflict)
		return
	}

	// Все проверки пройдены, станция и колонка готовы к использованию
	w.WriteHeader(http.StatusOK)
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

// func main() {
// 	apiKey := "your_api_key"
// 	orderID := "123456"
// 	reason := "Customer request"
// 	litre := 50.0
// 	extendedOrderID := "ext123456"
// 	extendedDate := "01.01.2020 12:00:00"

// 	if err := handleAccept(apiKey, orderID); err != nil {
// 		log.Printf("Failed to send accept status: %v", err)
// 	} else {
// 		fmt.Println("Accept status sent successfully.")
// 	}

// 	if err := handleFueling(apiKey, orderID); err != nil {
// 		log.Printf("Failed to send fueling status: %v", err)
// 	} else {
// 		fmt.Println("Fueling status sent successfully.")
// 	}

// 	if err := handleCanceled(apiKey, orderID, reason); err != nil {
// 		log.Printf("Failed to send canceled status: %v", err)
// 	} else {
// 		fmt.Println("Canceled status sent successfully.")
// 	}

// 	if err := handleCompleted(apiKey, orderID, litre, extendedOrderID, extendedDate); err != nil {
// 		log.Printf("Failed to send completed status: %v", err)
// 	} else {
// 		fmt.Println("Completed status sent successfully.")
// 	}
// }
