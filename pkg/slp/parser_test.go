package slp

import (
	"fmt"
	"testing"

	"github.com/bosley/slpx/pkg/object"
)

func TestParser(t *testing.T) {
	testCases := []struct {
		input    string
		expected object.Obj
	}{
		{
			input: "(set key value)",
			expected: object.Obj{
				Type: object.OBJ_TYPE_LIST,
				D: object.List{
					{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("set")},
					{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("key")},
					{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("value")},
				},
			},
		},
		{
			input: "(get key)",
			expected: object.Obj{
				Type: object.OBJ_TYPE_LIST,
				D: object.List{
					{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("get")},
					{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("key")},
				},
			},
		},
		{
			input: "_",
			expected: object.Obj{
				Type: object.OBJ_TYPE_NONE,
				D:    object.None{},
			},
		},
		{
			input: `(set "key with spaces" "value with spaces")`,
			expected: object.Obj{
				Type: object.OBJ_TYPE_LIST,
				D: object.List{
					{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("set")},
					{Type: object.OBJ_TYPE_STRING, D: "key with spaces"},
					{Type: object.OBJ_TYPE_STRING, D: "value with spaces"},
				},
			},
		},
		{
			input: `(set "key\nwith\ttabs" "value\r\nwith\nlines")`,
			expected: object.Obj{
				Type: object.OBJ_TYPE_LIST,
				D: object.List{
					{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("set")},
					{Type: object.OBJ_TYPE_STRING, D: "key\nwith\ttabs"},
					{Type: object.OBJ_TYPE_STRING, D: "value\r\nwith\nlines"},
				},
			},
		},
		{
			input: `(set "quote \"inside\"" "backslash \\ here")`,
			expected: object.Obj{
				Type: object.OBJ_TYPE_LIST,
				D: object.List{
					{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("set")},
					{Type: object.OBJ_TYPE_STRING, D: `quote "inside"`},
					{Type: object.OBJ_TYPE_STRING, D: `backslash \ here`},
				},
			},
		},
		{
			input: `'42`,
			expected: object.Obj{
				Type: object.OBJ_TYPE_SOME,
				D:    object.Some(object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(42)}),
			},
		},
		{
			input: `'(set key value)`,
			expected: object.Obj{
				Type: object.OBJ_TYPE_SOME,
				D: object.Some(object.Obj{
					Type: object.OBJ_TYPE_LIST,
					D: object.List{
						{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("set")},
						{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("key")},
						{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("value")},
					},
				}),
			},
		},
		{
			input: `'hello`,
			expected: object.Obj{
				Type: object.OBJ_TYPE_SOME,
				D:    object.Some(object.Obj{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("hello")}),
			},
		},

		{
			input: `(set (get key) (nested (deep value)))`,
			expected: object.Obj{
				Type: object.OBJ_TYPE_LIST,
				D: object.List{
					{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("set")},
					{Type: object.OBJ_TYPE_LIST, D: object.List{
						{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("get")},
						{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("key")},
					}},
					{Type: object.OBJ_TYPE_LIST, D: object.List{
						{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("nested")},
						{Type: object.OBJ_TYPE_LIST, D: object.List{
							{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("deep")},
							{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("value")},
						}},
					}},
				},
			},
		},
		{
			input: `(unknown_keyword arg1 arg2)`,
			expected: object.Obj{
				Type: object.OBJ_TYPE_LIST,
				D: object.List{
					{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("unknown_keyword")},
					{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("arg1")},
					{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("arg2")},
				},
			},
		},
		{
			input: `invalid_keyword`,
			expected: object.Obj{
				Type: object.OBJ_TYPE_IDENTIFIER,
				D:    object.Identifier("invalid_keyword"),
			},
		},
		{
			input: `(set key 42)`,
			expected: object.Obj{
				Type: object.OBJ_TYPE_LIST,
				D: object.List{
					{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("set")},
					{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("key")},
					{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(42)},
				},
			},
		},
		{
			input: `3.14`,
			expected: object.Obj{
				Type: object.OBJ_TYPE_REAL,
				D:    object.Real(3.14),
			},
		},
		{
			input: `-123`,
			expected: object.Obj{
				Type: object.OBJ_TYPE_INTEGER,
				D:    object.Integer(-123)},
		},
		{
			input: `+`,
			expected: object.Obj{
				Type: object.OBJ_TYPE_IDENTIFIER,
				D:    object.Identifier("+")},
		},
		{
			input: `-`,
			expected: object.Obj{
				Type: object.OBJ_TYPE_IDENTIFIER,
				D:    object.Identifier("-")},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			parser := &Parser{
				Target:   tc.input,
				Position: 0,
			}

			// TODO: Implement Parse method
			result, err := parser.Parse()
			if err != nil {
				t.Errorf("error parsing: %v", err)
				return
			}

			if result.Type != tc.expected.Type {
				t.Errorf("expected type %s, got %s", tc.expected.Type, result.Type)
			}

			if result.Type == object.OBJ_TYPE_INTEGER {
				expectedInt, ok := tc.expected.D.(object.Integer)
				if !ok {
					t.Errorf("expected integer but got non-integer expected type")
					return
				}
				actualInt, ok := result.D.(object.Integer)
				if !ok {
					t.Errorf("got integer type but result.D is not Integer")
					return
				}
				if actualInt != expectedInt {
					t.Errorf("expected integer %d, got %d", expectedInt, actualInt)
				}
			}

			if result.Type == object.OBJ_TYPE_REAL {
				expectedReal, ok := tc.expected.D.(object.Real)
				if !ok {
					t.Errorf("expected real but got non-real expected type")
					return
				}
				actualReal, ok := result.D.(object.Real)
				if !ok {
					t.Errorf("got real type but result.D is not Real")
					return
				}
				if !realsEqual(float64(actualReal), float64(expectedReal)) {
					t.Errorf("expected real %f, got %f", expectedReal, actualReal)
				}
			}

			// For identifier types, check string value
			if result.Type == object.OBJ_TYPE_IDENTIFIER {
				expectedId, ok := tc.expected.D.(object.Identifier)
				if !ok {
					t.Errorf("expected identifier but got non-identifier expected type")
					return
				}
				actualId, ok := result.D.(object.Identifier)
				if !ok {
					t.Errorf("got identifier type but result.D is not Identifier")
					return
				}
				if string(actualId) != string(expectedId) {
					t.Errorf("expected identifier %s, got %s", expectedId, actualId)
				}
			}

			// For some types (quoted expressions), check the contained object
			if result.Type == object.OBJ_TYPE_SOME {
				expectedSome, ok := tc.expected.D.(object.Some)
				if !ok {
					t.Errorf("expected some but got non-some expected type")
					return
				}
				actualSome, ok := result.D.(object.Some)
				if !ok {
					t.Errorf("got some type but result.D is not Some")
					return
				}
				// Compare the contained objects
				if expectedSome.Type != actualSome.Type {
					t.Errorf("expected some containing %s, got some containing %s", expectedSome.Type, actualSome.Type)
				}
				// For now, just check the type matches - could add deeper comparison if needed
			}

			// For now, just print the result to see what we get
			fmt.Printf("Input: %s\nResult: %+v\n\n", tc.input, result)
		})
	}
}

