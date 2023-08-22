// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import "forge-std/Test.sol";
import {CallProxy} from "../src/CallProxy.sol";

contract CallProxyTest is Test {
    event TargetSet(address target);

    address constant MOCK_TARGET_ADDRESS = 0x1337133713371337133713371337133713371337;

    CallProxy s_callProxy;

    function setUp() public virtual {
        s_callProxy = new CallProxy(MOCK_TARGET_ADDRESS);
    }

    function testConstructor() public {
        vm.expectEmit(true, true, true, true);
        emit TargetSet(MOCK_TARGET_ADDRESS);
        new CallProxy(MOCK_TARGET_ADDRESS);
    }

    function testCall_fuzz(bool expectedSuccess, bytes memory call, bytes memory ret) public {
        if (expectedSuccess) {
            vm.mockCall(MOCK_TARGET_ADDRESS, 0, call, ret);
        } else {
            vm.mockCallRevert(MOCK_TARGET_ADDRESS, 0, call, ret);
        }
        (bool actualSuccess, bytes memory result) = address(s_callProxy).call(call);
        vm.clearMockedCalls();

        assertEq(result, ret);
        assertEq(expectedSuccess, actualSuccess);
    }
}
