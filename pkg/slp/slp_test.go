package slp

import (
	"testing"

	"github.com/bosley/slpx/pkg/object"
)

func TestDeepCopy(t *testing.T) {
	tests := []struct {
		name     string
		original object.Obj
		modify   func(*object.Obj)
		check    func(*testing.T, object.Obj, object.Obj)
	}{
		{
			name:     "OBJ_TYPE_NONE",
			original: object.Obj{Type: object.OBJ_TYPE_NONE, D: object.None{}},
			modify:   func(o *object.Obj) {},
			check: func(t *testing.T, original, copied object.Obj) {
				if copied.Type != object.OBJ_TYPE_NONE {
					t.Errorf("Expected NONE type, got %v", copied.Type)
				}
			},
		},
		{
			name:     "OBJ_TYPE_STRING",
			original: object.Obj{Type: object.OBJ_TYPE_STRING, D: "test string"},
			modify:   func(o *object.Obj) {},
			check: func(t *testing.T, original, copied object.Obj) {
				if copied.D.(string) != "test string" {
					t.Errorf("String copy failed: got %v", copied.D)
				}
			},
		},
		{
			name:     "OBJ_TYPE_INTEGER",
			original: object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(42)},
			modify:   func(o *object.Obj) {},
			check: func(t *testing.T, original, copied object.Obj) {
				if copied.D.(object.Integer) != 42 {
					t.Errorf("Integer copy failed: got %v", copied.D)
				}
			},
		},
		{
			name:     "OBJ_TYPE_REAL",
			original: object.Obj{Type: object.OBJ_TYPE_REAL, D: object.Real(3.14159)},
			modify:   func(o *object.Obj) {},
			check: func(t *testing.T, original, copied object.Obj) {
				if copied.D.(object.Real) != 3.14159 {
					t.Errorf("Real copy failed: got %v", copied.D)
				}
			},
		},
		{
			name:     "OBJ_TYPE_IDENTIFIER",
			original: object.Obj{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("test_id")},
			modify:   func(o *object.Obj) {},
			check: func(t *testing.T, original, copied object.Obj) {
				if copied.D.(object.Identifier) != "test_id" {
					t.Errorf("Identifier copy failed: got %v", copied.D)
				}
			},
		},
		{
			name: "OBJ_TYPE_LIST",
			original: object.Obj{Type: object.OBJ_TYPE_LIST, D: object.List{
				object.Obj{Type: object.OBJ_TYPE_STRING, D: "item1"},
				object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(42)},
			}},
			modify: func(o *object.Obj) {
				list := o.D.(object.List)
				list[0] = object.Obj{Type: object.OBJ_TYPE_STRING, D: "modified"}
			},
			check: func(t *testing.T, original, copied object.Obj) {
				copiedList := copied.D.(object.List)
				if copiedList[0].D.(string) != "item1" {
					t.Error("List deep copy failed: original modification affected copy")
				}
			},
		},
		{
			name:     "OBJ_TYPE_SOME",
			original: object.Obj{Type: object.OBJ_TYPE_SOME, D: object.Obj{Type: object.OBJ_TYPE_STRING, D: "nested"}},
			modify: func(o *object.Obj) {
				some := o.D.(object.Obj)
				some.D = object.Obj{Type: object.OBJ_TYPE_STRING, D: "modified"}
			},
			check: func(t *testing.T, original, copied object.Obj) {
				copiedSome := copied.D.(object.Obj)
				if copiedSome.D.(string) != "nested" {
					t.Error("Some deep copy failed: original modification affected copy")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copied := tt.original.DeepCopy()

			if copied.Type != tt.original.Type {
				t.Errorf("Type mismatch: got %v, want %v", copied.Type, tt.original.Type)
			}

			tt.modify(&tt.original)
			tt.check(t, tt.original, copied)
		})
	}

	t.Run("nested structures isolation", func(t *testing.T) {
		original := object.Obj{Type: object.OBJ_TYPE_LIST, D: object.List{
			object.Obj{Type: object.OBJ_TYPE_STRING, D: "hello"},
			object.Obj{Type: object.OBJ_TYPE_LIST, D: object.List{
				object.Obj{Type: object.OBJ_TYPE_SOME, D: object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(42)}},
			}},
		}}

		copied := original.DeepCopy()

		originalList := original.D.(object.List)
		originalList[0] = object.Obj{Type: object.OBJ_TYPE_STRING, D: "modified"}

		nestedList := originalList[1].D.(object.List)
		nestedList[0] = object.Obj{Type: object.OBJ_TYPE_SOME, D: object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(999)}}

		copiedList := copied.D.(object.List)
		if copiedList[0].D.(string) != "hello" {
			t.Error("Deep copy isolation failed at top level")
		}

		copiedNestedList := copiedList[1].D.(object.List)
		copiedSome := copiedNestedList[0].D.(object.Obj)
		if copiedSome.D.(object.Integer) != 42 {
			t.Error("Deep copy isolation failed at nested level")
		}
	})
}
