package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/florentsorel/subtitle"
	"github.com/florentsorel/subtitle/srt/internal/lexer"
	"github.com/florentsorel/subtitle/srt/internal/token"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	nextToken    token.Token
}

// New creates a new Parser instance and initializes it with the first two tokens.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	p.readToken()
	p.readToken()
	return p
}

// readToken advances the parser to the next token, updating currentToken and nextToken.
func (p *Parser) readToken() {
	p.currentToken = p.nextToken
	p.nextToken = p.lexer.NextToken()
}

// Parse processes the input tokens and returns a slice of Cue.
func (p *Parser) Parse() ([]subtitle.Cue, error) {
	var subs []subtitle.Cue

	for p.currentToken.Type != token.EOF {
		cue, err := p.parseCue()
		if err != nil {
			return nil, err
		}
		subs = append(subs, cue)
	}

	return subs, nil
}

// parseCue parses a single subtitle cue from the current token stream.
func (p *Parser) parseCue() (subtitle.Cue, error) {
	var s subtitle.Cue

	// Index
	if p.currentToken.Type != token.INDEX {
		return s, fmt.Errorf("expected INDEX, got %s at line %d, column %d", p.currentToken.Type, p.currentToken.Line, p.currentToken.Column)
	}
	index, err := strconv.Atoi(p.currentToken.Literal)
	if err != nil {
		return s, err
	}
	s.Index = index
	p.readToken()

	// End of line
	if p.currentToken.Type != token.EOL {
		return s, fmt.Errorf("expected EOL, got %s at line %d, column %d", p.currentToken.Type, p.currentToken.Line, p.currentToken.Column)
	}
	p.readToken()

	// Start time
	if p.currentToken.Type != token.TIMESTAMP {
		return s, fmt.Errorf("expected TIMESTAMP, got %s at line %d, column %d", p.currentToken.Type, p.currentToken.Line, p.currentToken.Column)
	}
	startTime, err := ParseTimestamp(p.currentToken.Literal)
	if err != nil {
		return s, err
	}
	s.Start = startTime

	p.readToken()

	// Arrow
	if p.currentToken.Type != token.ARROW {
		return s, fmt.Errorf("expected ARROW, got %s at line %d, column %d", p.currentToken.Type, p.currentToken.Line, p.currentToken.Column)
	}

	p.readToken()

	// End time
	if p.currentToken.Type != token.TIMESTAMP {
		return s, fmt.Errorf("expected TIMESTAMP, got %s at line %d, column %d", p.currentToken.Type, p.currentToken.Line, p.currentToken.Column)
	}
	endTime, err := ParseTimestamp(p.currentToken.Literal)
	if err != nil {
		return s, err
	}
	s.End = endTime

	p.readToken()

	// Text
	var textLines []string
	for p.currentToken.Type == token.TEXT || p.currentToken.Type == token.EOL {
		if p.currentToken.Type == token.TEXT {
			textLines = append(textLines, p.currentToken.Literal)
		}

		p.readToken()
	}

	s.Text = strings.Join(textLines, "\n")

	// End of cue
	if p.currentToken.Type != token.EOC && p.currentToken.Type != token.EOF {
		return s, fmt.Errorf("expected EOC, got %s at line %d, column %d", p.currentToken.Type, p.currentToken.Line, p.currentToken.Column)
	}

	p.readToken()

	return s, nil
}

func ParseTimestamp(s string) (time.Time, error) {
	normalized := strings.Replace(s, ",", ".", 1)
	const layout = "15:04:05.000"

	t, err := time.Parse(layout, normalized)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
