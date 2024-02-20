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

  // 4 accounts
  const rollup = await Rollup.deploy("3065278908848025531261417432545150983396361703868879929916746847153394764839", burnVerifierAddress, claimVerifierAddress);

  // 16 accounts
  //const rollup = await Rollup.deploy("18992696670841424534069916451357995059717643186344737472388881927587660171461", burnVerifierAddress, claimVerifierAddress);

  await rollup.waitForDeployment();

  console.log("Rollup deployed to:", await rollup.getAddress());
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
