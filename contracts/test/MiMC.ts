import { ethers } from "hardhat";
import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";
import { expect } from "chai";

describe("MiMC", function () {
    async function deployFixture() {
        const MiMC = await ethers.getContractFactory("MiMC");
        const miMC = await MiMC.deploy();
        return { miMC };
    }

    it("hash", async function () {
        const { miMC } = await deployFixture();
        const hash = await miMC.hash([1,2])
        console.log(hash.toString());
    });
})