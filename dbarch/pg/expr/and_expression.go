package expr

import (
	"pg/tuple"
)

func NewAndExpression(lhs Expression, rhs Expression) Expression {
	return &andExpression{
		lhs: lhs,
		rhs: rhs,
	}
}

type andExpression struct {
	lhs Expression
	rhs Expression
}

func (a andExpression) Execute(tup tuple.Tuple) bool {
	return a.lhs.Execute(tup) && a.rhs.Execute(tup)
}
