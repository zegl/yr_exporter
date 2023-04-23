package main

import "time"

type NowCastResponse struct {
	Type       string     `json:"type"`
	Geometry   Geometry   `json:"geometry"`
	Properties Properties `json:"properties"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Units struct {
	AirTemperature      string `json:"air_temperature"`
	PrecipitationAmount string `json:"precipitation_amount"`
	PrecipitationRate   string `json:"precipitation_rate"`
	RelativeHumidity    string `json:"relative_humidity"`
	WindFromDirection   string `json:"wind_from_direction"`
	WindSpeed           string `json:"wind_speed"`
	WindSpeedOfGust     string `json:"wind_speed_of_gust"`
}

type Meta struct {
	UpdatedAt     time.Time `json:"updated_at"`
	Units         Units     `json:"units"`
	RadarCoverage string    `json:"radar_coverage"`
}

type InstantDetails struct {
	AirTemperature    float64 `json:"air_temperature"`
	PrecipitationRate float64 `json:"precipitation_rate"`
	RelativeHumidity  float64 `json:"relative_humidity"`
	WindFromDirection float64 `json:"wind_from_direction"`
	WindSpeed         float64 `json:"wind_speed"`
	WindSpeedOfGust   float64 `json:"wind_speed_of_gust"`
}

type Instant struct {
	Details InstantDetails `json:"details"`
}

type Summary struct {
	SymbolCode string `json:"symbol_code"`
}

type Details struct {
	PrecipitationAmount float64 `json:"precipitation_amount"`
}

type Next1Hours struct {
	Summary Summary `json:"summary"`
	Details Details `json:"details"`
}

type Data struct {
	Instant    Instant     `json:"instant"`
	Next1Hours *Next1Hours `json:"next_1_hours,omitempty"`
}

type Timeseries struct {
	Time time.Time `json:"time"`
	Data Data      `json:"data,omitempty"`
}

type Properties struct {
	Meta       Meta         `json:"meta"`
	Timeseries []Timeseries `json:"timeseries"`
}
