package jsparser // import "fknsrs.biz/p/jsparser"

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

type TokeniserError struct {
	Message string
	Offset  int
	Mode    LexicalState
}

func (t TokeniserError) Error() string {
	return fmt.Sprintf("TokeniserError (state %s at offset %d): %s", t.Mode, t.Offset, t.Message)
}

type TokenKind int

const (
	TokenKindWhitespace TokenKind = iota
	TokenKindBinaryAssignment
	TokenKindBinaryBitwiseAnd
	TokenKindBinaryBitwiseAndAssignment
	TokenKindBinaryBitwiseOr
	TokenKindBinaryBitwiseOrAssignment
	TokenKindBinaryBitwiseXor
	TokenKindBinaryBitwiseXorAssignment
	TokenKindBinaryDivide
	TokenKindBinaryDivideEquals
	TokenKindBinaryEquals
	TokenKindBinaryExponent
	TokenKindBinaryExponentAssignment
	TokenKindBinaryGreater
	TokenKindBinaryGreaterOrEqual
	TokenKindBinaryLess
	TokenKindBinaryLessOrEqual
	TokenKindBinaryLogicalAnd
	TokenKindBinaryLogicalOr
	TokenKindBinaryMinus
	TokenKindBinaryMinusAssignment
	TokenKindBinaryModulo
	TokenKindBinaryModuloAssignment
	TokenKindBinaryNotEquals
	TokenKindBinaryPlus
	TokenKindBinaryPlusAssignment
	TokenKindBinaryShiftLeft
	TokenKindBinaryShiftLeftAssignment
	TokenKindBinaryShiftRight
	TokenKindBinaryShiftRightAssignment
	TokenKindBinaryShiftRightUnsigned
	TokenKindBinaryShiftRightUnsignedAssignment
	TokenKindBinaryStar
	TokenKindBinaryStarAssignment
	TokenKindBinaryStrictEquals
	TokenKindBinaryStrictNotEquals
	TokenKindIdentifier
	TokenKindKeyword
	TokenKindMetaShebangLine
	TokenKindMultipleLineComment
	TokenKindNumber
	TokenKindPuncAt
	TokenKindPuncBacktick
	TokenKindPuncColon
	TokenKindPuncComma
	TokenKindPuncFatArrow
	TokenKindPuncLeftBrace
	TokenKindPuncLeftBracket
	TokenKindPuncLeftParen
	TokenKindPuncPeriod
	TokenKindPuncQuestion
	TokenKindPuncRightBrace
	TokenKindPuncRightBracket
	TokenKindPuncRightParen
	TokenKindPuncSemicolon
	TokenKindPuncSpread
	TokenKindRegexp
	TokenKindSingleLineComment
	TokenKindString
	TokenKindTemplateHead
	TokenKindTemplateMiddle
	TokenKindTemplateNoSubstitution
	TokenKindTemplateTail
	TokenKindUnaryBang
	TokenKindUnaryDecrement
	TokenKindUnaryIncrement
	TokenKindUnaryTilde
)

