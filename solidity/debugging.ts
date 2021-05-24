import { createProtobufRpcClient, QueryClient } from "@cosmjs/stargate";
import { Tendermint34Client } from "@cosmjs/tendermint-rpc";
import {
  Query,
  QueryClientImpl,
  SignerSetTxsRequest,
} from "./gen/gravity/v1/query";
import { SignerSetTx } from "./gen/gravity/v1/gravity";
import Long from "long";
import { exit } from "process";

async function getQueryService(): Promise<Query> {
  const cosmosNode = "http://localhost:26657";
  const tendermintClient = await Tendermint34Client.connect(cosmosNode);
  const queryClient = new QueryClient(tendermintClient);
  const rpcClient = createProtobufRpcClient(queryClient);
  return new QueryClientImpl(rpcClient);
}

async function getParams() {
  let queryService = await getQueryService();
  const res = await queryService.Params({});
  if (!res.params) {
    console.log("Could not retrieve params");
    exit(1);
  }
  return res.params;
}

async function getLatestValset(): Promise<SignerSetTx> {
  let queryService = await getQueryService();
  const req = { nonce: Long.fromInt(0) };
  const res = await queryService.SignerSetTx(req);
  if (!res.signerSet) {
    console.log("Could not retrieve signer set");
    exit(1);
  }
  return res.signerSet;
}

async function getAllValsets() {
  let queryService = await getQueryService();
  const res = await queryService.SignerSetTxs({});

  return res;
}

async function getDelegateKeys() {
  let queryService = await getQueryService();
  const res = await queryService.DelegateKeysByOrchestrator({
    orchestratorAddress: "cosmos14uvqun482ydhljwtvacy5grvgh23xrmgymg0wd",
  });
  return res;
}

(async function () {
  //   console.log(await getDelegateKeys());
  console.log(JSON.stringify(await getAllValsets()));
})();
