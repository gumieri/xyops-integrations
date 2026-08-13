#!/bin/sh
set -e

systemctl stop --no-block xyops-mcp
systemctl disable xyops-mcp 2>/dev/null || true
rm -f /etc/xyops/xyops-mcp.env
