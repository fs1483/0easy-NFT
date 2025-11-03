// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IERC721} from "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {EIP712Upgradeable} from "@openzeppelin/contracts-upgradeable/utils/cryptography/EIP712Upgradeable.sol";
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {IMarketplace} from "./interfaces/IMarketplace.sol";
import {OrderLib} from "./libraries/OrderLib.sol";

/// @title OeasyMarketplace
/// @notice Upgradeable orderbook-based NFT marketplace executing pre-signed orders on-chain.
/// @dev Implements the core settlement logic for a hybrid off-chain orderbook with on-chain atomic settlement.
contract OeasyMarketplace is
    Initializable,
    OwnableUpgradeable,
    UUPSUpgradeable,
    EIP712Upgradeable,
    PausableUpgradeable,
    ReentrancyGuardUpgradeable
{
    using ECDSA for bytes32;

    uint256 private constant PERCENTAGE_SCALE = 10_000; // Basis points denominator (10000 = 100%)

    /// @notice Mapping storing order nonces that have either been executed or cancelled.
    mapping(address => mapping(uint256 => bool)) public consumedNonces;

    /// @notice Address receiving platform fees. Zero address disables fee collection.
    address public feeRecipient;

    /// @notice Fee rate expressed in basis points (1% = 100). Applied to taker payment.
    uint96 public feeBps;

    /// @notice Emitted when a trade is successfully executed on-chain.
    event TradeExecuted(
        address indexed maker,
        address indexed taker,
        address indexed nft,
        uint256 tokenId,
        address paymentToken,
        uint256 price,
        IMarketplace.OrderSide side,
        uint256 fee
    );

    /// @notice Emitted when an order nonce is cancelled by its maker.
    event OrderCancelled(address indexed maker, uint256 indexed nonce);

    /// @notice Emitted when fee configuration changes.
    event FeeUpdated(address indexed recipient, uint96 feeBps);

    error InvalidSignature();
    error InvalidOrder();
    error NonceConsumed();
    error MakerCannotBeTaker();
    error UnsupportedSide();
    error InvalidFeeConfiguration();
    error TransferFailed();

    /// @notice Initializes the upgradeable contract components.
    /// @param owner_ Address that obtains contract ownership rights.
    function initialize(address owner_) external initializer {
        if (owner_ == address(0)) revert InvalidOrder();
        __Ownable_init(owner_);
        __UUPSUpgradeable_init();
        __ReentrancyGuard_init();
        __Pausable_init();
        __EIP712_init("Oeasy Marketplace", "1");
    }

    /// @notice Owner can configure platform fee recipient and rate.
    /// @param recipient Address receiving fees.
    /// @param feeBps_ Fee in basis points (max 100%).
    function setFeeConfiguration(address recipient, uint96 feeBps_) external onlyOwner {
        if (feeBps_ > PERCENTAGE_SCALE) revert InvalidFeeConfiguration();
        feeRecipient = recipient;
        feeBps = feeBps_;
        emit FeeUpdated(recipient, feeBps_);
    }

    /// @notice Pauses settlement in case of emergency.
    function pause() external onlyOwner {
        _pause();
    }

    /// @notice Resumes settlement operations.
    function unpause() external onlyOwner {
        _unpause();
    }

    /// @notice Cancels a previously signed order by marking its nonce as consumed.
    /// @param nonce Unique order nonce to invalidate.
    function cancelOrder(uint256 nonce) external whenNotPaused {
        _consumeNonce(msg.sender, nonce);
        emit OrderCancelled(msg.sender, nonce);
    }

    /// @notice Executes a trade between a signed maker order and a taker order.
    /// @param makerOrder Fully signed order from the maker (seller for ask).
    /// @param takerOrder Order from the taker (buyer for bid), must match maker fields (except maker address and side).
    /// @param makerSignature Off-chain signature corresponding to `makerOrder`.
    function executeTrade(
        IMarketplace.Order calldata makerOrder,
        IMarketplace.Order calldata takerOrder,
        bytes calldata makerSignature
    ) external whenNotPaused nonReentrant {
        _validateOrders(makerOrder, takerOrder);
        _consumeNonce(makerOrder.maker, makerOrder.nonce);
        _consumeNonce(takerOrder.maker, takerOrder.nonce);  // 消费 taker 的 nonce

        // taker 地址从订单中获取，不是 msg.sender
        // 这样允许第三方（如执行服务）代为提交交易
        address takerAddr = takerOrder.maker;
        if (takerAddr == makerOrder.maker) revert MakerCannotBeTaker();

        // Verify maker signature
        bytes32 makerDigest = _hashTypedDataV4(OrderLib.hash(makerOrder));
        address recoveredMaker = makerDigest.recover(makerSignature);
        if (recoveredMaker != makerOrder.maker) revert InvalidSignature();
        
        // 注意：taker 订单的签名已在链下验证（订单服务）
        // 链上不再验证 taker 签名以节省 Gas

        uint256 feeAmount = _settleTrade(makerOrder, takerAddr);

        emit TradeExecuted(
            makerOrder.maker,
            takerAddr,
            makerOrder.nft,
            makerOrder.tokenId,
            makerOrder.paymentToken,
            makerOrder.price,
            makerOrder.side,
            feeAmount
        );
    }

    /// @notice Computes the EIP-712 digest for an order (utility for off-chain services).
    function hashOrder(IMarketplace.Order calldata order) external view returns (bytes32 digest) {
        // TODO: [Security] - Evaluate necessity of exposing this helper in production deployments; backend services can replicate hashing logic independently.
        digest = _hashTypedDataV4(OrderLib.hash(order));
    }

    function _consumeNonce(address maker, uint256 nonce) internal {
        mapping(uint256 => bool) storage makerNonces = consumedNonces[maker];
        if (makerNonces[nonce]) revert NonceConsumed();
        makerNonces[nonce] = true;
    }

    function _validateOrders(IMarketplace.Order calldata maker, IMarketplace.Order calldata taker) internal view {
        // 验证 NFT 和支付代币必须相同
        if (maker.nft != taker.nft || maker.tokenId != taker.tokenId) revert InvalidOrder();
        if (maker.paymentToken != taker.paymentToken) revert InvalidOrder();
        
        // 验证订单未过期
        if (maker.expiry < block.timestamp || taker.expiry < block.timestamp) revert InvalidOrder();

        // 验证订单方向：maker 必须是 Ask（卖单），taker 必须是 Bid（买单）
        if (maker.side == taker.side) revert UnsupportedSide();
        if (maker.side != IMarketplace.OrderSide.Ask || taker.side != IMarketplace.OrderSide.Bid) {
            revert UnsupportedSide();
        }
        
        // 验证价格匹配：买单价格必须 >= 卖单价格（订单簿标准规则）
        // 企业级标准：bid >= ask，成交价按 ask 价格
        if (taker.price < maker.price) revert InvalidOrder();
    }

    function _settleTrade(IMarketplace.Order calldata makerOrder, address taker)
        internal
        returns (uint256 feeAmount)
    {
        IERC20 paymentToken = IERC20(makerOrder.paymentToken);
        uint256 price = makerOrder.price;

        if (feeRecipient != address(0) && feeBps > 0) {
            feeAmount = (price * feeBps) / PERCENTAGE_SCALE;
            if (feeAmount > 0) {
                bool feeOk = paymentToken.transferFrom(taker, feeRecipient, feeAmount);
                if (!feeOk) revert TransferFailed();
            }
        }

        uint256 netAmount = price - feeAmount;
        bool paymentSuccess = paymentToken.transferFrom(taker, makerOrder.maker, netAmount);
        if (!paymentSuccess) revert TransferFailed();

        IERC721(makerOrder.nft).safeTransferFrom(makerOrder.maker, taker, makerOrder.tokenId);
    }

    /// @inheritdoc UUPSUpgradeable
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}
