package numbers

import (
	"math"
	"math/rand/v2"
	"strconv"

	"github.com/bosley/slpx/pkg/slp/env"
	"github.com/bosley/slpx/pkg/slp/object"
)

type arithFunctions struct{}

func NewArithFunctions() env.FunctionGroup {
	return &arithFunctions{}
}

func (a *arithFunctions) Name() string {
	return "arith"
}

func (a *arithFunctions) Functions() map[object.Identifier]env.EnvFunction {
	return map[object.Identifier]env.EnvFunction{
		"int/add": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "a", Type: object.OBJ_TYPE_INTEGER},
				{Name: "b", Type: object.OBJ_TYPE_INTEGER},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdIntAdd,
		},
		"int/sub": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "a", Type: object.OBJ_TYPE_INTEGER},
				{Name: "b", Type: object.OBJ_TYPE_INTEGER},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdIntSub,
		},
		"int/mul": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "a", Type: object.OBJ_TYPE_INTEGER},
				{Name: "b", Type: object.OBJ_TYPE_INTEGER},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdIntMul,
		},
		"int/div": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "a", Type: object.OBJ_TYPE_INTEGER},
				{Name: "b", Type: object.OBJ_TYPE_INTEGER},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdIntDiv,
		},
		"int/mod": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "a", Type: object.OBJ_TYPE_INTEGER},
				{Name: "b", Type: object.OBJ_TYPE_INTEGER},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdIntMod,
		},
		"int/pow": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "a", Type: object.OBJ_TYPE_INTEGER},
				{Name: "b", Type: object.OBJ_TYPE_INTEGER},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdIntPow,
		},
		"int/sum": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "values", Type: object.OBJ_TYPE_INTEGER},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Variadic:   true,
			Body:       cmdIntSum,
		},
		"real/add": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "a", Type: object.OBJ_TYPE_REAL},
				{Name: "b", Type: object.OBJ_TYPE_REAL},
			},
			ReturnType: object.OBJ_TYPE_REAL,
			Body:       cmdRealAdd,
		},
		"real/sub": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "a", Type: object.OBJ_TYPE_REAL},
				{Name: "b", Type: object.OBJ_TYPE_REAL},
			},
			ReturnType: object.OBJ_TYPE_REAL,
			Body:       cmdRealSub,
		},
		"real/mul": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "a", Type: object.OBJ_TYPE_REAL},
				{Name: "b", Type: object.OBJ_TYPE_REAL},
			},
			ReturnType: object.OBJ_TYPE_REAL,
			Body:       cmdRealMul,
		},
		"real/div": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "a", Type: object.OBJ_TYPE_REAL},
				{Name: "b", Type: object.OBJ_TYPE_REAL},
			},
			ReturnType: object.OBJ_TYPE_REAL,
			Body:       cmdRealDiv,
		},
		"real/pow": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "a", Type: object.OBJ_TYPE_REAL},
				{Name: "b", Type: object.OBJ_TYPE_REAL},
			},
			ReturnType: object.OBJ_TYPE_REAL,
			Body:       cmdRealPow,
		},
		"real/sum": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "values", Type: object.OBJ_TYPE_REAL},
			},
			ReturnType: object.OBJ_TYPE_REAL,
			Variadic:   true,
			Body:       cmdRealSum,
		},
		"int/real": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "value", Type: object.OBJ_TYPE_INTEGER},
			},
			ReturnType: object.OBJ_TYPE_REAL,
			Body:       cmdIntToReal,
		},
		"real/int": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "value", Type: object.OBJ_TYPE_REAL},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdRealToInt,
		},
		"int/eq": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "a", Type: object.OBJ_TYPE_INTEGER},
				{Name: "b", Type: object.OBJ_TYPE_INTEGER},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdIntEq,
		},
		"int/gt": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "a", Type: object.OBJ_TYPE_INTEGER},
				{Name: "b", Type: object.OBJ_TYPE_INTEGER},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdIntGt,
		},
		"int/gte": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "a", Type: object.OBJ_TYPE_INTEGER},
				{Name: "b", Type: object.OBJ_TYPE_INTEGER},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdIntGte,
		},
		"int/lt": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "a", Type: object.OBJ_TYPE_INTEGER},
				{Name: "b", Type: object.OBJ_TYPE_INTEGER},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdIntLt,
		},
		"int/lte": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "a", Type: object.OBJ_TYPE_INTEGER},
				{Name: "b", Type: object.OBJ_TYPE_INTEGER},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdIntLte,
		},
		"real/eq": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "a", Type: object.OBJ_TYPE_REAL},
				{Name: "b", Type: object.OBJ_TYPE_REAL},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdRealEq,
		},
		"real/gt": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "a", Type: object.OBJ_TYPE_REAL},
				{Name: "b", Type: object.OBJ_TYPE_REAL},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdRealGt,
		},
		"real/gte": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "a", Type: object.OBJ_TYPE_REAL},
				{Name: "b", Type: object.OBJ_TYPE_REAL},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdRealGte,
		},
		"real/lt": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "a", Type: object.OBJ_TYPE_REAL},
				{Name: "b", Type: object.OBJ_TYPE_REAL},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdRealLt,
		},
		"real/lte": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "a", Type: object.OBJ_TYPE_REAL},
				{Name: "b", Type: object.OBJ_TYPE_REAL},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdRealLte,
		},
		"int/rand": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "lower", Type: object.OBJ_TYPE_INTEGER},
				{Name: "upper", Type: object.OBJ_TYPE_INTEGER},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdIntRand,
		},
		"real/rand": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "lower", Type: object.OBJ_TYPE_REAL},
				{Name: "upper", Type: object.OBJ_TYPE_REAL},
			},
			ReturnType: object.OBJ_TYPE_REAL,
			Body:       cmdRealRand,
		},
		"real/sqrt": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "value", Type: object.OBJ_TYPE_REAL},
			},
			ReturnType: object.OBJ_TYPE_REAL,
			Body:       cmdRealSqrt,
		},
		"real/exp": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "value", Type: object.OBJ_TYPE_REAL},
			},
			ReturnType: object.OBJ_TYPE_REAL,
			Body:       cmdRealExp,
		},
		"real/log": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "value", Type: object.OBJ_TYPE_REAL},
			},
			ReturnType: object.OBJ_TYPE_REAL,
			Body:       cmdRealLog,
		},
		"real/ceil": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "value", Type: object.OBJ_TYPE_REAL},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdRealCeil,
		},
		"real/round": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "value", Type: object.OBJ_TYPE_REAL},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdRealRound,
		},
		"real/is-nan": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "value", Type: object.OBJ_TYPE_REAL},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdRealIsNaN,
		},
		"real/is-inf": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "value", Type: object.OBJ_TYPE_REAL},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdRealIsInf,
		},
		"real/is-finite": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "value", Type: object.OBJ_TYPE_REAL},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdRealIsFinite,
		},
		"int/abs": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "value", Type: object.OBJ_TYPE_INTEGER},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdIntAbs,
		},
		"real/abs": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "value", Type: object.OBJ_TYPE_REAL},
			},
			ReturnType: object.OBJ_TYPE_REAL,
			Body:       cmdRealAbs,
		},
	}
}

