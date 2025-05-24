package srt

import (
	"github.com/florentsorel/subtitle"
	"github.com/florentsorel/subtitle/srt/internal/lexer"
	"github.com/florentsorel/subtitle/srt/internal/parser"
)

// Parse parses the input string and returns a slice of subtitle.Cue.
func Parse(input string) ([]subtitle.Cue, error) {
	l := lexer.New(input)
	p := parser.New(l)
	return p.Parse()
}
