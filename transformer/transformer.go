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
	case ast.RangeFilter:
		rf, _ := f.(ast.RangeFilter)
		res, err = FilterByRange(i, rf)
	case ast.BinaryOp:
		bo, _ := f.(ast.BinaryOp)
		res, err = ExecuteBinaryOp(i, bo)
	case ast.IgnoreErrorHandler:
		h, _ := f.(ast.IgnoreErrorHandler)
		res, err = Filter(i, h.Filter)
		if err != nil {
			return []map[string]interface{}{}, nil
		}
	}
	return res, err
}

func FilterByKey(i interface{}, f ast.KeyFilter) (interface{}, error) {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil, errors.New("transform error: Object must consist of key:value pairs")
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
		return nil, errors.New("transform error: Object must consist of array")
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

func FilterByRange(i interface{}, f ast.RangeFilter) (interface{}, error) {
	a, ok := i.([]map[string]interface{})
	if !ok {
		return nil, errors.New("transform error: Object must consist of array")
	}

	l := 0
	h := len(a) + 1
	if f.Low != "" {
		n, err := strconv.Atoi(f.Low)
		if err != nil {
			return nil, err
		}
		if n > l {
			l = n
		}
	}
	if f.High != "" {
		n, err := strconv.Atoi(f.High)
		if err != nil {
			return nil, err
		}
		if n < h {
			h = n
		}
	}

	if l > h {
		return []map[string]interface{}{}, nil
	}

	return a[l:h], nil
}

func ExecuteBinaryOp(i interface{}, bo ast.BinaryOp) (res interface{}, err error) {
	switch bo.Op.Literal {
	case "|":
		res, err = Filter(i, bo.Left)
		if err != nil {
			return res, err
		}
		res, err = Filter(res, bo.Right)
	case ",":
		lRes, err := Filter(i, bo.Left)
		if err != nil {
			return res, err
		}
		rRes, err := Filter(i, bo.Right)
		if err != nil {
			return res, err
		}
		res = []interface{}{lRes, rRes}
	}
	return res, err
}
