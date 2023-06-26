package webbr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValueType(t *testing.T) {
	testCases := []struct {
		name          string
		input         any
		expectedType  ValueType
		expectedErr   bool
		expectedValue any
	}{
		{
			name:          "string input",
			input:         "test",
			expectedType:  ValueTypeString,
			expectedErr:   false,
		},
		{
			name:         "boolean input",
			input:        true,
			expectedType: ValueTypeBool,
			expectedErr:  false,
		},
		{
			name:         "integer input",
			input:        2,
			expectedType: ValueTypeInt,
			expectedErr:  false,
		},
		{
			name:         "float input",
			input:        2.9,
			expectedType: ValueTypeFloat,
			expectedErr:  false,
		},
		{
			name:         "nil input",
			input:        nil,
			expectedType: ValueTypeUnknown,
			expectedErr:  true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			info, err := getValueTypeInfo(testCase.input)
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, info.valueType, testCase.expectedType)

			if testCase.expectedValue == nil {
				assert.Equal(t, info.underlying, testCase.expectedValue)
			}
		})
	}
}