func realsEqual(a, b float64) bool {
	const epsilon = 1e-10
	return a-b < epsilon && b-a < epsilon
}

func TestRoundTrip(t *testing.T) {
	testCases := []string{
		`(set key "hello world")`,
		`(nested (list (with "quoted strings" and numbers 42 3.14 -123)))`,
		`((lambda (x) (add x 1)) 5)`,
		`_`,
		`hello`,
		`42`,
		`3.14`,
		`-123`,
		`"string with \"quotes\" and\ttabs\nand newlines"`,
		`'42`,
		`'(set key value)`,
		`'(nested (quoted list))`,
		`+`,
		`-`,
		`<=`,
		`my_var_123`,
		`((deeply (nested (structure (with (many (levels (of (parentheses)))))))))`,
	}

	for i, input := range testCases {
		t.Run(fmt.Sprintf("roundtrip_%d", i), func(t *testing.T) {
			// Parse the input
			parser := &Parser{Target: input, Position: 0}
			original, err := parser.Parse()
			if err != nil {
				t.Errorf("error parsing: %v", err)
				return
			}

			// Encode it back to string
			encoded := original.Encode()

			// Parse the encoded string again
			parser2 := &Parser{Target: encoded, Position: 0}
			decoded, err := parser2.Parse()
			if err != nil {
				t.Errorf("error parsing: %v", err)
				return
			}

			// Verify they match
			if !objectsEqual(original, decoded) {
				t.Errorf("Round-trip failed for input: %s", input)
				t.Errorf("Original: %+v", original)
				t.Errorf("Encoded: %s", encoded)
				t.Errorf("Decoded: %+v", decoded)
			}
		})
	}
}

