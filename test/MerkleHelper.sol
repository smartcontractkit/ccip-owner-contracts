// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import "../src/ManyChainMultiSig.sol";

contract MerkleHelper is Test {
    uint256 constant ROOT_METADATA_LEAF_INDEX = 0;

    // TODO: fix copy-paste redundancy
    struct SetRootArgs {
        bytes32 root;
        uint32 validUntil;
        ManyChainMultiSig.RootMetadata metadata;
        bytes32[] metadataProof;
        ManyChainMultiSig.Signature[] signatures;
    }

    function build(
        uint256[] memory privateKeys,
        uint32 validUntil,
        ManyChainMultiSig.RootMetadata memory metadata,
        ManyChainMultiSig.Op[] memory ops
    ) public returns (SetRootArgs memory setRootArgs, bytes32[][] memory opProofs) {
        bytes32[] memory leaves = constructLeaves(ops, metadata);

        (
            bytes32 root,
            bytes32[] memory metadataProof,
            ManyChainMultiSig.Signature[] memory signatures
        ) = constructAnsSignRootAndProof(leaves, validUntil, privateKeys);

        setRootArgs = SetRootArgs({
            root: root,
            validUntil: validUntil,
            metadata: metadata,
            metadataProof: metadataProof,
            signatures: signatures
        });

        opProofs = new bytes32[][](ops.length);
        for (uint256 i = 0; i < ops.length; i++) {
            opProofs[i] = computeProofForLeaf(leaves, getLeafIndexOfOp(i));
        }
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
        result = new ManyChainMultiSig.Signature[](privateKeys.length);
        for (uint256 i = 0; i < privateKeys.length; i++) {
            (uint8 v, bytes32 r, bytes32 s) = vm.sign(privateKeys[i], hash);
            address signer = ecrecover(hash, v, r, s);
            assertTrue(signer == vm.addr(privateKeys[i]), "invalid signature");
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
