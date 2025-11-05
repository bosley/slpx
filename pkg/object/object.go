package object

import (
	"fmt"
)

type ObjType string

const (
	OBJ_TYPE_NONE       ObjType = "none"
	OBJ_TYPE_SOME       ObjType = "some"
	OBJ_TYPE_ANY        ObjType = "any"
	OBJ_TYPE_LIST       ObjType = "list"
	OBJ_TYPE_ERROR      ObjType = "error"
	OBJ_TYPE_STRING     ObjType = "string"
	OBJ_TYPE_INTEGER    ObjType = "integer"
	OBJ_TYPE_REAL       ObjType = "real"
	OBJ_TYPE_IDENTIFIER ObjType = "identifier"
	OBJ_TYPE_FUNCTION   ObjType = "function"
)

type List []Obj
type Some = Obj
type None struct{}
type Error struct {
	Position int
	Message  string
}
type Integer int64
type Real float64
type Identifier string

type Parameter struct {
	Name Identifier
	Type ObjType
}

type Function struct {
	Parameters []Parameter
	ReturnType ObjType
	Variadic   bool
	Body       List
	Self       Obj
}

type Obj struct {
	Type ObjType
	D    any

	// optional data storage for complex types built on-top of
	// the Obj. It is recommended that complex types be stored in a
	// "some" ObjType
	C any
}

func (o Obj) Encode() string {
	switch o.Type {
	case OBJ_TYPE_NONE:
		return "_"
	case OBJ_TYPE_SOME:
		quoted := o.D.(Some)
		return "'" + quoted.Encode()
	case OBJ_TYPE_LIST:
		list := o.D.(List)
		if len(list) == 0 {
			return "()"
		}
		result := "("
		for i, item := range list {
			if i > 0 {
				result += " "
			}
			result += item.Encode()
		}
		result += ")"
		return result
	case OBJ_TYPE_STRING:
		str := o.D.(string)
		return escapeString(str)
	case OBJ_TYPE_INTEGER:
		return fmt.Sprintf("%d", o.D.(Integer))
	case OBJ_TYPE_REAL:
		return fmt.Sprintf("%g", float64(o.D.(Real)))
	case OBJ_TYPE_IDENTIFIER:
		return string(o.D.(Identifier))
	case OBJ_TYPE_ERROR:
		err := o.D.(Error)
		return fmt.Sprintf("ERROR:%d:%s", err.Position, err.Message)
	case OBJ_TYPE_FUNCTION:
		function := o.D.(Function)
		return fmt.Sprintf("FUNCTION:LEN:%d", len(function.Body))
	default:
		return fmt.Sprintf("UNKNOWN_TYPE:%s", o.Type)
	}
}

func (o Obj) DeepCopy() Obj {
	switch o.Type {
	case OBJ_TYPE_LIST:
		originalList := o.D.(List)
		newList := make(List, len(originalList))
		for i, item := range originalList {
			newList[i] = item.DeepCopy()
		}
		return Obj{Type: OBJ_TYPE_LIST, D: newList}
	case OBJ_TYPE_SOME:
		originalSome := o.D.(Some)
		return Obj{Type: OBJ_TYPE_SOME, D: originalSome.DeepCopy()}
	case OBJ_TYPE_NONE:
		return Obj{Type: OBJ_TYPE_NONE, D: None{}}
	case OBJ_TYPE_ERROR:
		originalErr := o.D.(Error)
		return Obj{Type: OBJ_TYPE_ERROR, D: Error{Position: originalErr.Position, Message: originalErr.Message}}
	case OBJ_TYPE_STRING:
		return Obj{Type: OBJ_TYPE_STRING, D: o.D.(string)}
	case OBJ_TYPE_INTEGER:
		return Obj{Type: OBJ_TYPE_INTEGER, D: o.D.(Integer)}
	case OBJ_TYPE_REAL:
		return Obj{Type: OBJ_TYPE_REAL, D: o.D.(Real)}
	case OBJ_TYPE_IDENTIFIER:
		return Obj{Type: OBJ_TYPE_IDENTIFIER, D: o.D.(Identifier)}
	case OBJ_TYPE_FUNCTION:
		originalFunction := o.D.(Function)
		newBody := make(List, len(originalFunction.Body))
		newParameters := make([]Parameter, len(originalFunction.Parameters))
		for i, parameter := range originalFunction.Parameters {
			newParameters[i] = Parameter{Name: parameter.Name, Type: parameter.Type}
		}
		for i, instruction := range originalFunction.Body {
			newBody[i] = instruction.DeepCopy()
		}
		return Obj{Type: OBJ_TYPE_FUNCTION, D: Function{
			Parameters: newParameters,
			ReturnType: originalFunction.ReturnType,
			Variadic:   originalFunction.Variadic,
			Body:       newBody,
			Self:       originalFunction.Self,
		}, C: o.C}
	default:
		return o
	}
}

func escapeString(s string) string {
	result := "\""
	for _, r := range s {
		switch r {
		case '"':
			result += "\\\""
		case '\\':
			result += "\\\\"
		case '\n':
			result += "\\n"
		case '\t':
			result += "\\t"
		case '\r':
			result += "\\r"
		default:
			result += string(r)
		}
	}
	result += "\""
	return result
}
