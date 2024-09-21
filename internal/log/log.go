package log

import (
	"context"
	"fmt"
	"os"

	"log/slog"
)

type handler struct {
	slog.Handler
}

const (
	LevelDebug   = slog.LevelDebug
	LevelInfo    = slog.LevelInfo
	LevelWarning = slog.LevelWarn
	LevelError   = slog.LevelError
	LevelFatal   = slog.Level(12)
)

func replaceAttr(_ []string, a slog.Attr) slog.Attr {
	if a.Key != slog.LevelKey {
		return a
	}

	level := a.Value.Any().(slog.Level)
	switch {
	case level < LevelDebug:
		a.Value = slog.StringValue("trace")
	case level < LevelInfo:
		a.Value = slog.StringValue("debug")
	case level < LevelWarning:
		a.Value = slog.StringValue("info")
	case level < LevelError:
		a.Value = slog.StringValue("warning")
	case level < LevelFatal:
		a.Value = slog.StringValue("error")
	default:
		a.Value = slog.StringValue("fatal")
	}
	return a
}

func wrapErr(msg string, err error) error {
	return fmt.Errorf("%s: '%w'", msg, err)
}

func fmtMsgErr(msg string, err error) string {
	return wrapErr(msg, err).Error()
}

var Logger = newLogger(LevelDebug)

func logCtx(ctx context.Context, level slog.Level, msg string, args ...interface{}) {
	Logger.Log(ctx, level, fmt.Sprintf(msg, args...))
}

func logErrCtx(ctx context.Context, level slog.Level, msg string, err error) {
	Logger.Log(ctx, level, fmtMsgErr(msg, err))
}

func Debug(msg string, args ...interface{}) { logCtx(context.Background(), LevelDebug, msg, args...) }
func Info(msg string, args ...interface{})  { logCtx(context.Background(), LevelInfo, msg, args...) }

func Warning(msg string, err error) {
	logErrCtx(context.Background(), LevelWarning, msg, err)
}
func Error(msg string, err error) {
	logErrCtx(context.Background(), LevelError, msg, err)
}
func Fatal(msg string, err error) {
	logErrCtx(context.Background(), LevelFatal, msg, err)
	panic(err)
}

func With(ctx context.Context) struct {
	Debug   func(msg string, args ...interface{})
	Info    func(msg string, args ...interface{})
	Warning func(msg string, err error)
	Error   func(msg string, err error)
	Fatal   func(msg string, err error)
} {
	return struct {
		Debug   func(msg string, args ...interface{})
		Info    func(msg string, args ...interface{})
		Warning func(msg string, err error)
		Error   func(msg string, err error)
		Fatal   func(msg string, err error)
	}{
		Debug:   func(msg string, args ...interface{}) { logCtx(ctx, LevelDebug, msg, args...) },
		Info:    func(msg string, args ...interface{}) { logCtx(ctx, LevelInfo, msg, args...) },
		Warning: func(msg string, err error) { logErrCtx(ctx, LevelWarning, msg, err) },
		Error: func(msg string, err error) {
			logErrCtx(ctx, LevelError, msg, err)
		},
		Fatal: func(msg string, err error) {
			logErrCtx(ctx, LevelFatal, msg, err)
			panic(err)
		},
	}
}

func newLogger(level slog.Level) *slog.Logger {
	return slog.New(
		handler{
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level:       level,
					ReplaceAttr: replaceAttr,
				},
			),
		},
	)
}

func Setup() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		switch logLevel {
		case "debug":
			Logger = newLogger(LevelDebug)
		case "info":
			Logger = newLogger(LevelInfo)
		case "warning":
			Logger = newLogger(LevelWarning)
		case "error":
			Logger = newLogger(LevelError)
		case "fatal":
			Logger = newLogger(LevelFatal)
		default:
			Warning("Skipping level override", fmt.Errorf("unknow level %s", logLevel))
		}
	}
}
