package cache

type Cache interface {
	Put(key string, weather *CachedWeather) error
	Get(key string) (*CachedWeather, error)
}
