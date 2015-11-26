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
	{". | .", ast.BinaryOp{
		Left:  ast.EmptyFilter{},
		Op:    token.Token{Token: PIPE, Literal: "|"},
		Right: ast.EmptyFilter{}}},
}

func TestYyParse(t *testing.T) {
	for i, test := range parseTests {
		r := strings.NewReader(test.text)
		l := new(Lexer)
		l.Init(r)
		yyParse(l)
		if l.result != test.ast {
			t.Errorf("case %d: got %#v; expected %#v", i, l.result, test.ast)
		}
	}
}
