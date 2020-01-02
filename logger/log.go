package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"time"
)

var logger *zap.Logger

//var logLevel = zap.NewAtomicLevel()

func init() {
	proEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		TimeKey:     "ts",
		CallerKey:   "caller",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.000000-07:00"))
		},
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	debugEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		TimeKey:     "ts",
		CallerKey:   "caller",
		EncodeLevel: zapcore.CapitalColorLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.000000-07:00"))
		},
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel
	})

	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})

	infoHook := getWriter("./log/demo")
	warnHook := getWriter("./log/demo_error")

	core := zapcore.NewTee(
		zapcore.NewCore(proEncoder, zapcore.AddSync(infoHook), infoLevel),
		zapcore.NewCore(proEncoder, zapcore.AddSync(warnHook), warnLevel),
		zapcore.NewCore(debugEncoder, os.Stdout, zap.DebugLevel),
	)
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

type Level int8

//const (
//	DebugLevel Level = iota - 1
//
//	InfoLevel
//
//	WarnLevel
//
//	ErrorLevel
//
//	DPanicLevel
//
//	PanicLevel
//
//	FatalLevel
//)
//
//func SetLevel(level Level) {
//	logLevel.SetLevel(zapcore.Level(level))
//}

func getWriter(filename string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每1小时(整点)分割一次日志
	hook, err := rotatelogs.New(
		filename+".%Y%m%d.log", // 没有使用go风格反人类的format格式
		//rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	if err != nil {
		panic(err)
	}
	return hook
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}
func DPanic(msg string, fields ...zap.Field) {
	logger.DPanic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
