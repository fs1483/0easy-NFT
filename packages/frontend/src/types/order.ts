// 订单相关类型定义
// 集中管理所有订单相关的 TypeScript 类型

import { type Address, type Hex } from 'viem'

/**
 * 订单方向枚举
 */
export enum OrderSide {
  ASK = 0, // 卖单
  BID = 1, // 买单
}

/**
 * 订单状态枚举
 */
export enum OrderStatus {
  ACTIVE = 'active',       // 活跃
  FILLED = 'filled',       // 已成交
  CANCELLED = 'cancelled', // 已取消
}

/**
 * 订单数据结构（符合智能合约定义）
 * 用于 EIP-712 签名
 */
export interface Order {
  maker: Address        // 订单创建者地址
  nft: Address          // NFT 合约地址
  tokenId: bigint       // NFT Token ID
  paymentToken: Address // 支付代币地址（如 USDC）
  price: bigint         // 价格（wei 单位）
  expiry: bigint        // 过期时间戳（秒）
  nonce: bigint         // 唯一 nonce（防重放）
  side: OrderSide       // 订单方向（0=Ask, 1=Bid）
}

/**
 * 后端 API 返回的订单数据结构
 */
export interface OrderResponse {
  id: number
  maker: Address
  nftAddress: Address
  tokenId: string
  paymentToken: Address
  price: string
  expiry: string // ISO 时间字符串
  nonce: string
  side: 'ask' | 'bid'
  status: OrderStatus
  signature: Hex
  hash: Hex
  createdAt: string
  updatedAt: string
}

/**
 * 创建订单请求参数
 */
export interface CreateOrderRequest {
  maker: Address
  nft: Address
  tokenId: string
  paymentToken: Address
  price: string
  expiry: number // Unix 时间戳（秒）
  nonce: string
  side: OrderSide
  signature: Hex
}

/**
 * 订单查询过滤参数
 */
export interface OrderFilters {
  side?: 'ask' | 'bid'
  collection?: Address
  status?: OrderStatus
  maker?: Address
}

/**
 * 取消订单请求参数
 */
export interface CancelOrderRequest {
  signature: Hex
}

/**
 * 订单统计数据
 */
export interface OrderStats {
  totalOrders: number
  activeOrders: number
  filledOrders: number
  cancelledOrders: number
  totalVolume: string
  avgPrice: string
}

