// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import "../src/ManyChainMultiSig.sol";
import "../mock/ExposedManyChainMultiSig.sol";
import {Receiver} from "../mock/Receiver.sol";

contract ManyChainMultiSigBaseTest is Test {
    // the test instance of ManyChainMultiSig. We expose the internal fields
    // of the contract via ExposedManyChainMultiSig so can be tested easily.
    ExposedManyChainMultiSig s_testExposedManyChainMultiSig;

    // the test config
    uint8 constant SIGNERS_NUM = 9;
    address[] s_testSigners;
    uint256[] s_testPrivateKeys;
    uint8 constant NUM_SUBGROUPS = 3; // SIGNERS_NUM/3 in each group
    uint8 constant MAX_NUM_GROUPS = 32;
    uint8 constant GROUP0_QUORUM = 2;
    uint8 constant GROUP1_QUORUM = 3;
    uint8 constant GROUP2_QUORUM = 2;
    uint8 constant GROUP3_QUORUM = 1;
    // all groups have the root as parent
    uint8 constant GROUP0_PARENT = 0;
    uint8 constant GROUP1_PARENT = 0;
    uint8 constant GROUP2_PARENT = 0;
    uint8 constant GROUP3_PARENT = 0;
    uint8[MAX_NUM_GROUPS] s_testGroupQuorums;
    uint8[MAX_NUM_GROUPS] s_testGroupParents;
    uint8[] s_signerGroups;
    ManyChainMultiSig.Config s_testConfig;

    // the test addresses
    address internal constant MULTISIG_OWNER = 0x1232117A6dEd3a4B844206D8f892e2733A71c843;
    address internal constant EXTERNAL_CALLER = 0x89205A3A3b2A69De6Dbf7f01ED13B2108B2c43e7;
    address internal s_multiSigContractAddress;

    function setUp() public virtual {
        // fill signers' addresses and their corresponding private keys
        (s_testSigners, s_testPrivateKeys) = addressesWithPrivateKeys(SIGNERS_NUM);

        // assign the required quorum in each group
        s_testGroupQuorums[0] = GROUP0_QUORUM;
        s_testGroupQuorums[1] = GROUP1_QUORUM;
        s_testGroupQuorums[2] = GROUP2_QUORUM;
        s_testGroupQuorums[3] = GROUP3_QUORUM;

        // assign signers to groups
        for (uint8 i = 1; i <= SIGNERS_NUM; i++) {
            // plus one because we don't want signers in root group
            s_signerGroups.push((i % NUM_SUBGROUPS) + 1);
        }

        for (uint8 i = 0; i < SIGNERS_NUM; i++) {
            s_testConfig.signers.push(
                ManyChainMultiSig.Signer({
                    addr: s_testSigners[i],
                    index: i,
                    group: s_signerGroups[i]
                })
            );
        }
        s_testConfig.groupQuorums = s_testGroupQuorums;
        s_testConfig.groupParents = s_testGroupParents;

        // sets the owner of the multiSig and construct a test ManyChainMultiSig instance
        vm.prank(MULTISIG_OWNER);
        s_testExposedManyChainMultiSig = new ExposedManyChainMultiSig();
        s_multiSigContractAddress = address(s_testExposedManyChainMultiSig);
    }

    function addressesWithPrivateKeys(uint8 numSigners)
        public
        returns (address[] memory addresses, uint256[] memory privateKeys)
    {
        string[] memory cmd = new string[](4);
        cmd[0] = "go";
        cmd[1] = "run";
        // must be executed from the parent package
        cmd[2] = "./testCommands/generateAddressesAndKeys.go";
        cmd[3] = vm.toString(numSigners);

        bytes memory result = vm.ffi(cmd);
        (privateKeys, addresses) = abi.decode(result, (uint256[], address[]));
    }
}

