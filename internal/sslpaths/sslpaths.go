package sslpaths

import (
	"fmt"
	"path"
)

// SSLPaths is a struct that holds the paths to the SSL certificate and key
type SSLPaths struct {
	cert                string
	dir                 string
	domain              string
	key                 string
	nginxConfigCertLine string
	nginxConfigKeyLine  string
}

func NewSSLPaths(dir, domain string) *SSLPaths {
	sslP := &SSLPaths{
		cert:   path.Join(dir, domain+".crt"),
		dir:    dir,
		domain: domain,
		key:    path.Join(dir, domain+".key"),
	}
	sslP.nginxConfigCertLine = fmt.Sprintf("ssl_certificate %s;", sslP.GetCertPath())
	sslP.nginxConfigKeyLine = fmt.Sprintf("ssl_certificate_key %s;", sslP.GetKeyPath())

	return sslP
}

func (c *SSLPaths) GetCertPath() string {
	return c.cert
}

func (c *SSLPaths) GetKeyPath() string {
	return c.key
}

func (c *SSLPaths) GetDir() string {
	return c.dir
}

func (c *SSLPaths) GetDomain() string {
	return c.domain
}

func (c *SSLPaths) GetNginxConfigCertLine() string {
	return c.nginxConfigCertLine
}

func (c *SSLPaths) GetNginxConfigKeyLine() string {
	return c.nginxConfigKeyLine
}
