package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
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
	} `json:"sys"`
}

func InitRouter(r *gin.Engine) {
	r.GET("/temperature/:city", func(c *gin.Context) {
		city := c.Param("city")
		data, err := GetWeatherData(city)
		if err != nil {
			c.JSON(200, nil)
		}
		c.JSON(200, data)
	})
}

func GetWeatherData(cityId string) (OpenWeatherMapData, error) {
	data := OpenWeatherMapData{}

	API := "a52f1f3ec1f62e44817ceb3e22ce2e03"
	res, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q="+cityId+"&units=metric&appid=" + API)
	if err != nil || res.StatusCode != 200 {
		return data, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}