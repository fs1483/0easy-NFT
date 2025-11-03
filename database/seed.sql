-- ============================================
-- Oeasy-NFT 测试数据脚本
-- ============================================

\c oeasy_nft;

-- 使用实际部署的合约地址
-- Marketplace: 0x5FbDB2315678afecb367f032d93F642f64180aa3
-- OeasyNFT:    0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0
-- MockUSDC:    0x5FC8d32690cc91D4c39d9d3abcBD16989F875707

-- 卖单 1: User1 出售 NFT #1，价格 100 USDC
INSERT INTO orders (maker, nft_address, token_id, payment_token, price, expiry, nonce, side, status, signature, hash)
VALUES (
    '0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266',
    '0x9fe46736679d2d9a65f0992f2272de9f3c7fa6e0',
    1,
    '0x5fc8d32690cc91d4c39d9d3abcbd16989f875707',
    100000000,
    CURRENT_TIMESTAMP + INTERVAL '7 days',
    1729000000000001,
    'ask', 'active',
    '0x' || repeat('a', 130),
    '0x' || repeat('1', 64)
);

-- 卖单 2: User1 出售 NFT #2，价格 200 USDC
INSERT INTO orders (maker, nft_address, token_id, payment_token, price, expiry, nonce, side, status, signature, hash)
VALUES (
    '0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266',
    '0x9fe46736679d2d9a65f0992f2272de9f3c7fa6e0',
    2,
    '0x5fc8d32690cc91d4c39d9d3abcbd16989f875707',
    200000000,
    CURRENT_TIMESTAMP + INTERVAL '7 days',
    1729000000000002,
    'ask', 'active',
    '0x' || repeat('b', 130),
    '0x' || repeat('2', 64)
);

-- 买单 1: User2 购买 NFT #1，出价 95 USDC (低于卖价，不会匹配)
INSERT INTO orders (maker, nft_address, token_id, payment_token, price, expiry, nonce, side, status, signature, hash)
VALUES (
    '0x70997970c51812dc3a010c7d01b50e0d17dc79c8',
    '0x9fe46736679d2d9a65f0992f2272de9f3c7fa6e0',
    1,
    '0x5fc8d32690cc91d4c39d9d3abcbd16989f875707',
    95000000,
    CURRENT_TIMESTAMP + INTERVAL '7 days',
    1729000000000003,
    'bid', 'active',
    '0x' || repeat('c', 130),
    '0x' || repeat('3', 64)
);

-- 买单 2: User3 购买 NFT #2，出价 210 USDC (高于卖价，会匹配)
INSERT INTO orders (maker, nft_address, token_id, payment_token, price, expiry, nonce, side, status, signature, hash)
VALUES (
    '0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc',
    '0x9fe46736679d2d9a65f0992f2272de9f3c7fa6e0',
    2,
    '0x5fc8d32690cc91d4c39d9d3abcbd16989f875707',
    210000000,
    CURRENT_TIMESTAMP + INTERVAL '3 days',
    1729000000000004,
    'bid', 'active',
    '0x' || repeat('d', 130),
    '0x' || repeat('4', 64)
);

-- 已取消的订单
INSERT INTO orders (maker, nft_address, token_id, payment_token, price, expiry, nonce, side, status, signature, hash)
VALUES (
    '0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266',
    '0x9fe46736679d2d9a65f0992f2272de9f3c7fa6e0',
    3,
    '0x5fc8d32690cc91d4c39d9d3abcbd16989f875707',
    150000000,
    CURRENT_TIMESTAMP + INTERVAL '7 days',
    1729000000000005,
    'ask', 'cancelled',
    '0x' || repeat('e', 130),
    '0x' || repeat('5', 64)
);

-- 已成交的订单
INSERT INTO orders (maker, nft_address, token_id, payment_token, price, expiry, nonce, side, status, signature, hash)
VALUES (
    '0x70997970c51812dc3a010c7d01b50e0d17dc79c8',
    '0x9fe46736679d2d9a65f0992f2272de9f3c7fa6e0',
    10,
    '0x5fc8d32690cc91d4c39d9d3abcbd16989f875707',
    50000000,
    CURRENT_TIMESTAMP + INTERVAL '7 days',
    1729000000000006,
    'bid', 'filled',
    '0x' || repeat('f', 130),
    '0x' || repeat('6', 64)
);

-- 交易事件
INSERT INTO trade_events (transaction_hash, log_index, block_number, maker, taker, nft_address, token_id, payment_token, price, side, fee)
VALUES ('0x' || repeat('a', 64), 0, 12345, '0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266', '0x70997970c51812dc3a010c7d01b50e0d17dc79c8', '0x9fe46736679d2d9a65f0992f2272de9f3c7fa6e0', 10, '0x5fc8d32690cc91d4c39d9d3abcbd16989f875707', 50000000, 0, 500000);

INSERT INTO trade_events (transaction_hash, log_index, block_number, maker, taker, nft_address, token_id, payment_token, price, side, fee)
VALUES ('0x' || repeat('b', 64), 0, 12346, '0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc', '0x90f79bf6eb2c4f870365e785982e1f101e93b906', '0x9fe46736679d2d9a65f0992f2272de9f3c7fa6e0', 11, '0x5fc8d32690cc91d4c39d9d3abcbd16989f875707', 75000000, 1, 750000);

UPDATE indexer_status SET last_processed_block = 12346 WHERE id = 1;

RAISE NOTICE '测试数据插入完成！';
