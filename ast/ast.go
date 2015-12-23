package ast

import (
	"github.com/zoncoen/tq/token"
)

type Filter interface{}

type EmptyFilter struct {
}

type KeyFilter struct {
	Key string
}

type IndexFilter struct {
	Index string
}

type RangeFilter struct {
	Low  string
	High string
}

type BinaryOp struct {
	Left  Filter
	Op    token.Token
	Right Filter
}

type IgnoreErrorHandler struct {
	Filter Filter
}
