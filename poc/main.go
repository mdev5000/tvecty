package main

import (
	"github.com/dave/dst/decorator"
	"github.com/mdev5000/tvecty"
)

type filler struct {
	line  int
	ftype string
	value string
}

func ignore(i interface{}) {
}

func main() {
	fillers := map[int]filler{
		1: {line: 20, ftype: "html-return", value: `<div>
	{s:msg}
	{something()}
</div>`},
		2: {line: 20, ftype: "html-inline", value: `<div>
	{s:fmt.Sprintf("element %d", i)}
</div>
`},
	}
	ignore(fillers)

	code := `package thing

/*!!htmlfunc:*/ func RenderThing(msg string) {
	tvecty.Filler("1")
	//!!filler
	//!!filler
}

func Another() vecty.HTMLOrComponent {
	RenderThing()
}

func List() []vecty.HTMLOrComponent {
	var arr []vecty.HTMLOrComponent
	for i := 0; i < 10; i++ {
		/*!!start:htmltemplate*/ arr[0] = tvecty.Filler("2")
		//!!filler
		//!!filler
	}
	return i
}
`
	f, err := decorator.Parse(code)
	if err != nil {
		panic(err)
	}

	if err := tvecty.FinishHtmlFuncDefinitions(f); err != nil {
		panic(err)
	}
	//fmt.Println(len(f.Decls))

	if err := decorator.Print(f); err != nil {
		panic(err)
	}
}
