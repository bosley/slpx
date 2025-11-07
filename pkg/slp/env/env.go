package env

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/bosley/slpx/pkg/slp/object"
)

type IdentifiedType string

const (
	IdentifiedType_None       IdentifiedType = ":_"
	IdentifiedType_Some       IdentifiedType = ":Q"
	IdentifiedType_Any        IdentifiedType = ":*"
	IdentifiedType_List       IdentifiedType = ":L"
	IdentifiedType_Error      IdentifiedType = ":E"
	IdentifiedType_String     IdentifiedType = ":S"
	IdentifiedType_Integer    IdentifiedType = ":I"
	IdentifiedType_Real       IdentifiedType = ":R"
	IdentifiedType_Identifier IdentifiedType = ":X"
	IdentifiedType_Function   IdentifiedType = ":F"
)

type EnvParameter struct {
	Name string
	Type object.ObjType
}

type EnvFunction struct {
	EvaluateArgs bool
	Parameters   []EnvParameter
	ReturnType   object.ObjType
	Variadic     bool
	Body         func(ctx EvaluationContext, args object.List) (object.Obj, error)
}

type FunctionGroup interface {
	Name() string
	Functions() map[object.Identifier]EnvFunction
}

type EvaluationContext interface {
	AddFunctionGroup(group FunctionGroup)
	RemoveFunctionGroup(name string)

	Evaluate(obj object.Obj) (object.Obj, error)
	Execute(list object.List) (object.Obj, error)

	SetCurrentFilePath(path string)
	GetCurrentFilePath() string

	GetRuntime() Runtime
}

type IO interface {
	ReadLine() (string, error)
	ReadAll() ([]byte, error)
	Write(data []byte) (int, error)
	WriteString(s string) (int, error)
	WriteError(data []byte) (int, error)
	WriteErrorString(s string) (int, error)
	Flush() error
	SetStdin(r io.Reader)
	SetStdout(w io.Writer)
	SetStderr(w io.Writer)
}

type FS interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte, perm os.FileMode) error
	AppendFile(path string, data []byte) error
	DeleteFile(path string) error
	RemoveDir(path string) error
	RemoveDirAll(path string) error
	Exists(path string) bool
	IsDir(path string) bool
	IsFile(path string) bool
	ListDir(path string) ([]string, error)
	MkDir(path string, perm os.FileMode) error
	MkDirAll(path string, perm os.FileMode) error
	WorkingDir() string
	SetWorkingDir(path string) error
}

var (
	ErrUndefinedIdentifier = errors.New("undefined identifier")
	ErrEmptyList           = errors.New("cannot execute empty list")
	ErrNotCallable         = errors.New("first element of list is not callable")
	ErrWrongArgCount       = errors.New("wrong number of arguments")
)

type MEM interface {
	Get(key object.Identifier, searchParent bool) (object.Obj, error)
	Set(key object.Identifier, value object.Obj, searchParent bool) error
	Delete(key object.Identifier, searchParent bool) error
	Keys() []object.Identifier
	Values() []object.Obj
	Len() int
	IsEmpty() bool
	Clear()
	GetAll() map[object.Identifier]object.Obj

	// Get a child, with this MEM as parent
	Fork() MEM
}

type Runtime interface {
	GetMEM() MEM
	GetFS() FS
	GetIO() IO
	GetStartPath() string
}

type EvalBuilder struct {
	logger *slog.Logger

	maxRecursionDepth int

	io  IO
	fs  FS
	mem MEM

	functionGroups []FunctionGroup
}

func NewEvalBuilder(logger *slog.Logger) *EvalBuilder {
	return &EvalBuilder{
		logger: logger,
	}
}

func (x *EvalBuilder) WithIO(io IO) *EvalBuilder {
	x.io = io
	return x
}

func (x *EvalBuilder) WithFS(fs FS) *EvalBuilder {
	x.fs = fs
	return x
}

func (x *EvalBuilder) WithMEM(mem MEM) *EvalBuilder {
	x.mem = mem
	return x
}

func (x *EvalBuilder) WithMaxRecursionDepth(depth int) *EvalBuilder {
	x.maxRecursionDepth = depth
	return x
}

func (x *EvalBuilder) WithFunctionGroup(group FunctionGroup) *EvalBuilder {
	x.functionGroups = append(x.functionGroups, group)
	return x
}

