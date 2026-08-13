// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.13;

import "../src/RBACTimelock.sol";
import "./BaseTest.sol";

contract RBACTimelockRenounceRoleTest is BaseTest {
    function test_callProxyHasExecutorRole() public {
        // s_proxy is granted the executor role in BaseTest.setUp().
        assertTrue(s_timelock.hasRole(s_timelock.EXECUTOR_ROLE(), address(s_proxy)));
    }

    function test_callProxyCannotRenounceExecutorRole() public {
        bytes32 executorRole = s_timelock.EXECUTOR_ROLE();
        assertTrue(s_timelock.hasRole(executorRole, address(s_proxy)));

        // Forwarding a renounceRole call through the proxy is blocked.
        bytes memory data =
            abi.encodeCall(s_timelock.renounceRole, (executorRole, address(s_proxy)));

        vm.prank(EXTERNAL_CALLER);
        (bool success, bytes memory ret) = address(s_proxy).call(data);
        assertFalse(success);
        assertEq(
            ret,
            abi.encodeWithSignature("Error(string)", "CallProxy: renounceRole is blocked")
        );

        // The proxy still holds the executor role.
        assertTrue(s_timelock.hasRole(executorRole, address(s_proxy)));
    }

    function test_callProxyCanStillExecuteAfterBlockedRenounce() public {
        // Schedule an operation that the proxy can execute.
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

        // Attempt (and fail) to renounce the executor role through the proxy.
        vm.prank(EXTERNAL_CALLER);
        (bool renounceSuccess,) = address(s_proxy).call(
            abi.encodeCall(s_timelock.renounceRole, (s_timelock.EXECUTOR_ROLE(), address(s_proxy)))
        );
        assertFalse(renounceSuccess);

        // The proxy retains the role and can still execute.
        bytes memory executeData = abi.encodeCall(
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
        vm.prank(EXTERNAL_CALLER);
        (bool executeSuccess,) = address(s_proxy).call(executeData);
        assertTrue(executeSuccess);
        assertEq(s_counter.number(), 1);
    }
}