func objectsEqual(a, b object.Obj) bool {
	if a.Type != b.Type {
		return false
	}

	switch a.Type {
	case object.OBJ_TYPE_NONE:
		return true
	case object.OBJ_TYPE_STRING:
		return a.D.(string) == b.D.(string)
	case object.OBJ_TYPE_INTEGER:
		return a.D.(object.Integer) == b.D.(object.Integer)
	case object.OBJ_TYPE_REAL:
		return realsEqual(float64(a.D.(object.Real)), float64(b.D.(object.Real)))
	case object.OBJ_TYPE_IDENTIFIER:
		return string(a.D.(object.Identifier)) == string(b.D.(object.Identifier))
	case object.OBJ_TYPE_LIST:
		aList := a.D.(object.List)
		bList := b.D.(object.List)
		if len(aList) != len(bList) {
			return false
		}
		for i := range aList {
			if !objectsEqual(aList[i], bList[i]) {
				return false
			}
		}
		return true
	case object.OBJ_TYPE_SOME:
		return objectsEqual(object.Obj(a.D.(object.Some)), object.Obj(b.D.(object.Some)))

	default:
		return false
	}
}

func TestErrorLiteral(t *testing.T) {
	testCases := []struct {
		input       string
		expectError bool
		checkResult func(*testing.T, object.Obj)
	}{
		{
			input:       `@(Something went wrong)`,
			expectError: false,
			checkResult: func(t *testing.T, obj object.Obj) {
				if obj.Type != object.OBJ_TYPE_ERROR {
					t.Errorf("expected ERROR type, got %v", obj.Type)
				}
				errData := obj.D.(object.Error)
				if errData.Message != "Something went wrong" {
					t.Errorf("expected message 'Something went wrong', got %q", errData.Message)
				}
				if errData.Position != 0 {
					t.Errorf("expected position 0, got %d", errData.Position)
				}
			},
		},
		{
			input:       `@(Division by zero at 42)`,
			expectError: false,
			checkResult: func(t *testing.T, obj object.Obj) {
				if obj.Type != object.OBJ_TYPE_ERROR {
					t.Errorf("expected ERROR type, got %v", obj.Type)
				}
				errData := obj.D.(object.Error)
				if errData.Message != "Division by zero at 42" {
					t.Errorf("expected specific message, got %q", errData.Message)
				}
			},
		},
		{
			input:       `@(nested (error (message)))`,
			expectError: false,
			checkResult: func(t *testing.T, obj object.Obj) {
				if obj.Type != object.OBJ_TYPE_ERROR {
					t.Errorf("expected ERROR type, got %v", obj.Type)
				}
				errData := obj.D.(object.Error)
				if errData.Message != "nested (error (message))" {
					t.Errorf("expected nested list message, got %q", errData.Message)
				}
			},
		},
		{
			input:       `@`,
			expectError: true,
			checkResult: nil,
		},
		{
			input:       `@ `,
			expectError: true,
			checkResult: nil,
		},
		{
			input:       `@hello`,
			expectError: true,
			checkResult: nil,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("error_literal_%d", i), func(t *testing.T) {
			parser := &Parser{Target: tc.input, Position: 0}
			result, err := parser.Parse()

			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.checkResult != nil {
				tc.checkResult(t, result)
			}
		})
	}
}

