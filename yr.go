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

type location struct {
	lat  string
	long string
	name string
}

type yrCollector struct {
	logger    *zap.Logger
	locations []location

	nowcastAirTemperature    *prometheus.GaugeVec
	nowcastPrecipitationRate *prometheus.GaugeVec
	nowcastRelativeHumidity  *prometheus.GaugeVec
	nowcastWindFromDirection *prometheus.GaugeVec
	nowcastWindSpeed         *prometheus.GaugeVec
	nowcastWindSpeedOfGust   *prometheus.GaugeVec

	forecastAirTemperature                        *prometheus.GaugeVec
	forecastCloudAreaFraction                     *prometheus.GaugeVec
	forecastUltravioletIndexClearSky              *prometheus.GaugeVec
	forecastNextOneHoursSymbol                    *prometheus.GaugeVec
	forecastNextOneHoursProbabilityOfPecipitation *prometheus.GaugeVec
	forecastNext12HoursProbabilityOfPecipitation  *prometheus.GaugeVec

	sunriseSunriseSecondsAfterMidnight       *prometheus.GaugeVec
	sunriseSunsetSecondsAfterMidnight        *prometheus.GaugeVec
	sunriseSolarNoonSecondsAfterMidnight     *prometheus.GaugeVec
	sunriseSolarNoonElevation                *prometheus.GaugeVec
	sunriseSolarMidnightSecondsAfterMidnight *prometheus.GaugeVec
	sunriseSolarMidnightElevation            *prometheus.GaugeVec

	nowcastScrapesFailed prometheus.Counter
}

var nowcastGroupLabelNames = []string{
	"coordinates",
	"name",
}

var forecastGroupLabelNames = []string{
	"coordinates",
	"name",
	"in_hours",
}

