package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"weather-api/src/models/geonorge"
	"weather-api/src/models/weather"
)

func GetWeather(location *geonorge.Location) (*weather.WeatherResponse, error) {
	u, err := url.Parse("https://api.met.no/weatherapi/locationforecast/2.0/compact")
	if err != nil {
		return nil, err
	}

	query := u.Query()
	query.Add("lat", fmt.Sprintf("%f", location.Lat))
	query.Add("lon", fmt.Sprintf("%f", location.Lon))
	u.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "weather-api weather-api.azurewebsites.net")

	client := &http.Client{}

	weatherResp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer weatherResp.Body.Close()

	weatherData, err := io.ReadAll(weatherResp.Body)
	if err != nil {
		return nil, err
	}

	var weather weather.WeatherResponse
	err = json.Unmarshal(weatherData, &weather)
	if err != nil {
		return nil, err
	}

	return &weather, nil
}
