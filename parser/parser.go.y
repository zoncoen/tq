%{
package parser

import (
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

%token <token> PERIOD STRING PIPE

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
    | filter PIPE filter
    {
        $$ = ast.BinaryOp{Left: $1, Op: $2, Right: $3}
    }

%%
