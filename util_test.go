package tvecty

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func requireEqStr(t *testing.T, actual, expected string) {
	requireEq(t, actual, strings.TrimSpace(expected)+"\n")
}

func requireEq(t *testing.T, actual, expected interface{}) {
	require.Equal(t, expected, actual)
}
