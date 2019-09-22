var fs = require("fs");

/*******************************************
 *** Constants
 ******************************************/
const IN_PATH = "../build/contracts/Peggy.json";
const OUT_PATH = "../cmd/ebrelayer/contract/abi/Peggy.abi";

/*******************************************
 *** Reading and writing
 ******************************************/
const PeggyContract = require(IN_PATH);

fs.writeFile(OUT_PATH, JSON.stringify(PeggyContract.abi), "utf8", function(
  err
) {
  if (err) return console.log(err);
});
