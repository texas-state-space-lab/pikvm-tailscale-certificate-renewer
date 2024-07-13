package pikvm

import (
	"fmt"
	"os/exec"
)

// SetFSReadOnly sets the filesystem to read-only mode
func SetFSReadOnly() error {
	cmd := exec.Command("ro")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed enable read-only mode: %w", err)
	}

	return nil
}

// SetFSReadWrite sets the filesystem to read-write mode
func SetFSReadWrite() error {
	cmd := exec.Command("rw")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed enable read/write mode: %w", err)
	}

	return nil
}
