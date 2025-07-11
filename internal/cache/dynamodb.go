package cache

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sirupsen/logrus"
	"time"
	"weather-service/internal/logging"
)

//go:generate mockgen --source=dynamodb.go --destination mocks/dynamodb.go --package mocks

type DynamoDBClient interface {
	GetItem(ctx context.Context, input *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, input *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

type DynamoDBCache struct {
	client     DynamoDBClient
	tableName  string
	ttlMinutes int
}

func NewDynamoDBCache(client DynamoDBClient, tableName string, ttl int) *DynamoDBCache {
	return &DynamoDBCache{
		client:     client,
		tableName:  tableName,
		ttlMinutes: ttl,
	}
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

	if key == "" {
		return nil, fmt.Errorf("empty key provided")
	}

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
