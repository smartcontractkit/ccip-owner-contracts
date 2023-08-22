// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import "../src/RBACTimelock.sol";
import "../mock/Counter.sol";
import "../src/CallProxy.sol";

contract BaseTest is Test {
    uint256 internal constant MIN_DELAY = 7 days;
    uint256 internal constant DONE_TIMESTAMP = 1;

    address internal constant ADMIN = address(1);
    address internal constant PROPOSER_ONE = address(2);
    address internal constant PROPOSER_TWO = address(3);

    address internal constant EXECUTOR_ONE = address(4);
    address internal constant EXECUTOR_TWO = address(5);

    address internal constant CANCELLER_ONE = address(6);
    address internal constant CANCELLER_TWO = address(7);

    address internal constant BYPASSER_ONE = address(8);
    address internal constant BYPASSER_TWO = address(9);

    bytes32 internal constant NO_PREDECESSOR = bytes32("");
    bytes32 internal constant EMPTY_SALT = bytes32("");

    address[] internal PROPOSERS = new address[](2);
    address[] internal EXECUTORS = new address[](2);
    address[] internal CANCELLERS = new address[](2);
    address[] internal BYPASSERS = new address[](2);

    Counter internal s_counter;
    RBACTimelock internal s_timelock;

    address internal EXTERNAL_CALLER = 0x89205A3A3b2A69De6Dbf7f01ED13B2108B2c43e7;
    CallProxy internal s_proxy;

    function setUp() public virtual {
        PROPOSERS[0] = PROPOSER_ONE;
        PROPOSERS[1] = PROPOSER_TWO;

        EXECUTORS[0] = EXECUTOR_ONE;
        EXECUTORS[1] = EXECUTOR_TWO;

        CANCELLERS[0] = CANCELLER_ONE;
        CANCELLERS[1] = CANCELLER_TWO;

        BYPASSERS[0] = BYPASSER_ONE;
        BYPASSERS[1] = BYPASSER_TWO;

        s_timelock = new RBACTimelock(
            MIN_DELAY,
            ADMIN,
            PROPOSERS,
            EXECUTORS,
            CANCELLERS,
            BYPASSERS
        );
        s_proxy = new CallProxy(address(s_timelock));
        vm.startPrank(ADMIN);
        s_timelock.grantRole(s_timelock.EXECUTOR_ROLE(), address(s_proxy));
        vm.stopPrank();

        s_counter = new Counter(address(s_timelock));

        vm.deal(ADMIN, 1 ether);
    }

    function checkRoleNotSetForAddresses(
        RBACTimelock rbacTimelock,
        bytes32 role,
        address[] memory addresses
    ) internal {
        for (uint256 i = 0; i < addresses.length; i++) {
            assertFalse(rbacTimelock.hasRole(role, addresses[i]));
        }
    }

    function _getExpectedMissingRoleErrorMessage(address account, bytes32 role)
        internal
        pure
        returns (bytes memory)
    {
        return abi.encodePacked(
            "AccessControl: account ",
            Strings.toHexString(account),
            " is missing role ",
            Strings.toHexString(uint256(role), 32)
        );
    }

    // helper function that turns a single RBACTimelock.Call into a singleton
    // slice RBACTimelock.Call[]
    function _singletonCalls(RBACTimelock.Call memory call)
        internal
        pure
        returns (RBACTimelock.Call[] memory)
    {
        RBACTimelock.Call[] memory calls = new RBACTimelock.Call[](1);
        calls[0] = call;
        return calls;
    }
}
