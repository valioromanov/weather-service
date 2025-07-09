package handler

type WeatherServiceResponse struct {
	Date            string  `json:"date"`
	Temperature     float64 `json:"temperature"`
	UVIndex         float64 `json:"uvIndex"`
	RainProbability float64 `json:"rainProbability"`
}
