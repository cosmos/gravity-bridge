var fs = require("fs");

/*******************************************
 *** Contract type
 ******************************************/
var Contract = function(inPath, outPath) {
  this.inPath = inPath;
  this.outPath = outPath;
};

/*******************************************
 *** Helpers
 ******************************************/
function writeABIs(contracts) {
  for (var i = 0; i < contracts.length; i++) {
    // Set up contract instance using build.json
    const ContractInstance = require(contracts[i].inPath);

    // Write the contract's ABI to ebrelayer/contract/abi
    fs.writeFile(
      contracts[i].outPath,
      JSON.stringify(ContractInstance.abi),
      "utf8",
      function(err) {
        if (err) return console.log(err);
      }
    );
  }
}

/*******************************************
 *** Set up contracts
 ******************************************/
const VALSET_IN_PATH = "../build/contracts/Valset.json";
const VALSET_OUT_PATH = "../cmd/ebrelayer/contract/abi/Valset.abi";

const ORACLE_IN_PATH = "../build/contracts/Oracle.json";
const ORACLE_OUT_PATH = "../cmd/ebrelayer/contract/abi/Oracle.abi";

const COSMOS_BRIDGE_IN_PATH = "../build/contracts/CosmosBridge.json";
const COSMOS_BRIDGE_OUT_PATH = "../cmd/ebrelayer/contract/abi/CosmosBridge.abi";

const BRIDGE_BANK_IN_PATH = "../build/contracts/BridgeBank.json";
const BRIDGE_BANK_OUT_PATH = "../cmd/ebrelayer/contract/abi/BridgeBank.abi";

/*******************************************
 *** Write
 ******************************************/
// Build contracts
const ValsetContract = new Contract(VALSET_IN_PATH, VALSET_OUT_PATH);

const OracleContract = new Contract(ORACLE_IN_PATH, ORACLE_OUT_PATH);

const CosmosBridgeContract = new Contract(
  COSMOS_BRIDGE_IN_PATH,
  COSMOS_BRIDGE_OUT_PATH
);

const BridgeBankContract = new Contract(
  BRIDGE_BANK_IN_PATH,
  BRIDGE_BANK_OUT_PATH
);

// Combine contracts
const contracts = [
  ValsetContract,
  OracleContract,
  CosmosBridgeContract,
  BridgeBankContract
];

// Write contracts to relayer
writeABIs(contracts);
