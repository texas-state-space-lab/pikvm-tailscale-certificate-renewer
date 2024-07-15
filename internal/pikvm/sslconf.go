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
)

const (
	nginxSSLConf      = "/etc/kvmd/nginx/ssl.conf"
	nginxSSLConfPerms = 0o644
)

// SetCertsInNginxConfig sets the certificates in the nginx config file
func SetCertsInNginxConfig(ssl *sslpaths.SSLPaths) error {
	certLine := fmt.Sprintf("ssl_certificate %s;", ssl.GetCertPath())
	keyLine := fmt.Sprintf("ssl_certificate_key %s;", ssl.GetKeyPath())

	b, err := os.ReadFile(nginxSSLConf)
	if err != nil {
		return fmt.Errorf("failed to read ssl.conf: %w", err)
	}

	lines := strings.Split(string(b), "\n")

	if !slices.Contains(lines, certLine) || !slices.Contains(lines, keyLine) {
		slog.Warn("cert or key line not found in nginx config, adding")

		if err := writeNginxConfig(lines, certLine, keyLine); err != nil {
			return fmt.Errorf("failed to write nginx config: %w", err)
		}
	}

	return nil
}

// writeNginxConfig writes the cert and key lines to the nginx config
func writeNginxConfig(lines []string, certLine string, keyLine string) error {
	lines = setLine(lines, certlineRegex, certLine)
	lines = setLine(lines, keylineRegex, keyLine)

	if err := SetFSReadWrite(); err != nil {
		return fmt.Errorf("failed filesystem mode change: %w", err)
	}

	defer func() {
		if err := SetFSReadOnly(); err != nil {
			slog.Error("failed filesystem mode change", "error", err)
		}
	}()

	if err := os.WriteFile(nginxSSLConf, []byte(strings.Join(lines, "\n")), nginxSSLConfPerms); err != nil {
		return fmt.Errorf("failed to write ssl.conf: %w", err)
	}

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
