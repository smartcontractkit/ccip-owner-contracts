// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.13;

import "../src/RBACTimelock.sol";
import "./BaseTest.sol";
import "../mock/Counter.sol";

contract RBACTimelockBypasserExecuteBatchTest is BaseTest {
    Counter internal s_counterTwo;
    RBACTimelock.Call[] internal s_calls;

    // BypasserCallExecuted as defined in RBACTimelock.sol. Redefine it here for testing.
    event BypasserCallExecuted(uint256 indexed index, address target, uint256 value, bytes data);

    function setUp() public override {
        BaseTest.setUp();

        s_counterTwo = new Counter(address(s_timelock));

        s_calls.push(
            RBACTimelock.Call({
                target: address(s_counter),
                value: 0,
                data: abi.encodeWithSelector(Counter.increment.selector)
            })
        );
        s_calls.push(
            RBACTimelock.Call({
                target: address(s_counterTwo),
                value: 0,
                data: abi.encodeWithSelector(Counter.setNumber.selector, 10)
            })
        );
    }

    function test_cannotExecuteBatchIfNotBypasserOrAdmin() public {
        vm.expectRevert(
            _getExpectedMissingRoleErrorMessage(address(this), s_timelock.BYPASSER_ROLE())
        );
        s_timelock.bypasserExecuteBatch(s_calls);
    }

    function test_cannotExecuteIfOneTargetReverts() public {
        // Schedule a job where one target will revert
        s_calls[1].data = abi.encodeWithSelector(Counter.mockRevert.selector);

        vm.warp(block.timestamp + MIN_DELAY + 2 days);
        vm.prank(ADMIN);
        vm.expectRevert("RBACTimelock: underlying transaction reverted");
        s_timelock.bypasserExecuteBatch(s_calls);
    }

    function test_operationsAreBypasserBatchExecuted() public {
        // Batch execute two operations by bypasser
        // 1) Increment s_counter's number from 0 to 1
        // 2) Set s_counterTwo's number to 10
        vm.prank(BYPASSER_ONE);

        require(s_calls.length >= 1, "s_calls must not be empty for this test");
        for (uint256 i = 0; i < s_calls.length; ++i) {
            vm.expectEmit(true, true, true, true);
            emit BypasserCallExecuted(i, s_calls[i].target, s_calls[i].value, s_calls[i].data);
        }
        s_timelock.bypasserExecuteBatch(s_calls);
        assertEq(s_counter.number(), 1);
        assertEq(s_counterTwo.number(), 10);
    }
}

