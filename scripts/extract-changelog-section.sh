#!/usr/bin/env bash

set -euo pipefail

CHANGELOG_FILE=""
VERSION=""

usage() {
    echo "Uso: $(basename "$0") --file <CHANGELOG.md> --version <vX.Y.Z|X.Y.Z>" >&2
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        --file)
            CHANGELOG_FILE="$2"
            shift 2
            ;;
        --version)
            VERSION="$2"
            shift 2
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo "Argumento no soportado: $1" >&2
            usage
            exit 1
            ;;
    esac
done

if [[ -z "$CHANGELOG_FILE" || -z "$VERSION" ]]; then
    usage
    exit 1
fi

if [[ ! -f "$CHANGELOG_FILE" ]]; then
    echo "No existe el changelog: $CHANGELOG_FILE" >&2
    exit 1
fi

VERSION_NO_V="${VERSION#v}"

# One awk pass: print from the "## [X.Y.Z]" heading (also "## [vX.Y.Z]") up to the line before
# the next "## [" heading, or to EOF for the last section. The version is compared as a literal
# string, not as a regex, so "0.900.1" never matches "0.90011". Trailing blank lines are dropped.
# Exit code 3 = version not found (a plain "exit" inside a rule would still run END).
if awk -v version="$VERSION_NO_V" '
    BEGIN { plain = "## [" version "]"; prefixed = "## [v" version "]" }
    /^## \[/ {
        if (found) { exit }
        if (index($0, plain) == 1 || index($0, prefixed) == 1) { found = 1 }
    }
    found {
        if ($0 ~ /^[[:space:]]*$/) { pending++; next }
        for (; pending > 0; pending--) { print "" }
        print
    }
    END { if (!found) { exit 3 } }
' "$CHANGELOG_FILE"; then
    exit 0
fi

echo "No se encontro la version $VERSION_NO_V en $CHANGELOG_FILE" >&2
exit 1
