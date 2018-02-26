var Peggy = artifacts.require("Peggy")

module.exports = function(deployer, network, accounts) {
  // Use deployer to state migration tasks.
  var valset
  switch (network) {
    case "ganache":
    default:
      valset = [[accounts[0]], [100]]
      break
  }
  deployer.deploy(Peggy, valset[0], valset[1])
}
