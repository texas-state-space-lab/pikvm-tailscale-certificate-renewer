#!/usr/bin/env bash

# PiKVM Tailscale Certificate Renewer Installer
# This script can be curled and piped to bash to install the latest version

set -e

function cleanup {
    # Set FS to read-only
    ro
}

trap cleanup EXIT

owner="texas-state-space-lab"
name="pikvm-tailscale-certificate-renewer"
repo="${owner}/${name}"

# later we may determine the architecture and download the correct binary
# right now all PiKVMs are armv7
tar_name="${name}_Linux_armv7"

# Get latest release
latest_release=$(curl -s "https://api.github.com/repos/${repo}/releases/latest" | grep "tag_name" | cut -d '"' -f 4)
echo "Latest release: ${latest_release}"

# Set FS to read/write
rw

# Download binary and move to /usr/local/bin
curl -L -s "https://github.com/${repo}/releases/download/${latest_release}/${tar_name}.tar.gz" -o /tmp/${tar_name}.tar.gz
tar -xzf /tmp/${tar_name}.tar.gz -C /tmp
mv "/tmp/${name}" /usr/local/bin/

# Download systemd service file and move to /etc/systemd/system
curl -L -s "https://raw.githubusercontent.com/${repo}/${latest_release}/${name}.service" -o "/etc/systemd/system/${name}.service"

# Reload systemd and enable/start the service
systemctl daemon-reload
systemctl enable "${name}"
systemctl start "${name}"

echo "Installed ${name} ${latest_release}"
