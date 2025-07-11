package env

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type AppConfig struct {
	OpenMateoURL string `envconfig:"OPEN_MATEO_URL"`
	DynamoDBName string `envconfig:"DYNAMODB_TABLE"`
	TTL          int    `envconfig:"TTL_MINUTES"`
}

func LoadAppConfig() (AppConfig, error) {
	var config AppConfig
	if err := envconfig.Process("", &config); err != nil {
		logrus.Error("error while binding to 'AppConfig': ", err)
		return AppConfig{}, fmt.Errorf("failed to parse configuration from environment: %w", err)
	}

	return config, nil
}
