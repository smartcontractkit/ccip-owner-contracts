// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.13;

import {ManyChainMultiSig} from "../src/ManyChainMultiSig.sol";
import {ManyChainMultiSigBaseSetRootAndExecuteTest} from "./ManyChainMultiSigBaseTest.t.sol";

contract ManyChainMultiSigSetRootTest is ManyChainMultiSigBaseSetRootAndExecuteTest {
    event NewRoot(bytes32 indexed root, uint32 validUntil, ManyChainMultiSig.RootMetadata metadata);

    struct SetRootArgs {
        bytes32 root;
        uint32 validUntil;
        ManyChainMultiSig.RootMetadata metadata;
        bytes32[] metadataProof;
        ManyChainMultiSig.Signature[] signatures;
    }

    // Constructs validly signed SetRootArgs.
    // However, they don't contain an actual executable op.
    // Useful as a helper for testing setRoot logic.
    function makeSetRootArgs(uint32 validUntil, ManyChainMultiSig.RootMetadata memory metadata)
        internal
        returns (SetRootArgs memory)
    {
        ManyChainMultiSig.Op[] memory ops = new ManyChainMultiSig.Op[](1);
        ops[0] = s_testOps[1];

        // we didn't set overridePreviousRoot and there are previous ops
        bytes32[] memory leaves = constructLeaves(ops, metadata);

        (
            bytes32 root,
            bytes32[] memory metadataProof,
            ManyChainMultiSig.Signature[] memory signatures
        ) = constructAnsSignRootAndProof(leaves, validUntil, s_testPrivateKeys);

        return SetRootArgs({
            root: root,
            validUntil: validUntil,
            metadata: metadata,
            metadataProof: metadataProof,
            signatures: signatures
        });
    }

    function callSetRoot(SetRootArgs memory args) internal {
        s_testExposedManyChainMultiSig.setRoot(
            args.root, args.validUntil, args.metadata, args.metadataProof, args.signatures
        );
    }
}

