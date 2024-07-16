package sslpaths_test

import (
	"testing"

	"github.com/nateinaction/pikvm-tailscale-cert-renewer/internal/sslpaths"
)

func TestNewSSLPaths(t *testing.T) {
	t.Parallel()

	dir := "/path/to/dir"
	domain := "example.com"

	sslP := sslpaths.NewSSLPaths(dir, domain)

	if sslP.GetCertPath() != "/path/to/dir/example.com.crt" {
		t.Errorf("Expected cert path '/path/to/dir/example.com.crt', but got '%s'", sslP.GetCertPath())
	}

	if sslP.GetKeyPath() != "/path/to/dir/example.com.key" {
		t.Errorf("Expected key path '/path/to/dir/example.com.key', but got '%s'", sslP.GetKeyPath())
	}

	if sslP.GetDir() != "/path/to/dir" {
		t.Errorf("Expected dir '/path/to/dir', but got '%s'", sslP.GetDir())
	}

	if sslP.GetDomain() != "example.com" {
		t.Errorf("Expected domain 'example.com', but got '%s'", sslP.GetDomain())
	}

	if sslP.GetNginxConfigCertLine() != "ssl_certificate /path/to/dir/example.com.crt;" {
		t.Errorf("Expected nginx config cert line 'ssl_certificate /path/to/dir/example.com.crt;', but got '%s'",
			sslP.GetNginxConfigCertLine())
	}

	if sslP.GetNginxConfigKeyLine() != "ssl_certificate_key /path/to/dir/example.com.key;" {
		t.Errorf("Expected nginx config key line 'ssl_certificate_key /path/to/dir/example.com.key;', but got '%s'",
			sslP.GetNginxConfigKeyLine())
	}
}
