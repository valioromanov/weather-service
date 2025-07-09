package handler

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
	"time"
	"weather-service/internal/weather"
)

type WeatherService struct {
	WeatherClient weather.Client
}

func NewWeatherService(wc weather.Client) *WeatherService {
	return &WeatherService{
		WeatherClient: wc,
	}
}

func (wsvc *WeatherService) HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	lat := req.QueryStringParameters["lat"]
	lon := req.QueryStringParameters["lon"]
	date := req.QueryStringParameters["date"]

	if lat == "" || lon == "" || date == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Missing lat/lon/date"}, nil
	}

	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid date"}, nil
	}
	
	if parsedDate.Before(time.Now()) {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid date: Date could not be older than today"}, nil
	}

	sevenDaysLater := time.Now().Add(7 * 24 * time.Hour)
	if parsedDate.After(sevenDaysLater) {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid date: Date could not be 7 day from today"}, nil
	}

	forecastRes, err := wsvc.WeatherClient.GetForecast(lat, lon)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Weather api error"}, nil
	}
	if _, ok := forecastRes[date]; !ok {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Weather forecast not found for this date"}, nil
	}
	wsr := WeatherServiceResponse{date,
		forecastRes[date].Temp2max,
		forecastRes[date].UvIndexMax,
		forecastRes[date].PrecipProbability,
	}

	wsrBytes, err := json.Marshal(wsr)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Weather forecast not found for this date"}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(wsrBytes),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}
