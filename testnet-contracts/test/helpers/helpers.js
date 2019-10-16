function fixSignature(signature) {
  // in geth its always 27/28, in ganache its 0/1. Change to 27/28 to prevent
  // signature malleability if version is 0/1
  // see https://github.com/ethereum/go-ethereum/blob/v1.8.23/internal/ethapi/api.go#L465
  let v = parseInt(signature.slice(130, 132), 16);
  if (v < 27) {
    v += 27;
  }
  const vHex = v.toString(16);
  return signature.slice(0, 130) + vHex;
}

function toEthSignedMessageHash(messageHex) {
  const messageBuffer = Buffer.from(messageHex.substring(2), "hex");
  const prefix = Buffer.from(
    `\u0019Ethereum Signed Message:\n${messageBuffer.length}`
  );
  return web3.utils.sha3(Buffer.concat([prefix, messageBuffer]));
}

module.exports = {
  toEthSignedMessageHash,
  fixSignature
};
