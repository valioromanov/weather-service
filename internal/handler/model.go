package handler

type WeatherServiceResponse struct {
	Date            string  `json:"date"`
	Latitude        string  `json:"latitude"`
	Longitude       string  `json:"longitude"`
	Temperature     float64 `json:"temperature"`
	UVIndex         float64 `json:"uvIndex"`
	RainProbability float64 `json:"rainProbability"`
}
