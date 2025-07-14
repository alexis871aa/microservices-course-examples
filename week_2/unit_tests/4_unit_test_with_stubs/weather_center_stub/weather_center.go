package weather_center_stub

type WeatherCenter struct {
	temperatureByCity map[string]float32
}

func NewWeatherCenter() *WeatherCenter {
	return &WeatherCenter{
		temperatureByCity: make(map[string]float32),
	}
}

func (w *WeatherCenter) SetTemperature(city string, temperature float32) {
}

func (w *WeatherCenter) GetTemperature(city string) (float32, error) {
	return 25.0, nil
}
