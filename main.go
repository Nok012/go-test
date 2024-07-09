package main

import (
	"encoding/json"
	"fmt"
	"reflect"
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

func calculateAverage(readings []AirQualityReading) map[string]float64 {
	if len(readings) == 0 {
		return nil
	}

	totals := make(map[string]float64)
	count := float64(len(readings))

	for _, reading := range readings {
		val := reflect.ValueOf(reading)
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			if field.Type.Kind() == reflect.Float64 {
				totals[field.Name] += val.Field(i).Float()
			}
		}
	}

	averages := make(map[string]float64)
	for k, v := range totals {

		averages[k] = v / count
	}
	return averages
}

func findHighestPollutantByHour(readings []AirQualityReading) map[int]string {

	if len(readings) == 0 {
		return nil
	}

	result := make(map[int]string)
	averagePollutants := make(map[int]map[string]float64)
	count := make(map[int]map[string]float64)

	for _, reading := range readings {

		hour := reading.Timestamp.Hour()
		if _, ok := averagePollutants[hour]; !ok {
			averagePollutants[hour] = make(map[string]float64)
			count[hour] = make(map[string]float64)
		}

		val := reflect.ValueOf(reading)
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			if field.Type.Kind() == reflect.Float64 {

				averagePollutants[hour][field.Name] += val.Field(i).Float()
				count[hour][field.Name]++
			}
		}
	}

	for hour, pollutants := range averagePollutants {
		highestPollutant := ""
		highestAverage := 0.0

		for pollutant, total := range pollutants {
			average := total / count[hour][pollutant]
			if average > highestAverage {
				highestAverage = average
				highestPollutant = pollutant
			}
		}

		result[hour] = highestPollutant
	}

	return result
}

func main() {

	jsonData := []byte(`[
		{"sensor_id": "S001", "timestamp": "2023-12-28T10:00:00Z", "pm25": 25.5, "co2": 410.2},
		{"sensor_id": "S002", "timestamp": "2023-12-28T10:05:00Z", "pm25": 30.8, "co2": 405.7},
		{"sensor_id": "S001", "timestamp": "2023-12-28T11:00:00Z", "pm25": 18.2, "co2": 395.1}
	]`)

	readings, err := parseReadings(jsonData)
	if err != nil {
		fmt.Println("Error parsing readings:", err)
		return
	}

	averages := calculateAverage(readings)
	for pollutant, avg := range averages {
		fmt.Printf("Average %s: %.2f\n", pollutant, avg)
	}

	result := findHighestPollutantByHour(readings)
	for hour, pollutant := range result {
		fmt.Printf("Hour %d: Highest Pollutant - %s\n", hour, pollutant)
	}

}
