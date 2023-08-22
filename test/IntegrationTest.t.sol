pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import "../mock/Counter.sol";
import "../src/CallProxy.sol";
import "../src/ManyChainMultiSig.sol";
import "../src/RBACTimelock.sol";
import "./ManyChainMultiSigBaseTest.t.sol";
import "./MerkleHelper.sol";
import "openzeppelin-contracts/access/Ownable2Step.sol";

interface IExecuteBatch {
    function executeBatch(RBACTimelock.Call[] calldata calls, bytes32 predecessor, bytes32 salt)
        external;
}

// Integration test that closely resembles planned deployment for ManyChainMultiSig + RBACTimelock
contract IntegrationTest is Test {
    uint8 constant MCMS_NUM_GROUPS = 32;

    uint8 constant PROPOSE_COUNT = 8;
    uint8 constant PROPOSE_QUORUM = 4;

    uint8 constant VETO_COUNT = 22 + 7;
    uint8 constant VETO_QUORUM = ((VETO_COUNT - 1) / 3) + 1;

    uint256 constant MIN_DELAY = 24 hours;

    MerkleHelper s_merkleHelper;

    ManyChainMultiSigBaseTest s_manyChainMultiSigBaseTest;

    address[] s_signerAddresses;
    uint256[] s_signerPrivateKeys;

    RBACTimelock s_timelock;

    ManyChainMultiSig s_proposeMultiSig;
    ManyChainMultiSig s_vetoMultiSig;
    ManyChainMultiSig s_bypassMultiSig;

    CallProxy s_callProxy;

    Counter s_counter;

    function setUp() public virtual {
        s_merkleHelper = new MerkleHelper();
        s_manyChainMultiSigBaseTest = new ManyChainMultiSigBaseTest();
        (s_signerAddresses, s_signerPrivateKeys) =
            s_manyChainMultiSigBaseTest.addressesWithPrivateKeys(PROPOSE_COUNT + VETO_COUNT);

        address timelockAddress =
            computeCreateAddress(address(this), vm.getNonce(address(this)) + 0);
        address proposeMultiSigAddress =
            computeCreateAddress(address(this), vm.getNonce(address(this)) + 1);
        address vetoMultiSigAddress =
            computeCreateAddress(address(this), vm.getNonce(address(this)) + 2);
        address bypassMultiSigAddress =
            computeCreateAddress(address(this), vm.getNonce(address(this)) + 3);
        address callProxyAddress =
            computeCreateAddress(address(this), vm.getNonce(address(this)) + 4);

        s_timelock = new RBACTimelock(
            MIN_DELAY, 
            timelockAddress, 
            oneAddress(proposeMultiSigAddress),
            oneAddress(callProxyAddress),
            oneAddress(vetoMultiSigAddress), 
            oneAddress(bypassMultiSigAddress)
        );

        {
            s_proposeMultiSig = new ManyChainMultiSig();
            uint8[] memory proposeGroups = new uint8[](PROPOSE_COUNT);
            uint8[MCMS_NUM_GROUPS] memory proposeGroupQuorums;
            proposeGroupQuorums[0] = PROPOSE_QUORUM;
            uint8[MCMS_NUM_GROUPS] memory proposeGroupParents;
            s_proposeMultiSig.setConfig(
                proposeAddresses(), proposeGroups, proposeGroupQuorums, proposeGroupParents, false
            );
            transferOwner(s_proposeMultiSig, timelockAddress);
        }
        {
            s_vetoMultiSig = new ManyChainMultiSig();
            uint8[] memory vetoGroups = new uint8[](VETO_COUNT);
            uint8[MCMS_NUM_GROUPS] memory vetoGroupQuorums;
            vetoGroupQuorums[0] = VETO_QUORUM;
            uint8[MCMS_NUM_GROUPS] memory vetoGroupParents;
            s_vetoMultiSig.setConfig(
                vetoAddresses(), vetoGroups, vetoGroupQuorums, vetoGroupParents, false
            );
            transferOwner(s_vetoMultiSig, timelockAddress);
        }
        {
            s_bypassMultiSig = new ManyChainMultiSig();
            uint8[] memory bypassGroups = new uint8[](PROPOSE_COUNT + VETO_COUNT);
            for (uint256 i = 0; i < PROPOSE_COUNT; i++) {
                bypassGroups[i] = 1;
            }
            for (uint256 i = PROPOSE_COUNT; i < PROPOSE_COUNT + VETO_COUNT; i++) {
                bypassGroups[i] = 2;
            }
            uint8[MCMS_NUM_GROUPS] memory bypassGroupQuorums;
            bypassGroupQuorums[0] = 2;
            (bypassGroupQuorums[1], bypassGroupQuorums[2]) = (PROPOSE_QUORUM, VETO_QUORUM);
            uint8[MCMS_NUM_GROUPS] memory bypassGroupParents;
            bypassGroupParents[0] = 0;
            (bypassGroupParents[1], bypassGroupParents[2]) = (0, 0);

            assertEq(s_signerAddresses.length, bypassGroups.length);
            s_bypassMultiSig.setConfig(
                s_signerAddresses, bypassGroups, bypassGroupQuorums, bypassGroupParents, false
            );
            transferOwner(s_bypassMultiSig, timelockAddress);
        }
        {
            s_callProxy = new CallProxy(address(s_timelock));
        }
        {
            s_counter = new Counter(timelockAddress);
            assertEq(s_counter.number(), 0);
        }

        assertEq(timelockAddress, address(s_timelock));
        assertEq(proposeMultiSigAddress, address(s_proposeMultiSig));
        assertEq(vetoMultiSigAddress, address(s_vetoMultiSig));
        assertEq(bypassMultiSigAddress, address(s_bypassMultiSig));
        assertEq(callProxyAddress, address(s_callProxy));
    }

    function oneAddress(address a) internal pure returns (address[] memory) {
        address[] memory result = new address[](1);
        result[0] = a;
        return result;
    }

    function proposeAddresses() internal view returns (address[] memory) {
        address[] memory result = new address[](PROPOSE_COUNT);
        for (uint256 i = 0; i < PROPOSE_COUNT; i++) {
            result[i] = s_signerAddresses[i];
        }
        return result;
    }

    function proposePrivateKeys() internal view returns (uint256[] memory) {
        uint256[] memory result = new uint256[](PROPOSE_COUNT);
        for (uint256 i = 0; i < PROPOSE_COUNT; i++) {
            result[i] = s_signerPrivateKeys[i];
        }
        return result;
    }

    function vetoAddresses() internal view returns (address[] memory) {
        address[] memory result = new address[](VETO_COUNT);
        for (uint256 i = 0; i < VETO_COUNT; i++) {
            result[i] = s_signerAddresses[PROPOSE_COUNT + i];
        }
        return result;
    }

    function vetoPrivateKeys() internal view returns (uint256[] memory) {
        uint256[] memory result = new uint256[](VETO_COUNT);
        for (uint256 i = 0; i < VETO_COUNT; i++) {
            result[i] = s_signerPrivateKeys[PROPOSE_COUNT + i];
        }
        return result;
    }

    function transferOwner(Ownable2Step o, address newOwner) public {
        o.transferOwnership(newOwner);
        vm.prank(newOwner);
        o.acceptOwnership();
        assertEq(o.owner(), newOwner);
    }

    event Cancelled(bytes32 indexed id);

    function test_chainOfActions() public {
        RBACTimelock.Call[] memory calls;
        bytes32 callsHash;
        bytes32 proposePredecessor;

        //
        // increment twice through regular flow
        //
        {
            calls = new RBACTimelock.Call[](2);
            calls[0] = RBACTimelock.Call({
                target: address(s_counter),
                value: 0,
                data: abi.encodeWithSelector(Counter.increment.selector)
            });
            calls[1] = RBACTimelock.Call({
                target: address(s_counter),
                value: 0,
                data: abi.encodeWithSelector(Counter.increment.selector)
            });
            callsHash = s_timelock.hashOperationBatch(calls, proposePredecessor, bytes32(0));
        }

        {
            ManyChainMultiSig.Op[] memory ops = new ManyChainMultiSig.Op[](1);
            ops[0] = ManyChainMultiSig.Op({
                chainId: block.chainid,
                multiSig: address(s_proposeMultiSig),
                nonce: 0,
                to: address(s_timelock),
                value: 0,
                data: abi.encodeWithSelector(
                    RBACTimelock.scheduleBatch.selector,
                    calls,
                    proposePredecessor,
                    bytes32(0),
                    MIN_DELAY
                    )
            });

            (MerkleHelper.SetRootArgs memory setRootArgs, bytes32[][] memory opProofs) =
            s_merkleHelper.build(
                proposePrivateKeys(),
                uint32(block.timestamp + 2 hours),
                ManyChainMultiSig.RootMetadata({
                    chainId: block.chainid,
                    multiSig: address(s_proposeMultiSig),
                    preOpCount: 0,
                    postOpCount: 1,
                    overridePreviousRoot: false
                }),
                ops
            );

            s_proposeMultiSig.setRoot(
                setRootArgs.root,
                setRootArgs.validUntil,
                setRootArgs.metadata,
                setRootArgs.metadataProof,
                setRootArgs.signatures
            );
            s_proposeMultiSig.execute(ops[0], opProofs[0]);

            // fails if minDelay hasn't elapsed
            vm.expectRevert("RBACTimelock: operation is not ready");
            IExecuteBatch(address(s_callProxy)).executeBatch(calls, bytes32(0), bytes32(0));

            vm.warp(block.timestamp + MIN_DELAY);

            IExecuteBatch(address(s_callProxy)).executeBatch(calls, bytes32(0), bytes32(0));
            assertEq(s_counter.number(), 2);
        }

        proposePredecessor = callsHash;

        //
        // again, increment twice through regular flow
        //
        {
            ManyChainMultiSig.Op[] memory ops = new ManyChainMultiSig.Op[](1);
            ops[0] = ManyChainMultiSig.Op({
                chainId: block.chainid,
                multiSig: address(s_proposeMultiSig),
                nonce: 1,
                to: address(s_timelock),
                value: 0,
                data: abi.encodeWithSelector(
                    RBACTimelock.scheduleBatch.selector,
                    calls,
                    proposePredecessor,
                    bytes32(0),
                    MIN_DELAY
                    )
            });

            (MerkleHelper.SetRootArgs memory setRootArgs, bytes32[][] memory opProofs) =
            s_merkleHelper.build(
                proposePrivateKeys(),
                uint32(block.timestamp + 2 hours),
                ManyChainMultiSig.RootMetadata({
                    chainId: block.chainid,
                    multiSig: address(s_proposeMultiSig),
                    preOpCount: 1,
                    postOpCount: 2,
                    overridePreviousRoot: false
                }),
                ops
            );

            s_proposeMultiSig.setRoot(
                setRootArgs.root,
                setRootArgs.validUntil,
                setRootArgs.metadata,
                setRootArgs.metadataProof,
                setRootArgs.signatures
            );
            s_proposeMultiSig.execute(ops[0], opProofs[0]);

            vm.warp(block.timestamp + MIN_DELAY);

            // fails if predecessor isn't right
            bytes32 wrongPredecessor = bytes32(uint256(proposePredecessor) + 1);
            vm.expectRevert("RBACTimelock: operation is not ready");
            IExecuteBatch(address(s_callProxy)).executeBatch(calls, wrongPredecessor, bytes32(0));

            // succeeds once we use right predecessor
            IExecuteBatch(address(s_callProxy)).executeBatch(calls, proposePredecessor, bytes32(0));
            assertEq(s_counter.number(), 4);
        }

        proposePredecessor = callsHash;

        //
        // halve minDelay from bypasser
        //
        {
            calls = new RBACTimelock.Call[](1);
            calls[0] = RBACTimelock.Call({
                target: address(s_timelock),
                value: 0,
                data: abi.encodeWithSelector(RBACTimelock.updateDelay.selector, MIN_DELAY / 2)
            });
            callsHash = s_timelock.hashOperationBatch(calls, proposePredecessor, bytes32(0));

            ManyChainMultiSig.Op[] memory ops = new ManyChainMultiSig.Op[](1);
            ops[0] = ManyChainMultiSig.Op({
                chainId: block.chainid,
                multiSig: address(s_bypassMultiSig),
                nonce: 0,
                to: address(s_timelock),
                value: 0,
                data: abi.encodeWithSelector(RBACTimelock.bypasserExecuteBatch.selector, calls)
            });

            (MerkleHelper.SetRootArgs memory setRootArgs, bytes32[][] memory opProofs) =
            s_merkleHelper.build(
                s_signerPrivateKeys,
                uint32(block.timestamp + 2 hours),
                ManyChainMultiSig.RootMetadata({
                    chainId: block.chainid,
                    multiSig: address(s_bypassMultiSig),
                    preOpCount: 0,
                    postOpCount: 1,
                    overridePreviousRoot: false
                }),
                ops
            );

            s_bypassMultiSig.setRoot(
                setRootArgs.root,
                setRootArgs.validUntil,
                setRootArgs.metadata,
                setRootArgs.metadataProof,
                setRootArgs.signatures
            );
            s_bypassMultiSig.execute(ops[0], opProofs[0]);

            assertEq(s_timelock.getMinDelay(), MIN_DELAY / 2);
        }

        //
        // propose a malicious timelock owner, who is then vetoed
        //
        {
            address evil = 0x2991C067BB8E27b078Ed2c086Af0a13c81013B93;

            calls = new RBACTimelock.Call[](1);
            calls[0] = RBACTimelock.Call({
                target: address(s_timelock),
                value: 0,
                data: abi.encodeWithSelector(
                    s_timelock.grantRole.selector, s_timelock.ADMIN_ROLE(), evil
                    )
            });
            callsHash = s_timelock.hashOperationBatch(calls, proposePredecessor, bytes32(0));

            ManyChainMultiSig.Op[] memory ops = new ManyChainMultiSig.Op[](1);
            ops[0] = ManyChainMultiSig.Op({
                chainId: block.chainid,
                multiSig: address(s_proposeMultiSig),
                nonce: 2,
                to: address(s_timelock),
                value: 0,
                data: abi.encodeWithSelector(
                    RBACTimelock.scheduleBatch.selector,
                    calls,
                    proposePredecessor,
                    bytes32(0),
                    MIN_DELAY
                    )
            });

            (MerkleHelper.SetRootArgs memory setRootArgs, bytes32[][] memory opProofs) =
            s_merkleHelper.build(
                proposePrivateKeys(),
                uint32(block.timestamp + 2 hours),
                ManyChainMultiSig.RootMetadata({
                    chainId: block.chainid,
                    multiSig: address(s_proposeMultiSig),
                    preOpCount: 2,
                    postOpCount: 3,
                    overridePreviousRoot: false
                }),
                ops
            );

            s_proposeMultiSig.setRoot(
                setRootArgs.root,
                setRootArgs.validUntil,
                setRootArgs.metadata,
                setRootArgs.metadataProof,
                setRootArgs.signatures
            );
            s_proposeMultiSig.execute(ops[0], opProofs[0]);

            vm.expectRevert("RBACTimelock: operation is not ready");
            IExecuteBatch(address(s_callProxy)).executeBatch(calls, proposePredecessor, bytes32(0));

            vm.warp(block.timestamp + MIN_DELAY / 4);

            // veto bad proposal!
            ops = new ManyChainMultiSig.Op[](1);
            ops[0] = ManyChainMultiSig.Op({
                chainId: block.chainid,
                multiSig: address(s_vetoMultiSig),
                nonce: 0,
                to: address(s_timelock),
                value: 0,
                data: abi.encodeWithSelector(RBACTimelock.cancel.selector, callsHash)
            });

            (setRootArgs, opProofs) = s_merkleHelper.build(
                vetoPrivateKeys(),
                uint32(block.timestamp + 2 hours),
                ManyChainMultiSig.RootMetadata({
                    chainId: block.chainid,
                    multiSig: address(s_vetoMultiSig),
                    preOpCount: 0,
                    postOpCount: 1,
                    overridePreviousRoot: false
                }),
                ops
            );

            s_vetoMultiSig.setRoot(
                setRootArgs.root,
                setRootArgs.validUntil,
                setRootArgs.metadata,
                setRootArgs.metadataProof,
                setRootArgs.signatures
            );
            vm.expectEmit(true, true, true, true);
            emit Cancelled(callsHash);
            s_vetoMultiSig.execute(ops[0], opProofs[0]);

            vm.warp(block.timestamp + MIN_DELAY);

            vm.expectRevert("RBACTimelock: operation is not ready");
            IExecuteBatch(address(s_callProxy)).executeBatch(calls, proposePredecessor, bytes32(0));
        }

        //
        // decrease quorum for vetoers & proposers
        //
        {
            uint8[] memory proposeGroups = new uint8[](PROPOSE_COUNT);
            uint8[MCMS_NUM_GROUPS] memory proposeGroupQuorums;
            proposeGroupQuorums[0] = PROPOSE_QUORUM - 1;
            uint8[MCMS_NUM_GROUPS] memory proposeGroupParents;

            uint8[] memory vetoGroups = new uint8[](VETO_COUNT);
            uint8[MCMS_NUM_GROUPS] memory vetoGroupQuorums;
            vetoGroupQuorums[0] = VETO_QUORUM - 1;
            uint8[MCMS_NUM_GROUPS] memory vetoGroupParents;

            calls = new RBACTimelock.Call[](2);
            calls[0] = RBACTimelock.Call({
                target: address(s_proposeMultiSig),
                value: 0,
                data: abi.encodeWithSelector(
                    ManyChainMultiSig.setConfig.selector,
                    proposeAddresses(),
                    proposeGroups,
                    proposeGroupQuorums,
                    proposeGroupParents,
                    false
                    )
            });
            calls[1] = RBACTimelock.Call({
                target: address(s_vetoMultiSig),
                value: 0,
                data: abi.encodeWithSelector(
                    ManyChainMultiSig.setConfig.selector,
                    vetoAddresses(),
                    vetoGroups,
                    vetoGroupQuorums,
                    vetoGroupParents,
                    false
                    )
            });
            callsHash = s_timelock.hashOperationBatch(calls, proposePredecessor, bytes32(0));

            ManyChainMultiSig.Op[] memory ops = new ManyChainMultiSig.Op[](1);
            ops[0] = ManyChainMultiSig.Op({
                chainId: block.chainid,
                multiSig: address(s_proposeMultiSig),
                nonce: 3,
                to: address(s_timelock),
                value: 0,
                data: abi.encodeWithSelector(
                    RBACTimelock.scheduleBatch.selector,
                    calls,
                    proposePredecessor,
                    bytes32(0),
                    MIN_DELAY
                    )
            });

            (MerkleHelper.SetRootArgs memory setRootArgs, bytes32[][] memory opProofs) =
            s_merkleHelper.build(
                proposePrivateKeys(),
                uint32(block.timestamp + 2 hours),
                ManyChainMultiSig.RootMetadata({
                    chainId: block.chainid,
                    multiSig: address(s_proposeMultiSig),
                    preOpCount: 3,
                    postOpCount: 4,
                    overridePreviousRoot: false
                }),
                ops
            );

            s_proposeMultiSig.setRoot(
                setRootArgs.root,
                setRootArgs.validUntil,
                setRootArgs.metadata,
                setRootArgs.metadataProof,
                setRootArgs.signatures
            );
            s_proposeMultiSig.execute(ops[0], opProofs[0]);

            vm.warp(block.timestamp + MIN_DELAY);

            IExecuteBatch(address(s_callProxy)).executeBatch(calls, proposePredecessor, bytes32(0));

            assertEq(s_proposeMultiSig.getConfig().groupQuorums[0], PROPOSE_QUORUM - 1);
            assertEq(s_vetoMultiSig.getConfig().groupQuorums[0], VETO_QUORUM - 1);
        }

        proposePredecessor = callsHash;
    }
}
