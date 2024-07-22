package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	systemd "github.com/coreos/go-systemd/v22/daemon"
	"github.com/texas-state-space-lab/pikvm-tailscale-certificate-renewer/internal/certmanager"
	"github.com/texas-state-space-lab/pikvm-tailscale-certificate-renewer/internal/pikvm"
	"github.com/texas-state-space-lab/pikvm-tailscale-certificate-renewer/internal/sslpaths"
	"github.com/texas-state-space-lab/pikvm-tailscale-certificate-renewer/internal/tailscale"
)

const (
	timeToSleep = 1 * time.Minute
)

func main() {
	ctx := context.Background()

	slog.Info("starting tailscale cert renewer")

	systemdOk, err := systemd.SdNotify(false, systemd.SdNotifyReady)
	if !systemdOk {
		if err != nil {
			slog.Error("failed to notify systemd, notify socket is unset")
			os.Exit(1)
		}

		slog.Error("failed to notify systemd", "error", err)
		os.Exit(1)
	}

	for {
		systemdOk, err := systemd.SdNotify(false, systemd.SdNotifyWatchdog)
		if !systemdOk {
			if err != nil {
				slog.Error("failed to notify systemd, notify socket is unset")
			}

			slog.Error("failed to notify systemd", "error", err)
		}

		if err := doCertCheckAndRenewal(ctx); err != nil {
			slog.Error("failed to check or renew cert", "time_until_retry", timeToSleep, "error", err)
		}

		time.Sleep(timeToSleep)
	}
}

func doCertCheckAndRenewal(ctx context.Context) error {
	domain, err := tailscale.GetDomain(ctx)
	if err != nil {
		return fmt.Errorf("failed to get tailscale domain: %w", err)
	}

	if domain == "" {
		return errors.New("tailscale domain is empty")
	}

	ssl := sslpaths.NewSSLPaths("/etc/kvmd/nginx/ssl/", domain)

	certManager := certmanager.NewCertManager(ssl)

	sslCertsChanged, err := changedSSLCerts(ctx, certManager)
	if err != nil {
		return fmt.Errorf("failed to check ssl certs: %w", err)
	}

	nginxConfigChanged, err := changedNginxConfig(ssl)
	if err != nil {
		return fmt.Errorf("failed to check nginx config: %w", err)
	}

	if sslCertsChanged || nginxConfigChanged {
		if err := pikvm.RestartNginx(); err != nil {
			return fmt.Errorf("failed to restart nginx: %w", err)
		}
	}

	return nil
}

func changedSSLCerts(ctx context.Context, certManager *certmanager.CertManager) (bool, error) {
	if err := certManager.CheckCert(ctx); errors.Is(err, certmanager.ErrCertDoesNotExist) ||
		errors.Is(err, certmanager.ErrKeyDoesNotExist) ||
		errors.Is(err, certmanager.ErrCertDoesNotMatch) ||
		errors.Is(err, certmanager.ErrKeyDoesNotMatch) {
		if err := certManager.GenerateCert(ctx); err != nil {
			return false, fmt.Errorf("failed to generate cert: %w", err)
		}

		return true, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to check cert: %w", err)
	}

	return false, nil
}

func changedNginxConfig(ssl *sslpaths.SSLPaths) (bool, error) {
	if err := pikvm.CheckNginxConfig(ssl); errors.Is(err, pikvm.ErrNginxConfigMissingSSLDetails) {
		if err := pikvm.WriteNginxConfig(ssl); err != nil {
			return false, fmt.Errorf("failed to write nginx config: %w", err)
		}

		return true, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to set certs in nginx config: %w", err)
	}

	return false, nil
}
