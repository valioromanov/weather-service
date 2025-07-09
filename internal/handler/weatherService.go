package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
	"time"
	"weather-service/internal/cache"
	"weather-service/internal/weather"
)

type WeatherService struct {
	WeatherClient weather.Client
	WeatherCache  cache.Cache
}

func NewWeatherService(clnt weather.Client, wc cache.Cache) *WeatherService {
	return &WeatherService{
		WeatherClient: clnt,
		WeatherCache:  wc,
	}
}

func (wsvc *WeatherService) HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	lat := req.QueryStringParameters["lat"]
	lon := req.QueryStringParameters["lon"]
	date := req.QueryStringParameters["date"]

	if lat == "" || lon == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Missing lat/lon/date"}, nil
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
		return respond(CachedDataToWeatherServiceResponse(*cachedWeather))
	}

	forecastRes, err := wsvc.WeatherClient.GetForecast(lat, lon)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Weather api error"}, nil
	}
	if _, ok := forecastRes[date]; !ok {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Weather forecast not found for this date"}, nil
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
		_ = wsvc.WeatherCache.Put(keyStore, data)
	}
}

func respond(w WeatherServiceResponse) (events.APIGatewayProxyResponse, error) {
	wsrBytes, err := json.Marshal(w)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Weather forecast not found for this date"}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(wsrBytes),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}
