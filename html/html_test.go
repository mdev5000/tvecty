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

func TestRegression1(t *testing.T) {
	src := `<div class="border-solid border-2 border-light-grey-500 p-3 mb-4">
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
	remaining, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, `

		fields = append(fields, fieldC)
	}

	return <div>
		{fields}
		<button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded" click="c.onAddField">Add field</button>
	</div>
}
`, string(remaining))
}
