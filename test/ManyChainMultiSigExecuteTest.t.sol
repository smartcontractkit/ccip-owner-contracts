// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.13;

import {ManyChainMultiSig} from "../src/ManyChainMultiSig.sol";
import {ManyChainMultiSigBaseSetRootAndExecuteTest} from "./ManyChainMultiSigBaseTest.t.sol";

contract ManyChainMultiSigExecuteTest is ManyChainMultiSigBaseSetRootAndExecuteTest {
    event OpExecuted(uint40 indexed nonce, address to, bytes data, uint256 value);

    function setUp() public override {
        ManyChainMultiSigBaseSetRootAndExecuteTest.setUp();
        // set an initial root
        s_testExposedManyChainMultiSig.setRoot(
            s_testInitialRoot,
            s_testValidUntil,
            s_initialTestRootMetadata,
            s_metadataProof,
            s_signatures
        );
    }

    function test_revertsOnPostOpCountReached() public {
        s_testExposedManyChainMultiSig.setOpCount(s_initialTestRootMetadata.postOpCount);
        // op and proof aren't even checked by the contract in this situation
        ManyChainMultiSig.Op memory fakeOp = s_testOps[0];
        bytes32[] memory fakeProof = new bytes32[](0);
        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.PostOpCountReached.selector));
        s_testExposedManyChainMultiSig.execute(fakeOp, fakeProof);
    }

    function test_revertsOnBadProof() public {
        // modify the first op
        ManyChainMultiSig.Op memory t = s_testOps[0];
        t.value++;

        bytes32[] memory emptyProof;

        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.ProofCannotBeVerified.selector));
        s_testExposedManyChainMultiSig.execute(t, emptyProof);

        // send a proof for the original op before the modification
        bytes32[] memory proof = computeProofForLeaf(s_testLeavesInTree, 0);
        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.ProofCannotBeVerified.selector));
        s_testExposedManyChainMultiSig.execute(t, proof);
    }

    function test_revertsOnBadOpData() public {
        bytes32[] memory proof = new bytes32[](5);
        ManyChainMultiSig.Op memory t1 = s_testOps[0];

        t1.chainId++;
        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.WrongChainId.selector));
        s_testExposedManyChainMultiSig.execute(t1, proof);

        ManyChainMultiSig.Op memory t2 = s_testOps[0];
        t2.multiSig = MULTISIG_OWNER;
        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.WrongMultiSig.selector));
        s_testExposedManyChainMultiSig.execute(t2, proof);

        ManyChainMultiSig.Op memory t3 = s_testOps[0];
        t3.nonce++;
        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.WrongNonce.selector));

        s_testExposedManyChainMultiSig.execute(t3, proof);

        ManyChainMultiSig.Op memory t4 = s_testOps[0];
        vm.warp(s_testValidUntil + 1);
        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.RootExpired.selector));

        s_testExposedManyChainMultiSig.execute(t4, proof);
    }

    function test_opsAreExecutedInOrder() public {
        bytes32[] memory proof1 = computeProofForLeaf(s_testLeavesInTree, getLeafIndexOfOp(0));
        vm.expectEmit(true, true, true, true);
        emit OpExecuted(s_testOps[0].nonce, s_testOps[0].to, s_testOps[0].data, s_testOps[0].value);
        vm.expectCall(s_testOps[0].to, s_testOps[0].value, s_testOps[0].data);
        s_testExposedManyChainMultiSig.execute(s_testOps[0], proof1);
        assertEq(s_testExposedManyChainMultiSig.getOpCount(), s_testOps[0].nonce + 1);

        // try to re-execute the op
        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.WrongNonce.selector));
        s_testExposedManyChainMultiSig.execute(s_testOps[0], proof1);

        // try to execute the third op instead of the second
        bytes32[] memory proof2 = computeProofForLeaf(s_testLeavesInTree, getLeafIndexOfOp(2));
        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.WrongNonce.selector));
        s_testExposedManyChainMultiSig.execute(s_testOps[2], proof2);

        // execute the second op
        bytes32[] memory proof3 = computeProofForLeaf(s_testLeavesInTree, getLeafIndexOfOp(1));
        vm.expectEmit(true, true, true, true);
        emit OpExecuted(s_testOps[1].nonce, s_testOps[1].to, s_testOps[1].data, s_testOps[1].value);
        vm.expectCall(s_testOps[1].to, s_testOps[1].value, s_testOps[1].data);
        s_testExposedManyChainMultiSig.execute(s_testOps[1], proof3);
        assertEq(s_testExposedManyChainMultiSig.getOpCount(), s_testOps[1].nonce + 1);
    }

    function test_revertsOnFailedOp() public {
        s_testExposedManyChainMultiSig.setOpCount(uint40(REVERTING_OP_INDEX));

        bytes32[] memory proof =
            computeProofForLeaf(s_testLeavesInTree, getLeafIndexOfOp(REVERTING_OP_INDEX));

        ManyChainMultiSig.Op memory t = s_testOps[REVERTING_OP_INDEX];

        (bool success, bytes memory expectedRet) = t.to.call{value: t.value}(t.data);
        assertFalse(success);
        vm.expectRevert(
            abi.encodeWithSelector(ManyChainMultiSig.CallReverted.selector, expectedRet)
        );
        s_testExposedManyChainMultiSig.execute(t, proof);
    }

    function test_value() public {
        assertEq(address(s_testExposedManyChainMultiSig).balance, 0);

        s_testExposedManyChainMultiSig.setOpCount(s_testOps[VALUE_OP_INDEX].nonce);
        bytes32[] memory proof =
            computeProofForLeaf(s_testLeavesInTree, getLeafIndexOfOp(VALUE_OP_INDEX));

        // No ether present in ManyChainMultiSig
        vm.expectRevert();
        s_testExposedManyChainMultiSig.execute(s_testOps[VALUE_OP_INDEX], proof);

        // Send 1 wei to ManyChainMultiSig
        vm.deal(address(this), 1);
        (bool success,) = address(s_testExposedManyChainMultiSig).call{value: 1}("");
        assertTrue(success);
        assertEq(address(s_testExposedManyChainMultiSig).balance, 1);
        assertEq(s_testOps[VALUE_OP_INDEX].to.balance, 0);

        // Execute op sending 1 wei to receiver
        s_testExposedManyChainMultiSig.execute(s_testOps[VALUE_OP_INDEX], proof);
        assertEq(address(s_testExposedManyChainMultiSig).balance, 0);
        assertEq(s_testOps[VALUE_OP_INDEX].to.balance, 1);
    }
}
