package logging

import (
	"context"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/grpclog"
)

var _ grpclog.LoggerV2 = (*SugaredZapLogger)(nil)

func init() {
	logger, err := NewLogger()
	if err != nil {
		panic(err)
	}

	otelzap.ReplaceGlobals(logger.sugared.Desugar())
}

type Logger interface {
	Log(ctx context.Context, msg string, keyvals ...interface{})

	// Sync flushes the buffer and writes any pending logs.
	Sync() error
}

func Log(msg string, keyvals ...interface{}) {
	otelzap.S().Infow(msg, keyvals...)
}

// NewLogger returns a new logger safe for concurrent use.
func NewLogger() (SugaredZapLogger, error) {
	logger, err := zap.NewProduction(zap.AddCaller(), zap.AddCallerSkip(2))
	if err != nil {
		return SugaredZapLogger{}, err
	}

	return SugaredZapLogger{
		sugared: otelzap.New(logger).Sugar(),
	}, nil
}

// SugaredZapLogger is a wrapper for otelzap.SugaredLogger and implements Logger and grpclog.LoggerV2 interfaces.
type SugaredZapLogger struct {
	sugared *otelzap.SugaredLogger
}

func (logger SugaredZapLogger) Log(ctx context.Context, msg string, keyvals ...interface{}) {
	logger.sugared.InfowContext(ctx, msg, keyvals...)
}

func (logger SugaredZapLogger) Sync() error {
	return logger.sugared.Sync()
}

// Info logs to INFO log. Arguments are handled in the manner of fmt.Print.
func (logger SugaredZapLogger) Info(args ...any) {
	logger.sugared.Info(args...)
}

// Infoln logs to INFO log. Arguments are handled in the manner of fmt.Println.
func (logger SugaredZapLogger) Infoln(args ...any) {
	logger.sugared.Infoln(args...)
}

// Infof logs to INFO log. Arguments are handled in the manner of fmt.Printf.
func (logger SugaredZapLogger) Infof(format string, args ...any) {
	logger.sugared.Infof(format, args...)
}

// Warning logs to WARNING log. Arguments are handled in the manner of fmt.Print.
func (logger SugaredZapLogger) Warning(args ...any) {
	logger.sugared.Warn(args...)
}

// Warningln logs to WARNING log. Arguments are handled in the manner of fmt.Println.
func (logger SugaredZapLogger) Warningln(args ...any) {
	logger.sugared.Warnln(args...)
}

// Warningf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
func (logger SugaredZapLogger) Warningf(format string, args ...any) {
	logger.sugared.Warnf(format, args...)
}

// Error logs to ERROR log. Arguments are handled in the manner of fmt.Print.
func (logger SugaredZapLogger) Error(args ...any) {
	logger.sugared.Error(args...)
}

// Errorln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
func (logger SugaredZapLogger) Errorln(args ...any) {
	logger.sugared.Errorln(args...)
}

// Errorf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
func (logger SugaredZapLogger) Errorf(format string, args ...any) {
	logger.sugared.Errorf(format, args...)
}

// Fatal logs to ERROR log. Arguments are handled in the manner of fmt.Print.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (logger SugaredZapLogger) Fatal(args ...any) {
	logger.sugared.Fatal(args...)
}

// Fatalln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (logger SugaredZapLogger) Fatalln(args ...any) {
	logger.sugared.Fatalln(args...)
}

// Fatalf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (logger SugaredZapLogger) Fatalf(format string, args ...any) {
	logger.sugared.Fatalf(format, args...)
}

// V reports whether verbosity level l is at least the requested verbose level.
func (logger SugaredZapLogger) V(l int) bool {
	return logger.sugared.Level().Enabled(zapcore.Level(l))
}
