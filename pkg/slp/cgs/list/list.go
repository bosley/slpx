package list

import (
	"strconv"

	"github.com/bosley/slpx/pkg/slp/env"
	"github.com/bosley/slpx/pkg/slp/object"
)

type listFunctions struct{}

func NewListFunctions() env.FunctionGroup {
	return &listFunctions{}
}

func (l *listFunctions) Name() string {
	return "list"
}

func (l *listFunctions) Functions() map[object.Identifier]env.EnvFunction {
	return map[object.Identifier]env.EnvFunction{
		"list/new": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "length", Type: object.OBJ_TYPE_INTEGER},
				{Name: "default", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_LIST,
			Body:       cmdListNew,
		},
		"list/len": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "list", Type: object.OBJ_TYPE_LIST},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdListLen,
		},
		"list/get": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "list", Type: object.OBJ_TYPE_LIST},
				{Name: "index", Type: object.OBJ_TYPE_INTEGER},
			},
			ReturnType: object.OBJ_TYPE_ANY,
			Body:       cmdListGet,
		},
		"list/set": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "list", Type: object.OBJ_TYPE_LIST},
				{Name: "index", Type: object.OBJ_TYPE_INTEGER},
				{Name: "value", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_LIST,
			Body:       cmdListSet,
		},
		"list/push": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "list", Type: object.OBJ_TYPE_LIST},
				{Name: "element", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_LIST,
			Body:       cmdListPush,
		},
		"list/pop": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "list", Type: object.OBJ_TYPE_LIST},
			},
			ReturnType: object.OBJ_TYPE_ANY,
			Body:       cmdListPop,
		},
		"list/clear": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "list", Type: object.OBJ_TYPE_LIST},
			},
			ReturnType: object.OBJ_TYPE_LIST,
			Body:       cmdListClear,
		},
		"list/fill": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "list", Type: object.OBJ_TYPE_LIST},
				{Name: "value", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_LIST,
			Body:       cmdListFill,
		},
		"list/subset": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "list", Type: object.OBJ_TYPE_LIST},
				{Name: "start", Type: object.OBJ_TYPE_INTEGER},
				{Name: "end", Type: object.OBJ_TYPE_INTEGER},
			},
			ReturnType: object.OBJ_TYPE_LIST,
			Body:       cmdListSubset,
		},
		"list/iter": {
			EvaluateArgs: false,
			Parameters: []env.EnvParameter{
				{Name: "list", Type: object.OBJ_TYPE_ANY},
				{Name: "callback", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdListIter,
		},
		"list/contains": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "list", Type: object.OBJ_TYPE_LIST},
				{Name: "element", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdListContains,
		},
		"list/index": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "list", Type: object.OBJ_TYPE_LIST},
				{Name: "element", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdListIndex,
		},
		"list/concat": {
			EvaluateArgs: true,
			Parameters:   []env.EnvParameter{},
			ReturnType:   object.OBJ_TYPE_LIST,
			Variadic:     true,
			Body:         cmdListConcat,
		},
		"list/empty": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "list", Type: object.OBJ_TYPE_LIST},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdListEmpty,
		},
		"list/first": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "list", Type: object.OBJ_TYPE_LIST},
			},
			ReturnType: object.OBJ_TYPE_ANY,
			Body:       cmdListFirst,
		},
		"list/last": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "list", Type: object.OBJ_TYPE_LIST},
			},
			ReturnType: object.OBJ_TYPE_ANY,
			Body:       cmdListLast,
		},
		"list/reverse": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "list", Type: object.OBJ_TYPE_LIST},
			},
			ReturnType: object.OBJ_TYPE_LIST,
			Body:       cmdListReverse,
		},
		"list/join": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "list", Type: object.OBJ_TYPE_LIST},
				{Name: "separator", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_STRING,
			Body:       cmdListJoin,
		},
		"list/slice": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "list", Type: object.OBJ_TYPE_LIST},
				{Name: "start", Type: object.OBJ_TYPE_INTEGER},
				{Name: "end", Type: object.OBJ_TYPE_INTEGER},
			},
			ReturnType: object.OBJ_TYPE_LIST,
			Body:       cmdListSlice,
		},
		"list/map": {
			EvaluateArgs: false,
			Parameters: []env.EnvParameter{
				{Name: "list", Type: object.OBJ_TYPE_ANY},
				{Name: "mapper", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_LIST,
			Body:       cmdListMap,
		},
		"list/filter": {
			EvaluateArgs: false,
			Parameters: []env.EnvParameter{
				{Name: "list", Type: object.OBJ_TYPE_ANY},
				{Name: "predicate", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_LIST,
			Body:       cmdListFilter,
		},
		"list/reduce": {
			EvaluateArgs: false,
			Parameters: []env.EnvParameter{
				{Name: "list", Type: object.OBJ_TYPE_ANY},
				{Name: "initial", Type: object.OBJ_TYPE_ANY},
				{Name: "reducer", Type: object.OBJ_TYPE_ANY},
			},
			ReturnType: object.OBJ_TYPE_ANY,
			Body:       cmdListReduce,
		},
	}
}

func cmdListNew(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	length := int(args[0].D.(object.Integer))
	if length < 0 {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "list/new: length must be non-negative",
			},
		}, nil
	}

	defaultValue := args[1]
	result := make(object.List, length)
	for i := 0; i < length; i++ {
		result[i] = defaultValue.DeepCopy()
	}

	return object.Obj{Type: object.OBJ_TYPE_LIST, D: result}, nil
}

