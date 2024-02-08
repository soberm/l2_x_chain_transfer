import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";
import "@nomiclabs/hardhat-solhint";
import "hardhat-abi-exporter";
import "@nomiclabs/hardhat-ganache";

const config: HardhatUserConfig = {
  solidity: "0.8.19",
  networks: {
    localganache: {
        url: "HTTP://127.0.0.1:7545",
        //chainId: 1337,
    }
  }
};

export default config;
