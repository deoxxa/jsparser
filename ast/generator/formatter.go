package main

import (
	"fmt"
	"io"
	"strings"
)

type Formatter struct {
	w io.Writer
}

type formattingContext struct {
	p *Parser
	w io.Writer
	n int
	e error

	ts map[string]bool
	ti map[string]bool
	u  map[string]bool
}

func (c *formattingContext) f(format string, a ...interface{}) {
	if c.e == nil {
		n, err := fmt.Fprintf(c.w, format, a...)
		c.n += n
		c.e = err
	}
}

func NewFormatter(w io.Writer) *Formatter {
	return &Formatter{w: w}
}

func (f *Formatter) Format(p *Parser) (int, error) {
	c := formattingContext{
		p: p,
		w: f.w,

		ts: make(map[string]bool),
		ti: make(map[string]bool),
		u:  make(map[string]bool),
	}

	for _, e := range p.enums {
		f.formatEnum(&c, e)
	}

	for _, t := range p.types {
		f.formatTypeAsStruct(&c, t)
	}

	// for _, t := range p.types {
	// 	f.formatTypeAsInterface(&c, t)
	// }

	return c.n, c.e
}

func (f *Formatter) formatEnum(c *formattingContext, e esEnum) {
	w := c.f

	w("type %s string\n\n", e.name.Value())

	w("const (\n")
	for i, v := range e.values {
		w("  %s_%d = %q\n", e.name.Value(), i, v.Value())
	}
	w(")\n\n")

	w("func (v %s) Valid() bool { return ", e.name.Value())
	for i := range e.values {
		w("v == %s_%d", e.name.Value(), i)
		if i != len(e.values)-1 {
			w(" || ")
		}
	}
	w(" }\n\n")
}

func (f *Formatter) formatTypeAsInterface(c *formattingContext, t esType) {
	if c.ti[t.name] {
		return
	}

	c.ti[t.name] = true

	w := c.f

	w("type Is%s interface {\n", strings.Title(t.name))
	for _, e := range t.extends {
		w("  Is%s\n", strings.Title(e))
	}

	for _, tf := range t.fields {
		w("  Get%s() %s\n", strings.Title(tf.name), f.formatFieldType(c, tf))
	}
	w("}\n\n")
}

func (f *Formatter) formatTypeAsStruct(c *formattingContext, t esType) {
	if c.ts[t.name] {
		return
	}

	c.ts[t.name] = true

	for _, tf := range t.fields {
		f.maybeMakeUnion(c, tf)
	}

	w := c.f

	w("type %s struct {\n", t.name)
	for _, e := range t.extends {
		w("  %s\n", e)
	}

	for _, tf := range t.fields {
		if tf.Static() {
			continue
		}

		w("  %s %s `json:\"%s\"`\n", strings.Title(tf.name), f.formatFieldType(c, tf), tf.name)
	}
	w("}\n\n")
}

func (f *Formatter) maybeMakeUnion(c *formattingContext, tf esTypeField) {
	m := make(map[string]bool)
	var l []string

	for _, o := range tf.opts {
		var t string
		if o.Kind() == TokenKindString {
			t = "string"
		} else {
			t = o.(IdentifierToken).Value()
		}

		if t == "null" {
			continue
		}

		if _, ok := m[t]; !ok {
			m[t] = true
			l = append(l, t)
		}
	}

	if len(l) < 2 {
		return
	}

	n := strings.Join(l, "Or")

	if c.u[n] {
		return
	}

	c.u[n] = true

	w := c.f

	w("type %s interface { Is%s() bool }\n", n, n)
	for _, s := range l {
		w("func (%s) Is%s() bool { return true }\n", s, n)
	}
	w("\n")
}

func (f *Formatter) formatFieldType(c *formattingContext, tf esTypeField) string {
	m := make(map[string]int)
	var l []string

	maybeNull := false

	for _, o := range tf.opts {
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
	if tf.list {
		p = "[]"
	}

	if len(l) == 1 {
		if maybeNull {
			return p + "*" + l[0]
		}

		return p + l[0]
	}

	return p + strings.Join(l, "Or")
}
