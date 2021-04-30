package tvecty

import (
	"bytes"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/stretchr/testify/require"
	"testing"
)

// Replaces the body contents with the expr passed to this function.
func tWrapExpr(t *testing.T, expr dst.Expr) string {
	mn := `
package thing

func RenderThing(msg string) vecty.HTMLOrComponent {
}
`
	ast, err := decorator.Parse(mn)
	require.NoError(t, err)
	fd := ast.Decls[0].(*dst.FuncDecl)
	fd.Body = &dst.BlockStmt{
		List: []dst.Stmt{
			&dst.ExprStmt{
				X: expr,
				Decs: dst.ExprStmtDecorations{
					NodeDecs: dst.NodeDecs{
						Before: dst.SpaceType(1),
						Start:  nil,
						End:    nil,
						After:  dst.SpaceType(1),
					},
				},
			},
		},
		RbraceHasNoPos: false,
		Decs:           dst.BlockStmtDecorations{},
	}
	b := bytes.NewBuffer(nil)
	require.NoError(t, decorator.Fprint(b, ast))
	return b.String()
}

func TestCanParseHtmlToAst1(t *testing.T) {
	htmlS := `<div>
	{s:msg}
	{something()}
</div>`
	expr, err := htmlToDst(htmlS)
	require.NoError(t, err)
	requireEqStr(t, tWrapExpr(t, expr), `
package thing

func RenderThing(msg string) vecty.HTMLOrComponent {
	elem.Div(
		vecty.Text(msg),
		something())
}`)
}
