package ast

import (
	"encoding/json"
)

type UnaryOperator string

const (
	UnaryOperatorMinus  = "-"
	UnaryOperatorPlus   = "+"
	UnaryOperatorBang   = "!"
	UnaryOperatorTilde  = "~"
	UnaryOperatorTypeof = "typeof"
	UnaryOperatorVoid   = "void"
	UnaryOperatorDelete = "delete"
)

func (v UnaryOperator) Valid() bool {
	return v == UnaryOperatorMinus || v == UnaryOperatorPlus || v == UnaryOperatorBang || v == UnaryOperatorTilde || v == UnaryOperatorTypeof || v == UnaryOperatorVoid || v == UnaryOperatorDelete
}

type UpdateOperator string

const (
	UpdateOperatorIncrement = "++"
	UpdateOperatorDecrement = "--"
)

func (v UpdateOperator) Valid() bool {
	return v == UpdateOperatorIncrement || v == UpdateOperatorDecrement
}

type BinaryOperator string

const (
	BinaryOperatorEqual            = "=="
	BinaryOperatorNotEqual         = "!="
	BinaryOperatorStrictEqual      = "==="
	BinaryOperatorStrictNotEqual   = "!=="
	BinaryOperatorLess             = "<"
	BinaryOperatorLessOrEqual      = "<="
	BinaryOperatorGreater          = ">"
	BinaryOperatorGreaterOrEqual   = ">="
	BinaryOperatorShiftLeft        = "<<"
	BinaryOperatorShiftRight       = ">>"
	BinaryOperatorShiftRightSigned = ">>>"
	BinaryOperatorPlus             = "+"
	BinaryOperatorMinus            = "-"
	BinaryOperatorMultiply         = "*"
	BinaryOperatorDivide           = "/"
	BinaryOperatorModulo           = "%"
	BinaryOperatorOr               = "|"
	BinaryOperatorXor              = "^"
	BinaryOperatorAnd              = "&"
	BinaryOperatorIn               = "in"
	BinaryOperatorInstanceof       = "instanceof"
)

func (v BinaryOperator) Valid() bool {
	return v == BinaryOperatorEqual || v == BinaryOperatorNotEqual || v == BinaryOperatorStrictEqual || v == BinaryOperatorStrictNotEqual || v == BinaryOperatorLess || v == BinaryOperatorLessOrEqual || v == BinaryOperatorGreater || v == BinaryOperatorGreaterOrEqual || v == BinaryOperatorShiftLeft || v == BinaryOperatorShiftRight || v == BinaryOperatorShiftRightSigned || v == BinaryOperatorPlus || v == BinaryOperatorMinus || v == BinaryOperatorMultiply || v == BinaryOperatorDivide || v == BinaryOperatorModulo || v == BinaryOperatorOr || v == BinaryOperatorXor || v == BinaryOperatorAnd || v == BinaryOperatorIn || v == BinaryOperatorInstanceof
}

type AssignmentOperator string

const (
	AssignmentOperatorEquals           = "="
	AssignmentOperatorAdd              = "+="
	AssignmentOperatorSubtract         = "-="
	AssignmentOperatorMultiply         = "*="
	AssignmentOperatorDivide           = "/="
	AssignmentOperatorModulo           = "%="
	AssignmentOperatorShiftLeft        = "<<="
	AssignmentOperatorShiftRight       = ">>="
	AssignmentOperatorShiftRightSigned = ">>>="
	AssignmentOperatorOr               = "|="
	AssignmentOperatorXor              = "^="
	AssignmentOperatorAnd              = "&="
)

func (v AssignmentOperator) Valid() bool {
	return v == AssignmentOperatorEquals || v == AssignmentOperatorAdd || v == AssignmentOperatorSubtract || v == AssignmentOperatorMultiply || v == AssignmentOperatorDivide || v == AssignmentOperatorModulo || v == AssignmentOperatorShiftLeft || v == AssignmentOperatorShiftRight || v == AssignmentOperatorShiftRightSigned || v == AssignmentOperatorOr || v == AssignmentOperatorXor || v == AssignmentOperatorAnd
}