contract ManyChainMultiSigSetRootSanityChecks is ManyChainMultiSigSetRootTest {
    function test_revertsOnInvalidChainID() public {
        ManyChainMultiSig.RootMetadata memory corruptedRootMetadata = s_initialTestRootMetadata;
        corruptedRootMetadata.chainId++;

        SetRootArgs memory args = makeSetRootArgs(s_testValidUntil, corruptedRootMetadata);
        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.WrongChainId.selector));
        callSetRoot(args);
    }

    function test_revertsOnInvalidMultiSig() public {
        ManyChainMultiSig.RootMetadata memory corruptedRootMetadata = s_initialTestRootMetadata;
        corruptedRootMetadata.multiSig = MULTISIG_OWNER;

        SetRootArgs memory args = makeSetRootArgs(s_testValidUntil, corruptedRootMetadata);
        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.WrongMultiSig.selector));
        callSetRoot(args);
    }

    function test_revertsOnIncorrectPreOpCount() public {
        // preOpCount > opCount
        {
            ManyChainMultiSig.RootMetadata memory corruptedRootMetadata = s_initialTestRootMetadata;
            corruptedRootMetadata.overridePreviousRoot = true;
            corruptedRootMetadata.preOpCount = s_testExposedManyChainMultiSig.getOpCount() + 1;

            SetRootArgs memory args = makeSetRootArgs(s_testValidUntil, corruptedRootMetadata);
            vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.WrongPreOpCount.selector));
            callSetRoot(args);
        }

        // opCount > preOpCount
        {
            ManyChainMultiSig.RootMetadata memory corruptedRootMetadata = s_initialTestRootMetadata;
            corruptedRootMetadata.overridePreviousRoot = true;

            s_testExposedManyChainMultiSig.setOpCount(corruptedRootMetadata.preOpCount + 1);
            SetRootArgs memory args = makeSetRootArgs(s_testValidUntil, corruptedRootMetadata);
            vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.WrongPreOpCount.selector));
            callSetRoot(args);
        }
    }

    function test_revertsOnIncorrectPostOpCount() public {
        callSetRoot(makeSetRootArgs(s_testValidUntil, s_initialTestRootMetadata));
        s_testExposedManyChainMultiSig.setOpCount(s_initialTestRootMetadata.postOpCount);

        ManyChainMultiSig.RootMetadata memory corruptedRootMetadata = s_initialTestRootMetadata;
        corruptedRootMetadata.preOpCount = s_initialTestRootMetadata.postOpCount;
        corruptedRootMetadata.postOpCount = corruptedRootMetadata.preOpCount - 1;

        SetRootArgs memory args = makeSetRootArgs(s_testValidUntil, corruptedRootMetadata);
        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.WrongPostOpCount.selector));
        callSetRoot(args);
    }

    function test_revertsOnExpiredValidUntil() public {
        vm.warp(s_testValidUntil + 1);

        SetRootArgs memory args = makeSetRootArgs(s_testValidUntil, s_initialTestRootMetadata);
        vm.expectRevert(
            abi.encodeWithSelector(ManyChainMultiSig.ValidUntilHasAlreadyPassed.selector)
        );
        callSetRoot(args);
    }

    function test_revertsOnRepeatedRootAndValidUntil() public {
        ManyChainMultiSig.RootMetadata memory rootMetadata = s_initialTestRootMetadata;
        rootMetadata.overridePreviousRoot = true;
        SetRootArgs memory args = makeSetRootArgs(s_testValidUntil, rootMetadata);
        callSetRoot(args);
        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.SignedHashAlreadySeen.selector));
        callSetRoot(args);
        // modify validUntil and setRoot works
        args = makeSetRootArgs(s_testValidUntil + 1, rootMetadata);
        callSetRoot(args);
    }

    function test_revertsOnNoConfig() public {
        // initialize the config to default values.
        s_testExposedManyChainMultiSig.clearConfig();

        SetRootArgs memory args = makeSetRootArgs(s_testValidUntil, s_initialTestRootMetadata);
        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.MissingConfig.selector));
        callSetRoot(args);
    }
}

