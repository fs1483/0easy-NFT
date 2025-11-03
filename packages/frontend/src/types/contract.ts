// 智能合约相关类型定义

import { type Address } from 'viem'

/**
 * 合约地址配置
 */
export interface ContractAddresses {
  marketplace: Address
  mockNFT: Address
  mockUSDC: Address
}

/**
 * 支持的区块链网络
 */
export enum SupportedChainId {
  LOCALHOST = 31337,   // Anvil 本地测试网
  SEPOLIA = 11155111,  // Sepolia 测试网
  MAINNET = 1,         // 以太坊主网
}

/**
 * TradeExecuted 事件数据
 */
export interface TradeExecutedEvent {
  maker: Address
  taker: Address
  nft: Address
  tokenId: bigint
  paymentToken: Address
  price: bigint
  side: number
  fee: bigint
  transactionHash: string
  blockNumber: bigint
}

/**
 * OrderCancelled 事件数据
 */
export interface OrderCancelledEvent {
  maker: Address
  nonce: bigint
  transactionHash: string
  blockNumber: bigint
}

/**
 * 交易状态
 */
export enum TransactionStatus {
  PENDING = 'pending',     // 待确认
  SUCCESS = 'success',     // 成功
  FAILED = 'failed',       // 失败
  CANCELLED = 'cancelled', // 已取消
}

/**
 * NFT 元数据
 */
export interface NFTMetadata {
  name: string
  description: string
  image: string
  attributes?: Array<{
    trait_type: string
    value: string | number
  }>
}

/**
 * NFT 数据
 */
export interface NFT {
  contractAddress: Address
  tokenId: string
  owner: Address
  metadata?: NFTMetadata
  approved?: Address
}

