package lexer

import (
	"unicode"

	"github.com/florentsorel/subtitle/srt/internal/token"
)

type Lexer struct {
	input        string
	position     int
	nextPosition int
	ch           rune
	line         int
	column       int
}

// New creates a new Lexer instance with the provided input string.
func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar()
	return l
}

// readChar reads the next character from the input and updates the position and line/column counters.
func (l *Lexer) readChar() {
	if l.nextPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = rune(l.input[l.nextPosition])
	}

	l.position = l.nextPosition
	l.nextPosition++

	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

// NextToken returns the next token from the input string.
func (l *Lexer) NextToken() token.Token {
	l.skipWhitespaceExpectNewLine()

	if l.ch == 0 {
		return token.Token{Type: token.EOF, Literal: "", Line: l.line, Column: 0}
	}

	switch {
	case unicode.IsDigit(l.ch):
		return l.lexNumberOrTimestamp()
	case l.ch == '\n' && l.nextChar() == '\n':
		line := l.line
		l.readChar()
		l.readChar()
		return token.Token{Type: token.EOC, Literal: "\n\n", Line: line, Column: 0}
	case l.ch == '\n':
		line := l.line
		l.readChar()
		return token.Token{Type: token.EOL, Literal: "\n", Line: line, Column: 0}
	case l.ch == '-' && l.nextChar() == '-' && l.peekAhead(2) == '>':
		startPos := l.position
		line := l.line
		col := l.column
		l.readChar()
		l.readChar()
		l.readChar()
		literal := l.input[startPos:l.position]
		return token.Token{Type: token.ARROW, Literal: literal, Line: line, Column: col}
	default:
		return l.lexText()
	}
}

// lexNumberOrTimestamp lexes a number or timestamp from the input string.
func (l *Lexer) lexNumberOrTimestamp() token.Token {
	startPos := l.position
	line := l.line
	col := l.column
	hasColon := false
	hasComma := false

	for unicode.IsDigit(l.ch) || l.ch == ':' || l.ch == ',' {
		if l.ch == ':' {
			hasColon = true
		}
		if l.ch == ',' {
			hasComma = true
		}
		l.readChar()
	}

	literal := l.input[startPos:l.position]
	if hasColon && hasComma {
		return token.Token{Type: token.TIMESTAMP, Literal: literal, Line: line, Column: col}
	}
	return token.Token{Type: token.INDEX, Literal: literal, Line: line, Column: col}
}

// lexText lexes a text string from the input string.
func (l *Lexer) lexText() token.Token {
	startPos := l.position
	line := l.line
	col := l.column

	for l.ch != 0 && l.ch != '\n' {
		l.readChar()
	}
	literal := l.input[startPos:l.position]
	return token.Token{Type: token.TEXT, Literal: literal, Line: line, Column: col}
}

// nextChar returns the next character in the input without advancing the position.
func (l *Lexer) nextChar() rune {
	if l.nextPosition >= len(l.input) {
		return 0
	}
	return rune(l.input[l.nextPosition])
}

// peekAhead returns the character at the specified offset from the current position without advancing the position.
func (l *Lexer) peekAhead(n int) rune {
	if l.nextPosition+n-1 >= len(l.input) {
		return 0
	}
	return rune(l.input[l.nextPosition+n-1])
}

// skipWhitespaceExpectNewLine skips whitespace characters until a newline character is encountered.
func (l *Lexer) skipWhitespaceExpectNewLine() {
	for unicode.IsSpace(l.ch) && l.ch != '\n' {
		l.readChar()
	}
}
