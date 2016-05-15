package main

import (
	"fmt"
	"strings"
)

type production struct {
	name   string
	params []string
	rules  []rule
}

func (p production) String() string {
	s := p.name

	if len(p.params) > 0 {
		s += "[" + strings.Join(p.params, ", ") + "]"
	}

	s += " ::\n"

	for _, r := range p.rules {
		s += "  " + r.String() + "\n"
	}

	return s
}

func (p production) EBNF() []string {
	var a []string

	for _, v := range variants(p.name, p.params) {
		s := fmt.Sprintf("%s", strings.Join(v, "_"))

		m := make(map[string]bool)
		for _, e := range v[1:] {
			m[e] = true
		}

		i := 0
		for _, r := range p.rules {
			if rs := r.EBNF(m); rs != "" {
				if i == 0 {
					s += "\n  ::= "
				} else {
					s += "\n    | "
				}

				s += rs

				i++
			}
		}

		a = append(a, s)
	}

	return a
}

func (p production) Go(m map[string]*production) string {
	var args []string
	for _, a := range p.params {
		args = append(args, "flag"+a+" bool")
	}

	s := fmt.Sprintf("func parse%s(input string, %s) (string, bool, error) {\n", p.name, strings.Join(args, ", "))

	for _, r := range p.rules {
		s += fmt.Sprintf("  { %s }\n", r.Go(m, p))
	}

	s += "\n  return input, false, nil\n"

	s += "}"

	return s
}
