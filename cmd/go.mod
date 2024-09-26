module github.com/smartcontractkit/ccip-owner-contracts/tools/signing

go 1.22.2

// Required for ledger library.
// NOTE MUST BE > 1.14 for this fix
// https://github.com/ethereum/go-ethereum/pull/28945
require (
	github.com/ethereum/go-ethereum v1.14.7
	github.com/smartcontractkit/ccip-owner-contracts v0.0.0-20240919155713-f4bf4ae0b9c6
)

require (
	github.com/bits-and-blooms/bitset v1.10.0 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/consensys/bavard v0.1.13 // indirect
	github.com/consensys/gnark-crypto v0.12.1 // indirect
	github.com/crate-crypto/go-kzg-4844 v1.0.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/ethereum/c-kzg-4844 v1.0.0 // indirect
	github.com/holiman/uint256 v1.3.0 // indirect
	github.com/mmcloughlin/addchain v0.4.0 // indirect
	github.com/supranational/blst v0.3.11 // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	rsc.io/tmplfunc v0.0.3 // indirect
)


