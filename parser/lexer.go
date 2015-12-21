package parser

import (
	"strconv"
	"text/scanner"

	"github.com/zoncoen/tq/ast"
	"github.com/zoncoen/tq/token"
)

type Lexer struct {
	scanner.Scanner
	result ast.Filter
}

func (l *Lexer) Lex(lval *yySymType) int {
	tok := int(l.Scan())
	lit := l.TokenText()
	if tok == scanner.Ident {
		tok = STRING
	}
	if tok == scanner.String {
		tok = STRING
		lit, _ = strconv.Unquote(lit)
	}
	if tok == scanner.Int {
		tok = INT
	}
	if tok == int('.') {
		tok = PERIOD
	}
	if tok == int('|') {
		tok = PIPE
	}
	if tok == int('[') {
		tok = LBRACK
	}
	if tok == int(']') {
		tok = RBRACK
	}
	if tok == int(':') {
		tok = COLON
	}
	lval.token = token.Token{Token: tok, Literal: lit}
	return tok
}

func (l *Lexer) Error(e string) {
	panic(e)
}
