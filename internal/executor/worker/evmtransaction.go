package worker

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"

	op "github.com/smartcontractkit/ccip-owner-contracts/internal/executor/operation"
	opdb "github.com/smartcontractkit/ccip-owner-contracts/internal/executor/operationdb"
	chainselectors "github.com/smartcontractkit/chain-selectors"
)

func SendTxn(ctx context.Context, o opdb.Operation, odb opdb.Store) {
	fmt.Printf("[%s] Starting Worker...", o.ID)
	// TODO:
	// privatekey should be obtained from context
	// rpc and eth client should be generated from an rpclib
	client, err := ethclient.Dial("")
	if err != nil {
		log.Fatal(err)
	}

	// [M1]TODO: check if operation is already complete, or check nonce, or check status
	chainID, err := chainselectors.ChainIdFromSelector(o.Operation.GetOpChainID())
	if err != nil {
		log.Fatal(err)
	}
	opData := o.Operation.GetOpData().(op.EvmOperation)

	// get private key
	privateKey, err := crypto.HexToECDSA("")
	if err != nil {
		log.Fatal(err)
	}

	// get nonce
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasLimit := uint64(21000) // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(0)

	tx := types.NewTransaction(nonce, opData.To, value, gasLimit, gasPrice, []byte(opData.Data))
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(int64(chainID))), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	ts := types.Transactions{signedTx}
	rawTxBytes, _ := rlp.EncodeToBytes(ts[0])
	rawTxHex := hex.EncodeToString(rawTxBytes)
	fmt.Printf("[%s] rawTxHex...%s", o.ID, rawTxHex)

	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[%s]tx sent: %s", o.ID, tx.Hash().Hex())

	fmt.Printf("[%s] Op finished, writing to store...", o.ID)
	newO := o
	newO.Status = opdb.Success
	newO.Timestamp = time.Now().Unix()
	_, err = odb.WriteToStore(newO)
	if err != nil {
		log.Fatal(err)
	}
}
