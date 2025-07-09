package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type OpenMateoClient struct {
	Url string //"https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&daily=temperature_2m_max,uv_index_max,precipitation_probability_max&timezone=auto"
}

func (c *OpenMateoClient) GetForecast(lat, long string) (ForecastMap, error) {
	url := fmt.Sprintf(c.Url, lat, long)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var opr OpenMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&opr); err != nil {
		return nil, err
	}
	fm := make(ForecastMap)
	for i := 0; i < len(opr.Daily.Time); i++ {
		//dateToString := opr.Daily.Time[i].Format("2006-01-02")
		fm[opr.Daily.Time[i]] = Forecast{
			Temp2max:          opr.Daily.Temperature2mMax[i],
			UvIndexMax:        opr.Daily.UVIndexMax[i],
			PrecipProbability: opr.Daily.PrecipitationProbabilityMax[i],
		}
	}
	return fm, nil
}
