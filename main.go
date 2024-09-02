package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Vadosss63/t-azs/internal/application"
	"github.com/Vadosss63/t-azs/internal/controllers/trbl_controller"
	"github.com/Vadosss63/t-azs/internal/controllers/update_app_controller"
	"github.com/Vadosss63/t-azs/internal/controllers/user_controller"
	"github.com/Vadosss63/t-azs/internal/controllers/ya_controller"
	"github.com/Vadosss63/t-azs/internal/repository"
	"github.com/julienschmidt/httprouter"
)

type Settings struct {
	Port  int    `json:"port"`
	Token string `json:"token"`
}

func readSettings(filename string) (*Settings, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var settings Settings
	err = json.Unmarshal(data, &settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil
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

	ctx := context.Background()
	dbpool, err := repository.InitDBConn(ctx)
	if err != nil {
		log.Fatalf("Failed to init DB connection: %v", err)
	}
	defer dbpool.Close()

	a := application.NewApp(ctx, dbpool, settings.Token, settings.Port)
	r := httprouter.New()

	yaController := ya_controller.NewController(a)
	updateAppController := update_app_controller.NewController(a)
	userController := user_controller.NewController(a)
	trblController := trbl_controller.NewController(a)

	a.Routes(r)
	yaController.Routes(r)
	updateAppController.Routes(r)
	userController.Routes(r)
	trblController.Routes(r)

	fmt.Printf("It's alive! Try http://t-azs.ru:%d/ or http://127.0.0.1:%d\n", settings.Port, settings.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", settings.Port), r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
