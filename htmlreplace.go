package tvecty

import (
	"fmt"
	"github.com/mdev5000/tvecty/internal/strtokenizer"
	"io"
)

type htmlReplaceStateFn = func(w io.Writer, t *strtokenizer.StringCharacterIndex) error

type htmlReplaceState struct {
	htmlTagDepth int
	prevStateFn htmlReplaceStateFn
	stateFn htmlReplaceStateFn
}

// @todo better handle unmatched html tags.
func SourceHtmlReplace(w io.Writer, src string) error {
	t := strtokenizer.NewStringCharacterIndex(src)
	rs := htmlReplaceState{}
	rs.stateFn = rs.htmlStateReadChars
	for t.HasNext() {
		if err := rs.stateFn(w, t); err != nil {
			return err
		}
	}
	return nil
}

func (rs *htmlReplaceState) htmlStateReadChars(w io.Writer, t *strtokenizer.StringCharacterIndex) error {
	for t.HasNext() {
		c := t.NextValue()
		switch c {
		case '"':
			err := writeRune(w, c)
			rs.prevStateFn = rs.htmlStateReadChars
			rs.stateFn = rs.htmlStateReadStringQuote
			return err
		case '`':
			err := writeRune(w, c)
			rs.prevStateFn = rs.htmlStateReadChars
			rs.stateFn = rs.htmlStateReadStringTick
			return err
		case '<':
			_, err := fmt.Fprint(w, "tvecty.Html(`")
			if err != nil {
				return err
			}
			err = writeRune(w, c)
			rs.htmlTagDepth = 1
			rs.stateFn = rs.htmlStateReadHtml
			return err
		default:
			err := writeRune(w, c)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (rs *htmlReplaceState) htmlStateReadStringQuote(w io.Writer, t *strtokenizer.StringCharacterIndex) error {
	return rs.htmlStateReadString(w, t, '"')
}

func (rs *htmlReplaceState) htmlStateReadStringTick(w io.Writer, t *strtokenizer.StringCharacterIndex) error {
	return rs.htmlStateReadString(w, t, '`')
}

func (rs *htmlReplaceState) htmlStateReadString(w io.Writer, t *strtokenizer.StringCharacterIndex, quoteChar rune) error {
	for t.HasNext() {
		c := t.NextValue()
		switch c {
		case '\\':
			// Next character will be an escaped character, so read this character
			// and then next one, that way the next character is not parsed.
			err := writeRune(w, c)
			if err != nil {
				return err
			}
			if t.HasNext() {
				c := t.NextValue()
				err := writeRune(w, c)
				if err != nil {
					return err
				}
			}
		case quoteChar:
			err := writeRune(w, c)
			rs.stateFn = rs.prevStateFn
			return err
		default:
			err := writeRune(w, c)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (rs *htmlReplaceState) htmlStateReadHtml(w io.Writer, t *strtokenizer.StringCharacterIndex) error {
	for t.HasNext() {
		c := t.NextValue()
		switch c {
		case '"':
			err := writeRune(w, c)
			rs.prevStateFn = rs.htmlStateReadHtml
			rs.stateFn = rs.htmlStateReadStringQuote
			return err
		case '<':
			err := writeRune(w, c)
			if err != nil {
				return err
			}
			if !t.HasNext() {
				continue
			}
			c = t.NextValue()
			// Check if is opening on closing (ex. <tag> or </tag>)
			if c == '/' {
				rs.htmlTagDepth -= 1
			} else {
				rs.htmlTagDepth += 1
			}
			err = writeRune(w, c)
			rs.stateFn = rs.htmlStateReadHtml
			return err
		case '>':
			err := writeRune(w, c)
			if err != nil {
				return err
			}
			if rs.htmlTagDepth < 0 {
				panic("invalid html tag depth")
			} else if rs.htmlTagDepth == 0 {
				_, err = fmt.Fprint(w, "`)")
				rs.stateFn = rs.htmlStateReadChars
				return err
			}
		default:
			err := writeRune(w, c)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func writeRune(w io.Writer, c rune) error {
	_, err := fmt.Fprintf(w, "%c", c)
	return err
}