func (t TokenKind) String() string {
	switch t {
	case TokenKindBinaryAssignment:
		return "binaryAssignment"
	case TokenKindBinaryBitwiseAnd:
		return "binaryBitwiseAnd"
	case TokenKindBinaryBitwiseAndAssignment:
		return "binaryBitwiseAndAssignment"
	case TokenKindBinaryBitwiseOr:
		return "binaryBitwiseOr"
	case TokenKindBinaryBitwiseOrAssignment:
		return "binaryBitwiseOrAssignment"
	case TokenKindBinaryBitwiseXor:
		return "binaryBitwiseXor"
	case TokenKindBinaryBitwiseXorAssignment:
		return "binaryBitwiseXorAssignment"
	case TokenKindBinaryDivide:
		return "binaryDivide"
	case TokenKindBinaryDivideEquals:
		return "binaryDivideEquals"
	case TokenKindBinaryEquals:
		return "binaryEquals"
	case TokenKindBinaryExponent:
		return "binaryExponent"
	case TokenKindBinaryExponentAssignment:
		return "binaryExponentAssignment"
	case TokenKindBinaryGreater:
		return "binaryGreater"
	case TokenKindBinaryGreaterOrEqual:
		return "binaryGreaterOrEqual"
	case TokenKindBinaryLess:
		return "binaryLess"
	case TokenKindBinaryLessOrEqual:
		return "binaryLessOrEqual"
	case TokenKindBinaryLogicalAnd:
		return "binaryLogicalAnd"
	case TokenKindBinaryLogicalOr:
		return "binaryLogicalOr"
	case TokenKindBinaryMinus:
		return "binaryMinus"
	case TokenKindBinaryMinusAssignment:
		return "binaryMinusAssignment"
	case TokenKindBinaryModulo:
		return "binaryModulo"
	case TokenKindBinaryModuloAssignment:
		return "binaryModuloAssignment"
	case TokenKindBinaryNotEquals:
		return "binaryNotEquals"
	case TokenKindBinaryPlus:
		return "binaryPlus"
	case TokenKindBinaryPlusAssignment:
		return "binaryPlusAssignment"
	case TokenKindBinaryShiftLeft:
		return "binaryShiftLeft"
	case TokenKindBinaryShiftLeftAssignment:
		return "binaryShiftLeftAssignment"
	case TokenKindBinaryShiftRight:
		return "binaryShiftRight"
	case TokenKindBinaryShiftRightAssignment:
		return "binaryShiftRightAssignment"
	case TokenKindBinaryShiftRightUnsigned:
		return "binaryShiftRightUnsigned"
	case TokenKindBinaryShiftRightUnsignedAssignment:
		return "binaryShiftRightUnsignedAssignment"
	case TokenKindBinaryStar:
		return "binaryStar"
	case TokenKindBinaryStarAssignment:
		return "binaryStarAssignment"
	case TokenKindBinaryStrictEquals:
		return "binaryStrictEquals"
	case TokenKindBinaryStrictNotEquals:
		return "binaryStrictNotEquals"
	case TokenKindIdentifier:
		return "identifier"
	case TokenKindKeyword:
		return "keyword"
	case TokenKindMetaShebangLine:
		return "metaShebangLine"
	case TokenKindMultipleLineComment:
		return "multipleLineComment"
	case TokenKindNumber:
		return "number"
	case TokenKindPuncAt:
		return "puncAt"
	case TokenKindPuncBacktick:
		return "puncBacktick"
	case TokenKindPuncColon:
		return "puncColon"
	case TokenKindPuncComma:
		return "puncComma"
	case TokenKindPuncFatArrow:
		return "puncFatArrow"
	case TokenKindPuncLeftBrace:
		return "puncLeftBrace"
	case TokenKindPuncLeftBracket:
		return "puncLeftBracket"
	case TokenKindPuncLeftParen:
		return "puncLeftParen"
	case TokenKindPuncPeriod:
		return "puncPeriod"
	case TokenKindPuncQuestion:
		return "puncQuestion"
	case TokenKindPuncRightBrace:
		return "puncRightBrace"
	case TokenKindPuncRightBracket:
		return "puncRightBracket"
	case TokenKindPuncRightParen:
		return "puncRightParen"
	case TokenKindPuncSemicolon:
		return "puncSemicolon"
	case TokenKindPuncSpread:
		return "puncSpread"
	case TokenKindRegexp:
		return "regexp"
	case TokenKindSingleLineComment:
		return "singleLineComment"
	case TokenKindString:
		return "string"
	case TokenKindTemplateHead:
		return "templateHead"
	case TokenKindTemplateMiddle:
		return "templateMiddle"
	case TokenKindTemplateNoSubstitution:
		return "templateNoSubstitution"
	case TokenKindTemplateTail:
		return "templateTail"
	case TokenKindUnaryBang:
		return "unaryBang"
	case TokenKindUnaryDecrement:
		return "unaryDecrement"
	case TokenKindUnaryIncrement:
		return "unaryIncrement"
	case TokenKindUnaryTilde:
		return "unaryTilde"
	case TokenKindWhitespace:
		return "whitespace"
	default:
		return "unknown"
	}
}