func cmdIntAdd(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	a := args[0].D.(object.Integer)
	b := args[1].D.(object.Integer)
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: a + b}, nil
}

func cmdIntSub(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	a := args[0].D.(object.Integer)
	b := args[1].D.(object.Integer)
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: a - b}, nil
}

func cmdIntMul(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	a := args[0].D.(object.Integer)
	b := args[1].D.(object.Integer)
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: a * b}, nil
}

func cmdIntDiv(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	a := args[0].D.(object.Integer)
	b := args[1].D.(object.Integer)
	if b == 0 {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "int/div: division by zero",
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: a / b}, nil
}

func cmdIntMod(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	a := args[0].D.(object.Integer)
	b := args[1].D.(object.Integer)
	if b == 0 {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "int/mod: modulo by zero",
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: a % b}, nil
}

func cmdIntPow(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	a := args[0].D.(object.Integer)
	b := args[1].D.(object.Integer)
	if b < 0 {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "int/pow: negative exponent not supported for integer power",
			},
		}, nil
	}
	result := object.Integer(1)
	base := a
	exp := b
	for exp > 0 {
		if exp%2 == 1 {
			result *= base
		}
		base *= base
		exp /= 2
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: result}, nil
}

func cmdIntSum(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	if len(args) == 0 {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
	}
	sum := object.Integer(0)
	for i, arg := range args {
		if arg.Type != object.OBJ_TYPE_INTEGER {
			return object.Obj{
				Type: object.OBJ_TYPE_ERROR,
				D: object.Error{
					Position: 0,
					Message:  "int/sum: all arguments must be integers, got " + string(arg.Type) + " at position " + strconv.Itoa(i),
				},
			}, nil
		}
		sum += arg.D.(object.Integer)
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: sum}, nil
}

func cmdRealAdd(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	a := args[0].D.(object.Real)
	b := args[1].D.(object.Real)
	return object.Obj{Type: object.OBJ_TYPE_REAL, D: a + b}, nil
}

func cmdRealSub(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	a := args[0].D.(object.Real)
	b := args[1].D.(object.Real)
	return object.Obj{Type: object.OBJ_TYPE_REAL, D: a - b}, nil
}

func cmdRealMul(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	a := args[0].D.(object.Real)
	b := args[1].D.(object.Real)
	return object.Obj{Type: object.OBJ_TYPE_REAL, D: a * b}, nil
}

func cmdRealDiv(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	a := args[0].D.(object.Real)
	b := args[1].D.(object.Real)
	if b == 0.0 {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "real/div: division by zero",
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_REAL, D: a / b}, nil
}

func cmdRealPow(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	a := args[0].D.(object.Real)
	b := args[1].D.(object.Real)
	result := math.Pow(float64(a), float64(b))
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "real/pow: invalid result (NaN or Inf)",
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_REAL, D: object.Real(result)}, nil
}

