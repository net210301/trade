package tlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// Config 設定 Logger
type Config struct{
	Level		string
	Fields 		[]zap.Field
	AddCaller 	bool
	Core 		zapcore.Core
	CallerSkip	int
	Debug 		bool
	EncoderConfig *zapcore.EncoderConfig
}

func (config Config) Build() *Logger {
	// 如果要做一些前置的調整在這裡做
	if config.EncoderConfig == nil{
		config.EncoderConfig = DefaultZapConfig()
	}
	logger := newLogger(&config)
	return logger
}

func newLogger(config *Config) *Logger{
	zapOptions := make([]zap.Option, 0)
	zapOptions = append(zapOptions, zap.AddStacktrace(zap.DPanicLevel))
	if config.AddCaller {
		zapOptions = append(zapOptions, zap.AddCaller(), zap.AddCallerSkip(config.CallerSkip))
	}
	if len(config.Fields) > 0 {
		zapOptions = append(zapOptions, zap.Fields(config.Fields...))
	}
	var ws zapcore.WriteSyncer
	ws = os.Stdout

	lv := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if err := lv.UnmarshalText([]byte(config.Level)); err != nil {
		panic(err)
	}

	encoderConfig := *config.EncoderConfig
	core := config.Core
	if core == nil {
		core = zapcore.NewCore(
			func() zapcore.Encoder {
				if config.Debug {
					return zapcore.NewConsoleEncoder(encoderConfig)
				}
				return zapcore.NewJSONEncoder(encoderConfig)
			}(),
			ws,
			lv,
		)
	}

	zapLogger := zap.New(
		core,
		zapOptions...,
	)
	return &Logger{
		desugar: zapLogger,
		lv:      &lv,
		config:  *config,
		sugar:   zapLogger.Sugar(),
	}
}


// DefaultZapConfig 預設的 Zap 格式，及解碼器
func DefaultZapConfig()*zapcore.EncoderConfig{
	return &zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "lv",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}