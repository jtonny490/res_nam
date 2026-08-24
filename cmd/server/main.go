package main

import (
	"log"
	"net/http"
	"res_nam/internal/app"
	"res_nam/internal/config"
)

func main() {
	c := config.Load()
	db, err := app.New(c)
	if err != nil {
		log.Fatal(err)
	}
	r := app.BuildRouter(c, db)
	log.Printf("listening on :%s", c.Port)
	log.Fatal(http.ListenAndServe(":"+c.Port, r))
}
