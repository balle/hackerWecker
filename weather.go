// Get information about the weather
package hackerWecker

import (
	"fmt"
	"time"

	"github.com/EricNeid/go-openweather"
)

const unknownTemp int = 99999

type Weather struct {
	Description string
	Temp        int
	TempMin     int
	TempMax     int
	WindSpeed   int
	Rain        int
	Sunrise     int
	Sunset      int
	TempUnit    string
}

func GetWeather() Weather {
	var result Weather

	if len(config.WeatherLocation) > 0 && len(config.OpenWeatherAPIKey) > 0 {
		qry := openweather.NewQueryForCity(config.OpenWeatherAPIKey, config.WeatherLocation)
		resp, err := qry.Weather()

		if err != nil {
			LogError(fmt.Sprintf("Cannot get weather data for %s: %v\n", config.WeatherLocation, err))
		} else {
			forecast, err := qry.DailyForecast5()

			if err != nil && len(forecast.List) > 0 {
				result.TempMin = int(forecast.List[0].Temp.Min)
				result.TempMax = int(forecast.List[0].Temp.Max)
			} else {
				result.TempMin = unknownTemp
				result.TempMax = unknownTemp
			}

			result.Description = resp.Weather[0].Description
			result.Temp = int(resp.Main.Temp)
			result.WindSpeed = int(resp.Wind.Speed)
			result.Rain = resp.Rain.ThreeH
			result.Sunrise = resp.Sys.Sunrise
			result.Sunset = resp.Sys.Sunset
		}
	}

	return result
}

func ReadWeather() {
	weather := GetWeather()
	tempUnit := "degrees"

	Speak(weather.Description)

	if config.WeatherUnit != "metric" {
		tempUnit = "fahrenheit"
	}

	Speak(fmt.Sprintf("Temparature is at %d %s", weather.Temp, tempUnit))

	if weather.TempMax != unknownTemp {
		Speak(fmt.Sprintf("Maximum %d %s", weather.TempMax, tempUnit))
	}

	sunsetTime := time.Unix(int64(weather.Sunset), 0)
	Speak(fmt.Sprintf("Sunset is at %d %d", sunsetTime.Hour(), sunsetTime.Minute()))
}