contract ManyChainMultiSigSetOverrideRootTest is ManyChainMultiSigSetRootTest {
    function setUp() public override {
        ManyChainMultiSigBaseSetRootAndExecuteTest.setUp();
        // set an initial root
        vm.prank(MULTISIG_OWNER);
        s_testExposedManyChainMultiSig.setRoot(
            s_testInitialRoot,
            s_testValidUntil,
            s_initialTestRootMetadata,
            s_metadataProof,
            s_signatures
        );
        assertTrue(s_testExposedManyChainMultiSig.getOpCount() == 0);
    }

    function test_overridingRootSuccess() public {
        // we already have a root set during setUp(), which we're overriding
        // in this test
        s_initialTestRootMetadata.overridePreviousRoot = true;

        SetRootArgs memory args = makeSetRootArgs(s_testValidUntil, s_initialTestRootMetadata);
        vm.expectEmit(true, true, true, true);
        emit NewRoot(args.root, args.validUntil, args.metadata);
        callSetRoot(args);
        assertTrue(
            s_testExposedManyChainMultiSig.getOpCount() == s_initialTestRootMetadata.preOpCount
        );
    }

    function test_setRootAfterClearingSuccess() public {
        // "execute" all ops except one
        s_testExposedManyChainMultiSig.setOpCount(s_initialTestRootMetadata.postOpCount - 1);
        // setConfig with clearRoot = true
        vm.prank(MULTISIG_OWNER);
        s_testExposedManyChainMultiSig.setConfig(
            s_testSigners, s_signerGroups, s_testGroupQuorums, s_testGroupParents, true
        );

        ManyChainMultiSig.RootMetadata memory newRootMetadata = s_initialTestRootMetadata;
        newRootMetadata.preOpCount = s_initialTestRootMetadata.postOpCount - 1;
        callSetRoot(makeSetRootArgs(s_testValidUntil, newRootMetadata));
    }

    function test_executeOverrideRootExecute() public {
        {
            bytes32[] memory proof = computeProofForLeaf(s_testLeavesInTree, getLeafIndexOfOp(0));
            s_testExposedManyChainMultiSig.execute(s_testOps[0], proof);
        }

        ManyChainMultiSig.Op[] memory nextRootOps = new ManyChainMultiSig.Op[](1);
        nextRootOps[0] = s_testOps[1];

        ManyChainMultiSig.RootMetadata memory nextRootMetadata = ManyChainMultiSig.RootMetadata({
            chainId: s_testChainId,
            multiSig: s_multiSigContractAddress,
            preOpCount: 1,
            postOpCount: 2,
            overridePreviousRoot: false
        });

        {
            // we didn't set overridePreviousRoot and there are previous ops
            bytes32[] memory nextLeaves = constructLeaves(nextRootOps, nextRootMetadata);

            (
                bytes32 nextRoot,
                bytes32[] memory nextMetadataProof,
                ManyChainMultiSig.Signature[] memory nextSignatures
            ) = constructAnsSignRootAndProof(nextLeaves, s_testValidUntil, s_testPrivateKeys);

            vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.PendingOps.selector));
            s_testExposedManyChainMultiSig.setRoot(
                nextRoot, s_testValidUntil, nextRootMetadata, nextMetadataProof, nextSignatures
            );
        }

        {
            // set overridePreviousRoot, should work now
            nextRootMetadata.overridePreviousRoot = true;

            bytes32[] memory nextLeaves = constructLeaves(nextRootOps, nextRootMetadata);

            (
                bytes32 nextRoot,
                bytes32[] memory nextMetadataProof,
                ManyChainMultiSig.Signature[] memory nextSignatures
            ) = constructAnsSignRootAndProof(nextLeaves, s_testValidUntil, s_testPrivateKeys);

            s_testExposedManyChainMultiSig.setRoot(
                nextRoot, s_testValidUntil, nextRootMetadata, nextMetadataProof, nextSignatures
            );

            // execute
            bytes32[] memory proof = computeProofForLeaf(nextLeaves, getLeafIndexOfOp(0));
            s_testExposedManyChainMultiSig.execute(nextRootOps[0], proof);

            assertTrue(s_testExposedManyChainMultiSig.getOpCount() == nextRootMetadata.postOpCount);
        }
    }

    function test_revertsWhenNoOverrideAndThereIsPendingOps() public {
        ManyChainMultiSig.RootMetadata memory newRootMetadata = s_initialTestRootMetadata;
        // we don't override the root, so we expect a revert because there are still
        // unexecuted op in the previous root
        newRootMetadata.overridePreviousRoot = false;

        SetRootArgs memory args = makeSetRootArgs(s_testValidUntil, newRootMetadata);
        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.PendingOps.selector));
        callSetRoot(args);
    }

    function test_successWhenNoOverrideAfterEmptyRoot() public {
        ManyChainMultiSig.RootMetadata memory newRootMetadata = s_initialTestRootMetadata;
        // we already have a root set during setUp(), which we're overriding
        // in this test
        newRootMetadata.overridePreviousRoot = true;
        newRootMetadata.postOpCount = newRootMetadata.preOpCount;

        callSetRoot(makeSetRootArgs(s_testValidUntil, newRootMetadata));

        newRootMetadata.overridePreviousRoot = false; // there are no pending ops
        newRootMetadata.postOpCount = newRootMetadata.preOpCount + 1;

        callSetRoot(makeSetRootArgs(s_testValidUntil, newRootMetadata));
    }

    function test_successWhenNoOverrideAfterEverythingExecuted() public {
        assertLt(0, s_testExposedManyChainMultiSig.getRootMetadata().postOpCount);

        ManyChainMultiSig.RootMetadata memory newRootMetadata = s_initialTestRootMetadata;
        // we don't want to override the root.
        newRootMetadata.overridePreviousRoot = false;
        newRootMetadata.preOpCount = s_initialTestRootMetadata.postOpCount;
        newRootMetadata.postOpCount = newRootMetadata.preOpCount + 10;

        // advance opCount so now newRootMetadata.postOpCount == opCount, as
        // if everything in the previous root had been executed
        s_testExposedManyChainMultiSig.setOpCount(newRootMetadata.preOpCount);

        callSetRoot(makeSetRootArgs(s_testValidUntil, newRootMetadata));
    }
}