type LogicalOperator string

const (
	LogicalOperatorOr  = "||"
	LogicalOperatorAnd = "&&"
)

func (v LogicalOperator) Valid() bool { return v == LogicalOperatorOr || v == LogicalOperatorAnd }

type Node struct {
	Type string          `json:"type"`
	Loc  *SourceLocation `json:"loc"`
}

type SourceLocation struct {
	Source *string  `json:"source"`
	Start  Position `json:"start"`
	End    Position `json:"end"`
}

type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type Identifier struct {
	Expression
	Pattern
	Name string `json:"name"`
}

type Literal struct {
	Expression
}

type RegExpLiteral struct {
	Literal
	Pattern string `json:"pattern"`
	Flags   string `json:"flags"`
}

type NullLiteral struct {
	Literal
}

type StringLiteral struct {
	Literal
	Value string `json:"value"`
}

type BooleanLiteral struct {
	Literal
	Value bool `json:"value"`
}

type NumericLiteral struct {
	Literal
	Value float64 `json:"value"`
}

type StatementOrModuleDeclaration struct {
	Statement         *Statement
	ModuleDeclaration *ModuleDeclaration
}

func (u StatementOrModuleDeclaration) MarshalJSON() ([]byte, error) {
	switch {
	case u.Statement != nil:
		return json.Marshal(u.Statement)
	case u.ModuleDeclaration != nil:
		return json.Marshal(u.ModuleDeclaration)
	default:
		return []byte("null"), nil
	}
}

type Program struct {
	Node
	SourceType string                       `json:"sourceType"`
	Body       StatementOrModuleDeclaration `json:"body"`
	Directives []Directive                  `json:"directives"`
}

type Function struct {
	Node
	Id        *Identifier    `json:"id"`
	Params    []Pattern      `json:"params"`
	Body      BlockStatement `json:"body"`
	Generator bool           `json:"generator"`
	Async     bool           `json:"async"`
}

type Statement struct {
	Node
}

type ExpressionStatement struct {
	Statement
	Expression Expression `json:"expression"`
}

type BlockStatement struct {
	Statement
	Body       []Statement `json:"body"`
	Directives []Directive `json:"directives"`
}

type EmptyStatement struct {
	Statement
}

type DebuggerStatement struct {
	Statement
}

type WithStatement struct {
	Statement
	Object Expression `json:"object"`
	Body   Statement  `json:"body"`
}

type ReturnStatement struct {
	Statement
	Argument *Expression `json:"argument"`
}

type LabeledStatement struct {
	Statement
	Label Identifier `json:"label"`
	Body  Statement  `json:"body"`
}

type BreakStatement struct {
	Statement
	Label *Identifier `json:"label"`
}

type ContinueStatement struct {
	Statement
	Label *Identifier `json:"label"`
}

type IfStatement struct {
	Statement
	Test       Expression `json:"test"`
	Consequent Statement  `json:"consequent"`
	Alternate  *Statement `json:"alternate"`
}

type SwitchStatement struct {
	Statement
	Discriminant Expression   `json:"discriminant"`
	Cases        []SwitchCase `json:"cases"`
}

type SwitchCase struct {
	Node
	Test       *Expression `json:"test"`
	Consequent []Statement `json:"consequent"`
}

type ThrowStatement struct {
	Statement
	Argument Expression `json:"argument"`
}

type TryStatement struct {
	Statement
	Block     BlockStatement  `json:"block"`
	Handler   *CatchClause    `json:"handler"`
	Finalizer *BlockStatement `json:"finalizer"`
}

type CatchClause struct {
	Node
	Param Pattern        `json:"param"`
	Body  BlockStatement `json:"body"`
}

type WhileStatement struct {
	Statement
	Test Expression `json:"test"`
	Body Statement  `json:"body"`
}

type DoWhileStatement struct {
	Statement
	Body Statement  `json:"body"`
	Test Expression `json:"test"`
}

type VariableDeclarationOrExpression struct {
	VariableDeclaration *VariableDeclaration
	Expression          *Expression
}

