package application

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Vadosss63/t-azs/internal/repository/ya_azs"
	"github.com/julienschmidt/httprouter"
)

type PointData struct {
	IdAzs int     `json:"id_azs"`
	Lat   float64 `json:"lat"`
	Lng   float64 `json:"lng"`
}

func (a app) savePointHandler(w http.ResponseWriter, r *http.Request, par httprouter.Params) {
	var p PointData

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.repo.YaAzsRepo.UpdateYaAzsInfoLocation(a.ctx, p.IdAzs, ya_azs.Location{Lat: p.Lat, Lon: p.Lng})

	if err != nil {
		http.Error(w, "Ошибка обновления", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})
}

func (a app) pointsHandler(w http.ResponseWriter, r *http.Request, par httprouter.Params) {

	id := strings.TrimSpace(r.FormValue("id_azs"))
	id_azs, ok := getIntVal(id)

	if !ok {
		http.Error(w, "Error user", http.StatusBadRequest)
		return
	}

	point, err := a.repo.YaAzsRepo.GetYaAzsInfoLocation(a.ctx, id_azs)

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(PointData{IdAzs: id_azs, Lat: point.Lat, Lng: point.Lon})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
