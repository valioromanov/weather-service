package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"weather-service/internal/handler"
	"weather-service/internal/weather"
)

func main() {
	weatherClient := &weather.OpenMateoClient{Url: "https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&daily=temperature_2m_max,uv_index_max,precipitation_probability_max&timezone=auto"}
	service := handler.NewWeatherService(weatherClient)

	lambda.Start(service.HandleRequest)
}
