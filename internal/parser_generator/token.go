package main

import (
	"fmt"
	"strings"
	"unicode"
)

func quoted(s string) string {
	if strings.Contains(s, "\"") {
		return "'" + strings.NewReplacer("\\", "\\\\", "'", "\\'").Replace(s) + "'"
	}

	return "\"" + strings.NewReplacer("\\", "\\\\", "\"", "\\\"").Replace(s) + "\""
}

type token struct {
	value    string
	params   []string
	optional bool
	terminal bool
	oneof    bool
	values   []string
}

func (t token) String() string {
	s := t.value
	if t.terminal {
		s = fmt.Sprintf("%q", t.value)
	}
	if t.oneof {
		s = fmt.Sprintf("one of %v", t.values)
	}

	if len(t.params) > 0 {
		s += "[" + strings.Join(t.params, ", ") + "]"
	}

	if t.optional {
		s += "?"
	}

	return s
}

func (t token) EBNF(options map[string]bool) string {
	s := t.value
	switch {
	case t.terminal:
		s = fmt.Sprintf("%s", quoted(t.value))
	case t.oneof:
		c := true
		for _, e := range t.values {
			if len(e) > 1 || (!unicode.IsLetter(rune(e[0])) && !unicode.IsNumber(rune(e[0]))) {
				c = false
			}
		}

		if c {
			s = fmt.Sprintf("[%s]", strings.Join(t.values, ""))
		} else {
			a := make([]string, len(t.values))
			for i, e := range t.values {
				a[i] = fmt.Sprintf("%s", quoted(e))
			}

			s = fmt.Sprintf("( %s )", strings.Join(a, " | "))
		}
	default:
		for _, p := range t.params {
			maybe := false
			if p[0] == '?' {
				maybe = true
				p = p[1:]
			}

			if maybe {
				if options[p] {
					s = s + "_" + p
				}
			} else {
				s = s + "_" + p
			}
		}

		if t.optional {
			s = s + "?"
		}
	}

	return s
}

func (t token) Go(m map[string]*production, p production, inner string) string {
	ok := "ok"
	ifok := "else if ok"
	if t.optional {
		ok = "_"
		ifok = "else"
	}

	switch {
	case t.terminal:
		return fmt.Sprintf(
			"if r, %s, err := acceptString(r, %q); err != nil { return input, false, err } %s { %s }\n",
			ok,
			t.value,
			ifok,
			inner,
		)
	case t.oneof:
		return fmt.Sprintf(
			"if r, %s, err := acceptOne(r, %#v); err != nil { return input, false, err } %s { %s }\n",
			ok,
			t.values,
			ifok,
			inner,
		)
	default:
		args := make([]string, len(m[t.value].params)+1)
		args[0] = "r"
		for i, v := range m[t.value].params {
			passThrough := false
			for _, s := range p.params {
				if s == v {
					passThrough = true
				}
			}

			setTrue := false
			for _, s := range t.params {
				if s == v {
					setTrue = true
				}
			}

			if setTrue {
				args[i+1] = "true"
			} else if passThrough {
				args[i+1] = "flag" + v
			} else {
				args[i+1] = "false"
			}
		}

		return fmt.Sprintf(
			"if r, %s, err := parse%s(%s); err != nil { return input, false, err } %s { %s }\n",
			ok,
			t.value,
			strings.Join(args, ", "),
			ifok,
			inner,
		)
	}
}
