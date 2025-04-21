package logging

import (
	"log"
	"os"
	"path/filepath"
)

var Logger *log.Logger

var LogPath = filepath.Join(os.Getenv("HOME"), ".local", "share", "bugbox", "bugbox.log")

// SetupLogger sets up loggers to write to bugbox.log
func SetupLogger() error {
	dir := filepath.Dir(LogPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	file, err := os.OpenFile(LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	Logger = log.New(file, "[BugBox] ", log.Ldate|log.Ltime)
	return nil
}

// Info logs an info-level message
func Info(message string) {
	if Logger != nil {
		Logger.Println("[INFO] " + message)
	}
}

// Error logs an error-level message
func Error(message string) {
	if Logger != nil {
		Logger.Println("[ERROR] " + message)
	}
}
