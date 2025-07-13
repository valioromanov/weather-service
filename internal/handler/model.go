package handler

type WeatherServiceResponse struct {
	Date            string  `json:"date"`
	Latitude        string  `json:"latitude"`
	Longitude       string  `json:"longitude"`
	Temperature     float64 `json:"temperature"`
	UVIndex         float64 `json:"uvIndex"`
	RainProbability float64 `json:"rainProbability"`
}

type Forecast struct {
	Longitude         string  `json:"longitude"`
	Latitude          string  `json:"latitude"`
	Temp2max          float64 `json:"temperature_2m_max"`
	UvIndexMax        float64 `json:"uv_index_max"`
	PrecipProbability float64 `json:"precipitation_probability_max"`
}

type ForecastMap map[string]Forecast

type CachedWeather struct {
	Key      string  `dynamodbav:"Key"`
	TempMax  float64 `dynamodbav:"TempMax"`
	UVIndex  float64 `dynamodbav:"UVIndex"`
	RainProb float64 `dynamodbav:"RainProb"`
	TTL      int64   `dynamodbav:"TTL"`
}