contract ManyChainMultiSigBaseSetRootAndExecuteTest is ManyChainMultiSigBaseTest {
    // the test ops, rootMetadata, merkle tree, and the corresponding
    // signatures.
    ManyChainMultiSig.Op[] s_testOps;
    ManyChainMultiSig.RootMetadata s_initialTestRootMetadata;
    uint32 constant s_testValidUntil = 1000000;
    bytes32 s_testInitialRoot;
    bytes32[] s_metadataProof;
    ManyChainMultiSig.Signature[] s_signatures;
    uint256 constant OPS_NUM = 7;
    uint256 constant REVERTING_OP_INDEX = 5;
    uint256 constant VALUE_OP_INDEX = 6;
    uint256 constant LEAVES_NUM = 8;
    uint256 constant ROOT_METADATA_LEAF_INDEX = 0;
    bytes32[] s_testLeavesInTree = new bytes32[](LEAVES_NUM);

    // the assumed chainId in tests
    uint256 constant s_testChainId = 2;

    function setUp() public virtual override {
        ManyChainMultiSigBaseTest.setUp();

        // construct the merkle tree which contains a set of ops and the
        // rootMetadata.

        s_initialTestRootMetadata = ManyChainMultiSig.RootMetadata({
            chainId: s_testChainId,
            multiSig: s_multiSigContractAddress,
            preOpCount: 0,
            postOpCount: uint40(OPS_NUM),
            overridePreviousRoot: false
        });

        constructAndStoreOps();

        s_testLeavesInTree = constructLeaves(s_testOps, s_initialTestRootMetadata);

        {
            ManyChainMultiSig.Signature[] memory signatures;
            (s_testInitialRoot, s_metadataProof, signatures) = constructAnsSignRootAndProof(
                s_testLeavesInTree, s_testValidUntil, s_testPrivateKeys
            );

            assert(s_signatures.length == 0);
            for (uint256 i = 0; i < signatures.length; i++) {
                s_signatures.push(signatures[i]);
            }
        }

        // assign chainId and timestamp permanently
        vm.chainId(s_testChainId);
        // we start from timestamp 0
        vm.warp(0);

        // set a new config
        vm.prank(MULTISIG_OWNER);
        s_testExposedManyChainMultiSig.setConfig(
            s_testSigners, s_signerGroups, s_testGroupQuorums, s_testGroupParents, false
        );
    }

    function computeRoot(bytes32[] memory leaves) internal pure returns (bytes32 root) {
        bytes32[] memory proof;
        bool[] memory proofFlags = new bool[](leaves.length - 1);
        for (uint256 i = 0; i < proofFlags.length; i++) {
            // indicate that the intermediate hash should be computed, i.e., the internal
            // node (since we don't provide it)
            proofFlags[i] = true;
        }
        // processMultiProof computes the root from the set of leaves and returns it
        return MerkleProof.processMultiProof(proof, proofFlags, leaves);
    }

    function XXXTestOnly_fillSignatures(
        bytes32 root,
        uint32 validUntil,
        uint256[] memory privateKeys
    ) internal returns (ManyChainMultiSig.Signature[] memory result) {
        bytes32 hash = ECDSA.toEthSignedMessageHash(keccak256(abi.encode(root, validUntil)));
        result = new ManyChainMultiSig.Signature[](SIGNERS_NUM);
        for (uint256 i = 0; i < SIGNERS_NUM; i++) {
            (uint8 v, bytes32 r, bytes32 s) = vm.sign(privateKeys[i], hash);
            address signer = ecrecover(hash, v, r, s);
            assertTrue(signer == s_testSigners[i], "invalid signature");
            result[i] = ManyChainMultiSig.Signature({v: v, r: r, s: s});
        }
        return result;
    }

    // equivalent to ceil(log2(leafCount))
    function proofLen(uint256 leafCount) internal pure returns (uint256) {
        uint256 power = 1;
        uint256 exp = 0;
        while (power < leafCount) {
            power *= 2;
            exp += 1;
        }
        return exp;
    }

    function computeProofForLeaf(bytes32[] memory data, uint256 index)
        internal
        pure
        returns (bytes32[] memory)
    {
        // this method assumes that there is an even number of leaves.
        assert(data.length % 2 == 0);
        bytes32[] memory proof = new bytes32[](proofLen(data.length));
        uint256 currentPos = 0;
        while (data.length > 1) {
            if (index & 0x1 == 1) {
                proof[currentPos] = data[index - 1];
            } else {
                proof[currentPos] = data[index + 1];
            }
            index = index / 2;
            data = hashLevel(data);
            currentPos++;
        }
        return proof;
    }

    function hashLevel(bytes32[] memory data) internal pure returns (bytes32[] memory) {
        uint256 currentPos = 0;
        uint256 length = data.length;
        bytes32[] memory newData = new bytes32[](length / 2);
        for (uint256 i = 0; i < length - 1; i += 2) {
            if (data[i] < data[i + 1]) {
                newData[currentPos] = keccak256(abi.encode(data[i], data[i + 1]));
            } else {
                newData[currentPos] = keccak256(abi.encode(data[i + 1], data[i]));
            }

            currentPos++;
        }
        return newData;
    }

    function constructAndStoreOps() internal {
        for (uint256 i = 0; i < OPS_NUM; i++) {
            // the last op reverts
            bool reverts = (i == REVERTING_OP_INDEX ? true : false);
            uint256 value = (i == VALUE_OP_INDEX ? 1 : 0);
            s_testOps.push(
                ManyChainMultiSig.Op({
                    chainId: s_testChainId,
                    multiSig: s_multiSigContractAddress,
                    nonce: uint40(i),
                    to: address(new Receiver()),
                    value: value,
                    data: abi.encodeWithSignature("executableMethod(bool)", reverts)
                })
            );
        }
    }

    function constructLeaves(
        ManyChainMultiSig.Op[] memory ops,
        ManyChainMultiSig.RootMetadata memory rootMetadata
    ) internal pure returns (bytes32[] memory leaves) {
        leaves = new bytes32[](ops.length + 1);
        leaves[ROOT_METADATA_LEAF_INDEX] = keccak256(leafMetadataPreimage(rootMetadata));
        for (uint256 i = 0; i < ops.length; i++) {
            uint256 leafIndex;
            bytes32 leaf = keccak256(leafOpPreimage(ops[i]));
            leafIndex = (i >= ROOT_METADATA_LEAF_INDEX ? i + 1 : i);
            leaves[leafIndex] = leaf;
        }
        return leaves;
    }

    function leafMetadataPreimage(ManyChainMultiSig.RootMetadata memory rootMetadata)
        internal
        pure
        returns (bytes memory)
    {
        return abi.encode(MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_METADATA, rootMetadata);
    }

    function leafOpPreimage(ManyChainMultiSig.Op memory op) internal pure returns (bytes memory) {
        return abi.encode(MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP, op);
    }

    function getLeafIndexOfOp(uint256 opIndex) internal pure returns (uint256) {
        return opIndex < ROOT_METADATA_LEAF_INDEX ? opIndex : opIndex + 1;
    }

    function constructAnsSignRootAndProof(
        bytes32[] memory leaves,
        uint32 validUntil,
        uint256[] memory privateKeys
    )
        internal
        returns (
            bytes32 root,
            bytes32[] memory metadataProof,
            ManyChainMultiSig.Signature[] memory signatures
        )
    {
        root = computeRoot(leaves);
        metadataProof = computeProofForLeaf(leaves, ROOT_METADATA_LEAF_INDEX);
        signatures = XXXTestOnly_fillSignatures(root, validUntil, privateKeys);
    }
}
