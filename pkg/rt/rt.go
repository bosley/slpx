package rt

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"sync"
	"time"

	"github.com/bosley/slpx/pkg/planar"
	"github.com/bosley/slpx/pkg/planar/goba"
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
	ForegroundDefaultColor string
	BackgroundDefaultColor string
	CmdToggleEditor        string
	CmdToggleOutput        string
	CmdClear               string
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

	kvProvider planar.KVProvider
	rootKV     planar.KV
}

type Config struct {
	Logger          *slog.Logger
	SLPXHome        string
	LaunchDirectory string
	SetupContent    string
}

func New(config Config) (Runtime, error) {

	fmt.Println("Opening badger backend", filepath.Join(config.SLPXHome, "root.db"))
	kvp := goba.OpenBadgerBackend(filepath.Join(config.SLPXHome, "root.db"))

	rootKV, err := kvp.LoadOrCreate(config.Logger, "slpx-root")
	if err != nil {
		return nil, err
	}

	exists, err := rootKV.Exists([]byte("slpx-root-init-complete"))
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := firstTimeInit(config.Logger, rootKV); err != nil {
			return nil, err
		}
	}

	return &runtimeImpl{
		logger:          config.Logger,
		slpxHome:        config.SLPXHome,
		launchDirectory: config.LaunchDirectory,
		kvProvider:      kvp,
		activeContexts:  make(map[string]activeContext),
		acMutex:         sync.Mutex{},
		rootKV:          rootKV,
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

	repl := repl.NewSessionBuilder(r.logger).Build(r.launchDirectory)

	fs := r.getFsForNewActiveContext()
	io := r.getIoForNewActiveContext()

	configuration, err := slpxcfg.LoadFromContent(r.logger, r.launchDirectory, r.setupContent, 10*time.Second, []slpxcfg.Variable{
		{Identifier: "text_foreground", Type: object.OBJ_TYPE_STRING, Required: true},
		{Identifier: "text_background", Type: object.OBJ_TYPE_STRING, Required: true},
		{Identifier: "cmd_toggle_editor", Type: object.OBJ_TYPE_STRING, Required: true},
		{Identifier: "cmd_toggle_output", Type: object.OBJ_TYPE_STRING, Required: true},
		{Identifier: "cmd_clear", Type: object.OBJ_TYPE_STRING, Required: true},
	}, fs, io)
	if err != nil {
		return nil, err
	}

	tuiConfig := TuiConfig{
		ForegroundDefaultColor: configuration["text_foreground"].D.(string),
		BackgroundDefaultColor: configuration["text_background"].D.(string),
		CmdToggleEditor:        configuration["cmd_toggle_editor"].D.(string),
		CmdToggleOutput:        configuration["cmd_toggle_output"].D.(string),
		CmdClear:               configuration["cmd_clear"].D.(string),
	}

	ac := activeContext{
		id:          uuid,
		displayName: displayName,
		env:         r.getEvalBuilderForNewActiveContext(uuid),
		fs:          fs,
		io:          io,
		mem:         r.getMemForNewActiveContext(),
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

func firstTimeInit(logger *slog.Logger, rootKV planar.KV) error {
	logger.Info("performing first time initialization")

	nowStr := time.Now().Format(time.RFC3339)

	res, err := rootKV.Set([]byte("slpx-root-init-complete"), []byte(nowStr))
	if err != nil {
		return err
	}
	if res.UniqueKeyBytesDelta == 0 {
		return errors.New("failed to set slpx-root-init-complete")
	}

	return nil
}

func (r *runtimeImpl) Stop() error {

	r.acMutex.Lock()
	defer r.acMutex.Unlock()
	for _, ac := range r.activeContexts {
		ac.Close()
	}

	return r.kvProvider.Close()
}