func (x *EvalBuilder) Build() EvaluationContext {

	if x.io == nil {
		x.io = DefaultIO()
	}
	if x.fs == nil {
		x.fs = DefaultFS()
	}
	if x.mem == nil {
		x.mem = DefaultMEM()
	}

	functionGroupsMap := make(map[string]FunctionGroup)
	for _, group := range x.functionGroups {
		functionGroupsMap[group.Name()] = group
	}

	return &evalCtx{
		mem:             x.mem,
		io:              x.io,
		fs:              x.fs,
		functionGroups:  functionGroupsMap,
		currentFilePath: "",
		importedFiles:   make(map[string]bool),
	}
}

type evalCtx struct {
	mem MEM
	io  IO
	fs  FS

	functionGroups map[string]FunctionGroup

	currentFilePath string
	importedFiles   map[string]bool
}

var _ EvaluationContext = &evalCtx{}
var _ Runtime = &evalCtx{}

func (e *evalCtx) AddFunctionGroup(group FunctionGroup) {
	if e.functionGroups == nil {
		e.functionGroups = make(map[string]FunctionGroup)
	}
	e.functionGroups[group.Name()] = group
}

func (e *evalCtx) RemoveFunctionGroup(name string) {
	delete(e.functionGroups, name)
}

func (e *evalCtx) SetCurrentFilePath(path string) {
	e.currentFilePath = path
}

func (e *evalCtx) GetCurrentFilePath() string {
	return e.currentFilePath
}

func (e *evalCtx) makeError(pos uint16, message string) object.Obj {
	return object.Obj{
		Type: object.OBJ_TYPE_ERROR,
		D: object.Error{
			File:     e.currentFilePath,
			Position: int(pos),
			Message:  message,
		},
		Pos: pos,
	}
}

func (e *evalCtx) makeErrorFromObj(obj object.Obj, message string) object.Obj {
	return e.makeError(obj.Pos, message)
}

func (e *evalCtx) Evaluate(obj object.Obj) (object.Obj, error) {
	switch obj.Type {
	case object.OBJ_TYPE_NONE, object.OBJ_TYPE_STRING,
		object.OBJ_TYPE_INTEGER, object.OBJ_TYPE_REAL,
		object.OBJ_TYPE_ERROR, object.OBJ_TYPE_FUNCTION:
		return obj, nil

	case object.OBJ_TYPE_SOME:
		quoted := obj.D.(object.Some)
		return quoted, nil

	case object.OBJ_TYPE_IDENTIFIER:
		ident := obj.D.(object.Identifier)
		return e.lookupIdentifier(obj, ident)

	case object.OBJ_TYPE_LIST:
		list := obj.D.(object.List)
		return e.Execute(list)

	default:
		return e.makeErrorFromObj(obj, "unknown object type: "+string(obj.Type)), nil
	}
}

func (e *evalCtx) Execute(list object.List) (object.Obj, error) {
	if len(list) == 0 {
		return object.Obj{Type: object.OBJ_TYPE_NONE, D: object.None{}}, nil
	}

	firstEval, err := e.Evaluate(list[0])
	if err != nil {
		return object.Obj{}, err
	}
	if firstEval.Type == object.OBJ_TYPE_ERROR {
		return firstEval, nil
	}

	switch firstEval.Type {
	case object.OBJ_TYPE_FUNCTION:
		return e.executeObjectFunction(firstEval, list[1:])

	case object.OBJ_TYPE_IDENTIFIER:
		ident := firstEval.D.(object.Identifier)
		envFunction, found := e.lookupEnvFunction(ident)
		if !found {
			return e.makeErrorFromObj(list[0], "function not found: "+string(ident)), nil
		}
		return e.executeEnvFunction(envFunction, list[1:])

	default:
		return e.makeErrorFromObj(list[0], "first element is not callable: "+string(firstEval.Type)), nil
	}
}

func (e *evalCtx) lookupIdentifier(identObj object.Obj, ident object.Identifier) (object.Obj, error) {
	obj, err := e.mem.Get(ident, true)
	if err == nil {
		return obj, nil
	}

	_, found := e.lookupEnvFunction(ident)
	if found {
		return object.Obj{Type: object.OBJ_TYPE_IDENTIFIER, D: ident, Pos: identObj.Pos}, nil
	}

	return e.makeErrorFromObj(identObj, "undefined identifier: "+string(ident)), nil
}

