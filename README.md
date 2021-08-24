# Elvia Configuration Library (Golang)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Golang implementation of helper libraries for Elvia organization, which supports HashiCorp Vault, [ElvID](https://github.com/3lvia/elvid) and additional related services.

## Prerequisites

Setup required environment variables.

- `VAULT_ADDR`: The address of the vault. If it is not set (or empty), it will default to localhost.
- `GITHUB_TOKEN`: A GitHub token. If it is not set (or empty), it will default to use Kubernetes to login.
- `VAULT_CACERT`: If the Vault does not have a publicly signed CA certificate, you may set `VAULT_CACERT` as the file location of the self-signed certificate for the vault server with .pem format.

## Examples

See [example_test.go](example_test.go) file.

See also [demo.go](demo.go) for detailed example.
