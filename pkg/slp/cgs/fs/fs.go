package fs

import (
	"log/slog"

	"github.com/bosley/slpx/pkg/slp/env"
	"github.com/bosley/slpx/pkg/slp/object"
)

type fsFunctions struct {
	logger     *slog.Logger
	fs         env.FS
	workingDir string
}

func NewFsFunctions(logger *slog.Logger) *fsFunctions {
	return &fsFunctions{
		logger: logger,
	}
}

func (f *fsFunctions) Setup(runtime env.Runtime) {
	f.fs = runtime.GetFS()
	f.workingDir = runtime.GetStartPath()
}

func (f *fsFunctions) Name() string {
	return "fs"
}

func (f *fsFunctions) Functions() map[object.Identifier]env.EnvFunction {
	return map[object.Identifier]env.EnvFunction{
		"fs/exists?": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "path", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       f.cmdExists,
		},
		"fs/dir?": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "path", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       f.cmdIsDir,
		},
		"fs/file?": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "path", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       f.cmdIsFile,
		},
		"fs/read_file": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "path", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_STRING,
			Body:       f.cmdReadFile,
		},
		"fs/write_file": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "path", Type: object.OBJ_TYPE_STRING},
				{Name: "data", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       f.cmdWriteFile,
		},
		"fs/append_file": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "path", Type: object.OBJ_TYPE_STRING},
				{Name: "data", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       f.cmdAppendFile,
		},
		"fs/rm_file": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "path", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       f.cmdRemoveFile,
		},
		"fs/rm_dir": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "path", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       f.cmdRemoveDir,
		},
		"fs/rm_dir_all": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "path", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       f.cmdRemoveDirAll,
		},
		"fs/mk_dir": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "path", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       f.cmdMkDir,
		},
		"fs/mk_dir_all": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "path", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       f.cmdMkDirAll,
		},
		"fs/list_dir": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "path", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_LIST,
			Body:       f.cmdListDir,
		},
		"fs/working_dir": {
			EvaluateArgs: true,
			Parameters:   []env.EnvParameter{},
			ReturnType:   object.OBJ_TYPE_STRING,
			Body:         f.cmdWorkingDir,
		},
		"fs/set_working_dir": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "path", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       f.cmdSetWorkingDir,
		},
	}
}

func (f *fsFunctions) cmdExists(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	path := args[0].D.(string)
	if f.fs.Exists(path) {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func (f *fsFunctions) cmdIsDir(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	path := args[0].D.(string)
	if f.fs.IsDir(path) {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func (f *fsFunctions) cmdIsFile(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	path := args[0].D.(string)
	if f.fs.IsFile(path) {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func (f *fsFunctions) cmdReadFile(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	path := args[0].D.(string)
	data, err := f.fs.ReadFile(path)
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "failed to read file: " + err.Error(),
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_STRING, D: string(data)}, nil
}

func (f *fsFunctions) cmdWriteFile(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	path := args[0].D.(string)
	data := args[1].D.(string)
	err := f.fs.WriteFile(path, []byte(data), 0644)
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "failed to write file: " + err.Error(),
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
}

func (f *fsFunctions) cmdAppendFile(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	path := args[0].D.(string)
	data := args[1].D.(string)
	err := f.fs.AppendFile(path, []byte(data))
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "failed to append to file: " + err.Error(),
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
}

func (f *fsFunctions) cmdRemoveFile(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	path := args[0].D.(string)
	err := f.fs.DeleteFile(path)
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "failed to remove file: " + err.Error(),
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
}

func (f *fsFunctions) cmdRemoveDir(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	path := args[0].D.(string)
	err := f.fs.RemoveDir(path)
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "failed to remove directory: " + err.Error(),
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
}

func (f *fsFunctions) cmdRemoveDirAll(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	path := args[0].D.(string)
	err := f.fs.RemoveDirAll(path)
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "failed to remove directory tree: " + err.Error(),
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
}

func (f *fsFunctions) cmdMkDir(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	path := args[0].D.(string)
	err := f.fs.MkDir(path, 0755)
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "failed to create directory: " + err.Error(),
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
}

func (f *fsFunctions) cmdMkDirAll(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	path := args[0].D.(string)
	err := f.fs.MkDirAll(path, 0755)
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "failed to create directory tree: " + err.Error(),
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
}

func (f *fsFunctions) cmdListDir(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	path := args[0].D.(string)
	entries, err := f.fs.ListDir(path)
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "failed to list directory: " + err.Error(),
			},
		}, nil
	}
	list := make(object.List, len(entries))
	for i, entry := range entries {
		list[i] = object.Obj{Type: object.OBJ_TYPE_STRING, D: entry}
	}
	return object.Obj{Type: object.OBJ_TYPE_LIST, D: list}, nil
}

func (f *fsFunctions) cmdWorkingDir(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	return object.Obj{Type: object.OBJ_TYPE_STRING, D: f.fs.WorkingDir()}, nil
}

func (f *fsFunctions) cmdSetWorkingDir(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	path := args[0].D.(string)
	err := f.fs.SetWorkingDir(path)
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "failed to set working directory: " + err.Error(),
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
}
