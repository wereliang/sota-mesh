package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

type zeroLogger struct {
}

func (l *zeroLogger) Trace(format string, args ...interface{}) {
	zlog.Trace().Msgf(format, args...)
}

func (l *zeroLogger) Debug(format string, args ...interface{}) {
	zlog.Debug().Msgf(format, args...)
}

func (l *zeroLogger) Info(format string, args ...interface{}) {
	zlog.Info().Msgf(format, args...)
}

func (l *zeroLogger) Warn(format string, args ...interface{}) {
	zlog.Warn().Msgf(format, args...)
}

func (l *zeroLogger) Error(format string, args ...interface{}) {
	zlog.Error().Msgf(format, args...)
}

var zlogLevelMapping = map[Level]zerolog.Level{
	ErrorLevel: zerolog.ErrorLevel,
	WarnLevel:  zerolog.WarnLevel,
	InfoLevel:  zerolog.InfoLevel,
	DebugLevel: zerolog.DebugLevel,
	TraceLevel: zerolog.TraceLevel,
	Disabled:   zerolog.Disabled,
}

func NewZeroLogger(l Level) Logger {
	zl := zlogLevelMapping[l]
	zlog.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}).
		Level(zl).
		With().
		Timestamp().
		Caller().
		Logger()
	return &zeroLogger{}
}
