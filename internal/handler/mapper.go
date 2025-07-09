package handler

import (
	"strings"
	"weather-service/internal/cache"
	"weather-service/internal/weather"
)

func CachedDataToWeatherServiceResponse(cachedData cache.CachedWeather) WeatherServiceResponse {
	//key = lat_lon_date
	keySplit := strings.Split(cachedData.Key, "_")

	wsr := WeatherServiceResponse{
		keySplit[2],
		keySplit[0],
		keySplit[1],
		cachedData.TempMax,
		cachedData.UVIndex,
		cachedData.RainProb,
	}
	return wsr
}

func ForecastToCachedData(forecast weather.Forecast) *cache.CachedWeather {
	return &cache.CachedWeather{
		TempMax:  forecast.Temp2max,
		UVIndex:  forecast.UvIndexMax,
		RainProb: forecast.PrecipProbability,
	}
}
