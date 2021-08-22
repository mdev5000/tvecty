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
	src := bytes.NewReader([]byte(in))
	srcOut := bytes.NewBuffer(nil)

	_, err := sourceHtmlReplace(newHtmlTracker(), srcOut, src)
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

func TestSourceHtmlReplace_regression1(t *testing.T) {
	in := `package pages

import (
    "github.com/hexops/vecty"
    "github.com/hexops/vecty/elem"
    "github.com/hexops/vecty/event"
    "github.com/mdev5000/csvtransform/frontend/ajax"
    "github.com/mdev5000/csvtransform/frontend/comps"
    "github.com/mdev5000/csvtransform/frontend/router"
    "github.com/mdev5000/csvtransform/requests"
)

type Output struct {
	vecty.Core
	Output *requests.Output
}

func (c *Output) onAddField(*vecty.Event) {
	c.Output.Fields = append(c.Output.Fields, requests.Field{})
	vecty.Rerender(c)
}

func (c *Output) Render() vecty.ComponentOrHTML {
	fields := vecty.List{}
	for _, f := range c.Output.Fields {
		fieldC := <div class="border-solid border-2 border-light-grey-500 p-3 mb-4">
			<div>
        		<input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="username" type="text" placeholder="Field Name" value="{f.FieldName}"/>
			</div>
			<div>
        		<input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="username" type="text" placeholder="Display Name" value="{f.DisplayName}"/>
			</div>
		</div>

		fields = append(fields, fieldC)
	}

	return <div>
		{fields}
		<button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded" click="c.onAddField">Add field</button>
	</div>
}
`
	src := bytes.NewReader([]byte(in))
	srcOut := bytes.NewBuffer(nil)
	_, err := sourceHtmlReplace(newHtmlTracker(), srcOut, src)
	require.NoError(t, err)
	requireEqStr(t, srcOut.String(), strings.Replace(`package pages

import (
    "github.com/hexops/vecty"
    "github.com/hexops/vecty/elem"
    "github.com/hexops/vecty/event"
    "github.com/mdev5000/csvtransform/frontend/ajax"
    "github.com/mdev5000/csvtransform/frontend/comps"
    "github.com/mdev5000/csvtransform/frontend/router"
    "github.com/mdev5000/csvtransform/requests"
)

type Output struct {
	vecty.Core
	Output *requests.Output
}

func (c *Output) onAddField(*vecty.Event) {
	c.Output.Fields = append(c.Output.Fields, requests.Field{})
	vecty.Rerender(c)
}

func (c *Output) Render() vecty.ComponentOrHTML {
	fields := vecty.List{}
	for _, f := range c.Output.Fields {
		fieldC := tvecty.Html(1, :tick:<div class="border-solid border-2 border-light-grey-500 p-3 mb-4">
			<div>
        		<input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="username" type="text" placeholder="Field Name" value="{f.FieldName}"/>
			</div>
			<div>
        		<input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="username" type="text" placeholder="Display Name" value="{f.DisplayName}"/>
			</div>
		</div>:tick:)

		fields = append(fields, fieldC)
	}

	return tvecty.Html(2, :tick:<div>
		{fields}
		<button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded" click="c.onAddField">Add field</button>
	</div>:tick:)
}
`, ":tick:", "`", -1))
}
