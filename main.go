package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Vadosss63/t-azs/internal/application"
	"github.com/Vadosss63/t-azs/internal/controllers/admin_controller"
	"github.com/Vadosss63/t-azs/internal/controllers/azs_controller"
	"github.com/Vadosss63/t-azs/internal/controllers/customer_controller"
	"github.com/Vadosss63/t-azs/internal/controllers/map_controller"
	"github.com/Vadosss63/t-azs/internal/controllers/trbl_controller"
	"github.com/Vadosss63/t-azs/internal/controllers/update_app_controller"
	"github.com/Vadosss63/t-azs/internal/controllers/user_controller"
	"github.com/Vadosss63/t-azs/internal/controllers/ya_controller"
	"github.com/Vadosss63/t-azs/internal/repository"
	"github.com/julienschmidt/httprouter"
)

func readSettings(filename string) (application.Settings, error) {
	var s application.Settings

	data, err := os.ReadFile(filename)
	if err != nil {
		return s, err
	}

	err = json.Unmarshal(data, &s)
	if err != nil {
		return s, err
	}

	return s, nil
}

func main() {

	settings, err := readSettings("settings.json")

	if err != nil {
		log.Fatalf("Failed to read settings: %v", err)
	}

	if len(settings.Token) != 32 {
		log.Fatalf("Invalid token length: got %d, want 32", len(settings.Token))
	}

	fmt.Println("Port:", settings.Port)
	fmt.Println("Token:", settings.Token)
	fmt.Println("ya_api_key:", settings.YaApiKey)

	ctx := context.Background()
	dbpool, err := repository.InitDBConn(ctx)
	if err != nil {
		log.Fatalf("Failed to init DB connection: %v", err)
	}
	defer dbpool.Close()

	a := application.NewApp(ctx, dbpool, settings)
	r := httprouter.New()

	yaController := ya_controller.NewController(a)
	updateAppController := update_app_controller.NewController(a)
	userController := user_controller.NewController(a)
	trblController := trbl_controller.NewController(a)
	mapController := map_controller.NewController(a)
	customerController := customer_controller.NewController(a)
	azsController := azs_controller.NewController(a)
	adminController := admin_controller.NewController(a)

	a.Routes(r)
	yaController.Routes(r)
	updateAppController.Routes(r)
	userController.Routes(r)
	trblController.Routes(r)
	mapController.Routes(r)
	customerController.Routes(r)
	azsController.Routes(r)
	adminController.Routes(r)

	r.GET("/", a.Authorized(func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

		userId, ok := r.Context().Value("userId").(int)

		if !ok {
			http.Error(rw, "Error user", http.StatusBadRequest)
			return
		}
		u, err := a.Repo.UserRepo.Get(a.Ctx, userId)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		if u.Login == "admin" {
			adminController.AdminPage(rw, r, p, u, -1)
			return
		}

		customerController.UserPage(rw, r, p, u)
	}))

	fmt.Printf("It's alive! Try http://t-azs.ru:%d/ or http://127.0.0.1:%d\n", settings.Port, settings.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", settings.Port), r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