func (u VariableDeclarationOrExpression) MarshalJSON() ([]byte, error) {
	switch {
	case u.VariableDeclaration != nil:
		return json.Marshal(u.VariableDeclaration)
	case u.Expression != nil:
		return json.Marshal(u.Expression)
	default:
		return []byte("null"), nil
	}
}

type ForStatement struct {
	Statement
	Init   VariableDeclarationOrExpression `json:"init"`
	Test   *Expression                     `json:"test"`
	Update *Expression                     `json:"update"`
	Body   Statement                       `json:"body"`
}

type ForInStatement struct {
	Statement
	Left  VariableDeclarationOrExpression `json:"left"`
	Right Expression                      `json:"right"`
	Body  Statement                       `json:"body"`
}

type ForOfStatement struct {
	ForInStatement
}

type Declaration struct {
	Statement
}

type FunctionDeclaration struct {
	Function
	Declaration
	Id Identifier `json:"id"`
}

type VariableDeclaration struct {
	Declaration
	Declarations []VariableDeclarator `json:"declarations"`
	Kind         string               `json:"kind"`
}

type VariableDeclarator struct {
	Node
	Id   Pattern     `json:"id"`
	Init *Expression `json:"init"`
}

type Decorator struct {
	Node
	Expression Expression `json:"expression"`
}

type Directive struct {
	Node
	Value DirectiveLiteral `json:"value"`
}

type DirectiveLiteral struct {
	StringLiteral
}

type Expression struct {
	Node
}

type Super struct {
	Node
}

type ThisExpression struct {
	Expression
}

type BlockStatementOrExpression struct {
	BlockStatement *BlockStatement
	Expression     *Expression
}

func (u BlockStatementOrExpression) MarshalJSON() ([]byte, error) {
	switch {
	case u.BlockStatement != nil:
		return json.Marshal(u.BlockStatement)
	case u.Expression != nil:
		return json.Marshal(u.Expression)
	default:
		return []byte("null"), nil
	}
}

type ArrowFunctionExpression struct {
	Function
	Expression
	Body         BlockStatementOrExpression `json:"body"`
	IsExpression bool                       `json:"expression"`
}

type YieldExpression struct {
	Expression
	Argument *Expression `json:"argument"`
	Delegate bool        `json:"delegate"`
}

type AwaitExpression struct {
	Expression
	Argument *Expression `json:"argument"`
}

type ExpressionOrSpreadElement struct {
	Expression    *Expression
	SpreadElement *SpreadElement
}

func (u ExpressionOrSpreadElement) MarshalJSON() ([]byte, error) {
	switch {
	case u.Expression != nil:
		return json.Marshal(u.Expression)
	case u.SpreadElement != nil:
		return json.Marshal(u.SpreadElement)
	default:
		return []byte("null"), nil
	}
}

type ArrayExpression struct {
	Expression
	Elements ExpressionOrSpreadElement `json:"elements"`
}

type ObjectPropertyOrObjectMethodOrSpreadProperty struct {
	ObjectProperty *ObjectProperty
	ObjectMethod   *ObjectMethod
	SpreadProperty *SpreadProperty
}

func (u ObjectPropertyOrObjectMethodOrSpreadProperty) MarshalJSON() ([]byte, error) {
	switch {
	case u.ObjectProperty != nil:
		return json.Marshal(u.ObjectProperty)
	case u.ObjectMethod != nil:
		return json.Marshal(u.ObjectMethod)
	case u.SpreadProperty != nil:
		return json.Marshal(u.SpreadProperty)
	default:
		return []byte("null"), nil
	}
}

type ObjectExpression struct {
	Expression
	Properties ObjectPropertyOrObjectMethodOrSpreadProperty `json:"properties"`
}

type ObjectMember struct {
	Node
	Key        Expression  `json:"key"`
	Computed   bool        `json:"computed"`
	Value      Expression  `json:"value"`
	Decorators []Decorator `json:"decorators"`
}

type ObjectProperty struct {
	ObjectMember
	Shorthand bool `json:"shorthand"`
}

type ObjectMethod struct {
	ObjectMember
	Function
	Kind string `json:"kind"`
}

type RestProperty struct {
	Node
	Argument Expression `json:"argument"`
}