func cmdListLen(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	list := args[0].D.(object.List)
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(len(list))}, nil
}

func cmdListGet(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	list := args[0].D.(object.List)
	index := int(args[1].D.(object.Integer))

	if index < 0 || index >= len(list) {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "list/get: index out of bounds: " + strconv.Itoa(index),
			},
		}, nil
	}

	return list[index], nil
}

func cmdListSet(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	list := args[0].D.(object.List)
	index := int(args[1].D.(object.Integer))
	value := args[2]

	if index < 0 || index >= len(list) {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "list/set: index out of bounds: " + strconv.Itoa(index),
			},
		}, nil
	}

	list[index] = value
	return object.Obj{Type: object.OBJ_TYPE_LIST, D: list}, nil
}

func cmdListPush(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	list := args[0].D.(object.List)
	element := args[1]

	list = append(list, element)
	return object.Obj{Type: object.OBJ_TYPE_LIST, D: list}, nil
}

func cmdListPop(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	list := args[0].D.(object.List)

	if len(list) == 0 {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "list/pop: cannot pop from empty list",
			},
		}, nil
	}

	return list[len(list)-1], nil
}

func cmdListClear(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	return object.Obj{Type: object.OBJ_TYPE_LIST, D: object.List{}}, nil
}

func cmdListFill(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	list := args[0].D.(object.List)
	value := args[1]

	for i := range list {
		list[i] = value.DeepCopy()
	}

	return object.Obj{Type: object.OBJ_TYPE_LIST, D: list}, nil
}

func cmdListSubset(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	list := args[0].D.(object.List)
	start := int(args[1].D.(object.Integer))
	end := int(args[2].D.(object.Integer))

	if start < 0 || start >= len(list) {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "list/subset: start index out of bounds: " + strconv.Itoa(start),
			},
		}, nil
	}

	if end < 0 || end >= len(list) {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "list/subset: end index out of bounds: " + strconv.Itoa(end),
			},
		}, nil
	}

	if start > end {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "list/subset: start index must be <= end index",
			},
		}, nil
	}

	result := make(object.List, end-start+1)
	for i := start; i <= end; i++ {
		result[i-start] = list[i].DeepCopy()
	}

	return object.Obj{Type: object.OBJ_TYPE_LIST, D: result}, nil
}

