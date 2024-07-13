package certmanager

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/nateinaction/tailscale-cert-renewer/internal/sslpaths"
	"github.com/nateinaction/tailscale-cert-renewer/internal/tailscale"
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
	closeToExpire = -7 * 24 * time.Hour
)

var (
	ErrExpiringSoon = errors.New("connection error")
	ErrDoesNotExist = errors.New("cert does not exist")
)

// CheckCert checks if the cert exists and is not expiring soon
func (c *CertManager) CheckCert() error {
	_, err := os.Stat(c.ssl.GetCertPath())
	if errors.Is(err, os.ErrNotExist) {
		return ErrDoesNotExist
	}

	if err != nil {
		return fmt.Errorf("failed to stat cert file: %w", err)
	}

	b, err := os.ReadFile(c.ssl.GetCertPath())
	if err != nil {
		return fmt.Errorf("failed to read cert file: %w", err)
	}

	cert, err := x509.ParseCertificate(b)
	if err != nil {
		return fmt.Errorf("failed to parse cert: %w", err)
	}

	renewIfAfter := time.Now().Add(closeToExpire)
	if cert.NotAfter.After(renewIfAfter) {
		return ErrExpiringSoon
	}

	return nil
}

// GenerateCert generates a new cert
func (c *CertManager) GenerateCert(ctx context.Context) error {
	cert, key, err := tailscale.CertPair(ctx, c.ssl.GetDomain())
	if err != nil {
		return fmt.Errorf("failed to get tailscale cert pair: %w", err)
	}

	if _, err := os.Stat(c.ssl.GetDir()); os.IsNotExist(err) {
		if err := os.MkdirAll(c.ssl.GetDir(), certDirPerms); err != nil {
			return fmt.Errorf("failed to create cert path: %w", err)
		}
	}

	if err := os.WriteFile(c.ssl.GetCertPath(), cert, certFilePerms); err != nil {
		return fmt.Errorf("failed to write cert file: %w", err)
	}

	if err := os.WriteFile(c.ssl.GetKeyPath(), key, certFilePerms); err != nil {
		return fmt.Errorf("failed to write key file: %w", err)
	}

	return nil
}
