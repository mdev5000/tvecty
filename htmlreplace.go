package tvecty

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/mdev5000/tvecty/html"
	"io"
)

type htmlReplaceReader = bytes.Reader

type htmlReplaceStateFn = func(w io.Writer) (bool, error)

type htmlReplaceState struct {
	ht          htmlTracker
	prevStateFn htmlReplaceStateFn
	stateFn     htmlReplaceStateFn
	r           *htmlReplaceReader
}

func sourceHtmlReplace(ht htmlTracker, w io.Writer, src *bytes.Reader) (htmlTracker, error) {
	rs := htmlReplaceState{ht: ht, r: src}
	rs.stateFn = rs.htmlStateReadChars
	for {
		if done, err := rs.stateFn(w); err != nil {
			return nil, err
		} else if done {
			return rs.ht, nil
		}
	}
}

func (rs *htmlReplaceState) htmlStateReadChars(w io.Writer) (bool, error) {
	for {
		c, _, err := rs.r.ReadRune()
		if err == io.EOF {
			return true, nil
		}
		if err != nil {
			return true, err
		}
		switch c {
		case '<': // Parse html: ex: <div>
			// Unread so the html can be correctly parsing, ex. <div> instead of div>.
			err = rs.r.UnreadRune()
			rs.stateFn = rs.htmlStateReadHtml
			return false, err
		case '/': // Parse comments, ex: // testing
			err := writeRune(w, c)
			rs.prevStateFn = rs.htmlStateReadChars
			rs.stateFn = rs.htmlStateReadComment
			return false, err
		case '"': // Parse strings
			err := writeRune(w, c)
			rs.prevStateFn = rs.htmlStateReadChars
			rs.stateFn = rs.htmlStateReadStringQuote
			return false, err
		case '`': // Parse tick strings
			err := writeRune(w, c)
			rs.prevStateFn = rs.htmlStateReadChars
			rs.stateFn = rs.htmlStateReadStringTick
			return false, err
		default:
			if err := writeRune(w, c); err != nil {
				return false, err
			}
		}
	}
}

func (rs *htmlReplaceState) htmlStateReadStringQuote(w io.Writer) (bool, error) {
	return rs.htmlStateReadString(w, '"')
}

func (rs *htmlReplaceState) htmlStateReadStringTick(w io.Writer) (bool, error) {
	return rs.htmlStateReadString(w, '`')
}

func (rs *htmlReplaceState) htmlStateReadString(w io.Writer, quoteChar rune) (bool, error) {
	for {
		c, n, err := rs.r.ReadRune()
		if err != nil {
			return true, err
		}
		if n == 0 {
			return true, nil
		}
		switch c {
		case '\\':
			// Next character will be an escaped character, so read this character
			// and then next one, that way the next character is not parsed.
			err := writeRune(w, c)
			if err != nil {
				return true, err
			}
			c, n, err := rs.r.ReadRune()
			if err != nil {
				return true, err
			}
			if n == 0 {
				return true, nil
			}
			err = writeRune(w, c)
			if err != nil {
				return true, err
			}
		case quoteChar:
			err := writeRune(w, c)
			rs.stateFn = rs.prevStateFn
			return false, err
		default:
			err := writeRune(w, c)
			if err != nil {
				return true, err
			}
		}
	}
}

func (rs *htmlReplaceState) htmlStateReadComment(w io.Writer) (bool, error) {
	c, _, err := rs.r.ReadRune()
	if err == io.EOF {
		return true, nil
	}
	if err != nil {
		return true, err
	}
	if err := writeRune(w, c); err != nil {
		return false, err
	}
	// Read a line comment, ex: // testing
	if c == '/' {
		return rs.readUntil(w, []rune{'\n'})
	}
	// Read a multiline comment, ex: /* testing */
	if c == '*' {
		return rs.readUntil(w, []rune{'*', '/'})
	}
	rs.stateFn = rs.prevStateFn
	return false, nil
}

func (rs *htmlReplaceState) readUntil(w io.Writer, escapeCharacters []rune) (bool, error) {
	if len(escapeCharacters) == 0 {
		return false, errors.New("escapeCharacters cannot be empty")
	}
	remaining := escapeCharacters
	for {
		c, n, err := rs.r.ReadRune()
		if err != nil {
			return true, err
		}
		if n == 0 {
			return true, nil
		}
		if err := writeRune(w, c); err != nil {
			return false, err
		}
		// Check if we've hit the next escape character sequence, if not reset.
		if remaining[0] == c {
			remaining = remaining[1:]
			if len(remaining) == 0 {
				rs.stateFn = rs.prevStateFn
				return false, err
			}
		} else {
			remaining = escapeCharacters
			if remaining[0] == c {
				remaining = remaining[1:]
			}
		}
	}
}

func (rs *htmlReplaceState) htmlStateReadHtml(w io.Writer) (bool, error) {
	tag, htmlSrc, err := html.ParseHtml(rs.r)
	if err != nil {
		return true, err
	}
	// If the source did not have html (ex. 1 < 2, or 10 << 2), then the tag will be nil and there's nothing to add.
	// However we still need to remove the leading < character, to avoid an infinite loop.
	if tag == nil {
		c, n, err := rs.r.ReadRune()
		if err != nil {
			return true, err
		}
		if n == 0 {
			return true, nil
		}
		if err := writeRune(w, c); err != nil {
			return false, err
		}
	} else {
		var tagId htmlTrackingId
		rs.ht, tagId = rs.ht.add(tag)
		_, err = fmt.Fprintf(w, "tvecty.Html(%d, `%s`)", tagId, htmlSrc)
	}
	rs.stateFn = rs.htmlStateReadChars
	return false, err
}

func writeRune(w io.Writer, c rune) error {
	_, err := fmt.Fprintf(w, "%c", c)
	return err
}
