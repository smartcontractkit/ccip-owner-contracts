// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.13;

import "./BaseTest.sol";
import "../src/RBACTimelock.sol";
import "../mock/Counter.sol";

contract RBACTimelockBlockFunctionTest is BaseTest {
    function test_failsIfNotAdmin() public {
        vm.expectRevert(_getExpectedMissingRoleErrorMessage(address(this), s_timelock.ADMIN_ROLE()));
        s_timelock.blockFunctionSelector(Counter.increment.selector);
    }

    function test_blocksFunction() public {
        // Schedule operation should succeed
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

        // Block function selector
        vm.prank(ADMIN);
        s_timelock.blockFunctionSelector(Counter.increment.selector);
        uint256 blockedFnSelectorCount = s_timelock.getBlockedFunctionSelectorCount();
        assertEq(blockedFnSelectorCount, 1);
        bytes4 blockedFnSelector = s_timelock.getBlockedFunctionSelectorAt(0);
        assertEq(blockedFnSelector, bytes4(Counter.increment.selector));

        // Make sure blocked function cannot be scheduled
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
            bytes32("salt"),
            MIN_DELAY
        );
    }

    function test_blocksFunction_c4issue41() public {
        bytes4 zeroSelector = 0x00000000;
        // Block function selector
        vm.prank(ADMIN);
        s_timelock.blockFunctionSelector(zeroSelector);
        uint256 blockedFnSelectorCount = s_timelock.getBlockedFunctionSelectorCount();
        assertEq(blockedFnSelectorCount, 1);
        bytes4 blockedFnSelector = s_timelock.getBlockedFunctionSelectorAt(0);
        assertEq(blockedFnSelector, bytes4(zeroSelector));

        // Make sure that zero selector cannot be scheduled
        vm.expectRevert("RBACTimelock: selector is blocked");
        vm.prank(PROPOSER_ONE);
        s_timelock.scheduleBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: bytes.concat(zeroSelector)
                })
            ),
            NO_PREDECESSOR,
            bytes32("salt"),
            MIN_DELAY
        );

        // Make sure that zero selector plus another zero cannot be scheduled
        vm.expectRevert("RBACTimelock: selector is blocked");
        vm.prank(PROPOSER_ONE);
        s_timelock.scheduleBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: bytes.concat(zeroSelector, bytes1(0))
                })
            ),
            NO_PREDECESSOR,
            bytes32("salt"),
            MIN_DELAY
        );

        // Make sure that empty call *can* be scheduled
        vm.prank(PROPOSER_ONE);
        s_timelock.scheduleBatch(
            _singletonCalls(
                RBACTimelock.Call({target: address(s_counter), value: 0, data: bytes.concat()})
            ),
            NO_PREDECESSOR,
            bytes32("salt"),
            MIN_DELAY
        );

        // Make sure that three zero bytes can be scheduled.
        vm.prank(PROPOSER_ONE);
        s_timelock.scheduleBatch(
            _singletonCalls(
                RBACTimelock.Call({
                    target: address(s_counter),
                    value: 0,
                    data: bytes.concat(bytes3(0x000000))
                })
            ),
            NO_PREDECESSOR,
            bytes32("salt"),
            MIN_DELAY
        );
    }

    function test_unblocksFunction() public {
        vm.prank(ADMIN);

        // Block Function
        s_timelock.blockFunctionSelector(Counter.increment.selector);

        // Try schedule blocked function and expect it to revert
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

        // Unblock Function
        vm.prank(ADMIN);
        s_timelock.unblockFunctionSelector(Counter.increment.selector);
        uint256 blockedFnSelectorCount = s_timelock.getBlockedFunctionSelectorCount();

        assertEq(blockedFnSelectorCount, 0);

        // Make sure unblocked function can be scheduled
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
        assertTrue(s_timelock.isOperation(operationID));
    }
}
