// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Make sure solidity compiler artifacts are up-to-date. Only output stdout on failure.
//go:generate ./generation/compile_contracts.sh

//go:generate go run ./generation/generate/wrap.go abi/CallProxy/CallProxy.abi abi/CallProxy/CallProxy.bin CallProxy call_proxy_wrapper
//go:generate go run ./generation/generate/wrap.go abi/ManyChainMultiSig/ManyChainMultiSig.abi abi/ManyChainMultiSig/ManyChainMultiSig.bin ManyChainMultiSig many_chain_multi_sig_wrapper
//go:generate go run ./generation/generate/wrap.go abi/RBACTimelock/RBACTimelock.abi abi/RBACTimelock/RBACTimelock.bin RBACTimelock rbac_timelock_wrapper
