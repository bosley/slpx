package bits

import (
	"fmt"
	"math"

	"github.com/bosley/slpx/pkg/env"
	"github.com/bosley/slpx/pkg/object"
)

type bitsFunctions struct {
}

func NewBitsFunctions() *bitsFunctions {
	return &bitsFunctions{}
}

func (x *bitsFunctions) Setup(runtime env.Runtime) {
}

func (x *bitsFunctions) Name() string {
	return "bits"
}

func (x *bitsFunctions) Functions() map[object.Identifier]env.EnvFunction {
	return map[object.Identifier]env.EnvFunction{
		"bits/explode": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "value", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_LIST,
			Body:       x.cmdExplode,
		},
		"bits/int": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "bits", Type: object.OBJ_TYPE_LIST},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       x.cmdInt,
		},
		"bits/real": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "bits", Type: object.OBJ_TYPE_LIST},
			},
			ReturnType: object.OBJ_TYPE_REAL,
			Body:       x.cmdReal,
		},
	}
}

func (x *bitsFunctions) cmdExplode(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	value := args[0]

	var bits []object.Obj

	switch value.Type {
	case object.OBJ_TYPE_INTEGER:
		intVal := value.D.(object.Integer)
		uint64Val := uint64(intVal)
		bits = make([]object.Obj, 64)
		for i := 0; i < 64; i++ {
			bit := (uint64Val >> i) & 1
			bits[i] = object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(bit)}
		}

	case object.OBJ_TYPE_REAL:
		realVal := value.D.(object.Real)
		uint64Val := math.Float64bits(float64(realVal))
		bits = make([]object.Obj, 64)
		for i := 0; i < 64; i++ {
			bit := (uint64Val >> i) & 1
			bits[i] = object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(bit)}
		}

	default:
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  fmt.Sprintf("bits/explode: unsupported type %s, expected integer or real", value.Type),
			},
		}, nil
	}

	return object.Obj{Type: object.OBJ_TYPE_LIST, D: object.List(bits)}, nil
}

func (x *bitsFunctions) cmdInt(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	bitsList := args[0].D.(object.List)

	if len(bitsList) != 64 {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  fmt.Sprintf("bits/int: expected 64 bits, got %d", len(bitsList)),
			},
		}, nil
	}

	var result uint64
	for i := 0; i < 64; i++ {
		if bitsList[i].Type != object.OBJ_TYPE_INTEGER {
			return object.Obj{
				Type: object.OBJ_TYPE_ERROR,
				D: object.Error{
					Position: 0,
					Message:  fmt.Sprintf("bits/int: bit at position %d is not an integer", i),
				},
			}, nil
		}
		bit := bitsList[i].D.(object.Integer)
		if bit != 0 && bit != 1 {
			return object.Obj{
				Type: object.OBJ_TYPE_ERROR,
				D: object.Error{
					Position: 0,
					Message:  fmt.Sprintf("bits/int: bit at position %d must be 0 or 1, got %d", i, bit),
				},
			}, nil
		}
		if bit == 1 {
			result |= (1 << i)
		}
	}

	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(int64(result))}, nil
}

func (x *bitsFunctions) cmdReal(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	bitsList := args[0].D.(object.List)

	if len(bitsList) != 64 {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  fmt.Sprintf("bits/real: expected 64 bits, got %d", len(bitsList)),
			},
		}, nil
	}

	var result uint64
	for i := 0; i < 64; i++ {
		if bitsList[i].Type != object.OBJ_TYPE_INTEGER {
			return object.Obj{
				Type: object.OBJ_TYPE_ERROR,
				D: object.Error{
					Position: 0,
					Message:  fmt.Sprintf("bits/real: bit at position %d is not an integer", i),
				},
			}, nil
		}
		bit := bitsList[i].D.(object.Integer)
		if bit != 0 && bit != 1 {
			return object.Obj{
				Type: object.OBJ_TYPE_ERROR,
				D: object.Error{
					Position: 0,
					Message:  fmt.Sprintf("bits/real: bit at position %d must be 0 or 1, got %d", i, bit),
				},
			}, nil
		}
		if bit == 1 {
			result |= (1 << i)
		}
	}

	realVal := math.Float64frombits(result)
	return object.Obj{Type: object.OBJ_TYPE_REAL, D: object.Real(realVal)}, nil
}
