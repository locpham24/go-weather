package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/locpham24/go-weather/models"
)

type PgDb struct {
	DB *gorm.DB
}

func (pg *PgDb) Connect() {
	db, err := gorm.Open("postgres",
		"host=localhost port=5432 user=default dbname=weather password=secret sslmode=disable")
	if err != nil {
		panic("Failed to connect to database!")
	}
	pg.DB = db
}

func (pg *PgDb) Migrate() {
	pg.DB.AutoMigrate(&models.Location{})
}

func (pg *PgDb) Close() {
	pg.DB.Close()
}
