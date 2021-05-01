package tvecty

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestPlacesHtmlInAVectyWrapper(t *testing.T) {
	in := `package somepackage

func MyRender() vecty.HTMLOrComponent {
	return <div>
		<span>
			Content
		</span>
	</div>
}
`

	tracker := newHtmlTracker()
	src := bytes.NewReader([]byte(in))
	srcOut := bytes.NewBuffer(nil)

	tracker, err := sourceHtmlReplace(tracker, srcOut, src)
	require.NoError(t, err)
	requireEqStr(t, srcOut.String(), strings.Replace(`
package somepackage

func MyRender() vecty.HTMLOrComponent {
	return tvecty.Html(1, :tick:<div>
		<span>
			Content
		</span>
	</div>:tick:)
}
`, ":tick:", "`", -1))
}

func TestCorrectlyExtractsHtmlTagInformation(t *testing.T) {
	in := `package somepackage

func MyRender() vecty.HTMLOrComponent {
	return <div>
		<span>
			Content
		</span>
	</div>
}
`
	tracker := newHtmlTracker()
	src := bytes.NewReader([]byte(in))
	srcOut := bytes.NewBuffer(nil)
	tracker, err := sourceHtmlReplace(tracker, srcOut, src)
	require.NoError(t, err)
	require.Len(t, tracker, 1)
	require.Equal(t, `
div 
  span 
    Content
`, tracker[0].DebugString())
}
