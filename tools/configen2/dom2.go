package configen2

type Dom2root struct {
	PackageName string
	Imports     map[string]string         // map[ path ] alias
	Components  map[string]*Dom2component // map[com_name] com
}

type Dom2component struct {
	StructName   string
	Attributes   map[string]string
	InjectionMap map[string]*Dom2injection

	ID      string
	Class   string
	Scope   string
	Aliases string
	ComType string

	// do://inject@injection
	InjectionMainMethod string

	// do://inject@instance (optional)
	InjectMethod string

	InitMethod    string
	DestroyMethod string
}

type Dom2injection struct {
	FieldName string
	FieldType string
	FieldTag  string

	Auto     bool
	Selector string

	// do://inject@injection
	InjectionGetterMethod string

	Attributes map[string]string
}
