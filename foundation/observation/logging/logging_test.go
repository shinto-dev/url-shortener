package logging

import (
	"context"
	"errors"
	"testing"

	"github.com/shinto-dev/url-shortener/foundation/observation"
	"github.com/shinto-dev/url-shortener/foundation/observation/trace"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLogCtx(t *testing.T) {
	t.Run("debug log should include log fields", func(t *testing.T) {
		observedLogs := initLogger()

		ctx := WithLogger(context.Background())
		observation.Add(ctx, LField("name", "sample test"))
		FromContext(ctx).Debug("some debug log message")

		firstLog := observedLogs.All()[0]
		assertLog(
			t, firstLog, zapcore.DebugLevel, "some debug log message",
			zap.String("name", "sample test"),
		)
	})

	t.Run("info log should include log fields", func(t *testing.T) {
		observedLogs := initLogger()

		ctx := WithLogger(context.Background())
		observation.Add(ctx, LField("name", "sample test"))
		FromContext(ctx).Info("some info log message")

		firstLog := observedLogs.All()[0]
		assertLog(
			t, firstLog, zapcore.InfoLevel, "some info log message",
			zap.String("name", "sample test"),
		)
	})

	t.Run("warn log should include log fields", func(t *testing.T) {
		observedLogs := initLogger()

		ctx := WithLogger(context.Background())
		observation.Add(ctx, LField("name", "sample test"))
		FromContext(ctx).Warn("some warn log message")

		firstLog := observedLogs.All()[0]
		assertLog(
			t, firstLog, zapcore.WarnLevel, "some warn log message",
			zap.String("name", "sample test"),
		)
	})

	t.Run("error log should include log fields", func(t *testing.T) {
		observedLogs := initLogger()

		ctx := WithLogger(context.Background())
		observation.Add(ctx, LField("name", "sample test"))

		FromContext(ctx).Error("some error log message", errors.New("some error"))

		firstLog := observedLogs.All()[0]
		assertLog(
			t, firstLog, zapcore.ErrorLevel, "some error log message",
			zap.String("name", "sample test"),
			zap.Error(errors.New("some error")),
		)
	})

	t.Run("Fatal log should include log fields", func(t *testing.T) {
		t.Skip()
		observedLogs := initLogger()

		ctx := WithLogger(context.Background())
		observation.Add(ctx, LField("name", "sample test"))
		FromContext(ctx).Fatal("some fatal log message", errors.New("some error"))

		firstLog := observedLogs.All()[0]
		assert.Equal(t, "some fatal log message", firstLog.Message)
		assert.Equal(t, zapcore.Level(2), firstLog.Level)
		assert.Equal(t, []zap.Field{zap.String("name", "sample test")}, firstLog.Context)
	})

	t.Run("with fields should include log fields and temporary log fields", func(t *testing.T) {
		observedLogs := initLogger()

		ctx := WithLogger(context.Background())
		observation.Add(ctx, LField("name", "sample test"))

		FromContext(ctx).WithFields(zap.String("name2", "test2")).
			Info("some error log message")

		firstLog := observedLogs.All()[0]
		assertLog(
			t, firstLog, zapcore.InfoLevel, "some error log message",
			zap.String("name2", "test2"),
			zap.String("name", "sample test"),
		)
	})

	t.Run("should include multiple fields", func(t *testing.T) {
		observedLogs := initLogger()

		ctx := WithLogger(context.Background())
		observation.Add(ctx,
			LField("name", "sample test"),
			LField("context", "test"),
			LField("x", 2),
		)
		observation.Add(ctx,
			LField("y", 4),
		)
		FromContext(ctx).Info("some log message")

		firstLog := observedLogs.All()[0]
		assertLog(t, firstLog, zapcore.InfoLevel, "some log message",
			zap.String("name", "sample test"),
			zap.String("context", "test"),
			zap.Int("x", 2),
			zap.Int("y", 4),
		)
	})

	t.Run("should include trace id", func(t *testing.T) {
		observedLogs := initLogger()

		ctx := WithLogger(context.Background())
		ctx = trace.WithTraceID(ctx, "test-trace-id")
		observation.Add(ctx, LField("name", "sample test"))
		FromContext(ctx).Info("some log message")

		firstLog := observedLogs.All()[0]
		assertLog(t, firstLog, zapcore.InfoLevel, "some log message",
			zap.String("name", "sample test"),
			zap.String("trace_id", "test-trace-id"),
		)
	})
}

func initLogger() *observer.ObservedLogs {
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)

	zapLogger = observedLogger
	return observedLogs
}

func TestLogWithFields(t *testing.T) {
	t.Run("debug log should include log fields", func(t *testing.T) {
		observedLogs := initLogger()

		WithFields(zap.String("name", "sample test")).
			Debug("some debug log message")

		firstLog := observedLogs.All()[0]
		assertLog(
			t, firstLog, zapcore.DebugLevel, "some debug log message",
			zap.String("name", "sample test"),
		)
	})

	t.Run("info log should include log fields", func(t *testing.T) {
		observedLogs := initLogger()

		WithFields(zap.String("name", "sample test")).
			Info("some info log message")

		firstLog := observedLogs.All()[0]
		assertLog(
			t, firstLog, zapcore.InfoLevel, "some info log message",
			zap.String("name", "sample test"),
		)
	})

	t.Run("warn log should include log fields", func(t *testing.T) {
		observedLogs := initLogger()

		WithFields(zap.String("name", "sample test")).
			Warn("some warn log message")

		firstLog := observedLogs.All()[0]
		assertLog(
			t, firstLog, zapcore.WarnLevel, "some warn log message",
			zap.String("name", "sample test"),
		)
	})

	t.Run("error log should include log fields", func(t *testing.T) {
		observedLogs := initLogger()

		WithFields(
			zap.String("name", "sample test"),
			zap.Error(errors.New("some error")),
		).Error("some error log message")

		firstLog := observedLogs.All()[0]
		assertLog(
			t, firstLog, zapcore.ErrorLevel, "some error log message",
			zap.String("name", "sample test"),
			zap.Error(errors.New("some error")),
		)
	})

	t.Run("Fatal log should include log fields", func(t *testing.T) {
		t.Skip()
		observedLogs := initLogger()

		WithFields(
			zap.String("name", "sample test"),
			zap.Error(errors.New("some error")),
		).Error("some fatal log message")

		firstLog := observedLogs.All()[0]
		assert.Equal(t, "some fatal log message", firstLog.Message)
		assert.Equal(t, zapcore.Level(2), firstLog.Level)
		assert.Equal(t, []zap.Field{zap.String("name", "sample test")}, firstLog.Context)
	})
}

func assertLog(t *testing.T, log observer.LoggedEntry, level zapcore.Level, msg string, fields ...zap.Field) {
	assert.Equal(t, msg, log.Message)
	assert.Equal(t, level, log.Level)
	assert.Equal(t, fields, log.Context)
}
