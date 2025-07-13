package handler

import (
	"strings"
)

func CachedDataToWeatherServiceResponse(cachedData CachedWeather) WeatherServiceResponse {
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

func ForecastToCachedData(forecast Forecast) *CachedWeather {
	return &CachedWeather{
		TempMax:  forecast.Temp2max,
		UVIndex:  forecast.UvIndexMax,
		RainProb: forecast.PrecipProbability,
	}
}
