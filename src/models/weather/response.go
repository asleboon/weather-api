package weather

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
