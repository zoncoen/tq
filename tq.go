package main

import (
	"flag"
	"fmt"
	"io"
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

func write(w io.Writer, t interface{}) error {
	if t == nil {
		_, err := fmt.Fprint(w, "\n")
		if err != nil {
			return err
		}
	} else {
		switch t.(type) {
		case int64:
			_, err := fmt.Fprintf(w, "%d\n", t)
			if err != nil {
				return err
			}
		case string:
			_, err := fmt.Fprintf(w, "\"%s\"\n", t)
			if err != nil {
				return err
			}
		case []map[string]interface{}:
			e := toml.NewEncoder(w)
			s, _ := t.([]map[string]interface{})
			for _, v := range s {
				err := e.Encode(v)
				if err != nil {
					return err
				}
			}
		case []interface{}:
			s, _ := t.([]interface{})
			for _, v := range s {
				err := write(w, v)
				if err != nil {
					return err
				}
			}
		default:
			e := toml.NewEncoder(os.Stdout)
			err := e.Encode(t)
			if err != nil {
				return err
			}
		}
	}
	return nil
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

	err = write(os.Stdout, t)
	if err != nil {
		log.Fatalf("%s", err)
	}
}
