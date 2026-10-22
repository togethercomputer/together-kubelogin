FROM ghcr.io/dexidp/dex:v2.46.1
# CVE-2026-111943: no dex release ships the fixed ca-certificates yet; upgrade it from the Alpine repo.
USER root
RUN apk upgrade --no-cache ca-certificates-bundle \
 && v="$(apk info -e -v ca-certificates-bundle | sed 's/^ca-certificates-bundle-//')" \
 && [ "$(apk version -t "$v" 20261020-r0)" != "<" ] \
 || { echo "ca-certificates-bundle $v is older than 20261020-r0" >&2; exit 1; }
USER dex:dex
