package models

type LocationReq struct {
	ID        uint    `json:"id"`
	Latitude  float64 `json:"lat" binding:"required"`
	Longitude float64 `json:"lng" binding:"required"`
	Name      string  `json:"name" binding:"required"`
}
type Location struct {
	ID        uint    `json:"id" gorm:"primary_key"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	Count     int     `json:"count" gorm:"default:0"`
}
