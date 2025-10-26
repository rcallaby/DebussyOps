package agent

import (
	"context"
	"fmt"
)

// WeatherInfo represents weather data
type WeatherInfo struct {
	Temperature float64
	Conditions string
	Location   string
}

// GetWeather retrieves weather information for a given location
func GetWeather(ctx context.Context, location string) (*WeatherInfo, error) {
	// TODO: Implement actual weather fetching logic
	// This is just a placeholder structure
	return &WeatherInfo{
		Temperature: 0.0,
		Conditions: "Unknown",
		Location:   location,
	}, nil
}