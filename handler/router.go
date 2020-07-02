package handler

import (
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/locpham24/go-weather/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/locpham24/go-weather/db"
)

type OpenWeatherMapData struct {
	CityName string `json:"name"`
	Current  struct {
		Temp     float64 `json:"temp"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
		Humidity float64 `json:"humidity"`
	} `json:"main"`
	Sys struct {
		Country string `json:"country"`
		ID      int    `json:"id"`
	} `json:"sys"`
	Coord struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lon"`
	} `json:"coord"`
}

type WeatherHandler struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

type LocationHandler struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

func InitRouter(pg *db.PgDb, redis *redis.Client) *gin.Engine {
	locHandler := LocationHandler{
		DB:          pg.DB,
		RedisClient: redis,
	}
	weatherHandler := WeatherHandler{
		DB:          pg.DB,
		RedisClient: redis,
	}
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.GET("/temperature/:city", weatherHandler.getWeather)
	r.POST("/location", locHandler.createLocation)
	r.GET("/location/:city", locHandler.getOneLocation)
	return r
}

func (w *WeatherHandler) getWeather(c *gin.Context) {
	city := c.Param("city")
	data := OpenWeatherMapData{}

	weatherInfo, err := w.RedisClient.Get(city).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if weatherInfo != "" {
		err := json.Unmarshal([]byte(weatherInfo), &data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		data, err = GetWeatherData(city)
		jsonData, err := json.Marshal(data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		w.RedisClient.Set(city, string(jsonData), 30*time.Second)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if data == (OpenWeatherMapData{}) {
		c.JSON(http.StatusNotFound, gin.H{"error": errors.New("No data")})
	}

	var location models.Location
	w.DB.Where("name = ?", city).First(&location)
	if location.Name != "" {
		w.DB.Model(&location).Updates(models.Location{
			Latitude:  data.Coord.Lat,
			Longitude: data.Coord.Lng,
			Count:     location.Count + 1,
		})
	} else {
		loc := models.Location{
			Latitude:  data.Coord.Lat,
			Longitude: data.Coord.Lng,
			Name:      city,
		}

		w.DB.Create(&loc)
	}

	c.JSON(http.StatusOK, data)
}
func GetWeatherData(cityId string) (OpenWeatherMapData, error) {
	data := OpenWeatherMapData{}

	API := "a52f1f3ec1f62e44817ceb3e22ce2e03"
	res, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + cityId + "&units=metric&appid=" + API)
	if err != nil || res.StatusCode != 200 {
		return data, err
	}
	defer res.Body.Close()

	time.Sleep(5 * time.Second)
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (loc *LocationHandler) createLocation(c *gin.Context) {
	var location models.LocationReq
	if err := c.ShouldBindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m := models.Location{
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
		Name:      location.Name,
	}

	loc.DB.Create(&m)
}

func (loc *LocationHandler) getOneLocation(c *gin.Context) {
	city := c.Param("city")
	var location models.Location

	loc.DB.Where("name = ?", city).First(&location)

	c.JSON(http.StatusOK, location)
}

func (loc *LocationHandler) deleteLocation(c *gin.Context) {
	// soft delete
}