type Token struct {
	Kind   TokenKind
	Value  string
	Raw    string
	Offset int
}

type TokenSet []Token

type LexicalState int

const (
	InputElementDiv LexicalState = iota
	InputElementRegExp
	InputElementRegExpOrTemplateTail
	InputElementTemplateTail
)

func (s LexicalState) String() string {
	switch s {
	case InputElementDiv:
		return "InputElementDiv"
	case InputElementRegExp:
		return "InputElementRegExp"
	case InputElementRegExpOrTemplateTail:
		return "InputElementRegExpOrTemplateTail"
	case InputElementTemplateTail:
		return "InputElementTemplateTail"
	default:
		return "unknown"
	}
}

type Tokeniser struct {
	state  LexicalState
	rd     *bufio.Reader
	buf    []rune
	regexp bool
	saved  int
	pos    int
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
		Mode:    t.state,
	}
}

func (t *Tokeniser) Mode() LexicalState {
	return t.state
}

func (t *Tokeniser) ReadAll() (TokenSet, error) {
	var a TokenSet

	for {
		tk, err := t.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		a = append(a, *tk)
	}

	return a, nil
}

func (t *Tokeniser) Read() (*Token, error) {
	if t.pos == 0 {
		r0, err := t.readRune()
		if err != nil {
			return nil, err
		}

		switch r0 {
		case '#':
			r1, err := t.readRune()
			if err != nil {
				return nil, err
			}

			switch r1 {
			case '!':
				b := []rune{r0, r1}
				for {
					r2, err := t.readRune()
					if err != nil {
						return nil, err
					}

					if r2 == '\n' {
						return &Token{
							Kind:   TokenKindMetaShebangLine,
							Value:  string(b),
							Raw:    string(b),
							Offset: t.save(),
						}, nil
					}

					b = append(b, r2)
				}
			default:
				t.unreadRune(r1, r0)
			}
		default:
			t.unreadRune(r0)
		}
	}

	var ws []rune
	// skip whitespace
	for {
		r, err := t.readRune()
		if err != nil {
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
		return &Token{Kind: TokenKindWhitespace, Raw: string(ws), Offset: t.save()}, nil
	}

	// try to read a static token
	r0, err := t.readRune()
	if err != nil {
		return nil, err
	}

	switch r0 {
	case '!':
		r1, err := t.readRune()
		if err != nil {
			return nil, err
		}

		switch r1 {
		case '=':
			r2, err := t.readRune()
			if err != nil {
				return nil, err
			}

			switch r2 {
			case '=':
				return &Token{Kind: TokenKindBinaryStrictNotEquals, Raw: "!==", Offset: t.save()}, nil
			}

			t.unreadRune(r2)

			return &Token{Kind: TokenKindBinaryNotEquals, Raw: "!=", Offset: t.save()}, nil
		}

		t.unreadRune(r1)

		return &Token{Kind: TokenKindUnaryBang, Raw: "!", Offset: t.save()}, nil
	case '%':
		r1, err := t.readRune()
		if err != nil {
			return nil, err
		}

		switch r1 {
		case '=':
			return &Token{Kind: TokenKindBinaryModuloAssignment, Raw: "%=", Offset: t.save()}, nil
		}

		t.unreadRune(r1)

		return &Token{Kind: TokenKindBinaryModulo, Raw: "%", Offset: t.save()}, nil
	case '&':
		r1, err := t.readRune()
		if err != nil {
			return nil, err
		}

		switch r1 {
		case '&':
			return &Token{Kind: TokenKindBinaryLogicalAnd, Raw: "&&", Offset: t.save()}, nil
		case '=':
			return &Token{Kind: TokenKindBinaryBitwiseAndAssignment, Raw: "&=", Offset: t.save()}, nil
		}

		t.unreadRune(r1)

		return &Token{Kind: TokenKindBinaryBitwiseAnd, Raw: "&", Offset: t.save()}, nil
	case '(':
		return &Token{Kind: TokenKindPuncLeftParen, Raw: "(", Offset: t.save()}, nil
	case ')':
		return &Token{Kind: TokenKindPuncRightParen, Raw: ")", Offset: t.save()}, nil
	case '*':
		r1, err := t.readRune()
		if err != nil {
			return nil, err
		}

		switch r1 {
		case '*':
			r2, err := t.readRune()
			if err != nil {
				return nil, err
			}

			switch r2 {
			case '=':
				return &Token{Kind: TokenKindBinaryExponentAssignment, Raw: "**=", Offset: t.save()}, nil
			}

			t.unreadRune(r2)

			return &Token{Kind: TokenKindBinaryExponent, Raw: "**", Offset: t.save()}, nil
		case '=':
			return &Token{Kind: TokenKindBinaryStarAssignment, Raw: "*=", Offset: t.save()}, nil
		}

		t.unreadRune(r1)

		return &Token{Kind: TokenKindBinaryStar, Raw: "*", Offset: t.save()}, nil
	case '+':
		r1, err := t.readRune()
		if err != nil {
			return nil, err
		}

		switch r1 {
		case '+':
			return &Token{Kind: TokenKindUnaryIncrement, Raw: "++", Offset: t.save()}, nil
		case '=':
			return &Token{Kind: TokenKindBinaryPlusAssignment, Raw: "+=", Offset: t.save()}, nil
		}

		t.unreadRune(r1)

		return &Token{Kind: TokenKindBinaryPlus, Raw: "+", Offset: t.save()}, nil
	case ',':
		return &Token{Kind: TokenKindPuncComma, Raw: ",", Offset: t.save()}, nil
	case '-':
		r1, err := t.readRune()
		if err != nil {
			return nil, err
		}

		switch r1 {
		case '-':
			return &Token{Kind: TokenKindUnaryDecrement, Raw: "--", Offset: t.save()}, nil
		case '=':
			return &Token{Kind: TokenKindBinaryMinusAssignment, Raw: "-=", Offset: t.save()}, nil
		}

		t.unreadRune(r1)

		return &Token{Kind: TokenKindBinaryMinus, Raw: "-", Offset: t.save()}, nil
	case '.':
		r1, err := t.readRune()
		if err != nil {
			return nil, err
		}

		switch r1 {
		case '.':
			r2, err := t.readRune()
			if err != nil {
				return nil, err
			}

			switch r2 {
			case '.':
				return &Token{Kind: TokenKindPuncSpread, Raw: "...", Offset: t.save()}, nil
			}

			t.unreadRune(r2)
		}

		t.unreadRune(r1)

		return &Token{Kind: TokenKindPuncPeriod, Raw: ".", Offset: t.save()}, nil
	case '/':
		r1, err := t.readRune()
		if err != nil {
			return nil, err
		}

		switch {
		case r1 == '/':
			t.unreadRune(r1, r0)
			return t.lexSingleLineComment()
		case r1 == '*':
			t.unreadRune(r1, r0)
			return t.lexMultipleLineComment()
		case t.state == InputElementRegExp || t.state == InputElementRegExpOrTemplateTail:
			t.unreadRune(r1, r0)
			return t.lexRegexp()
		case r1 == '=':
			return &Token{Kind: TokenKindBinaryDivideEquals, Raw: "/=", Offset: t.save()}, nil
		}

		t.unreadRune(r1)

		return &Token{Kind: TokenKindBinaryDivide, Raw: "/", Offset: t.save()}, nil
	case ':':
		return &Token{Kind: TokenKindPuncColon, Raw: ":", Offset: t.save()}, nil
	case ';':
		return &Token{Kind: TokenKindPuncSemicolon, Raw: ";", Offset: t.save()}, nil
	case '<':
		r1, err := t.readRune()
		if err != nil {
			return nil, err
		}

		switch r1 {
		case '<':
			r2, err := t.readRune()
			if err != nil {
				return nil, err
			}

			switch r2 {
			case '=':
				return &Token{Kind: TokenKindBinaryShiftLeftAssignment, Raw: "<<=", Offset: t.save()}, nil
			}

			t.unreadRune(r2)

			return &Token{Kind: TokenKindBinaryShiftLeft, Raw: "<<", Offset: t.save()}, nil
		case '=':
			return &Token{Kind: TokenKindBinaryLessOrEqual, Raw: "<=", Offset: t.save()}, nil
		}

		t.unreadRune(r1)

		return &Token{Kind: TokenKindBinaryLess, Raw: "<", Offset: t.save()}, nil
	case '=':
		r1, err := t.readRune()
		if err != nil {
			return nil, err
		}

		switch r1 {
		case '=':
			r2, err := t.readRune()
			if err != nil {
				return nil, err
			}

			switch r2 {
			case '=':
				return &Token{Kind: TokenKindBinaryStrictEquals, Raw: "===", Offset: t.save()}, nil
			}

			t.unreadRune(r2)

			return &Token{Kind: TokenKindBinaryEquals, Raw: "==", Offset: t.save()}, nil
		case '>':
			return &Token{Kind: TokenKindPuncFatArrow, Raw: "=>", Offset: t.save()}, nil
		}

		t.unreadRune(r1)

		return &Token{Kind: TokenKindBinaryAssignment, Raw: "=", Offset: t.save()}, nil
	case '>':
		r1, err := t.readRune()
		if err != nil {
			return nil, err
		}

		switch r1 {
		case '=':
			return &Token{Kind: TokenKindBinaryGreaterOrEqual, Raw: ">=", Offset: t.save()}, nil
		case '>':
			r2, err := t.readRune()
			if err != nil {
				return nil, err
			}

			switch r2 {
			case '=':
				return &Token{Kind: TokenKindBinaryShiftRightAssignment, Raw: ">>=", Offset: t.save()}, nil
			case '>':
				r3, err := t.readRune()
				if err != nil {
					return nil, err
				}

				switch r3 {
				case '=':
					return &Token{Kind: TokenKindBinaryShiftRightUnsignedAssignment, Raw: ">>>=", Offset: t.save()}, nil
				}

				t.unreadRune(r3)

				return &Token{Kind: TokenKindBinaryShiftRightUnsigned, Raw: ">>>", Offset: t.save()}, nil
			}

			t.unreadRune(r2)

			return &Token{Kind: TokenKindBinaryShiftRight, Raw: ">>", Offset: t.save()}, nil
		}

		t.unreadRune(r1)

		return &Token{Kind: TokenKindBinaryGreater, Raw: ">", Offset: t.save()}, nil
	case '?':
		return &Token{Kind: TokenKindPuncQuestion, Raw: "?", Offset: t.save()}, nil
	case '@':
		return &Token{Kind: TokenKindPuncAt, Raw: "@", Offset: t.save()}, nil
	case '[':
		return &Token{Kind: TokenKindPuncLeftBracket, Raw: "[", Offset: t.save()}, nil
	case ']':
		return &Token{Kind: TokenKindPuncRightBracket, Raw: "]", Offset: t.save()}, nil
	case '`':
		t.unreadRune(r0)
		return t.lexTemplateHead()
	case '^':
		r1, err := t.readRune()
		if err != nil {
			return nil, err
		}

		switch r1 {
		case '=':
			return &Token{Kind: TokenKindBinaryBitwiseXorAssignment, Raw: "^=", Offset: t.save()}, nil
		}

		t.unreadRune(r1)

		return &Token{Kind: TokenKindBinaryBitwiseXor, Raw: "^", Offset: t.save()}, nil
	case '{':
		return &Token{Kind: TokenKindPuncLeftBrace, Raw: "{", Offset: t.save()}, nil
	case '|':
		r1, err := t.readRune()
		if err != nil {
			return nil, err
		}

		switch r1 {
		case '=':
			return &Token{Kind: TokenKindBinaryBitwiseOrAssignment, Raw: "|=", Offset: t.save()}, nil
		case '|':
			return &Token{Kind: TokenKindBinaryLogicalOr, Raw: "||", Offset: t.save()}, nil
		}

		t.unreadRune(r1)

		return &Token{Kind: TokenKindBinaryBitwiseOr, Raw: "|", Offset: t.save()}, nil
	case '}':
		switch t.state {
		case InputElementTemplateTail, InputElementRegExpOrTemplateTail:
			t.unreadRune(r0)
			return t.lexTemplateTail()
		default:
			return &Token{Kind: TokenKindPuncRightBrace, Raw: "}", Offset: t.save()}, nil
		}
	case '~':
		return &Token{Kind: TokenKindUnaryTilde, Raw: "~", Offset: t.save()}, nil
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
		case unicode.Is(unicode.N, r):
			t.unreadRune(r)

			return t.lexNumber()
		case r == '$' || r == '_' || unicode.Is(unicode.L, r) || unicode.Is(unicode.M, r):
			t.unreadRune(r)

			return t.lexIdentifier()
		default:
			return nil, t.errf("unexpected character %q", r)
		}
	}
}

