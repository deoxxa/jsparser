package main

type constraint struct {
	inverse bool
	value   string
}

func (c constraint) String() string {
	s := "+"
	if c.inverse {
		s = "~"
	}

	s += c.value

	return s
}
