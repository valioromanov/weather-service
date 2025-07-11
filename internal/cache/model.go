package cache

type CachedWeather struct {
	Key      string  `dynamodbav:"Key"`
	TempMax  float64 `dynamodbav:"TempMax"`
	UVIndex  float64 `dynamodbav:"UVIndex"`
	RainProb float64 `dynamodbav:"RainProb"`
	TTL      int64   `dynamodbav:"TTL"`
}
