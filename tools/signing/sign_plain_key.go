package signing

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	// NOTE MUST BE > 1.14 for this fix
	// https://github.com/ethereum/go-ethereum/pull/28945

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/executable"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/managed"
)

// Just run this locally to sign from the ledger.
func signPlainKey(privateKeyHex string) {
	// Load file
	proposal, _ := ProposalFromFile(managed.MCMSProposalTypeMap[os.Args[0]], os.Args[1])
	err := proposal.Validate()
	if err != nil {
		fmt.Println(err)
		return
	}

	executableProposal, err := proposal.ToExecutableMCMSProposal()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get the signing hash
	payload, err := executableProposal.SigningHash()
	if err != nil {
		log.Fatal(err)
	}

	// Load private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	// Sign the payload
	sig, err := crypto.Sign(payload, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Sign the payload
	unmarshalledSig := executable.Signature{}
	json.Unmarshal(sig, &unmarshalledSig)

	// Add signature to proposal
	proposal.AddSignature(unmarshalledSig)

	// Write proposal to file
	WriteProposalToFile(proposal, os.Args[0])
}