func (t *Tokeniser) lexTemplateHead() (*Token, error) {
	r0, err := t.readRune()
	if err != nil {
		return nil, err
	}

	if r0 != '`' {
		return nil, t.errf("unexpected character %q", r0)
	}

	b := []rune{r0}
	v := ""

	for {
		r1, err := t.readRune()
		if err != nil {
			return nil, err
		}

		switch r1 {
		case '`':
			b = append(b, r1)

			return &Token{
				Kind:   TokenKindTemplateNoSubstitution,
				Raw:    string(b),
				Value:  v,
				Offset: t.save(),
			}, nil
		case '$':
			r2, err := t.readRune()
			if err != nil {
				return nil, err
			}

			switch r2 {
			case '{':
				b = append(b, r1, r2)

				return &Token{
					Kind:   TokenKindTemplateHead,
					Raw:    string(b),
					Value:  v,
					Offset: t.save(),
				}, nil
			}

			t.unreadRune(r2)
		}

		b = append(b, r1)
		v = v + string(r1)
	}
}

func (t *Tokeniser) lexTemplateTail() (*Token, error) {
	r0, err := t.readRune()
	if err != nil {
		return nil, err
	}

	if r0 != '}' {
		return nil, t.errf("unexpected character %q", r0)
	}

	b := []rune{r0}
	v := ""

	for {
		r1, err := t.readRune()
		if err != nil {
			return nil, err
		}

		switch r1 {
		case '`':
			b = append(b, r1)

			return &Token{
				Kind:   TokenKindTemplateTail,
				Raw:    string(b),
				Value:  v,
				Offset: t.save(),
			}, nil
		case '$':
			r2, err := t.readRune()
			if err != nil {
				return nil, err
			}

			switch r2 {
			case '{':
				b = append(b, r1, r2)

				return &Token{
					Kind:   TokenKindTemplateMiddle,
					Raw:    string(b),
					Value:  v,
					Offset: t.save(),
				}, nil
			}

			t.unreadRune(r2)

			b = append(b, r1)
			v = v + string(r1)
		}
	}
}

