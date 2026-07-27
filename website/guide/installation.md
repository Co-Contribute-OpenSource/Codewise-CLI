# Installation

## Requirements

- Go 1.25 or newer to build the current module
- Docker for image builds and pushes
- `kubectl` plus a configured cluster for Kubernetes operations
- Helm for Helm-based deployments

Only Go is required for local scaffolding, encoding, help, and version commands.

## Install a signed release

The installer supports Linux and macOS on x86-64 and ARM64. It verifies the
release SHA-256 checksum before extracting the binary. When Cosign is installed,
it also verifies the checksum file's keyless Sigstore signature.

Review and run the installer:

```bash
curl -fsSLO https://raw.githubusercontent.com/AryanSharma9917/Codewise-CLI/main/install.sh
less install.sh
bash install.sh
```

The default destination is `~/.local/bin`. Override it or select a release:

```bash
CODEWISE_INSTALL_DIR="$HOME/bin" CODEWISE_VERSION="v1.9.1" bash install.sh
```

Windows users should download the ZIP archive and `checksums.txt` from the
GitHub release page, verify it with `Get-FileHash -Algorithm SHA256`, and place
`codewise.exe` on `PATH`.

## Build from source

```bash
git clone https://github.com/aryansharma9917/codewise-cli.git
cd codewise-cli
make build
./codewise-cli version
```

You can also build directly:

```bash
go build -o codewise-cli main.go
```

## Install on your PATH

Move the compiled binary to a directory already included in your `PATH`:

```bash
sudo install -m 0755 codewise-cli /usr/local/bin/codewise
codewise --help
```

## Verify external tools

```bash
docker --version
kubectl version --client
helm version
```

Codewise reports a clear error when a command requires a missing external tool.
