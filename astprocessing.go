package tvecty

import (
	"github.com/dave/dst"
	"strings"
)

// FinishHtmlFuncDefinitions trims the comments indicating a function is an html
// function and adds the correct return type.
func FinishHtmlFuncDefinitions(f *dst.File) error {
	for _, d := range f.Decls {
		df, ok := d.(*dst.FuncDecl)
		if !ok {
			continue
		}
		comments := df.Decs.NodeDecs.Start
		var lastComment string
		if len(comments) == 0 {
			lastComment = ""
		} else {
			lastComment = comments[len(comments)-1]
		}
		if strings.HasPrefix(lastComment, "/*!!htmlfunc") {
			df.Decs.NodeDecs.Start = comments[:len(comments)-1]
			if df.Type.Results == nil {
				df.Type.Results = &dst.FieldList{
					Opening: false,
					List:    nil,
					Closing: false,
					Decs:    dst.FieldListDecorations{},
				}
			}
			df.Type.Results.List = append(df.Type.Results.List, &dst.Field{
				Names: nil,
				Type: &dst.SelectorExpr{
					X:    dst.NewIdent("vecty"),
					Sel:  dst.NewIdent("HTMLOrComponent"),
					Decs: dst.SelectorExprDecorations{},
				},
				Tag:  nil,
				Decs: dst.FieldDecorations{},
			})
		}
	}
	return nil
}
