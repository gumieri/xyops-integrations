#!/bin/sh
set -e

systemctl daemon-reload
systemctl enable --now xyops-mcp
