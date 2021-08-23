package html

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

func TestCanParseHtml(t *testing.T) {
	tag, err := ParseHtmlString(`
<nav class="navbar" role="navigation" aria-label="dropdown navigation">
  <div class="navbar-item has-dropdown">
    <a class="navbar-link">
      Docs
    </a>

    <div class="navbar-dropdown">
      <a class="navbar-item">
		<img src="some/source.png"/>
        Overview
      </a>
      <a class="navbar-item">
        Elements
      </a>
      <a class="navbar-item">
        Components
      </a>
      <hr class="navbar-divider" />
      <div class="navbar-item">
        Version 0.9.1
      </div>
    </div>
  </div>
</nav>
`)
	require.Nil(t, err)
	require.Equal(t, tag.DebugString(), `
nav class="navbar" role="navigation" aria-label="dropdown navigation"
  div class="navbar-item has-dropdown"
    a class="navbar-link"
      Docs
    div class="navbar-dropdown"
      a class="navbar-item"
        img src="some/source.png"
        Overview
      a class="navbar-item"
        Elements
      a class="navbar-item"
        Components
      hr class="navbar-divider"
      div class="navbar-item"
        Version 0.9.1
`)
}

func TestCanParseEmbeddedCode(t *testing.T) {
	tag, err := ParseHtmlString(`
<div class="some-thing">
  <div class="another thing">
	{&Header{}}
  </div>
  {s.someComp()}
  <img src="some/path.png" />
</div>
`)
	require.Nil(t, err)
	require.Equal(t, tag.DebugString(), `
div class="some-thing"
  div class="another thing"
    embed:{&Header{}}
  embed:{s.someComp()}
  img src="some/path.png"
`)
}

func TestReadUntilDepthIsZero(t *testing.T) {
	src := `<div class="some-thing">
  <div class="another thing">
	{&Header{}}
  </div>
  {s.someComp()}
  <img src="some/path.png" />
</div>

some stuff after
`
	r := bytes.NewReader([]byte(src))
	tag, htmlSrc, err := ParseHtml(r)
	require.Nil(t, err)
	require.NotNil(t, tag)
	require.Equal(t, `<div class="some-thing">
  <div class="another thing">
	{&Header{}}
  </div>
  {s.someComp()}
  <img src="some/path.png" />
</div>`, string(htmlSrc))
	remaining, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, "\n\nsome stuff after\n", string(remaining))
}

func TestDoesNothingWhenTheresNoHtml(t *testing.T) {
	src := `< 5
}

func Another() {
}
`
	r := bytes.NewReader([]byte(src))
	tag, htmlSrc, err := ParseHtml(r)
	require.NoError(t, err)
	require.Nil(t, tag)
	require.Nil(t, htmlSrc)
	remaining, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, "< 5\n}\n\nfunc Another() {\n}\n", string(remaining))
}

func TestCanDebugWhenHitEndOfStatement(t *testing.T) {
	_, err := ParseHtmlString(`
<div class="some-thing">
  <div class="another thing">
	{&Header{}}
</div>;

stm := "more"
`)
	require.EqualError(t, err, `hit end-of-statement should have terminated before this
expected closing tag: div
remaining text:
;

stm := "more"`)
}