func (t *Tokeniser) lexNumber() (*Token, error) {
	var b []rune

	seenPeriod := false

	for {
		r, err := t.readRune()
		if err != nil {
			return nil, err
		}

		if unicode.Is(unicode.N, r) {
			b = append(b, r)
		} else if r == '.' && !seenPeriod {
			b = append(b, r)
			seenPeriod = true
		} else {
			t.unreadRune(r)

			break
		}
	}

	return &Token{
		Kind:   TokenKindNumber,
		Value:  string(b),
		Raw:    string(b),
		Offset: t.save(),
	}, nil
}

func (t *Tokeniser) lexIdentifier() (*Token, error) {
	var b []rune

	for {
		r, err := t.readRune()
		if err != nil {
			return nil, err
		}

		if r == '$' || r == '_' || unicode.Is(unicode.L, r) || unicode.Is(unicode.M, r) || unicode.Is(unicode.N, r) {
			b = append(b, r)
		} else {
			t.unreadRune(r)

			break
		}
	}

	tk := Token{
		Kind:   TokenKindIdentifier,
		Value:  string(b),
		Raw:    string(b),
		Offset: t.save(),
	}

	return &tk, nil
}

func (t *Tokeniser) lexString() (*Token, error) {
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

	return &Token{
		Kind:   TokenKindString,
		Value:  buf,
		Raw:    string(raw),
		Offset: t.save(),
	}, nil
}

