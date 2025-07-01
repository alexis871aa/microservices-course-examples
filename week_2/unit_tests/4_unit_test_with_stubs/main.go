package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	weatherCenter "github.com/olezhek28/microservices-course-examples/week_2/unit_tests/4_unit_test_with_stubs/weather_center"
)

type weatherCenterClient interface {
	SetTemperature(city string, temperature float32)
	GetTemperature(city string) (float32, error)
}

func main() {
	fmt.Println("Привет! Я погодный помощник:)")
	fmt.Print("Хочешь узнать погоду у себя в городе? Тогда введи его название: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		log.Fatalln("Что-то пошло не так :(")
	}

	city := strings.TrimSpace(scanner.Text())

	client := weatherCenter.NewWeatherCenter()
	client.SetTemperature("Москва", 25.0)

	weatherAdvice, err := getWeatherAdvice(client, city)
	if err != nil {
		log.Fatalln("Произошла ошибка: ", err.Error())
	}

	fmt.Println(weatherAdvice)
}

func getWeatherAdvice(client weatherCenterClient, city string) (string, error) {
	temperature, err := client.GetTemperature(city)
	if err != nil {
		return "", err
	}

	switch {
	case temperature >= -71.0 && temperature <= -40.0:
		return "Лучше не суйся на улицу", nil
	case temperature > 40.0 && temperature <= -20.0:
		return "Можно идти гулять, но одевайся теплее", nil
	case temperature > -20.0 && temperature <= 0.0:
		return "Прохладно, но можно выйти на улицу", nil
	case temperature > 0.0 && temperature <= 15.0:
		return "Температура нормальная, можно гулять", nil
	case temperature > 15.0 && temperature <= 25.0:
		return "Отличная погода для прогулок", nil
	case temperature > 25.0 && temperature <= 35.0:
		return "Жарковато, но можно выйти на улицу", nil
	case temperature > 35.0 && temperature <= 52.0:
		return "Жарко, лучше остаться дома", nil
	default:
		return "А ты точно на Земле?", nil
	}
}