func (e *evalCtx) lookupEnvFunction(ident object.Identifier) (EnvFunction, bool) {
	for _, group := range e.functionGroups {
		functions := group.Functions()
		if function, exists := functions[ident]; exists {
			return function, true
		}
	}
	return EnvFunction{}, false
}

func (e *evalCtx) executeObjectFunction(functionObj object.Obj, args object.List) (object.Obj, error) {
	function := functionObj.D.(object.Function)

	var childMem MEM
	if functionObj.C != nil {
		closureMem := functionObj.C.(MEM)
		childMem = closureMem.Fork()
	} else {
		childMem = e.mem.Fork()
	}

	if function.Variadic {
		return e.executeVariadicFunction(function, args, childMem)
	}

	return e.executeNormalFunction(function, args, childMem)
}

func (e *evalCtx) executeVariadicFunction(function object.Function, args object.List, childMem MEM) (object.Obj, error) {
	evaledArgs := make(object.List, len(args))
	firstArgPos := uint16(0)
	for i, arg := range args {
		if i == 0 && len(args) > 0 {
			firstArgPos = args[0].Pos
		}
		evaledArg, err := e.Evaluate(arg)
		if err != nil {
			return object.Obj{}, err
		}
		if evaledArg.Type == object.OBJ_TYPE_ERROR {
			return evaledArg, nil
		}
		evaledArgs[i] = evaledArg
	}

	argsObj := object.Obj{
		Type: object.OBJ_TYPE_LIST,
		D:    evaledArgs,
		Pos:  firstArgPos,
	}
	childMem.Set("$args", argsObj, false)

	childCtx := &evalCtx{
		mem:             childMem,
		io:              e.io,
		fs:              e.fs,
		functionGroups:  e.functionGroups,
		currentFilePath: e.currentFilePath,
		importedFiles:   e.importedFiles,
	}

	var result object.Obj
	var err error
	for _, instruction := range function.Body {
		result, err = childCtx.Evaluate(instruction)
		if err != nil {
			return object.Obj{}, err
		}
		if result.Type == object.OBJ_TYPE_ERROR {
			return result, nil
		}
	}

	if function.ReturnType != object.OBJ_TYPE_ANY && function.ReturnType != result.Type {
		resultPos := result.Pos
		if resultPos == 0 && len(function.Body) > 0 {
			resultPos = function.Body[len(function.Body)-1].Pos
		}
		return e.makeError(resultPos, fmt.Sprintf("return type mismatch: expected %s, got %s", function.ReturnType, result.Type)), nil
	}

	return result, nil
}

func (e *evalCtx) executeNormalFunction(function object.Function, args object.List, childMem MEM) (object.Obj, error) {
	if len(args) != len(function.Parameters) {
		argPos := uint16(0)
		if len(args) > 0 {
			argPos = args[0].Pos
		}
		return e.makeError(argPos, "wrong number of arguments"), nil
	}

	evaledArgs := make(object.List, len(args))
	for i, arg := range args {
		evaledArg, err := e.Evaluate(arg)
		if err != nil {
			return object.Obj{}, err
		}
		if evaledArg.Type == object.OBJ_TYPE_ERROR {
			return evaledArg, nil
		}
		evaledArgs[i] = evaledArg
	}

	for i, arg := range evaledArgs {
		param := function.Parameters[i]

		if param.Type != object.OBJ_TYPE_ANY && param.Type != arg.Type {
			return e.makeErrorFromObj(arg, fmt.Sprintf("type mismatch for parameter '%s': expected %s, got %s", param.Name, param.Type, arg.Type)), nil
		}

		childMem.Set(param.Name, arg, false)
	}

	childCtx := &evalCtx{
		mem:             childMem,
		io:              e.io,
		fs:              e.fs,
		functionGroups:  e.functionGroups,
		currentFilePath: e.currentFilePath,
		importedFiles:   e.importedFiles,
	}

	var result object.Obj
	var err error
	for _, instruction := range function.Body {
		result, err = childCtx.Evaluate(instruction)
		if err != nil {
			return object.Obj{}, err
		}
		if result.Type == object.OBJ_TYPE_ERROR {
			return result, nil
		}
	}

	if function.ReturnType != object.OBJ_TYPE_ANY && function.ReturnType != result.Type {
		resultPos := result.Pos
		if resultPos == 0 && len(function.Body) > 0 {
			resultPos = function.Body[len(function.Body)-1].Pos
		}
		return e.makeError(resultPos, fmt.Sprintf("return type mismatch: expected %s, got %s", function.ReturnType, result.Type)), nil
	}

	return result, nil
}

