package application

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func (a app) savePointHandler(w http.ResponseWriter, r *http.Request, par httprouter.Params) {
	var p Point

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.repo.UpdateYaAzsInfoLocation(a.ctx, requestData.IdAzs, requestData.IsEnabled)

	if err != nil {
		http.Error(w, "Ошибка обновления", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})

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
