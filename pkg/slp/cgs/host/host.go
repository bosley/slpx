package host

import (
	"log/slog"
	"os"
	"runtime"

	"github.com/bosley/slpx/pkg/slp/env"
	"github.com/bosley/slpx/pkg/slp/object"
)

type hostFunctions struct {
	logger *slog.Logger
}

func NewHostFunctions(logger *slog.Logger) env.FunctionGroup {
	return &hostFunctions{
		logger: logger,
	}
}

func (h *hostFunctions) Name() string {
	return "host"
}

func (h *hostFunctions) Functions() map[object.Identifier]env.EnvFunction {
	return map[object.Identifier]env.EnvFunction{
		"host/env/get": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "name", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_STRING,
			Body:       cmdEnvGet,
		},
		"host/env/set": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "name", Type: object.OBJ_TYPE_STRING},
				{Name: "value", Type: object.OBJ_TYPE_STRING},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdEnvSet,
		},
		"host/dir/home": {
			EvaluateArgs: true,
			Parameters:   []env.EnvParameter{},
			ReturnType:   object.OBJ_TYPE_STRING,
			Body:         cmdDirHome,
		},
		"host/dir/config": {
			EvaluateArgs: true,
			Parameters:   []env.EnvParameter{},
			ReturnType:   object.OBJ_TYPE_STRING,
			Body:         cmdDirConfig,
		},
		"host/dir/temp": {
			EvaluateArgs: true,
			Parameters:   []env.EnvParameter{},
			ReturnType:   object.OBJ_TYPE_STRING,
			Body:         cmdDirTemp,
		},
		"host/dir/cache": {
			EvaluateArgs: true,
			Parameters:   []env.EnvParameter{},
			ReturnType:   object.OBJ_TYPE_STRING,
			Body:         cmdDirCache,
		},
		"host/os": {
			EvaluateArgs: true,
			Parameters:   []env.EnvParameter{},
			ReturnType:   object.OBJ_TYPE_STRING,
			Body:         cmdOS,
		},
		"host/hw/mem/total": {
			EvaluateArgs: true,
			Parameters:   []env.EnvParameter{},
			ReturnType:   object.OBJ_TYPE_INTEGER,
			Body:         cmdHwMemTotal,
		},
		"host/hw/mem/available": {
			EvaluateArgs: true,
			Parameters:   []env.EnvParameter{},
			ReturnType:   object.OBJ_TYPE_INTEGER,
			Body:         cmdHwMemAvailable,
		},
		"host/hw/mem/used": {
			EvaluateArgs: true,
			Parameters:   []env.EnvParameter{},
			ReturnType:   object.OBJ_TYPE_INTEGER,
			Body:         cmdHwMemUsed,
		},
		"host/hw/mem/percent": {
			EvaluateArgs: true,
			Parameters:   []env.EnvParameter{},
			ReturnType:   object.OBJ_TYPE_REAL,
			Body:         cmdHwMemPercent,
		},
		"host/hw/disk/total": {
			EvaluateArgs: true,
			Parameters:   []env.EnvParameter{},
			ReturnType:   object.OBJ_TYPE_INTEGER,
			Body:         cmdHwDiskTotal,
		},
		"host/hw/disk/used": {
			EvaluateArgs: true,
			Parameters:   []env.EnvParameter{},
			ReturnType:   object.OBJ_TYPE_INTEGER,
			Body:         cmdHwDiskUsed,
		},
		"host/hw/disk/percent": {
			EvaluateArgs: true,
			Parameters:   []env.EnvParameter{},
			ReturnType:   object.OBJ_TYPE_REAL,
			Body:         cmdHwDiskPercent,
		},
		"host/hw/cpu/percent": {
			EvaluateArgs: true,
			Parameters:   []env.EnvParameter{},
			ReturnType:   object.OBJ_TYPE_REAL,
			Body:         cmdHwCpuPercent,
		},
		"host/hw/cpu/count": {
			EvaluateArgs: true,
			Parameters:   []env.EnvParameter{},
			ReturnType:   object.OBJ_TYPE_INTEGER,
			Body:         cmdHwCpuCount,
		},
		"host/hw/cpu/percent/at": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "idx", Type: object.OBJ_TYPE_INTEGER},
			},
			ReturnType: object.OBJ_TYPE_REAL,
			Body:       cmdHwCpuPercentAt,
		},
		"host/hw/cpu/model/at": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "idx", Type: object.OBJ_TYPE_INTEGER},
			},
			ReturnType: object.OBJ_TYPE_STRING,
			Body:       cmdHwCpuModelAt,
		},
		"host/hw/cpu/mhz/at": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "idx", Type: object.OBJ_TYPE_INTEGER},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdHwCpuMhzAt,
		},
		"host/hw/cpu/cache/at": {
			EvaluateArgs: true,
			Parameters: []env.EnvParameter{
				{Name: "idx", Type: object.OBJ_TYPE_INTEGER},
			},
			ReturnType: object.OBJ_TYPE_INTEGER,
			Body:       cmdHwCpuCacheAt,
		},
	}
}

