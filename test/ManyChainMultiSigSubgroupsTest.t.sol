pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import "../mock/Counter.sol";
import "../src/ManyChainMultiSig.sol";
import "./ManyChainMultiSigBaseTest.t.sol";
import "./MerkleHelper.sol";

contract ManyChainMultiSigSubgroupsTest is Test {
    uint8 constant MCMS_NUM_GROUPS = 32;

    MerkleHelper s_merkleHelper = new MerkleHelper();
    ManyChainMultiSigBaseTest s_manyChainMultiSigBaseTest;

    uint8 constant NUM_SIGNERS = 20;

    address[] s_signerAddresses;
    uint256[] s_signerPrivateKeys;

    ManyChainMultiSig s_multisig;

    function setUp() public virtual {
        s_manyChainMultiSigBaseTest = new ManyChainMultiSigBaseTest();
        (s_signerAddresses, s_signerPrivateKeys) =
            s_manyChainMultiSigBaseTest.addressesWithPrivateKeys(NUM_SIGNERS);
        s_multisig = new ManyChainMultiSig();
    }

    function randomBetween(bytes32 randomState, uint8 lower, uint8 upper)
        internal
        pure
        returns (uint8, bytes32)
    {
        uint256 range = upper - lower;
        uint256 n = uint256(randomState);
        uint8 randomNumber = uint8((n % range) + lower);
        return (randomNumber, keccak256(bytes.concat(randomState)));
    }

    function removeIndex(ManyChainMultiSig.Signature[] memory signatures, uint8 index)
        internal
        pure
        returns (ManyChainMultiSig.Signature[] memory)
    {
        if (index >= signatures.length) {
            return signatures;
        }

        ManyChainMultiSig.Signature[] memory result =
            new ManyChainMultiSig.Signature[](signatures.length-1);
        uint256 offset = 0;
        for (uint256 i = 0; i < signatures.length; i++) {
            if (i != index) {
                result[offset] = signatures[i];
                offset++;
            }
        }

        return result;
    }

    function test_setConfig_chain() public {
        uint8[] memory signerGroups = new uint8[](NUM_SIGNERS);

        uint8[MCMS_NUM_GROUPS] memory groupQuorums;
        uint8[MCMS_NUM_GROUPS] memory groupParents;

        // all signers are in the last group
        for (uint256 i = 0; i < signerGroups.length; i++) {
            signerGroups[i] = MCMS_NUM_GROUPS - 1;
        }
        // form a chain of groups from the last group to the root
        for (uint8 i = 0; i < MCMS_NUM_GROUPS; i++) {
            if (i != 0) {
                groupParents[i] = i - 1;
            }
            groupQuorums[i] = 1;
        }
        groupQuorums[MCMS_NUM_GROUPS - 1] = NUM_SIGNERS - 1;

        s_multisig.setConfig(s_signerAddresses, signerGroups, groupQuorums, groupParents, false);

        ManyChainMultiSig.Op[] memory ops = new ManyChainMultiSig.Op[](1);
        (MerkleHelper.SetRootArgs memory setRootArgs,) = s_merkleHelper.build(
            s_signerPrivateKeys,
            uint32(block.timestamp + 2 hours),
            ManyChainMultiSig.RootMetadata({
                chainId: block.chainid,
                multiSig: address(s_multisig),
                preOpCount: 0,
                postOpCount: 1,
                overridePreviousRoot: true
            }),
            ops
        );

        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.InsufficientSigners.selector));
        s_multisig.setRoot(
            setRootArgs.root,
            setRootArgs.validUntil,
            setRootArgs.metadata,
            setRootArgs.metadataProof,
            removeIndex(removeIndex(setRootArgs.signatures, 0), 0)
        );

        s_multisig.setRoot(
            setRootArgs.root,
            setRootArgs.validUntil,
            setRootArgs.metadata,
            setRootArgs.metadataProof,
            removeIndex(setRootArgs.signatures, 0)
        );
    }

    function test_setConfig_fuzz(bytes32 randomState) public {
        uint8[] memory signerGroups = new uint8[](NUM_SIGNERS);

        uint8[MCMS_NUM_GROUPS] memory groupChildrenCounts;
        uint8[MCMS_NUM_GROUPS] memory groupQuorums;
        uint8[MCMS_NUM_GROUPS] memory groupParents;

        for (uint256 i = 0; i < signerGroups.length; i++) {
            (signerGroups[i], randomState) = randomBetween(randomState, 0, MCMS_NUM_GROUPS);
            groupChildrenCounts[signerGroups[i]]++;
        }

        for (uint8 j = 0; j < MCMS_NUM_GROUPS; j++) {
            uint8 i = MCMS_NUM_GROUPS - 1 - j;

            if (groupChildrenCounts[i] == 0) continue;

            (groupQuorums[i], randomState) = randomBetween(randomState, 0, groupChildrenCounts[i]);
            groupQuorums[i]++;

            if (i != 0) {
                (groupParents[i], randomState) = randomBetween(randomState, 0, i);
                groupChildrenCounts[groupParents[i]]++;
            }
        }

        bool allSignersNeeded = true;
        for (uint8 i = 0; i < MCMS_NUM_GROUPS; i++) {
            allSignersNeeded = allSignersNeeded && groupQuorums[i] == groupChildrenCounts[i];
        }

        s_multisig.setConfig(s_signerAddresses, signerGroups, groupQuorums, groupParents, false);

        {
            ManyChainMultiSig.Op[] memory ops = new ManyChainMultiSig.Op[](1);
            (MerkleHelper.SetRootArgs memory setRootArgs,) = s_merkleHelper.build(
                s_signerPrivateKeys,
                uint32(block.timestamp + 2 hours),
                ManyChainMultiSig.RootMetadata({
                    chainId: block.chainid,
                    multiSig: address(s_multisig),
                    preOpCount: 0,
                    postOpCount: 1,
                    overridePreviousRoot: false
                }),
                ops
            );

            if (!allSignersNeeded) {
                bool success = false;
                // can remove at least some signature and setRoot still works
                for (uint8 i = 0; i < setRootArgs.signatures.length; i++) {
                    try s_multisig.setRoot(
                        setRootArgs.root,
                        setRootArgs.validUntil,
                        setRootArgs.metadata,
                        setRootArgs.metadataProof,
                        removeIndex(setRootArgs.signatures, i)
                    ) {
                        success = true;
                        break;
                    } catch {}
                }
                assertTrue(success);
            }
        }

        {
            ManyChainMultiSig.Op[] memory ops = new ManyChainMultiSig.Op[](1);
            (MerkleHelper.SetRootArgs memory setRootArgs,) = s_merkleHelper.build(
                s_signerPrivateKeys,
                uint32(block.timestamp + 2 hours),
                ManyChainMultiSig.RootMetadata({
                    chainId: block.chainid,
                    multiSig: address(s_multisig),
                    preOpCount: 0,
                    postOpCount: 1,
                    overridePreviousRoot: true
                }),
                ops
            );

            s_multisig.setRoot(
                setRootArgs.root,
                setRootArgs.validUntil,
                setRootArgs.metadata,
                setRootArgs.metadataProof,
                setRootArgs.signatures
            );
        }
    }

    function test_setConfig_c4issue16() public {
        uint8[] memory signerGroups = new uint8[](NUM_SIGNERS);
        signerGroups[0] = MCMS_NUM_GROUPS - 1; // put one signer in last group

        uint8[MCMS_NUM_GROUPS] memory groupQuorums = [
            1,
            1,
            1,
            1,
            1,
            1,
            1,
            1,
            1,
            1,
            1,
            1,
            1,
            1,
            1,
            1,
            1,
            1,
            1,
            1,
            1,
            1,
            1,
            1,
            1,
            1,
            1,
            1,
            1,
            1,
            1,
            1
        ];
        uint8[MCMS_NUM_GROUPS] memory groupParents = [
            1,
            2,
            3,
            4,
            5,
            6,
            7,
            8,
            9,
            10,
            11,
            12,
            13,
            14,
            15,
            16,
            17,
            18,
            19,
            20,
            21,
            22,
            23,
            24,
            25,
            26,
            27,
            28,
            29,
            30,
            31,
            32
        ];

        // The following setConfig reverts, contradicting the issue author's claim.
        vm.expectRevert(abi.encodeWithSelector(ManyChainMultiSig.GroupTreeNotWellFormed.selector));
        s_multisig.setConfig(s_signerAddresses, signerGroups, groupQuorums, groupParents, false);

        // now, let's fix the indizes in groupParents, everything should work just fine
        groupParents = [
            0,
            0,
            1,
            2,
            3,
            4,
            5,
            6,
            7,
            8,
            9,
            10,
            11,
            12,
            13,
            14,
            15,
            16,
            17,
            18,
            19,
            20,
            21,
            22,
            23,
            24,
            25,
            26,
            27,
            28,
            29,
            30
        ];
        s_multisig.setConfig(s_signerAddresses, signerGroups, groupQuorums, groupParents, false);
        ManyChainMultiSig.Config memory config = s_multisig.getConfig();
        assertEq(abi.encode(groupParents), abi.encode(config.groupParents));
    }
}
