package srt

import (
	"github.com/florentsorel/subtitle"
	"github.com/florentsorel/subtitle/srt/internal/lexer"
	"github.com/florentsorel/subtitle/srt/internal/parser"
)

func Parse(input string) ([]subtitle.Subtitle, error) {
	l := lexer.New(input)
	p := parser.New(l)

	return p.Parse()
}
