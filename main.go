package main

import (
	"github.com/go-redis/redis"
	"github.com/locpham24/go-weather/db"
	"github.com/locpham24/go-weather/handler"
)

func main() {
	pg := db.PgDb{}
	pg.Connect()

	defer pg.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	router := handler.InitRouter(&pg, redisClient)
	router.Run(":8080")
}
