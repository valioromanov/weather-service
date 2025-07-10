package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
	"weather-service/internal/cache"
	"weather-service/internal/logging"
	"weather-service/internal/weather"
)

//go:generate mockgen --source=weatherService.go --destination mocks/weatherService.go --package mocks

type ForecastClient interface {
	GetForecast(lat, long string) (weather.ForecastMap, error)
}

type Cache interface {
	Put(key string, weather *cache.CachedWeather) error
	Get(key string) (*cache.CachedWeather, error)
}

type WeatherService struct {
	WeatherClient ForecastClient
	WeatherCache  Cache
}

func NewWeatherService(clnt ForecastClient, wc Cache) *WeatherService {
	return &WeatherService{
		WeatherClient: clnt,
		WeatherCache:  wc,
	}
}

func (wsvc *WeatherService) HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	lat := req.QueryStringParameters["lat"]
	lon := req.QueryStringParameters["lon"]
	date := req.QueryStringParameters["date"]

	logrus.WithFields(logrus.Fields{
		"lat":  lat,
		"lon":  lon,
		"date": date,
	}).Info("Going to handle request")

	if lat == "" || lon == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Missing lat/lon"}, nil
	}

	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid date"}, nil
	}

	today := time.Now().Truncate(24 * time.Hour)
	if parsedDate.Before(today) {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid date: Date could not be older than today"}, nil
	}

	sevenDaysLater := time.Now().Add(7 * 24 * time.Hour)
	if parsedDate.After(sevenDaysLater) {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid date: Date could not be 7 day from today"}, nil
	}

	key := fmt.Sprintf("%s_%s_%s", lat, lon, date)
	if cachedWeather, err := wsvc.WeatherCache.Get(key); err == nil && cachedWeather != nil {
		logrus.WithFields(logrus.Fields{
			"key": key,
		}).Info("Got weather from cache")
		return respond(CachedDataToWeatherServiceResponse(*cachedWeather))
	}

	logrus.WithFields(logrus.Fields{
		"key": key,
	}).Info("Did not find weather from cache, will fetch from third party provider")
	forecastRes, err := wsvc.WeatherClient.GetForecast(lat, lon)
	if err != nil {
		errId := logging.LogError(err, map[string]interface{}{"lat": lat, "lon": lon})
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("[%s] Weather api error", errId)}, nil
	}
	if _, ok := forecastRes[date]; !ok {
		errId := logging.LogError(err, map[string]interface{}{"lat": lat, "lon": lon, "date": date})
		return events.APIGatewayProxyResponse{StatusCode: 404, Body: fmt.Sprintf("[%s] Weather forecast not found for this date", errId)}, nil
	}

	wsr := WeatherServiceResponse{date,
		forecastRes[date].Latitude,
		forecastRes[date].Longitude,
		forecastRes[date].Temp2max,
		forecastRes[date].UvIndexMax,
		forecastRes[date].PrecipProbability,
	}

	batchPutToCacheStore(wsvc, forecastRes)

	return respond(wsr)
}

func batchPutToCacheStore(wsvc *WeatherService, fm weather.ForecastMap) {
	for key, value := range fm {
		keyStore := fmt.Sprintf("%s_%s_%s", value.Latitude, value.Longitude, key)
		data := ForecastToCachedData(value)
		err := wsvc.WeatherCache.Put(keyStore, data)
		if err != nil {
			logging.LogError(err, map[string]interface{}{"key": key, "data": data})
		}
	}
}

func respond(w WeatherServiceResponse) (events.APIGatewayProxyResponse, error) {
	wsrBytes, err := json.Marshal(w)
	if err != nil {
		errId := logging.LogError(err, map[string]interface{}{"weatherServiceResponse": w})
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: fmt.Sprintf("[%s] Error while generating response", errId)}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(wsrBytes),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}
