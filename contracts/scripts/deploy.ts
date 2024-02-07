import { ethers } from "hardhat";

async function main() {
  const BurnVerifier = await ethers.getContractFactory("BurnVerifier");
  const ClaimVerifier = await ethers.getContractFactory("ClaimVerifier");
  const Rollup = await ethers.getContractFactory("Rollup");

  const burnVerifier = await BurnVerifier.deploy();
  const claimVerifier = await ClaimVerifier.deploy();

  await burnVerifier.waitForDeployment();
  await claimVerifier.waitForDeployment();

  const burnVerifierAddress = await burnVerifier.getAddress();
  const claimVerifierAddress = await claimVerifier.getAddress();

  console.log("BurnVerifier deployed to:", burnVerifierAddress);
  console.log("ClaimVerifier deployed to:", claimVerifierAddress);

  const rollup = await Rollup.deploy(0, burnVerifierAddress, claimVerifierAddress);
  await rollup.waitForDeployment();

  const result = await burnVerifier.compressProof([0,0,0,0,0,0,0,0]);
  console.log("Result:", result);

  console.log("Rollup deployed to:", await rollup.getAddress());
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