func NewYrCollector(namespace string, logger *zap.Logger, locations []location) prometheus.Collector {
	c := &yrCollector{
		logger:    logger,
		locations: locations,

		nowcastAirTemperature: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Subsystem: "nowcast", Name: "air_temperature"},
			nowcastGroupLabelNames,
		),

		nowcastPrecipitationRate: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Subsystem: "nowcast", Name: "precipitation_rate"},
			nowcastGroupLabelNames,
		),

		nowcastRelativeHumidity: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Subsystem: "nowcast", Name: "relative_humidity"},
			nowcastGroupLabelNames,
		),

		nowcastWindFromDirection: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Subsystem: "nowcast", Name: "wind_from_direction"},
			nowcastGroupLabelNames,
		),

		nowcastWindSpeed: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Subsystem: "nowcast", Name: "wind_speed"},
			nowcastGroupLabelNames,
		),

		nowcastWindSpeedOfGust: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Subsystem: "nowcast", Name: "wind_speed_of_gust"},
			nowcastGroupLabelNames,
		),

		forecastAirTemperature: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Subsystem: "forecast", Name: "air_temperature"},
			forecastGroupLabelNames,
		),

		forecastCloudAreaFraction: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Subsystem: "forecast", Name: "cloud_area_fraction"},
			forecastGroupLabelNames,
		),

		forecastUltravioletIndexClearSky: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Subsystem: "forecast", Name: "ultraviolet_index_clear_sky"},
			forecastGroupLabelNames,
		),

		forecastNextOneHoursSymbol: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Subsystem: "forecast", Name: "next_1_hours_symbol"},
			[]string{"coordinates", "name", "in_hours", "code"},
		),

		forecastNextOneHoursProbabilityOfPecipitation: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Subsystem: "forecast", Name: "next_1_hours_probability_of_precipitation"},
			forecastGroupLabelNames,
		),

		forecastNext12HoursProbabilityOfPecipitation: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Subsystem: "forecast", Name: "next_12_hours_probability_of_precipitation"},
			forecastGroupLabelNames,
		),

		sunriseSunriseSecondsAfterMidnight: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Subsystem: "sunrise", Name: "sunrise_seconds_after_midnight"},
			nowcastGroupLabelNames,
		),
		sunriseSunsetSecondsAfterMidnight: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Subsystem: "sunrise", Name: "sunset_seconds_after_midnight"},
			nowcastGroupLabelNames,
		),
		sunriseSolarNoonSecondsAfterMidnight: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Subsystem: "sunrise", Name: "sunset_solar_noon_seconds_after_midnight"},
			nowcastGroupLabelNames,
		),
		sunriseSolarNoonElevation: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Subsystem: "sunrise", Name: "sunset_solar_noon_elevation"},
			nowcastGroupLabelNames,
		),
		sunriseSolarMidnightSecondsAfterMidnight: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Subsystem: "sunrise", Name: "sunset_solar_midnight_seconds_after_midnight"},
			nowcastGroupLabelNames,
		),
		sunriseSolarMidnightElevation: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Subsystem: "sunrise", Name: "sunset_solar_midnight_elevation"},
			nowcastGroupLabelNames,
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

	c.forecastAirTemperature.Reset()
	c.forecastCloudAreaFraction.Reset()
	c.forecastUltravioletIndexClearSky.Reset()
	c.forecastNextOneHoursSymbol.Reset()
	c.forecastNextOneHoursProbabilityOfPecipitation.Reset()
	c.forecastNext12HoursProbabilityOfPecipitation.Reset()

	c.sunriseSunriseSecondsAfterMidnight.Reset()
	c.sunriseSunsetSecondsAfterMidnight.Reset()
	c.sunriseSolarNoonSecondsAfterMidnight.Reset()
	c.sunriseSolarNoonElevation.Reset()
	c.sunriseSolarMidnightSecondsAfterMidnight.Reset()
	c.sunriseSolarMidnightElevation.Reset()

	for _, loc := range c.locations {
		labels := prometheus.Labels{
			"coordinates": fmt.Sprintf("%s,%s", loc.lat, loc.long),
			"name":        loc.name,
		}

		if nowcast, err := c.getNowCast(loc); err != nil {
			c.logger.Error("Failed to update nowcast", zap.Error(err))
			c.nowcastScrapesFailed.Inc()
		} else {
			now := nowcast.Properties.Timeseries[0].Data.Instant.Details
			c.nowcastAirTemperature.With(labels).Set(now.AirTemperature)
			c.nowcastPrecipitationRate.With(labels).Set(now.PrecipitationRate)
			c.nowcastRelativeHumidity.With(labels).Set(now.RelativeHumidity)
			c.nowcastWindFromDirection.With(labels).Set(now.WindFromDirection)
			c.nowcastWindSpeed.With(labels).Set(now.WindSpeed)
			c.nowcastWindSpeedOfGust.With(labels).Set(now.WindSpeedOfGust)
		}

		if forecast, err := c.getForecast(loc); err != nil {
			c.logger.Error("Failed to update forecast", zap.Error(err))
			c.nowcastScrapesFailed.Inc()
		} else {
			for k, ts := range forecast.Properties.Timeseries {
				if k > 24 {
					break
				}
				d := ts.Data.Instant.Details

				forecastLabels := labels
				forecastLabels["in_hours"] = fmt.Sprintf("%d", k)

				c.forecastAirTemperature.With(forecastLabels).Set(d.AirTemperature)
				c.forecastCloudAreaFraction.With(forecastLabels).Set(d.CloudAreaFraction)
				c.forecastUltravioletIndexClearSky.With(forecastLabels).Set(d.UltravioletIndexClearSky)

				{
					gv, _ := c.forecastNextOneHoursSymbol.CurryWith(forecastLabels)
					gv.With(prometheus.Labels{"code": ts.Data.Next1Hours.Summary.SymbolCode}).Set(1)
				}

				c.forecastNextOneHoursProbabilityOfPecipitation.With(forecastLabels).Set(ts.Data.Next1Hours.Details.PrecipitationAmount)
				c.forecastNext12HoursProbabilityOfPecipitation.With(forecastLabels).Set(ts.Data.Next12Hours.Details.ProbabilityOfPrecipitation)
			}
		}

		if sunrise, err := c.getSunrise(loc); err != nil {
			c.logger.Error("Failed to update forecast", zap.Error(err))
			c.nowcastScrapesFailed.Inc()
		} else {
			labels := prometheus.Labels{
				"coordinates": fmt.Sprintf("%s,%s", loc.lat, loc.long),
				"name":        loc.name,
			}
			if (len(sunrise.Location.Time)) > 1 {
				s := sunrise.Location.Time[0]
				c.sunriseSunriseSecondsAfterMidnight.With(labels).Set(float64(secondsAfterMidnight(s.Sunrise.Time)))
				c.sunriseSunsetSecondsAfterMidnight.With(labels).Set(float64(secondsAfterMidnight(s.Sunset.Time)))
				c.sunriseSolarNoonSecondsAfterMidnight.With(labels).Set(float64(secondsAfterMidnight(s.Solarnoon.Time)))
				c.sunriseSolarNoonElevation.With(labels).Set(s.Solarnoon.Elevation)
				c.sunriseSolarMidnightSecondsAfterMidnight.With(labels).Set(float64(secondsAfterMidnight((s.Solarmidnight.Time))))
				c.sunriseSolarMidnightElevation.With(labels).Set(s.Solarmidnight.Elevation)
			}
		}
	}

	c.nowcastAirTemperature.Collect(ch)
	c.nowcastPrecipitationRate.Collect(ch)
	c.nowcastRelativeHumidity.Collect(ch)
	c.nowcastWindFromDirection.Collect(ch)
	c.nowcastWindSpeed.Collect(ch)
	c.nowcastWindSpeedOfGust.Collect(ch)

	c.forecastAirTemperature.Collect(ch)
	c.forecastCloudAreaFraction.Collect(ch)
	c.forecastUltravioletIndexClearSky.Collect(ch)
	c.forecastNextOneHoursSymbol.Collect(ch)
	c.forecastNextOneHoursProbabilityOfPecipitation.Collect(ch)
	c.forecastNext12HoursProbabilityOfPecipitation.Collect(ch)

	c.sunriseSunriseSecondsAfterMidnight.Collect(ch)
	c.sunriseSunsetSecondsAfterMidnight.Collect(ch)
	c.sunriseSolarNoonSecondsAfterMidnight.Collect(ch)
	c.sunriseSolarNoonElevation.Collect(ch)
	c.sunriseSolarMidnightSecondsAfterMidnight.Collect(ch)
	c.sunriseSolarMidnightElevation.Collect(ch)
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

func (c *yrCollector) getForecast(loc location) (*ForecastResponse, error) {
	url := fmt.Sprintf("https://api.met.no/weatherapi/locationforecast/2.0/complete?lat=%s&lon=%s", loc.lat, loc.long)
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

	var res ForecastResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal nowcast: %w", err)
	}

	return &res, nil
}

func (c *yrCollector) getSunrise(loc location) (*SunriseResponse, error) {
	date := time.Now().Format("2006-01-02")
	url := fmt.Sprintf("https://api.met.no/weatherapi/sunrise/2.0/.json?lat=%s&lon=%s&date=%s&offset=00:00", loc.lat, loc.long, date)
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

	var res SunriseResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal nowcast: %w", err)
	}

	return &res, nil
}

func secondsAfterMidnight(ts string) int {
	const layout = "2006-01-02T15:04:05-07:00"
	t, err := time.Parse(layout, ts)
	if err != nil {
		log.Println(ts, err)
		return 0
	}
	return t.Second() + t.Minute()*60 + t.Hour()*60*60

}