type SpreadProperty struct {
	Node
	Argument Expression `json:"argument"`
}

type FunctionExpression struct {
	Function
	Expression
}

type UnaryExpression struct {
	Expression
	Operator UnaryOperator `json:"operator"`
	Prefix   bool          `json:"prefix"`
	Argument Expression    `json:"argument"`
}

type UpdateExpression struct {
	Expression
	Operator UpdateOperator `json:"operator"`
	Argument Expression     `json:"argument"`
	Prefix   bool           `json:"prefix"`
}

type BinaryExpression struct {
	Expression
	Operator BinaryOperator `json:"operator"`
	Left     Expression     `json:"left"`
	Right    Expression     `json:"right"`
}

type PatternOrExpression struct {
	Pattern    *Pattern
	Expression *Expression
}

func (u PatternOrExpression) MarshalJSON() ([]byte, error) {
	switch {
	case u.Pattern != nil:
		return json.Marshal(u.Pattern)
	case u.Expression != nil:
		return json.Marshal(u.Expression)
	default:
		return []byte("null"), nil
	}
}

type AssignmentExpression struct {
	Expression
	Operator AssignmentOperator  `json:"operator"`
	Left     PatternOrExpression `json:"left"`
	Right    Expression          `json:"right"`
}

type LogicalExpression struct {
	Expression
	Operator LogicalOperator `json:"operator"`
	Left     Expression      `json:"left"`
	Right    Expression      `json:"right"`
}

type SpreadElement struct {
	Node
	Argument Expression `json:"argument"`
}

type ExpressionOrSuper struct {
	Expression *Expression
	Super      *Super
}

func (u ExpressionOrSuper) MarshalJSON() ([]byte, error) {
	switch {
	case u.Expression != nil:
		return json.Marshal(u.Expression)
	case u.Super != nil:
		return json.Marshal(u.Super)
	default:
		return []byte("null"), nil
	}
}

type MemberExpression struct {
	Expression
	Pattern
	Object   ExpressionOrSuper `json:"object"`
	Property Expression        `json:"property"`
	Computed bool              `json:"computed"`
}

type BindExpression struct {
	Expression
	Object []*Expression `json:"object"`
	Callee []Expression  `json:"callee"`
}

type ConditionalExpression struct {
	Expression
	Test       Expression `json:"test"`
	Alternate  Expression `json:"alternate"`
	Consequent Expression `json:"consequent"`
}

type CallExpression struct {
	Expression
	Callee    ExpressionOrSuper         `json:"callee"`
	Arguments ExpressionOrSpreadElement `json:"arguments"`
}

type NewExpression struct {
	CallExpression
}

type SequenceExpression struct {
	Expression
	Expressions []Expression `json:"expressions"`
}

type TemplateLiteral struct {
	Expression
	Quasis      []TemplateElement `json:"quasis"`
	Expressions []Expression      `json:"expressions"`
}

type TaggedTemplateExpression struct {
	Expression
	Tag   Expression      `json:"tag"`
	Quasi TemplateLiteral `json:"quasi"`
}

type TemplateElement struct {
	Node
	Tail   bool   `json:"tail"`
	Cooked string `json:"cooked"`
	Raw    string `json:"raw"`
}

type Pattern struct {
	Node
}

type AssignmentProperty struct {
	ObjectProperty
	Value Pattern `json:"value"`
}

type AssignmentPropertyOrRestProperty struct {
	AssignmentProperty *AssignmentProperty
	RestProperty       *RestProperty
}

func (u AssignmentPropertyOrRestProperty) MarshalJSON() ([]byte, error) {
	switch {
	case u.AssignmentProperty != nil:
		return json.Marshal(u.AssignmentProperty)
	case u.RestProperty != nil:
		return json.Marshal(u.RestProperty)
	default:
		return []byte("null"), nil
	}
}

type ObjectPattern struct {
	Pattern
	Properties AssignmentPropertyOrRestProperty `json:"properties"`
}

type ArrayPattern struct {
	Pattern
	Elements []*Pattern `json:"elements"`
}

type RestElement struct {
	Pattern
	Argument Pattern `json:"argument"`
}

