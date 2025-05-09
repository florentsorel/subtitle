package lexer

import (
	"testing"

	"github.com/florentsorel/subtitle/srt/internal/token"
)

func TestNewLexer(t *testing.T) {
	input := `1
00:00:01,123 --> 00:00:01,456
Text block éè
Second line
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.INDEX, "1", 1, 1},
		{token.EOL, "\n", 2, 0},
		{token.TIMESTAMP, "00:00:01,123", 2, 1},
		{token.ARROW, "-->", 2, 14},
		{token.TIMESTAMP, "00:00:01,456", 2, 18},
		{token.EOL, "\n", 3, 0},
		{token.TEXT, "Text block éè", 3, 1},
		{token.EOL, "\n", 4, 0},
		{token.TEXT, "Second line", 4, 1},
		{token.EOL, "\n", 5, 0},
		{token.EOF, "", 5, 0},
	}

	lexer := New(input)

	for i, tt := range tests {
		tok := lexer.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - TokenType. expected=%q, got=%q.", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal. expected=%q, got=%q.", i, tt.expectedLiteral, tok.Literal)
		}

		if tok.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - Line. expected=%d, got=%d.", i, tt.expectedLine, tok.Line)
		}

		if tok.Column != tt.expectedColumn {
			t.Fatalf("tests[%d] - Column. expected=%d, got=%d.", i, tt.expectedColumn, tok.Column)
		}
	}
}
