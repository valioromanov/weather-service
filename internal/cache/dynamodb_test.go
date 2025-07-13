package cache_test

import (
	"errors"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"weather-service/helper/mockutil"
	"weather-service/internal/cache/mocks"
	"weather-service/internal/handler"

	"weather-service/internal/cache"
)

var _ = Describe("Dynamodb", mockutil.Mockable(func(helper *mockutil.Helper) {
	var (
		mockDynamoDBClient *mocks.MockDynamoDBClient
		dynamoDBClient     *cache.DynamoDBCache
	)

	BeforeEach(func() {
		mockDynamoDBClient = mocks.NewMockDynamoDBClient(helper.Controller())
		dynamoDBClient = cache.NewDynamoDBCache(mockDynamoDBClient, "WeatherCache", 10)
	})

	Context("GetItem", func() {
		When("everything works", func() {
			var item handler.CachedWeather
			BeforeEach(func() {
				item = handler.CachedWeather{
					Key:      "42.0_23.0_2025-07-10",
					TempMax:  30.5,
					UVIndex:  7.8,
					RainProb: 40.0,
					TTL:      123621653216,
				}
				av, err := attributevalue.MarshalMap(item)
				Expect(err).To(BeNil())
				mockDynamoDBClient.EXPECT().GetItem(gomock.Any(), gomock.Any()).Return(&dynamodb.GetItemOutput{
					Item: av,
				}, nil).Times(1)
			})

			It("should return the item", func() {
				res, err := dynamoDBClient.Get("42.0_23.0_2025-07-10")
				Expect(err).To(BeNil())
				Expect(res.RainProb).To(Equal(item.RainProb))
				Expect(res.TTL).To(Equal(item.TTL))
				Expect(res.Key).To(Equal("42.0_23.0_2025-07-10"))
				Expect(res.TempMax).To(Equal(item.TempMax))
				Expect(res.UVIndex).To(Equal(item.UVIndex))
			})
		})

		When("dynamodb returns an error", func() {
			BeforeEach(func() {
				mockDynamoDBClient.EXPECT().GetItem(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).Times(1)
			})

			It("should return the error", func() {
				_, err := dynamoDBClient.Get("42.0_23.0_2025-07-10")
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("error"))
			})
		})

		When("dynamodb returns nil item", func() {
			BeforeEach(func() {
				mockDynamoDBClient.EXPECT().GetItem(gomock.Any(), gomock.Any()).Return(&dynamodb.GetItemOutput{
					Item: nil,
				}, nil).Times(1)
			})

			It("should return nil", func() {
				res, err := dynamoDBClient.Get("42.0_23.0_2025-07-10")
				Expect(err).To(BeNil())
				Expect(res).To(BeNil())
			})
		})

		When("no key is provided", func() {
			BeforeEach(func() {
				mockDynamoDBClient.EXPECT().GetItem(gomock.Any(), gomock.Any()).Return(nil, nil).Times(0)
			})

			It("should return error", func() {
				_, err := dynamoDBClient.Get("")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("empty key provided"))
			})
		})
	})

	Context("PutItem", func() {
		cachedWeather := &handler.CachedWeather{
			TempMax:  30.5,
			UVIndex:  7.8,
			RainProb: 40.0,
		}
		When("everything works", func() {
			BeforeEach(func() {
				mockDynamoDBClient.EXPECT().PutItem(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
			})

			It("should return no error", func() {
				err := dynamoDBClient.Put("43.0_23.9_2025_07_10", cachedWeather)
				Expect(err).ToNot(HaveOccurred())
				Expect(err).To(BeNil())
			})
		})

		When("dynamodb returns an error", func() {
			BeforeEach(func() {
				mockDynamoDBClient.EXPECT().PutItem(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).Times(1)
			})
			It("should return the error", func() {
				err := dynamoDBClient.Put("43.0_23.9_2025_07_10", cachedWeather)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("error"))
			})
		})
	})

}))
