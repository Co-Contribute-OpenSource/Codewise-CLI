#!/usr/bin/env bash
set -euo pipefail

repository="AryanSharma9917/Codewise-CLI"
install_dir="${CODEWISE_INSTALL_DIR:-$HOME/.local/bin}"
requested_version="${CODEWISE_VERSION:-latest}"

case "$(uname -s)" in
  Linux) operating_system="Linux" ;;
  Darwin) operating_system="Darwin" ;;
  *)
    echo "Unsupported operating system. Download a release manually." >&2
    exit 1
    ;;
esac

case "$(uname -m)" in
  x86_64|amd64) architecture="x86_64" ;;
  arm64|aarch64) architecture="arm64" ;;
  *)
    echo "Unsupported architecture: $(uname -m)" >&2
    exit 1
    ;;
esac

asset="codewise_${operating_system}_${architecture}.tar.gz"
if [[ "$requested_version" == "latest" ]]; then
  release_url="https://github.com/${repository}/releases/latest/download"
else
  release_url="https://github.com/${repository}/releases/download/${requested_version}"
fi

temporary_dir="$(mktemp -d)"
trap 'rm -rf "$temporary_dir"' EXIT

echo "Downloading ${asset}..."
curl --fail --location --silent --show-error \
  --output "$temporary_dir/$asset" "$release_url/$asset"
curl --fail --location --silent --show-error \
  --output "$temporary_dir/checksums.txt" "$release_url/checksums.txt"

expected_line="$(awk -v asset="$asset" '$2 == asset { print; exit }' "$temporary_dir/checksums.txt")"
if [[ -z "$expected_line" ]]; then
  echo "Release checksum does not contain ${asset}" >&2
  exit 1
fi

if command -v sha256sum >/dev/null 2>&1; then
  (cd "$temporary_dir" && printf '%s\n' "$expected_line" | sha256sum --check --strict -)
else
  expected_hash="${expected_line%% *}"
  actual_hash="$(shasum -a 256 "$temporary_dir/$asset" | awk '{print $1}')"
  if [[ "$expected_hash" != "$actual_hash" ]]; then
    echo "Checksum verification failed for ${asset}" >&2
    exit 1
  fi
  echo "${asset}: OK"
fi

if command -v cosign >/dev/null 2>&1; then
  curl --fail --location --silent --show-error \
    --output "$temporary_dir/checksums.txt.sigstore.json" \
    "$release_url/checksums.txt.sigstore.json"
  cosign verify-blob \
    --bundle "$temporary_dir/checksums.txt.sigstore.json" \
    --certificate-identity-regexp "https://github.com/${repository}/.github/workflows/release.yml@refs/tags/v.*" \
    --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
    "$temporary_dir/checksums.txt"
else
  echo "Cosign is not installed; SHA-256 verified, Sigstore verification skipped."
fi

tar -xzf "$temporary_dir/$asset" -C "$temporary_dir"
if [[ ! -f "$temporary_dir/codewise" ]]; then
  echo "Archive does not contain the codewise binary" >&2
  exit 1
fi

mkdir -p "$install_dir"
install -m 0755 "$temporary_dir/codewise" "$install_dir/codewise"

echo "Installed codewise to $install_dir/codewise"
if [[ ":$PATH:" != *":$install_dir:"* ]]; then
  echo "Add $install_dir to PATH to invoke codewise directly."
fi
"$install_dir/codewise" version
