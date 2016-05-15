package main

import (
	"fmt"
	"strings"
)

type rule struct {
	constraints []constraint
	tokens      []token
}

func (r rule) String() string {
	s := ""

	for _, c := range r.constraints {
		s += c.String() + " "
	}

	for _, t := range r.tokens {
		s += t.String() + " "
	}

	return strings.TrimSpace(s)
}

func (r rule) EBNF(options map[string]bool) string {
	s := ""

	for _, c := range r.constraints {
		if options[c.value] != c.inverse {
			return ""
		}
	}

	for _, t := range r.tokens {
		s += t.EBNF(options) + " "
	}

	return strings.TrimSpace(s)
}

func (r rule) Go(m map[string]*production, p production) string {
	s := "r := input\n\n"

	b := "return r, true, nil"
	for i := len(r.tokens) - 1; i >= 0; i-- {
		b = r.tokens[i].Go(m, p, b)
	}

	s += b

	for _, c := range r.constraints {
		if c.inverse {
			s = fmt.Sprintf("    if !flag%s {\n%s\n    }", c.value, strings.Replace(s, "\n", "\n  ", -1))
		} else {
			s = fmt.Sprintf("    if flag%s {\n%s\n    }", c.value, strings.Replace(s, "\n", "\n  ", -1))
		}
	}

	return s
}
