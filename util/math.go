package util

import "golang.org/x/exp/constraints"

//type Integer constraints.Integer
//type Float constraints.Float

type Numeric interface {
	constraints.Integer | constraints.Float
}
