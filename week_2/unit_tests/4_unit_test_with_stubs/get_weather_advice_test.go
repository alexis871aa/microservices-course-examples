package main

import (
	"testing"

	weatherCenterStub "github.com/olezhek28/microservices-course-examples/week_2/unit_tests/4_unit_test_with_stubs/weather_center_stub"
	"github.com/stretchr/testify/require"
)

func TestGetWeatherAdvice(t *testing.T) {
	tests := []struct {
		name     string
		city     string
		expected string
	}{
		{
			name:     "Температура +25 градусов",
			city:     "Москва",
			expected: "Отличная погода для прогулок",
		},
	}

	clientStub := weatherCenterStub.NewWeatherCenter()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := getWeatherAdvice(clientStub, test.city)
			require.NoError(t, err)
			require.Equal(t, test.expected, res)
		})
	}
}
