// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.13;

import "./BaseTest.sol";

contract RBACTimelockReceivableTest is BaseTest {
    function test_canReceiveETH() public {
        vm.prank(ADMIN);
        payable(address(s_timelock)).transfer(0.5 ether);
        assertEq(address(s_timelock).balance, 0.5 ether);
    }
}
