package ya_controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vadosss63/t-azs/internal/application"
	"github.com/julienschmidt/httprouter"
)

func TestGetPriceListHandler(t *testing.T) {
	// Create a new YaController instance
	app := &application.App{}
	controller := NewController(app)

	// Create a new httprouter.Router instance
	router := httprouter.New()

	// Register the GetPriceListHandler route
	controller.Routes(router)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/tanker/price?apikey=expected_api_key", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new HTTP response recorder
	rr := httptest.NewRecorder()

	// Serve the HTTP request using the router
	router.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// TODO: Add more assertions to validate the response body
}
