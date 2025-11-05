package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/bosley/slpx/pkg/cgs/fs"
	"github.com/bosley/slpx/pkg/cgs/list"
	"github.com/bosley/slpx/pkg/cgs/numbers"
	"github.com/bosley/slpx/pkg/cgs/reflection"
	"github.com/bosley/slpx/pkg/cgs/str"
	"github.com/bosley/slpx/pkg/env"
	"github.com/bosley/slpx/pkg/object"
	"github.com/bosley/slpx/pkg/slp"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [--run] <file1> [file2] [...]\n", os.Args[0])
		os.Exit(1)
	}

	runMode := false
	files := os.Args[1:]

	if len(files) > 0 && files[0] == "--run" {
		runMode = true
		files = files[1:]
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	for _, filePath := range files {
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filePath, err)
			os.Exit(1)
		}

		parser := slp.NewParser(string(content))
		items, err := parser.ParseAll()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing file %s: %v\n", filePath, err)
			os.Exit(1)
		}

		if runMode {
			fmt.Printf("Running: %s\n", filePath)

			fsFunctions := fs.NewFsFunctions(logger)

			evalCtx := env.NewEvalBuilder(logger.WithGroup("eval")).
				WithFunctionGroup(env.NewCoreFunctions()).
				WithFunctionGroup(numbers.NewArithFunctions()).
				WithFunctionGroup(str.NewStrFunctions()).
				WithFunctionGroup(list.NewListFunctions()).
				WithFunctionGroup(reflection.NewReflectionFunctions()).
				WithFunctionGroup(fsFunctions).
				Build()

			absFilePath, err := filepath.Abs(filePath)
			if err != nil {
				absFilePath = filePath
			}

			evalCtx.SetCurrentFilePath(absFilePath)

			// setup after file path set
			fsFunctions.Setup(evalCtx.GetRuntime())

			var result object.Obj = object.Obj{Type: object.OBJ_TYPE_NONE, D: object.None{}}
			for _, item := range items {
				res, err := evalCtx.Evaluate(item)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
					os.Exit(1)
				}
				result = res
			}

			fmt.Printf("Result: %s\n", result.Encode())
		} else {
			fmt.Printf("File: %s\n", filePath)
			for i, item := range items {
				fmt.Printf("[%d] %s\n", i, item.Encode())
				printObjDetails(item, "    ")
				fmt.Println()
			}
		}
	}
}

func printObjDetails(obj object.Obj, indent string) {
	fmt.Printf("%sType: %s\n", indent, obj.Type)
	switch obj.Type {
	case object.OBJ_TYPE_LIST:
		list := obj.D.(object.List)
		fmt.Printf("%sElements: %d\n", indent, len(list))
		for i, item := range list {
			fmt.Printf("%s[%d]:\n", indent, i)
			printObjDetails(item, indent+"    ")
		}
	case object.OBJ_TYPE_INTEGER:
		fmt.Printf("%sValue: %d (integer)\n", indent, obj.D.(object.Integer))
	case object.OBJ_TYPE_REAL:
		fmt.Printf("%sValue: %g (real)\n", indent, obj.D.(object.Real))
	case object.OBJ_TYPE_STRING:
		fmt.Printf("%sValue: %q\n", indent, obj.D.(string))
	case object.OBJ_TYPE_IDENTIFIER:
		fmt.Printf("%sValue: %s\n", indent, obj.D.(object.Identifier))
	case object.OBJ_TYPE_SOME:
		fmt.Printf("%sQuoted:\n", indent)
		printObjDetails(obj.D.(object.Some), indent+"    ")
	case object.OBJ_TYPE_FUNCTION:
		function := obj.D.(object.Function)
		fmt.Printf("%sFunction:\n", indent)
		fmt.Printf("%sParameters: %d\n", indent, len(function.Parameters))
		for i, param := range function.Parameters {
			fmt.Printf("%sParameter[%d]: %s\n", indent, i, param.Name)
		}
		fmt.Printf("%sBody: %d\n", indent, len(function.Body))
		for i, item := range function.Body {
			fmt.Printf("%sBody[%d]:\n", indent, i)
			printObjDetails(item, indent+"    ")
		}
	case object.OBJ_TYPE_ERROR:
		err := obj.D.(object.Error)
		fmt.Printf("%sPosition: %d\n", indent, err.Position)
		fmt.Printf("%sMessage: %s\n", indent, err.Message)
	}
}