func cmdListIter(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	listObj, err := ctx.Evaluate(args[0])
	if err != nil {
		return object.Obj{}, err
	}
	if listObj.Type == object.OBJ_TYPE_ERROR {
		return listObj, nil
	}

	if listObj.Type != object.OBJ_TYPE_LIST {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "list/iter: first argument must be a list, got " + string(listObj.Type),
			},
		}, nil
	}

	callbackObj, err := ctx.Evaluate(args[1])
	if err != nil {
		return object.Obj{}, err
	}
	if callbackObj.Type == object.OBJ_TYPE_ERROR {
		return callbackObj, nil
	}

	list := listObj.D.(object.List)

	for _, element := range list {
		callList := object.List{callbackObj, element}
		result, err := ctx.Execute(callList)
		if err != nil {
			return object.Obj{}, err
		}
		if result.Type == object.OBJ_TYPE_ERROR {
			return result, nil
		}

		if result.Type != object.OBJ_TYPE_INTEGER {
			return object.Obj{
				Type: object.OBJ_TYPE_ERROR,
				D: object.Error{
					Position: 0,
					Message:  "list/iter: callback must return integer (1 to continue, 0 to stop)",
				},
			}, nil
		}

		continueIter := result.D.(object.Integer)
		if continueIter == 0 {
			return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
		}
	}

	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
}

func cmdListContains(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	list := args[0].D.(object.List)
	element := args[1]

	elementEncoded := element.Encode()
	for _, item := range list {
		if item.Encode() == elementEncoded {
			return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
		}
	}

	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdListIndex(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	list := args[0].D.(object.List)
	element := args[1]

	elementEncoded := element.Encode()
	for i, item := range list {
		if item.Encode() == elementEncoded {
			return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(i)}, nil
		}
	}

	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(-1)}, nil
}

func cmdListConcat(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	if len(args) == 0 {
		return object.Obj{Type: object.OBJ_TYPE_LIST, D: object.List{}}, nil
	}

	totalLen := 0
	for i, arg := range args {
		if arg.Type != object.OBJ_TYPE_LIST {
			return object.Obj{
				Type: object.OBJ_TYPE_ERROR,
				D: object.Error{
					Position: 0,
					Message:  "list/concat: all arguments must be lists, got " + string(arg.Type) + " at position " + strconv.Itoa(i),
				},
			}, nil
		}
		totalLen += len(arg.D.(object.List))
	}

	result := make(object.List, 0, totalLen)
	for _, arg := range args {
		list := arg.D.(object.List)
		for _, item := range list {
			result = append(result, item.DeepCopy())
		}
	}

	return object.Obj{Type: object.OBJ_TYPE_LIST, D: result}, nil
}

func cmdListEmpty(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	list := args[0].D.(object.List)
	if len(list) == 0 {
		return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(0)}, nil
}

func cmdListFirst(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	list := args[0].D.(object.List)

	if len(list) == 0 {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "list/first: cannot get first element of empty list",
			},
		}, nil
	}

	return list[0], nil
}

func cmdListLast(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	list := args[0].D.(object.List)

	if len(list) == 0 {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "list/last: cannot get last element of empty list",
			},
		}, nil
	}

	return list[len(list)-1], nil
}

func cmdListReverse(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	list := args[0].D.(object.List)

	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}

	return object.Obj{Type: object.OBJ_TYPE_LIST, D: list}, nil
}

func cmdListJoin(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	list := args[0].D.(object.List)
	separator := args[1].D.(string)

	if len(list) == 0 {
		return object.Obj{Type: object.OBJ_TYPE_STRING, D: ""}, nil
	}

	result := ""
	for i, item := range list {
		if i > 0 {
			result += separator
		}
		if item.Type == object.OBJ_TYPE_STRING {
			result += item.D.(string)
		} else {
			result += item.Encode()
		}
	}

	return object.Obj{Type: object.OBJ_TYPE_STRING, D: result}, nil
}

func cmdListSlice(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	list := args[0].D.(object.List)
	start := int(args[1].D.(object.Integer))
	end := int(args[2].D.(object.Integer))

	if start < 0 {
		start = 0
	}
	if end > len(list) {
		end = len(list)
	}
	if start > end {
		start = end
	}

	result := make(object.List, end-start)
	for i := start; i < end; i++ {
		result[i-start] = list[i].DeepCopy()
	}

	return object.Obj{Type: object.OBJ_TYPE_LIST, D: result}, nil
}

