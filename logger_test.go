package observability

import (
	"os"
	"os/exec"
	"testing"
)

func TestNewLogger(t *testing.T) {
	t.Run("Default Config", func(t *testing.T) {
		cfg := &BaseConfig{
			ServiceName: "test-logger",
			LogLevel:    "info",
		}
		l := NewLogger(cfg)
		if l == nil {
			t.Fatal("NewLogger returned nil")
		}
		// Since l.SugaredLogger is embedded, we can't easily check internal level
		// without using unsafe or verifying behavior.
		// But verify methods don't panic
		l.Info("test info message")
	})

	t.Run("Nil Config", func(t *testing.T) {
		l := NewLogger(nil)
		if l == nil {
			t.Fatal("NewLogger(nil) returned nil")
		}
		l.Info("test info message with nil config")
	})

	t.Run("Invalid Log Level Falls back to Info", func(t *testing.T) {
		// This unit test is tricky because NewLogger logic:
		// if parsed, err := zapcore.ParseLevel(cfg.LogLevel); err == nil { level = parsed }
		// if err != nil, it keeps default 'InfoLevel'.
		// We want to confirm it doesn't crash.
		cfg := &BaseConfig{
			ServiceName: "test-logger",
			LogLevel:    "invalid", // Should trigger error in zapcore.ParseLevel
		}
		l := NewLogger(cfg)
		if l == nil {
			t.Fatal("NewLogger returned nil")
		}
		l.Info("should be info level")
	})
}

// Check if Logger implements expected methods
func TestLoggerMethods(t *testing.T) {
	cfg := &BaseConfig{
		ServiceName: "test-methods",
		LogLevel:    "debug",
	}
	l := NewLogger(cfg)

	// Just calling them to Ensure no panics
	l.Debug("debug msg", "key", "val")
	l.Info("info msg", "key", "val")
	l.Warn("warn msg", "key", "val")
	l.Error("error msg", "key", "val")
	// l.Fatal will exit the program, so we skip it or mock os.Exit if possible (hard with Zap)
}

func TestLoggerSync(t *testing.T) {
	cfg := &BaseConfig{
		ServiceName: "test-sync",
		LogLevel:    "info",
	}
	l := NewLogger(cfg)

	// Ensure Sync does not panic and returns (it's a void method)
	l.Sync()
}

// This test verifies that calling Fatal results in process exit. It runs
// the helper in a subprocess so the test binary itself doesn't exit.
func TestLoggerFatalSubprocess(t *testing.T) {
	if os.Getenv("TEST_FATAL_HELPER") == "1" {
		// In helper mode: create logger and call Fatal which should exit
		cfg := &BaseConfig{ServiceName: "helper-fatal"}
		l := NewLogger(cfg)
		l.Fatal("fatal called from helper")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestLoggerFatalSubprocess")
	cmd.Env = append(os.Environ(), "TEST_FATAL_HELPER=1")
	err := cmd.Run()
	if err == nil {
		t.Fatalf("expected subprocess to exit with non-zero status")
	}
	// We expect an ExitError when the subprocess calls Fatal/exit
	if _, ok := err.(*exec.ExitError); !ok {
		t.Fatalf("expected *exec.ExitError, got %T: %v", err, err)
	}
}
