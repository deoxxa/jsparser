type UnaryOperator string

const (
	UnaryOperator_0 = "-"
	UnaryOperator_1 = "+"
	UnaryOperator_2 = "!"
	UnaryOperator_3 = "~"
	UnaryOperator_4 = "typeof"
	UnaryOperator_5 = "void"
	UnaryOperator_6 = "delete"
)

func (v UnaryOperator) Valid() bool {
	return v == UnaryOperator_0 || v == UnaryOperator_1 || v == UnaryOperator_2 || v == UnaryOperator_3 || v == UnaryOperator_4 || v == UnaryOperator_5 || v == UnaryOperator_6
}

type UpdateOperator string

const (
	UpdateOperator_0 = "++"
	UpdateOperator_1 = "--"
)

func (v UpdateOperator) Valid() bool { return v == UpdateOperator_0 || v == UpdateOperator_1 }

type BinaryOperator string

const (
	BinaryOperator_0  = "=="
	BinaryOperator_1  = "!="
	BinaryOperator_2  = "==="
	BinaryOperator_3  = "!=="
	BinaryOperator_4  = "<"
	BinaryOperator_5  = "<="
	BinaryOperator_6  = ">"
	BinaryOperator_7  = ">="
	BinaryOperator_8  = "<<"
	BinaryOperator_9  = ">>"
	BinaryOperator_10 = ">>>"
	BinaryOperator_11 = "+"
	BinaryOperator_12 = "-"
	BinaryOperator_13 = "*"
	BinaryOperator_14 = "/"
	BinaryOperator_15 = "%"
	BinaryOperator_16 = "|"
	BinaryOperator_17 = "^"
	BinaryOperator_18 = "&"
	BinaryOperator_19 = "in"
	BinaryOperator_20 = "instanceof"
)

func (v BinaryOperator) Valid() bool {
	return v == BinaryOperator_0 || v == BinaryOperator_1 || v == BinaryOperator_2 || v == BinaryOperator_3 || v == BinaryOperator_4 || v == BinaryOperator_5 || v == BinaryOperator_6 || v == BinaryOperator_7 || v == BinaryOperator_8 || v == BinaryOperator_9 || v == BinaryOperator_10 || v == BinaryOperator_11 || v == BinaryOperator_12 || v == BinaryOperator_13 || v == BinaryOperator_14 || v == BinaryOperator_15 || v == BinaryOperator_16 || v == BinaryOperator_17 || v == BinaryOperator_18 || v == BinaryOperator_19 || v == BinaryOperator_20
}

type AssignmentOperator string

const (
	AssignmentOperator_0  = "="
	AssignmentOperator_1  = "+="
	AssignmentOperator_2  = "-="
	AssignmentOperator_3  = "*="
	AssignmentOperator_4  = "/="
	AssignmentOperator_5  = "%="
	AssignmentOperator_6  = "<<="
	AssignmentOperator_7  = ">>="
	AssignmentOperator_8  = ">>>="
	AssignmentOperator_9  = "|="
	AssignmentOperator_10 = "^="
	AssignmentOperator_11 = "&="
)

func (v AssignmentOperator) Valid() bool {
	return v == AssignmentOperator_0 || v == AssignmentOperator_1 || v == AssignmentOperator_2 || v == AssignmentOperator_3 || v == AssignmentOperator_4 || v == AssignmentOperator_5 || v == AssignmentOperator_6 || v == AssignmentOperator_7 || v == AssignmentOperator_8 || v == AssignmentOperator_9 || v == AssignmentOperator_10 || v == AssignmentOperator_11
}

type LogicalOperator string

const (
	LogicalOperator_0 = "||"
	LogicalOperator_1 = "&&"
)

func (v LogicalOperator) Valid() bool { return v == LogicalOperator_0 || v == LogicalOperator_1 }

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
	Line   number `json:"line"`
	Column number `json:"column"`
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
	Value boolean `json:"value"`
}

type NumericLiteral struct {
	Literal
	Value number `json:"value"`
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
	Generator boolean        `json:"generator"`
	Async     boolean        `json:"async"`
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
	Body       BlockStatementOrExpression `json:"body"`
	Expression boolean                    `json:"expression"`
}

type YieldExpression struct {
	Expression
	Argument *Expression `json:"argument"`
	Delegate boolean     `json:"delegate"`
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
	Computed   boolean     `json:"computed"`
	Value      Expression  `json:"value"`
	Decorators []Decorator `json:"decorators"`
}

type ObjectProperty struct {
	ObjectMember
	Shorthand boolean `json:"shorthand"`
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
	Prefix   boolean       `json:"prefix"`
	Argument Expression    `json:"argument"`
}

type UpdateExpression struct {
	Expression
	Operator UpdateOperator `json:"operator"`
	Argument Expression     `json:"argument"`
	Prefix   boolean        `json:"prefix"`
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
	Computed boolean           `json:"computed"`
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
	Tail   boolean `json:"tail"`
	Cooked string  `json:"cooked"`
	Raw    string  `json:"raw"`
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
	Computed   boolean            `json:"computed"`
	Static     boolean            `json:"static"`
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

