package main

import (
	"fmt"
)

type Parser struct {
	buf   []Token
	t     *Tokeniser
	ls    LexicalState
	enums []esEnum
	types []esType
}

func NewParser(t *Tokeniser) *Parser {
	return &Parser{t: t}
}

func (p *Parser) Type(name string) *esType {
	for _, t := range p.types {
		if t.name == name {
			return &t
		}
	}

	return nil
}

func (p *Parser) read() (Token, error) {
	if len(p.buf) > 0 {
		tk := p.buf[len(p.buf)-1]
		p.buf = p.buf[0 : len(p.buf)-1]
		return tk, nil
	}

	tk, err := p.t.Read(p.ls)
	if err != nil {
		return nil, err
	}

	if tk.Kind() == TokenKindWhitespace {
		return p.read()
	}

	return tk, nil
}

func (p *Parser) unread(tk Token) {
	p.buf = append(p.buf, tk)
}

func (p *Parser) expect(kinds ...TokenKind) (Token, error) {
	tk, err := p.read()
	if err != nil {
		return nil, err
	}

	for _, k := range kinds {
		if tk.Kind() == k {
			return tk, nil
		}
	}

	p.unread(tk)

	return nil, fmt.Errorf("unexpected %v at offset %v (expected one of %v)", tk.Kind(), tk.Position(), kinds)
}

func (p *Parser) accept(kinds ...TokenKind) (Token, error) {
	tk, err := p.read()
	if err != nil {
		return nil, err
	}

	for _, k := range kinds {
		if tk.Kind() == k {
			return tk, nil
		}
	}

	p.unread(tk)

	return nil, nil
}

func (p *Parser) parse() error {
loop:
	for {
		switch p.ls {
		case TextState:
			tk, err := p.expect(TokenKindEOF, TokenKindText, TokenKindFence)
			if err != nil {
				return err
			}

			switch tk.Kind() {
			case TokenKindEOF:
				break loop
			case TokenKindFence:
				p.ls = CodeState
			}
		case CodeState:
			tk, err := p.expect(TokenKindFence, TokenKindKeywordEnum, TokenKindKeywordInterface)
			if err != nil {
				return err
			}

			switch tk.Kind() {
			case TokenKindFence:
				p.ls = TextState
			case TokenKindKeywordEnum:
				p.unread(tk)
				if err := p.parseEnum(); err != nil {
					return err
				}
			case TokenKindKeywordInterface:
				p.unread(tk)
				if err := p.parseInterface(); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (p *Parser) parseEnum() error {
	if _, err := p.expect(TokenKindKeywordEnum); err != nil {
		return err
	}

	id, err := p.expect(TokenKindIdentifier)
	if err != nil {
		return err
	}

	if _, err := p.expect(TokenKindLeftBrace); err != nil {
		return err
	}

	var a []StringToken
	for {
		s, err := p.accept(TokenKindString)
		if err != nil {
			return err
		}

		if s == nil {
			break
		}

		a = append(a, s.(StringToken))

		if _, err := p.accept(TokenKindPipe); err != nil {
			return err
		}
	}

	if _, err := p.expect(TokenKindRightBrace); err != nil {
		return err
	}

	p.enums = append(p.enums, esEnum{
		name:   id.(IdentifierToken),
		values: a,
	})

	return nil
}

func (p *Parser) parseInterface() error {
	if _, err := p.expect(TokenKindKeywordInterface); err != nil {
		return err
	}

	id, err := p.expect(TokenKindIdentifier)
	if err != nil {
		return err
	}

	var extends []string
	if t, err := p.accept(TokenKindInherits); err != nil {
		return err
	} else if t != nil {
		for {
			id, err := p.expect(TokenKindIdentifier)
			if err != nil {
				return err
			}

			extends = append(extends, id.(IdentifierToken).Value())

			if t, err := p.accept(TokenKindComma); err != nil {
				return err
			} else if t == nil {
				break
			}
		}
	}

	if _, err := p.expect(TokenKindLeftBrace); err != nil {
		return err
	}

	var fields []esTypeField
	for {
		if t, err := p.accept(TokenKindRightBrace); err != nil {
			return err
		} else if t != nil {
			p.unread(t)
			break
		}

		var f esTypeField

		id, err := p.expect(TokenKindIdentifier)
		if err != nil {
			return err
		}
		f.name = id.(IdentifierToken).Value()

		if _, err := p.expect(TokenKindColon); err != nil {
			return err
		}

		if t, err := p.accept(TokenKindLeftBracket); err != nil {
			return err
		} else if t != nil {
			f.list = true
		}

		expected := []TokenKind{TokenKindString, TokenKindIdentifier, TokenKindSemicolon}
		if f.list {
			expected = append(expected, TokenKindRightBracket)
		}

	loop:
		for {
			tk, err := p.expect(expected...)
			if err != nil {
				return err
			}

			switch tk.Kind() {
			case TokenKindString, TokenKindIdentifier:
				f.opts = append(f.opts, tk)
			case TokenKindSemicolon, TokenKindRightBracket:
				p.unread(tk)
				break loop
			}

			if tk, err := p.accept(TokenKindPipe); err != nil {
				return err
			} else if tk != nil {
				if tk, err := p.expect(TokenKindIdentifier, TokenKindString); err != nil {
					return err
				} else if tk != nil {
					p.unread(tk)
				}

				continue
			}

			break
		}

		if f.list {
			if _, err := p.expect(TokenKindRightBracket); err != nil {
				return err
			}
		}

		if _, err := p.expect(TokenKindSemicolon); err != nil {
			return err
		}

		fields = append(fields, f)
	}

	if _, err := p.expect(TokenKindRightBrace); err != nil {
		return err
	}

	p.types = append(p.types, esType{
		name:    id.(IdentifierToken).Value(),
		extends: extends,
		fields:  fields,
	})

	return nil
}
