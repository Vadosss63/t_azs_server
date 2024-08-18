package application

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
)

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

var points []Point
var mu sync.Mutex

func savePointHandler(w http.ResponseWriter, r *http.Request, par httprouter.Params) {
	var p Point

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	points = []Point{p}
	mu.Unlock()

	w.WriteHeader(http.StatusOK)
}

func pointsHandler(w http.ResponseWriter, r *http.Request, par httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	mu.Lock()
	err := json.NewEncoder(w).Encode(points)
	mu.Unlock()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
