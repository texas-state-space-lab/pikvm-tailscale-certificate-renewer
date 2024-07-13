package main

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/nateinaction/tailscale-cert-renewer/internal/certmanager"
	"github.com/nateinaction/tailscale-cert-renewer/internal/pikvm"
	"github.com/nateinaction/tailscale-cert-renewer/internal/sslpaths"
	"github.com/nateinaction/tailscale-cert-renewer/internal/tailscale"
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
				slog.Info("cert is missing or expiring soon, generating new cert", "reason", err)

				if err := pikvm.SetFSReadWrite(); err != nil {
					slog.Error("failed filesystem mode change", "error", err)
					continue
				}

				genCert(ctx, certManager)
			} else {
				slog.Error("failed to check cert", "error", err)
			}
		}

		time.Sleep(timeToSleep)
	}
}

func genCert(ctx context.Context, certManager *certmanager.CertManager) {
	defer func() {
		if err := pikvm.SetFSReadOnly(); err != nil {
			slog.Error("failed filesystem mode change", "error", err)
		}
	}()

	if err := certManager.GenerateCert(ctx); err != nil {
		slog.Error("failed to generate cert", "error", err)
	}
}