func (e *evalCtx) executeEnvFunction(function EnvFunction, args object.List) (object.Obj, error) {
	evaledArgs := args
	if function.EvaluateArgs {
		evaledArgs = make(object.List, len(args))
		for i, arg := range args {
			evaledArg, err := e.Evaluate(arg)
			if err != nil {
				return object.Obj{}, err
			}
			if evaledArg.Type == object.OBJ_TYPE_ERROR {
				return evaledArg, nil
			}
			evaledArgs[i] = evaledArg
		}
	}

	if len(function.Parameters) > 0 {
		if errObj := e.validateEnvArgCount(function, evaledArgs); errObj.Type == object.OBJ_TYPE_ERROR {
			return errObj, nil
		}

		if errObj := e.validateEnvArgTypes(function, evaledArgs); errObj.Type == object.OBJ_TYPE_ERROR {
			return errObj, nil
		}
	}

	result, err := function.Body(e, evaledArgs)
	if err != nil {
		return object.Obj{}, err
	}

	if function.ReturnType != "" && function.ReturnType != object.OBJ_TYPE_ANY {
		if errObj := e.validateEnvReturnType(function, result); errObj.Type == object.OBJ_TYPE_ERROR {
			return errObj, nil
		}
	}

	return result, nil
}

func (e *evalCtx) validateEnvArgCount(fn EnvFunction, args object.List) object.Obj {
	minArgs := len(fn.Parameters)
	argPos := uint16(0)
	if len(args) > 0 {
		argPos = args[0].Pos
	}

	if fn.Variadic {
		if len(args) < minArgs {
			return e.makeError(argPos, fmt.Sprintf("insufficient arguments: expected at least %d, got %d", minArgs, len(args)))
		}
	} else {
		if len(args) != minArgs {
			return e.makeError(argPos, fmt.Sprintf("wrong number of arguments: expected %d, got %d", minArgs, len(args)))
		}
	}

	return object.Obj{Type: object.OBJ_TYPE_NONE, D: object.None{}, Pos: argPos}
}

func (e *evalCtx) validateEnvArgTypes(fn EnvFunction, args object.List) object.Obj {
	argPos := uint16(0)
	if len(args) > 0 {
		argPos = args[0].Pos
	}

	for i, param := range fn.Parameters {
		if i >= len(args) {
			break
		}

		arg := args[i]

		if param.Type != object.OBJ_TYPE_ANY && param.Type != arg.Type {
			return e.makeErrorFromObj(arg, fmt.Sprintf("type mismatch for parameter '%s': expected %s, got %s", param.Name, param.Type, arg.Type))
		}
	}

	if fn.Variadic && len(args) > len(fn.Parameters) && len(fn.Parameters) > 0 {
		lastParam := fn.Parameters[len(fn.Parameters)-1]
		for i := len(fn.Parameters); i < len(args); i++ {
			arg := args[i]
			if lastParam.Type != object.OBJ_TYPE_ANY && lastParam.Type != arg.Type {
				return e.makeErrorFromObj(arg, fmt.Sprintf("type mismatch for variadic parameter '%s' at position %d: expected %s, got %s", lastParam.Name, i, lastParam.Type, arg.Type))
			}
		}
	}

	return object.Obj{Type: object.OBJ_TYPE_NONE, D: object.None{}, Pos: argPos}
}

func (e *evalCtx) validateEnvReturnType(fn EnvFunction, result object.Obj) object.Obj {
	if result.Type == object.OBJ_TYPE_ERROR {
		return result
	}

	if fn.ReturnType != result.Type {
		return e.makeErrorFromObj(result, fmt.Sprintf("return type mismatch: expected %s, got %s", fn.ReturnType, result.Type))
	}

	return object.Obj{Type: object.OBJ_TYPE_NONE, D: object.None{}, Pos: result.Pos}
}

func (e *evalCtx) GetMEM() MEM {
	return e.mem
}

func (e *evalCtx) GetFS() FS {
	return e.fs
}

func (e *evalCtx) GetIO() IO {
	return e.io
}

func (e *evalCtx) GetStartPath() string {
	dir := filepath.Dir(e.currentFilePath)
	return dir
}

func (e *evalCtx) GetRuntime() Runtime {
	return e
}
