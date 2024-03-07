package shorten

import (
	"testing"
)

func TestShorten(t *testing.T) {
	t.Run("Return shorten id", func(t *testing.T) {
		type testCase struct {
			id       uint32
			expected string
		}

		testCases := []testCase{
			{id: 0, expected: ""},
			{id: 1, expected: "n"},
			{id: 10, expected: "H"},
			{id: 100, expected: "n6"},
			{id: 12345, expected: "Jv8"},
			{id: 9876543, expected: "ihzV"},
			{id: 99999, expected: "L6g"},
			{id: 54321, expected: "Ed3"},
			{id: 1024, expected: "Mv"},
		}

		for _, tc := range testCases {
			actual := Shorten(tc.id)
			if actual != tc.expected {
				t.Fatalf("Incorrect shorten result for id: %d, actual: %s, expected: %s", tc.id, actual, tc.expected)
			}
		}

	})

	t.Run("Return idempotent", func(t *testing.T) {
		for i := 0; i < 250; i++ {
			if Shorten(1024) != "Mv" {
				t.Fatalf("Result is not idempotent")
			}
		}
	})
}
