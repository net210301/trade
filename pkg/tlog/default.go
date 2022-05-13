package tlog


var WatchDog = Config{
	Level: "info",
	AddCaller: true,
	Debug: false,
	CallerSkip: 0,
}.Build()


func Debug(msg string, fields ...Field) {
	WatchDog.desugar.Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	WatchDog.desugar.Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	WatchDog.desugar.Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
	WatchDog.desugar.Error(msg, fields...)
}

func Panic(msg string, fields ...Field) {
	WatchDog.desugar.Panic(msg, fields...)
}

func DPanic(msg string, fields ...Field) {
	WatchDog.desugar.DPanic(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	WatchDog.desugar.Fatal(msg, fields...)
}