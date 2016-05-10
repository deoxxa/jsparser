package jsparser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokeniser(t *testing.T) {
	a := assert.New(t)

	o := NewTokeniser(strings.NewReader("var what = 'test'; console.log(`this is a ${/* yep */what/* nope */}`);"))

	r, err := o.ReadAll()
	a.NoError(err)
	a.Equal(TokenSet{
		Token{Kind: TokenKindIdentifier, Value: "var", Raw: "var", Offset: 0},
		Token{Kind: TokenKindWhitespace, Raw: " ", Offset: 3},
		Token{Kind: TokenKindIdentifier, Value: "what", Raw: "what", Offset: 4},
		Token{Kind: TokenKindWhitespace, Raw: " ", Offset: 8},
		Token{Kind: TokenKindBinaryAssignment, Raw: "=", Offset: 9},
		Token{Kind: TokenKindWhitespace, Raw: " ", Offset: 10},
		Token{Kind: TokenKindString, Value: "test", Raw: "'test'", Offset: 11},
		Token{Kind: TokenKindPuncSemicolon, Raw: ";", Offset: 17},
		Token{Kind: TokenKindWhitespace, Raw: " ", Offset: 18},
		Token{Kind: TokenKindIdentifier, Value: "console", Raw: "console", Offset: 19},
		Token{Kind: TokenKindPuncPeriod, Raw: ".", Offset: 26},
		Token{Kind: TokenKindIdentifier, Value: "log", Raw: "log", Offset: 27},
		Token{Kind: TokenKindPuncLeftParen, Raw: "(", Offset: 30},
		Token{Kind: TokenKindPuncBacktick, Raw: "`", Offset: 31},
		Token{Kind: TokenKindTemplateText, Value: "this is a ", Raw: "this is a ", Offset: 32},
		Token{Kind: TokenKindPuncTemplateVariableOpen, Raw: "${", Offset: 42},
		Token{Kind: TokenKindMultipleLineComment, Value: "/* yep */", Raw: "/* yep */", Offset: 44},
		Token{Kind: TokenKindIdentifier, Value: "what", Raw: "what", Offset: 53},
		Token{Kind: TokenKindMultipleLineComment, Value: "/* nope */", Raw: "/* nope */", Offset: 57},
		Token{Kind: TokenKindPuncTemplateVariableClose, Raw: "}", Offset: 67},
		Token{Kind: TokenKindPuncBacktick, Raw: "`", Offset: 68},
		Token{Kind: TokenKindPuncRightParen, Raw: ")", Offset: 69},
		Token{Kind: TokenKindPuncSemicolon, Raw: ";", Offset: 70},
	}, r)
}
