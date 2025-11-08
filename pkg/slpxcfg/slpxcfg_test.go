package slpxcfg

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bosley/slpx/pkg/slp/object"
)

func TestLoad_AllTypes(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.slpx")

	configContent := `
(set str_var "hello world")
(set int_var 42)
(set real_var 3.14)
(set list_var '(1 2 3))
(set none_var _)
(set func_var (fn () :S "result"))
(set some_var (qu 123))
`

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	variables := []Variable{
		{Identifier: "str_var", Type: object.OBJ_TYPE_STRING, Required: true},
		{Identifier: "int_var", Type: object.OBJ_TYPE_INTEGER, Required: true},
		{Identifier: "real_var", Type: object.OBJ_TYPE_REAL, Required: true},
		{Identifier: "list_var", Type: object.OBJ_TYPE_LIST, Required: true},
		{Identifier: "none_var", Type: object.OBJ_TYPE_NONE, Required: true},
		{Identifier: "func_var", Type: object.OBJ_TYPE_FUNCTION, Required: true},
		{Identifier: "some_var", Type: object.OBJ_TYPE_SOME, Required: true},
	}

	result, err := Load(logger, configFile, variables, 5*time.Second)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(result) != len(variables) {
		t.Errorf("Expected %d variables, got %d", len(variables), len(result))
	}

	if obj, ok := result["str_var"]; !ok {
		t.Error("str_var not found")
	} else if obj.Type != object.OBJ_TYPE_STRING {
		t.Errorf("str_var wrong type: expected %s, got %s", object.OBJ_TYPE_STRING, obj.Type)
	} else if obj.D.(string) != "hello world" {
		t.Errorf("str_var wrong value: expected 'hello world', got '%s'", obj.D.(string))
	}

	if obj, ok := result["int_var"]; !ok {
		t.Error("int_var not found")
	} else if obj.Type != object.OBJ_TYPE_INTEGER {
		t.Errorf("int_var wrong type: expected %s, got %s", object.OBJ_TYPE_INTEGER, obj.Type)
	} else if obj.D.(object.Integer) != 42 {
		t.Errorf("int_var wrong value: expected 42, got %d", obj.D.(object.Integer))
	}

	if obj, ok := result["real_var"]; !ok {
		t.Error("real_var not found")
	} else if obj.Type != object.OBJ_TYPE_REAL {
		t.Errorf("real_var wrong type: expected %s, got %s", object.OBJ_TYPE_REAL, obj.Type)
	}

	if obj, ok := result["list_var"]; !ok {
		t.Error("list_var not found")
	} else if obj.Type != object.OBJ_TYPE_LIST {
		t.Errorf("list_var wrong type: expected %s, got %s", object.OBJ_TYPE_LIST, obj.Type)
	}

	if obj, ok := result["none_var"]; !ok {
		t.Error("none_var not found")
	} else if obj.Type != object.OBJ_TYPE_NONE {
		t.Errorf("none_var wrong type: expected %s, got %s", object.OBJ_TYPE_NONE, obj.Type)
	}

	if obj, ok := result["func_var"]; !ok {
		t.Error("func_var not found")
	} else if obj.Type != object.OBJ_TYPE_FUNCTION {
		t.Errorf("func_var wrong type: expected %s, got %s", object.OBJ_TYPE_FUNCTION, obj.Type)
	}

	if obj, ok := result["some_var"]; !ok {
		t.Error("some_var not found")
	} else if obj.Type != object.OBJ_TYPE_SOME {
		t.Errorf("some_var wrong type: expected %s, got %s", object.OBJ_TYPE_SOME, obj.Type)
	}
}

func TestLoad_RequiredVariableMissing(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.slpx")

	configContent := `
(set var1 "exists")
`

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	variables := []Variable{
		{Identifier: "var1", Type: object.OBJ_TYPE_STRING, Required: true},
		{Identifier: "missing_var", Type: object.OBJ_TYPE_STRING, Required: true},
	}

	_, err := Load(logger, configFile, variables, 5*time.Second)
	if err == nil {
		t.Fatal("Expected error for missing required variable, got nil")
	}

	expectedErrMsg := "required variable 'missing_var' not found in config"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

func TestLoad_OptionalVariableMissing(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.slpx")

	configContent := `
(set var1 "exists")
`

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	variables := []Variable{
		{Identifier: "var1", Type: object.OBJ_TYPE_STRING, Required: true},
		{Identifier: "optional_var", Type: object.OBJ_TYPE_STRING, Required: false},
	}

	result, err := Load(logger, configFile, variables, 5*time.Second)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 variable, got %d", len(result))
	}

	if _, ok := result["var1"]; !ok {
		t.Error("var1 not found")
	}

	if _, ok := result["optional_var"]; ok {
		t.Error("optional_var should not be in result")
	}
}

