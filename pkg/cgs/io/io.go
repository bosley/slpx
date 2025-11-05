package io

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bosley/slpx/pkg/env"
	"github.com/bosley/slpx/pkg/object"
)

type ioFunctions struct {
	io        env.IO
	precision int
}

func NewIoFunctions() *ioFunctions {
	return &ioFunctions{
		precision: 6,
	}
}

func (i *ioFunctions) Setup(runtime env.Runtime) {
	i.io = runtime.GetIO()
}

func (i *ioFunctions) Name() string {
	return "io"
}

func (i *ioFunctions) Functions() map[object.Identifier]env.EnvFunction {
	return map[object.Identifier]env.EnvFunction{
		"io/out": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "args", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_NONE,
			Variadic:   true,
			Body:       i.cmdOut,
		},
		"io/color/fg": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "color", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_STRING,
			Body:       i.cmdColorFg,
		},
		"io/color/bg": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "color", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_STRING,
			Body:       i.cmdColorBg,
		},
		"io/color/reset": {
			EvaluateArgs: true,
			Parameters:   []env.EnvParameter{},
			ReturnType:   object.OBJ_TYPE_STRING,
			Body:         i.cmdColorReset,
		},
		"io/in": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "prompt", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_STRING,
			Body:       i.cmdIn,
		},
		"io/in/int": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "prompt", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       i.cmdInInt,
		},
		"io/in/real": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "prompt", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_REAL,
			Body:       i.cmdInReal,
		},
		"io/out/set_precision": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "precision", Type: object.OBJ_TYPE_INTEGER},
			},
			ReturnType: object.OBJ_TYPE_NONE,
			Body:       i.cmdSetPrecision,
		},
		"io/flush": {
			EvaluateArgs: true,
			Parameters:   []env.EnvParameter{},
			ReturnType:   object.OBJ_TYPE_NONE,
			Body:         i.cmdFlush,
		},
	}
}

func (i *ioFunctions) cmdOut(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	for _, arg := range args {
		var str string
		switch arg.Type {
		case object.OBJ_TYPE_STRING:
			str = arg.D.(string)
		case object.OBJ_TYPE_INTEGER:
			str = fmt.Sprintf("%d", arg.D.(object.Integer))
		case object.OBJ_TYPE_REAL:
			str = fmt.Sprintf("%.*f", i.precision, float64(arg.D.(object.Real)))
		default:
			str = arg.Encode()
		}

		i.io.WriteString(str)
		i.io.Flush()
	}
	return object.Obj{Type: object.OBJ_TYPE_NONE, D: object.None{}}, nil
}

func (i *ioFunctions) cmdColorFg(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	hexColor := args[0].D.(string)
	r, g, b, err := i.parseHexColor(hexColor)
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "invalid hex color: " + err.Error(),
			},
		}, nil
	}
	ansiCode := fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
	return object.Obj{Type: object.OBJ_TYPE_STRING, D: ansiCode}, nil
}

func (i *ioFunctions) cmdColorBg(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	hexColor := args[0].D.(string)
	r, g, b, err := i.parseHexColor(hexColor)
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "invalid hex color: " + err.Error(),
			},
		}, nil
	}
	ansiCode := fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b)
	return object.Obj{Type: object.OBJ_TYPE_STRING, D: ansiCode}, nil
}

func (i *ioFunctions) cmdColorReset(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	return object.Obj{Type: object.OBJ_TYPE_STRING, D: "\033[0m"}, nil
}

func (i *ioFunctions) cmdIn(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	prompt := args[0].D.(string)
	i.io.WriteString(prompt)
	i.io.Flush()

	line, err := i.io.ReadLine()
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "failed to read input: " + err.Error(),
			},
		}, nil
	}

	return object.Obj{Type: object.OBJ_TYPE_STRING, D: line}, nil
}

func (i *ioFunctions) cmdInInt(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	prompt := args[0].D.(string)
	i.io.WriteString(prompt)
	i.io.Flush()

	line, err := i.io.ReadLine()
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "failed to read input: " + err.Error(),
			},
		}, nil
	}

	line = strings.TrimSpace(line)
	intVal, err := strconv.ParseInt(line, 10, 64)
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "input is not a valid integer: " + line,
			},
		}, nil
	}

	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(intVal)}, nil
}

func (i *ioFunctions) cmdInReal(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	prompt := args[0].D.(string)
	i.io.WriteString(prompt)
	i.io.Flush()

	line, err := i.io.ReadLine()
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "failed to read input: " + err.Error(),
			},
		}, nil
	}

	line = strings.TrimSpace(line)
	realVal, err := strconv.ParseFloat(line, 64)
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "input is not a valid real number: " + line,
			},
		}, nil
	}

	return object.Obj{Type: object.OBJ_TYPE_REAL, D: object.Real(realVal)}, nil
}

func (i *ioFunctions) cmdSetPrecision(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	precision := int(args[0].D.(object.Integer))
	if precision < 0 {
		precision = 0
	}
	if precision > 20 {
		precision = 20
	}
	i.precision = precision
	return object.Obj{Type: object.OBJ_TYPE_NONE, D: object.None{}}, nil
}

func (i *ioFunctions) cmdFlush(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	if err := i.io.Flush(); err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "failed to flush output: " + err.Error(),
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_NONE, D: object.None{}}, nil
}

func (i *ioFunctions) parseHexColor(hexColor string) (r, g, b int, err error) {
	hexColor = strings.TrimPrefix(hexColor, "#")

	if len(hexColor) != 6 {
		return 0, 0, 0, fmt.Errorf("hex color must be 6 characters (got %d)", len(hexColor))
	}

	val, err := strconv.ParseUint(hexColor, 16, 32)
	if err != nil {
		return 0, 0, 0, err
	}

	r = int((val >> 16) & 0xFF)
	g = int((val >> 8) & 0xFF)
	b = int(val & 0xFF)

	return r, g, b, nil
}
