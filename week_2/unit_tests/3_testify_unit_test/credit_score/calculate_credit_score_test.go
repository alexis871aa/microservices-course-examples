package credit_score

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateCreditScore(t *testing.T) {
	t.Run("valid male engineer", func(t *testing.T) {
		client := Client{"male", 30, "engineer", 7, 60000}
		expected := 750

		result, err := CalculateCreditScore(client)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("valid female teacher", func(t *testing.T) {
		client := Client{"female", 22, "teacher", 3, 40000}
		expected := 510

		result, err := CalculateCreditScore(client)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("valid male doctor", func(t *testing.T) {
		client := Client{"male", 45, "doctor", 20, 120000}
		expected := 850

		result, err := CalculateCreditScore(client)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("valid female retired", func(t *testing.T) {
		client := Client{"female", 55, "retired", 30, 20000}
		expected := 460

		result, err := CalculateCreditScore(client)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("invalid age", func(t *testing.T) {
		client := Client{"male", -1, "engineer", 7, 60000}

		_, err := CalculateCreditScore(client)
		require.Error(t, err)
		assert.Equal(t, "invalid age", err.Error())
	})

	t.Run("invalid salary", func(t *testing.T) {
		client := Client{"male", 30, "engineer", 7, -5000}

		_, err := CalculateCreditScore(client)
		require.Error(t, err)
		assert.Equal(t, "invalid salary", err.Error())
	})
}