func TestLoad_TypeMismatch(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.slpx")

	configContent := `
(set var1 "string value")
`

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	variables := []Variable{
		{Identifier: "var1", Type: object.OBJ_TYPE_INTEGER, Required: true},
	}

	_, err := Load(logger, configFile, variables, 5*time.Second)
	if err == nil {
		t.Fatal("Expected error for type mismatch, got nil")
	}

	expectedErrMsg := "type mismatch for variable 'var1': expected integer, got string"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

func TestLoad_AnyTypeAcceptsAll(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.slpx")

	configContent := `
(set str_var "hello")
(set int_var 123)
(set real_var 4.56)
`

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	variables := []Variable{
		{Identifier: "str_var", Type: object.OBJ_TYPE_ANY, Required: true},
		{Identifier: "int_var", Type: object.OBJ_TYPE_ANY, Required: true},
		{Identifier: "real_var", Type: object.OBJ_TYPE_ANY, Required: true},
	}

	result, err := Load(logger, configFile, variables, 5*time.Second)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 variables, got %d", len(result))
	}

	if obj, ok := result["str_var"]; !ok {
		t.Error("str_var not found")
	} else if obj.Type != object.OBJ_TYPE_STRING {
		t.Errorf("str_var wrong type: got %s", obj.Type)
	}

	if obj, ok := result["int_var"]; !ok {
		t.Error("int_var not found")
	} else if obj.Type != object.OBJ_TYPE_INTEGER {
		t.Errorf("int_var wrong type: got %s", obj.Type)
	}

	if obj, ok := result["real_var"]; !ok {
		t.Error("real_var not found")
	} else if obj.Type != object.OBJ_TYPE_REAL {
		t.Errorf("real_var wrong type: got %s", obj.Type)
	}
}

func TestLoad_ParseError(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.slpx")

	configContent := `
(set var1 "unclosed string
`

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	variables := []Variable{
		{Identifier: "var1", Type: object.OBJ_TYPE_STRING, Required: true},
	}

	_, err := Load(logger, configFile, variables, 5*time.Second)
	if err == nil {
		t.Fatal("Expected parse error, got nil")
	}
}

func TestLoad_EvaluationError(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.slpx")

	configContent := `
(undefined_function "arg")
`

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	variables := []Variable{}

	_, err := Load(logger, configFile, variables, 5*time.Second)
	if err == nil {
		t.Fatal("Expected evaluation error, got nil")
	}
}

func TestLoad_EmptyVariableList(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.slpx")

	configContent := `
(set var1 "value1")
(set var2 42)
`

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	variables := []Variable{}

	result, err := Load(logger, configFile, variables, 5*time.Second)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected empty result map, got %d items", len(result))
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	variables := []Variable{
		{Identifier: "var1", Type: object.OBJ_TYPE_STRING, Required: true},
	}

	_, err := Load(logger, "/nonexistent/path/config.slpx", variables, 5*time.Second)
	if err == nil {
		t.Fatal("Expected file not found error, got nil")
	}
}

func TestLoad_MixedRequiredOptional(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.slpx")

	configContent := `
(set required1 "value1")
(set optional1 "value2")
(set required2 100)
`

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	variables := []Variable{
		{Identifier: "required1", Type: object.OBJ_TYPE_STRING, Required: true},
		{Identifier: "required2", Type: object.OBJ_TYPE_INTEGER, Required: true},
		{Identifier: "optional1", Type: object.OBJ_TYPE_STRING, Required: false},
		{Identifier: "optional2", Type: object.OBJ_TYPE_STRING, Required: false},
	}

	result, err := Load(logger, configFile, variables, 5*time.Second)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 variables, got %d", len(result))
	}

	if _, ok := result["required1"]; !ok {
		t.Error("required1 not found")
	}
	if _, ok := result["required2"]; !ok {
		t.Error("required2 not found")
	}
	if _, ok := result["optional1"]; !ok {
		t.Error("optional1 not found")
	}
	if _, ok := result["optional2"]; ok {
		t.Error("optional2 should not be in result")
	}
}

func TestLoad_ComplexListsAndFunctions(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.slpx")

	configContent := `
(set nested_list '((1 2) (3 4) (5 6)))
(set func_with_logic (fn (x :I) :I (int/add x 10)))
(set mixed_list '(1 "two" 3.0 _ (nested)))
`

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	variables := []Variable{
		{Identifier: "nested_list", Type: object.OBJ_TYPE_LIST, Required: true},
		{Identifier: "func_with_logic", Type: object.OBJ_TYPE_FUNCTION, Required: true},
		{Identifier: "mixed_list", Type: object.OBJ_TYPE_LIST, Required: true},
	}

	result, err := Load(logger, configFile, variables, 5*time.Second)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 variables, got %d", len(result))
	}

	if obj, ok := result["nested_list"]; !ok {
		t.Error("nested_list not found")
	} else if obj.Type != object.OBJ_TYPE_LIST {
		t.Errorf("nested_list wrong type: expected %s, got %s", object.OBJ_TYPE_LIST, obj.Type)
	}

	if obj, ok := result["func_with_logic"]; !ok {
		t.Error("func_with_logic not found")
	} else if obj.Type != object.OBJ_TYPE_FUNCTION {
		t.Errorf("func_with_logic wrong type: expected %s, got %s", object.OBJ_TYPE_FUNCTION, obj.Type)
	}

	if obj, ok := result["mixed_list"]; !ok {
		t.Error("mixed_list not found")
	} else if obj.Type != object.OBJ_TYPE_LIST {
		t.Errorf("mixed_list wrong type: expected %s, got %s", object.OBJ_TYPE_LIST, obj.Type)
	}
}
