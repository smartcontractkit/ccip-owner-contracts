// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import "./BaseTest.sol";
import "../src/RBACTimelock.sol";
import "../mock/Counter.sol";

contract RBACTimelockIsOperationTest is BaseTest {
    function test_falseIfNotAnOperation() public {
        bool isOperation = s_timelock.isOperation(bytes32("non-op"));
        assertEq(isOperation, false);
    }

    function test_trueIfAnOperation() public {
        vm.prank(PROPOSER_ONE);
        s_timelock.scheduleBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(0),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT,
            MIN_DELAY
        );

        bytes32 operationID = s_timelock.hashOperationBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(0),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT
        );

        bool isOperation = s_timelock.isOperation(operationID);
        assertEq(isOperation, true);
    }
}

contract RBACTimelockIsOperationPendingTest is BaseTest {
    function test_falseIfNotAnOperation() public {
        bool isOperationPending = s_timelock.isOperationPending(bytes32("non-op"));
        assertEq(isOperationPending, false);
    }

    function test_trueIfScheduledOperatonNotYetExecuted() public {
        vm.prank(PROPOSER_ONE);
        s_timelock.scheduleBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT,
            MIN_DELAY
        );

        bytes32 operationID = s_timelock.hashOperationBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT
        );

        bool isOperationPending = s_timelock.isOperationPending(operationID);
        assertEq(isOperationPending, true);
    }

    function test_falseIfOperationHasBeenExecuted() public {
        vm.prank(PROPOSER_ONE);
        s_timelock.scheduleBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT,
            MIN_DELAY
        );

        vm.warp(block.timestamp + MIN_DELAY);
        vm.prank(EXECUTOR_ONE);
        s_timelock.executeBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT
        );

        bytes32 operationID = s_timelock.hashOperationBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT
        );

        bool isOperationPending = s_timelock.isOperationPending(operationID);
        assertEq(isOperationPending, false);
    }
}

contract RBACTimelockIsOperationReadyTest is BaseTest {
    function test_falseIfNotAnOperation() public {
        bool isOperationReady = s_timelock.isOperationReady(bytes32("non-op"));
        assertEq(isOperationReady, false);
    }

    function test_trueIfOnTheDelayedExecutionTime() public {
        vm.prank(PROPOSER_ONE);
        s_timelock.scheduleBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT,
            MIN_DELAY
        );

        vm.warp(block.timestamp + MIN_DELAY);

        bytes32 operationID = s_timelock.hashOperationBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT
        );

        bool isOperationReady = s_timelock.isOperationReady(operationID);
        assertEq(isOperationReady, true);
    }

    function test_trueIfAfterTheDelayedExecutionTime() public {
        vm.prank(PROPOSER_ONE);
        s_timelock.scheduleBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT,
            MIN_DELAY
        );

        vm.warp(block.timestamp + MIN_DELAY + 1 days);

        bytes32 operationID = s_timelock.hashOperationBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT
        );

        bool isOperationReady = s_timelock.isOperationReady(operationID);
        assertEq(isOperationReady, true);
    }

    function test_falseIfBeforeTheDelayedExecutionTime() public {
        vm.prank(PROPOSER_ONE);
        s_timelock.scheduleBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT,
            MIN_DELAY
        );

        vm.warp(block.timestamp + MIN_DELAY - 1 days);

        bytes32 operationID = s_timelock.hashOperationBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT
        );

        bool isOperationReady = s_timelock.isOperationReady(operationID);
        assertEq(isOperationReady, false);
    }

    function test_falseIfOperationHasBeenExecuted() public {
        vm.prank(PROPOSER_ONE);
        s_timelock.scheduleBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT,
            MIN_DELAY
        );

        vm.warp(block.timestamp + MIN_DELAY);
        vm.prank(EXECUTOR_ONE);
        s_timelock.executeBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT
        );

        bytes32 operationID = s_timelock.hashOperationBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT
        );

        bool isOperationReady = s_timelock.isOperationReady(operationID);
        assertEq(isOperationReady, false);
    }
}

contract RBACTimelockIsOperationDoneTest is BaseTest {
    function test_falseIfNotAnOperation() public {
        bool isOperationDone = s_timelock.isOperationDone(bytes32("non-op"));
        assertEq(isOperationDone, false);
    }

    function test_falseItTheOperationHasNotBeenExecuted() public {
        vm.prank(PROPOSER_ONE);
        s_timelock.scheduleBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT,
            MIN_DELAY
        );

        bytes32 operationID = s_timelock.hashOperationBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT
        );

        bool isOperationDone = s_timelock.isOperationDone(operationID);
        assertEq(isOperationDone, false);
    }

    function test_trueIfOperationHasBeenExecuted() public {
        vm.prank(PROPOSER_ONE);
        s_timelock.scheduleBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT,
            MIN_DELAY
        );

        vm.warp(block.timestamp + MIN_DELAY);
        vm.prank(EXECUTOR_ONE);
        s_timelock.executeBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT
        );

        bytes32 operationID = s_timelock.hashOperationBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT
        );

        bool isOperationDone = s_timelock.isOperationDone(operationID);
        assertEq(isOperationDone, true);
    }
}

contract RBACTimelockGetTimestampTest is BaseTest {
    function test_returnsZeroIfNotAnOperation() public {
        uint256 operationTimestamp = s_timelock.getTimestamp(bytes32("non-op"));
        assertEq(operationTimestamp, 0);
    }

    function test_returnsTheCorrectTimestampIfTheOperationHasNotBeenExecuted() public {
        vm.prank(PROPOSER_ONE);
        s_timelock.scheduleBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT,
            MIN_DELAY
        );

        bytes32 operationID = s_timelock.hashOperationBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT
        );

        uint256 operationTimestamp = s_timelock.getTimestamp(operationID);
        assertEq(operationTimestamp, block.timestamp + MIN_DELAY);
    }

    function test_returnsOneIfOperationHasBeenExecuted() public {
        vm.prank(PROPOSER_ONE);
        s_timelock.scheduleBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT,
            MIN_DELAY
        );

        vm.warp(block.timestamp + MIN_DELAY);
        vm.prank(EXECUTOR_ONE);
        s_timelock.executeBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT
        );

        bytes32 operationID = s_timelock.hashOperationBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.increment.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT
        );

        uint256 operationTimestamp = s_timelock.getTimestamp(operationID);
        assertEq(operationTimestamp, DONE_TIMESTAMP);
    }
}
