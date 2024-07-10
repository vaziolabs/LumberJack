package engine

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Level int

const (
	TraceLevel Level = iota // 0
	DebugLevel              // 1
	InfoLevel               // 2
	ErrorLevel              // 3
	FatalLevel              // 4
)

var (
	instance *Logger
	once     sync.Once
)

type Logger struct {
	mu          sync.Mutex
	logger      *log.Logger
	file        *os.File
	minLogLevel Level
	toConsole   bool
}

type Config struct {
	Filename    string
	MinLogLevel Level
	ToConsole   bool
}

func InitLogger(config Config) error {
	var err error
	once.Do(func() {
		instance, err = newLogger(config)
	})
	return err
}

func newLogger(config Config) (*Logger, error) {
	file, err := os.OpenFile(config.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	return &Logger{
		logger:      log.New(file, "", log.LstdFlags),
		file:        file,
		minLogLevel: config.MinLogLevel,
		toConsole:   config.ToConsole,
	}, nil
}

func Clear() {
	if instance != nil {
		instance.logger.SetOutput(os.Stdout)
	}
}

func Break() {
	println(" ")
}

func Log(level Level, format string, args ...interface{}) {
	if instance == nil {
		fmt.Println("Logger not initialized")
		err := InitLogger(DefaultConfig)

		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if level < instance.minLogLevel {
		return
	}

	instance.log(level, format, args...)
}

func (l *Logger) log(level Level, format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	message := fmt.Sprintf(format, args...)
	output := fmt.Sprintf("%s [%s] %s", time.Now().Format(time.RFC3339), levelToString(level), message)

	if l.toConsole {
		// Log to console
		fmt.Println(output)
	}

	if level >= InfoLevel {
		// Log to file
		l.logger.Println(output)
	}

	if level == FatalLevel {
		os.Exit(1)
	}
}

func Close() {
	if instance != nil && instance.file != nil {
		instance.file.Close()
	}
}

func levelToString(level Level) string {
	switch level {
	case TraceLevel:
		return "TRACE"
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}
