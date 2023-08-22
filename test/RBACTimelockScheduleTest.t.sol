// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.13;

import "../src/RBACTimelock.sol";
import "./BaseTest.sol";
import "../mock/Counter.sol";

contract RBACTimelockScheduleBatchTest is BaseTest {
    Counter internal s_counterTwo;
    RBACTimelock.Call[] internal s_calls;

    function setUp() public override {
        BaseTest.setUp();

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

    function test_cannotScheduleBatchIfNotProposer() public {
        vm.expectRevert(
            _getExpectedMissingRoleErrorMessage(address(this), s_timelock.PROPOSER_ROLE())
        );
        s_timelock.scheduleBatch(s_calls, NO_PREDECESSOR, EMPTY_SALT, MIN_DELAY);
    }

    function test_cannotScheduleIfBatchContainsBlockedFunction() public {
        // Block function
        vm.prank(ADMIN);
        s_timelock.blockFunctionSelector(Counter.increment.selector);

        // Expect revert
        vm.expectRevert("RBACTimelock: selector is blocked");
        vm.prank(PROPOSER_ONE);
        s_timelock.scheduleBatch(s_calls, NO_PREDECESSOR, EMPTY_SALT, MIN_DELAY);
    }

    function test_proposerCanBatchSchedule() public {
        _scheduleBatchedOperation(PROPOSER_ONE);
    }

    function test_adminCanBatchSchedule() public {
        _scheduleBatchedOperation(ADMIN);
    }

    event CallScheduled(
        bytes32 indexed id,
        uint256 indexed index,
        address target,
        uint256 value,
        bytes data,
        bytes32 predecessor,
        bytes32 salt,
        uint256 delay
    );

    function _scheduleBatchedOperation(address proposer) internal {
        bytes32 batchedOperationID =
            s_timelock.hashOperationBatch(s_calls, NO_PREDECESSOR, EMPTY_SALT);

        assertEq(s_timelock.isOperation(batchedOperationID), false);

        for (uint256 i = 0; i < s_calls.length; i++) {
            vm.expectEmit(true, true, true, true, address(s_timelock));

            emit CallScheduled(
                batchedOperationID,
                i,
                s_calls[i].target,
                s_calls[i].value,
                s_calls[i].data,
                NO_PREDECESSOR,
                EMPTY_SALT,
                MIN_DELAY
            );
        }

        vm.prank(proposer);
        s_timelock.scheduleBatch(s_calls, NO_PREDECESSOR, EMPTY_SALT, MIN_DELAY);

        assertEq(s_timelock.isOperation(batchedOperationID), true);
    }
}

contract RBACTimelockScheduleTest is BaseTest {
    function test_nonProposerCannotSchedule() public {
        vm.expectRevert(
            _getExpectedMissingRoleErrorMessage(address(this), s_timelock.PROPOSER_ROLE())
        );
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
    }

    function test_cannotScheduleABlockedFunction() public {
        vm.prank(ADMIN);
        s_timelock.blockFunctionSelector(Counter.increment.selector);
        vm.expectRevert("RBACTimelock: selector is blocked");
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
    }

    function test_cannotScheduleIfOperationScheduled() public {
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
        vm.expectRevert("RBACTimelock: operation already scheduled");
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
    }

    function test_cannotScheduleIfDelayLessThanMinDelay() public {
        vm.expectRevert("RBACTimelock: insufficient delay");
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
            3 days
        );
    }

    function test_proposerCanScheduleOperation() public {
        _scheduleOperation(PROPOSER_ONE);
    }

    function test_adminCanScheduleOperation() public {
        _scheduleOperation(ADMIN);
    }

    function _scheduleOperation(address proposer) internal {
        vm.prank(proposer);
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
        assertTrue(s_timelock.isOperation(operationID));
    }
}
