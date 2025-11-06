package repl

import (
	"log/slog"
	"path/filepath"

	"github.com/bosley/slpx/pkg/cgs/bits"
	"github.com/bosley/slpx/pkg/cgs/fs"
	"github.com/bosley/slpx/pkg/cgs/io"
	"github.com/bosley/slpx/pkg/cgs/list"
	"github.com/bosley/slpx/pkg/cgs/numbers"
	"github.com/bosley/slpx/pkg/cgs/reflection"
	"github.com/bosley/slpx/pkg/cgs/str"
	"github.com/bosley/slpx/pkg/env"
	"github.com/bosley/slpx/pkg/object"
	"github.com/bosley/slpx/pkg/slp"
)

type sessionEnv struct {
	io      env.IO
	fs      env.FS
	mem     env.MEM
	evalCtx env.EvaluationContext
}

type Session struct {
	logger *slog.Logger
	env    sessionEnv

	pathOnFS string
}

type SessionBuilder struct {
	logger *slog.Logger
	env    sessionEnv

	fgs []env.FunctionGroup
}

func NewSessionBuilder(logger *slog.Logger) *SessionBuilder {
	return &SessionBuilder{
		logger: logger.WithGroup("session"),
		env:    sessionEnv{},
		fgs:    []env.FunctionGroup{},
	}
}

func (b *SessionBuilder) WithIO(io env.IO) *SessionBuilder {
	b.env.io = io
	return b
}

func (b *SessionBuilder) WithFS(fs env.FS) *SessionBuilder {
	b.env.fs = fs
	return b
}

func (b *SessionBuilder) WithMEM(mem env.MEM) *SessionBuilder {
	b.env.mem = mem
	return b
}

func (b *SessionBuilder) WithFunctionGroup(group env.FunctionGroup) *SessionBuilder {
	b.fgs = append(b.fgs, group)
	return b
}

// Path is the "session path" on-disk (in fs) - (likely the path of the main.splx
// file as the user would expect to read/write relative to their launch point)
func (b *SessionBuilder) Build(forPathOnFS string) *Session {
	if b.env.io == nil {
		b.env.io = env.DefaultIO()
	}
	if b.env.fs == nil {
		b.env.fs = env.DefaultFS()
	}
	if b.env.mem == nil {
		b.env.mem = env.DefaultMEM()
	}

	cleanedPathOnFS := filepath.Clean(forPathOnFS)

	session := &Session{
		logger:   b.logger,
		pathOnFS: cleanedPathOnFS,
		env: sessionEnv{
			io:  b.env.io,
			fs:  b.env.fs,
			mem: b.env.mem,
		},
	}

	fsFunctions := fs.NewFsFunctions(b.logger.WithGroup("fs"))
	ioFunctions := io.NewIoFunctions()
	bitsFunctions := bits.NewBitsFunctions()

	session.env.evalCtx = env.NewEvalBuilder(b.logger.WithGroup("eval")).
		WithIO(b.env.io).
		WithFS(b.env.fs).
		WithMEM(b.env.mem).
		WithFunctionGroup(env.NewCoreFunctions()).
		WithFunctionGroup(numbers.NewArithFunctions()).
		WithFunctionGroup(str.NewStrFunctions()).
		WithFunctionGroup(list.NewListFunctions()).
		WithFunctionGroup(reflection.NewReflectionFunctions()).
		WithFunctionGroup(fsFunctions).
		WithFunctionGroup(ioFunctions).
		WithFunctionGroup(bitsFunctions).
		Build()

	session.env.evalCtx.SetCurrentFilePath(
		session.pathOnFS,
	)

	fsFunctions.Setup(session.env.evalCtx.GetRuntime())
	ioFunctions.Setup(session.env.evalCtx.GetRuntime())
	bitsFunctions.Setup(session.env.evalCtx.GetRuntime())

	for _, fg := range b.fgs {
		session.env.evalCtx.AddFunctionGroup(fg)
	}

	return session
}

func (x *Session) Evaluate(source string) (object.Obj, error) {
	parser := slp.NewParser(source)
	items, err := parser.ParseAll()
	if err != nil {
		return object.Obj{}, err
	}

	var result object.Obj = object.Obj{Type: object.OBJ_TYPE_NONE, D: object.None{}}
	for _, item := range items {
		res, err := x.env.evalCtx.Evaluate(item)
		if err != nil {
			return object.Obj{}, err
		}
		if res.Type == object.OBJ_TYPE_ERROR {
			return res, nil
		}
		result = res
	}

	return result, nil
}

func (x *Session) GetIO() env.IO {
	return x.env.io
}

func (x *Session) GetFS() env.FS {
	return x.env.fs
}

func (x *Session) GetMEM() env.MEM {
	return x.env.mem
}
