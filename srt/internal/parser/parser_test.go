package parser

import (
	"testing"

	"github.com/florentsorel/subtitle/srt/internal/lexer"
)

func TestNewParser(t *testing.T) {
	tests := []struct {
		input         string
		expectedError string
	}{
		{
			input:         `-->`,
			expectedError: "expected INDEX, got ARROW at line 1, column 1",
		},
		{
			input:         `First line`,
			expectedError: "expected INDEX, got TEXT at line 1, column 1",
		},
		{
			input:         `00:00:01,123`,
			expectedError: "expected INDEX, got TIMESTAMP at line 1, column 1",
		},
		{
			input:         `1 00:00:01,123`,
			expectedError: "expected EOL, got TIMESTAMP at line 1, column 3",
		},
		{
			input: `1
-->`,
			expectedError: "expected TIMESTAMP, got ARROW at line 2, column 1",
		},
		{
			input: `1
00:00:01,123 --> 00:00:01,456
First line
Second line
2
00:00:01,123 --> 00:00:01,456
First line`,
			expectedError: "expected EOC, got INDEX at line 5, column 1",
		},
		{
			input: `1
00:00:01,123 --> 00:00:01,456
First line
Second line
-->`,
			expectedError: "expected EOC, got ARROW at line 5, column 1",
		},
		{
			input: `1
00:00:01,123 --> 00:00:01,456
First line
Second line
00:00:01,123`,
			expectedError: "expected EOC, got TIMESTAMP at line 5, column 1",
		},
		{
			input: `1
00:00:01,123 --> 00:00:01,456
First line
Second line

2
00:00:01,123 --> 00:00:01,456
First line

text`,
			expectedError: "expected INDEX, got TEXT at line 10, column 1",
		},
		{
			input: `1
00:00:01,123 --> 00:00:01,456
First line
Second line

2
00:00:01,123 --> 00:00:01,456
First line

3
00:00:01,123 --> 00:00:01,456
text
toto`,
			expectedError: "",
		},
	}

	for i, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		_, err := p.Parse()

		if tt.expectedError == "" {
			if err != nil {
				t.Fatalf("test[%d] - expected no error, got=%q.", i, err)
			}
		} else if err == nil || err.Error() != tt.expectedError {
			t.Fatalf("test[%d] - expected=%q, got=%q.", i, tt.expectedError, err)
		}
	}
}
