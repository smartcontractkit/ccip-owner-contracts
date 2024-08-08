#!/bin/bash

set -e

abigen_version="1.13.8"
abigen_package_path="github.com/ethereum/go-ethereum/cmd/abigen@v$abigen_version"

abigen() {
    go run "$abigen_package_path" "$@"
}

abigen --version | grep -F "abigen version $abigen_version-stable" >/dev/null || (
    echo "ASSERTION: abigen version should have been $abigen_version." 1>&2
    exit 1
)

abigen_owner_contracts() {
  abigen --pkg gethwrappers \
  --abi <(jq .abi ../out/"$1".sol/"$1".json) \
  --bin  <(jq --raw-output .bytecode.object ../out/"$1".sol/"$1".json) \
  --type "$1" \
  --out "$1".go
}

forge build
abigen_owner_contracts ManyChainMultiSig
abigen_owner_contracts RBACTimelock
abigen_owner_contracts CallProxy

