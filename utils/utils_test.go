package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetOnuSN(t *testing.T) {
	testCases := []struct {
		hexFlag        bool
		input          string
		expectedOutput string
		expectedError  string
	}{
		{
			hexFlag:        true,
			input:          "SCOM00001B6D",
			expectedOutput: "53434f4d00001B6D",
		},

		{
			hexFlag:       false,
			input:         "SCO00001B6D",
			expectedOutput: "SCO00001B6D",
		},
	}

	for _, testCase := range testCases {
		OnuSNhex = testCase.hexFlag
		hexFormat := GetOnuSN(testCase.input)
		assert.Equal(t, testCase.expectedOutput, hexFormat)
	}
}


