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
00:00:01,123 00:00:01,456`,
			expectedError: "expected ARROW, got TIMESTAMP at line 2, column 14",
		},
		{
			input: `1
00:00:01,123 --> hello`,
			expectedError: "expected TIMESTAMP, got TEXT at line 2, column 18",
		},
		{
			input: `1
00:00:01,123 --> 00:00:01,456
First line
Second line

-->`,
			expectedError: "expected INDEX, got ARROW at line 6, column 1",
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
First line`,
			expectedError: "expected EOL (EOB), got TIMESTAMP at line 5, column 1",
		},
	}

	for i, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		_, err := p.Parse()

		if err == nil || err.Error() != tt.expectedError {
			t.Fatalf("test[%d] - expected=%q, got=%q.", i, tt.expectedError, err)
		}
	}
}
