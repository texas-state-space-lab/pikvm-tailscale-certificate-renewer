package certmanager

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"reflect"

	"github.com/texas-state-space-lab/pikvm-tailscale-certificate-renewer/internal/pikvm"
	"github.com/texas-state-space-lab/pikvm-tailscale-certificate-renewer/internal/sslpaths"
	"github.com/texas-state-space-lab/pikvm-tailscale-certificate-renewer/internal/tailscale"
)

type CertManager struct {
	ssl *sslpaths.SSLPaths
}

func NewCertManager(ssl *sslpaths.SSLPaths) *CertManager {
	return &CertManager{
		ssl: ssl,
	}
}

const (
	certDirPerms  = 0o755
	certFilePerms = 0o644
)

var (
	ErrCertDoesNotExist = errors.New("cert does not exist")
	ErrKeyDoesNotExist  = errors.New("key does not exist")
	ErrCertDoesNotMatch = errors.New("cert does not match")
	ErrKeyDoesNotMatch  = errors.New("key does not match")
)

// CheckCert checks the cert and key files to see if they exist and match the tailscale cert
func (c *CertManager) CheckCert(ctx context.Context) error {
	if _, err := os.Stat(c.ssl.GetCertPath()); os.IsNotExist(err) {
		slog.Warn("cert file does not exist", "path", c.ssl.GetCertPath())

		return ErrCertDoesNotExist
	}

	if _, err := os.Stat(c.ssl.GetKeyPath()); os.IsNotExist(err) {
		slog.Warn("key file does not exist", "path", c.ssl.GetKeyPath())

		return ErrKeyDoesNotExist
	}

	tsCert, tsKey, err := tailscale.CertPair(ctx, c.ssl.GetDomain())
	if err != nil {
		return fmt.Errorf("failed to get tailscale cert pair: %w", err)
	}

	fsCert, err := os.ReadFile(c.ssl.GetCertPath())
	if err != nil {
		return fmt.Errorf("failed to read cert file: %w", err)
	}

	fsKey, err := os.ReadFile(c.ssl.GetKeyPath())
	if err != nil {
		return fmt.Errorf("failed to read key file: %w", err)
	}

	if !reflect.DeepEqual(tsCert, fsCert) {
		slog.Warn("tailscale and filesystem certs do not match", "path", c.ssl.GetCertPath())

		return ErrCertDoesNotMatch
	}

	if !reflect.DeepEqual(tsKey, fsKey) {
		slog.Warn("tailscale and filesystem keys do not match", "path", c.ssl.GetCertPath())

		return ErrKeyDoesNotMatch
	}

	return nil
}

// GenerateCert generates a new cert
func (c *CertManager) GenerateCert(ctx context.Context) error {
	cert, key, err := tailscale.CertPair(ctx, c.ssl.GetDomain())
	if err != nil {
		return fmt.Errorf("failed to get tailscale cert pair: %w", err)
	}

	if err := pikvm.SetFSReadWrite(); err != nil {
		return fmt.Errorf("failed filesystem mode change: %w", err)
	}

	defer func() {
		if err := pikvm.SetFSReadOnly(); err != nil {
			slog.Error("failed filesystem mode change", "error", err)
		}
	}()

	if _, err := os.Stat(c.ssl.GetDir()); os.IsNotExist(err) {
		if err := os.MkdirAll(c.ssl.GetDir(), certDirPerms); err != nil {
			return fmt.Errorf("failed to create cert path: %w", err)
		}
	}

	if err := os.WriteFile(c.ssl.GetCertPath(), cert, certFilePerms); err != nil {
		return fmt.Errorf("failed to write cert file: %w", err)
	}

	slog.Info("wrote cert file", "path", c.ssl.GetCertPath())

	if err := os.WriteFile(c.ssl.GetKeyPath(), key, certFilePerms); err != nil {
		return fmt.Errorf("failed to write key file: %w", err)
	}

	slog.Info("wrote key file", "path", c.ssl.GetKeyPath())

	return nil
}
