package strtokenizer

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStringCharacterIndex_IndicatesWhenHasNext(t *testing.T) {
	s := NewStringCharacterIndex("my")
	require.True(t, s.HasNext())
	require.Equal(t, s.NextValue(), 'm')

	require.True(t, s.HasNext())
	require.Equal(t, s.NextValue(), 'y')

	require.False(t, s.HasNext())
}
