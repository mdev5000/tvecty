package tvecty

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestSourceHtmlReplace_PlacesHtmlInAVectyWrapper(t *testing.T) {
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

func TestSourceHtmlReplace_CorrectlyExtractsHtmlTagInformation(t *testing.T) {
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

func TestSourceHtmlReplace_IgnoresHtmlInLineComments(t *testing.T) {
	in := `package somepackage

// <div>test
func MyRender() int {
	// <div>another
	return 10 / 2
}

//</div>
`
	src := bytes.NewReader([]byte(in))
	srcOut := bytes.NewBuffer(nil)
	_, err := sourceHtmlReplace(newHtmlTracker(), srcOut, src)
	require.NoError(t, err)
	requireEqStr(t, srcOut.String(), `
package somepackage

// <div>test
func MyRender() int {
	// <div>another
	return 10 / 2
}

//</div>
`)
}

func TestSourceHtmlReplace_IgnoreHtmlInMultilineComments(t *testing.T) {
	in := `package somepackage

/**
 * <div>here
 **/
func MyRender(t *Thing) int {
	/* <span>single */
	return t.First / 2
}
/*
 <div>another
*/
`
	src := bytes.NewReader([]byte(in))
	srcOut := bytes.NewBuffer(nil)
	_, err := sourceHtmlReplace(newHtmlTracker(), srcOut, src)
	require.NoError(t, err)
	requireEqStr(t, srcOut.String(), `
package somepackage

/**
 * <div>here
 **/
func MyRender(t *Thing) int {
	/* <span>single */
	return t.First / 2
}
/*
 <div>another
*/
`)
}

func TestSourceHtmlReplace_IgnoreLessThanSymbol(t *testing.T) {
	in := `package somepackage

func MyRender() int {
	if 2 < 3 {
		return 10 << 2
	}
	return 0
}
`
	src := bytes.NewReader([]byte(in))
	srcOut := bytes.NewBuffer(nil)
	_, err := sourceHtmlReplace(newHtmlTracker(), srcOut, src)
	require.NoError(t, err)
	requireEqStr(t, srcOut.String(), `
package somepackage

func MyRender() int {
	if 2 < 3 {
		return 10 << 2
	}
	return 0
}`)
}
