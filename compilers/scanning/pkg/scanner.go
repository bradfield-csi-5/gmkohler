package scanner

type Scanner interface {
	ScanTokens() []Token
}
