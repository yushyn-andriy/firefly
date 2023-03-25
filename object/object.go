package object

type ObjectType string

const (
	CLASS            = "CLASS"
	INSTANCE         = "INSTANCE"
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
	TYPE_OBJ         = "TYPE"
	FORLOOP_OBJ      = "FORLOOP"
	NULL_OBJ         = "NULL"
)

// magic methods
const (
	MAGIC_METHOD_INIT = "__init__"
	MAGIC_METHOD_LEN  = "__len__"
	MAGIC_METHOD_REPR = "__repr__"
	MAGIC_METHOD_STR  = "__str__"
)

// magic attributes
const (
	MAGIC_ATTR_NAME = "__name__"
	MAGIC_ATTR_DOC  = "__doc__"
)

var (
	NULL  = &Null{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

type Object interface {
	Type() ObjectType
	Inspect() string
	SetAttr(key string, value Object) Object
	GetAttr(key string) Object
}

type Hashable interface {
	HashKey() HashKey
}