type AssignmentPattern struct {
	Pattern
	Left  Pattern    `json:"left"`
	Right Expression `json:"right"`
}

type Class struct {
	Node
	Id         *Identifier `json:"id"`
	SuperClass *Expression `json:"superClass"`
	Body       ClassBody   `json:"body"`
	Decorators []Decorator `json:"decorators"`
}

type ClassMethodOrClassProperty struct {
	ClassMethod   *ClassMethod
	ClassProperty *ClassProperty
}

func (u ClassMethodOrClassProperty) MarshalJSON() ([]byte, error) {
	switch {
	case u.ClassMethod != nil:
		return json.Marshal(u.ClassMethod)
	case u.ClassProperty != nil:
		return json.Marshal(u.ClassProperty)
	default:
		return []byte("null"), nil
	}
}

type ClassBody struct {
	Node
	Body ClassMethodOrClassProperty `json:"body"`
}

type ClassMethod struct {
	Node
	Key        Expression         `json:"key"`
	Value      FunctionExpression `json:"value"`
	Kind       string             `json:"kind"`
	Computed   bool               `json:"computed"`
	Static     bool               `json:"static"`
	Decorators []Decorator        `json:"decorators"`
}

type ClassProperty struct {
	Node
	Key   Identifier `json:"key"`
	Value Expression `json:"value"`
}

type ClassDeclaration struct {
	Class
	Declaration
	Id Identifier `json:"id"`
}

type ClassExpression struct {
	Class
	Expression
}

type MetaProperty struct {
	Expression
	Meta     Identifier `json:"meta"`
	Property Identifier `json:"property"`
}

type ModuleDeclaration struct {
	Node
}

type ModuleSpecifier struct {
	Node
	Local Identifier `json:"local"`
}

type ImportSpecifierOrImportDefaultSpecifierOrImportNamespaceSpecifier struct {
	ImportSpecifier          *ImportSpecifier
	ImportDefaultSpecifier   *ImportDefaultSpecifier
	ImportNamespaceSpecifier *ImportNamespaceSpecifier
}

func (u ImportSpecifierOrImportDefaultSpecifierOrImportNamespaceSpecifier) MarshalJSON() ([]byte, error) {
	switch {
	case u.ImportSpecifier != nil:
		return json.Marshal(u.ImportSpecifier)
	case u.ImportDefaultSpecifier != nil:
		return json.Marshal(u.ImportDefaultSpecifier)
	case u.ImportNamespaceSpecifier != nil:
		return json.Marshal(u.ImportNamespaceSpecifier)
	default:
		return []byte("null"), nil
	}
}

type ImportDeclaration struct {
	ModuleDeclaration
	Specifiers ImportSpecifierOrImportDefaultSpecifierOrImportNamespaceSpecifier `json:"specifiers"`
	Source     Literal                                                           `json:"source"`
}

type ImportSpecifier struct {
	ModuleSpecifier
	Imported Identifier `json:"imported"`
}

type ImportDefaultSpecifier struct {
	ModuleSpecifier
}

type ImportNamespaceSpecifier struct {
	ModuleSpecifier
}

type ExportNamedDeclaration struct {
	ModuleDeclaration
	Declaration *Declaration      `json:"declaration"`
	Specifiers  []ExportSpecifier `json:"specifiers"`
	Source      *Literal          `json:"source"`
}

type ExportSpecifier struct {
	ModuleSpecifier
	Exported Identifier `json:"exported"`
}

type DeclarationOrExpression struct {
	Declaration *Declaration
	Expression  *Expression
}

func (u DeclarationOrExpression) MarshalJSON() ([]byte, error) {
	switch {
	case u.Declaration != nil:
		return json.Marshal(u.Declaration)
	case u.Expression != nil:
		return json.Marshal(u.Expression)
	default:
		return []byte("null"), nil
	}
}

type ExportDefaultDeclaration struct {
	ModuleDeclaration
	Declaration DeclarationOrExpression `json:"declaration"`
}

type ExportAllDeclaration struct {
	ModuleDeclaration
	Source Literal `json:"source"`
}
