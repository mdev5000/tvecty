# tvecty

An experiment for embedding html templates in Go files that
can rendered into **Vecty** code
(https://github.com/hexops/vecty).

`example.vtpl`
```
package comps

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

func ExampleComp() vecty.ComponentOrHTML {
    embedAsText := "embed this as text"
    child := <div>this is a child</div>
	
    return <div class="root">
    	<p>this is some test</p>
    	<p>{s:embedAsText}</p>
        {child}
        {AnotherComp()}
    </div>
}

func AnotherComp() vecty.ComponentOrHTML {
    //...
}
```


# Installation

```bash
go install github.com/mdev5000/tvecty/cmd/tvecty@latest
```


# Usage

Compiling the example above via:

```bash
tvecty compile file example.vtpl
```

Will output the following:

```go
package comps

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

func ExampleComp() vecty.ComponentOrHTML {
	embedAsText := "embed this as text"
	child :=
		elem.Div(
			vecty.Text("this is a child"),
		)
	return elem.Div(
		vecty.Markup(
			vecty.Class("root"),
		),
		elem.Paragraph(
			vecty.Text("this is some test"),
		),
		elem.Paragraph(
			vecty.Text(embedAsText),
		),
		child,
		AnotherComp())
}

func AnotherComp() vecty.ComponentOrHTML {
	//...
}
```

You can also output directly to file or compile based on
directory matching.

```bash
# compile a single file
tvecty compile file example.vtpl example.vtpl.go

# compile a directory file
tvecty compile dir ./tpl/**/*.vtpl
```

Or setup to use with Go generate. 

`tpl/gen.go`
```go
//go:generate tvecty compile directory **/*.vtpl
```

```bash
go gen ./tpl
```