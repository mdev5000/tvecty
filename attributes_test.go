package tvecty

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTokenizeParts_CorrectlyTokenizesParts(t *testing.T) {
	parts, err := tokenizeParts(`{first} Second value {third}
And some more {fourth} stuff
And another thing of text
{fifth}
`)
	require.NoError(t, err)
	require.Equal(t, parts, []attributeToken{
		{"first", true},
		{"Second value", false},
		{"third", true},
		{"And some more", false},
		{"fourth", true},
		{"stuff", false},
		{"And another thing of text", false},
		{"fifth", true},
	})
}

func TestParseExpressionWrapper2_CanParseExpressions(t *testing.T) {
	expr, err := parseExpressionWrapper2("first", true)
	require.NoError(t, err)
	requireEqStr(t, tWrapExpr(t, expr), `
package thing

func RenderThing(msg string) vecty.HTMLOrComponent {
	first
}`)
}

func TestParseExpressionWrapper2_CanParseExpressionsWithStringModifiers(t *testing.T) {
	expr, err := parseExpressionWrapper2("s:wrapped", true)
	require.NoError(t, err)
	requireEqStr(t, tWrapExpr(t, expr), `
package thing

func RenderThing(msg string) vecty.HTMLOrComponent {
	vecty.Text(wrapped)
}`)
}

func TestParseExpressionWrapper2_CanParseStrings(t *testing.T) {
	expr, err := parseExpressionWrapper2("some string", false)
	require.NoError(t, err)
	requireEqStr(t, tWrapExpr(t, expr), `
package thing

func RenderThing(msg string) vecty.HTMLOrComponent {
	vecty.Text("some string")
}`)
}

// @todo test this
//func TestParseExpressionWrapper2_WhenQuotesAreIllegalReturnsError(t *testing.T) { }
