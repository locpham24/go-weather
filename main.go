package main

import (
	"github.com/gin-gonic/gin"
	"github.com/locpham24/go-weather/handler"
)

func main() {
	r := gin.Default()
	handler.InitRouter(r)
	r.Run()
}
