// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/// @title IMarketplace
/// @notice Interface exposing shared structs used by the Oeasy marketplace ecosystem.
interface IMarketplace {
    /// @notice Order side enumerates maker intent.
    enum OrderSide {
        Ask,
        Bid
    }

    /// @notice Order struct representing a signed intent to trade an asset.
    /// @dev Both buy and sell orders share the same structure. Side is determined by context (maker/taker).
    struct Order {
        address maker; // Order creator (seller for ask, buyer for bid)
        address nft; // ERC721 collection address
        uint256 tokenId; // Token being traded
        address paymentToken; // ERC20 token used for payment
        uint256 price; // Price for the entire token (no partial fills in MVP)
        uint256 expiry; // Expiration timestamp
        uint256 nonce; // Unique nonce per maker to support cancellation
        OrderSide side; // Ask or bid
    }
}
