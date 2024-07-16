package pikvm

import (
	"fmt"
	"log/slog"
	"os/exec"
)

// SetFSReadOnly sets the filesystem to read-only mode
func SetFSReadOnly() error {
	out, err := exec.Command("ro").Output()
	if err != nil {
		return fmt.Errorf("failed enable read-only mode with output: %s: %w", out, err)
	}

	slog.Info("filesystem mode changed to read-only")

	return nil
}

// SetFSReadWrite sets the filesystem to read-write mode
func SetFSReadWrite() error {
	out, err := exec.Command("rw").Output()
	if err != nil {
		return fmt.Errorf("failed enable read/write mode with output: %s: %w", out, err)
	}

	slog.Info("filesystem mode changed to read/write")

	return nil
}
