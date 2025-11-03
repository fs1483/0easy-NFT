// 智能合约常量配置
// 集中管理合约地址和 ABI

import { type Address } from 'viem'
import { SupportedChainId, type ContractAddresses } from '../types/contract'

/**
 * 各网络的合约地址配置
 */
export const CONTRACT_ADDRESSES: Record<SupportedChainId, ContractAddresses> = {
  // Anvil 本地测试网
  [SupportedChainId.LOCALHOST]: {
    marketplace: '0x5FbDB2315678afecb367f032d93F642f64180aa3' as Address,
    mockNFT: '0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512' as Address,
    mockUSDC: '0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0' as Address,
  },
  
  // Sepolia 测试网（部署后更新）
  [SupportedChainId.SEPOLIA]: {
    marketplace: '0x0000000000000000000000000000000000000000' as Address,
    mockNFT: '0x0000000000000000000000000000000000000000' as Address,
    mockUSDC: '0x0000000000000000000000000000000000000000' as Address,
  },
  
  // 以太坊主网（暂不支持）
  [SupportedChainId.MAINNET]: {
    marketplace: '0x0000000000000000000000000000000000000000' as Address,
    mockNFT: '0x0000000000000000000000000000000000000000' as Address,
    mockUSDC: '0x0000000000000000000000000000000000000000' as Address,
  },
}

/**
 * 获取当前网络的合约地址
 */
export function getContractAddresses(chainId: number): ContractAddresses {
  const addresses = CONTRACT_ADDRESSES[chainId as SupportedChainId]
  
  if (!addresses) {
    throw new Error(`不支持的网络 Chain ID: ${chainId}`)
  }
  
  return addresses
}

/**
 * EIP-712 Domain 名称
 */
export const EIP712_DOMAIN_NAME = 'Oeasy Marketplace'

/**
 * EIP-712 Domain 版本
 */
export const EIP712_DOMAIN_VERSION = '1'

/**
 * 代币精度配置
 */
export const TOKEN_DECIMALS = {
  USDC: 6,  // USDC 使用 6 位精度
  ETH: 18,  // ETH 使用 18 位精度
} as const

/**
 * 默认订单有效期（天）
 */
export const DEFAULT_ORDER_EXPIRY_DAYS = 7

/**
 * 最大订单有效期（天）
 */
export const MAX_ORDER_EXPIRY_DAYS = 90

/**
 * 订单刷新间隔（毫秒）
 */
export const ORDER_REFRESH_INTERVAL = 10000 // 10 秒

