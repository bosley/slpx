package str

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/bosley/slpx/pkg/slp/env"
	"github.com/bosley/slpx/pkg/slp/object"
)

type strFunctions struct {
	mu        sync.Mutex
	precision uint8
}

func NewStrFunctions() env.FunctionGroup {
	return &strFunctions{
		precision: 6,
	}
}

func (s *strFunctions) Name() string {
	return "str"
}

func (s *strFunctions) Functions() map[object.Identifier]env.EnvFunction {
	return map[object.Identifier]env.EnvFunction{
		"str/eq": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "a", Type: object.OBJ_TYPE_STRING},
				{Name: "b", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdStrEq,
		},
		"str/len": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "s", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdStrLen,
		},
		"str/clear": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "s", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_STRING,
			Body:       cmdStrClear,
		},
		"str/from": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "obj", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_STRING,
			Body:       s.cmdStrFrom,
		},
		"str/int": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "s", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdStrInt,
		},
		"str/real": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "s", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_REAL,
			Body:       cmdStrReal,
		},
		"str/list": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "s", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_LIST,
			Body:       cmdStrList,
		},
		"str/concat": {
			EvaluateArgs: true,
			Parameters:   []env.EnvParameter{},
			ReturnType:   object.OBJ_TYPE_STRING,
			Variadic:     true,
			Body:         cmdStrConcat,
		},
		"str/upper": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "s", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_STRING,
			Body:       cmdStrUpper,
		},
		"str/lower": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "s", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_STRING,
			Body:       cmdStrLower,
		},
		"str/trim": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "s", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_STRING,
			Body:       cmdStrTrim,
		},
		"str/contains": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "s", Type: object.OBJ_TYPE_STRING},
				{Name: "substr", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdStrContains,
		},
		"str/index": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "s", Type: object.OBJ_TYPE_STRING},
				{Name: "substr", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdStrIndex,
		},
		"str/slice": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "s", Type: object.OBJ_TYPE_STRING},
				{Name: "start", Type: object.OBJ_TYPE_INTEGER},
				{Name: "end", Type: object.OBJ_TYPE_INTEGER},
			},
			ReturnType: object.OBJ_TYPE_STRING,
			Body:       cmdStrSlice,
		},
		"str/split": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "s", Type: object.OBJ_TYPE_STRING},
				{Name: "sep", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_LIST,
			Body:       cmdStrSplit,
		},
		"str/replace": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "s", Type: object.OBJ_TYPE_STRING},
				{Name: "old", Type: object.OBJ_TYPE_STRING},
				{Name: "new", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_STRING,
			Body:       cmdStrReplace,
		},
		"str/precision": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "p", Type: object.OBJ_TYPE_INTEGER},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       s.cmdStrPrecision,
		},
	}
}

func cmdStrEq(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	a := args[0].D.(string)
	b := args[1].D.(string)
	if a == b {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdStrLen(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	s := args[0].D.(string)
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(len([]rune(s)))}, nil
}

func cmdStrClear(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	return object.Obj{Type: object.OBJ_TYPE_STRING, D: ""}, nil
}

func (s *strFunctions) cmdStrFrom(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	obj := args[0]
	if obj.Type == object.OBJ_TYPE_STRING {
		return obj, nil
	}

	if obj.Type == object.OBJ_TYPE_REAL {
		s.mu.Lock()
		prec := s.precision
		s.mu.Unlock()

		realVal := obj.D.(object.Real)
		format := fmt.Sprintf("%%.%df", prec)
		return object.Obj{Type: object.OBJ_TYPE_STRING, D: fmt.Sprintf(format, realVal)}, nil
	}

	return object.Obj{Type: object.OBJ_TYPE_STRING, D: obj.Encode()}, nil
}

func cmdStrInt(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	s := args[0].D.(string)
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "str/int: failed to parse integer: " + err.Error(),
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(val)}, nil
}

func cmdStrReal(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	s := args[0].D.(string)
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "str/real: failed to parse real: " + err.Error(),
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_REAL, D: object.Real(val)}, nil
}

func cmdStrList(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	s := args[0].D.(string)
	runes := []rune(s)
	result := make(object.List, len(runes))
	for i, r := range runes {
		result[i] = object.Obj{Type: object.OBJ_TYPE_STRING, D: string(r)}
	}
	return object.Obj{Type: object.OBJ_TYPE_LIST, D: result}, nil
}

func cmdStrConcat(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	if len(args) == 0 {
		return object.Obj{Type: object.OBJ_TYPE_STRING, D: ""}, nil
	}
	var builder strings.Builder
	for i, arg := range args {
		if arg.Type != object.OBJ_TYPE_STRING {
			return object.Obj{
				Type: object.OBJ_TYPE_ERROR,
				D: object.Error{
					Position: 0,
					Message:  "str/concat: all arguments must be strings, got " + string(arg.Type) + " at position " + strconv.Itoa(i),
				},
			}, nil
		}
		builder.WriteString(arg.D.(string))
	}
	return object.Obj{Type: object.OBJ_TYPE_STRING, D: builder.String()}, nil
}

func cmdStrUpper(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	s := args[0].D.(string)
	return object.Obj{Type: object.OBJ_TYPE_STRING, D: strings.ToUpper(s)}, nil
}

func cmdStrLower(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	s := args[0].D.(string)
	return object.Obj{Type: object.OBJ_TYPE_STRING, D: strings.ToLower(s)}, nil
}

func cmdStrTrim(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	s := args[0].D.(string)
	return object.Obj{Type: object.OBJ_TYPE_STRING, D: strings.TrimSpace(s)}, nil
}

func cmdStrContains(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	s := args[0].D.(string)
	substr := args[1].D.(string)
	if strings.Contains(s, substr) {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdStrIndex(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	s := args[0].D.(string)
	substr := args[1].D.(string)
	idx := strings.Index(s, substr)
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(idx)}, nil
}

func cmdStrSlice(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	s := args[0].D.(string)
	start := int(args[1].D.(object.Integer))
	end := int(args[2].D.(object.Integer))

	runes := []rune(s)
	length := len(runes)

	if start < 0 {
		start = 0
	}
	if end > length {
		end = length
	}
	if start > end {
		start = end
	}

	return object.Obj{Type: object.OBJ_TYPE_STRING, D: string(runes[start:end])}, nil
}

func cmdStrSplit(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	s := args[0].D.(string)
	sep := args[1].D.(string)
	parts := strings.Split(s, sep)
	result := make(object.List, len(parts))
	for i, part := range parts {
		result[i] = object.Obj{Type: object.OBJ_TYPE_STRING, D: part}
	}
	return object.Obj{Type: object.OBJ_TYPE_LIST, D: result}, nil
}

func cmdStrReplace(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	s := args[0].D.(string)
	old := args[1].D.(string)
	new := args[2].D.(string)
	return object.Obj{Type: object.OBJ_TYPE_STRING, D: strings.ReplaceAll(s, old, new)}, nil
}

func (s *strFunctions) cmdStrPrecision(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	prec := args[0].D.(object.Integer)

	if prec < 0 {
		prec = 0
	}
	if prec > 255 {
		prec = 255
	}

	s.mu.Lock()
	s.precision = uint8(prec)
	s.mu.Unlock()

	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: prec}, nil
}
