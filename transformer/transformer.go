package transformer

import (
	"errors"

	"github.com/zoncoen/tq/ast"
)

func Transform(i interface{}, f ast.Filter) (interface{}, error) {
	return Filter(i, f)
}

func Filter(i interface{}, f ast.Filter) (res interface{}, err error) {
	switch f.(type) {
	case ast.EmptyFilter:
		res = i
	case ast.KeyFilter:
		kf, _ := f.(ast.KeyFilter)
		res, err = FilterByKey(i, kf)
	case ast.BinaryOp:
		bo, _ := f.(ast.BinaryOp)
		switch bo.Op.Literal {
		case "|":
			res, err = Filter(i, bo.Left)
			if err != nil {
				return res, err
			}
			res, err = Filter(res, bo.Right)
		}
	}
	return res, err
}

func FilterByKey(i interface{}, f ast.KeyFilter) (interface{}, error) {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil, errors.New("parse error: Objects must consist of key:value pairs")
	}
	return m[f.Key], nil
}
