package html

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"strings"
)

func ParseHtmlString(htmlRaw string) (*TagOrText, error) {
	return ParseHtml([]byte(htmlRaw))
}

func ParseHtml(htmlRaw []byte) (*TagOrText, error) {
	stack := &tagStack{}
	b := bytes.NewBuffer(htmlRaw)
	z := html.NewTokenizer(b)
	var lastPop *TagOrText
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			err := z.Err()
			if err.Error() == "EOF" {
				return lastPop, nil
			} else {
				return lastPop, err
			}
		case html.TextToken:
			txt := string(z.Text())
			txtT := strings.TrimSpace(txt)
			if txtT == "" {
				continue
			}
			tag := &TagOrText{Text: txtT}
			if err := stack.pushChild(tag); err != nil {
				return lastPop, err
			}
		case html.SelfClosingTagToken:
			tnb, hasAttr := z.TagName()
			tag := &TagOrText{TagName: string(tnb)}
			if hasAttr {
				tag.Attr = parseAttributes(z)
			}
			if err := stack.pushChild(tag); err != nil {
				return lastPop, err
			}
		case html.StartTagToken:
			tnb, hasAttr := z.TagName()
			tag := &TagOrText{TagName: string(tnb)}
			if hasAttr {
				tag.Attr = parseAttributes(z)
			}
			stack.push(tag)
		case html.EndTagToken:
			var err error
			lastPop, err = stack.pop()
			if err != nil {
				return lastPop, err
			}
			tnb, _ := z.TagName()
			tn := string(tnb)
			if tn != lastPop.TagName {
				return lastPop, fmt.Errorf("expected closing tag '%s' but was '%s'", lastPop.TagName, tn)
			}
		}
	}
}

func ParseHtml2(r io.Reader) (*TagOrText, []byte, error) {
	stack := &tagStack{}
	z := html.NewTokenizer(r)
	var lastPop *TagOrText
	currentDepth := 0
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			err := z.Err()
			if err.Error() == "EOF" {
				if currentDepth == 0 {
					return lastPop, nil, nil
				}
				return lastPop, nil, fmt.Errorf("unexpected EOF, expected closing tag '%s'", lastPop.TagName)
			} else {
				return lastPop, nil, err
			}
		case html.TextToken:
			txt := string(z.Text())
			txtT := strings.TrimSpace(txt)
			if txtT == "" {
				continue
			}
			tag := &TagOrText{Text: txtT}
			if err := stack.pushChild(tag); err != nil {
				return lastPop, nil, err
			}
		case html.SelfClosingTagToken:
			tnb, hasAttr := z.TagName()
			tag := &TagOrText{TagName: string(tnb)}
			if hasAttr {
				tag.Attr = parseAttributes(z)
			}
			if err := stack.pushChild(tag); err != nil {
				return lastPop, nil, err
			}
		case html.StartTagToken:
			currentDepth += 1
			tnb, hasAttr := z.TagName()
			tag := &TagOrText{TagName: string(tnb)}
			if hasAttr {
				tag.Attr = parseAttributes(z)
			}
			stack.push(tag)
		case html.EndTagToken:
			currentDepth -= 1
			var err error
			lastPop, err = stack.pop()
			if err != nil {
				return lastPop, nil, err
			}
			tnb, _ := z.TagName()
			tn := string(tnb)
			if tn != lastPop.TagName {
				return lastPop, nil, fmt.Errorf("expected closing tag '%s' but was '%s'", lastPop.TagName, tn)
			}
			if currentDepth == 0 {
				return lastPop, z.Buffered(), nil
			}
		}
	}
}

func parseAttributes(z *html.Tokenizer) []*Attr {
	var out []*Attr
	for {
		key, val, more := z.TagAttr()
		out = append(out, &Attr{string(key), string(val)})
		if !more {
			break
		}
	}
	return out
}
