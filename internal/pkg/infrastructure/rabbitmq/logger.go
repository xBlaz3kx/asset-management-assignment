package rabbitmq

import (
	"github.com/xBlaz3kx/DevX/observability"
)

type logger struct {
	obs observability.Observability
}

func newLogger(obs observability.Observability) *logger {
	return &logger{}
}

func (l *logger) Fatalf(s string, i ...interface{}) {
	l.obs.Log().Fatal(s)
}

func (l *logger) Errorf(s string, i ...interface{}) {
	l.obs.Log().Error(s)
}

func (l *logger) Warnf(s string, i ...interface{}) {
	l.obs.Log().Warn(s)
}

func (l *logger) Infof(s string, i ...interface{}) {
	l.obs.Log().Info(s)
}

func (l *logger) Debugf(s string, i ...interface{}) {
	l.obs.Log().Debug(s)
}