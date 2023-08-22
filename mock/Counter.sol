// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.13;

contract Counter {
    address private s_timelock;
    uint256 private s_number;

    constructor(address timelock) {
        s_timelock = timelock;
    }

    function setNumber(uint256 newNumber) public onlyTimelock {
        s_number = newNumber;
    }

    function increment() public onlyTimelock {
        s_number++;
    }

    function mockRevert() public pure {
        require(false, "Transaction reverted");
    }

    function number() external view returns (uint256) {
        return s_number;
    }

    modifier onlyTimelock() {
        require(msg.sender == s_timelock, "Not timelock controller");
        _;
    }
}
