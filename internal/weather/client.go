package weather

type Client interface {
	GetForecast(lat, long string) (ForecastMap, error)
}
