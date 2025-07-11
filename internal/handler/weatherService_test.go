package handler_test

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
	"weather-service/helper/mockutil"
	"weather-service/internal/cache"
	"weather-service/internal/handler"
	"weather-service/internal/handler/mocks"
	"weather-service/internal/weather"
)

var _ = Describe("WeatherService", mockutil.Mockable(func(helper *mockutil.Helper) {
	var (
		mockCache          *mocks.MockCache
		mockForecastClient *mocks.MockForecastClient
		ws                 *handler.WeatherService
	)

	BeforeEach(func() {
		mockCache = mocks.NewMockCache(helper.Controller())
		mockForecastClient = mocks.NewMockForecastClient(helper.Controller())
		ws = handler.NewWeatherService(mockForecastClient, mockCache)
	})

	Context("Right query params", func() {
		today := time.Now().Format("2006-01-02")
		When("cache does not return data", func() {
			expectedRes := weather.ForecastMap{
				today: weather.Forecast{
					Latitude:          "42.0",
					Longitude:         "23.0",
					Temp2max:          23,
					UvIndexMax:        3,
					PrecipProbability: 0,
				},
			}
			BeforeEach(func() {
				mockCache.EXPECT().Get(gomock.Any()).Return(nil, nil).Times(1)
				mockForecastClient.EXPECT().GetForecast(gomock.Any(), gomock.Any()).Return(expectedRes, nil).Times(1)
				mockCache.EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			})

			It("should return date from forecast client", func() {

				req := events.APIGatewayProxyRequest{
					QueryStringParameters: map[string]string{
						"lat":  "42.0",
						"lon":  "23.0",
						"date": today,
					},
				}
				res, err := ws.HandleRequest(context.TODO(), req)
				Expect(err).ToNot(HaveOccurred())
				Expect(res.StatusCode).To(Equal(200))
				Expect(res.Body).To(Equal(fmt.Sprintf("{\"date\":\"%s\",\"latitude\":\"42.0\",\"longitude\":\"23.0\",\"temperature\":23,\"uvIndex\":3,\"rainProbability\":0}", today)))
			})
		})
		When("cache return data", func() {
			BeforeEach(func() {
				key := fmt.Sprintf("42.0_23.0_%s", today)
				expectedCachedResult := &cache.CachedWeather{
					Key:      key,
					TempMax:  23.0,
					UVIndex:  3,
					RainProb: 0,
					TTL:      1233312,
				}

				mockCache.EXPECT().Get(key).Return(expectedCachedResult, nil).Times(1)
				mockForecastClient.EXPECT().GetForecast(gomock.Any(), gomock.Any()).Return(nil, nil).Times(0)
				mockCache.EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil).Times(0)
			})

			It("should return date from cache client", func() {
				req := events.APIGatewayProxyRequest{
					QueryStringParameters: map[string]string{
						"lat":  "42.0",
						"lon":  "23.0",
						"date": today,
					},
				}
				res, err := ws.HandleRequest(context.TODO(), req)
				Expect(err).ToNot(HaveOccurred())
				Expect(res.StatusCode).To(Equal(200))
				Expect(res.Body).To(Equal(fmt.Sprintf("{\"date\":\"%s\",\"latitude\":\"42.0\",\"longitude\":\"23.0\",\"temperature\":23,\"uvIndex\":3,\"rainProbability\":0}", today)))
			})
		})

		When("cache return no data and forecast client returns error", func() {
			BeforeEach(func() {
				mockCache.EXPECT().Get(gomock.Any()).Return(nil, nil).Times(1)
				mockForecastClient.EXPECT().GetForecast(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("some error")).Times(1)
				mockCache.EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil).Times(0)
			})

			It("should return error response", func() {
				req := events.APIGatewayProxyRequest{
					QueryStringParameters: map[string]string{
						"lat":  "42.0",
						"lon":  "23.0",
						"date": today,
					},
				}
				res, err := ws.HandleRequest(context.TODO(), req)
				Expect(err).ToNot(HaveOccurred())
				Expect(res.StatusCode).To(Equal(500))
				Expect(res.Body).To(ContainSubstring("Weather api error"))
			})
		})

		When("forecast not found for this date", func() {
			BeforeEach(func() {
				resFromClient := weather.ForecastMap{
					"2025-07-10": weather.Forecast{
						Latitude:          "42.0",
						Longitude:         "23.0",
						Temp2max:          23,
						UvIndexMax:        3,
						PrecipProbability: 0,
					},
				}
				mockCache.EXPECT().Get(gomock.Any()).Return(nil, nil).Times(1)
				mockForecastClient.EXPECT().GetForecast(gomock.Any(), gomock.Any()).Return(resFromClient, nil).Times(1)
				mockCache.EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil).Times(0)
			})

			It("should return error response", func() {
				req := events.APIGatewayProxyRequest{
					QueryStringParameters: map[string]string{
						"lat":  "42.0",
						"lon":  "23.0",
						"date": "2025-07-11",
					},
				}

				resp, err := ws.HandleRequest(context.TODO(), req)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(404))
				Expect(resp.Body).To(ContainSubstring("Weather forecast not found for this date"))
			})
		})
	})
	Context("Wrong query params", func() {
		When("latitude or longitude is not provided", func() {
			It("should return error response", func() {
				req := events.APIGatewayProxyRequest{
					QueryStringParameters: map[string]string{
						"lon":  "23.0",
						"date": "2025-07-09",
					},
				}
				resp, err := ws.HandleRequest(context.TODO(), req)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(400))
				Expect(resp.Body).To(ContainSubstring("Missing lat/lon"))
			})
		})

		When("previous date provided", func() {
			It("should return error response", func() {
				req := events.APIGatewayProxyRequest{
					QueryStringParameters: map[string]string{
						"lat":  "42.0",
						"lon":  "23.0",
						"date": "2025-07-09",
					},
				}
				resp, err := ws.HandleRequest(context.TODO(), req)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(400))
				Expect(resp.Body).To(ContainSubstring("Invalid date: Date could not be older than today"))
			})
		})

		When("late date provided", func() {
			It("should return error response", func() {
				req := events.APIGatewayProxyRequest{
					QueryStringParameters: map[string]string{
						"lat":  "42.0",
						"lon":  "23.0",
						"date": "2025-07-28",
					},
				}
				resp, err := ws.HandleRequest(context.TODO(), req)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(400))
				Expect(resp.Body).To(ContainSubstring("Invalid date: Date could not be 7 day from today"))
			})
		})
	})
}))
