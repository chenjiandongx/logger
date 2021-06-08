package logger

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Level int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota - 1

	// InfoLevel is the default logging priority.
	InfoLevel

	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel

	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel

	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel

	// PanicLevel logs a message, then panics.
	PanicLevel

	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel
)

// Options is the option set for Logger.
type Options struct {
	// Stdout sets the writer as stdout if it is true.
	Stdout bool

	// ConsoleMode sets logger to be the console mode which claims the logger encoder type as console.
	ConsoleMode bool

	// Filename is the file to write logs to.  Backup log files will be retained
	// in the same directory.
	Filename string

	// MaxSize is the maximum size in megabytes of the log file before it gets rotated.
	MaxSize int

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int

	// MaxBackups is the maximum number of old log files to retain. The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int

	// Level is a logging priority. Higher levels are more important.
	Level Level
}

type Logger struct {
	sugared *zap.SugaredLogger
}

// With adds a variadic number of fields to the logging context. It accepts a
// mix of strongly-typed Field objects and loosely-typed key-value pairs. When
// processing pairs, the first element of the pair is used as the field key
// and the second as the field value.
func (l Logger) With(args ...interface{}) Logger {
	return Logger{sugared: l.sugared.With(args...)}
}

// Println is the alias for Info
func (l Logger) Println(args ...interface{}) {
	l.sugared.Info(args...)
}

// Printf is the alias for Infof
func (l Logger) Printf(template string, args ...interface{}) {
	l.sugared.Infof(template, args...)
}

// Debug uses fmt.Sprint to construct and log a message.
func (l Logger) Debug(args ...interface{}) {
	l.sugared.Debug(args...)
}

// Info uses fmt.Sprint to construct and log a message.
func (l Logger) Info(args ...interface{}) {
	l.sugared.Info(args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func (l Logger) Warn(args ...interface{}) {
	l.sugared.Warn(args...)
}

// Error uses fmt.Sprint to construct and log a message.
func (l Logger) Error(args ...interface{}) {
	l.sugared.Error(args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func (l Logger) Panic(args ...interface{}) {
	l.sugared.Panic(args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func (l Logger) Fatal(args ...interface{}) {
	l.sugared.Fatal(args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func (l Logger) Debugf(template string, args ...interface{}) {
	l.sugared.Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func (l Logger) Infof(template string, args ...interface{}) {
	l.sugared.Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (l Logger) Warnf(template string, args ...interface{}) {
	l.sugared.Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (l Logger) Errorf(template string, args ...interface{}) {
	l.sugared.Errorf(template, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func (l Logger) Panicf(template string, args ...interface{}) {
	l.sugared.Panicf(template, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (l Logger) Fatalf(template string, args ...interface{}) {
	l.sugared.Fatalf(template, args...)
}

// New returns the logger instance with Production Config by default.
func New(opt Options) Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Local().Format("2006-01-02 15:04:05.000"))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	encoder := zapcore.NewJSONEncoder(encoderConfig)
	if opt.ConsoleMode {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	if err := os.MkdirAll(filepath.Dir(opt.Filename), os.ModePerm); err != nil {
		panic(err)
	}

	var w zapcore.WriteSyncer
	if opt.Stdout {
		w = zapcore.AddSync(os.Stdout)
	} else {
		w = zapcore.AddSync(&lumberjack.Logger{
			Filename:   opt.Filename,
			MaxSize:    opt.MaxSize,
			MaxBackups: opt.MaxBackups,
			MaxAge:     opt.MaxAge,
			LocalTime:  true,
		})
	}

	core := zapcore.NewCore(encoder, w, zapcore.Level(opt.Level))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return Logger{sugared: logger.Sugar()}
}

var std = New(Options{Stdout: true})

// StandardLogger returns the standard logger with stdout output.
func StandardLogger() Logger {
	return std
}

// SetOptions sets the options for the standard logger.
func SetOptions(opt Options) {
	std = New(opt)
}

// With adds a variadic number of fields to the logging context. It accepts a
// mix of strongly-typed Field objects and loosely-typed key-value pairs. When
// processing pairs, the first element of the pair is used as the field key
// and the second as the field value.
func With(args ...interface{}) Logger {
	s := std
	s.sugared = std.sugared.With(args...)
	return s
}

// Println is the alias for Info
func Println(args ...interface{}) {
	std.sugared.Info(args...)
}

// Printf is the alias for Infof
func Printf(template string, args ...interface{}) {
	std.sugared.Infof(template, args...)
}

// Debug uses fmt.Sprint to construct and log a message.
func Debug(args ...interface{}) {
	std.sugared.Debug(args...)
}

// Info uses fmt.Sprint to construct and log a message.
func Info(args ...interface{}) {
	std.sugared.Info(args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(args ...interface{}) {
	std.sugared.Warn(args...)
}

// Error uses fmt.Sprint to construct and log a message.
func Error(args ...interface{}) {
	std.sugared.Error(args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func Panic(args ...interface{}) {
	std.sugared.Panic(args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	std.sugared.Fatal(args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	std.sugared.Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
	std.sugared.Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	std.sugared.Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	std.sugared.Errorf(template, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(template string, args ...interface{}) {
	std.sugared.Panicf(template, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(template string, args ...interface{}) {
	std.sugared.Fatalf(template, args...)
}
