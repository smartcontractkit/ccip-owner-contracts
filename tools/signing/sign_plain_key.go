package signing

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	// NOTE MUST BE > 1.14 for this fix
	// https://github.com/ethereum/go-ethereum/pull/28945

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/mcms"
)

// Just run this locally to sign from the ledger.
func signPlainKey(privateKeyHex string) {
	// Load file
	proposal, _ := ProposalFromFile(os.Args[0])
	err := proposal.Validate()
	if err != nil {
		fmt.Println(err)
		return
	}

	executor, err := proposal.ToExecutor() // TODO: pass in a real backend
	if err != nil {
		log.Fatal(err)
	}

	// Get the signing hash
	payload, err := executor.SigningHash()
	if err != nil {
		log.Fatal(err)
	}

	// Load private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	// Sign the payload
	sig, err := crypto.Sign(payload.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Sign the payload
	unmarshalledSig := mcms.Signature{}
	err = json.Unmarshal(sig, &unmarshalledSig)
	if err != nil {
		log.Fatal(err)
	}

	// Add signature to proposal
	proposal.Signatures = append(proposal.Signatures, unmarshalledSig)

	// Write proposal to file
	WriteProposalToFile(proposal, os.Args[0])
}
