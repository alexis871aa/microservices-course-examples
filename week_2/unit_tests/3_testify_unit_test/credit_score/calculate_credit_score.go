package credit_score

import "errors"

// Client структура для хранения данных клиента
type Client struct {
	Gender        string
	Age           int
	Profession    string
	Experience    int
	AverageSalary float64
}

// CalculateCreditScore функция для расчета кредитного рейтинга
func CalculateCreditScore(client Client) (int, error) {
	if client.Age < 0 {
		return 0, errors.New("invalid age")
	}

	score := 0

	// Пример логики расчета кредитного рейтинга
	if client.Gender == "male" {
		score += 50
	} else if client.Gender == "female" {
		score += 60
	}

	if client.Age >= 18 && client.Age <= 25 {
		score += 100
	} else if client.Age > 25 && client.Age <= 35 {
		score += 150
	} else if client.Age > 35 && client.Age <= 50 {
		score += 200
	} else if client.Age > 50 {
		score += 100
	}

	if client.Profession == "engineer" {
		score += 200
	} else if client.Profession == "teacher" {
		score += 150
	} else {
		score += 100
	}

	if client.Experience >= 1 && client.Experience <= 5 {
		score += 100
	} else if client.Experience > 5 && client.Experience <= 10 {
		score += 150
	} else if client.Experience > 10 {
		score += 200
	}

	err := checkSalary(client)
	if err != nil {
		return 0, err
	}

	if client.AverageSalary >= 30000 && client.AverageSalary <= 50000 {
		score += 100
	} else if client.AverageSalary > 50000 && client.AverageSalary <= 100000 {
		score += 200
	} else if client.AverageSalary > 100000 {
		score += 300
	}

	// Ограничение рейтинга от 0 до 1000
	if score > 1000 {
		score = 1000
	} else if score < 0 {
		score = 0
	}

	return score, nil
}

func checkSalary(client Client) error {
	if client.AverageSalary < 0 {
		return errors.New("invalid salary")
	}

	return nil
}
