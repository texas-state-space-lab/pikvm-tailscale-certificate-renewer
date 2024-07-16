package pikvm

import (
	"fmt"
	"log/slog"
	"os/exec"
)

// RestartNginx restarts the kvmd-nginx service
func RestartNginx() error {
	out, err := exec.Command("systemctl", "restart", "kvmd-nginx").Output()
	if err != nil {
		return fmt.Errorf("failed restart kvmd-nginx with output: %s: %w", out, err)
	}

	slog.Info("kvmd-nginx restarted")

	return nil
}
