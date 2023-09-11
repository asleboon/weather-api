package services

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"weather-api/src/models/geonorge"

	"github.com/gin-gonic/gin"
)

func GetLocation(c *gin.Context) (*geonorge.Location, error) {
	address := c.Query("address")
	city := c.Query("city")

	u, err := url.Parse("https://ws.geonorge.no/adresser/v1/sok?fuzzy=true&utkoordsys=4258&treffPerSide=1&side=0&asciiKompatibel=true")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("sok", address+", "+city)
	u.RawQuery = q.Encode()
	urlString := u.String()

	locationResp, err := http.Get(urlString)
	if err != nil {
		return nil, err
	}
	defer locationResp.Body.Close()

	locationData, err := io.ReadAll(locationResp.Body)
	if err != nil {
		return nil, err
	}

	var location geonorge.GeoNorgeResponse
	err = json.Unmarshal(locationData, &location)
	if err != nil {
		return nil, err
	}

	if len(location.Adresser) == 0 {
		return nil, errors.New("no location found")
	}

	data := &geonorge.Location{
		Lat: location.Adresser[0].Representasjonspunkt.Lat,
		Lon: location.Adresser[0].Representasjonspunkt.Lon,
	}

	return data, nil
}
