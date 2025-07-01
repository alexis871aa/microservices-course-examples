package models

import (
	"time"
)

// Weather представляет информацию о погоде для конкретного города
type Weather struct {
	// Название города
	City string `json:"city"`
	// Температура в градусах Цельсия
	Temperature float64 `json:"temperature"`
	// Время последнего обновления данных
	UpdatedAt time.Time `json:"updated_at"`
}
