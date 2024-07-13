package sslpaths

import "path"

// SSLPaths is a struct that holds the paths to the SSL certificate and key
type SSLPaths struct {
	cert   string
	dir    string
	domain string
	key    string
}

func NewSSLPaths(dir, domain string) *SSLPaths {
	return &SSLPaths{
		cert:   path.Join(dir, domain+".crt"),
		dir:    dir,
		domain: domain,
		key:    path.Join(dir, domain+".key"),
	}
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
