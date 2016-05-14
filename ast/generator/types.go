package main

import (
	"fmt"
	"strings"
)

type esEnum struct {
	name   IdentifierToken
	values []StringToken
}

func (e esEnum) String() string {
	body := make([]string, len(e.values))
	for i, v := range e.values {
		body[i] = v.source
	}

	return fmt.Sprintf("enum %s { %s }", e.name.source, strings.Join(body, " | "))
}

type esType struct {
	name    string
	extends []string
	fields  []esTypeField
}

func (t esType) String() string {
	body := make([]string, len(t.fields))
	for i, f := range t.fields {
		body[i] = "  " + f.String()
	}

	s := strings.Join(body, "\n")
	if len(s) > 0 {
		s = "\n" + s + "\n"
	}

	e := strings.Join(t.extends, ", ")
	if len(e) > 0 {
		e = " <: " + e
	}

	return fmt.Sprintf("interface %s%s {%s}", t.name, e, s)
}

type esTypeField struct {
	name string
	list bool
	opts []Token
}

func (f esTypeField) Static() bool {
	if f.list == false && len(f.opts) == 1 && f.opts[0].Kind() == TokenKindString {
		return true
	}

	return false
}

func (f esTypeField) Type() string {
	m := make(map[string]int)
	var l []string

	maybeNull := false

	for _, o := range f.opts {
		var t string
		if o.Kind() == TokenKindString {
			t = "string"
		} else {
			t = o.(IdentifierToken).Value()
		}

		if t == "null" {
			maybeNull = true
			continue
		}

		if _, ok := m[t]; !ok {
			m[t] = 0
			l = append(l, t)
		}

		m[t]++
	}

	p := ""
	if f.list {
		p = "[]"
	}

	if len(l) == 1 {
		if maybeNull {
			return p + "*" + l[0]
		}

		return p + l[0]
	}

	a := make([]string, len(l))
	for i, v := range l {
		a[i] = fmt.Sprintf("%s %s", strings.Title(v), strings.Title(v))
	}

	return p + fmt.Sprintf("struct { %s }", strings.Join(a, "; "))
}

func (f esTypeField) String() string {
	body := make([]string, len(f.opts))
	for i, o := range f.opts {
		body[i] = o.Source()
	}

	s := strings.Join(body, " | ")
	if f.list {
		s = "[ " + s + " ]"
	}

	return fmt.Sprintf("%s: %s;", f.name, s)
}
