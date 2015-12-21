package parser

import (
	"strings"
	"testing"

	"github.com/zoncoen/tq/ast"
	"github.com/zoncoen/tq/token"
)

var parseTests = []struct {
	text string
	ast  ast.Filter
}{
	{".", ast.EmptyFilter{}},
	{".key", ast.KeyFilter{Key: "key"}},
	{".[0]", ast.IndexFilter{Index: "0"}},
	{". | .", ast.BinaryOp{
		Left:  ast.EmptyFilter{},
		Op:    token.Token{Token: PIPE, Literal: "|"},
		Right: ast.EmptyFilter{}}},
}

func TestParse(t *testing.T) {
	for i, test := range parseTests {
		r := strings.NewReader(test.text)
		res := Parse(r)
		if res != test.ast {
			t.Errorf("case %d: got %#v; expected %#v", i, res, test.ast)
		}
	}
}
