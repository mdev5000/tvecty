package html

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"strings"
)

func ParseHtmlString(htmlRaw string) (*TagOrText, error) {
	r := bytes.NewReader([]byte(htmlRaw))
	tag, _, err := ParseHtml(r)
	return tag, err
}

// ParseHtml Reads and parses the html from the starting tag to the matching end tag. If the tag depth is inconsistent
// or the tags at the same depth are not the same type then an error is returned.
//
// After the html is parsed r is reset to the remaining bytes that have not been parsed.
//
func ParseHtml(r *bytes.Reader) (tag *TagOrText, htmlSrc []byte, err error) {
	stack := &tagStack{}
	remainingAtStart := int64(r.Len())
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
			// @todo Possibly consider a better way to do this, especially for something like '2 < 3' early in a program, where essentially the entire program has to be converted into a string :(
			txtb := z.Text()
			txt := string(txtb)
			txtT := strings.TrimSpace(txt)
			if txtT == "" {
				continue
			}

			if strings.HasPrefix(txtT, ";") {
				return nil, nil, fmt.Errorf(
					"hit end-of-statement should have terminated before this\nexpected closing tag: %s\nremaining text:\n%s",
					lastPop.TagName,
					txtT,
				)
			}

			// No html has been parsed, so everything in text should be returned to the reader.
			if lastPop == nil && stack.isEmpty() {
				r.Reset(txtb)
				return nil, nil, nil
			}
			tag := &TagOrText{Text: txtT}
			if err := stack.pushChild(tag); err != nil {
				return lastPop, nil, err
			}
		case html.SelfClosingTagToken:
			tnb, hasAttr := z.TagName()
			//fmt.Println(strings.Repeat("-", currentDepth), string(tnb), "(self-closing)", string(z.Raw()))
			tag := &TagOrText{TagName: string(tnb)}
			if hasAttr {
				tag.Attr = parseAttributes(z)
			}
			if currentDepth == 0 {
				return finalizeTagParsing(r, z, remainingAtStart, tag)
			}
			if err := stack.pushChild(tag); err != nil {
				return lastPop, nil, err
			}
		case html.StartTagToken:
			currentDepth += 1
			tnb, hasAttr := z.TagName()
			//fmt.Println(strings.Repeat("-", currentDepth - 1), string(tnb), "starting", string(z.Raw()))
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
			//fmt.Println(strings.Repeat("-", currentDepth), "/" + string(tnb), "closing")
			if currentDepth == 0 {
				return finalizeTagParsing(r, z, remainingAtStart, lastPop)
			}
		}
	}
}

func finalizeTagParsing(r *bytes.Reader, z *html.Tokenizer, remainingAtStart int64, lastPop *TagOrText) (*TagOrText, []byte, error) {
	// Extract the html content just parsed into html and return it to the caller.

	// [       orig buffer size                        ]
	//               [ remaining at start              ]
	//                            [ remaining buffer   ]
	// [ start index ]
	//               [ start html ]
	remainingBuffer := z.Buffered()
	remainingPastBuffer, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, err
	}
	remainingAll := make([]byte, len(remainingBuffer)+len(remainingPastBuffer))
	copy(remainingAll, remainingBuffer)
	copy(remainingAll[len(remainingBuffer):], remainingPastBuffer)

	origBufferSize := r.Size()
	startIndex := origBufferSize - remainingAtStart
	startHtmlLen := origBufferSize - int64(len(remainingAll)) - startIndex
	startHtml := make([]byte, startHtmlLen)
	_, err = r.ReadAt(startHtml, startIndex)
	// Then set the r to whatever bytes were remaining after parsing the html.
	r.Reset(remainingAll)
	return lastPop, startHtml, err
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
