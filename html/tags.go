package html

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Attr struct {
	Name  string
	Value string
}

type TagOrText struct {
	TagName  string
	Text     string
	Attr     []*Attr
	Children []*TagOrText
}

func (t *TagOrText) AppendChild(child *TagOrText) {
	t.Children = append(t.Children, child)
}

func (t *TagOrText) DebugString() string {
	b := bytes.NewBufferString("\n")
	t.debugString(b, "")
	return b.String()
}

func (t *TagOrText) debugString(w io.Writer, depth string) {
	if t.TagName == "" {
		fmt.Fprint(w, depth)
		if t.IsGoCodeEmbed() {
			fmt.Fprintln(w, "embed:"+t.Text)
		} else {
			fmt.Fprintln(w, t.Text)
		}
		return
	}
	fmt.Fprint(w, depth)
	fmt.Fprintln(w, t.TagName, t.debugAttr())
	for _, child := range t.Children {
		child.debugString(w, depth+"  ")
	}
}

func (t *TagOrText) debugAttr() string {
	var out []string
	for _, attr := range t.Attr {
		out = append(out, fmt.Sprintf(`%s="%s"`, attr.Name, attr.Value))
	}
	return strings.Join(out, " ")
}

func (t *TagOrText) IsGoCodeEmbed() bool {
	s := strings.TrimSpace(t.Text)
	return strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}")
}

//func (t *TagOrText) WriteText(w io.Writer) {
//	if t.IsComponent() {
//		fmt.Fprint(w, t.Text[1:])
//		fmt.Fprintln(w, ",")
//	} else {
//		fmt.Fprint(w, `vecty.Text("`)
//		fmt.Fprint(w, t.Text)
//		fmt.Fprintln(w, `"),`)
//	}
//}

//func (t *TagOrText) IsComponent() bool {
//	return strings.HasPrefix(t.Text, "!")
//}
