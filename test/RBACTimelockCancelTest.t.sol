// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.13;

import "../src/RBACTimelock.sol";
import "./BaseTest.sol";
import "../mock/Counter.sol";

contract RBACTimelockCancelTest is BaseTest {
    function test_nonCancellerCannotCancel() public {
        vm.expectRevert(
            _getExpectedMissingRoleErrorMessage(EXECUTOR_ONE, s_timelock.CANCELLER_ROLE())
        );
        vm.prank(EXECUTOR_ONE);
        s_timelock.cancel(EMPTY_SALT);
    }

    function test_cannotCancelFinishedOperation() public {
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
        vm.prank(CANCELLER_ONE);
        vm.expectRevert("RBACTimelock: operation cannot be cancelled");
        s_timelock.cancel(operationID);
    }

    function test_cancellerCanCancelOperation() public {
        _cancelOperation(CANCELLER_ONE);
    }

    function test_adminCanCancelOperation() public {
        _cancelOperation(ADMIN);
    }

    function _cancelOperation(address canceller) internal {
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
        vm.prank(canceller);
        s_timelock.cancel(operationID);
        assertFalse(s_timelock.isOperation(operationID));
    }
}
