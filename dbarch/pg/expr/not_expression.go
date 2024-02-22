package expr

import (
	"pg/tuple"
)

func NewNotExpression(subExpression Expression) Expression {
	return &notExpression{
		subExpression: subExpression,
	}
}

type notExpression struct {
	subExpression Expression
}

func (n notExpression) Execute(tup tuple.Tuple) bool {
	return !n.subExpression.Execute(tup)
}
