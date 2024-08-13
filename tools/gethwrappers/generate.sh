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

improved_abigen () {
  jq .abi ../../out/$1.sol/$1.json > ../../out/$1.sol/$1.abi
  jq --raw-output .bytecode.object ../../out/$1.sol/$1.json > ../../out/$1.sol/$1.bin
  go run ./generate/wrap.go ../../out/$1.sol/$1.abi ../../out/$1.sol/$1.bin $1 gethwrappers
}

forge build
improved_abigen CallProxy
improved_abigen ManyChainMultiSig
improved_abigen RBACTimelock
