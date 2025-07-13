package weather

type Daily struct {
	Time                        []string  `json:"time"`
	Temperature2mMax            []float64 `json:"temperature_2m_max"`
	UVIndexMax                  []float64 `json:"uv_index_max"`
	PrecipitationProbabilityMax []float64 `json:"precipitation_probability_max"`
}

type OpenMeteoResponse struct {
	Daily     Daily   `json:"daily"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
