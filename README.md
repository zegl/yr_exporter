# yr_exporter

Exports weather forecasts from YR.no

## Metrics

```
# HELP yr_nowcast_air_temperature 
# TYPE yr_nowcast_air_temperature gauge
yr_nowcast_air_temperature{coordinates="60.10,9.58",name="oslo"} 3.7
# HELP yr_nowcast_precipitation_rate 
# TYPE yr_nowcast_precipitation_rate gauge
yr_nowcast_precipitation_rate{coordinates="60.10,9.58",name="oslo"} 0
# HELP yr_nowcast_relative_humidity 
# TYPE yr_nowcast_relative_humidity gauge
yr_nowcast_relative_humidity{coordinates="60.10,9.58",name="oslo"} 100
# HELP yr_nowcast_wind_from_direction 
# TYPE yr_nowcast_wind_from_direction gauge
yr_nowcast_wind_from_direction{coordinates="60.10,9.58",name="oslo"} 107.1
# HELP yr_nowcast_wind_speed 
# TYPE yr_nowcast_wind_speed gauge
yr_nowcast_wind_speed{coordinates="60.10,9.58",name="oslo"} 0.8
# HELP yr_nowcast_wind_speed_of_gust 
# TYPE yr_nowcast_wind_speed_of_gust gauge
yr_nowcast_wind_speed_of_gust{coordinates="60.10,9.58",name="oslo"} 2

# HELP yr_forecast_air_temperature 
# TYPE yr_forecast_air_temperature gauge
yr_forecast_air_temperature{coordinates="60.10,9.58",in_hours="0",name="oslo"} 4.1
yr_forecast_air_temperature{coordinates="60.10,9.58",in_hours="1",name="oslo"} 3.5
yr_forecast_air_temperature{coordinates="60.10,9.58",in_hours="2",name="oslo"} 3.2
yr_forecast_air_temperature{coordinates="60.10,9.58",in_hours="3",name="oslo"} 2.9
# HELP yr_forecast_in_one_hour_symbol 
# TYPE yr_forecast_in_one_hour_symbol gauge
yr_forecast_in_one_hour_symbol{code="fog",coordinates="60.10,9.58",in_hours="0",name="oslo"} 1
yr_forecast_in_one_hour_symbol{code="fog",coordinates="60.10,9.58",in_hours="1",name="oslo"} 1
```
