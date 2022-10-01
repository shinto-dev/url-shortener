package logging

import (
	"context"
	"log"
	"time"
	"url-shortener/foundation/observation/trace"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logFieldsKeyType string

const logFieldsKey = logFieldsKeyType("key")

var zapLogger *zap.Logger

func init() {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	logger, err := loggerConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatal(err)
	}

	zapLogger = logger.Sugar().Desugar()
}

func WithLogger(ctx context.Context) context.Context {
	var fields []zap.Field
	return context.WithValue(ctx, logFieldsKey, &fields)
}

func LField(name string, value interface{}) func(ctx context.Context) {
	return func(ctx context.Context) {
		logFields := getLogFields(ctx)
		*logFields = append(*logFields, zap.Any(name, value))
	}
}

func getLogFields(ctx context.Context) *[]zap.Field {
	logFields := ctx.Value(logFieldsKey)
	if logFields == nil {
		return &[]zap.Field{}
	}

	return logFields.(*[]zap.Field)
}

func getLogFieldsWithTraceID(ctx context.Context) []zap.Field {
	logFields := *getLogFields(ctx)
	traceID := trace.GetTraceID(ctx)
	if traceID != "" {
		logFields = append(logFields, zap.String("trace_id", traceID))
	}
	return logFields
}

type LogCtx struct {
	ctx context.Context
}

func FromContext(ctx context.Context) *LogCtx {
	return &LogCtx{ctx: ctx}
}

func (l *LogCtx) Info(msg string) {
	zapLogger.Info(msg, getLogFieldsWithTraceID(l.ctx)...)
}

func (l *LogCtx) Error(msg string, err error) {
	logFields := append(getLogFieldsWithTraceID(l.ctx), zap.Error(err))
	zapLogger.Error(msg, logFields...)
}

func (l *LogCtx) Warn(msg string) {
	zapLogger.Warn(msg, getLogFieldsWithTraceID(l.ctx)...)
}

func (l *LogCtx) Fatal(msg string, err error) {
	logFields := append(getLogFieldsWithTraceID(l.ctx), zap.Error(err))
	zapLogger.Fatal(msg, logFields...)
}

func (l *LogCtx) Debug(msg string) {
	zapLogger.Debug(msg, getLogFieldsWithTraceID(l.ctx)...)
}

func (l *LogCtx) WithFields(fields ...zap.Field) *LogWithFields {
	return &LogWithFields{
		fields: append(fields, getLogFieldsWithTraceID(l.ctx)...),
	}
}

type LogWithFields struct {
	fields []zap.Field
}

func WithFields(fields ...zap.Field) *LogWithFields {
	return &LogWithFields{
		fields: fields,
	}
}

func (l *LogWithFields) Info(msg string) {
	zapLogger.Info(msg, l.fields...)
}

func (l *LogWithFields) Error(msg string) {
	zapLogger.Error(msg, l.fields...)
}

func (l *LogWithFields) Warn(msg string) {
	zapLogger.Warn(msg, l.fields...)
}

func (l *LogWithFields) Fatal(msg string) {
	zapLogger.Fatal(msg, l.fields...)
}

func (l *LogWithFields) Debug(msg string) {
	zapLogger.Debug(msg, l.fields...)
}
