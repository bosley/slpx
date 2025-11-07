package env

import (
	"os"
	"path/filepath"
)

type fsImpl struct {
	workingDir string
}

var _ FS = &fsImpl{}

func DefaultFS() FS {
	wd, _ := os.Getwd()
	return &fsImpl{
		workingDir: wd,
	}
}

func (f *fsImpl) ReadFile(path string) ([]byte, error) {
	fullPath := f.resolvePath(path)
	return os.ReadFile(fullPath)
}

func (f *fsImpl) WriteFile(path string, data []byte, perm os.FileMode) error {
	fullPath := f.resolvePath(path)
	return os.WriteFile(fullPath, data, perm)
}

func (f *fsImpl) AppendFile(path string, data []byte) error {
	fullPath := f.resolvePath(path)
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	return err
}

func (f *fsImpl) DeleteFile(path string) error {
	fullPath := f.resolvePath(path)
	return os.Remove(fullPath)
}

func (f *fsImpl) RemoveDir(path string) error {
	fullPath := f.resolvePath(path)
	return os.Remove(fullPath)
}

func (f *fsImpl) RemoveDirAll(path string) error {
	fullPath := f.resolvePath(path)
	return os.RemoveAll(fullPath)
}

func (f *fsImpl) Exists(path string) bool {
	fullPath := f.resolvePath(path)
	_, err := os.Stat(fullPath)
	return err == nil
}

func (f *fsImpl) IsDir(path string) bool {
	fullPath := f.resolvePath(path)
	info, err := os.Stat(fullPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func (f *fsImpl) IsFile(path string) bool {
	fullPath := f.resolvePath(path)
	info, err := os.Stat(fullPath)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func (f *fsImpl) ListDir(path string) ([]string, error) {
	fullPath := f.resolvePath(path)
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(entries))
	for i, entry := range entries {
		names[i] = entry.Name()
	}
	return names, nil
}

func (f *fsImpl) MkDir(path string, perm os.FileMode) error {
	fullPath := f.resolvePath(path)
	return os.Mkdir(fullPath, perm)
}

func (f *fsImpl) MkDirAll(path string, perm os.FileMode) error {
	fullPath := f.resolvePath(path)
	return os.MkdirAll(fullPath, perm)
}

func (f *fsImpl) WorkingDir() string {
	return f.workingDir
}

func (f *fsImpl) SetWorkingDir(path string) error {
	fullPath := f.resolvePath(path)
	if !f.IsDir(fullPath) {
		return os.ErrNotExist
	}
	f.workingDir = fullPath
	return nil
}

func (f *fsImpl) resolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(f.workingDir, path)
}
