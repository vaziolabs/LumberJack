package types

import (
	"fmt"
	"log"
	"strings"
)

type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Notice(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Critical(format string, args ...interface{})
	Alert(format string, args ...interface{})
	Emergency(format string, args ...interface{})
	Success(format string, args ...interface{}) // For Testing Purposes
	Failure(format string, args ...interface{}) // For Testing Purposes
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

func (l *LogInfo) Debug(format string, args ...interface{}) {
	l.log("🔍", format, args...)
}

func (l *LogInfo) Notice(format string, args ...interface{}) {
	l.log("📝", format, args...)
}

func (l *LogInfo) Warn(format string, args ...interface{}) {
	l.log("⚠", format, args...)
}

func (l *LogInfo) Error(format string, args ...interface{}) {
	l.log("❌", format, args...)
}

func (l *LogInfo) Critical(format string, args ...interface{}) {
	l.log("🔥", format, args...)
}

func (l *LogInfo) Alert(format string, args ...interface{}) {
	l.log("🚨", format, args...)
}

func (l *LogInfo) Emergency(format string, args ...interface{}) {
	l.log("💀", format, args...)
}
