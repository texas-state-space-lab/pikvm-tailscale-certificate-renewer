package pikvm

import (
	"fmt"
	"os/exec"
)

// SetFSReadOnly sets the filesystem to read-only mode
func SetFSReadOnly() error {
	out, err := exec.Command("ro").Output()
	if err != nil {
		return fmt.Errorf("failed enable read-only mode with output: %s: %w", out, err)
	}

	return nil
}

// SetFSReadWrite sets the filesystem to read-write mode
func SetFSReadWrite() error {
	out, err := exec.Command("rw").Output()
	if err != nil {
		return fmt.Errorf("failed enable read/write mode with output: %s: %w", out, err)
	}

	return nil
}
