package logging

import (
	"context"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

var global stdLogger

func init() {
	l, err := NewLogger()
	if err != nil {
		panic(err)
	}

	global = l.(stdLogger)
}

type Logger interface {
	Log(ctx context.Context, msg string, keyvals ...interface{})
}

// Log implements Logger
type stdLogger struct {
	*otelzap.SugaredLogger
}

func Log(msg string, keyvals ...interface{}) {
	global.Infow(msg, keyvals...)
}

// NewLogger returns an error on hardware error.
func NewLogger() (Logger, error) {
	logger, err := zap.NewProduction(zap.AddCaller(), zap.AddCallerSkip(2))
	otelLogger := otelzap.New(logger)
	sugar := otelLogger.Sugar()
	defer sugar.Sync()

	otelzap.ReplaceGlobals(otelLogger)

	return stdLogger{
		SugaredLogger: sugar,
	}, err
}

func (log stdLogger) Log(ctx context.Context, msg string, keyvals ...interface{}) {
	log.InfowContext(ctx, msg, keyvals...)
}
