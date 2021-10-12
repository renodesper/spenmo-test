package noop

import (
	"gitlab.com/renodesper/spenmo-test/util/logger"
)

type Noop struct {
}

// CreateLogger creates logger that does nothing
func CreateLogger() logger.Logger {
	return &Noop{}
}

// Debug implements Logger.Debug
func (ls *Noop) Debug(args ...interface{}) {
}

// Debugf implements Logger.Debugf
func (ls *Noop) Debugf(template string, args ...interface{}) {
}

// Debugw implements Logger.Debugw
func (ls *Noop) Debugw(msg string, keysAndValues ...interface{}) {
}

// Info implements Logger.Info
func (ls *Noop) Info(args ...interface{}) {
}

// Infof implements Logger.Infof
func (ls *Noop) Infof(template string, args ...interface{}) {
}

// Infow implements Logger.Infow
func (ls *Noop) Infow(msg string, keysAndValues ...interface{}) {
}

// Warn implements Logger.Warn
func (ls *Noop) Warn(args ...interface{}) {
}

// Warnf implements Logger.Warnf
func (ls *Noop) Warnf(template string, args ...interface{}) {
}

// Warnw implements Logger.Warnw
func (ls *Noop) Warnw(msg string, keysAndValues ...interface{}) {
}

// Error implements Logger.Error
func (ls *Noop) Error(args ...interface{}) {
}

// Errorf implements Logger.Errorf
func (ls *Noop) Errorf(template string, args ...interface{}) {
}

// Errorw implements Logger.Errorw
func (ls *Noop) Errorw(msg string, keysAndValues ...interface{}) {
}
