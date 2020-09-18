## Early MVP
Happy path implementations

### Oracle

#### Assumptions
- An orchestrator may want to submit multiple claims with a msg (withdrawal batch update + MultiSig Set update )
- Nonces are not unique without a context (withdrawal nonce and MultiSig Set update can have same nonce (=height))
- A nonce is unique in it's context and never reused
- Multiple claims by an orchestrator for the same ETH event are forbidden
- We know the ETH event types beforehand (and handle them as ClaimTypes) 
- For an **observation** status in Attestation the power AND count thresholds must be exceeded
- Fraction type allows higher precision math than %. For example with 2/3

A good start to follow the process would be the `x/peggy/handler_test.go` file


### Not covered/ implemented
- [ ] unhappy cases
- [ ] proper unit + integration tests
- [ ] message validation
- [ ] Genesis I/O
- [ ] Parameters
- [ ] authZ: EthereumChainID whitelisted
- [ ] authZ: bridge contract address whitelisted

