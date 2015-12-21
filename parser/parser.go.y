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

%token <token> PERIOD STRING INT PIPE LBRACK RBRACK

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
    | PERIOD STRING
    {
        $$ = ast.KeyFilter{Key: $2.Literal}
    }
    | PERIOD LBRACK INT RBRACK
    {
        $$ = ast.IndexFilter{Index: $3.Literal}
    }
    | filter PIPE filter
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
