package util

import (
	"fmt"
	"os"
	"time"
)

const (
	// LevelError logs errors.
	LevelError = iota
	// LevelWarning logs warnings.
	LevelWarning
	// LevelInformational logs informational messages.
	LevelInformational
	// LevelDebug logs debug messages.
	LevelDebug
)

var logger *Logger

// Logger is the logger.
type Logger struct {
	level int
}

// Println prints a message.
func (ll *Logger) Println(msg string) {
	fmt.Printf("%s %s", time.Now().Format("2006-01-02 15:04:05 -0700"), msg)
}

// Panic logs a fatal error and exits.
func (ll *Logger) Panic(format string, v ...interface{}) {
	if LevelError > ll.level {
		return
	}
	msg := fmt.Sprintf("[Panic] "+format, v...)
	ll.Println(msg)
	os.Exit(1)
}

// Error logs an error.
func (ll *Logger) Error(format string, v ...interface{}) {
	if LevelError > ll.level {
		return
	}
	msg := fmt.Sprintf("[E] "+format, v...)
	ll.Println(msg)
}

// Warning logs a warning.
func (ll *Logger) Warning(format string, v ...interface{}) {
	if LevelWarning > ll.level {
		return
	}
	msg := fmt.Sprintf("[W] "+format, v...)
	ll.Println(msg)
}

// Info logs an informational message.
func (ll *Logger) Info(format string, v ...interface{}) {
	if LevelInformational > ll.level {
		return
	}
	msg := fmt.Sprintf("[I] "+format, v...)
	ll.Println(msg)
}

// Debug logs a debug message.
func (ll *Logger) Debug(format string, v ...interface{}) {
	if LevelDebug > ll.level {
		return
	}
	msg := fmt.Sprintf("[D] "+format, v...)
	ll.Println(msg)
}

// BuildLogger builds the logger.
func BuildLogger(level string) {
	intLevel := LevelError
	switch level {
	case "error":
		intLevel = LevelError
	case "warning":
		intLevel = LevelWarning
	case "info":
		intLevel = LevelInformational
	case "debug":
		intLevel = LevelDebug
	}
	l := Logger{
		level: intLevel,
	}
	logger = &l
}

// Log returns the logger.
func Log() *Logger {
	if logger == nil {
		l := Logger{
			level: LevelDebug,
		}
		logger = &l
	}
	return logger
}