func (t *Tokeniser) lexRegexp() (*Token, error) {
	var raw []rune

	d, err := t.readRune()
	if err != nil {
		return nil, err
	}
	raw = append(raw, d)

	esc := false
loop1:
	for {
		r, err := t.readRune()
		if err != nil {
			return nil, err
		}

		if esc {
			esc = false
		} else {
			switch r {
			case '\\':
				esc = true
			case '/':
				t.unreadRune(r)
				break loop1
			}
		}

		raw = append(raw, r)
	}

	r, err := t.readRune()
	if err != nil {
		return nil, err
	}
	raw = append(raw, r)

	if r != d {
		return nil, t.errf("invalid regexp literal")
	}

loop2:
	for {
		r, err := t.readRune()
		if err != nil {
			return nil, err
		}

		if !unicode.IsLower(r) {
			t.unreadRune(r)
			break loop2
		}

		raw = append(raw, r)
	}

	return &Token{
		Kind:   TokenKindRegexp,
		Value:  string(raw),
		Raw:    string(raw),
		Offset: t.save(),
	}, nil
}

func (t *Tokeniser) lexSingleLineComment() (*Token, error) {
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

	return &Token{
		Kind:   TokenKindSingleLineComment,
		Raw:    string(buf),
		Value:  string(buf),
		Offset: t.save(),
	}, nil
}

func (t *Tokeniser) lexMultipleLineComment() (*Token, error) {
	r0, err := t.readRune()
	if err != nil {
		return nil, err
	}
	if r0 != '/' {
		return nil, t.errf("invalid multi-line first opening boundary character %q", r0)
	}

	r1, err := t.readRune()
	if err != nil {
		return nil, err
	}
	if r1 != '*' {
		return nil, t.errf("invalid multi-line second opening boundary character %q", r1)
	}

	buf := []rune{r0, r1}
loop:
	for {
		r0, err := t.readRune()
		if err != nil {
			return nil, err
		}

		switch r0 {
		case '*':
			r1, err := t.readRune()
			if err != nil {
				return nil, err
			}

			switch r1 {
			case '/':
				buf = append(buf, r0, r1)
				break loop
			}

			t.unreadRune(r1)
		}

		buf = append(buf, r0)
	}

	return &Token{
		Kind:   TokenKindMultipleLineComment,
		Raw:    string(buf),
		Value:  string(buf),
		Offset: t.save(),
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
