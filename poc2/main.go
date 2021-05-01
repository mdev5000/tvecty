package main

import (
	"bytes"
	"fmt"
	"github.com/mdev5000/tvecty"
)

func main() {
	src := `package something

func Render(msg string) vecty.ComponentOrHTML {
	another := <div>
		<div>
			{s:"Part"}
		</div>
		<div>
			{s:"Subpart"}
		</div>
	</div>
	return <div>
		<span>
			{s:msg}
			{another}
		</span>
	</div>
}
`
	b := bytes.NewBuffer(nil)
	if err := tvecty.ConvertToVecty("", b, []byte(src)); err != nil {
		panic(err)
	}
	fmt.Println(b.String())
}
