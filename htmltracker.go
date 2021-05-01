package tvecty

import (
	"github.com/dave/dst"
	"github.com/mdev5000/tvecty/html"
)

type htmlTrackingId = int

type htmlTracker []*html.TagOrText

type htmlTrackerParsed []dst.Expr

func newHtmlTracker() htmlTracker {
	return nil
}

func (h htmlTracker) add(tag *html.TagOrText) (htmlTracker, htmlTrackingId) {
	out := append(h, tag)
	return out, len(out)
}

func (h htmlTracker) parseAll() (htmlTrackerParsed, error) {
	out := make(htmlTrackerParsed, len(h)+1)
	for i, tag := range h {
		exprs, err := tagToAst(nil, tag)
		if err != nil {
			return nil, err
		}
		if len(exprs) == 0 {
			panic("exprs should never be empty")
		}
		out[i+1] = exprs[0]
	}
	return out, nil
}

func (h htmlTrackerParsed) Get(id int) (dst.Expr, bool) {
	if id < 0 || len(h) <= id {
		return nil, false
	}
	return h[id], true
}
