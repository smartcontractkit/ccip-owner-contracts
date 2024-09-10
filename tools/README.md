# Go-MCMS Library

The packages in this directory provide a set of tools that users can use to interact with the MCMS and Timelock contracts

## Deployment & Configuration

### MCMS

Deploy the MCMS contract using the [`DeployManyChainMultiSig`](./gethwrappers/ManyChainMultiSig.go#L76) function

The default [`ManyChainMultiSigConfig`](./gethwrappers/ManyChainMultiSig.go#L32) is very unintuitive to define and construct based on the desired group structure. As a result, this library provides a [`Config` Wrapper](./configwrappers/config.go#L13) which defines a more intuitive MCMS membership structure for ease of use.

The configuration is a nested tree structure where a given group is considered at `quorum` if the sum of `Signers` with Signatures and `GroupSigners` that individually are at `quorum` are greater than or equal to the top-level `quorum`

For example, given the following [`Config`](./configwrappers/config.go#L13):

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
1. [`0x1`, `0x2`, `0x4`]
1. [`0x1`, `0x2`, `0x5`]
1. [`0x1`, `0x2`, `0x6`]
1. [`0x1`, `0x3`, `0x5`]
1. [`0x1`, `0x3`, `0x6`]
1. [`0x1`, `0x4`, `0x5`]
1. [`0x1`, `0x4`, `0x6`]
1. [`0x2`, `0x3`, `0x5`]
1. [`0x2`, `0x3`, `0x6`]
1. [`0x2`, `0x4`, `0x5`]
1. [`0x2`, `0x4`, `0x6`]

Once a satisfactory MCMS Membership configuration is constructed, users can use the [`ExtractSetConfigInputs`](./configwrappers/config.go#L153) function to generate inputs and call [`SetConfig`](./gethwrappers/ManyChainMultiSig.go#L428)

Note: Signers cannot be repeated in this configuration (i.e. they cannot belong to multiple groups)

### Timelock

Deploy the RBACTimelock using the [DeployRBACTimelock](./gethwrappers/RBACTimelock.go#L47) function

Users can configure other addresses with certain roles using the [`GrantRole`](./gethwrappers/RBACTimelock.go#L667) and [`RevokeRole`](./gethwrappers/RBACTimelock.go#L727) functions

Note: These configurations can only be done by the admin, so it's probably easier to set the deployer as the admin until configuration is as desired, then use the [`RenounceRole`](./gethwrappers/RBACTimelock.go#L715) to give up `admin` privileges

### CallProxy

Deploy the call proxy using the [`DeployCallProxy`](./gethwrappers/CallProxy.go#L41) function

Note: the `target` in the CallProxy is only configurable during deployment and cannot be set after the fact

## Proposals

### Construction

### Proposal Validation

### Proposal Signing

### Proposal Execution

## Development

### Upgrade Go-Ethereum Wrappers

From within the [gethwrappers/](./gethwrappers/) directory, run the following command:
```
bash generate.sh
```