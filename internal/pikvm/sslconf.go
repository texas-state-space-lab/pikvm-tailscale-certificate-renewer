package pikvm

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/nateinaction/pikvm-tailscale-cert-renewer/internal/sslpaths"
)

var (
	certlineRegex = regexp.MustCompile(`^ssl_certificate\s+.*`)
	keylineRegex  = regexp.MustCompile(`^ssl_certificate_key\s+.*`)

	ErrNginxConfigMissingSSLDetails = fmt.Errorf("nginx config missing ssl details")
)

const (
	nginxSSLConf      = "/etc/kvmd/nginx/ssl.conf"
	nginxSSLConfPerms = 0o644
)

// CheckNginxConfig checks the nginx config for the cert and key lines
func CheckNginxConfig(ssl *sslpaths.SSLPaths) error {
	b, err := os.ReadFile(nginxSSLConf)
	if err != nil {
		return fmt.Errorf("failed to read ssl.conf: %w", err)
	}

	lines := strings.Split(string(b), "\n")

	if !slices.Contains(lines, ssl.GetNginxConfigCertLine()) ||
		!slices.Contains(lines, ssl.GetNginxConfigKeyLine()) {
		slog.Warn("cert or key line not found in nginx config", "path", nginxSSLConf)

		return ErrNginxConfigMissingSSLDetails
	}

	return nil
}

// WriteNginxConfig writes the cert and key lines to the nginx config
func WriteNginxConfig(ssl *sslpaths.SSLPaths) error {
	b, err := os.ReadFile(nginxSSLConf)
	if err != nil {
		return fmt.Errorf("failed to read ssl.conf: %w", err)
	}

	lines := strings.Split(string(b), "\n")

	lines = setLine(lines, certlineRegex, ssl.GetNginxConfigCertLine())
	lines = setLine(lines, keylineRegex, ssl.GetNginxConfigKeyLine())

	if err := SetFSReadWrite(); err != nil {
		return fmt.Errorf("failed filesystem mode change: %w", err)
	}

	defer func() {
		if err := SetFSReadOnly(); err != nil {
			slog.Error("failed filesystem mode change", "error", err)
		}
	}()

	if err := os.WriteFile(nginxSSLConf, []byte(strings.Join(lines, "\n")), nginxSSLConfPerms); err != nil {
		return fmt.Errorf("failed to write to nginx ssl config at %s: %w", nginxSSLConf, err)
	}

	slog.Info("wrote to nginx ssl config", "path", nginxSSLConf)

	return nil
}

// setLine sets the line in the contents if the regex matches
func setLine(contents []string, regex *regexp.Regexp, certLine string) []string {
	for i, line := range contents {
		if regex.MatchString(line) {
			contents[i] = certLine
			return contents
		}
	}

	return append(contents, fmt.Sprintf("%s\n", certLine))
}
