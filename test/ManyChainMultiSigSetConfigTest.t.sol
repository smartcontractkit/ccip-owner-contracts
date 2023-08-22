// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.13;

import {ManyChainMultiSig} from "../src/ManyChainMultiSig.sol";
import "./ManyChainMultiSigBaseTest.t.sol";

contract ManyChainMultiSigSetConfigTest is ManyChainMultiSigBaseTest {
    event ConfigSet(ManyChainMultiSig.Config config, bool isRootCleared);

    function test_revertsOnNonOwnerCaller() public {
        // call sefConfig from non-owner address
        vm.expectRevert("Ownable: caller is not the owner");
        vm.prank(EXTERNAL_CALLER);
        s_testExposedManyChainMultiSig.setConfig(
            s_testSigners, s_signerGroups, s_testGroupQuorums, s_testGroupParents, false
        );
    }

    function test_revertsOnInvalidConfig() public {
        // setConfig must be called from the multiSig owner only
        vm.startPrank(MULTISIG_OWNER);

        // signer's list must not be empty
        {
            address[] memory emptySignersList;
            uint8[] memory emptySignerGroupsList;
            vm.expectRevert(
                abi.encodeWithSelector(ManyChainMultiSig.OutOfBoundsNumOfSigners.selector)
            );
            s_testExposedManyChainMultiSig.setConfig(
                emptySignersList,
                emptySignerGroupsList,
                s_testGroupQuorums,
                s_testGroupParents,
                false
            );
        }

        // signers must be distinct
        {
            address[] memory signers = s_testSigners;

            signers[1] = signers[0];
            vm.expectRevert(
                abi.encodeWithSelector(
                    ManyChainMultiSig.SignersAddressesMustBeStrictlyIncreasing.selector
                )
            );
            s_testExposedManyChainMultiSig.setConfig(
                signers, s_signerGroups, s_testGroupQuorums, s_testGroupParents, false
            );
        }

        // test that setConfig reverts on out of bounds signer's group
        {
            uint8[] memory localSignerGroups = s_signerGroups;
            localSignerGroups[0] = MAX_NUM_GROUPS + 1;
            vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.OutOfBoundsGroup.selector));

            s_testExposedManyChainMultiSig.setConfig(
                s_testSigners, localSignerGroups, s_testGroupQuorums, s_testGroupParents, false
            );
        }

        // test that setConfig on too large group quorum
        {
            uint8[MAX_NUM_GROUPS] memory localGroupQuorums = s_testGroupQuorums;
            localGroupQuorums[0] = SIGNERS_NUM + 1;
            vm.expectRevert(
                abi.encodeWithSelector(ManyChainMultiSig.OutOfBoundsGroupQuorum.selector)
            );

            s_testExposedManyChainMultiSig.setConfig(
                s_testSigners, s_signerGroups, localGroupQuorums, s_testGroupParents, false
            );
        }

        // test that setConfig on non well-formed group tree: root doesn't have itself as parent
        {
            uint8[MAX_NUM_GROUPS] memory localGroupParents = s_testGroupParents;
            localGroupParents[0] = 1;
            vm.expectRevert(
                abi.encodeWithSelector(ManyChainMultiSig.GroupTreeNotWellFormed.selector)
            );

            s_testExposedManyChainMultiSig.setConfig(
                s_testSigners, s_signerGroups, s_testGroupQuorums, localGroupParents, false
            );
        }

        // test that setConfig reverts on non well-formed group tree: some non-root group has itself
        // as parent
        {
            uint8[MAX_NUM_GROUPS] memory localGroupParents = s_testGroupParents;
            localGroupParents[1] = 1;
            vm.expectRevert(
                abi.encodeWithSelector(ManyChainMultiSig.GroupTreeNotWellFormed.selector)
            );

            s_testExposedManyChainMultiSig.setConfig(
                s_testSigners, s_signerGroups, s_testGroupQuorums, localGroupParents, false
            );
        }

        // test that setConfig reverts on signer being included in disabled group
        {
            uint8[] memory localSignerGroups = s_signerGroups;
            localSignerGroups[1] = MAX_NUM_GROUPS - 1;
            vm.expectRevert(
                abi.encodeWithSelector(ManyChainMultiSig.SignerInDisabledGroup.selector)
            );

            s_testExposedManyChainMultiSig.setConfig(
                s_testSigners, localSignerGroups, s_testGroupQuorums, s_testGroupParents, false
            );
        }

        // test that setConfig reverts when signers.length != signerGroups.length
        {
            address[] memory signers = new address[](4);
            uint8[] memory signerGroups = new uint8[](3);
            vm.expectRevert(
                abi.encodeWithSelector(ManyChainMultiSig.SignerGroupsLengthMismatch.selector)
            );

            s_testExposedManyChainMultiSig.setConfig(
                signers, signerGroups, s_testGroupQuorums, s_testGroupParents, false
            );

            signerGroups = new uint8[](2);
            vm.expectRevert(
                abi.encodeWithSelector(ManyChainMultiSig.SignerGroupsLengthMismatch.selector)
            );

            s_testExposedManyChainMultiSig.setConfig(
                signers, signerGroups, s_testGroupQuorums, s_testGroupParents, false
            );
        }
    }

    function test_success() public {
        // setConfig must be called from the multiSig owner only
        vm.startPrank(MULTISIG_OWNER);
        vm.expectEmit(true, true, true, true);
        emit ConfigSet(s_testConfig, false);
        s_testExposedManyChainMultiSig.setConfig(
            s_testSigners, s_signerGroups, s_testGroupQuorums, s_testGroupParents, false
        );
        ManyChainMultiSig.Config memory config = s_testExposedManyChainMultiSig.getConfig();
        assertEq(abi.encode(config.groupParents), abi.encode(s_testGroupParents));
        assertEq(config.signers.length, s_testSigners.length);
        for (uint256 i = 0; i < s_testSigners.length; i++) {
            assertEq(s_testSigners[i], config.signers[i].addr);
        }

        // test clear root
        emit ConfigSet(s_testConfig, true);
        s_testExposedManyChainMultiSig.setConfig(
            s_testSigners, s_signerGroups, s_testGroupQuorums, s_testGroupParents, true
        );
        (bytes32 root, uint32 validUntil) = s_testExposedManyChainMultiSig.getRoot();
        assertEq(root, 0);
        assertEq(validUntil, 0);
        ManyChainMultiSig.RootMetadata memory rootMetadata =
            s_testExposedManyChainMultiSig.getRootMetadata();
        assertEq(rootMetadata.chainId, block.chainid);
        assertEq(rootMetadata.multiSig, address(s_testExposedManyChainMultiSig));
        assertEq(rootMetadata.preOpCount, s_testExposedManyChainMultiSig.getOpCount());
        assertEq(rootMetadata.postOpCount, s_testExposedManyChainMultiSig.getOpCount());
        assertEq(rootMetadata.overridePreviousRoot, true);
    }
}
