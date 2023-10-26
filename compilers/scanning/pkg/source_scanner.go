package scanner

import (
	"fmt"
	"os"
)

type sourceScanner struct {
	source  string // consider some sort of reader
	tokens  []Token
	start   int
	current int
	line    int
}

func NewSourceScanner(source string) Scanner {
	return &sourceScanner{
		source: source,
		line:   1,
	}
}

func (s *sourceScanner) ScanTokens() []Token {
	for !s.atEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{
		tType: TokenTypeEof,
		line:  s.line,
	})

	return s.tokens
}

func (s *sourceScanner) scanToken() {
	switch r := s.advance(); r {
	case '\n':
		s.line++
	case ' ', '\r', '\t':
	default:
		if isAlpha(r) {
			s.identifier()
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "unexpected character %c\n", r)
		}
	}
}

func (s *sourceScanner) addToken(tt TokenType) {
	s.addTokenWithLiteral(tt, nil)
}

func (s *sourceScanner) addTokenWithLiteral(tt TokenType, literal any) {
	s.tokens = append(
		s.tokens,
		Token{
			tType:   tt,
			lexeme:  s.source[s.start:s.current],
			literal: literal,
			line:    s.line,
		},
	)
}

func (s *sourceScanner) advance() rune {
	r := rune(s.source[s.current])
	s.current++
	return r
}

func (s *sourceScanner) match(expected rune) bool {
	if s.atEnd() {
		return false
	}
	if rune(s.source[s.current]) != expected {
		return false
	}
	s.current++
	return true
}

func (s *sourceScanner) identifier() {
	for isAlpha(s.peek()) {
		s.advance()
	}
	tt := TokenTypeFromKeyword(s.source[s.start:s.current])
	if tt == TokenTypeUnrecognized {
		tt = TokenTypeIdentifier
	}
	s.addToken(tt)
}

func (s *sourceScanner) peek() rune {
	if s.atEnd() {
		return rune(0)
	}
	return rune(s.source[s.current])
}

func (s *sourceScanner) atEnd() bool {
	return s.current >= len(s.source)
}
