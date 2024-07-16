# PiKVM Tailscale Cert Renewer

This is a tool to automatically renew tailscale certs for a PiKVM

This tool assumes you have setup your PiKVM and the [tailscale integration](https://docs.pikvm.org/tailscale/) using the [official docs](https://docs.pikvm.org/). This tool is designed around the following information from the docs:
>If you have a certificate (making a cert falls outside the scope of PiKVM - please reference OpenSSL documentation or use Let's Encrypt), replace keys in /etc/kvmd/nginx/ssl, edit /etc/kvmd/nginx/ssl.conf if necessary and restart kvmd-nginx service. *[PiKVM Common Questions](https://docs.pikvm.org/faq/#common-questions)*

This tool automatically discovers your tailscale domain, creates and renews certs for that domain, sets the cert path in the nginx config, and restarts NGINX.

```
curl -L -s "https://raw.githubusercontent.com/nateinaction/pikvm-tailscale-cert-renewer/main/install.sh" | bash
```
