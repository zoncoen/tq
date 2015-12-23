%{
package parser

import (
    "io"

    "github.com/zoncoen/tq/ast"
    "github.com/zoncoen/tq/token"
)
%}

%union{
    token token.Token
    expr  ast.Filter
}

%type <expr> program
%type <expr> filter
%type <expr> key_filter
%type <expr> index_filter
%type <expr> range_filter
%type <expr> binary_op

%token <token> PERIOD STRING INT PIPE LBRACK RBRACK COLON COMMA QUESTION

%left PIPE

%%

program
    : filter
    {
        $$ = $1
        yylex.(*Lexer).result = $$
    }

filter
    : PERIOD
    {
        $$ = ast.EmptyFilter{}
    }
    | key_filter
    {
        $$ = $1
    }
    | index_filter
    {
        $$ = $1
    }
    | range_filter
    {
        $$ = $1
    }
    | binary_op
    {
        $$ = $1
    }
    | filter QUESTION
    {
        $$ = ast.IgnoreErrorHandler{Filter: $1}
    }

key_filter
    : PERIOD STRING
    {
        $$ = ast.KeyFilter{Key: $2.Literal}
    }
    | PERIOD LBRACK STRING RBRACK
    {
        $$ = ast.KeyFilter{Key: $3.Literal}
    }

index_filter
    : PERIOD LBRACK INT RBRACK
    {
        $$ = ast.IndexFilter{Index: $3.Literal}
    }

range_filter
    : PERIOD LBRACK INT COLON INT RBRACK
    {
        $$ = ast.RangeFilter{Low: $3.Literal, High: $5.Literal}
    }
    | PERIOD LBRACK INT COLON RBRACK
    {
        $$ = ast.RangeFilter{Low: $3.Literal, High: ""}
    }
    | PERIOD LBRACK COLON INT RBRACK
    {
        $$ = ast.RangeFilter{Low: "", High: $4.Literal}
    }
    | PERIOD LBRACK RBRACK
    {
        $$ = ast.RangeFilter{Low: "", High: ""}
    }

binary_op
    : filter PIPE filter
    {
        $$ = ast.BinaryOp{Left: $1, Op: $2, Right: $3}
    }
    | filter PERIOD STRING
    {
        $$ = ast.BinaryOp{Left: $1, Op: token.Token{Token: PIPE, Literal: "|"}, Right: ast.KeyFilter{Key: $3.Literal}}
    }
    | filter COMMA filter
    {
        $$ = ast.BinaryOp{Left: $1, Op: $2, Right: $3}
    }
%%

func Parse(r io.Reader) ast.Filter {
    l := new(Lexer)
    l.Init(r)
    yyParse(l)
    return l.result
}
