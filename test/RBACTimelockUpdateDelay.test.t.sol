// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.13;

import "./BaseTest.sol";

contract RBACTimelockExecuteTest is BaseTest {
    function test_cannotUpdateIfNotAdminRole() public {
        vm.prank(0x89205A3A3b2A69De6Dbf7f01ED13B2108B2c43e7);
        vm.expectRevert(
            "AccessControl: account 0x89205a3a3b2a69de6dbf7f01ed13b2108b2c43e7 is missing role 0xa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775"
        );
        s_timelock.updateDelay(3 days);
    }

    function test_updatesMinDelay() public {
        vm.prank(ADMIN);
        s_timelock.updateDelay(3 days);
        uint256 minDelay = s_timelock.getMinDelay();
        assertEq(minDelay, 3 days);
    }
}
