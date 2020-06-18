package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	r := InitRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestGetTemperature(t *testing.T) {
	r := InitRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/temperature/saigon", nil)
	r.ServeHTTP(w, req)

	data := OpenWeatherMapData{}
	err := json.NewDecoder(w.Body).Decode(&data)
	if err != nil {
		t.Errorf("Convert JSON fail")
	}
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "9314", strconv.Itoa(data.Sys.ID))
	fmt.Println(w.Body)

}
