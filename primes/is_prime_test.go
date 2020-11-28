package primes

import (
	"fmt"
	"testing"
)

func TestIsPrime(t *testing.T) {
	testCases := []struct {
		nr       int
		expected bool
	}{
		{1, false},
		{2, true},
		{3, true},
		{4, false},
		{5, true},
		{6, false},
		{7, true},
		{8, false},
		{9, false},
		{10, false},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Testing %d", testCase.nr), func(t *testing.T) {
			actual := IsPrime(testCase.nr)
			if actual != testCase.expected {
				t.Errorf("Expected \"%d\" to be prime", testCase.nr)
			}
		})
	}
}
