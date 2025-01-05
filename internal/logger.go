package internal

import (
	"fmt"
	"log"
	"strings"
)

type Logger interface {
	Info(format string, args ...interface{})
	Success(format string, args ...interface{})
	Failure(format string, args ...interface{})
	Enter(name string)
	Exit(name string)
}

type LogInfo struct {
	depth int
}

func NewLogger() *LogInfo {
	return &LogInfo{depth: 0}
}

func (l *LogInfo) getIndent() string {
	if l.depth < 0 {
		l.depth = 0
	}
	return strings.Repeat("│  ", l.depth)
}

func (l *LogInfo) log(prefix, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	log.Printf("%s%s %s", l.getIndent(), prefix, message)
}

func (l *LogInfo) Info(format string, args ...interface{}) {
	l.log("ℹ", format, args...)
}

func (l *LogInfo) Success(format string, args ...interface{}) {
	l.log("✓", format, args...)
}

func (l *LogInfo) Failure(format string, args ...interface{}) {
	l.log("✗", format, args...)
}

func (l *LogInfo) Enter(name string) {
	l.log("┌─", "BEGIN: %s", name)
	l.depth++
}

func (l *LogInfo) Exit(name string) {
	if l.depth > 0 {
		l.depth--
	}
	l.log("└─", "END: %s", name)
}
