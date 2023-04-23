package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

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

type location struct {
	lat  string
	long string
	name string
}

type yrCollector struct {
	logger                   *zap.Logger
	locations                []location
	nowcastAirTemperature    *prometheus.GaugeVec
	nowcastPrecipitationRate *prometheus.GaugeVec
	nowcastRelativeHumidity  *prometheus.GaugeVec
	nowcastWindFromDirection *prometheus.GaugeVec
	nowcastWindSpeed         *prometheus.GaugeVec
	nowcastWindSpeedOfGust   *prometheus.GaugeVec
	nowcastScrapesFailed     prometheus.Counter
}

var variableGroupLabelNames = []string{
	"coordinates",
	"name",
}

func NewYrCollector(namespace string, logger *zap.Logger, locations []location) prometheus.Collector {
	c := &yrCollector{
		logger:    logger,
		locations: locations,

		nowcastAirTemperature: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Subsystem: "nowcast", Name: "air_temperature"},
			variableGroupLabelNames,
		),

		nowcastPrecipitationRate: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Subsystem: "nowcast", Name: "precipitation_rate"},
			variableGroupLabelNames,
		),

		nowcastRelativeHumidity: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Subsystem: "nowcast", Name: "relative_humidity"},
			variableGroupLabelNames,
		),

		nowcastWindFromDirection: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Subsystem: "nowcast", Name: "wind_from_direction"},
			variableGroupLabelNames,
		),

		nowcastWindSpeed: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Subsystem: "nowcast", Name: "wind_speed"},
			variableGroupLabelNames,
		),

		nowcastWindSpeedOfGust: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Subsystem: "nowcast", Name: "wind_speed_of_gust"},
			variableGroupLabelNames,
		),

		nowcastScrapesFailed: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "nowcast",
				Name:      "scrapes_failed",
				Help:      "Count of scrapes of group data from YR that have failed",
			},
		),
	}

	return c
}

func (c yrCollector) Describe(ch chan<- *prometheus.Desc) {
	c.nowcastAirTemperature.Describe(ch)
}

func (c *yrCollector) Collect(ch chan<- prometheus.Metric) {
	c.nowcastAirTemperature.Reset()
	c.nowcastPrecipitationRate.Reset()
	c.nowcastRelativeHumidity.Reset()
	c.nowcastWindFromDirection.Reset()
	c.nowcastWindSpeed.Reset()
	c.nowcastWindSpeedOfGust.Reset()

	for _, loc := range c.locations {

		if nowcast, err := c.getNowCast(loc); err != nil {
			c.logger.Error("Failed to update nowcast", zap.Error(err))
			c.nowcastScrapesFailed.Inc()
		} else {

			now := nowcast.Properties.Timeseries[0].Data.Instant.Details

			labels := prometheus.Labels{
				"coordinates": fmt.Sprintf("%s,%s", loc.lat, loc.long),
				"name":        loc.name,
			}

			c.nowcastAirTemperature.With(labels).Set(now.AirTemperature)
			c.nowcastPrecipitationRate.With(labels).Set(now.PrecipitationRate)
			c.nowcastRelativeHumidity.With(labels).Set(now.RelativeHumidity)
			c.nowcastWindFromDirection.With(labels).Set(now.WindFromDirection)
			c.nowcastWindSpeed.With(labels).Set(now.WindSpeed)
			c.nowcastWindSpeedOfGust.With(labels).Set(now.WindSpeedOfGust)
		}
	}

	c.nowcastAirTemperature.Collect(ch)
	c.nowcastPrecipitationRate.Collect(ch)
	c.nowcastRelativeHumidity.Collect(ch)
	c.nowcastWindFromDirection.Collect(ch)
	c.nowcastWindSpeed.Collect(ch)
	c.nowcastWindSpeedOfGust.Collect(ch)
}

func (c *yrCollector) getNowCast(loc location) (*NowCastResponse, error) {
	url := fmt.Sprintf("https://api.met.no/weatherapi/nowcast/2.0/complete?lat=%s&lon=%s", loc.lat, loc.long)
	log.Printf("Fetching %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("User-Agent", "https://github.com/zegl/yr_exporter")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get nowcast: %w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read nowcast: %w", err)
	}
	defer resp.Body.Close()

	var res NowCastResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal nowcast: %w", err)
	}

	return &res, nil
}
