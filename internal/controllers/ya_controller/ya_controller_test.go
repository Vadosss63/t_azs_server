package ya_controller

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vadosss63/t-azs/internal/repository"
)

func TestGetPriceListHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/tanker/price-list?apikey=expected_api_key", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		dbpool, err := repository.InitDBConn(ctx)
		if err != nil {
			log.Fatalf("Failed to init DB connection: %v", err)
		}
		defer dbpool.Close()

		a := NewApp(ctx, dbpool, "", 8086)
		a.getPriceListHandler(w, r, nil)
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// expected := `[{"StationId":"10111920","ProductId":"diesel","Price":50.24},{"StationId":"11111111","ProductId":"a95","Price":49.27},{"StationId":"11111111","ProductId":"a95","Price":51.01}]`
	// resp := rr.Body.String()
	//
	//	if resp != expected {
	//		t.Errorf("handler returned unexpected body: \ngot \n%v \nwant \n%v", resp, expected)
	//	}
}
