package tvecty

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestPlacesAllHtmlInAVectyWrapper(t *testing.T) {
	in := `package somepackage

html MyRender() {
	<div>
		<span>
			Content
		</span>
	</div>
}
`

	b := bytes.NewBuffer(nil)
	require.NoError(t, SourceHtmlReplace(b, in))
	requireEqStr(t, b.String(), strings.Replace(`
package somepackage

html MyRender() {
	tvecty.Html(:tick:<div>
		<span>
			Content
		</span>
	</div>:tick:)
}
`, ":tick:", "`", -1))
}
