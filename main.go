package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Vadosss63/t-azs/internal/application"
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
		panic(err)
	}

	if len(settings.Token) != 32 {
		panic(fmt.Sprintf("Error Token = %s", settings.Token))
	}

	fmt.Println("Port:", settings.Port)
	fmt.Println("Token:", settings.Token)

	ctx := context.Background()
	dbpool, err := repository.InitDBConn(ctx)
	if err != nil {
		log.Fatalf("%w failed to init DB connection", err)
	}
	defer dbpool.Close()
	a := application.NewApp(ctx, dbpool, settings.Token, settings.Port)
	r := httprouter.New()
	a.Routes(r)
	fmt.Printf("It's alive! Try http://t-azs.ru:%d/\n", settings.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", settings.Port), r)
	//http.ListenAndServeTLS(fmt.Sprintf(":%d", settings.Port), "cert.pem", "key.pem", r)
}
