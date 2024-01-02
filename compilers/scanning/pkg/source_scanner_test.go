package scanner

import (
	"slices"
	"testing"
)

type testCase struct {
	src      string
	expected []Token
}

var testCases = []testCase{
	{
		src: "hello AND world OR alice AND NOT bob",
		expected: []Token{
			{
				tType:  TokenTypeIdentifier,
				lexeme: "hello",
				line:   1,
			},
			{
				tType:  TokenTypeAnd,
				lexeme: "AND",
				line:   1,
			},
			{
				tType:  TokenTypeIdentifier,
				lexeme: "world",
				line:   1,
			},
			{
				tType:  TokenTypeOr,
				lexeme: "OR",
				line:   1,
			},
			{
				tType:  TokenTypeIdentifier,
				lexeme: "alice",
				line:   1,
			},
			{
				tType:  TokenTypeAnd,
				lexeme: "AND",
				line:   1,
			},
			{
				tType:  TokenTypeNot,
				lexeme: "NOT",
				line:   1,
			},
			{
				tType:  TokenTypeIdentifier,
				lexeme: "bob",
				line:   1,
			},
			{
				tType: TokenTypeEof,
				line:  1,
			},
		},
	},
}

func TestSourceScanner_ScanTokens(t *testing.T) {
	for _, tc := range testCases {
		s := NewSourceScanner(tc.src)
		tokens := s.ScanTokens()
		if !slices.Equal(tokens, tc.expected) {
			t.Fatalf("\nExpected: %+v\nGot: %+v", tc.expected, tokens)
		}
	}

}
