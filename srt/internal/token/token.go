package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

const (
	INDEX     = "INDEX"
	TIMESTAMP = "TIMESTAMP"
	ARROW     = "ARROW"
	TEXT      = "TEXT"
	EOL       = "EOL"
	EOC       = "EOC"
	EOF       = "EOF"
)