func cmdEnvGet(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	name := args[0].D.(string)
	value, exists := os.LookupEnv(name)
	if !exists {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "host/env/get: environment variable not found: " + name,
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_STRING, D: value}, nil
}

func cmdEnvSet(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	name := args[0].D.(string)
	value := args[1].D.(string)
	err := os.Setenv(name, value)
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "host/env/set: failed to set environment variable: " + err.Error(),
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(1)}, nil
}

func cmdDirHome(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "host/dir/home: failed to get home directory: " + err.Error(),
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_STRING, D: dir}, nil
}

func cmdDirConfig(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "host/dir/config: failed to get config directory: " + err.Error(),
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_STRING, D: dir}, nil
}

func cmdDirTemp(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	dir := os.TempDir()
	return object.Obj{Type: object.OBJ_TYPE_STRING, D: dir}, nil
}

func cmdDirCache(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "host/dir/cache: failed to get cache directory: " + err.Error(),
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_STRING, D: dir}, nil
}

func cmdOS(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	return object.Obj{Type: object.OBJ_TYPE_STRING, D: runtime.GOOS}, nil
}

func cmdHwMemTotal(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	hw, err := GetHardwareProfile()
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "host/hw/mem/total: failed to get hardware profile: " + err.Error(),
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(hw.Memory.Total)}, nil
}

func cmdHwMemAvailable(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	hw, err := GetHardwareProfile()
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "host/hw/mem/available: failed to get hardware profile: " + err.Error(),
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(hw.Memory.Available)}, nil
}

func cmdHwMemUsed(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	hw, err := GetHardwareProfile()
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "host/hw/mem/used: failed to get hardware profile: " + err.Error(),
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(hw.Memory.Used)}, nil
}

func cmdHwMemPercent(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	hw, err := GetHardwareProfile()
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "host/hw/mem/percent: failed to get hardware profile: " + err.Error(),
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_REAL, D: object.Real(hw.Memory.Percent)}, nil
}

func cmdHwDiskTotal(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	hw, err := GetHardwareProfile()
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "host/hw/disk/total: failed to get hardware profile: " + err.Error(),
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(hw.MainDisk.Total)}, nil
}

func cmdHwDiskUsed(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	hw, err := GetHardwareProfile()
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "host/hw/disk/used: failed to get hardware profile: " + err.Error(),
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(hw.MainDisk.Used)}, nil
}

func cmdHwDiskPercent(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	hw, err := GetHardwareProfile()
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "host/hw/disk/percent: failed to get hardware profile: " + err.Error(),
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_REAL, D: object.Real(hw.MainDisk.Percent)}, nil
}

func cmdHwCpuPercent(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	hw, err := GetHardwareProfile()
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "host/hw/cpu/percent: failed to get hardware profile: " + err.Error(),
			},
		}, nil
	}
	if len(hw.CPU) == 0 {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "host/hw/cpu/percent: no CPU information available",
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_REAL, D: object.Real(hw.CPU[0].Percent)}, nil
}

func cmdHwCpuCount(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	hw, err := GetHardwareProfile()
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "host/hw/cpu/count: failed to get hardware profile: " + err.Error(),
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(len(hw.CPU))}, nil
}

func cmdHwCpuPercentAt(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	idx := int(args[0].D.(object.Integer))
	hw, err := GetHardwareProfile()
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "host/hw/cpu/percent/at: failed to get hardware profile: " + err.Error(),
			},
		}, nil
	}
	if idx < 0 || idx >= len(hw.CPU) {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "host/hw/cpu/percent/at: index out of range",
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_REAL, D: object.Real(hw.CPU[idx].Percent)}, nil
}

func cmdHwCpuModelAt(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	idx := int(args[0].D.(object.Integer))
	hw, err := GetHardwareProfile()
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "host/hw/cpu/model/at: failed to get hardware profile: " + err.Error(),
			},
		}, nil
	}
	if idx < 0 || idx >= len(hw.CPU) {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "host/hw/cpu/model/at: index out of range",
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_STRING, D: hw.CPU[idx].ModelName}, nil
}

func cmdHwCpuMhzAt(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	idx := int(args[0].D.(object.Integer))
	hw, err := GetHardwareProfile()
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "host/hw/cpu/mhz/at: failed to get hardware profile: " + err.Error(),
			},
		}, nil
	}
	if idx < 0 || idx >= len(hw.CPU) {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "host/hw/cpu/mhz/at: index out of range",
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(hw.CPU[idx].MHz)}, nil
}

func cmdHwCpuCacheAt(ctx env.EvaluationContext, args object.List) (object.Obj, error) {
	idx := int(args[0].D.(object.Integer))
	hw, err := GetHardwareProfile()
	if err != nil {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "host/hw/cpu/cache/at: failed to get hardware profile: " + err.Error(),
			},
		}, nil
	}
	if idx < 0 || idx >= len(hw.CPU) {
		return object.Obj{
			Type: object.OBJ_TYPE_ERROR,
			D: object.Error{
				Position: 0,
				Message:  "host/hw/cpu/cache/at: index out of range",
			},
		}, nil
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(hw.CPU[idx].CacheSize)}, nil
}
