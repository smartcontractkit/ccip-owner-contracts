package main

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"sort"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// The method in this file is used in the foundry tests for computing the test
// addresses and their corresponding private keys. The addresses are sorted in
// ascending order as required in ManyChainMultiSig.sol.

var (
	uint256Array, _ = abi.NewType("uint256[]", "[]uint256",
		[]abi.ArgumentMarshaling{
			{Name: "val", Type: "uint256[]"},
		},
	)
	addressesArray, _ = abi.NewType("address[]", "[]address",
		[]abi.ArgumentMarshaling{
			{Name: "val", Type: "address[]"},
		},
	)
	encodingSigners = abi.Arguments{
		{Type: uint256Array, Name: "privateKeys"},
		{Type: addressesArray, Name: "signers"},
	}
)

type bytesArrayForSort [][]byte

func (a bytesArrayForSort) Len() int {
	return len(a)
}

func (a bytesArrayForSort) Less(i, j int) bool {
	switch bytes.Compare(a[i], a[j]) {
	case -1:
		return true
	case 0, 1:
		return false
	default:
		log.Panic("compare failed")
		return false
	}
}

func (a bytesArrayForSort) Swap(i, j int) {
	a[j], a[i] = a[i], a[j]
}

// main receives the required number of signers' addresses and prints the encoded
// (in HEX) list of private keys and (sorted) list of addresses. The encoding according to
// the supported format in Solidity.
func main() {
	if len(os.Args) < 2 {
		panic("should pass the required number of addresses")
	}
	var addresses [][]byte
	addressesToPrivateKeys := make(map[common.Address]*big.Int)

	numKeys, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	for i := 0; i < numKeys; i++ {
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			panic(err)
		}
		publicKey := privateKey.Public()
		publicKey, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			panic("publicKey is not of type *ecdsa.PublicKey")
		}
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			panic("publicKey is not of type *ecdsa.PublicKey")
		}
		address := crypto.PubkeyToAddress(*publicKeyECDSA)
		addresses = append(addresses, address[:])
		addressesToPrivateKeys[address] = privateKey.D
	}

	sort.Sort(bytesArrayForSort(addresses))

	var sortedPrivateKeys []*big.Int
	var sortedAddresses []common.Address
	for _, addrInBytes := range addresses {
		addr := common.BytesToAddress(addrInBytes)
		sortedPrivateKeys = append(sortedPrivateKeys, addressesToPrivateKeys[addr])
		sortedAddresses = append(sortedAddresses, addr)
	}
	packed, err := encodingSigners.Pack(
		sortedPrivateKeys, sortedAddresses)
	if err != nil {
		panic(err)
	}
	// Must NOT print a new line
	fmt.Print(common.Bytes2Hex(packed))

}
