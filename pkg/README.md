# Go-MCMS Library

The packages in this directory provide a set of tools that users can use to interact with the MCMS and Timelock contracts

## Deployment & Configuration

### MCMS

Deploy the MCMS contract using the [`DeployManyChainMultiSig`](./gethwrappers/ManyChainMultiSig.go#L76) function

The default [`ManyChainMultiSigConfig`](./gethwrappers/ManyChainMultiSig.go#L32) is very unintuitive to define and construct based on the desired group structure. As a result, this library provides a [`Config` Wrapper](./config/config.go#L13) which defines a more intuitive MCMS membership structure for ease of use.

The configuration is a nested tree structure where a given group is considered at `quorum` if the sum of `Signers` with Signatures and `GroupSigners` that individually are at `quorum` are greater than or equal to the top-level `quorum`

For example, given the following [`Config`](./config/config.go#L13):

```
Config{
    Quorum:       3,
    Signers:      ["0x1", "0x2"],
    GroupSigners: [
        {
            Quorum: 1,
            Signers: ["0x3","0x4"],
            GroupSigners: []
        },
        {
            Quorum: 1,
            Signers: ["0x5","0x6"],
            GroupSigners: []
        }
    ],
}
```

This configuration represents a membership structure that requires 3 entities to approve, in which any of the following combinations of signatures would satisfy the top-level quorum of `3`:

1. [`0x1`, `0x2`, `0x3`]
2. [`0x1`, `0x2`, `0x4`]
3. [`0x1`, `0x2`, `0x5`]
4. [`0x1`, `0x2`, `0x6`]
5. [`0x1`, `0x3`, `0x5`]
6. [`0x1`, `0x3`, `0x6`]
7. [`0x1`, `0x4`, `0x5`]
8. [`0x1`, `0x4`, `0x6`]
9. [`0x2`, `0x3`, `0x5`]
10. [`0x2`, `0x3`, `0x6`]
11. [`0x2`, `0x4`, `0x5`]
12. [`0x2`, `0x4`, `0x6`]

Once a satisfactory MCMS Membership configuration is constructed, users can use the [`ExtractSetConfigInputs`](./config/config.go#L153) function to generate inputs and call [`SetConfig`](./gethwrappers/ManyChainMultiSig.go#L428)

Note: Signers cannot be repeated in this configuration (i.e. they cannot belong to multiple groups)

### Timelock

Deploy the RBACTimelock using the [DeployRBACTimelock](./gethwrappers/RBACTimelock.go#L47) function

Users can configure other addresses with certain roles using the [`GrantRole`](./gethwrappers/RBACTimelock.go#L667) and [`RevokeRole`](./gethwrappers/RBACTimelock.go#L727) functions

Note: These configurations can only be done by the admin, so it's probably easier to set the deployer as the admin until configuration is as desired, then use the [`RenounceRole`](./gethwrappers/RBACTimelock.go#L715) to give up `admin` privileges

### CallProxy

Deploy the call proxy using the [`DeployCallProxy`](./gethwrappers/CallProxy.go#L41) function

Note: the `target` in the CallProxy is only configurable during deployment and cannot be set after the fact

## Proposals

Once relevant MCMS/RBACTimelock/CallProxy contracts are deployed, the way users can interact with these contracts is through a [`Proposal`](./proposal/mcms/proposal.go#L18). At it's core, a `Proposal` is just a list of (currently only EVM) operations that are to be executed through the MCMS, along with additional metadata about individual transactions and the proposal as a whole. Proposals come in two flavors:

1. [`MCMSProposal`](./proposal/mcms/proposal.go#L18): Represents a simple list of operations (`to`,`value`,`data`) that are to be executed through the mcms with no transformation
2. [`MCMSWithTimelockProposal`](./proposal/timelock/mcm_with_timelock.go#L24): Represents a list of operations that are to be wrapped in a given timelock operation (`Schedule`,`Cancel`,`Bypass`) before being executed through the MCMS. More details about this flavor can be found [below](#nuances-of-mcmswithtimelockproposals)

### Construction

Proposal types can be constructed with their respective `NewProposal...` functions. For example, [`NewMCMSWithTimelockProposal`](./proposal/timelock/mcm_with_timelock.go#L38) and [`NewMCMSProposal`](./proposal/mcms/proposal.go#L36)

### Proposal Validation

Proposal types all contain a relevant `Validate()` function that can be used to validate the proposal format. This function is executed by default when using the constructors but for proposals that are incrementally constructed, this function can be used to revalidate.

### Proposal Signing

`cd cmd && go run main.go --help`

### Proposal Execution

The library provides two functions to help with the execution of an MCMS Proposal:

1. [`SetRootOnChain`](./proposal/mcms/executor.go#L234): Given auth and a ChainIdentifer, calls `setRoot` on the target MCMS for that given chainIdentifier.
2. [`ExecuteOnChain`](./proposal/mcms/executor.go#L269): Given auth and an index, calls `execute` on the target MCMS for that given operation.

### Nuances of MCMSWithTimelockProposals

The [`MCMSWithTimelockProposal`](./proposal/timelock/mcm_with_timelock.go#L24) is an extension of the `MCMSProposal` and has the following additional fields:

1. `Operation`: One of <`Schedule` | `Cancel` | `Bypass`> which determines how to wrap each call in `Transactions`, wrapping calls in `scheduleBatch`,`cancel`, and `bypasserExecuteBatch` calls, respectively
2. `MinDelay` is a string representation of a Go `time.Duration` ("1s", "1h", "1d", etc.). This field is only required when `Operation == Schedule` and sets the delay for each transaction to be the provided value in seconds
3. `TimelockAddresses` is a map of `ChainIdentifier` to the target `RBACTimelock` address for each chain.
4. Each element in `Transactions` is now an array of operations that are all to be wrapped in a single `scheduleBatch` or `bypasserExecuteBatch` call and executed atomically. There is no concept of batching natively available in the MCMS contract which is why this is only available in the RBACTimelock flavor.

## Development

### Upgrade Go-Ethereum Wrappers

From within the [gethwrappers/](./gethwrappers/) directory, run the following command:
```
bash generate.sh
```