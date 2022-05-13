package tlog

import (
	"go.uber.org/zap"
)

const (
	// DebugLevel 在 production 環境中禁用，因為數量非常龐大
	DebugLevel = zap.DebugLevel
	// InfoLevel 預設的等級
	InfoLevel = zap.InfoLevel
	// WarnLevel 比 Info 更重要，但不需要團隊開會討論才能設定這個 level，下面等級請通過討論再定義
	WarnLevel = zap.WarnLevel
	// ErrorLevel 較高的優先處理順序，如果應用程式沒啥大錯，不應該跑這個等級
	ErrorLevel = zap.ErrorLevel
	// PanicLevel 記錄一個訊息後 panic
	PanicLevel = zap.PanicLevel
	// FatalLevel 記錄一條訊息，然後應用程式呼叫 os.Exit(1).
	FatalLevel = zap.FatalLevel
)

type Field  = zap.Field

type Logger struct {
	config 	Config
	sugar 	*zap.SugaredLogger
	lv 		*zap.AtomicLevel
	desugar *zap.Logger
}

func (logger *Logger) Debug(msg string, fields ...Field) {
	logger.desugar.Debug(msg, fields...)
}

func (logger *Logger) Info(msg string, fields ...Field) {
	logger.desugar.Info(msg, fields...)
}

func (logger *Logger) Warn(msg string, fields ...Field) {
	logger.desugar.Warn(msg, fields...)
}

func (logger *Logger) Error(msg string, fields ...Field) {
	logger.desugar.Error(msg, fields...)
}

func (logger *Logger) Panic(msg string, fields ...Field) {
	logger.desugar.Panic(msg, fields...)
}

func (logger *Logger) DPanic(msg string, fields ...Field) {
	logger.desugar.DPanic(msg, fields...)
}

func (logger *Logger) Fatal(msg string, fields ...Field) {
	logger.desugar.Fatal(msg, fields...)
}

// With ...
func (logger *Logger) With(fields ...Field) *Logger {
	desugarLogger := logger.desugar.With(fields...)
	return &Logger{
		desugar: desugarLogger,
		lv:      logger.lv,
		sugar:   desugarLogger.Sugar(),
		config:  logger.config,
	}
}