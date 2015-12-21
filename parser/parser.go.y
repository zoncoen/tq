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
    | key_filter
    {
        $$ = $1
    }
    | index_filter
    {
        $$ = $1
    }
    | filter PIPE filter
    {
        $$ = ast.BinaryOp{Left: $1, Op: $2, Right: $3}
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
%%

func Parse(r io.Reader) ast.Filter {
    l := new(Lexer)
    l.Init(r)
    yyParse(l)
    return l.result
}
