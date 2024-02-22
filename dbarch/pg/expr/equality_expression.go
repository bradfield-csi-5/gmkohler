package expr

import "pg/tuple"

func NewEqualityExpression(field string, value tuple.ColumnValue) Expression {
	return &equalityExpression{
		field: field,
		value: value,
	}
}

type equalityExpression struct {
	field string
	value tuple.ColumnValue
}

func (e equalityExpression) Execute(tup tuple.Tuple) bool {
	value, err := tup.GetColumnValue(e.field)
	if err != nil {
		panic(err) // or false
	}
	return value == e.value
}
