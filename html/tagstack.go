package html

import "fmt"

type stackItem struct {
	Tag    *TagOrText
	Parent *stackItem
}

type tagStack struct {
	Next *stackItem
}

func (s *tagStack) pop() (*TagOrText, error) {
	if s.Next == nil {
		return nil, fmt.Errorf("popped last item on the stack")
	}
	next := s.Next
	tag := next.Tag
	s.Next = next.Parent
	if s.Next != nil {
		s.Next.Tag.AppendChild(tag)
	}
	return tag, nil
}

func (s *tagStack) isEmpty() bool {
	return s.Next == nil
}

func (s *tagStack) pushChild(tag *TagOrText) error {
	s.Next.Tag.AppendChild(tag)
	return nil
}

func (s *tagStack) push(tag *TagOrText) {
	next := s.Next
	s.Next = &stackItem{
		Tag:    tag,
		Parent: next,
	}
}
