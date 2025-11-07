package reflection

import (
	"github.com/bosley/slpx/pkg/slp/env"
	"github.com/bosley/slpx/pkg/slp/object"
)

type reflectionFunctions struct{}

func NewReflectionFunctions() env.FunctionGroup {
	return &reflectionFunctions{}
}

func (r *reflectionFunctions) Name() string {
	return "reflect"
}

func (r *reflectionFunctions) Functions() map[object.Identifier]env.EnvFunction {
	return map[object.Identifier]env.EnvFunction{
		"reflect/type?": {
			EvaluateArgs: false,
			Parameters: []env.EnvParameter{
				{Name: "value", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_STRING,
			Body:       cmdReflectType,
		},
		"reflect/equal?": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "a", Type: object.OBJ_TYPE_ANY},
				{Name: "b", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdReflectEqual,
		},
		"reflect/int?": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "value", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdReflectIsInt,
		},
		"reflect/real?": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "value", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdReflectIsReal,
		},
		"reflect/str?": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "value", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdReflectIsStr,
		},
		"reflect/list?": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "value", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdReflectIsList,
		},
		"reflect/fn?": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "value", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdReflectIsFn,
		},
		"reflect/none?": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "value", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdReflectIsNone,
		},
		"reflect/error?": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "value", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdReflectIsError,
		},
		"reflect/some?": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "value", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdReflectIsSome,
		},
		"reflect/ident?": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "value", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdReflectIsIdent,
		},
	}
}

func cmdReflectType(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	value := args[0]

	if value.Type == object.OBJ_TYPE_IDENTIFIER {
		resolved, err := ctx.Evaluate(value)
		if err != nil {
			return object.Obj{Type: object.OBJ_TYPE_STRING, D: string(object.OBJ_TYPE_IDENTIFIER)}, nil
		}
		if resolved.Type == object.OBJ_TYPE_ERROR {
			return object.Obj{Type: object.OBJ_TYPE_STRING, D: string(object.OBJ_TYPE_IDENTIFIER)}, nil
		}
		value = resolved
	}

	return object.Obj{Type: object.OBJ_TYPE_STRING, D: string(value.Type)}, nil
}

func cmdReflectEqual(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	a := args[0]
	b := args[1]
	if a.Type == b.Type {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdReflectIsInt(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	value := args[0]
	if value.Type == object.OBJ_TYPE_INTEGER {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdReflectIsReal(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	value := args[0]
	if value.Type == object.OBJ_TYPE_REAL {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdReflectIsStr(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	value := args[0]
	if value.Type == object.OBJ_TYPE_STRING {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdReflectIsList(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	value := args[0]
	if value.Type == object.OBJ_TYPE_LIST {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdReflectIsFn(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	value := args[0]
	if value.Type == object.OBJ_TYPE_FUNCTION {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdReflectIsNone(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	value := args[0]
	if value.Type == object.OBJ_TYPE_NONE {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdReflectIsError(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	value := args[0]
	if value.Type == object.OBJ_TYPE_ERROR {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdReflectIsSome(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	value := args[0]
	if value.Type == object.OBJ_TYPE_SOME {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdReflectIsIdent(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	value := args[0]
	if value.Type == object.OBJ_TYPE_IDENTIFIER {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}
