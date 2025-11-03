// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {IMarketplace} from "../interfaces/IMarketplace.sol";

/// @title OrderLib
/// @notice EIP-712 helpers for encoding marketplace orders.
/// @dev This library is intentionally lightweight to keep calldata size minimal.
library OrderLib {
    bytes32 internal constant ORDER_TYPEHASH = keccak256(
        "Order(address maker,address nft,uint256 tokenId,address paymentToken,uint256 price,uint256 expiry,uint256 nonce,uint8 side)"
    );

    /// @notice Computes the EIP-712 struct hash for a marketplace order.
    /// @param order Struct containing order fields.
    /// @return digest Keccak-256 hash of the encoded order.
    function hash(IMarketplace.Order memory order) internal pure returns (bytes32 digest) {
        digest = keccak256(
            abi.encode(
                ORDER_TYPEHASH,
                order.maker,
                order.nft,
                order.tokenId,
                order.paymentToken,
                order.price,
                order.expiry,
                order.nonce,
                order.side
            )
        );
    }
}
