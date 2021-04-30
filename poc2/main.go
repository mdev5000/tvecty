package main

import (
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	//	code := `package thing
	//
	//
	//func (c *Component) List() []vecty.HTMLOrComponent {
	//	var arr []vecty.HTMLOrComponent
	//	for i := 0; i < 10; i++ {
	//		e := elem.Div(
	//			elem.Div(
	//				vecty.Markup(
	//					vecty.Class("container", c.MyClass),
	//					event.Click(c.MyFunc),
	//				),
	//				comps.NewNavBar(),
	//				p.child,
	//			),
	//		)
	//		arr = append(arr, e)
	//	}
	//	return arr
	//}
	//`
	//	f, err := decorator.Parse(code)
	//	if err != nil {
	//		panic(err)
	//	}
	//	for _, decl := range f.Decls {
	//		fd, ok := decl.(*dst.FuncDecl)
	//		if !ok {
	//			continue
	//		}
	//		spew.Dump(fd.Body)
	//	}
	v, err := parseExpression("first.Thing()")
	if err != nil {
		panic(err)
	}
	spew.Dump(v)
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
