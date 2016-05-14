package main

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

type LexicalState int

const (
	TextState LexicalState = iota
	CodeState
)

type TokenKind int

const (
	TokenKindWhitespace TokenKind = iota
	TokenKindEOF
	TokenKindText
	TokenKindFence
	TokenKindComment
	TokenKindIdentifier
	TokenKindKeywordEnum
	TokenKindKeywordInterface
	TokenKindString
	TokenKindComma
	TokenKindColon
	TokenKindSemicolon
	TokenKindInherits
	TokenKindLeftBracket
	TokenKindRightBracket
	TokenKindLeftBrace
	TokenKindPipe
	TokenKindRightBrace
)

func (t TokenKind) String() string {
	switch t {
	case TokenKindWhitespace:
		return "Whitespace"
	case TokenKindEOF:
		return "EOF"
	case TokenKindText:
		return "Text"
	case TokenKindFence:
		return "Fence"
	case TokenKindComment:
		return "Comment"
	case TokenKindIdentifier:
		return "Identifier"
	case TokenKindKeywordEnum:
		return "KeywordEnum"
	case TokenKindKeywordInterface:
		return "KeywordInterface"
	case TokenKindString:
		return "String"
	case TokenKindComma:
		return "Comma"
	case TokenKindColon:
		return "Colon"
	case TokenKindSemicolon:
		return "Semicolon"
	case TokenKindInherits:
		return "Inherits"
	case TokenKindLeftBracket:
		return "LeftBracket"
	case TokenKindRightBracket:
		return "RightBracket"
	case TokenKindLeftBrace:
		return "LeftBrace"
	case TokenKindPipe:
		return "Pipe"
	case TokenKindRightBrace:
		return "RightBrace"
	default:
		return "unknown"
	}
}

type Token interface {
	Kind() TokenKind
	Source() string
	Position() int
}

type EOFToken struct {
	position int
}

func (t EOFToken) Kind() TokenKind { return TokenKindEOF }
func (t EOFToken) Source() string  { return "" }
func (t EOFToken) Position() int   { return t.position }

type BasicToken struct {
	kind     TokenKind
	source   string
	position int
}

func (t BasicToken) Kind() TokenKind { return t.kind }
func (t BasicToken) Source() string  { return t.source }
func (t BasicToken) Position() int   { return t.position }

type FenceToken struct {
	BasicToken
	meta string
}

func (t FenceToken) Meta() string { return t.meta }

type IdentifierToken struct {
	BasicToken
	value string
}

func (t IdentifierToken) Value() string { return t.value }

var keywords = map[string]TokenKind{
	"enum":      TokenKindKeywordEnum,
	"interface": TokenKindKeywordInterface,
}

type StringToken struct {
	BasicToken
	value string
}

func (t StringToken) Value() string { return t.value }

type CommentToken struct {
	BasicToken
}

type TokenSet []Token

type TokeniserError struct {
	Message string
	Offset  int
}

func (t TokeniserError) Error() string {
	return fmt.Sprintf("TokeniserError (offset %d): %s", t.Offset, t.Message)
}

type Tokeniser struct {
	rd    *bufio.Reader
	buf   []rune
	saved int
	pos   int
}

func NewTokeniser(rd io.Reader) *Tokeniser {
	return &Tokeniser{rd: bufio.NewReader(rd)}
}

func (t *Tokeniser) save() int {
	s := t.saved
	t.saved = t.pos
	return s
}

func (t *Tokeniser) errf(format string, a ...interface{}) TokeniserError {
	return TokeniserError{
		Message: fmt.Sprintf(format, a...),
		Offset:  t.saved,
	}
}

func (t *Tokeniser) Read(s LexicalState) (Token, error) {
	if _, err := t.rd.Peek(1); err != nil {
		if err == io.EOF {
			return EOFToken{position: t.pos}, nil
		}

		return nil, err
	}

	// look for ```
	{
		r0, err := t.readRune()
		if err != nil {
			return nil, err
		}

		switch r0 {
		case '`':
			r1, err := t.readRune()
			if err != nil {
				return nil, err
			}

			switch r1 {
			case '`':
				r2, err := t.readRune()
				if err != nil {
					return nil, err
				}

				switch r2 {
				case '`':
					var b []rune

				loop:
					for {
						r, err := t.readRune()
						if err != nil {
							return nil, err
						}

						switch r {
						case '\r', '\n':
							t.unreadRune(r)
							break loop
						default:
							b = append(b, r)
						}
					}

					return FenceToken{
						BasicToken: BasicToken{kind: TokenKindFence, source: "```" + string(b), position: t.save()},
						meta:       string(b),
					}, nil
				}

				t.unreadRune(r2)
			}

			t.unreadRune(r1)
		}

		t.unreadRune(r0)
	}

	if s == TextState {
		return t.lexText()
	}

	var ws []rune
	// read whitespace
	for {
		r, err := t.readRune()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		if !unicode.IsSpace(r) {
			t.unreadRune(r)
			break
		} else {
			ws = append(ws, r)
		}
	}
	if len(ws) > 0 {
		return BasicToken{kind: TokenKindWhitespace, source: string(ws), position: t.save()}, nil
	}

	// try to read a static token
	r0, err := t.readRune()
	if err != nil {
		return nil, err
	}

	switch r0 {
	case ',':
		return BasicToken{kind: TokenKindComma, source: ",", position: t.save()}, nil
	case '/':
		r1, err := t.readRune()
		if err != nil {
			return nil, err
		}

		switch {
		case r1 == '/':
			t.unreadRune(r1, r0)
			return t.lexSingleLineComment()
		}

		t.unreadRune(r1)
	case ':':
		return BasicToken{kind: TokenKindColon, source: ":", position: t.save()}, nil
	case ';':
		return BasicToken{kind: TokenKindSemicolon, source: ";", position: t.save()}, nil
	case '<':
		r1, err := t.readRune()
		if err != nil {
			return nil, err
		}

		switch r1 {
		case ':':
			return BasicToken{kind: TokenKindInherits, source: "<:", position: t.save()}, nil
		}

		t.unreadRune(r1)
	case '[':
		return BasicToken{kind: TokenKindLeftBracket, source: "[", position: t.save()}, nil
	case ']':
		return BasicToken{kind: TokenKindRightBracket, source: "]", position: t.save()}, nil
	case '{':
		return BasicToken{kind: TokenKindLeftBrace, source: "{", position: t.save()}, nil
	case '|':
		return BasicToken{kind: TokenKindPipe, source: "|", position: t.save()}, nil
	case '}':
		return BasicToken{kind: TokenKindRightBrace, source: "}", position: t.save()}, nil
	}

	t.unreadRune(r0)

	for {
		r, err := t.readRune()
		if err != nil {
			return nil, err
		}

		switch {
		case r == '"' || r == '\'':
			t.unreadRune(r)

			return t.lexString()
		case unicode.Is(unicode.L, r):
			t.unreadRune(r)

			return t.lexIdentifier()
		default:
			return nil, t.errf("unexpected character %q", r)
		}
	}
}

