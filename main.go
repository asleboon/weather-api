package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type GeoNorgeResponse struct {
	Metadata Metadata   `json:"metadata"`
	Adresser []Adresser `json:"adresser"`
}

type Adresser struct {
	Adressenavn                         string               `json:"adressenavn"`
	Adressetekst                        string               `json:"adressetekst"`
	Adressetilleggsnavn                 interface{}          `json:"adressetilleggsnavn"` // using interface{} to allow null value
	Adressekode                         int                  `json:"adressekode"`
	Nummer                              int                  `json:"nummer"`
	Bokstav                             string               `json:"bokstav"`
	Kommunenummer                       string               `json:"kommunenummer"`
	Kommunenavn                         string               `json:"kommunenavn"`
	Gardsnummer                         int                  `json:"gardsnummer"`
	Bruksnummer                         int                  `json:"bruksnummer"`
	Festenummer                         int                  `json:"festenummer"`
	Undernummer                         interface{}          `json:"undernummer"` // using interface{} to allow null value
	Bruksenhetsnummer                   []string             `json:"bruksenhetsnummer"`
	Objtype                             string               `json:"objtype"`
	Poststed                            string               `json:"poststed"`
	Postnummer                          string               `json:"postnummer"`
	AdressetekstUtenAdressetilleggsnavn string               `json:"adressetekstutenadressetilleggsnavn"`
	Stedfestingverifisert               bool                 `json:"stedfestingverifisert"`
	Representasjonspunkt                Representasjonspunkt `json:"representasjonspunkt"`
	Oppdateringsdato                    string               `json:"oppdateringsdato"` // assuming the date is in a parsable format
}

type Representasjonspunkt struct {
	EPSG string  `json:"epsg"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

type Metadata struct {
	ViserFra          int    `json:"viserFra"`
	SokeStreng        string `json:"sokeStreng"`
	Side              int    `json:"side"`
	TreffPerSide      int    `json:"treffPerSide"`
	TotaltAntallTreff int    `json:"totaltAntallTreff"`
	AsciiKompatibel   bool   `json:"asciiKompatibel"`
	ViserTil          int    `json:"viserTil"`
}

// WeatherResponse represents the structure of the weather response
type WeatherResponse struct {
	Type       string     `json:"type"`
	Geometry   Geometry   `json:"geometry"`
	Properties Properties `json:"properties"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Properties struct {
	Meta       Meta       `json:"meta"`
	Timeseries []Timesery `json:"timeseries"`
}

type Meta struct {
	UpdatedAt string `json:"updated_at"`
	Units     Units  `json:"units"`
}

type Units struct {
	AirPressureAtSeaLevel string `json:"air_pressure_at_sea_level"`
	AirTemperature        string `json:"air_temperature"`
	CloudAreaFraction     string `json:"cloud_area_fraction"`
	PrecipitationAmount   string `json:"precipitation_amount"`
	RelativeHumidity      string `json:"relative_humidity"`
	WindFromDirection     string `json:"wind_from_direction"`
	WindSpeed             string `json:"wind_speed"`
}

type Timesery struct {
	Time string `json:"time"`
	Data Data   `json:"data"`
}

type Data struct {
	Instant     Instant   `json:"instant"`
	Next12Hours NextHours `json:"next_12_hours"`
	Next1Hours  NextHours `json:"next_1_hours"`
	Next6Hours  NextHours `json:"next_6_hours"`
}

type Instant struct {
	Details InstantDetails `json:"details"`
}

type InstantDetails struct {
	AirPressureAtSeaLevel float64 `json:"air_pressure_at_sea_level"`
	AirTemperature        float64 `json:"air_temperature"`
	CloudAreaFraction     float64 `json:"cloud_area_fraction"`
	RelativeHumidity      float64 `json:"relative_humidity"`
	WindFromDirection     float64 `json:"wind_from_direction"`
	WindSpeed             float64 `json:"wind_speed"`
}

type NextHours struct {
	Summary Summary           `json:"summary"`
	Details *NextHoursDetails `json:"details,omitempty"`
}

type NextHoursDetails struct {
	PrecipitationAmount float64 `json:"precipitation_amount"`
}

type Summary struct {
	SymbolCode string `json:"symbol_code"`
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.GET("/weather", func(c *gin.Context) {
		address := c.Query("address")
		city := c.Query("city")

		if address == "" || city == "" {
			c.JSON(400, gin.H{"error": "Bad request: missing query params"})
			return
		}
		// Step 1: Get the coordinates from the address
		u, err := url.Parse("https://ws.geonorge.no/adresser/v1/sok?fuzzy=true&utkoordsys=4258&treffPerSide=1&side=0&asciiKompatibel=true")

		if err != nil {
			return
		}

		q := u.Query()

		q.Add("sok", address+", "+city)

		u.RawQuery = q.Encode()

		urlString := u.String()

		locationResp, err := http.Get(urlString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		defer locationResp.Body.Close()

		locationData, err := io.ReadAll(locationResp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var location GeoNorgeResponse
		err = json.Unmarshal(locationData, &location)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(location.Adresser) == 0 {
			c.JSON(404, gin.H{"message": "Address not found"})
			return
		}

		lat := location.Adresser[0].Representasjonspunkt.Lat
		lon := location.Adresser[0].Representasjonspunkt.Lon

		u2, err := url.Parse("https://api.met.no/weatherapi/locationforecast/2.0/compact")

		if err != nil {
			return
		}

		q2 := u2.Query()

		q2.Add("lat", fmt.Sprintf("%f", lat))
		q2.Add("lon", fmt.Sprintf("%f", lon))

		u2.RawQuery = q2.Encode()

		urlString2 := u2.String()

		// Step 2: Get the weather using the coordinates
		// Add headers
		req, err := http.NewRequest("GET", urlString2, nil)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		req.Header.Set("User-Agent", "asledev asledev.vercel.app")

		client := &http.Client{}

		weatherResp, err := client.Do(req)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		defer weatherResp.Body.Close()

		weatherData, err := io.ReadAll(weatherResp.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		log.Println(string(weatherData))

		var weather WeatherResponse
		err = json.Unmarshal(weatherData, &weather)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"data": weather})
		return

	})

	r.Run(":8080")
}
