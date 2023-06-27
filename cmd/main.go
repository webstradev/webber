package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/westradev/webbr/api"
	"github.com/westradev/webbr/webbr"
)

func main() {
	db, err := webbr.New(webbr.WithDBName("main"))
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(db)

	e := echo.New()
	e.HideBanner = true

	e.POST("/:dbname/:collname", server.HandlePostInsert)
	e.GET("/:dbname/:collname", server.HandleGetQuery)
	log.Fatal(e.Start(":6969"))

}