func TestMacros(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []object.Obj
	}{
		{
			name:  "simple doubling macro",
			input: `$(twice ?x) (list ?x ?x) ($twice "hello")`,
			expected: []object.Obj{
				{
					Type: object.OBJ_TYPE_LIST,
					D: object.List{
						{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("list")},
						{Type: object.OBJ_TYPE_STRING, D: "hello"},
						{Type: object.OBJ_TYPE_STRING, D: "hello"},
					},
				},
			},
		},
		{
			name:  "when macro",
			input: `$(when ?cond ?body) (if ?cond ?body _) ($when test action)`,
			expected: []object.Obj{
				{
					Type: object.OBJ_TYPE_LIST,
					D: object.List{
						{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("if")},
						{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("test")},
						{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("action")},
						{Type: object.OBJ_TYPE_NONE, D: object.None{}},
					},
				},
			},
		},
		{
			name:  "nested macro expansion",
			input: `$(double ?x) (add ?x ?x) $(quad ?x) ($double ($double ?x)) ($quad 5)`,
			expected: []object.Obj{
				{
					Type: object.OBJ_TYPE_LIST,
					D: object.List{
						{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("add")},
						{
							Type: object.OBJ_TYPE_LIST,
							D: object.List{
								{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("add")},
								{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(5)},
								{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(5)},
							},
						},
						{
							Type: object.OBJ_TYPE_LIST,
							D: object.List{
								{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("add")},
								{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(5)},
								{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(5)},
							},
						},
					},
				},
			},
		},
		{
			name:  "macro with quoted parameter",
			input: `$(defconst ?name ?value) (def ?name '?value) ($defconst PI 3.14)`,
			expected: []object.Obj{
				{
					Type: object.OBJ_TYPE_LIST,
					D: object.List{
						{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("def")},
						{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("PI")},
						{Type: object.OBJ_TYPE_SOME, D: object.Some(object.Obj{Type: object.OBJ_TYPE_REAL, D: object.Real(3.14)})},
					},
				},
			},
		},
		{
			name:  "non-macro list unchanged",
			input: `(regular function call)`,
			expected: []object.Obj{
				{
					Type: object.OBJ_TYPE_LIST,
					D: object.List{
						{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("regular")},
						{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("function")},
						{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("call")},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewParser(tc.input)
			results, err := parser.ParseAll()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(results) != len(tc.expected) {
				t.Fatalf("expected %d results, got %d", len(tc.expected), len(results))
			}

			for i, result := range results {
				if !objectsEqual(result, tc.expected[i]) {
					t.Errorf("result %d mismatch:\nexpected: %+v\ngot: %+v\nexpected encoded: %s\ngot encoded: %s",
						i, tc.expected[i], result, tc.expected[i].Encode(), result.Encode())
				}
			}
		})
	}
}

func TestMacroErrors(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "macro with wrong arg count",
			input:       `$(test ?x ?y) (add ?x ?y) ($test 1)`,
			expectError: true,
		},
		{
			name:        "undefined macro",
			input:       `($undefined 1 2)`,
			expectError: true,
		},
		{
			name:        "macro parameter without ?",
			input:       `$(test x) (add x 1)`,
			expectError: true,
		},
		{
			name:        "empty macro pattern",
			input:       `$() (add 1 1)`,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewParser(tc.input)
			_, err := parser.ParseAll()
			if tc.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	testCases := []struct {
		obj      object.Obj
		expected string
	}{
		{
			obj:      object.Obj{Type: object.OBJ_TYPE_NONE, D: object.None{}},
			expected: "_",
		},
		{
			obj: object.Obj{
				Type: object.OBJ_TYPE_LIST,
				D: object.List{
					{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("set")},
					{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("key")},
					{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("value")},
				},
			},
			expected: "(set key value)",
		},
		{
			obj: object.Obj{
				Type: object.OBJ_TYPE_LIST,
				D:    object.List{},
			},
			expected: "()",
		},
		{
			obj:      object.Obj{Type: object.OBJ_TYPE_STRING, D: "hello world"},
			expected: "\"hello world\"",
		},
		{
			obj:      object.Obj{Type: object.OBJ_TYPE_STRING, D: "quote \"inside\""},
			expected: "\"quote \\\"inside\\\"\"",
		},
		{
			obj:      object.Obj{Type: object.OBJ_TYPE_STRING, D: "with\nnewlines"},
			expected: "\"with\\nnewlines\"",
		},
		{
			obj:      object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(42)},
			expected: "42",
		},
		{
			obj:      object.Obj{Type: object.OBJ_TYPE_REAL, D: object.Real(3.14)},
			expected: "3.14",
		},
		{
			obj:      object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(-123)},
			expected: "-123",
		},
		{
			obj:      object.Obj{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("hello")},
			expected: "hello",
		},
		{
			obj: object.Obj{
				Type: object.OBJ_TYPE_SOME,
				D:    object.Some(object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(42)}),
			},
			expected: "'42",
		},
		{
			obj: object.Obj{
				Type: object.OBJ_TYPE_SOME,
				D: object.Some(object.Obj{
					Type: object.OBJ_TYPE_LIST,
					D: object.List{
						{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("set")},
						{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier("key")},
					},
				}),
			},
			expected: "'(set key)",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("encode_%d", i), func(t *testing.T) {
			result := tc.obj.Encode()
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}

		})
	}
}
