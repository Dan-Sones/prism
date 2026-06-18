#!/bin/sh
set -e

cat > /app/dist/config.js <<EOF
window.__APP_CONFIG__ = {
  apiBaseUrl: "${API_BASE_URL:-http://localhost:8081/api}",
};
EOF

exec serve -s dist
