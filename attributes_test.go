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
