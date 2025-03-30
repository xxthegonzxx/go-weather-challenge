# Go Weather Check for Rain or Shine
Will It Rain or Shine?
I am terrible about remembering to bring an umbrella with me on rainy days and sunscreen on sunny days. Now I just need a script that runs every day at 6am, checks the weather and lets me know whether I need to pack an umbrella or sunscreen.

## About the Go program
It uses the https://ipinfo.io API to get the Geolocation (latitude and longitude) of the IP address on the machine and passes that information over to the openweathermap.org API and checks the daily weather.

### Build the Application
Use this option if you would like to build it locally on your computer and run it with a .env file.
```bash
make build
```
this will build a Go binary and dump it into the current directory in a bin folder. You can then add a .env file there and run it.

## Testing locally
To test the `rain` argument:
```bash
cd bin
./weather-app rain
```

To test the `shine` argument:
```bash
cd bin
./weather-app shine
```

## How to Run with Docker
```bash
make docker-build
```

### Run the Application
Check for rain:
```bash
make docker-run-rain
```

Check for shine:
```bash
make docker-run-shine
```

## Environment Setup
The API keys will be provided in a separate .env file via email or you can enter your own and save it to a .env file.

**NOTE**: You will need to get API keys from both https://api.openweathermap.org and https://ipinfo.io

Example `.env` file format:
```
OPEN_WEATHER_API_KEY=my-open-weather-key
IP_INFO_API_KEY=my-ip-info-key
OPEN_WEATHER_BASE_URL=https://api.openweathermap.org
IP_INFO_BASE_URL=https://ipinfo.io
```

## Available Make Commands
- `make build` - Build the Go application locally
- `make docker-build` - Build the Docker image
- `make docker-run-rain` - Run container to check rain forecast
- `make docker-run-shine` - Run container to check UV index
