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

func TestTokenizeParts_CanTokenizeEmptyValue(t *testing.T) {
	parts, err := tokenizeParts("")
	require.NoError(t, err)
	require.Equal(t, parts, []attributeToken{
		{"", false},
	})
}

func TestTokenizeParts_NewLineIsIllegalInExpressions(t *testing.T) {
	_, err := tokenizeParts(`{first
}
`)
	require.EqualError(t, err, "illegal character '\\n' in embedded code block in expressions '{first\n}\n'")
}

func TestTokenizeParts_ErrorWhenUnclosedExpressions(t *testing.T) {
	_, err := tokenizeParts(`{first`)
	require.EqualError(t, err, "missing closing '}' tag for embedded code in '{first'")
}

func TestTokenizeParts_ErrorWhenRandomClosingExpressions(t *testing.T) {
	_, err := tokenizeParts(`stuff}`)
	require.EqualError(t, err, "unexpected '}' in expressions 'stuff}'")
}

func TestParseExpressionWrapper2_CanParseExpressions(t *testing.T) {
	expr, err := parseExpressionWrapper("first", true, true)
	require.NoError(t, err)
	requireEqStr(t, tWrapExpr(t, expr), `
package thing

func RenderThing(msg string) vecty.HTMLOrComponent {
	first
}`)
}

func TestParseExpressionWrapper2_CanParseExpressionsWithStringModifiers(t *testing.T) {
	expr, err := parseExpressionWrapper("s:wrapped", true, true)
	require.NoError(t, err)
	requireEqStr(t, tWrapExpr(t, expr), `
package thing

func RenderThing(msg string) vecty.HTMLOrComponent {
	vecty.Text(wrapped)
}`)
}

func TestParseExpressionWrapper2_CanParseStrings(t *testing.T) {
	expr, err := parseExpressionWrapper("some string", false, true)
	require.NoError(t, err)
	requireEqStr(t, tWrapExpr(t, expr), `
package thing

func RenderThing(msg string) vecty.HTMLOrComponent {
	vecty.Text("some string")
}`)
}

func TestParseExpressionWrapper2_CanParseNonWrappedStrings(t *testing.T) {
	expr, err := parseExpressionWrapper("some string", false, false)
	require.NoError(t, err)
	requireEqStr(t, tWrapExpr(t, expr), `
package thing

func RenderThing(msg string) vecty.HTMLOrComponent {
	"some string"
}`)
}
