package tvecty

import (
	"bytes"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/stretchr/testify/require"
	"testing"
	"text/template"
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

func TestHTmlToDst_CanParseHtmlToAst1(t *testing.T) {
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

func TestHTmlToDst_ParsesAndInsertsMultivalueAttributesAttributes(t *testing.T) {
	htmlS := `<div class="some-class {myvar}">{s:"stuff"}</div>`
	expr, err := htmlToDst(htmlS)
	require.NoError(t, err)
	requireEqStr(t, tWrapExpr(t, expr), `
package thing

func RenderThing(msg string) vecty.HTMLOrComponent {
	elem.Div(
		vecty.Markup(
			vecty.Class("some-class", myvar),
		),
		vecty.Text("stuff"),
	)
}`)
}

func TestHTmlToDst_ParsesCustomAttributesAsSingleValue(t *testing.T) {
	htmlS := `<div data-ducks="this is all one argument">{s:"stuff"}</div>`
	expr, err := htmlToDst(htmlS)
	require.NoError(t, err)
	requireEqStr(t, tWrapExpr(t, expr), `
package thing

func RenderThing(msg string) vecty.HTMLOrComponent {
	elem.Div(
		vecty.Markup(
			vecty.Attribute("data-ducks", "this is all one argument"),
		),
		vecty.Text("stuff"),
	)
}`)
}

func TestHTmlToDst_ConvertsTagsWithSpecialNames(t *testing.T) {
	special := []struct {
		Tag string
		Out string
	}{
		{"h1", "Heading1"},
		{"h2", "Heading2"},
		{"h3", "Heading3"},
		{"h4", "Heading4"},
		{"h5", "Heading5"},
		{"h6", "Heading6"},
		{"nav", "Navigation"},
		{"img", "Image"},
		{"a", "Anchor"},
		{"hr", "HorizontalRule"},
		{"cite", "Citation"},
		//{"abbr", "Abbreviation"}, //?
	}
	inTpl := template.Must(template.New("in").Parse(`<div>
	<{{.Tag}}></{{.Tag}}>
</div>`))
	outTpl := template.Must(template.New("out").Parse(`
package thing

func RenderThing(msg string) vecty.HTMLOrComponent {
	elem.Div(
		elem.{{.Out}}(),
	)
}`))
	for _, s := range special {
		sinner := s
		t.Run("converts Tag "+s.Tag+" to "+s.Out, func(t *testing.T) {
			html := bytes.NewBuffer(nil)
			require.NoError(t, inTpl.Execute(html, sinner))
			expr, err := htmlToDst(html.String())
			require.NoError(t, err)
			expected := bytes.NewBuffer(nil)
			require.NoError(t, outTpl.Execute(expected, sinner))
			requireEqStr(t, tWrapExpr(t, expr), expected.String())

		})
	}
}

func TestHTmlToDst_UsesVectyTagWhenVectyDoesNotSupportTheTag(t *testing.T) {
	htmlS := `<customtag class="someclass">{s:"text"}</customtag>`
	expr, err := htmlToDst(htmlS)
	require.NoError(t, err)
	requireEqStr(t, tWrapExpr(t, expr), `
package thing

func RenderThing(msg string) vecty.HTMLOrComponent {
	vecty.Tag("customtag",
		vecty.Markup(
			vecty.Class("someclass"),
		),
		vecty.Text("text"),
	)
}`)
}

func TestHTmlToDst_SupportsTextInTags(t *testing.T) {
	htmlS := `<div>this is some text</div>`
	expr, err := htmlToDst(htmlS)
	require.NoError(t, err)
	requireEqStr(t, tWrapExpr(t, expr), `
package thing

func RenderThing(msg string) vecty.HTMLOrComponent {
	elem.Div(
		vecty.Text("this is some text"),
	)
}`)
}

func TestHTmlToDst_SupportsComplexTagTextEmbedding(t *testing.T) {
	htmlS := `<div>
this is some text
{firstValue} more text
and some here {secondValue}
{thirdValue} {fourthValue}
</div>`
	expr, err := htmlToDst(htmlS)
	require.NoError(t, err)
	requireEqStr(t, tWrapExpr(t, expr), `
package thing

func RenderThing(msg string) vecty.HTMLOrComponent {
	elem.Div(
		vecty.Text("this is some text"),
		firstValue,
		vecty.Text("more text"),
		vecty.Text("and some here"),
		secondValue,
		thirdValue,
		fourthValue,
	)
}`)
}

func TestHtmlToDst_SupportsEmbeddingMarkdownDirectly(t *testing.T) {
	htmlS := `<div markup="{someMap} {anotherMap}" class="first">text</div>`
	expr, err := htmlToDst(htmlS)
	require.NoError(t, err)
	requireEqStr(t, tWrapExpr(t, expr), `
package thing

func RenderThing(msg string) vecty.HTMLOrComponent {
	elem.Div(
		vecty.Markup(
			someMap,
			anotherMap,
			vecty.Class("first"),
		),
		vecty.Text("text"),
	)
}`)
}

func TestHtmlToDst_ParsesClassesSeparately(t *testing.T) {
	htmlS := `<div class="first second {third} fourth">text</div>`
	expr, err := htmlToDst(htmlS)
	require.NoError(t, err)
	requireEqStr(t, tWrapExpr(t, expr), `
package thing

func RenderThing(msg string) vecty.HTMLOrComponent {
	elem.Div(
		vecty.Markup(
			vecty.Class("first", "second", third, "fourth"),
		),
		vecty.Text("text"),
	)
}`)
}
