package logging

import (
	"fmt"
	"time"
)

// Logger handles structured logging for the simulation
type Logger struct {
	enabled bool
}

// NewLogger creates a new Logger instance
func NewLogger(enabled bool) *Logger {
	return &Logger{enabled: enabled}
}

// LogTick logs the start of a new time tick
func (l *Logger) LogTick(tick int) {
	if !l.enabled {
		return
	}
	fmt.Printf("\n========== TICK %d [%s] ==========\n", tick, time.Now().Format("15:04:05"))
}

// LogEvent logs a general event
func (l *Logger) LogEvent(message string) {
	if !l.enabled {
		return
	}
	fmt.Printf("  %s\n", message)
}

// LogEvents logs multiple events
func (l *Logger) LogEvents(messages []string) {
	if !l.enabled {
		return
	}
	for _, msg := range messages {
		l.LogEvent(msg)
	}
}

// LogSummary logs a summary section
func (l *Logger) LogSummary(title string, data map[string]interface{}) {
	if !l.enabled {
		return
	}
	fmt.Printf("\n--- %s ---\n", title)
	for key, value := range data {
		fmt.Printf("  %s: %v\n", key, value)
	}
}

// LogError logs an error
func (l *Logger) LogError(err error) {
	if !l.enabled {
		return
	}
	fmt.Printf("  ❌ ERROR: %v\n", err)
}
