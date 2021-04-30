package html

import (
	"github.com/stretchr/testify/require"
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