contract ManyChainMultiSigSetRootVerifyProofTest is ManyChainMultiSigSetRootTest {
    function test_failsWhenPostOpCountIsNotConsistentWithProof() public {
        // corrupted postOpCount. Now postOpCount is not consistent with s_metadataProof.
        s_initialTestRootMetadata.postOpCount++;
        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.ProofCannotBeVerified.selector));
        s_testExposedManyChainMultiSig.setRoot(
            s_testInitialRoot,
            s_testValidUntil,
            s_initialTestRootMetadata,
            s_metadataProof,
            s_signatures
        );
    }

    function test_failsWhenOverridePreviousRootIsNotConsistentWithProof() public {
        // now the new overridePreviousRoot is not consistent with newMetadataProof
        s_initialTestRootMetadata.overridePreviousRoot =
            !s_initialTestRootMetadata.overridePreviousRoot;

        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.ProofCannotBeVerified.selector));
        s_testExposedManyChainMultiSig.setRoot(
            s_testInitialRoot,
            s_testValidUntil,
            s_initialTestRootMetadata,
            s_metadataProof,
            s_signatures
        );
    }

    function test_failsWhenPreOpCountIsNotConsistentWithProof() public {
        // preOpCount is not consistent with s_metadataProof (which is the proof
        // of inclusion of s_initialTestRootMetadata)
        s_initialTestRootMetadata.preOpCount++;
        s_testExposedManyChainMultiSig.setOpCount(s_initialTestRootMetadata.preOpCount);
        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.ProofCannotBeVerified.selector));
        s_testExposedManyChainMultiSig.setRoot(
            s_testInitialRoot,
            s_testValidUntil,
            s_initialTestRootMetadata,
            s_metadataProof,
            s_signatures
        );
    }

    function test_failsWhenMultiSigIsNotConsistentWithProof() public {
        // multiSig is not consistent with s_metadataProof (which is the proof
        // of inclusion of s_initialTestRootMetadata)
        s_initialTestRootMetadata.multiSig =
            address(uint160(s_initialTestRootMetadata.multiSig) + 1);
        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.ProofCannotBeVerified.selector));
        s_testExposedManyChainMultiSig.setRoot(
            s_testInitialRoot,
            s_testValidUntil,
            s_initialTestRootMetadata,
            s_metadataProof,
            s_signatures
        );
    }

    function test_failsWhenChainIdIsNotConsistentWithProof() public {
        // chainId is not consistent with s_metadataProof (which is the proof
        // of inclusion of s_initialTestRootMetadata)
        s_initialTestRootMetadata.chainId++;
        vm.chainId(s_initialTestRootMetadata.chainId);
        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.ProofCannotBeVerified.selector));
        s_testExposedManyChainMultiSig.setRoot(
            s_testInitialRoot,
            s_testValidUntil,
            s_initialTestRootMetadata,
            s_metadataProof,
            s_signatures
        );
    }
}

