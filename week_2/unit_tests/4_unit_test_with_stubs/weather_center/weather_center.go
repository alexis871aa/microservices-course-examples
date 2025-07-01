package weather_center

import "errors"

var ErrCityNotFound = errors.New("city not found")

type WeatherCenter struct {
	temperatureByCity map[string]float32
}

func NewWeatherCenter() *WeatherCenter {
	return &WeatherCenter{
		temperatureByCity: make(map[string]float32),
	}
}

func (w *WeatherCenter) SetTemperature(city string, temperature float32) {
	w.temperatureByCity[city] = temperature
}

func (w *WeatherCenter) GetTemperature(city string) (float32, error) {
	temperature, ok := w.temperatureByCity[city]
	if !ok {
		return 0, ErrCityNotFound
	}

	return temperature, nil
}
