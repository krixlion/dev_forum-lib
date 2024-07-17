package logging

import (
	"context"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

func init() {
	logger, err := zap.NewProduction(zap.AddCaller(), zap.AddCallerSkip(2))
	if err != nil {
		panic(err)
	}

	otelLogger := otelzap.New(logger)
	otelzap.ReplaceGlobals(otelLogger)
}

type Logger interface {
	Log(ctx context.Context, msg string, keyvals ...interface{})

	// Sync flushes the buffer and writes any pending logs.
	Sync() error
}

func Log(msg string, keyvals ...interface{}) {
	otelzap.S().Infow(msg, keyvals...)
}

// NewLogger returns a new logger safe for concurrent use
// Returns an error on hardware error.
func NewLogger() (Logger, error) {
	logger, err := zap.NewProduction(zap.AddCaller(), zap.AddCallerSkip(2))
	if err != nil {
		return nil, err
	}

	return stdLogger{
		SugaredLogger: otelzap.New(logger).Sugar(),
	}, nil
}

// stdLogger is a wrapper for otelzap.SugaredLogger and implements Logger interface.
type stdLogger struct {
	*otelzap.SugaredLogger
}

func (logger stdLogger) Log(ctx context.Context, msg string, keyvals ...interface{}) {
	logger.InfowContext(ctx, msg, keyvals...)
}
