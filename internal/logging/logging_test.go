package logging

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestSetupLogger(t *testing.T) {
	t.Run("returns nil for valid log", func(t *testing.T) {
		tmpFile := createTmpFile(t)
		defer os.Remove(tmpFile.Name())
		LogPath = tmpFile.Name()

		if err := SetupLogger(); err != nil {
			t.Fatal("expected nil, got error")
		}
	})
}

func TestInfo(t *testing.T) {
	t.Run("returns nil for info log", func(t *testing.T) {
		tmpFile := createTmpFile(t)
		defer os.Remove(tmpFile.Name())
		LogPath = tmpFile.Name()

		if err := SetupLogger(); err != nil {
			t.Fatal("expected nil, got error")
		}

		Info("test info message")

		logLine, err := readLine()
		if err != nil {
			t.Fatal("expected nil, got error")
		}

		want := "[INFO] test info message"
		if !strings.HasSuffix(logLine, want) {
			t.Errorf("got %q, want %q", logLine, want)
		}
	})
}

func TestDebug(t *testing.T) {
	t.Run("returns nil for debug log", func(t *testing.T) {
		tmpFile := createTmpFile(t)
		defer os.Remove(tmpFile.Name())
		LogPath = tmpFile.Name()

		if err := SetupLogger(); err != nil {
			t.Fatal("expected nil, got error")
		}

		Debug("test debug message")

		logLine, err := readLine()
		if err != nil {
			t.Fatal("expected nil, got error")
		}

		want := "[DEBUG] test debug message"
		if !strings.HasSuffix(logLine, want) {
			t.Errorf("got %q, want %q", logLine, want)
		}
	})
}

func TestError(t *testing.T) {
	t.Run("returns nil for error log", func(t *testing.T) {
		tmpFile := createTmpFile(t)
		defer os.Remove(tmpFile.Name())
		LogPath = tmpFile.Name()

		if err := SetupLogger(); err != nil {
			t.Fatal("expected nil, got error")
		}

		Error("test error message")

		logLine, err := readLine()
		if err != nil {
			t.Fatal("expected nil, got error")
		}

		want := "[ERROR] test error message"
		if !strings.HasSuffix(logLine, want) {
			t.Errorf("got %q, want %q", logLine, want)
		}
	})
}

func createTmpFile(t *testing.T) *os.File {
	tmpFile, err := os.CreateTemp("", "bugbox-*.log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()

	return tmpFile
}

func readLine() (string, error) {
	data, err := os.ReadFile(LogPath)
	if err != nil {
		return "", fmt.Errorf("expected nil, got error")
	}
	lines := bytes.Split(data, []byte("\n"))
	return strings.TrimSpace(string(lines[len(lines)-2])), nil
}
