package loggger

import (
	"log"
)

// LogLevel represents the log level.
type LogLevel int

var Logs Logger

const (
	Debug LogLevel = iota
	Info
	Warn
	Error
)

// Logger defines the interface for logging.
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// ConfigurableLogger logs based on configuration.
type ConfigurableLogger struct {
	LogLevel LogLevel
}

// Debugf logs a debug message based on log level.
func (l ConfigurableLogger) Debugf(format string, args ...interface{}) {
	if l.LogLevel <= Debug {
		log.Printf("[DEBUG] "+format, args...)
	}
}

// Infof logs an info message based on log level.
func (l ConfigurableLogger) Infof(format string, args ...interface{}) {
	if l.LogLevel <= Info {
		log.Printf("[INFO] "+format, args...)
	}
}

// Warnf logs a warning message based on log level.
func (l ConfigurableLogger) Warnf(format string, args ...interface{}) {
	if l.LogLevel <= Warn {
		log.Printf("[WARN] "+format, args...)
	}
}

// Errorf logs an error message based on log level.
func (l ConfigurableLogger) Errorf(format string, args ...interface{}) {
	if l.LogLevel <= Error {
		log.Printf("[ERROR] "+format, args...)
	}
}
