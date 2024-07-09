package main

import (
	"encoding/json"
	"time"
)

type AirQualityReading struct {
	SensorID  string    `json:"sensor_id"`
	Timestamp time.Time `json:"timestamp"`
	PM25      float64   `json:"pm25"`
	CO2       float64   `json:"co2"`
}

func parseReadings(data []byte) ([]AirQualityReading, error) {

	var airQualityReading []AirQualityReading
	err := json.Unmarshal(data, &airQualityReading)
	if err != nil {
		return nil, err
	}

	return airQualityReading, nil
}
