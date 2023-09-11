package geonorge

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

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
