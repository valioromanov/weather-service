package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"weather-service/internal/cache"
	"weather-service/internal/handler"
	"weather-service/internal/weather"
)

func main() {
	weatherClient := &weather.OpenMateoClient{Url: "https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&daily=temperature_2m_max,uv_index_max,precipitation_probability_max&timezone=auto"}
	weatherCache, err := cache.NewDynamoDBCache()
	if err != nil {
		log.Fatal(err)
	}
	service := handler.NewWeatherService(weatherClient, weatherCache)

	lambda.Start(service.HandleRequest)
}
