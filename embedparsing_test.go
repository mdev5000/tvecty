package tvecty

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTokenizeExpressionParts_CorrectlyTokenizesParts(t *testing.T) {
	parts, err := tokenizeExpressionParts(`{first} Second value {third}
And some more {fourth} stuff
And another thing of text
{fifth}
`)
	require.NoError(t, err)
	require.Equal(t, parts, []embedToken{
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

func TestTokenizeExpressionParts_CanTokenizeEmptyValue(t *testing.T) {
	parts, err := tokenizeExpressionParts("")
	require.NoError(t, err)
	require.Equal(t, parts, []embedToken{
		{"", false},
	})
}

func TestTokenizeExpressionParts_NewLineIsIllegalInExpressions(t *testing.T) {
	_, err := tokenizeExpressionParts(`{first
}
`)
	require.EqualError(t, err, "illegal character '\\n' in embedded code block in expressions '{first\n}\n'")
}

func TestTokenizeExpressionParts_ErrorWhenUnclosedExpressions(t *testing.T) {
	_, err := tokenizeExpressionParts(`{first`)
	require.EqualError(t, err, "missing closing '}' tag for embedded code in '{first'")
}

func TestTokenizeExpressionParts_ErrorWhenRandomClosingExpressions(t *testing.T) {
	_, err := tokenizeExpressionParts(`stuff}`)
	require.EqualError(t, err, "unexpected '}' in expressions 'stuff}'")
}

func TestParseExpressionOrText_CanParseExpressions(t *testing.T) {
	expr, err := parseExpressionOrText("first", true, true)
	require.NoError(t, err)
	requireEqStr(t, tWrapExpr(t, expr), `
package thing

func RenderThing(msg string) vecty.HTMLOrComponent {
	first
}`)
}

func TestParseExpressionOrText_CanParseExpressionsWithStringModifiers(t *testing.T) {
	expr, err := parseExpressionOrText("s:wrapped", true, true)
	require.NoError(t, err)
	requireEqStr(t, tWrapExpr(t, expr), `
package thing

func RenderThing(msg string) vecty.HTMLOrComponent {
	vecty.Text(wrapped)
}`)
}

func TestParseExpressionOrText_CanParseStrings(t *testing.T) {
	expr, err := parseExpressionOrText("some string", false, true)
	require.NoError(t, err)
	requireEqStr(t, tWrapExpr(t, expr), `
package thing

func RenderThing(msg string) vecty.HTMLOrComponent {
	vecty.Text("some string")
}`)
}

func TestParseExpressionOrText_CanParseNonWrappedStrings(t *testing.T) {
	expr, err := parseExpressionOrText("some string", false, false)
	require.NoError(t, err)
	requireEqStr(t, tWrapExpr(t, expr), `
package thing

func RenderThing(msg string) vecty.HTMLOrComponent {
	"some string"
}`)
}
