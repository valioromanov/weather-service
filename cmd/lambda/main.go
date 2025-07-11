package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/sirupsen/logrus"
	"net/http"
	"weather-service/cmd/env"
	"weather-service/internal/cache"
	"weather-service/internal/handler"
	"weather-service/internal/weather"
)

func main() {

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	//Loading env vars
	appConfig, err := env.LoadAppConfig()
	if err != nil {
		logrus.Fatal(err)
	}
	//TODO Add url in env var
	//weatherClient := &weather.OpenMateoClient{Url: "https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&daily=temperature_2m_max,uv_index_max,precipitation_probability_max&timezone=auto"}
	// Initializing weather client
	httpClient := &http.Client{}
	weatherClient := weather.NewOpenMateoClient(httpClient, appConfig.OpenMateoURL)

	// Loading AWS config
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("eu-west-1"))
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create dynamoDB client")
	}

	// Initializing DynamoDB Client
	dynamoDBClient := dynamodb.NewFromConfig(cfg)

	// Initializing Cache
	weatherCache := cache.NewDynamoDBCache(dynamoDBClient, appConfig.DynamoDBName, appConfig.TTL)

	// Initializing handler
	service := handler.NewWeatherService(weatherClient, weatherCache)

	logrus.Info("Starting Weather api Lambda")
	lambda.Start(service.HandleRequest)
}
