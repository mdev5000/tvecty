package tvecty

import (
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"go/token"
)

func stringLit(s string) *dst.BasicLit {
	return &dst.BasicLit{
		Kind:  token.STRING,
		Value: s,
		Decs:  dst.BasicLitDecorations{},
	}
}

func simpleSelectorExpr(x, sel string) *dst.SelectorExpr {
	return &dst.SelectorExpr{
		X:   dst.NewIdent(x),
		Sel: dst.NewIdent(sel),
		Decs: dst.SelectorExprDecorations{
			NodeDecs: dst.NodeDecs{
				Before: dst.SpaceType(1),
				Start:  nil,
				End:    nil,
				After:  dst.SpaceType(1),
			},
		},
	}
}

func simpleCallExpr(x, sel string, args []dst.Expr) *dst.CallExpr {
	return &dst.CallExpr{
		Decs: dst.CallExprDecorations{
			NodeDecs: dst.NodeDecs{
				Before: dst.SpaceType(1),
				Start:  nil,
				End:    nil,
				After:  dst.SpaceType(1),
			},
		},
		Fun: &dst.SelectorExpr{
			X:   dst.NewIdent(x),
			Sel: dst.NewIdent(sel),
			Decs: dst.SelectorExprDecorations{
				NodeDecs: dst.NodeDecs{
					Before: dst.SpaceType(1),
					Start:  nil,
					End:    nil,
					After:  dst.SpaceType(1),
				},
			},
		},
		Args:     args,
		Ellipsis: false,
	}
}

func parseExpression(exprStr string) (dst.Expr, error) {
	// Basically cheat and make a mini go file for dst to parse.
	v, err := decorator.Parse(fmt.Sprintf("package tmp; var e = %s", exprStr))
	if err != nil {
		return nil, fmt.Errorf("error with expression '%s':\n%w", exprStr, err)
	}
	gd, ok := v.Decls[0].(*dst.GenDecl)
	if !ok {
		panic("failed to convert GenDecl")
	}
	vs, ok := gd.Specs[0].(*dst.ValueSpec)
	if !ok {
		panic("failed to convert ValueSpec")
	}
	return vs.Values[0], nil
}
