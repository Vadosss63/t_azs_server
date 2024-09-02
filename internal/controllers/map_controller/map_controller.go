package map_controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Vadosss63/t-azs/internal/application"
	"github.com/Vadosss63/t-azs/internal/repository/ya_azs"
	"github.com/julienschmidt/httprouter"
)

type PointData struct {
	IdAzs int     `json:"id_azs"`
	Lat   float64 `json:"lat"`
	Lng   float64 `json:"lng"`
}

type MapController struct {
	app *application.App
}

func NewController(app *application.App) *MapController {
	return &MapController{app: app}
}

func (c MapController) Routes(router *httprouter.Router) {
	router.POST("/save-point", c.app.Authorized(c.savePointHandler))
	router.GET("/points", c.app.Authorized(c.pointsHandler))
}

func (c MapController) savePointHandler(w http.ResponseWriter, r *http.Request, par httprouter.Params) {
	var p PointData

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.app.Repo.YaAzsRepo.UpdateLocation(c.app.Ctx, p.IdAzs, ya_azs.Location{Lat: p.Lat, Lon: p.Lng})

	if err != nil {
		http.Error(w, "Ошибка обновления", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})
}

func (c MapController) pointsHandler(w http.ResponseWriter, r *http.Request, par httprouter.Params) {

	id := strings.TrimSpace(r.FormValue("id_azs"))
	id_azs, ok := application.GetIntVal(id)

	if !ok {
		http.Error(w, "Error user", http.StatusBadRequest)
		return
	}

	point, _ := c.app.Repo.YaAzsRepo.GetLocation(c.app.Ctx, id_azs)

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(PointData{IdAzs: id_azs, Lat: point.Lat, Lng: point.Lon})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
