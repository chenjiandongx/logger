# logger

A lightweight, auto-rolling logger for Go.

### Options

```golang
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

	// Skip is the number of callers skipped by caller annotation
	Skip int
}
```

### Example

```golang
package main

import "github.com/chenjiandongx/logger"

// InitLogger initializes the logger
func InitLogger() {
	// feel free to config the options.
	logger.SetOptions(logger.Options{
		Filename:   "/data/log/awesome-project/applog",
		MaxSize:    1000, // 1GB
		MaxAge:     3,    // 3 days
		MaxBackups: 3,    // 3 backups
	})
}

func main() {
	// Note: init logger when you want to run your program in the production env.
	// for example:
	InitLogger()

	// use logger Method anywhere you want directly, such as Info/Warn/Error/...
	// logs will be displayed on the stdout stream by default.
	logger.Info("This is the info level message.")
	logger.Warnf("This is the warn level message. %s", "oop!")
	logger.Error("Something error here.")
}

```

### LICENSE

MIT [©chenjiandongx](https://github.com/chenjiandongx)
