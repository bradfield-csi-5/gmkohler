package expr

import (
	"pg/tuple"
)

type Expression interface {
	Execute(tuple.Tuple) bool
}
