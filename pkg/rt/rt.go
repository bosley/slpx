package rt

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/bosley/slpx/pkg/slp/env"
	"github.com/bosley/slpx/pkg/slp/object"
	"github.com/bosley/slpx/pkg/slp/repl"
	"github.com/bosley/slpx/pkg/slpxcfg"
	"github.com/google/uuid"
)

type ActiveContext interface {
	ID() string
	DisplayName() string
	GetRepl() *repl.Session
	GetTuiConfig() TuiConfig
	Close() error
}

type Runtime interface {
	SLPXHome() string

	/*
		Create a new context that can be used to get an isolated repl and
		be used as a primary interface into the runtime.

		The runtime can hand out an arbitrary number of sessions
	*/
	NewActiveContext(name string) (ActiveContext, error)

	Stop() error
}

type TuiConfig struct {
	CmdToggleEditor      string
	CmdToggleOutput      string
	CmdClear             string
	PromptColor          string
	ResultColor          string
	ErrorColor           string
	HelpColor            string
	FocusedBorderColor   string
	BlurredBorderColor   string
	SelectedItemColor    string
	HistoryItemColor     string
	DirtyPromptColor     string
	SecondaryActionColor string
	CommandRouter        object.Function
}

type activeContext struct {
	id          string
	displayName string
	env         env.EvaluationContext
	fs          env.FS
	io          env.IO
	mem         env.MEM

	repl *repl.Session

	onClose func() error

	tuiConfig TuiConfig
}

func (x *activeContext) ID() string {
	return x.id
}

func (x *activeContext) DisplayName() string {
	return x.displayName
}

func (x *activeContext) GetRepl() *repl.Session {
	return x.repl
}

func (x *activeContext) GetTuiConfig() TuiConfig {
	return x.tuiConfig
}

func (x *activeContext) Close() error {
	return x.onClose()
}

type runtimeImpl struct {
	logger          *slog.Logger
	slpxHome        string
	launchDirectory string
	setupContent    string

	activeContexts map[string]activeContext
	acMutex        sync.Mutex
}

type Config struct {
	Logger          *slog.Logger
	SLPXHome        string
	LaunchDirectory string
	SetupContent    string
}

func New(config Config) (Runtime, error) {

	return &runtimeImpl{
		logger:          config.Logger,
		slpxHome:        config.SLPXHome,
		launchDirectory: config.LaunchDirectory,
		setupContent:    config.SetupContent,
		activeContexts:  make(map[string]activeContext),
		acMutex:         sync.Mutex{},
	}, nil
}

func (r *runtimeImpl) SLPXHome() string {
	return r.slpxHome
}

func (r *runtimeImpl) onActiveContextClose(id string) {
	r.acMutex.Lock()
	defer r.acMutex.Unlock()
	delete(r.activeContexts, id)
}

func (r *runtimeImpl) getFsForNewActiveContext() env.FS {
	return env.DefaultFS()
}
func (r *runtimeImpl) getIoForNewActiveContext() env.IO {
	return env.DefaultIO()
}
func (r *runtimeImpl) getMemForNewActiveContext() env.MEM {
	return env.DefaultMEM()
}

func (r *runtimeImpl) getEvalBuilderForNewActiveContext(id string) env.EvaluationContext {
	return env.NewEvalBuilder(r.logger.WithGroup("ac:" + id)).Build()
}

