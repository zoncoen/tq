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
	{".\"key\"", ast.KeyFilter{Key: "key"}},
	{".[\"key\"]", ast.KeyFilter{Key: "key"}},
	{".[0]", ast.IndexFilter{Index: "0"}},
	{".[0:1]", ast.RangeFilter{Low: "0", High: "1"}},
	{".[0:]", ast.RangeFilter{Low: "0", High: ""}},
	{".[:1]", ast.RangeFilter{Low: "", High: "1"}},
	{".[]", ast.RangeFilter{Low: "", High: ""}},
	{". | .", ast.BinaryOp{
		Left:  ast.EmptyFilter{},
		Op:    token.Token{Token: PIPE, Literal: "|"},
		Right: ast.EmptyFilter{}}},
	{".first.second", ast.BinaryOp{
		Left:  ast.KeyFilter{Key: "first"},
		Op:    token.Token{Token: PIPE, Literal: "|"},
		Right: ast.KeyFilter{Key: "second"}}},
	{".first.second.third", ast.BinaryOp{
		Left: ast.BinaryOp{
			Left:  ast.KeyFilter{Key: "first"},
			Op:    token.Token{Token: PIPE, Literal: "|"},
			Right: ast.KeyFilter{Key: "second"}},
		Op:    token.Token{Token: PIPE, Literal: "|"},
		Right: ast.KeyFilter{Key: "third"}}},
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
