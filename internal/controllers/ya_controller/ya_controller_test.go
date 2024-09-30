package ya_controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vadosss63/t-azs/internal/application"
	"github.com/Vadosss63/t-azs/internal/repository"
	"github.com/Vadosss63/t-azs/internal/repository/azs"
	"github.com/Vadosss63/t-azs/internal/repository/ya_azs"
	"github.com/Vadosss63/t-azs/internal/repository/ya_pay"
	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

// mockgen -source=internal/repository/ya_pay/ya_pay.go -destination=internal/repository/ya_pay/moc_ya_pay.go -package=ya_pay
func TestGetPriceListHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockYaAzsRepo := ya_azs.NewMockYaAzsRepository(ctrl)
	mockAzsRepo := azs.NewMockAzsRepository(ctrl)

	mockYaAzsRepo.EXPECT().GetEnableList(gomock.Any()).Return([]int{1}, nil)

	mockAzsData := azs.AzsStatsData{
		Id:                  1,
		IdAzs:               1,
		IdUser:              1,
		IsAuthorized:        1,
		CountColum:          2,
		IsSecondPriceEnable: 1,
		Time:                "2024-09-07T10:00:00Z",
		Name:                "AZS Example",
		Address:             "123 Main St",
		Stats:               `{"azs_nodes":[{"averageTemperature":"0.00","commonLiters":"0.00","dailyLiters":"0.00","density":"0.00","fuelArrival":0,"fuelVolume":"0","fuelVolumePerc":"0.00","lockFuelValue":60,"price":4927,"priceCashless":5522,"typeFuel":"АИ-92"},{"averageTemperature":"0.00","commonLiters":"0.00","dailyLiters":"0.00","density":"0.00","fuelArrival":0,"fuelVolume":"0","fuelVolumePerc":"0.00","lockFuelValue":100,"price":5101,"priceCashless":5200,"typeFuel":"АИ-95"}],"main_info":{"commonCash":0,"commonCashless":0,"commonOnline":0,"dailyCash":0,"dailyCashless":0,"dailyOnline":0,"isBlock":false,"version":"1.0.2"}}`,
	}

	mockAzsRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(mockAzsData, nil).AnyTimes()

	mocRepo := repository.Repository{
		YaAzsRepo: mockYaAzsRepo,
		AzsRepo:   mockAzsRepo,
	}

	app := &application.App{
		Repo:        &mocRepo,
		YaPayApiKey: "expected_api_key",
	}

	controller := NewController(app)
	req, err := http.NewRequest("GET", "/tanker/price?apikey=expected_api_key", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := httprouter.New()
	router.GET("/tanker/price", controller.GetPriceListHandler)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedJSONresponse := `[{"StationId":"1","ProductId":"a92","Price":49.27},{"StationId":"1","ProductId":"a95","Price":51.01}]`
	assert.JSONEq(t, expectedJSONresponse, rr.Body.String())
}

func TestGetStationsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockYaAzsRepo := ya_azs.NewMockYaAzsRepository(ctrl)
	mockAzsRepo := azs.NewMockAzsRepository(ctrl)

	mocStation := ya_azs.Station{Id: "1", Enable: true, Location: ya_azs.Location{Lat: 11, Lon: 12}, Columns: map[int32]ya_azs.Column{1: {Fuels: []string{"a92", "a95"}}}}
	mockYaAzsRepo.EXPECT().GetEnableAll(gomock.Any()).Return([]ya_azs.Station{mocStation}, nil)

	mockAzsData := azs.AzsStatsData{
		Id:                  1,
		IdAzs:               1,
		IdUser:              1,
		IsAuthorized:        1,
		CountColum:          2,
		IsSecondPriceEnable: 1,
		Time:                "2024-09-07T10:00:00Z",
		Name:                "AZS Example",
		Address:             "123 Main St",
		Stats:               `{"azs_nodes":[{"averageTemperature":"0.00","commonLiters":"0.00","dailyLiters":"0.00","density":"0.00","fuelArrival":0,"fuelVolume":"0","fuelVolumePerc":"0.00","lockFuelValue":60,"price":4927,"priceCashless":5522,"typeFuel":"АИ-92"},{"averageTemperature":"0.00","commonLiters":"0.00","dailyLiters":"0.00","density":"0.00","fuelArrival":0,"fuelVolume":"0","fuelVolumePerc":"0.00","lockFuelValue":100,"price":5101,"priceCashless":5200,"typeFuel":"АИ-95"}],"main_info":{"commonCash":0,"commonCashless":0,"commonOnline":0,"dailyCash":0,"dailyCashless":0,"dailyOnline":0,"isBlock":false,"version":"1.0.2"}}`,
	}

	mockAzsRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(mockAzsData, nil).AnyTimes()

	mocRepo := repository.Repository{
		YaAzsRepo: mockYaAzsRepo,
		AzsRepo:   mockAzsRepo,
	}

	app := &application.App{
		Repo:        &mocRepo,
		YaPayApiKey: "expected_api_key",
	}

	controller := NewController(app)

	req, err := http.NewRequest("GET", "/tanker/station?apikey=expected_api_key", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := httprouter.New()
	router.GET("/tanker/station", controller.GetStationsHandler)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedJSONresponse := `[{"Id":"1","Enable":true,"Name":"AZS Example","Address":"123 Main St","Location":{"Lat":11,"Lon":12},"Columns":{"0":{"Fuels":["a92"]},"1":{"Fuels":["a92","a95","a95"]}}}]`
	assert.JSONEq(t, expectedJSONresponse, rr.Body.String())
}

func TestPingHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockYaAzsRepo := ya_azs.NewMockYaAzsRepository(ctrl)
	mockAzsRepo := azs.NewMockAzsRepository(ctrl)
	mockYaPayRepo := ya_pay.NewMockYaPayRepository(ctrl)

	azsId := 11111111
	colId := 0

	mockYaAzsRepo.EXPECT().GetEnable(gomock.Any(), azsId).Return(true, nil)
	mockYaPayRepo.EXPECT().Get(gomock.Any(), azsId, colId).Return(ya_pay.YaPay{IdAzs: azsId, ColumnId: colId, Status: 0, Data: ""}, nil)

	mocRepo := repository.Repository{
		YaAzsRepo: mockYaAzsRepo,
		AzsRepo:   mockAzsRepo,
		YaPayRepo: mockYaPayRepo,
	}

	app := &application.App{
		Repo:        &mocRepo,
		YaPayApiKey: "expected_api_key",
	}

	controller := NewController(app)

	req, err := http.NewRequest("GET", "/tanker/ping?apikey=expected_api_key&stationId=11111111&columnId=0", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := httprouter.New()
	router.GET("/tanker/ping", controller.PingHandler)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateYandexPayStatusHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	azsId := 11111111

	mockYaAzsRepo := ya_azs.NewMockYaAzsRepository(ctrl)

	mockYaAzsRepo.EXPECT().UpdateEnable(gomock.Any(), azsId, true).Return(nil)

	mocRepo := repository.Repository{
		YaAzsRepo: mockYaAzsRepo,
	}

	app := &application.App{
		Repo: &mocRepo,
	}

	controller := NewController(app)

	requestData := struct {
		IdAzs     int  `json:"idAzs"`
		IsEnabled bool `json:"isEnabled"`
	}{
		IdAzs:     11111111,
		IsEnabled: true,
	}
	body, err := json.Marshal(requestData)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/update_yandexpay_status", bytes.NewBuffer(body))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := httprouter.New()
	router.POST("/update_yandexpay_status", controller.UpdateYandexPayStatusHandler)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

// func TestUpdateOrderStatusHandler(t *testing.T) {
// 	app := &application.App{}
// 	controller := NewController(app)

// 	order := ya_azs.Order{
// 		Id:        "123",
// 		StationId: "11111111",
// 		Status:    "completed",
// 	}
// 	body, err := json.Marshal(order)
// 	assert.NoError(t, err)

// 	req, err := http.NewRequest("POST", "/tanker/order", bytes.NewBuffer(body))
// 	assert.NoError(t, err)
// 	req.Header.Set("Content-Type", "application/json")

// 	rr := httptest.NewRecorder()
// 	router := httprouter.New()
// 	router.POST("/tanker/order", controller.UpdateOrderStatusHandler)
// 	router.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusOK, rr.Code)
// }

//{"Id": "9DA356FB-3483-4FD4-B62C-7B85A81D003D", "DateCreate": "2023-08-23T12:26:51+03:00", "OrderType": "Liters", "OrderVolume": 2.0, "StationExtendedId": "190011", "ColumnId": 1, "FuelExtendedId": "a92", "PriceFuel": 50.0, "Sum": 100.0, "SumPaid": 100.0, "Litre": 2.0, "Status": "OrderCreated", "DateEnd": null, "ReasonId": null, "Reason": null, "LitreCompleted": null, "SumPaidCompleted": null}
