import { ethers } from "hardhat";
import fs from 'fs';
import {Config} from "./config";

async function main() {
  const BurnVerifier = await ethers.getContractFactory("BurnVerifier");
  const ClaimVerifier = await ethers.getContractFactory("ClaimVerifier");
  const OracleMock = await ethers.getContractFactory("OracleMock");
  const Rollup = await ethers.getContractFactory("Rollup");

  const burnVerifier = await BurnVerifier.deploy();
  const claimVerifier = await ClaimVerifier.deploy();
  const oracleMock = await OracleMock.deploy();

  await burnVerifier.waitForDeployment();
  await claimVerifier.waitForDeployment();
  await oracleMock.waitForDeployment();

  const burnVerifierAddress = await burnVerifier.getAddress();
  const claimVerifierAddress = await claimVerifier.getAddress();
  const oracleMockAddress = await oracleMock.getAddress();

  console.log("BurnVerifier deployed to:", burnVerifierAddress);
  console.log("ClaimVerifier deployed to:", claimVerifierAddress);
  console.log("OracleMock deployed to:", oracleMockAddress);

  // 4 accounts
  //const rollup = await Rollup.deploy("3065278908848025531261417432545150983396361703868879929916746847153394764839", burnVerifierAddress, claimVerifierAddress, oracleMockAddress);

  // 65536 accounts
  const rollup = await Rollup.deploy("5069628288876998019696692163615553505271603549100099079061297050627237284397", burnVerifierAddress, claimVerifierAddress, oracleMockAddress);

  // 16 accounts
  //const rollup = await Rollup.deploy("18992696670841424534069916451357995059717643186344737472388881927587660171461", burnVerifierAddress, claimVerifierAddress);

  await rollup.waitForDeployment();
  const rollupAddress = await rollup.getAddress();

  console.log("Rollup deployed to:", rollupAddress);

  const config = readConfig("../operator/configs/config.json");
  config.Ethereum.burnVerifierContract = burnVerifierAddress;
  config.Ethereum.claimVerifierContract = claimVerifierAddress;
  config.Ethereum.oracleMockContract = oracleMockAddress;
  config.Ethereum.rollupContract = rollupAddress;

  writeJsonToFile("../operator/configs/config.json", config);
}
function readConfig(filePath: string): Config {
  try {
    const rawData = fs.readFileSync(filePath, { encoding: 'utf8' });
    const config: Config = JSON.parse(rawData);
    return config;
  } catch (error) {
    console.error('Error reading the JSON file:', error);
  }
}

function writeJsonToFile(filePath: string, data: Config): void {
  try {
    const jsonData = JSON.stringify(data, null, 2);
    fs.writeFileSync(filePath, jsonData, { encoding: 'utf8' });
    console.log('Updated config file.');
  } catch (error) {
    console.error('Error writing the JSON file:', error);
  }
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
