package main

import (
	"github.com/dave/dst/decorator"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	src := `package something

func Render(msg string) vecty.ComponentOrHTML {
	return vecty.Class("some-class", myvar)
}
`
	f, err := decorator.Parse(src)
	if err != nil {
		panic(err)
	}
	spew.Dump(f.Decls[0])
	//b := bytes.NewBuffer(nil)
	//if err := tvecty.ConvertToVecty("", b, []byte(src)); err != nil {
	//	panic(err)
	//}
	//fmt.Println(b.String())
}
