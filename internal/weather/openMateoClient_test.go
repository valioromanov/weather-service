package weather_test

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
	"net/http"
	"weather-service/helper/mockutil"
	"weather-service/internal/weather"
	"weather-service/internal/weather/mocks"
)

var _ = Describe("OpenMateoClient", mockutil.Mockable(func(helper *mockutil.Helper) {

	var (
		mockHTTPClient *mocks.MockHttpRequester
		omc            *weather.OpenMateoClient
	)

	BeforeEach(func() {
		mockHTTPClient = mocks.NewMockHttpRequester(helper.Controller())
		omc = weather.NewOpenMateoClient(mockHTTPClient, "testurl.com/latitude=%s&longitude=%s")
	})

	Context("Get", func() {
		When("everything works", func() {
			BeforeEach(func() {
				response := "{\"latitude\":43.0,\"longitude\":23.0,\"daily\":{\"time\":[\"2025-07-10\"],\"temperature_2m_max\":[20.8],\"uv_index_max\":[5.3],\"precipitation_probability_max\":[0]}}"
				mockHTTPClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString(response)),
				}, nil).Times(1)
			})

			It("should return weather data", func() {
				resp, err := omc.GetForecast("43.0", "23.0")
				Expect(err).ToNot(HaveOccurred())
				Expect(len(resp)).To(Equal(1))
				Expect(resp["2025-07-10"].Latitude).To(Equal("43.0000"))
				Expect(resp["2025-07-10"].Longitude).To(Equal("23.0000"))
				Expect(resp["2025-07-10"].Temp2max).To(Equal(20.8))
				Expect(resp["2025-07-10"].UvIndexMax).To(Equal(5.3))
				Expect(resp["2025-07-10"].PrecipProbability).To(Equal(float64(0)))
			})
		})

		When("request fails", func() {
			BeforeEach(func() {
				mockHTTPClient.EXPECT().Do(gomock.Any()).Return(&http.Response{}, errors.New("error")).Times(1)
			})

			It("should return error", func() {
				resp, err := omc.GetForecast("43.0", "23.0")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("error"))
				Expect(resp).To(BeNil())
			})
		})
	})
}))
