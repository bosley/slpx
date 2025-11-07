package object

import "fmt"

/*
These are defined in object as it makes sense existentially to find them
here, but there is a "cross concern" given that its formalizing an aspect
of a representation of a concept in object simply because the implementation
of the concept is here, but its where I draw the line. I don't anticipate
using type names in any other way than with the object of the thing of which
that type indicator represents

they are set as "var" so we can, in practice, change the very central-and-important
symbols but that would make invalid every piece of code that leverages them

depending on the context however, it may be useful, say for tests, to have all of
manually tweakable to induce strange behaviors
*/
var (
	SYMBOL_ObjType_None       = ":_"
	SYMBOL_ObjType_Some       = ":Q"
	SYMBOL_ObjType_Any        = ":*"
	SYMBOL_ObjType_List       = ":L"
	SYMBOL_ObjType_Error      = ":E"
	SYMBOL_ObjType_String     = ":S"
	SYMBOL_ObjType_Integer    = ":I"
	SYMBOL_ObjType_Real       = ":R"
	SYMBOL_ObjType_Identifier = ":X"
	SYMBOL_ObjType_Function   = ":F"
)

func GetTypeFromIdentifier(target Identifier) (ObjType, error) {
	switch string(target) {
	case SYMBOL_ObjType_None:
		return OBJ_TYPE_NONE, nil
	case SYMBOL_ObjType_Some:
		return OBJ_TYPE_SOME, nil
	case SYMBOL_ObjType_Any:
		return OBJ_TYPE_ANY, nil
	case SYMBOL_ObjType_List:
		return OBJ_TYPE_LIST, nil
	case SYMBOL_ObjType_Error:
		return OBJ_TYPE_ERROR, nil
	case SYMBOL_ObjType_String:
		return OBJ_TYPE_STRING, nil
	case SYMBOL_ObjType_Integer:
		return OBJ_TYPE_INTEGER, nil
	case SYMBOL_ObjType_Real:
		return OBJ_TYPE_REAL, nil
	case SYMBOL_ObjType_Identifier:
		return OBJ_TYPE_IDENTIFIER, nil
	case SYMBOL_ObjType_Function:
		return OBJ_TYPE_FUNCTION, nil
	default:
		return "", fmt.Errorf("invalid type identifier: %s", target)
	}
}

func GetIdentifierFromType(target ObjType) Identifier {
	switch target {
	case OBJ_TYPE_NONE:
		return Identifier(SYMBOL_ObjType_None)
	case OBJ_TYPE_SOME:
		return Identifier(SYMBOL_ObjType_Some)
	case OBJ_TYPE_ANY:
		return Identifier(SYMBOL_ObjType_Any)
	case OBJ_TYPE_LIST:
		return Identifier(SYMBOL_ObjType_List)
	case OBJ_TYPE_ERROR:
		return Identifier(SYMBOL_ObjType_Error)
	case OBJ_TYPE_STRING:
		return Identifier(SYMBOL_ObjType_String)
	case OBJ_TYPE_INTEGER:
		return Identifier(SYMBOL_ObjType_Integer)
	case OBJ_TYPE_REAL:
		return Identifier(SYMBOL_ObjType_Real)
	case OBJ_TYPE_IDENTIFIER:
		return Identifier(SYMBOL_ObjType_Identifier)
	case OBJ_TYPE_FUNCTION:
		return Identifier(SYMBOL_ObjType_Function)
	default:
		return Identifier(SYMBOL_ObjType_None)
	}
}