contract ManyChainMultiSigSetRootVerifySignaturesTest is ManyChainMultiSigSetRootTest {
    function test_revertsOnInsufficientSignaturesNum() public {
        uint8 signersNum = 9;
        require(SIGNERS_NUM >= 9, "this test requires that there is at least 9 signers");

        // assign the required quorum in each group
        s_testGroupQuorums[0] = 3;
        s_testGroupQuorums[1] = 3;
        s_testGroupQuorums[2] = 2;
        s_testGroupQuorums[3] = 1;

        uint8[] memory signerGroups = new uint8[](signersNum);
        address[] memory signers = new address[](signersNum);
        // we assign 3 signers in each group, 0,1,2 in group 1, 3,4,5 in group 2, and
        // 6, 7, 8 in group 3.
        for (uint8 i = 0; i < signersNum; i++) {
            signers[i] = s_testSigners[i];
            signerGroups[i] = i / 3 + 1;
        }
        // we send 2 signatures from group 1, 1 from group 2, and 3 from group 3
        uint8 numSignatures = s_testGroupQuorums[1] - 1 + s_testGroupQuorums[2] - 1 + signersNum / 3;

        // set the new partition of signers to groups
        vm.prank(MULTISIG_OWNER);
        s_testExposedManyChainMultiSig.setConfig(
            signers, signerGroups, s_testGroupQuorums, s_testGroupParents, false
        );
        // we build signatures such that the we get 2 signatures from group 0, 1 from group 1, and 3 from group 2.
        ManyChainMultiSig.Signature[] memory signatures =
            new ManyChainMultiSig.Signature[](numSignatures);

        for (uint256 i = 0; i < signersNum / 3 - 1; i++) {
            signatures[i] = s_signatures[i];
        }
        signatures[signersNum / 3 - 1] = s_signatures[signersNum / 3];

        for (uint256 i = (2 * signersNum) / 3; i < signersNum; i++) {
            signatures[i - 3] = s_signatures[i];
        }
        // should revert as we have only 1 successful group
        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.InsufficientSigners.selector));
        s_testExposedManyChainMultiSig.setRoot(
            s_testInitialRoot,
            s_testValidUntil,
            s_initialTestRootMetadata,
            s_metadataProof,
            signatures
        );
    }

    function test_revertsOnNoSignatures() public {
        require(SIGNERS_NUM >= 2, "this test requires SIGNERS_NUM be at least 2");
        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.InsufficientSigners.selector));
        s_testExposedManyChainMultiSig.setRoot(
            s_testInitialRoot,
            s_testValidUntil,
            s_initialTestRootMetadata,
            s_metadataProof,
            new ManyChainMultiSig.Signature[](0)
        );
    }

    function test_revertsOnRepeatedSignatures() public {
        // there is two identical signatures
        require(SIGNERS_NUM >= 2, "this test requires SIGNERS_NUM be at least 2");
        s_signatures[0] = s_signatures[1];
        vm.expectRevert(
            abi.encodeWithSelector(
                ManyChainMultiSig.SignersAddressesMustBeStrictlyIncreasing.selector
            )
        );
        s_testExposedManyChainMultiSig.setRoot(
            s_testInitialRoot,
            s_testValidUntil,
            s_initialTestRootMetadata,
            s_metadataProof,
            s_signatures
        );
    }

    function test_revertsOnInvalidSignatureOnRoot() public {
        // modify a op leaf in the merkle tree and construct a new root
        s_testLeavesInTree[getLeafIndexOfOp(1)] = keccak256(abi.encode("1111111"));
        (
            bytes32 newRoot,
            bytes32[] memory newMetadataProof,
            /* signatures */
        ) = constructAnsSignRootAndProof(s_testLeavesInTree, s_testValidUntil, s_testPrivateKeys);

        // we send the old signatures' (the signatures on the old root), hence the
        // signatures are invalid.
        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.InvalidSigner.selector));
        s_testExposedManyChainMultiSig.setRoot(
            newRoot,
            s_testValidUntil,
            s_initialTestRootMetadata, // the rootMetadata didn't change
            newMetadataProof,
            s_signatures
        );
    }

    function test_revertsOnInconsistentValidUntilWithSignature() public {
        // we  pass a different validUntil to setRoot() instead of s_testValidUntil,
        // this should fail as (s_rs, s_ss, s_vs) are signatures on (s_testInitialRoot, s_testValidUntil)
        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.InvalidSigner.selector));
        s_testExposedManyChainMultiSig.setRoot(
            s_testInitialRoot,
            s_testValidUntil + 1,
            s_initialTestRootMetadata,
            s_metadataProof,
            s_signatures
        );
    }
}
