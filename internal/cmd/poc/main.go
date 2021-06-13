package main

import (
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/mdev5000/tvecty"
	"go/token"
	"strconv"
)

//type filler struct {
//	line  int
//	ftype string
//	value string
//}
//
//func ignore(i interface{}) {
//}

func stringLit(s string) *dst.BasicLit {
	return &dst.BasicLit{
		Kind:  token.STRING,
		Value: s,
		Decs:  dst.BasicLitDecorations{},
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

type fakeReplace struct {
}

func (fakeReplace) Get(id int) (dst.Expr, bool) {
	return simpleCallExpr("replaced", "This", []dst.Expr{stringLit(strconv.Itoa(id))}), true
}

func main() {
	code := `package thing

var thing = tvecty.Html(1, "another")

func RenderThing(msg string) vecty.HTMLOrComponent {
	return tvecty.Html(2, "<div>test</div>")
}

func Another() vecty.HTMLOrComponent {
	var thing vecty.HTMLOrComponent = tvecty.Html(3, "<div>another</div>")
	thing2 := tvecty.Html(4, "<div>another</div>")
}

func List() []vecty.HTMLOrComponent {
	var arr []vecty.HTMLOrComponent
	for i := 0; i < 10; i++ {
		arr[0] = tvecty.Html(5, "<div>another</div>")
	}
	return i
}

func InFunctionCall() []vecty.HTMLOrComponent {
	return someCall(tvecty.Html(5, "<div>another</div>"))
}
`
	f, err := decorator.Parse(code)
	if err != nil {
		panic(err)
	}
	//spew.Dump(f.Decls[4])

	err = tvecty.Replace(fakeReplace{}, f)
	if err != nil {
		panic(err)
	}

	//if err := tvecty.FinishHtmlFuncDefinitions(f); err != nil {
	//	panic(err)
	//}
	////fmt.Println(len(f.Decls))
	//
	if err := decorator.Print(f); err != nil {
		panic(err)
	}
}
