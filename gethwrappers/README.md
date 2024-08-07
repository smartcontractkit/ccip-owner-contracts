# Generating Geth Wrappers
## Prerequisites: 
If this is your first time in the repo you need these installed.
Solidity
```bash
brew install solidity
```
Python used to next install solc-select
```bash
brew install python3
```
solc-select is used in script to handle multiple solc versions
```bash
pip3 install solc-select
```
Foundry is for testing
```bash
curl -L https://foundry.paradigm.xyz | bash
```
re-`source` shell like `source /Users/<user>/.zshenv` then
```bash
foundryup
```
Run the foundry tests so you also get the `openzeppelin-contracts` imported to `/lib`
```bash
forge test --ffi
```
## Generating ABIs, BINs and Wrappers
First make your local abigen
```bash
make abigen
```
Run the generations
```bash
go generate gethwrappers/go_generate.go
```
