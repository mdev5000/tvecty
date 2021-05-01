package strtokenizer

const (
	EmptyRune = rune(0)
)

// StringCharacterIndex An iterator for moving character by character through a string.
type StringCharacterIndex struct {
	R     []rune
	Index int
}

func NewStringCharacterIndex(s string) *StringCharacterIndex {
	return &StringCharacterIndex{
		R:     []rune(s),
		Index: -1,
	}
}

func (s *StringCharacterIndex) MoveNext() bool {
	if !s.HasNext() {
		return false
	}
	s.Index += 1
	return true
}

func (s *StringCharacterIndex) MovePrev() bool {
	if !s.HasPrev() {
		return false
	}
	s.Index -= 1
	return true
}

func (s *StringCharacterIndex) Value() rune {
	if s.Index < len(s.R) {
		return s.R[s.Index]
	} else {
		return EmptyRune
	}
}

func (s *StringCharacterIndex) NextValue() rune {
	s.MoveNext()
	return s.Value()
}

func (s *StringCharacterIndex) HasNext() bool {
	return s.Index < (len(s.R) - 1)
}

func (s *StringCharacterIndex) HasPrev() bool {
	return 0 < s.Index
}

func (s *StringCharacterIndex) RemainingString() string {
	if s.Index == -1 {
		return string(s.R)
	}
	return string(s.R[s.Index:])
}
