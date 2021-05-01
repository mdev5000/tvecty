package tvecty

import "github.com/mdev5000/tvecty/html"

type htmlTrackingId = int

type htmlTracker []*html.TagOrText

func newHtmlTracker() htmlTracker {
	return nil
}

func (h htmlTracker) add(tag *html.TagOrText) (htmlTracker, htmlTrackingId) {
	out := append(h, tag)
	return out, len(out)
}