func (r *runtimeImpl) NewActiveContext(displayName string) (ActiveContext, error) {

	uuid := uuid.New().String()

	fs := r.getFsForNewActiveContext()
	io := r.getIoForNewActiveContext()
	mem := r.getMemForNewActiveContext()

	repl := repl.NewSessionBuilder(r.logger).WithFS(fs).WithIO(io).WithMEM(mem).Build(r.launchDirectory)

	configuration, err := slpxcfg.LoadFromContent(r.logger, r.launchDirectory, r.setupContent, 10*time.Second, []slpxcfg.Variable{
		{Identifier: "cmd_toggle_editor", Type: object.OBJ_TYPE_STRING, Required: true},
		{Identifier: "cmd_toggle_output", Type: object.OBJ_TYPE_STRING, Required: true},
		{Identifier: "cmd_clear", Type: object.OBJ_TYPE_STRING, Required: true},
		{Identifier: "color_prompt", Type: object.OBJ_TYPE_STRING, Required: true},
		{Identifier: "color_result", Type: object.OBJ_TYPE_STRING, Required: true},
		{Identifier: "color_error", Type: object.OBJ_TYPE_STRING, Required: true},
		{Identifier: "color_help", Type: object.OBJ_TYPE_STRING, Required: true},
		{Identifier: "color_focused_border", Type: object.OBJ_TYPE_STRING, Required: true},
		{Identifier: "color_blurred_border", Type: object.OBJ_TYPE_STRING, Required: true},
		{Identifier: "color_selected_item", Type: object.OBJ_TYPE_STRING, Required: true},
		{Identifier: "color_history_item", Type: object.OBJ_TYPE_STRING, Required: true},
		{Identifier: "color_dirty_prompt", Type: object.OBJ_TYPE_STRING, Required: true},
		{Identifier: "color_secondary_action", Type: object.OBJ_TYPE_STRING, Required: true},
		{Identifier: "command_router", Type: object.OBJ_TYPE_FUNCTION, Required: false},
		{Identifier: "environment_preload", Type: object.OBJ_TYPE_LIST, Required: false},
	}, fs, io)
	if err != nil {
		return nil, err
	}

	var commandRouter object.Function
	if routerObj, ok := configuration["command_router"]; ok {
		commandRouter = routerObj.D.(object.Function)
		r.logger.Info("command_router loaded from configuration")
	} else {
		r.logger.Warn("command_router not found in configuration - custom commands will not be available")
	}

	if preloadObj, ok := configuration["environment_preload"]; ok {
		r.logger.Info("evaluating environment_preload")
		encodedList := preloadObj.Encode()
		_, evalErr := repl.Evaluate(encodedList)
		if evalErr != nil {
			r.logger.Warn("failed to evaluate environment_preload", "error", evalErr)
		} else {
			r.logger.Info("environment_preload evaluated successfully")
		}
	}

	tuiConfig := TuiConfig{
		CmdToggleEditor:      configuration["cmd_toggle_editor"].D.(string),
		CmdToggleOutput:      configuration["cmd_toggle_output"].D.(string),
		CmdClear:             configuration["cmd_clear"].D.(string),
		PromptColor:          configuration["color_prompt"].D.(string),
		ResultColor:          configuration["color_result"].D.(string),
		ErrorColor:           configuration["color_error"].D.(string),
		HelpColor:            configuration["color_help"].D.(string),
		FocusedBorderColor:   configuration["color_focused_border"].D.(string),
		BlurredBorderColor:   configuration["color_blurred_border"].D.(string),
		SelectedItemColor:    configuration["color_selected_item"].D.(string),
		HistoryItemColor:     configuration["color_history_item"].D.(string),
		DirtyPromptColor:     configuration["color_dirty_prompt"].D.(string),
		SecondaryActionColor: configuration["color_secondary_action"].D.(string),
		CommandRouter:        commandRouter,
	}

	restrictedShortcuts := []string{
		"enter",
		"up",
		"down",
		"tab",
		"esc",
		"ctrl+c",
		"ctrl+q",
		"ctrl+d",
	}

	toCheck := []string{
		tuiConfig.CmdToggleEditor,
		tuiConfig.CmdToggleOutput,
		tuiConfig.CmdClear,
	}

	for _, shortcut := range toCheck {
		for _, shortcutRestricted := range restrictedShortcuts {
			if shortcut == shortcutRestricted {
				return nil, fmt.Errorf("shortcut %s is restricted", shortcut)
			}
		}
	}

	ac := activeContext{
		id:          uuid,
		displayName: displayName,
		env:         r.getEvalBuilderForNewActiveContext(uuid),
		fs:          fs,
		io:          io,
		mem:         mem,
		repl:        repl,
		tuiConfig:   tuiConfig,
		onClose: func() error {
			/*
				Note: In the future we may want to have logic to deny close, so
				we made sure to return an error in the fn sig, but not returning one yet
			*/
			r.onActiveContextClose(uuid)
			return nil
		},
	}

	return &ac, nil
}

func (r *runtimeImpl) Stop() error {

	r.acMutex.Lock()
	defer r.acMutex.Unlock()
	for _, ac := range r.activeContexts {
		ac.Close()
	}

	return nil
}