func TestRegression1(t *testing.T) {
	src := `<div class="border-solid border-2 border-light-grey-500 p-3 mb-4">
			<div>
        		<input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="username" type="text" placeholder="Field Name" value="{f.FieldName}"/>
			</div>
			<div>
        		<input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="username" type="text" placeholder="Display Name" value="{f.DisplayName}"/>
			</div>
		</div>;

		fields = append(fields, fieldC)
	}

	return <div>
		{fields}
		<button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded" click="c.onAddField">Add field</button>
	</div>
}


type CustomInputPage struct {
	vecty.Core

	router *router.Router
	api *ajax.AjaxApi

	input *requests.InputRq
	output *requests.InputDisplayRs
	outputError string

	filterCard *comps.InputCard
	outputCard *comps.InputCard
	outputComp *Output
	filterEditor *comps.IBlizeEditor
}

func NewCustomInputPage() *CustomInputPage {
	api := ajax.NewAjaxApi("/input/custom")
	input := &requests.InputRq{
		Condition: &requests.Condition{},
	}
	c := &CustomInputPage{
		api: api,

		input: input,
		output: nil,

		filterCard: &comps.InputCard{IsOpen: true, Title: "Filter", Child: nil},
		outputCard: &comps.InputCard{IsOpen: true, Title: "Output", Child: nil},
		outputComp: &Output{Output: &input.Output},
		filterEditor: comps.NewIBlizeEditor(),
	}
	c.filterEditor.OnChange = c.onFilterChange
	c.filterEditor.Value = "some value"
	return c
}

func (c *CustomInputPage) recompute() {
	if c.input.InputName == "" {
		return
	}
	c.api.Request("/refresh", c.input, func(result string) {
		rs := &requests.InputRs{}
		if err := json.Unmarshal([]byte(result), rs); err != nil {
			panic(err)
		}

		if rs.Error != "" {
			c.outputError = rs.Error
			c.output = nil
			vecty.Rerender(p)
			return
		}

		c.outputError = ""
		c.output = rs.Data
		vecty.Rerender(p)
	})
}


func (c *CustomInputPage) onFilterChange(s string) {
	c.input.Condition.JsCondition = s
	c.recompute()
}

// Render implements the vecty.Component interface.
func (c *CustomInputPage) Render() vecty.ComponentOrHTML {
	filter := <div>
		<label class="block text-gray-700 text-sm font-bold mb-2" for="username">
		Value:
		</label>
		{c.filterEditor}
	</div>
	c.filterCard.Child = filter

	c.outputCard.Child = c.outputComp

	return <div class="flex flex-wrap bg-gray-100 w-full h-screen"> 
		<div class="w-4/12">
			<h1>Input name</h1>

			<div class="mt-4 mb-4">
		        <label class="block text-gray-700 text-sm font-bold mb-2" for="username">
		          Input
		        </label>
	        	<input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="username" type="text" placeholder="Input"/>
	      	</div>

			{c.filterCard}
			{c.outputCard}
		</div>
		<div class="flex-auto w-8/12 pl-4">
			{c.outputTable()}
		</div>
	</div>
}

func (c *CustomInputPage) outputTable() vecty.ComponentOrHTML {
	if c.output == nil {
		return nil
	}
	headers := vecty.List{}
	for _, h := range c.output.Headers {
		hC := <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
			{s:h.DisplayName}
		</th>
		headers = append(headers, hC)
	}

	rows := vecty.List{}
	for _, row := range c.output.Rows {
		cols := vecty.List{}
		for _, h := range c.output.Headers {
			s, _ := row[h.FieldName]
			col := <td class="px-6 py-4 whitespace-nowrap">{s:s}</td>
			cols = append(cols, col)
		}
		rowC := <tr>{cols}</tr>
		rows = append(rows, rowC)
	}

	return <div class="flex flex-col">
	    <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
	      <div class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
	        <div class="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
	          <table class="min-w-full divide-y divide-gray-200">
		            <thead class="bg-gray-50">
		            	<tr>{headers}</tr>
		            </thead>
		            <tbody class="bg-white divide-y divide-gray-200">{rows}</tbody>
	          </table>
	      </div>
	    </div>
	  </div>
	</div>
}
`
	r := bytes.NewReader([]byte(src))
	tag, htmlSrc, err := ParseHtml(r)
	require.Nil(t, err)
	require.NotNil(t, tag)
	require.Equal(t, `<div class="border-solid border-2 border-light-grey-500 p-3 mb-4">
			<div>
        		<input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="username" type="text" placeholder="Field Name" value="{f.FieldName}"/>
			</div>
			<div>
        		<input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="username" type="text" placeholder="Display Name" value="{f.DisplayName}"/>
			</div>
		</div>`, string(htmlSrc))
	//remaining, err := io.ReadAll(r)
	//require.NoError(t, err)
	//require.Equal(t, `
	//
	//	fields = append(fields, fieldC)
	//}
	//
	//return <div>
	//	{fields}
	//	<button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded" click="c.onAddField">Add field</button>
	//</div>
	//}
	//`, string(remaining))
}
