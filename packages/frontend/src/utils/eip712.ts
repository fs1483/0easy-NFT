// EIP-712 签名工具函数
// 用于创建和签名订单数据，实现无 Gas 链下挂单

import { type Address, type TypedDataDomain } from 'viem'

/**
 * 订单方向枚举
 */
export enum OrderSide {
  ASK = 0, // 卖单
  BID = 1, // 买单
}

/**
 * 订单数据结构（符合智能合约定义）
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
 * EIP-712 Domain 定义
 * 必须与智能合约中的定义完全一致
 */
export function getEIP712Domain(
  chainId: number,
  verifyingContract: Address
): TypedDataDomain {
  return {
    name: 'Oeasy Marketplace',
    version: '1',
    chainId,
    verifyingContract,
  }
}

/**
 * EIP-712 类型定义
 * 必须与智能合约中的 Order 结构体完全一致
 */
export const EIP712_ORDER_TYPES = {
  Order: [
    { name: 'maker', type: 'address' },
    { name: 'nft', type: 'address' },
    { name: 'tokenId', type: 'uint256' },
    { name: 'paymentToken', type: 'address' },
    { name: 'price', type: 'uint256' },
    { name: 'expiry', type: 'uint256' },
    { name: 'nonce', type: 'uint256' },
    { name: 'side', type: 'uint8' },
  ],
} as const

/**
 * 生成随机 nonce
 * 使用时间戳 + 随机数确保唯一性
 */
export function generateNonce(): bigint {
  const timestamp = BigInt(Date.now())
  const random = BigInt(Math.floor(Math.random() * 1000000))
  return timestamp * 1000000n + random
}

/**
 * 创建 Ask 订单（卖单）
 * 
 * @param maker - 卖家地址
 * @param nft - NFT 合约地址
 * @param tokenId - NFT Token ID
 * @param paymentToken - 支付代币地址
 * @param price - 价格（wei）
 * @param expiryDays - 有效期（天数）
 */
export function createAskOrder(
  maker: Address,
  nft: Address,
  tokenId: bigint,
  paymentToken: Address,
  price: bigint,
  expiryDays: number = 7
): Order {
  const now = Math.floor(Date.now() / 1000)
  const expiry = BigInt(now + expiryDays * 24 * 60 * 60)

  return {
    maker,
    nft,
    tokenId,
    paymentToken,
    price,
    expiry,
    nonce: generateNonce(),
    side: OrderSide.ASK,
  }
}

/**
 * 创建 Bid 订单（买单）
 * 
 * @param maker - 买家地址
 * @param nft - NFT 合约地址
 * @param tokenId - NFT Token ID
 * @param paymentToken - 支付代币地址
 * @param price - 出价（wei）
 * @param expiryDays - 有效期（天数）
 */
export function createBidOrder(
  maker: Address,
  nft: Address,
  tokenId: bigint,
  paymentToken: Address,
  price: bigint,
  expiryDays: number = 7
): Order {
  const now = Math.floor(Date.now() / 1000)
  const expiry = BigInt(now + expiryDays * 24 * 60 * 60)

  return {
    maker,
    nft,
    tokenId,
    paymentToken,
    price,
    expiry,
    nonce: generateNonce(),
    side: OrderSide.BID,
  }
}

/**
 * 格式化价格显示
 * 
 * @param price - 价格（wei）
 * @param decimals - 代币精度（默认 6 for USDC）
 */
export function formatPrice(price: bigint, decimals: number = 6): string {
  const divisor = BigInt(10 ** decimals)
  const integerPart = price / divisor
  const decimalPart = price % divisor
  
  if (decimalPart === 0n) {
    return integerPart.toString()
  }
  
  const decimalStr = decimalPart.toString().padStart(decimals, '0')
  return `${integerPart}.${decimalStr.replace(/0+$/, '')}`
}

/**
 * 解析价格输入
 * 
 * @param priceStr - 价格字符串（如 "100.5"）
 * @param decimals - 代币精度（默认 6 for USDC）
 */
export function parsePrice(priceStr: string, decimals: number = 6): bigint {
  const parts = priceStr.split('.')
  const integerPart = parts[0] || '0'
  const decimalPart = (parts[1] || '').padEnd(decimals, '0').slice(0, decimals)
  
  return BigInt(integerPart) * BigInt(10 ** decimals) + BigInt(decimalPart)
}

/**
 * 格式化过期时间显示
 * @param expiry - bigint (Unix 时间戳秒) 或 string (ISO 日期字符串)
 */
export function formatExpiry(expiry: bigint | string): string {
  let expiryDate: Date
  
  if (typeof expiry === 'string') {
    // ISO 日期字符串 "2025-10-18T06:41:12.872278Z"
    expiryDate = new Date(expiry)
  } else {
    // Unix 时间戳（秒）
    expiryDate = new Date(Number(expiry) * 1000)
  }
  
  const now = new Date()
  
  if (expiryDate < now) {
    return '已过期'
  }
  
  const diffMs = expiryDate.getTime() - now.getTime()
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))
  const diffHours = Math.floor((diffMs % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
  
  if (diffDays > 0) {
    return `${diffDays} 天 ${diffHours} 小时后过期`
  } else if (diffHours > 0) {
    return `${diffHours} 小时后过期`
  } else {
    const diffMinutes = Math.floor((diffMs % (1000 * 60 * 60)) / (1000 * 60))
    return `${diffMinutes} 分钟后过期`
  }
}