func (t *Tokeniser) lexText() (Token, error) {
	var b []rune
loop:
	for {
		r0, err := t.readRune()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		switch r0 {
		case '`':
			r1, err := t.readRune()
			if err != nil {
				return nil, err
			}

			switch r1 {
			case '`':
				r2, err := t.readRune()
				if err != nil {
					return nil, err
				}

				switch r2 {
				case '`':
					t.unreadRune(r2, r1, r0)
					break loop
				}

				t.unreadRune(r2)
			}

			t.unreadRune(r1)
		default:
			b = append(b, r0)
		}
	}

	return BasicToken{
		kind:     TokenKindText,
		source:   string(b),
		position: t.save(),
	}, nil
}

func (t *Tokeniser) lexIdentifier() (Token, error) {
	var b []rune

	for {
		r, err := t.readRune()
		if err != nil {
			return nil, err
		}

		if unicode.Is(unicode.L, r) {
			b = append(b, r)
		} else {
			t.unreadRune(r)

			break
		}
	}

	k := TokenKindIdentifier
	if kw, ok := keywords[string(b)]; ok {
		k = kw
	}

	return IdentifierToken{
		BasicToken: BasicToken{
			kind:     k,
			source:   string(b),
			position: t.save(),
		},
		value: string(b),
	}, nil
}

func (t *Tokeniser) lexString() (Token, error) {
	var buf string
	var raw []rune

	q, err := t.readRune()
	if err != nil {
		return nil, err
	}
	raw = append(raw, q)

	esc := false
loop:
	for {
		r, err := t.readRune()
		if err != nil {
			return nil, err
		}

		if esc {
			esc = false

			switch r {
			case '\n':
				// nothing
			case q:
				buf += string(r)
			case '\\':
				buf += "\\"
			case 'n':
				buf += "\n"
			case 'r':
				buf += "\r"
			case 't':
				buf += "\t"
			}
		} else {
			switch r {
			case '\\':
				esc = true
			case q:
				t.unreadRune(r)
				break loop
			case '\n':
				t.unreadRune(r)
				break loop
			default:
				buf += string(r)
			}
		}

		raw = append(raw, r)
	}

	r, err := t.readRune()
	if err != nil {
		return nil, err
	}
	raw = append(raw, r)

	if r != q {
		return nil, t.errf("invalid string literal")
	}

	return StringToken{
		BasicToken: BasicToken{
			kind:     TokenKindString,
			source:   string(raw),
			position: t.save(),
		},
		value: buf,
	}, nil
}

func (t *Tokeniser) lexSingleLineComment() (Token, error) {
	r0, err := t.readRune()
	if err != nil {
		return nil, err
	}
	if r0 != '/' {
		return nil, t.errf("invalid single-line first opening character %q", r0)
	}

	r1, err := t.readRune()
	if err != nil {
		return nil, err
	}
	if r1 != '/' {
		return nil, t.errf("invalid single-line second opening character %q", r1)
	}

	buf := []rune{r0, r1}
	for {
		r, err := t.readRune()
		if err != nil {
			return nil, err
		}

		if r == '\n' {
			t.unreadRune(r)
			break
		}

		buf = append(buf, r)
	}

	return CommentToken{
		BasicToken: BasicToken{
			kind:     TokenKindComment,
			source:   string(buf),
			position: t.save(),
		},
	}, nil
}

func (t *Tokeniser) readRune() (rune, error) {
	if len(t.buf) > 0 {
		r := t.buf[len(t.buf)-1]
		t.buf = t.buf[0 : len(t.buf)-1]
		t.pos += utf8.RuneLen(r)
		return r, nil
	}

	r, n, err := t.rd.ReadRune()
	t.pos += n
	return r, err
}

func (t *Tokeniser) unreadRune(r ...rune) {
	for _, r := range r {
		t.buf = append(t.buf, r)
		t.pos -= utf8.RuneLen(r)
	}
}
