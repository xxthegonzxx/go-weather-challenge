// I am terrible about remembering to bring an umbrella with me on rainy days and sunscreen on sunny days. Now I just
// need a script that runs every day at 6am, checks the weather and lets me know whether I need to pack an umbrella or
// sunscreen. -->

// every day at 6am
// Based on command argument: rain/shine
// checks weather API based on current location and command line argument
// return Umbrella or Sunscreen

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/joho/godotenv"
)

// WeatherResponse represents the OpenWeather API response
type WeatherResponse struct {
	Daily []struct {
		Dt  int64   `json:"dt"`  // Time of forecasted data, unix timestamp
		Pop float64 `json:"pop"` // Probability of precipitation (0-1)
		Uvi float64 `json:"uvi"` // UV index
	} `json:"daily"`
}

// init initializes the application by loading environment variables from a .env file
// and verifying that all required API keys are present. If the .env file is not found,
// it falls back to using system environment variables.
func init() {
	// Only try to load .env file if not running in Docker
	if os.Getenv("DOCKER_CONTAINER") == "" {
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: .env file not found. Using environment variables.")
		}
	}

	required := []string{"IP_INFO_API_KEY", "OPEN_WEATHER_API_KEY"}
	for _, env := range required {
		if os.Getenv(env) == "" {
			log.Fatalf("Required environment variable %s is not set", env)
		}
	}
}

// checkWeather checks the weather based on the provided weather type (rain or shine).
func checkWeather(weatherType string) string {
	if weatherType != "rain" && weatherType != "shine" {
		return "Acceptable arguments are either rain or shine."
	}

	lat, lng, err := getLocation()
	if err != nil {
		log.Fatal(err)
	}

	needsProtection, err := getWeather(lat, lng, weatherType)
	if err != nil {
		log.Fatal(err)
	}

	if weatherType == "rain" {
		if needsProtection {
			return "Bring an umbrella!"
		}
		return "No umbrella needed today."
	} else {
		if needsProtection {
			return "Bring sunscreen!"
		}
		return "No sunscreen needed today."
	}
}

// getLocation retrieves the user's location using the IPinfo API (https://ipinfo.io).
// It returns the latitude and longitude of the user's location.
// If an error occurs during the API call, it returns an error.
func getLocation() (latitude, longitude float64, err error) {
	// Initialize the client with the API token
	client := ipinfo.NewClient(nil, nil, os.Getenv("IP_INFO_API_KEY"))

	loc, err := client.GetIPLocation(nil)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get IP location: %v", err)
	}

	coords := strings.Split(loc, ",")
	if len(coords) != 2 {
		return 0, 0, fmt.Errorf("invalid location format: %s", loc)
	}

	lat, err := strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse latitude: %v", err)
	}

	lng, err := strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse longitude: %v", err)
	}

	return lat, lng, nil
}

// getWeather retrieves the weather information based on the user's location coordinates.
// It checks the chance of rain and UV index for today.
// If the chance of rain is greater than 5% or the UV index is greater than 3, it returns true.
// Otherwise, it returns false.
// It uses the OpenWeather API (https://openweathermap.org/api) to get the weather data.
func getWeather(lat, lng float64, weatherType string) (bool, error) {
	baseURL := os.Getenv("OPEN_WEATHER_BASE_URL")
	apiKey := os.Getenv("OPEN_WEATHER_API_KEY")

	url := fmt.Sprintf("%s/data/3.0/onecall?lat=%f&lon=%f&appid=%s",
		baseURL, lat, lng, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var weather WeatherResponse
	if err := json.Unmarshal(body, &weather); err != nil {
		return false, err
	}

	if len(weather.Daily) == 0 {
		return false, fmt.Errorf("no daily forecast available")
	}

	today := weather.Daily[0]

	if weatherType == "rain" {
		// Chance of rain is given as a decimal (0-1), so multiply by 100
		chanceOfRain := today.Pop * 100
		fmt.Printf("Chance of rain today: %.1f%%\n", chanceOfRain)
		return chanceOfRain > 5.0, nil
	}
	// Must be "shine" due to earlier validation
	fmt.Printf("UV Index today: %.1f\n", today.Uvi)
	return today.Uvi > 3.0, nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a weather type (rain or shine).")
	}

	weatherType := os.Args[1]
	message := checkWeather(weatherType)
	fmt.Println(message)
}
