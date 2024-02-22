package expr

import (
	"pg/tuple"
)

func NewOrExpression(lhs Expression, rhs Expression) Expression {
	return &orExpression{lhs: lhs, rhs: rhs}
}

type orExpression struct {
	lhs Expression
	rhs Expression
}

func (o orExpression) Execute(tup tuple.Tuple) bool {
	return o.lhs.Execute(tup) || o.rhs.Execute(tup)
}
