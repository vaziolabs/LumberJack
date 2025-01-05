package internal

import (
	"fmt"
	"log"
	"strings"
)

// Logger interface defines all logging methods
type Logger interface {
	Info(format string, args ...interface{})
	Success(format string, args ...interface{})
	Failure(format string, args ...interface{})
	Enter(name string)
	Exit(name string)
}

// ProductionLogger implements Logger for production use
type ProductionLogger struct {
	depth int
}

func NewLogger() *ProductionLogger {
	return &ProductionLogger{depth: 0}
}

func (l *ProductionLogger) getIndent() string {
	if l.depth < 0 {
		l.depth = 0
	}
	return strings.Repeat("│  ", l.depth)
}

func (l *ProductionLogger) log(prefix, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	log.Printf("%s%s %s", l.getIndent(), prefix, message)
}

func (l *ProductionLogger) Info(format string, args ...interface{}) {
	l.log("ℹ", format, args...)
}

func (l *ProductionLogger) Success(format string, args ...interface{}) {
	l.log("✓", format, args...)
}

func (l *ProductionLogger) Failure(format string, args ...interface{}) {
	l.log("✗", format, args...)
}

func (l *ProductionLogger) Enter(name string) {
	l.log("┌─", "BEGIN: %s", name)
	l.depth++
}

func (l *ProductionLogger) Exit(name string) {
	if l.depth > 0 {
		l.depth--
	}
	l.log("└─", "END: %s", name)
}
