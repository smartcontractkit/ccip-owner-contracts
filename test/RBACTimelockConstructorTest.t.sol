// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.13;

import "../src/RBACTimelock.sol";
import "./BaseTest.sol";

contract RBACTimelockConstructorTest is BaseTest {
    function test_adminRoleSet() public {
        RBACTimelock rbacTimelock = new RBACTimelock(
            MIN_DELAY,
            ADMIN,
            PROPOSERS,
            EXECUTORS,
            CANCELLERS,
            BYPASSERS
        );
        bool hasAdminRole = rbacTimelock.hasRole(rbacTimelock.ADMIN_ROLE(), ADMIN);
        assertEq(hasAdminRole, true);
    }

    function test_proposersDoNotHaveAdminRole() public {
        RBACTimelock rbacTimelock = new RBACTimelock(
            MIN_DELAY,
            ADMIN,
            PROPOSERS,
            EXECUTORS,
            CANCELLERS,
            BYPASSERS
        );
        BaseTest.checkRoleNotSetForAddresses(rbacTimelock, rbacTimelock.ADMIN_ROLE(), PROPOSERS);
    }

    function test_executorsDoNotHaveAdminRole() public {
        RBACTimelock rbacTimelock = new RBACTimelock(
            MIN_DELAY,
            ADMIN,
            PROPOSERS,
            EXECUTORS,
            CANCELLERS,
            BYPASSERS
        );
        BaseTest.checkRoleNotSetForAddresses(rbacTimelock, rbacTimelock.ADMIN_ROLE(), EXECUTORS);
    }

    function test_cancellersDoNotHaveAdminRole() public {
        RBACTimelock rbacTimelock = new RBACTimelock(
            MIN_DELAY,
            ADMIN,
            PROPOSERS,
            EXECUTORS,
            CANCELLERS,
            BYPASSERS
        );
        BaseTest.checkRoleNotSetForAddresses(rbacTimelock, rbacTimelock.ADMIN_ROLE(), CANCELLERS);
    }

    function test_proposerRolesSet() public {
        RBACTimelock rbacTimelock = new RBACTimelock(
            MIN_DELAY,
            ADMIN,
            PROPOSERS,
            EXECUTORS,
            CANCELLERS,
            BYPASSERS
        );
        assertEq(rbacTimelock.hasRole(rbacTimelock.PROPOSER_ROLE(), PROPOSER_ONE), true);
        assertEq(rbacTimelock.hasRole(rbacTimelock.PROPOSER_ROLE(), PROPOSER_TWO), true);
    }

    function test_adminDoesNotHaveProposerRole() public {
        RBACTimelock rbacTimelock = new RBACTimelock(
            MIN_DELAY,
            ADMIN,
            PROPOSERS,
            EXECUTORS,
            CANCELLERS,
            BYPASSERS
        );
        assertFalse(rbacTimelock.hasRole(rbacTimelock.PROPOSER_ROLE(), ADMIN));
    }

    function test_executorsDoNotHaveProposerRole() public {
        RBACTimelock rbacTimelock = new RBACTimelock(
            MIN_DELAY,
            ADMIN,
            PROPOSERS,
            EXECUTORS,
            CANCELLERS,
            BYPASSERS
        );
        BaseTest.checkRoleNotSetForAddresses(rbacTimelock, rbacTimelock.PROPOSER_ROLE(), EXECUTORS);
    }

    function test_cancellersDoNotHaveProposerRole() public {
        RBACTimelock rbacTimelock = new RBACTimelock(
            MIN_DELAY,
            ADMIN,
            PROPOSERS,
            EXECUTORS,
            CANCELLERS,
            BYPASSERS
        );
        BaseTest.checkRoleNotSetForAddresses(rbacTimelock, rbacTimelock.PROPOSER_ROLE(), CANCELLERS);
    }

    function test_executorRolesSet() public {
        RBACTimelock rbacTimelock = new RBACTimelock(
            MIN_DELAY,
            ADMIN,
            PROPOSERS,
            EXECUTORS,
            CANCELLERS,
            BYPASSERS
        );
        assertEq(rbacTimelock.hasRole(rbacTimelock.EXECUTOR_ROLE(), EXECUTOR_ONE), true);
        assertEq(rbacTimelock.hasRole(rbacTimelock.EXECUTOR_ROLE(), EXECUTOR_TWO), true);
    }

    function test_adminDoesNotHaveExecutorRole() public {
        RBACTimelock rbacTimelock = new RBACTimelock(
            MIN_DELAY,
            ADMIN,
            PROPOSERS,
            EXECUTORS,
            CANCELLERS,
            BYPASSERS
        );
        assertFalse(rbacTimelock.hasRole(rbacTimelock.EXECUTOR_ROLE(), ADMIN));
    }

    function test_proposersDoNotHaveExecutorRole() public {
        RBACTimelock rbacTimelock = new RBACTimelock(
            MIN_DELAY,
            ADMIN,
            PROPOSERS,
            EXECUTORS,
            CANCELLERS,
            BYPASSERS
        );
        BaseTest.checkRoleNotSetForAddresses(rbacTimelock, rbacTimelock.EXECUTOR_ROLE(), PROPOSERS);
    }

    function test_cancellersDoNotHaveExecutorRole() public {
        RBACTimelock rbacTimelock = new RBACTimelock(
            MIN_DELAY,
            ADMIN,
            PROPOSERS,
            EXECUTORS,
            CANCELLERS,
            BYPASSERS
        );
        BaseTest.checkRoleNotSetForAddresses(rbacTimelock, rbacTimelock.EXECUTOR_ROLE(), CANCELLERS);
    }

    function test_cancellerRolesSet() public {
        RBACTimelock rbacTimelock = new RBACTimelock(
            MIN_DELAY,
            ADMIN,
            PROPOSERS,
            EXECUTORS,
            CANCELLERS,
            BYPASSERS
        );
        assertEq(rbacTimelock.hasRole(rbacTimelock.CANCELLER_ROLE(), CANCELLER_ONE), true);
        assertEq(rbacTimelock.hasRole(rbacTimelock.CANCELLER_ROLE(), CANCELLER_TWO), true);
    }

    function test_adminDoesNotHaveCancellerRole() public {
        RBACTimelock rbacTimelock = new RBACTimelock(
            MIN_DELAY,
            ADMIN,
            PROPOSERS,
            EXECUTORS,
            CANCELLERS,
            BYPASSERS
        );
        assertFalse(rbacTimelock.hasRole(rbacTimelock.CANCELLER_ROLE(), ADMIN));
    }

    function test_executorsDoNotHaveCancellerRole() public {
        RBACTimelock rbacTimelock = new RBACTimelock(
            MIN_DELAY,
            ADMIN,
            PROPOSERS,
            EXECUTORS,
            CANCELLERS,
            BYPASSERS
        );
        BaseTest.checkRoleNotSetForAddresses(rbacTimelock, rbacTimelock.CANCELLER_ROLE(), EXECUTORS);
    }

    function test_proposersDoNotHaveCancellerRole() public {
        RBACTimelock rbacTimelock = new RBACTimelock(
            MIN_DELAY,
            ADMIN,
            PROPOSERS,
            EXECUTORS,
            CANCELLERS,
            BYPASSERS
        );
        BaseTest.checkRoleNotSetForAddresses(rbacTimelock, rbacTimelock.CANCELLER_ROLE(), PROPOSERS);
    }

    function test_minDelaySet() public {
        RBACTimelock rbacTimelock = new RBACTimelock(
            MIN_DELAY,
            ADMIN,
            PROPOSERS,
            EXECUTORS,
            CANCELLERS,
            BYPASSERS
        );
        assertEq(rbacTimelock.getMinDelay(), MIN_DELAY);
    }

    function test_noBlockedFunctions() public {
        RBACTimelock rbacTimelock = new RBACTimelock(
            MIN_DELAY,
            ADMIN,
            PROPOSERS,
            EXECUTORS,
            CANCELLERS,
            BYPASSERS
        );
        uint256 numBlockedFns = rbacTimelock.getBlockedFunctionSelectorCount();
        assertEq(numBlockedFns, 0);
    }
}
