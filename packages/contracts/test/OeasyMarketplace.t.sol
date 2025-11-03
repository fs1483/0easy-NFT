// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Test} from "forge-std/Test.sol";
import {OeasyMarketplace} from "../src/OeasyMarketplace.sol";
import {OeasyNFT} from "../src/OeasyNFT.sol";
import {MockUSDC} from "../src/MockUSDC.sol";
import {IMarketplace} from "../src/interfaces/IMarketplace.sol";

contract OeasyMarketplaceTest is Test {
    OeasyMarketplace internal marketplace;
    OeasyNFT internal nft;
    MockUSDC internal usdc;

    address internal maker;
    uint256 internal makerKey;

    address internal taker;
    uint256 internal takerKey;

    function setUp() public {
        makerKey = 0xA11CE;
        takerKey = 0xB0B;
        maker = vm.addr(makerKey);
        taker = vm.addr(takerKey);

        marketplace = new OeasyMarketplace();
        marketplace.initialize(address(this));

        nft = new OeasyNFT();
        nft.initialize("Oeasy Mock", "OEASY");
        nft.transferOwnership(address(this));

        usdc = new MockUSDC();
        usdc.initialize("Mock USDC", "mUSDC");

        nft.mintWithId(maker, 1);
        usdc.mint(taker, 1_000 ether);

        vm.prank(maker);
        nft.approve(address(marketplace), 1);

        vm.prank(taker);
        usdc.approve(address(marketplace), type(uint256).max);
    }

    function _buildOrder(address orderMaker, IMarketplace.OrderSide side)
        internal
        view
        returns (IMarketplace.Order memory order)
    {
        order = IMarketplace.Order({
            maker: orderMaker,
            nft: address(nft),
            tokenId: 1,
            paymentToken: address(usdc),
            price: 100 ether,
            expiry: block.timestamp + 1 hours,
            nonce: 1,
            side: side
        });
    }

    function _signOrder(IMarketplace.Order memory order, uint256 privKey)
        internal
        view
        returns (bytes memory signature)
    {
        bytes32 digest = marketplace.hashOrder(order);
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(privKey, digest);
        signature = abi.encodePacked(r, s, v);
    }

    function testExecuteTradeSuccess() public {
        IMarketplace.Order memory makerOrder = _buildOrder(maker, IMarketplace.OrderSide.Ask);
        bytes memory signature = _signOrder(makerOrder, makerKey);

        vm.prank(taker);
        marketplace.executeTrade(makerOrder, _buildOrder(taker, IMarketplace.OrderSide.Bid), signature);

        assertEq(nft.ownerOf(1), taker, "NFT should transfer to taker");
        assertEq(usdc.balanceOf(maker), 100 ether, "Maker should receive payment");
        assertTrue(marketplace.consumedNonces(maker, 1), "Nonce should be consumed");
    }

    function testFeeChargedWhenConfigured() public {
        marketplace.setFeeConfiguration(address(this), 500); // 5%

        IMarketplace.Order memory makerOrder = _buildOrder(maker, IMarketplace.OrderSide.Ask);
        bytes memory signature = _signOrder(makerOrder, makerKey);

        uint256 makerBalanceBefore = usdc.balanceOf(maker);
        uint256 feeRecipientBefore = usdc.balanceOf(address(this));

        vm.prank(taker);
        marketplace.executeTrade(makerOrder, _buildOrder(taker, IMarketplace.OrderSide.Bid), signature);

        assertEq(usdc.balanceOf(maker) - makerBalanceBefore, 95 ether);
        assertEq(usdc.balanceOf(address(this)) - feeRecipientBefore, 5 ether);
    }

    function testPausePreventsExecution() public {
        marketplace.pause();
        IMarketplace.Order memory makerOrder = _buildOrder(maker, IMarketplace.OrderSide.Ask);
        bytes memory signature = _signOrder(makerOrder, makerKey);

        vm.prank(taker);
        vm.expectRevert();
        marketplace.executeTrade(makerOrder, _buildOrder(taker, IMarketplace.OrderSide.Bid), signature);
    }

    function testCannotReuseNonce() public {
        IMarketplace.Order memory makerOrder = _buildOrder(maker, IMarketplace.OrderSide.Ask);
        bytes memory signature = _signOrder(makerOrder, makerKey);

        vm.prank(taker);
        marketplace.executeTrade(makerOrder, _buildOrder(taker, IMarketplace.OrderSide.Bid), signature);

        vm.prank(taker);
        vm.expectRevert(OeasyMarketplace.NonceConsumed.selector);
        marketplace.executeTrade(makerOrder, _buildOrder(taker, IMarketplace.OrderSide.Bid), signature);
    }

    function testCancelOrder() public {
        vm.prank(maker);
        marketplace.cancelOrder(5);
        assertTrue(marketplace.consumedNonces(maker, 5));

        IMarketplace.Order memory makerOrder = _buildOrder(maker, IMarketplace.OrderSide.Ask);
        makerOrder.nonce = 5;
        bytes memory signature = _signOrder(makerOrder, makerKey);

        vm.prank(taker);
        vm.expectRevert(OeasyMarketplace.NonceConsumed.selector);
        marketplace.executeTrade(makerOrder, _buildOrder(taker, IMarketplace.OrderSide.Bid), signature);
    }

    function testCannotExecuteExpiredOrder() public {
        IMarketplace.Order memory makerOrder = _buildOrder(maker, IMarketplace.OrderSide.Ask);
        makerOrder.expiry = block.timestamp - 1;
        bytes memory signature = _signOrder(makerOrder, makerKey);

        vm.prank(taker);
        vm.expectRevert(OeasyMarketplace.InvalidOrder.selector);
        marketplace.executeTrade(makerOrder, _buildOrder(taker, IMarketplace.OrderSide.Bid), signature);
    }

    function testInvalidSignature() public {
        IMarketplace.Order memory makerOrder = _buildOrder(maker, IMarketplace.OrderSide.Ask);
        bytes memory signature = _signOrder(makerOrder, takerKey);

        vm.prank(taker);
        vm.expectRevert(OeasyMarketplace.InvalidSignature.selector);
        marketplace.executeTrade(makerOrder, _buildOrder(taker, IMarketplace.OrderSide.Bid), signature);
    }

    function testUnsupportedSideCombination() public {
        IMarketplace.Order memory makerOrder = _buildOrder(maker, IMarketplace.OrderSide.Bid);
        bytes memory signature = _signOrder(makerOrder, makerKey);

        vm.prank(taker);
        vm.expectRevert(OeasyMarketplace.UnsupportedSide.selector);
        marketplace.executeTrade(makerOrder, _buildOrder(taker, IMarketplace.OrderSide.Bid), signature);
    }
}
