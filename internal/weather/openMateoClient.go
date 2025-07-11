package weather

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"weather-service/internal/logging"
)

//go:generate mockgen --source=openMateoClient.go --destination mocks/openMateoClient.go --package mocks

type HttpRequester interface {
	Do(req *http.Request) (*http.Response, error)
}

type OpenMateoClient struct {
	HttpClient HttpRequester
	Url        string //"https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&daily=temperature_2m_max,uv_index_max,precipitation_probability_max&timezone=auto"
}

func NewOpenMateoClient(hc HttpRequester, url string) *OpenMateoClient {
	return &OpenMateoClient{
		HttpClient: hc,
		Url:        url,
	}
}

func (c *OpenMateoClient) GetForecast(lat, long string) (ForecastMap, error) {
	logrus.WithFields(logrus.Fields{
		"lat":  lat,
		"long": long,
	}).Info("Going to get forecast from OpenMateo")

	url := fmt.Sprintf(c.Url, lat, long)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var opr OpenMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&opr); err != nil {
		logging.LogError(err, map[string]interface{}{"lat": lat, "long": long})
		return nil, err
	}
	fm := make(ForecastMap)
	for i := 0; i < len(opr.Daily.Time); i++ {
		fm[opr.Daily.Time[i]] = Forecast{
			Latitude:          fmt.Sprintf("%.4f", opr.Latitude),
			Longitude:         fmt.Sprintf("%.4f", opr.Longitude),
			Temp2max:          opr.Daily.Temperature2mMax[i],
			UvIndexMax:        opr.Daily.UVIndexMax[i],
			PrecipProbability: opr.Daily.PrecipitationProbabilityMax[i],
		}
	}
	return fm, nil
}
