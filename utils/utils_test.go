// Copyright 2018 Open Networking Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//Package utils gives type conversion helpers


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
			expectedOutput: "53434F4D00001B6D",
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


