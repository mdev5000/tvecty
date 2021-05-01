package tvecty

import (
	"fmt"
	"github.com/dave/dst"
	"strconv"
)

type ReplacementFinder interface {
	Get(id int) (dst.Expr, bool)
}

func Replace(h ReplacementFinder, f *dst.File) error {
	var err error
	dst.Inspect(f, func(node dst.Node) bool {
		if err != nil {
			return false
		}
		switch expr := node.(type) {
		case *dst.ValueSpec:
			if ierr := tryConvertExprsToVectyCall(h, expr.Values); ierr != nil {
				err = ierr
				return false
			}
		case *dst.ReturnStmt:
			if ierr := tryConvertExprsToVectyCall(h, expr.Results); ierr != nil {
				err = ierr
				return false
			}
		case *dst.AssignStmt:
			if ierr := tryConvertExprsToVectyCall(h, expr.Rhs); ierr != nil {
				err = ierr
				return false
			}
		}
		return true
	})
	return err
}

func tryConvertExprsToVectyCall(h ReplacementFinder, exprs []dst.Expr) error {
	for i, expr := range exprs {
		outExpr, err := tryConvertTVectyCall(h, expr)
		if err != nil {
			return err
		}
		if outExpr == nil {
			continue
		}
		exprs[i] = outExpr
	}
	return nil
}

func tryConvertTVectyCall(h ReplacementFinder, expr dst.Expr) (dst.Expr, error) {
	cExpr, ok := expr.(*dst.CallExpr)
	if !ok {
		return nil, nil
	}
	sExpr, ok := cExpr.Fun.(*dst.SelectorExpr)
	if !ok {
		return nil, nil
	}
	x, ok := sExpr.X.(*dst.Ident)
	if !ok || x.Name != "tvecty" {
		return nil, nil
	}
	if sExpr.Sel.Name != "Html" {
		return nil, nil
	}
	if len(cExpr.Args) < 2 {
		return nil, fmt.Errorf("invalid tvecty.Html call, must have 2 arguments:\n%v", cExpr)
	}
	idLit, ok := cExpr.Args[0].(*dst.BasicLit)
	if !ok {
		return nil, fmt.Errorf("invalid tvecty.Html call, expected arg 1 to be a BasicLit:\n%v", cExpr)
	}
	id, errAtoi := strconv.Atoi(idLit.Value)
	if errAtoi != nil {
		return nil, fmt.Errorf("invalid tvecty.Html call, expected arg 1 to be an int:\n%v", cExpr)
	}
	replacement, found := h.Get(id)
	if !found {
		return nil, fmt.Errorf("invalid tvecty.Html call, html id %d not found:\n%v", id, cExpr)
	}
	return replacement, nil
}
