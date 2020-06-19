package main

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/locpham24/go-weather/db"
	"github.com/locpham24/go-weather/handler"
)

func main() {
	pg := db.PgDb{}
	pg.Connect()

	defer pg.Close()

	router := handler.InitRouter(&pg)
	router.Run(":8080")
}
