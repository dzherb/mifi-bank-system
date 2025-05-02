#!/usr/bin/env bash
set -e

PROJECT_ROOT="$(dirname "$(realpath "$0")")/../"

export $(grep DATABASE_URL "$PROJECT_ROOT/.env")

eval $(python3 - <<EOF
import os
from urllib.parse import urlparse

p = urlparse(os.environ["DATABASE_URL"])
print(f"""
user={p.username}
pass={p.password}
host={p.hostname}
port={p.port or 5432}
dbname={p.path[1:]}
""")
EOF
)

export PGPASSWORD=$pass

npx pg-mermaid \
--host $host \
--port $port \
--dbname $dbname \
--username $user \
--excluded-tables schema_migrations \
--output-path "$PROJECT_ROOT/database.md"
