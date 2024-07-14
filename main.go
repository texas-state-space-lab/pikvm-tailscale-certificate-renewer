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

const timeToSleep = 24 * time.Hour

func main() {
	ctx := context.Background()

	for {
		domain, err := tailscale.GetDomain(ctx)
		if err != nil {
			slog.Error("failed to get domain", "error", err)
			continue
		}

		ssl := sslpaths.NewSSLPaths("/etc/kvmd/nginx/ssl/", domain)

		certManager := certmanager.NewCertManager(ssl)

		if err := certManager.CheckCert(); err != nil {
			if errors.Is(err, certmanager.ErrDoesNotExist) || errors.Is(err, certmanager.ErrExpiringSoon) {
				slog.Warn("cert is missing or expiring soon, generating new cert", "reason", err)

				if err := doCertRenewal(ctx, certManager, ssl); err != nil {
					slog.Error("failed to renew cert", "error", err)
				}
			} else {
				slog.Error("failed to check cert", "error", err, "cert_path", ssl.GetCertPath())
			}
		}

		slog.Info("sleeping", "duration", timeToSleep)
		time.Sleep(timeToSleep)
	}
}

func doCertRenewal(ctx context.Context, certManager *certmanager.CertManager, ssl *sslpaths.SSLPaths) error {
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
