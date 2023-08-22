// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import "./ManyChainMultiSigBaseTest.t.sol";
import "../src/ManyChainMultiSig.sol";
import "../mock/ExposedManyChainMultiSig.sol";
import {Receiver} from "../mock/Receiver.sol";

contract ManyChainMultiSigDomainSeparationTest is ManyChainMultiSigBaseSetRootAndExecuteTest {
    function setUp() public virtual override {
        ManyChainMultiSigBaseTest.setUp();
    }

    function test_merkleTreePreimageDomainSeparation() public {
        // We store three kinds of items in the Merkle tree:
        // - inner nodes which are of size 64
        //   (see openzeppelin-contracts/contracts/utils/cryptography/MerkleProof.sol:15)
        // - RootMetadata
        // - Op
        // RootMetadata and Op are both hashed with a domain separator,
        // so we here we just need to ensure that both their pre-images have
        // a length different from that of inner nodes (64 bytes)

        ManyChainMultiSig.RootMetadata memory rootMetadata;
        bytes memory rootMetadataPreimage = leafMetadataPreimage(rootMetadata);
        assertLt(64, rootMetadataPreimage.length);

        ManyChainMultiSig.Op memory op;
        bytes memory opPreimage = leafOpPreimage(op);
        assertLt(64, opPreimage.length);
    }
}
