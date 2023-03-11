package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Vadosss63/t-azs/internal/application"
	"github.com/Vadosss63/t-azs/internal/repository"
	"github.com/julienschmidt/httprouter"
)

func main() {
	ctx := context.Background()
	dbpool, err := repository.InitDBConn(ctx)
	if err != nil {
		log.Fatalf("%w failed to init DB connection", err)
	}
	defer dbpool.Close()
	a := application.NewApp(ctx, dbpool)
	r := httprouter.New()
	a.Routes(r)
	fmt.Println("It's alive! Try http://t-azs.ru")
	http.ListenAndServe(":80", r)
	//fmt.Println("It is alive! Try http://t-azs.ru:8080")
	//http.ListenAndServe(":8080", r)
	//http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", r)

}
