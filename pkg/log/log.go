package log

var DefaultLog Logger

// Level type
type Level uint8

// Log Level
const (
	ErrorLevel Level = iota
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
	Disabled
)

type Logger interface {
	Trace(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
}

func Trace(format string, args ...interface{}) {
	DefaultLog.Trace(format, args...)
}

func Debug(format string, args ...interface{}) {
	DefaultLog.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	DefaultLog.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	DefaultLog.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	DefaultLog.Error(format, args...)
}

func init() {
	DefaultLog = NewSimpleLogger(DebugLevel, true)
	// DefaultLog = NewZeroLogger(DebugLevel)
}
