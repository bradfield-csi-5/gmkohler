package scanner

type TokenType int

func TokenTypeFromKeyword(kw string) TokenType {
	switch kw {
	case "OR":
		return TokenTypeOr
	case "AND":
		return TokenTypeAnd
	case "NOT":
		return TokenTypeNot
	default:
		return TokenTypeUnrecognized
	}
}

func (t TokenType) String() string {
	switch t {
	case TokenTypeOr:
		return "OR"
	case TokenTypeAnd:
		return "AND"
	case TokenTypeNot:
		return "NOT"
	case TokenTypeIdentifier:
		return "IDENTIFIER"
	case TokenTypeEof:
		return "EOF"
	default:
		return "unrecognized"
	}
}

const (
	// literals
	TokenTypeUnrecognized TokenType = iota
	// keywords
	TokenTypeOr
	TokenTypeAnd
	TokenTypeNot
	TokenTypeIdentifier
	TokenTypeEof
)

type Token struct {
	tType   TokenType
	lexeme  string
	literal any
	line    int
}
