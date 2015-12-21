package parser

import (
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
	if tok == scanner.Ident {
		tok = STRING
	}
	if tok == scanner.String {
		tok = STRING
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
	lval.token = token.Token{Token: tok, Literal: l.TokenText()}
	return tok
}

func (l *Lexer) Error(e string) {
	panic(e)
}
