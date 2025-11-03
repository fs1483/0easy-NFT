// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Script, console2} from "forge-std/Script.sol";
import {OeasyMarketplace} from "../src/OeasyMarketplace.sol";
import {OeasyNFT} from "../src/OeasyNFT.sol";
import {MockUSDC} from "../src/MockUSDC.sol";

/// @title DeployScript
/// @notice Foundry deployment script for spinning up marketplace components on Sepolia or local forks.
contract DeployScript is Script {
    function run() external {
        // TODO: [Security] - Load deployer keys via encrypted key management (e.g. environment variables with sops) for production readiness.
        uint256 deployerKey = vm.envUint("DEPLOYER_PRIVATE_KEY");
        address deployer = vm.addr(deployerKey);

        vm.startBroadcast(deployerKey);

        OeasyMarketplace marketplace = new OeasyMarketplace();
        marketplace.initialize(deployer);

        OeasyNFT nft = new OeasyNFT();
        nft.initialize("Oeasy Mock", "OEASY");
        nft.transferOwnership(deployer);

        MockUSDC usdc = new MockUSDC();
        usdc.initialize("Mock USDC", "mUSDC");
        usdc.transferOwnership(deployer);

        vm.stopBroadcast();

        console2.log("OeasyMarketplace", address(marketplace));
        console2.log("OeasyNFT", address(nft));
        console2.log("MockUSDC", address(usdc));
    }
}
