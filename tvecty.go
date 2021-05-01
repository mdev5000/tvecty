package tvecty

import (
	"bytes"
	"github.com/dave/dst/decorator"
	"github.com/mdev5000/tvecty/html"
	"io"
)

func ExtractHtml(filename string, w io.Writer, src []byte) ([]*html.TagOrText, error) {
	return sourceHtmlReplace(newHtmlTracker(), w, bytes.NewReader(src))
}

func ConvertToVecty(filename string, w io.Writer, src []byte) error {
	srcWithoutHtml := bytes.NewBuffer(nil)
	tracker, err := sourceHtmlReplace(newHtmlTracker(), srcWithoutHtml, bytes.NewReader(src))
	if err != nil {
		return err
	}
	parsed, err := tracker.parseAll()
	if err != nil {
		return err
	}
	f, err := decorator.Parse(srcWithoutHtml)
	if err := Replace(parsed, f); err != nil {
		return err
	}
	return decorator.Fprint(w, f)
}
