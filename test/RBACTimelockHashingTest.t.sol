// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.13;

import "../src/RBACTimelock.sol";
import "./BaseTest.sol";

contract RBACTimelockHashingTest is BaseTest {
    function test_hashesBatchedOperationsCorrectly() public {
        RBACTimelock.Call[] memory calls = new RBACTimelock.Call[](2);
        calls[0] = RBACTimelock.Call({
            target: address(s_counter),
            value: 0,
            data: abi.encodeWithSelector(Counter.increment.selector)
        });
        calls[1] = RBACTimelock.Call({
            target: address(s_counter),
            value: 1,
            data: abi.encodeWithSelector(Counter.setNumber.selector, 10)
        });
        bytes32 predecessor = NO_PREDECESSOR;
        bytes32 salt = EMPTY_SALT;

        bytes32 hashedOperation = s_timelock.hashOperationBatch(calls, predecessor, salt);
        bytes32 expectedHash = keccak256(abi.encode(calls, predecessor, salt));
        assertEq(hashedOperation, expectedHash);
    }
}
