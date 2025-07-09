package cache

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sirupsen/logrus"
	"time"
	"weather-service/internal/logging"
)

type DynamoDBCache struct {
	client     *dynamodb.Client
	tableName  string
	ttlMinutes int
}

func NewDynamoDBCache() (*DynamoDBCache, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-west-1"))
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)
	return &DynamoDBCache{
		client:     client,
		tableName:  "WeatherCache",
		ttlMinutes: 10,
	}, nil
}

func (c *DynamoDBCache) Put(key string, weather *CachedWeather) error {
	weather.Key = key
	weather.TTL = time.Now().Add(time.Duration(c.ttlMinutes) * time.Minute).Unix()

	item, err := attributevalue.MarshalMap(weather)
	if err != nil {
		return err
	}

	_, err = c.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(c.tableName),
		Item:      item,
	})

	return err
}

func (c *DynamoDBCache) Get(key string) (*CachedWeather, error) {
	logrus.WithFields(logrus.Fields{
		"key": key,
	}).Info("Going to get a weather from cache")

	resp, err := c.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(c.tableName),
		Key: map[string]types.AttributeValue{
			"Key": &types.AttributeValueMemberS{Value: key},
		},
	})
	if err != nil {
		logging.LogError(fmt.Errorf("error while getting weather from cache: %w", err), map[string]interface{}{"key": key})
		return nil, err
	}

	if resp.Item == nil {
		return nil, nil // not found, and we will make a request if nothing was found so we do not need an error
	}

	var data CachedWeather
	err = attributevalue.UnmarshalMap(resp.Item, &data)
	if err != nil {
		logging.LogError(fmt.Errorf("error while unmarshaling weather from cache: %w", err), map[string]interface{}{"key": key})
		return nil, err
	}

	// check for expire, if so ignore
	if data.TTL < time.Now().Unix() {
		return nil, nil
	}

	return &data, nil
}
