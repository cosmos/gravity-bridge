import chai from "chai";
import { ethers } from "hardhat";
import { solidity } from "ethereum-waffle";

import { deployContracts } from "../test-utils";
import {
    getSignerAddresses,
    signHash,
    examplePowers
} from "../test-utils/pure";

chai.use(solidity);
const { expect } = chai;

describe("Gas tests", function () {
    it("makeCheckpoint in isolation", async function () {
        const signers = await ethers.getSigners();
        const peggyId = ethers.utils.formatBytes32String("foo");

        // This is the power distribution on the Cosmos hub as of 7/14/2020
        let powers = examplePowers();
        let validators = signers.slice(0, powers.length);

        const powerThreshold = 6666;

        const {
            peggy,
            testERC20,
            checkpoint: deployCheckpoint
        } = await deployContracts(peggyId, validators, powers, powerThreshold);

        await peggy.testMakeCheckpoint(
            await getSignerAddresses(validators),
            powers,
            0,
            peggyId
        );
    });

    it("checkValidatorSignatures in isolation", async function () {
        const signers = await ethers.getSigners();
        const peggyId = ethers.utils.formatBytes32String("foo");

        // This is the power distribution on the Cosmos hub as of 7/14/2020
        let powers = examplePowers();
        let validators = signers.slice(0, powers.length);

        const powerThreshold = 6666;

        const {
            peggy,
            testERC20,
            checkpoint: deployCheckpoint
        } = await deployContracts(peggyId, validators, powers, powerThreshold);

        let sigs = await signHash(
            validators,
            "0x7bc422a00c175cae98cf2f4c36f2f8b63ec51ab8c57fecda9bccf0987ae2d67d"
        );

        await peggy.testCheckValidatorSignatures(
            await getSignerAddresses(validators),
            powers,
            sigs.v,
            sigs.r,
            sigs.s,
            "0x7bc422a00c175cae98cf2f4c36f2f8b63ec51ab8c57fecda9bccf0987ae2d67d",
            6666
        );
    });
});