func cmdRealSum(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	if len(args) == 0 {
		return object.Obj{Type: object.OBJ_TYPE_REAL, D: object.Real(0.0)}, nil
	}
	sum := object.Real(0.0)
	for i, arg := range args {
		if arg.Type != object.OBJ_TYPE_REAL {
			return object.Obj{
				Type: object.OBJ_TYPE_ERROR,
				D: object.Error{
					Position: 0,
					Message:  "real/sum: all arguments must be reals, got " + string(arg.Type) + " at position " + strconv.Itoa(i),
				},
			}, nil
		}
		sum += arg.D.(object.Real)
	}
	return object.Obj{Type: object.OBJ_TYPE_REAL, D: sum}, nil
}

func cmdIntToReal(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	i := args[0].D.(object.Integer)
	return object.Obj{Type: object.OBJ_TYPE_REAL, D: object.Real(i)}, nil
}

func cmdRealToInt(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	r := args[0].D.(object.Real)
	floored := math.Floor(float64(r))
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(floored)}, nil
}

func cmdIntEq(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	a := args[0].D.(object.Integer)
	b := args[1].D.(object.Integer)
	if a == b {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdIntGt(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	a := args[0].D.(object.Integer)
	b := args[1].D.(object.Integer)
	if a > b {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdIntGte(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	a := args[0].D.(object.Integer)
	b := args[1].D.(object.Integer)
	if a >= b {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdIntLt(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	a := args[0].D.(object.Integer)
	b := args[1].D.(object.Integer)
	if a < b {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdIntLte(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	a := args[0].D.(object.Integer)
	b := args[1].D.(object.Integer)
	if a <= b {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdRealEq(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	a := args[0].D.(object.Real)
	b := args[1].D.(object.Real)
	if a == b {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdRealGt(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	a := args[0].D.(object.Real)
	b := args[1].D.(object.Real)
	if a > b {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdRealGte(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	a := args[0].D.(object.Real)
	b := args[1].D.(object.Real)
	if a >= b {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdRealLt(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	a := args[0].D.(object.Real)
	b := args[1].D.(object.Real)
	if a < b {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdRealLte(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	a := args[0].D.(object.Real)
	b := args[1].D.(object.Real)
	if a <= b {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdIntRand(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	lower := args[0].D.(object.Integer)
	upper := args[1].D.(object.Integer)
	if lower > upper {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "int/rand: lower bound must be less than or equal to upper bound",
			},
		}, nil
	}
	if lower == upper {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: lower}, nil
	}
	rangeSize := upper - lower + 1
	result := lower + object.Integer(rand.IntN(int(rangeSize)))
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: result}, nil
}

func cmdRealRand(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	lower := args[0].D.(object.Real)
	upper := args[1].D.(object.Real)
	if lower > upper {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "real/rand: lower bound must be less than or equal to upper bound",
			},
		}, nil
	}
	if lower == upper {
		return object.Obj{Type: object.OBJ_TYPE_REAL, D: lower}, nil
	}
	rangeSize := upper - lower
	result := lower + object.Real(rand.Float64())*rangeSize
	return object.Obj{Type: object.OBJ_TYPE_REAL, D: result}, nil
}

func cmdRealSqrt(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	value := args[0].D.(object.Real)
	if value < 0.0 {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "real/sqrt: cannot compute square root of negative number",
			},
		}, nil
	}
	result := math.Sqrt(float64(value))
	return object.Obj{Type: object.OBJ_TYPE_REAL, D: object.Real(result)}, nil
}

func cmdRealExp(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	value := args[0].D.(object.Real)
	result := math.Exp(float64(value))
	if math.IsInf(result, 0) {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "real/exp: result overflow (infinity)",
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_REAL, D: object.Real(result)}, nil
}

func cmdRealLog(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	value := args[0].D.(object.Real)
	if value <= 0.0 {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "real/log: logarithm undefined for non-positive numbers",
			},
		}, nil
	}
	result := math.Log(float64(value))
	return object.Obj{Type: object.OBJ_TYPE_REAL, D: object.Real(result)}, nil
}

func cmdRealCeil(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	value := args[0].D.(object.Real)
	ceiled := math.Ceil(float64(value))
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(ceiled)}, nil
}

func cmdRealRound(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	value := args[0].D.(object.Real)
	rounded := math.Round(float64(value))
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(rounded)}, nil
}

func cmdRealIsNaN(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	value := args[0].D.(object.Real)
	if math.IsNaN(float64(value)) {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdRealIsInf(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	value := args[0].D.(object.Real)
	if math.IsInf(float64(value), 0) {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdRealIsFinite(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	value := args[0].D.(object.Real)
	if !math.IsNaN(float64(value)) && !math.IsInf(float64(value), 0) {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdIntAbs(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	value := args[0].D.(object.Integer)
	if value < 0 {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: -value}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: value}, nil
}

func cmdRealAbs(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	value := args[0].D.(object.Real)
	result := math.Abs(float64(value))
	return object.Obj{Type: object.OBJ_TYPE_REAL, D: object.Real(result)}, nil
}
