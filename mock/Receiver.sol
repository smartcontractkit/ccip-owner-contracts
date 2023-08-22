// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.13;

contract Receiver {
    function executableMethod(bool reverts) public payable {
        if (reverts) {
            require(false, "transaction failed");
        }
    }
}
