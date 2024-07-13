package tailscale

import (
	"context"
	"fmt"

	ts "tailscale.com/client/tailscale"
)

// GetDomain gets the domain of the current machine from the tailscale client
func GetDomain(ctx context.Context) (string, error) {
	client := &ts.LocalClient{}

	statusResp, err := client.Status(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get status: %w", err)
	}

	return statusResp.Self.DNSName, nil
}

// CertPair generates the cert pair for the given domain
func CertPair(ctx context.Context, domain string) ([]byte, []byte, error) {
	client := &ts.LocalClient{}

	cert, key, err := client.CertPair(ctx, domain)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get cert pair: %w", err)
	}

	return cert, key, nil
}
