package models

import (
	"sync"
)

// WeatherStorage представляет потокобезопасное хранилище данных о погоде
type WeatherStorage struct {
	mu       sync.RWMutex
	weathers map[string]*Weather
}

// NewWeatherStorage создает новое хранилище данных о погоде
func NewWeatherStorage() *WeatherStorage {
	return &WeatherStorage{
		weathers: make(map[string]*Weather),
	}
}

// GetWeather возвращает информацию о погоде по имени города.
// Если город не найден, возвращает nil.
func (s *WeatherStorage) GetWeather(city string) *Weather {
	s.mu.RLock()
	defer s.mu.RUnlock()

	weather, ok := s.weathers[city]
	if !ok {
		return nil
	}

	return weather
}

// UpdateWeather обновляет данные о погоде для указанного города.
// Если города нет в хранилище, создает новую запись.
func (s *WeatherStorage) UpdateWeather(weather *Weather) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.weathers[weather.City] = weather
}
