package env

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bosley/slpx/pkg/object"
	"github.com/bosley/slpx/pkg/slp"
)

type coreFunctions struct{}

func NewCoreFunctions() FunctionGroup {
	return &coreFunctions{}
}

func (c *coreFunctions) Name() string {
	return "core"
}

func (c *coreFunctions) Functions() map[object.Identifier]EnvFunction {
	return map[object.Identifier]EnvFunction{
		"set": {
			EvaluateArgs: false,
			Parameters: []EnvParameter{
				{Name: "name", Type: object.OBJ_TYPE_IDENTIFIER},
				{Name: "value", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_ANY,
			Body:       cmdSet,
		},
		"putln": {
			EvaluateArgs: true,
			Parameters: []EnvParameter{
				{Name: "args", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_NONE,
			Variadic:   true,
			Body:       cmdPutln,
		},
		"fn": {
			EvaluateArgs: false,
			Parameters: []EnvParameter{
				{Name: "params", Type: object.OBJ_TYPE_LIST},
				{Name: "body", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_FUNCTION,
			Variadic:   true,
			Body:       cmdFn,
		},
		"try": {
			EvaluateArgs: false,
			Parameters: []EnvParameter{
				{Name: "expr", Type: object.OBJ_TYPE_ANY},
				{Name: "handler", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_ANY,
			Body:       cmdTry,
		},
		"do": {
			EvaluateArgs: false,
			Parameters: []EnvParameter{
				{Name: "exprs", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_ANY,
			Variadic:   true,
			Body:       cmdDo,
		},
		"drop": {
			EvaluateArgs: false,
			Parameters: []EnvParameter{
				{Name: "name", Type: object.OBJ_TYPE_IDENTIFIER},
			},
			ReturnType: object.OBJ_TYPE_NONE,
			Body:       cmdDrop,
		},
		"qu": {
			EvaluateArgs: false,
			Parameters: []EnvParameter{
				{Name: "expr", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_SOME,
			Body:       cmdQu,
		},
		"uq": {
			EvaluateArgs: false,
			Parameters: []EnvParameter{
				{Name: "quoted", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_ANY,
			Body:       cmdUq,
		},
		"use": {
			EvaluateArgs: true,
			Parameters: []EnvParameter{
				{Name: "paths", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_NONE,
			Variadic:   true,
			Body:       cmdUse,
		},
		"exit": {
			EvaluateArgs: false,
			Parameters: []EnvParameter{
				{Name: "code", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_NONE,
			Body:       cmdExit,
		},
		"if": {
			EvaluateArgs: false,
			Parameters: []EnvParameter{
				{Name: "condition", Type: object.OBJ_TYPE_ANY},
				{Name: "true_body", Type: object.OBJ_TYPE_ANY},
				{Name: "false_body", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_ANY,
			Body:       cmdIf,
		},
	}
}

func cmdSet(ctx EvaluationContext, args object.List) (object.Obj, error) {
	if len(args) != 2 {
		return object.Obj{}, fmt.Errorf("set: requires 2 arguments, got %d", len(args))
	}

	evalCtx := ctx.(*evalCtx)

	if args[0].Type != object.OBJ_TYPE_IDENTIFIER {
		return object.Obj{}, fmt.Errorf("set: first argument must be identifier, got %s", args[0].Type)
	}

	name := args[0].D.(object.Identifier)

	value, err := ctx.Evaluate(args[1])
	if err != nil {
		return object.Obj{}, err
	}

	evalCtx.mem.Set(name, value, true)
	return value, nil
}

func cmdPutln(ctx EvaluationContext, args object.List) (object.Obj, error) {
	evalCtx := ctx.(*evalCtx)

	for i, arg := range args {
		if i > 0 {
			evalCtx.io.WriteString(" ")
		}

		switch arg.Type {
		case object.OBJ_TYPE_STRING:
			evalCtx.io.WriteString(arg.D.(string))
		default:
			evalCtx.io.WriteString(arg.Encode())
		}
	}

	evalCtx.io.WriteString("\n")
	return object.Obj{Type: object.OBJ_TYPE_NONE, D: object.None{}}, nil
}

func cmdFn(ctx EvaluationContext, args object.List) (object.Obj, error) {
	if len(args) < 2 {
		return object.Obj{}, fmt.Errorf("fn: requires at least 2 arguments (params, body...)")
	}

	if args[0].Type != object.OBJ_TYPE_LIST {
		return object.Obj{}, fmt.Errorf("fn: parameter list must be a list, got %s", args[0].Type)
	}

	paramList := args[0].D.(object.List)

	isVariadic := false
	var parameters []object.Parameter

	if len(paramList) == 1 && paramList[0].Type == object.OBJ_TYPE_IDENTIFIER {
		ident := paramList[0].D.(object.Identifier)
		if ident == ".." {
			isVariadic = true
			parameters = []object.Parameter{}
		} else {
			return object.Obj{}, fmt.Errorf("fn: single parameter must be '..' for variadic or name-type pair")
		}
	} else if len(paramList) == 0 {
		parameters = []object.Parameter{}
	} else {
		if len(paramList)%2 != 0 {
			return object.Obj{}, fmt.Errorf("fn: parameters must be name-type pairs")
		}

		parameters = make([]object.Parameter, len(paramList)/2)
		for i := 0; i < len(paramList); i += 2 {
			nameObj := paramList[i]
			typeObj := paramList[i+1]

			if nameObj.Type != object.OBJ_TYPE_IDENTIFIER {
				return object.Obj{}, fmt.Errorf("fn: parameter name must be identifier, got %s", nameObj.Type)
			}
			if typeObj.Type != object.OBJ_TYPE_IDENTIFIER {
				return object.Obj{}, fmt.Errorf("fn: parameter type must be identifier, got %s", typeObj.Type)
			}

			name := nameObj.D.(object.Identifier)
			typeIdent := typeObj.D.(object.Identifier)

			objType, err := object.GetTypeFromIdentifier(typeIdent)
			if err != nil {
				return object.Obj{}, err
			}

			parameters[i/2] = object.Parameter{
				Name: name,
				Type: objType,
			}
		}
	}

	returnType := object.OBJ_TYPE_ANY
	bodyStartIdx := 1

	if len(args) > 1 && args[1].Type == object.OBJ_TYPE_IDENTIFIER {
		returnTypeIdent := args[1].D.(object.Identifier)
		parsedReturnType, err := object.GetTypeFromIdentifier(returnTypeIdent)
		if err == nil {
			returnType = parsedReturnType
			bodyStartIdx = 2
		}
	}

	if bodyStartIdx >= len(args) {
		return object.Obj{}, fmt.Errorf("fn: function body cannot be empty")
	}

	body := args[bodyStartIdx:]

	evalCtx := ctx.(*evalCtx)

	return object.Obj{
		Type: object.OBJ_TYPE_FUNCTION,
		D: object.Function{
			Parameters: parameters,
			ReturnType: returnType,
			Variadic:   isVariadic,
			Body:       body,
		},
		C: evalCtx.mem,
	}, nil
}

func cmdTry(ctx EvaluationContext, args object.List) (object.Obj, error) {
	evalCtx := ctx.(*evalCtx)
	if len(args) != 2 {
		argPos := uint16(0)
		if len(args) > 0 {
			argPos = args[0].Pos
		}
		return evalCtx.makeError(argPos, fmt.Sprintf("try: requires 2 arguments, got %d", len(args))), nil
	}

	result, err := ctx.Evaluate(args[0])
	if err != nil {
		return object.Obj{}, err
	}

	if result.Type == object.OBJ_TYPE_ERROR {

		errorObj := result.D.(object.Error)
		errorString := object.Obj{
			Type: object.OBJ_TYPE_STRING,
			D:    errorObj.Message,
		}

		evalCtx.mem.Set("$error", errorString, false)

		handlerResult, handlerErr := ctx.Evaluate(args[1])

		evalCtx.mem.Delete("$error", false)

		return handlerResult, handlerErr
	}

	return result, nil
}

func cmdDo(ctx EvaluationContext, args object.List) (object.Obj, error) {
	if len(args) == 0 {
		return object.Obj{Type: object.OBJ_TYPE_NONE, D: object.None{}}, nil
	}

	var result object.Obj
	var err error
	for _, arg := range args {
		result, err = ctx.Evaluate(arg)
		if err != nil {
			return object.Obj{}, err
		}
		if result.Type == object.OBJ_TYPE_ERROR {
			return result, nil
		}
	}

	return result, nil
}

func cmdDrop(ctx EvaluationContext, args object.List) (object.Obj, error) {
	evalCtx := ctx.(*evalCtx)
	if len(args) != 1 {
		argPos := uint16(0)
		if len(args) > 0 {
			argPos = args[0].Pos
		}
		return evalCtx.makeError(argPos, fmt.Sprintf("drop: requires 1 argument, got %d", len(args))), nil
	}

	if args[0].Type != object.OBJ_TYPE_IDENTIFIER {
		return evalCtx.makeErrorFromObj(args[0], fmt.Sprintf("drop: argument must be identifier, got %s", args[0].Type)), nil
	}

	name := args[0].D.(object.Identifier)
	evalCtx.mem.Delete(name, true)
	return object.Obj{Type: object.OBJ_TYPE_NONE, D: object.None{}}, nil
}

func cmdQu(ctx EvaluationContext, args object.List) (object.Obj, error) {
	evalCtx := ctx.(*evalCtx)
	if len(args) != 1 {
		argPos := uint16(0)
		if len(args) > 0 {
			argPos = args[0].Pos
		}
		return evalCtx.makeError(argPos, fmt.Sprintf("qu: requires 1 argument, got %d", len(args))), nil
	}

	return object.Obj{
		Type: object.OBJ_TYPE_SOME,
		D:    args[0],
	}, nil
}

func cmdUq(ctx EvaluationContext, args object.List) (object.Obj, error) {
	evalCtx := ctx.(*evalCtx)
	if len(args) != 1 {
		argPos := uint16(0)
		if len(args) > 0 {
			argPos = args[0].Pos
		}
		return evalCtx.makeError(argPos, fmt.Sprintf("uq: requires 1 argument, got %d", len(args))), nil
	}

	evaluated, err := ctx.Evaluate(args[0])
	if err != nil {
		return object.Obj{}, err
	}

	if evaluated.Type == object.OBJ_TYPE_ERROR {
		return evaluated, nil
	}

	if evaluated.Type != object.OBJ_TYPE_SOME {
		return evalCtx.makeErrorFromObj(args[0], fmt.Sprintf("uq: argument must be quoted (type 'some'), got %s", evaluated.Type)), nil
	}

	return evaluated.D.(object.Some), nil
}

func cmdUse(ctx EvaluationContext, args object.List) (object.Obj, error) {
	evalCtx := ctx.(*evalCtx)
	if len(args) == 0 {
		return evalCtx.makeError(0, "use: requires at least 1 argument"), nil
	}

	for _, arg := range args {
		if arg.Type != object.OBJ_TYPE_STRING {
			return evalCtx.makeErrorFromObj(arg, fmt.Sprintf("use: argument must be string, got %s", arg.Type)), nil
		}

		filePath := arg.D.(string)

		var fullPath string
		if filepath.IsAbs(filePath) {
			fullPath = filePath
		} else {
			if evalCtx.currentFilePath != "" {
				currentDir := filepath.Dir(evalCtx.currentFilePath)
				fullPath = filepath.Join(currentDir, filePath)
			} else {
				fullPath = filePath
			}
		}

		absPath, err := filepath.Abs(fullPath)
		if err != nil {
			absPath = fullPath
		}

		if evalCtx.importedFiles[absPath] {
			continue
		}

		evalCtx.importedFiles[absPath] = true

		content, err := evalCtx.fs.ReadFile(fullPath)
		if err != nil {
			return evalCtx.makeErrorFromObj(arg, fmt.Sprintf("use: failed to read file %s: %v", fullPath, err)), nil
		}

		parser := slp.NewParser(string(content))
		items, err := parser.ParseAll()
		if err != nil {
			return evalCtx.makeErrorFromObj(arg, fmt.Sprintf("use: failed to parse file %s: %v", fullPath, err)), nil
		}

		previousFilePath := evalCtx.currentFilePath
		evalCtx.currentFilePath = fullPath

		for itemIdx, item := range items {
			result, err := ctx.Evaluate(item)
			if err != nil {
				evalCtx.currentFilePath = previousFilePath
				return evalCtx.makeErrorFromObj(item, fmt.Sprintf("use: error evaluating file %s at item %d: %v", fullPath, itemIdx, err)), nil
			}
			if result.Type == object.OBJ_TYPE_ERROR {
				evalCtx.currentFilePath = previousFilePath
				errObj := result.D.(object.Error)
				return evalCtx.makeErrorFromObj(item, fmt.Sprintf("use: file %s item %d produced error: %s", fullPath, itemIdx, errObj.Message)), nil
			}
		}

		evalCtx.currentFilePath = previousFilePath
	}

	return object.Obj{Type: object.OBJ_TYPE_NONE, D: object.None{}}, nil
}

func cmdExit(ctx EvaluationContext, args object.List) (object.Obj, error) {
	evalCtx := ctx.(*evalCtx)
	if len(args) != 1 {
		argPos := uint16(0)
		if len(args) > 0 {
			argPos = args[0].Pos
		}
		return evalCtx.makeError(argPos, fmt.Sprintf("exit: requires 1 argument, got %d", len(args))), nil
	}

	arg := args[0]

	if arg.Type != object.OBJ_TYPE_INTEGER && arg.Type != object.OBJ_TYPE_IDENTIFIER {
		return evalCtx.makeErrorFromObj(arg, fmt.Sprintf("exit: argument must be integer or identifier, got %s", arg.Type)), nil
	}

	var exitCode int

	if arg.Type == object.OBJ_TYPE_INTEGER {
		exitCode = int(arg.D.(object.Integer))
	} else {
		result, err := ctx.Evaluate(arg)
		if err != nil {
			os.Exit(1)
		}

		if result.Type != object.OBJ_TYPE_INTEGER {
			os.Exit(1)
		}

		exitCode = int(result.D.(object.Integer))
	}

	os.Exit(exitCode)

	return object.Obj{Type: object.OBJ_TYPE_NONE, D: object.None{}}, nil
}

func cmdIf(ctx EvaluationContext, args object.List) (object.Obj, error) {
	evalCtx := ctx.(*evalCtx)
	if len(args) != 3 {
		argPos := uint16(0)
		if len(args) > 0 {
			argPos = args[0].Pos
		}
		return evalCtx.makeError(argPos, fmt.Sprintf("if: requires 3 arguments, got %d", len(args))), nil
	}

	condition, err := ctx.Evaluate(args[0])
	if err != nil {
		return object.Obj{}, err
	}

	if condition.Type == object.OBJ_TYPE_ERROR {
		return condition, nil
	}

	if condition.Type != object.OBJ_TYPE_INTEGER {
		return evalCtx.makeErrorFromObj(args[0], fmt.Sprintf("if: condition must evaluate to integer, got %s", condition.Type)), nil
	}

	condValue := condition.D.(object.Integer)

	if condValue > 0 {
		return ctx.Evaluate(args[1])
	}

	return ctx.Evaluate(args[2])
}
