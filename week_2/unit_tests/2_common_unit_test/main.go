package main

import (
	"fmt"

	"github.com/olezhek28/microservices-course-examples/week_2/unit_tests/2_common_unit_test/credit_score"
)

func main() {
	client := credit_score.Client{
		Gender:        "male",
		Age:           30,
		Profession:    "engineer",
		Experience:    7,
		AverageSalary: 60000,
	}

	creditScore := credit_score.CalculateCreditScore(client)
	fmt.Printf("Credit Score: %d\n", creditScore)
}
