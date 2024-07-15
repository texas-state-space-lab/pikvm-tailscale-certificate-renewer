package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/nateinaction/pikvm-tailscale-cert-renewer/internal/certmanager"
	"github.com/nateinaction/pikvm-tailscale-cert-renewer/internal/pikvm"
	"github.com/nateinaction/pikvm-tailscale-cert-renewer/internal/sslpaths"
	"github.com/nateinaction/pikvm-tailscale-cert-renewer/internal/tailscale"
)

const (
	timeToSleep = 1 * time.Minute
)

func main() {
	ctx := context.Background()

	for {
		if err := doCertCheckAndRenewal(ctx); err != nil {
			slog.Error("failed to check or renew cert", "time_until_retry", timeToSleep, "error", err)
			time.Sleep(timeToSleep)

			continue
		}

		time.Sleep(timeToSleep)
	}
}

func doCertCheckAndRenewal(ctx context.Context) error {
	domain, err := tailscale.GetDomain(ctx)
	if err != nil {
		return fmt.Errorf("failed to get tailscale domain: %w", err)
	}

	ssl := sslpaths.NewSSLPaths("/etc/kvmd/nginx/ssl/", domain)

	certManager := certmanager.NewCertManager(ssl)

	if err := certManager.CheckCert(ctx); !errors.Is(err, certmanager.ErrCertDoesNotExist) &&
		!errors.Is(err, certmanager.ErrKeyDoesNotExist) &&
		!errors.Is(err, certmanager.ErrCertDoesNotMatch) &&
		!errors.Is(err, certmanager.ErrKeyDoesNotMatch) {
		return fmt.Errorf("failed to check cert: %w", err)
	}

	if err := certManager.GenerateCert(ctx); err != nil {
		return fmt.Errorf("failed to generate cert: %w", err)
	}

	if err := pikvm.SetCertsInNginxConfig(ssl); err != nil {
		return fmt.Errorf("failed to set certs in nginx config: %w", err)
	}

	if err := pikvm.RestartNginx(); err != nil {
		return fmt.Errorf("failed to restart nginx: %w", err)
	}

	return nil
}
