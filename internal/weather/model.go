package weather

type Daily struct {
	Time                        []string  `json:"time"`
	Temperature2mMax            []float64 `json:"temperature_2m_max"`
	UVIndexMax                  []float64 `json:"uv_index_max"`
	PrecipitationProbabilityMax []float64 `json:"precipitation_probability_max"`
}

type OpenMeteoResponse struct {
	Daily Daily `json:"daily"`
}
type Forecast struct {
	Temp2max          float64 `json:"temperature_2m_max"`
	UvIndexMax        float64 `json:"uv_index_max"`
	PrecipProbability float64 `json:"precipitation_probability_max"`
}

type ForecastMap map[string]Forecast
