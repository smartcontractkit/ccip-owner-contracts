#!/bin/bash

set -e

# notice: set the envvar GO_BINPATH to the path of a go v1.22
# binary; for instance:
# GO_BINPATH="/opt/homebrew/Cellar/go@1.22/1.22.12/bin/go" ./generate.sh

abigen_owner_contracts() {
  jq .abi ../out/"$1".sol/"$1".json > ../out/"$1".sol/"$1".abi
  jq --raw-output .bytecode.object ../out/"$1".sol/"$1".json > ../out/"$1".sol/"$1".bin
  go run ./abigen/... ../out/"$1".sol/"$1".abi ../out/"$1".sol/"$1".bin "$1" gethwrappers
}

forge build
abigen_owner_contracts ManyChainMultiSig
abigen_owner_contracts RBACTimelock
abigen_owner_contracts CallProxy
