package credit_score

import "testing"

func TestCalculateCreditScore(t *testing.T) {
	tests := []struct {
		client   Client
		expected int
	}{
		{
			client:   Client{"male", 30, "engineer", 7, 60000, 0},
			expected: 750,
		},
		{
			client:   Client{"female", 22, "teacher", 3, 40000, 1},
			expected: 460,
		},
		{
			client:   Client{"male", 45, "doctor", 20, 120000, 0},
			expected: 850,
		},
		{
			client:   Client{"female", 55, "retired", 30, 20000, 2},
			expected: 410,
		},
	}

	for _, test := range tests {
		result := CalculateCreditScore(test.client)
		if result != test.expected {
			t.Errorf("CalculateCreditScore(%#v) = %d; want %d", test.client, result, test.expected)
		}
	}
}
