package repl

import (
	"log/slog"

	"github.com/bosley/slpx/pkg/env"
	"github.com/bosley/slpx/pkg/object"
)

type replFunctions struct {
	logger  *slog.Logger
	session *Session
}

func newReplFunctions(logger *slog.Logger, session *Session) env.FunctionGroup {
	return &replFunctions{
		logger:  logger,
		session: session,
	}
}

func (r *replFunctions) Name() string {
	return "repl"
}

func (r *replFunctions) Functions() map[object.Identifier]env.EnvFunction {
	/*
		repl/quit
	*/
	return map[object.Identifier]env.EnvFunction{}
}
