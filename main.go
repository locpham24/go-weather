package main

import (
	"github.com/locpham24/go-weather/handler"
)

func main() {
	router := handler.InitRouter()
	router.Run(":8080")
}
