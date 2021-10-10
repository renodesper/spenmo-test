package logger

// Logger interface
type Logger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
}

// Loggers contains one or more logger
type Loggers struct {
	Ls []Logger
}

const (
	// Debug ...
	Debug = "debug"
	// Info ...
	Info = "info"
	// Warn ...
	Warn = "warn"
	// Error ...
	Error = "error"
)

// Debug implements Logger.Debug
func (ls *Loggers) Debug(args ...interface{}) {
	for _, l := range ls.Ls {
		l.Debug(args...)
	}
}

// Debugf implements Logger.Debugf
func (ls *Loggers) Debugf(template string, args ...interface{}) {
	for _, l := range ls.Ls {
		l.Debugf(template, args...)
	}
}

// Debugw implements Logger.Debugw
func (ls *Loggers) Debugw(msg string, keysAndValues ...interface{}) {
	for _, l := range ls.Ls {
		l.Debugw(msg, keysAndValues...)
	}
}

// Info implements Logger.Info
func (ls *Loggers) Info(args ...interface{}) {
	for _, l := range ls.Ls {
		l.Info(args...)
	}
}

// Infof implements Logger.Infof
func (ls *Loggers) Infof(template string, args ...interface{}) {
	for _, l := range ls.Ls {
		l.Infof(template, args...)
	}
}

// Infow implements Logger.Infow
func (ls *Loggers) Infow(msg string, keysAndValues ...interface{}) {
	for _, l := range ls.Ls {
		l.Infow(msg, keysAndValues...)
	}
}

// Warn implements Logger.Warn
func (ls *Loggers) Warn(args ...interface{}) {
	for _, l := range ls.Ls {
		l.Warn(args...)
	}
}

// Warnf implements Logger.Warnf
func (ls *Loggers) Warnf(template string, args ...interface{}) {
	for _, l := range ls.Ls {
		l.Warnf(template, args...)
	}
}

// Warnw implements Logger.Warnw
func (ls *Loggers) Warnw(msg string, keysAndValues ...interface{}) {
	for _, l := range ls.Ls {
		l.Warnw(msg, keysAndValues...)
	}
}

// Error implements Logger.Error
func (ls *Loggers) Error(args ...interface{}) {
	for _, l := range ls.Ls {
		l.Error(args...)
	}
}

// Errorf implements Logger.Errorf
func (ls *Loggers) Errorf(template string, args ...interface{}) {
	for _, l := range ls.Ls {
		l.Errorf(template, args...)
	}
}

// Errorw implements Logger.Errorw
func (ls *Loggers) Errorw(msg string, keysAndValues ...interface{}) {
	for _, l := range ls.Ls {
		l.Errorw(msg, keysAndValues...)
	}
}

// New loggers
func New(ls ...Logger) *Loggers {
	return &Loggers{ls}
}
