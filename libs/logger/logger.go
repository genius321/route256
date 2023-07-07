package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// переменные отвечающие за сам логер
var (
	global     *zap.SugaredLogger
	infoLevel  = zap.NewAtomicLevelAt(zap.InfoLevel)
	errorLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
)

// SetLoggerByEnvironment - not thread safe
func SetLoggerByEnvironment(environment string) {
	if environment == "DEVELOPMENT" {
		global = New(infoLevel, os.Stdout)
	} else if environment == "PRODUCTION" {
		// если изменить info на errorLevel, то в логах не будет info
		global = New(errorLevel, os.Stdout,
			// стектрейсы для zap.PanicLevel, для других уровней стектрейсы перебор
			zap.AddStacktrace(zap.NewAtomicLevelAt(zap.PanicLevel)),
		)
	}
}

func New(level zap.AtomicLevel, sink io.Writer, opts ...zap.Option) *zap.SugaredLogger {
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:  "ts",
			LevelKey: "level",
			// NameKey:        "logger",
			// CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.TimeEncoderOfLayout(time.DateTime),
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
		zapcore.AddSync(sink),
		level,
	)

	return zap.New(core, opts...).Sugar()
}

func Info(args ...interface{}) {
	global.Info(args...)
}

func Infof(template string, args ...interface{}) {
	global.Infof(template, args...)
}

func Infoln(args ...interface{}) {
	global.Infoln(args...)
}

func Error(args ...interface{}) {
	global.Error(args...)
}

func Errorf(ctx context.Context, method, template string, args ...interface{}) {
	withTraceID(ctx).Desugar().
		With(zap.String("method", method)).Sugar().Errorf(template, args...)
}

func Fatal(args ...interface{}) {
	global.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	global.Fatalf(template, args...)
}

func Panic(args ...interface{}) {
	global.Panic()
}

func Panicf(template string, args ...interface{}) {
	global.Panicf(template, args...)
}

func withTraceID(ctx context.Context) *zap.SugaredLogger {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return global
	}

	if sc, ok := span.Context().(jaeger.SpanContext); ok {
		return global.Desugar().With(
			zap.Stringer("trace_id", sc.TraceID()),
			zap.Stringer("span_id", sc.SpanID()),
		).Sugar()
	}

	return global
}