func cmdListMap(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	listObj, err := ctx.Evaluate(args[0])
	if err != nil {
		return object.Obj{}, err
	}
	if listObj.Type == object.OBJ_TYPE_ERROR {
		return listObj, nil
	}

	if listObj.Type != object.OBJ_TYPE_LIST {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "list/map: first argument must be a list, got " + string(listObj.Type),
			},
		}, nil
	}

	mapperObj, err := ctx.Evaluate(args[1])
	if err != nil {
		return object.Obj{}, err
	}
	if mapperObj.Type == object.OBJ_TYPE_ERROR {
		return mapperObj, nil
	}

	list := listObj.D.(object.List)
	result := make(object.List, len(list))

	for i, element := range list {
		callList := object.List{mapperObj, element}
		mapped, err := ctx.Execute(callList)
		if err != nil {
			return object.Obj{}, err
		}
		if mapped.Type == object.OBJ_TYPE_ERROR {
			return mapped, nil
		}
		result[i] = mapped
	}

	return object.Obj{Type: object.OBJ_TYPE_LIST, D: result}, nil
}

func cmdListFilter(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	listObj, err := ctx.Evaluate(args[0])
	if err != nil {
		return object.Obj{}, err
	}
	if listObj.Type == object.OBJ_TYPE_ERROR {
		return listObj, nil
	}

	if listObj.Type != object.OBJ_TYPE_LIST {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "list/filter: first argument must be a list, got " + string(listObj.Type),
			},
		}, nil
	}

	predicateObj, err := ctx.Evaluate(args[1])
	if err != nil {
		return object.Obj{}, err
	}
	if predicateObj.Type == object.OBJ_TYPE_ERROR {
		return predicateObj, nil
	}

	list := listObj.D.(object.List)
	result := make(object.List, 0, len(list))

	for _, element := range list {
		callList := object.List{predicateObj, element}
		testResult, err := ctx.Execute(callList)
		if err != nil {
			return object.Obj{}, err
		}
		if testResult.Type == object.OBJ_TYPE_ERROR {
			return testResult, nil
		}

		if testResult.Type != object.OBJ_TYPE_INTEGER {
			return object.Obj{
				Type: object.OBJ_TYPE_ERROR,
				D: object.Error{
					Position: 0,
					Message:  "list/filter: predicate must return integer (1 to include, 0 to exclude)",
				},
			}, nil
		}

		include := testResult.D.(object.Integer)
		if include != 0 {
			result = append(result, element)
		}
	}

	return object.Obj{Type: object.OBJ_TYPE_LIST, D: result}, nil
}

func cmdListReduce(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	listObj, err := ctx.Evaluate(args[0])
	if err != nil {
		return object.Obj{}, err
	}
	if listObj.Type == object.OBJ_TYPE_ERROR {
		return listObj, nil
	}

	if listObj.Type != object.OBJ_TYPE_LIST {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "list/reduce: first argument must be a list, got " + string(listObj.Type),
			},
		}, nil
	}

	initialObj, err := ctx.Evaluate(args[1])
	if err != nil {
		return object.Obj{}, err
	}
	if initialObj.Type == object.OBJ_TYPE_ERROR {
		return initialObj, nil
	}

	reducerObj, err := ctx.Evaluate(args[2])
	if err != nil {
		return object.Obj{}, err
	}
	if reducerObj.Type == object.OBJ_TYPE_ERROR {
		return reducerObj, nil
	}

	list := listObj.D.(object.List)
	accumulator := initialObj

	for _, element := range list {
		callList := object.List{reducerObj, accumulator, element}
		result, err := ctx.Execute(callList)
		if err != nil {
			return object.Obj{}, err
		}
		if result.Type == object.OBJ_TYPE_ERROR {
			return result, nil
		}
		accumulator = result
	}

	return accumulator, nil
}
