package tvecty

import (
	"bytes"
	"github.com/dave/dst/decorator"
	"io"
)

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
