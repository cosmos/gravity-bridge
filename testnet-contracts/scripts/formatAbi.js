var fs = require("fs");

/*******************************************
 *** Constants
 ******************************************/
const IN_PATH = "../build/contracts/BridgeBank.json";
const OUT_PATH = "../cmd/ebrelayer/contract/abi/BridgeBank.abi";

/*******************************************
 *** Reading and writing
 ******************************************/
const ContractInstance = require(IN_PATH);

fs.writeFile(OUT_PATH, JSON.stringify(ContractInstance.abi), "utf8", function(
  err
) {
  if (err) return console.log(err);
});
