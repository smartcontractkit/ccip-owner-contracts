pragma solidity ^0.8.0;

import {ManyChainMultiSig} from "../src/ManyChainMultiSig.sol";

contract ExposedManyChainMultiSig is ManyChainMultiSig {
    function setOpCount(uint40 opCount) public {
        s_expiringRootAndOpCount.opCount = opCount;
    }

    function clearRoot() public {
        s_expiringRootAndOpCount = ExpiringRootAndOpCount({
            root: 0,
            validUntil: 0,
            opCount: s_expiringRootAndOpCount.opCount
        });
    }

    function clearConfig() public {
        delete s_config.signers;
        uint8[NUM_GROUPS] memory groupQuorums;
        s_config.groupQuorums = groupQuorums;
        uint8[NUM_GROUPS] memory groupParents;
        s_config.groupQuorums = groupParents;

    }
}
