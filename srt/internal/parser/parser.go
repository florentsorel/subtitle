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

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	p.readToken()
	p.readToken()
	return p
}

func (p *Parser) readToken() {
	p.currentToken = p.nextToken
	p.nextToken = p.lexer.NextToken()
}

func (p *Parser) Parse() ([]subtitle.Subtitle, error) {
	var subs []subtitle.Subtitle

	for p.currentToken.Type != token.EOF {
		cue, err := p.parseCue()
		if err != nil {
			return nil, err
		}
		subs = append(subs, cue)
	}

	return subs, nil
}

func (p *Parser) parseCue() (subtitle.Subtitle, error) {
	var s subtitle.Subtitle

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
	startTime, err := parseTimestamp(p.currentToken.Literal)
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
	endTime, err := parseTimestamp(p.currentToken.Literal)
	if err != nil {
		return s, err
	}
	s.End = endTime

	p.readToken()

	var textLines []string
	for p.currentToken.Type != token.EOF && p.currentToken.Type != token.INDEX {
		switch p.currentToken.Type {
		case token.TEXT:
			textLines = append(textLines, p.currentToken.Literal)
		case token.EOL:
			if p.nextToken.Type == token.TEXT {
				textLines = append(textLines, "")
			}
		}
		p.readToken()
	}

	if p.currentToken.Type == token.INDEX {
		p.currentToken = p.nextToken
		p.nextToken = p.lexer.NextToken()
	}

	s.Text = strings.Join(textLines, "\n")

	return s, nil
}

func parseTimestamp(s string) (time.Time, error) {
	normalized := strings.Replace(s, ",", ".", 1)
	const layout = "15:04:05.000"

	t, err := time.Parse(layout, normalized)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
