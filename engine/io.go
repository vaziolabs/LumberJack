package engine

import (
	"os"
)

// IO utility functions
func NewFile(filename string) (*os.File, error) {
	file := "log.txt"
	if FileExists(file) {
		err := DeleteFile(file)
		if err != nil {
			Log(ErrorLevel, "Error deleting file: %v", err)
			return nil, err
		}
	}
	return os.Create(file)
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func DeleteFile(filename string) error {
	return os.Remove(filename)
}

// Default Config
var DefaultConfig = Config{
	Filename:    "log.txt",
	MinLogLevel: TraceLevel,
	ToConsole:   true,
}
