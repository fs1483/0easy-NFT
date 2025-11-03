// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {ERC721Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC721/ERC721Upgradeable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

/// @title OeasyNFT
/// @notice Upgradeable ERC721 token used for integration tests and demonstration within the marketplace MVP.
/// @dev Simplified NFT collection with owner-controlled minting. Not intended for production deployment without further hardening.
contract OeasyNFT is Initializable, ERC721Upgradeable, OwnableUpgradeable, UUPSUpgradeable {
    /// @notice Counter used to generate incremental token identifiers when using {mintTo}.
    uint256 public nextTokenId;

    /// @notice Initializes the upgradeable ERC721 contract.
    /// @param _name Collection name.
    /// @param _symbol Collection ticker symbol.
    function initialize(string memory _name, string memory _symbol) external initializer {
        __ERC721_init(_name, _symbol);
        __Ownable_init(msg.sender);
        __UUPSUpgradeable_init();

        nextTokenId = 1;
    }

    /// @notice Mints a new token to the specified recipient with an auto-incrementing identifier.
    /// @dev Only callable by the contract owner. This helper is primarily for testing and seeding liquidity.
    /// @param recipient Address receiving the freshly minted NFT.
    /// @return tokenId Newly minted token identifier.
    function mintTo(address recipient) external onlyOwner returns (uint256 tokenId) {
        tokenId = nextTokenId++;
        _safeMint(recipient, tokenId);
    }

    /// @notice Mints a specific token ID to the provided recipient.
    /// @dev Allows pre-minting deterministic identifiers for test scenarios.
    /// @param recipient Address receiving the minted NFT.
    /// @param tokenId Identifier to be minted. Must not already exist.
    function mintWithId(address recipient, uint256 tokenId) external onlyOwner {
        // TODO: [Validation] - Consider restricting tokenId ranges or enforcing uniqueness constraints beyond the default `_mint` checks for production collections.
        if (tokenId >= nextTokenId) {
            nextTokenId = tokenId + 1;
        }
        _safeMint(recipient, tokenId);
    }

    /// @inheritdoc UUPSUpgradeable
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}
