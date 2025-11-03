// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {ERC20Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

/// @title MockUSDC
/// @notice Upgradeable mock stablecoin used to settle trades within the MVP environment.
/// @dev Owner-controlled minting allows seeding trader balances in local and testnet environments.
contract MockUSDC is Initializable, ERC20Upgradeable, OwnableUpgradeable {
    /// @notice Initializes the upgradeable ERC20 mock contract with name and symbol.
    /// @param _name Token name (e.g. "Mock USDC").
    /// @param _symbol Token ticker (e.g. "mUSDC").
    function initialize(string memory _name, string memory _symbol) external initializer {
        __ERC20_init(_name, _symbol);
        __Ownable_init(msg.sender);
    }

    /// @notice Mints mock stablecoins to the supplied account.
    /// @dev For testnet/development: Anyone can mint for easy testing.
    ///      Production contracts should use onlyOwner or role-based access control.
    /// @param account Recipient address.
    /// @param amount Token amount in the smallest denomination (6 decimals for USDC).
    function mint(address account, uint256 amount) external {
        // 测试环境：允许任何人铸造
        // 生产环境：应该使用 onlyOwner 或 AccessControl
        _mint(account, amount);
    }
    
    /// @notice Returns the number of decimals (USDC uses 6, not 18).
    function decimals() public pure override returns (uint8) {
        return 6;
    }
}
