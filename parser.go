package jsparser // import "fknsrs.biz/p/jsparser"

import (
	"bytes"
	"io"
	"strings"
)

type Keyword int

const (
	KeywordBreak Keyword = iota
	KeywordCase
	KeywordCatch
	KeywordClass
	KeywordConst
	KeywordContinue
	KeywordDebugger
	KeywordDefault
	KeywordDelete
	KeywordDo
	KeywordElse
	KeywordExport
	KeywordExtends
	KeywordFinally
	KeywordFor
	KeywordFunction
	KeywordIf
	KeywordImport
	KeywordIn
	KeywordInstanceof
	KeywordNew
	KeywordReturn
	KeywordSuper
	KeywordSwitch
	KeywordThis
	KeywordThrow
	KeywordTry
	KeywordTypeof
	KeywordVar
	KeywordVoid
	KeywordWhile
	KeywordWith
	KeywordYield
)

var keywords = map[string]Keyword{
	"break":      KeywordBreak,
	"case":       KeywordCase,
	"catch":      KeywordCatch,
	"class":      KeywordClass,
	"const":      KeywordConst,
	"continue":   KeywordContinue,
	"debugger":   KeywordDebugger,
	"default":    KeywordDefault,
	"delete":     KeywordDelete,
	"do":         KeywordDo,
	"else":       KeywordElse,
	"export":     KeywordExport,
	"extends":    KeywordExtends,
	"finally":    KeywordFinally,
	"for":        KeywordFor,
	"function":   KeywordFunction,
	"if":         KeywordIf,
	"import":     KeywordImport,
	"in":         KeywordIn,
	"instanceof": KeywordInstanceof,
	"new":        KeywordNew,
	"return":     KeywordReturn,
	"super":      KeywordSuper,
	"switch":     KeywordSwitch,
	"this":       KeywordThis,
	"throw":      KeywordThrow,
	"try":        KeywordTry,
	"typeof":     KeywordTypeof,
	"var":        KeywordVar,
	"void":       KeywordVoid,
	"while":      KeywordWhile,
	"with":       KeywordWith,
	"yield":      KeywordYield,
}

func ParseString(s string) (TokenSet, error) {
	return Parse(strings.NewReader(s))
}

func ParseBytes(b []byte) (TokenSet, error) {
	return Parse(bytes.NewReader(b))
}

func Parse(rd io.Reader) (TokenSet, error) {
	t := NewTokeniser(rd)

	var state stateList

	var a TokenSet
	for {
		tk, err := t.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		switch tk.Kind {
		case TokenKindTemplateHead, TokenKindTemplateMiddle:
			state.pushMode(t, InputElementRegExpOrTemplateTail)
		case TokenKindTemplateTail:
			state.popMode(t, InputElementDiv)
		case TokenKindPuncLeftBrace:
			state.pushMode(t, InputElementRegExp)
		case TokenKindPuncRightBrace:
			state.popMode(t, InputElementRegExp)
		}

		if tk.Kind == TokenKindIdentifier {
			if _, ok := keywords[tk.Value]; ok {
				tk = &Token{
					Kind:   TokenKindKeyword,
					Value:  tk.Value,
					Raw:    tk.Raw,
					Offset: tk.Offset,
				}
			}
		}

		a = append(a, *tk)
	}

	return a, nil
}
