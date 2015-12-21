package transformer

import (
	"errors"
	"strconv"

	"github.com/zoncoen/tq/ast"
)

func Transform(i interface{}, f ast.Filter) (interface{}, error) {
	return Filter(i, f)
}

func Filter(i interface{}, f ast.Filter) (res interface{}, err error) {
	if i == nil {
		return nil, nil
	}
	switch f.(type) {
	case ast.EmptyFilter:
		res = i
	case ast.KeyFilter:
		kf, _ := f.(ast.KeyFilter)
		res, err = FilterByKey(i, kf)
	case ast.IndexFilter:
		inf, _ := f.(ast.IndexFilter)
		res, err = FilterByIndex(i, inf)
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
		return nil, errors.New("transform error: Objects must consist of key:value pairs")
	}
	v, ok := m[f.Key]
	if !ok {
		return nil, nil
	}
	return v, nil
}

func FilterByIndex(i interface{}, f ast.IndexFilter) (interface{}, error) {
	a, ok := i.([]map[string]interface{})
	if !ok {
		return nil, errors.New("transform error: Objects must consist of array")
	}
	index, err := strconv.Atoi(f.Index)
	if err != nil {
		return nil, err
	}
	if index >= len(a) {
		return nil, nil
	}
	return a[index], nil
}