contract RBACTimelockExecuteBatchTest is BaseTest {
    Counter internal s_counterTwo;
    RBACTimelock.Call[] internal s_calls;

    function setUp() public override {
        BaseTest.setUp();
        s_counterTwo = new Counter(address(s_timelock));

        s_calls.push(
            RBACTimelock.Call({
                target: address(s_counter),
                value: 0,
                data: abi.encodeWithSelector(Counter.increment.selector)
            })
        );
        s_calls.push(
            RBACTimelock.Call({
                target: address(s_counterTwo),
                value: 0,
                data: abi.encodeWithSelector(Counter.setNumber.selector, 10)
            })
        );
    }

    function test_cannotExecuteBatchIfNotExecutor() public {
        vm.expectRevert(
            _getExpectedMissingRoleErrorMessage(address(this), s_timelock.EXECUTOR_ROLE())
        );
        s_timelock.executeBatch(s_calls, NO_PREDECESSOR, EMPTY_SALT);
    }

    function test_cannotBeExecutedIfOperationNotReady() public {
        vm.prank(PROPOSER_ONE);
        s_timelock.scheduleBatch(s_calls, NO_PREDECESSOR, EMPTY_SALT, MIN_DELAY);
        vm.warp(block.timestamp + MIN_DELAY - 2 days);
        vm.expectRevert("RBACTimelock: operation is not ready");
        vm.prank(EXECUTOR_ONE);
        s_timelock.executeBatch(s_calls, NO_PREDECESSOR, EMPTY_SALT);
    }

    function test_cannotBeExecutedIfPredecessorOperationNotExecuted() public {
        vm.prank(PROPOSER_ONE);

        // Schedule predecessor job
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
        bytes32 operationOneID = s_timelock.hashOperationBatch(
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

        // Schedule dependent job
        vm.prank(PROPOSER_ONE);
        s_timelock.scheduleBatch(s_calls, operationOneID, EMPTY_SALT, MIN_DELAY);

        // Check that executing the dependent job reverts
        vm.warp(block.timestamp + MIN_DELAY + 2 days);
        vm.expectRevert("RBACTimelock: missing dependency");
        vm.prank(EXECUTOR_ONE);
        s_timelock.executeBatch(s_calls, operationOneID, EMPTY_SALT);
    }

    function test_cannotExecuteIfOneTargetReverts() public {
        vm.prank(PROPOSER_ONE);

        // Schedule a job where one target will revert
        s_calls[1].data = abi.encodeWithSelector(Counter.mockRevert.selector);
        s_timelock.scheduleBatch(s_calls, NO_PREDECESSOR, EMPTY_SALT, MIN_DELAY);

        vm.warp(block.timestamp + MIN_DELAY + 2 days);
        vm.prank(EXECUTOR_ONE);
        vm.expectRevert("RBACTimelock: underlying transaction reverted");
        s_timelock.executeBatch(s_calls, NO_PREDECESSOR, EMPTY_SALT);
    }

    function test_executorCanBatchExecuteOperation() public {
        _executeBatchedOperation(EXECUTOR_ONE);
    }

    function test_adminCanBatchExecuteOperation() public {
        _executeBatchedOperation(ADMIN);
    }

    function _executeBatchedOperation(address executor) internal {
        vm.prank(PROPOSER_ONE);

        // Schedule batch executon
        s_timelock.scheduleBatch(s_calls, NO_PREDECESSOR, EMPTY_SALT, MIN_DELAY);

        vm.warp(block.timestamp + MIN_DELAY);

        vm.prank(executor);
        s_timelock.executeBatch(s_calls, NO_PREDECESSOR, EMPTY_SALT);

        bytes32 operationID = s_timelock.hashOperationBatch(s_calls, NO_PREDECESSOR, EMPTY_SALT);
        uint256 operationTimestamp = s_timelock.getTimestamp(operationID);
        assertEq(operationTimestamp, DONE_TIMESTAMP);
    }
}

contract RBACTimelockExecuteTest is BaseTest {
    // CallExecuted as defined in RBACTimelock.sol. Redefine it here for testing.
    event CallExecuted(
        bytes32 indexed id, uint256 indexed index, address target, uint256 value, bytes data
    );

    function test_cannotBeExecutedByNonExecutorIfRestrictionsSet() public {
        vm.expectRevert(
            _getExpectedMissingRoleErrorMessage(address(this), s_timelock.EXECUTOR_ROLE())
        );
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
    }

    function test_cannotBeExecutedIfOperationNotReady() public {
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
        vm.warp(block.timestamp + MIN_DELAY - 2 days);
        vm.expectRevert("RBACTimelock: operation is not ready");
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
    }

    function test_cannotBeExecutedIfPredecessorOperationNotExecuted() public {
        vm.prank(PROPOSER_ONE);

        // Schedule predecessor job
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
        bytes32 operationOneID = s_timelock.hashOperationBatch(
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

        // Schedule dependent job
        vm.prank(PROPOSER_ONE);
        s_timelock.scheduleBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.setNumber.selector, 1)
                })
            ),
            operationOneID,
            EMPTY_SALT,
            MIN_DELAY
        );

        // Check that executing the dependent job reverts
        vm.warp(block.timestamp + MIN_DELAY + 2 days);
        vm.expectRevert("RBACTimelock: missing dependency");
        vm.prank(EXECUTOR_ONE);
        s_timelock.executeBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.setNumber.selector, 1)
                })
            ),
            operationOneID,
            EMPTY_SALT
        );
    }

    function test_cannotExecuteIfTargetReverts() public {
        vm.prank(PROPOSER_ONE);

        // Schedule predecessor job
        s_timelock.scheduleBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.mockRevert.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT,
            MIN_DELAY
        );

        vm.warp(block.timestamp + MIN_DELAY + 2 days);
        vm.expectRevert("RBACTimelock: underlying transaction reverted");
        vm.prank(EXECUTOR_ONE);
        s_timelock.executeBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.mockRevert.selector)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT
        );
    }

    function test_executorCanExecuteOperation() public {
        _executeOperation(EXECUTOR_ONE);
    }

    function test_adminCanExecuteOperation() public {
        _executeOperation(ADMIN);
    }

    function _executeOperation(address executor) internal {
        vm.prank(PROPOSER_ONE);
        uint256 num = 10;
        s_timelock.scheduleBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.setNumber.selector, num)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT,
            MIN_DELAY
        );

        vm.warp(block.timestamp + MIN_DELAY + 2 days);
        vm.prank(executor);
        s_timelock.executeBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: abi.encodeWithSelector(Counter.setNumber.selector, num)
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
                    data: abi.encodeWithSelector(Counter.setNumber.selector, num)
                })
            ),
            NO_PREDECESSOR,
            EMPTY_SALT
        );
        uint256 operationTimestamp = s_timelock.getTimestamp(operationID);
        assertEq(operationTimestamp, DONE_TIMESTAMP);
        uint256 counterNumber = s_counter.number();
        assertEq(counterNumber, num);
    }

    function test_executeThroughCallProxy() public {
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
        vm.warp(block.timestamp + MIN_DELAY + 1);
        RBACTimelock.Call[] memory calls = _singletonCalls(
            RBACTimelock.Call({
                target: address(s_counter),
                value: 0,
                data: abi.encodeCall(Counter.increment, ())
            })
        );
        bytes32 predecessor = NO_PREDECESSOR;
        bytes32 salt = EMPTY_SALT;

        bytes memory data = abi.encodeCall(s_timelock.executeBatch, (calls, predecessor, salt));
        bytes32 id = s_timelock.hashOperationBatch(calls, predecessor, salt);

        // call execute from non-executor address
        vm.prank(EXTERNAL_CALLER);

        vm.expectEmit(true, true, true, true);
        emit CallExecuted(id, 0, calls[0].target, calls[0].value, calls[0].data);
        // s_proxy is set as executor in setUp()
        (bool success,) = address(s_proxy).call(data);
        assertTrue(success);
    }

    function test_executeThroughInvalidCallProxy() public {
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
        vm.warp(block.timestamp + MIN_DELAY + 1);
        bytes memory data = abi.encodeCall(
            s_timelock.executeBatch,
            (
                _singletonCalls(
                    RBACTimelock.Call({
                        target: address(s_counter),
                        value: 0,
                        data: abi.encodeCall(Counter.increment, ())
                    })
                    ),
                NO_PREDECESSOR,
                EMPTY_SALT
            )
        );
        // call execute from non-executor address
        vm.prank(EXTERNAL_CALLER);
        // now try to execute from another CallProxy (not the one set as executor)
        CallProxy faultyCallProxy = new CallProxy(address(s_timelock));
        vm.expectRevert();
        (bool success,) = address(faultyCallProxy).call(data);
        assertTrue(success);
    }
}
