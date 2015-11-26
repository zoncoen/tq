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

type BinaryOp struct {
	Left  Filter
	Op    token.Token
	Right Filter
}
