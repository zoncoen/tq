package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/zoncoen/tq/ast"
	"github.com/zoncoen/tq/parser"
	"github.com/zoncoen/tq/transformer"
)

type option struct {
	help bool
}

var opt option

func init() {
	flag.BoolVar(&opt.help, "h", false, "print help message")
	flag.BoolVar(&opt.help, "help", false, "print help message")
}

func main() {
	flag.Parse()

	if opt.help {
		flag.Usage()
		os.Exit(0)
	}

	var t interface{}
	_, err := toml.DecodeReader(os.Stdin, &t)
	if err != nil {
		log.Fatalf("parse error: %s", err)
	}

	var filter ast.Filter
	if flag.NArg() > 0 {
		str := flag.Arg(0)
		r := strings.NewReader(str)
		filter = parser.Parse(r)
		t, err = transformer.Transform(t, filter)
		if err != nil {
			log.Fatalf("%s", err)
		}
	}

	if t == nil {
		fmt.Print("\n")
	} else {
		switch t.(type) {
		case int64:
			fmt.Printf("%d\n", t)
		case string:
			fmt.Printf("\"%s\"\n", t)
		default:
			e := toml.NewEncoder(os.Stdout)
			err = e.Encode(t)
			if err != nil {
				log.Fatalf("%s", err)
			}
		}
	}
}
