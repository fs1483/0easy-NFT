// ============================================
// 前端配置管理 - 企业级最佳实践
// ============================================
// 使用环境变量而非硬编码，支持多环境部署

import { type Address } from 'viem'
import { localhost, sepolia, mainnet } from 'wagmi/chains'

/**
 * 支持的链配置
 */
const SUPPORTED_CHAINS = {
  31337: localhost,   // Anvil 本地测试链
  11155111: sepolia,  // Sepolia 测试网
  1: mainnet,         // 以太坊主网
} as const

/**
 * 从环境变量获取配置
 */
export const config = {
  // 链 ID（从环境变量读取，默认本地链）
  chainId: Number(import.meta.env.VITE_CHAIN_ID || 31337),
  
  // RPC URL（从环境变量读取）
  rpcUrl: import.meta.env.VITE_RPC_URL || 'http://localhost:8545',
  
  // 合约地址（从环境变量读取）
  contracts: {
    marketplace: (import.meta.env.VITE_MARKETPLACE_ADDRESS || 
      '0x5FbDB2315678afecb367f032d93F642f64180aa3') as Address,
    mockNFT: (import.meta.env.VITE_NFT_ADDRESS || 
      '0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512') as Address,
    mockUSDC: (import.meta.env.VITE_USDC_ADDRESS || 
      '0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0') as Address,
  },
  
  // API 端点（从环境变量读取）
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081',
  
  // 调试模式
  debug: import.meta.env.VITE_ENABLE_DEBUG === 'true',
} as const

/**
 * 获取当前链配置
 */
export function getCurrentChain() {
  const chainId = config.chainId as keyof typeof SUPPORTED_CHAINS
  const chain = SUPPORTED_CHAINS[chainId]
  
  if (!chain) {
    throw new Error(`不支持的链 ID: ${config.chainId}`)
  }
  
  return chain
}

/**
 * 验证配置是否有效
 */
export function validateConfig() {
  const errors: string[] = []
  
  // 检查合约地址是否为零地址
  if (config.contracts.marketplace === '0x0000000000000000000000000000000000000000') {
    errors.push('Marketplace 合约地址未配置')
  }
  
  if (config.contracts.mockNFT === '0x0000000000000000000000000000000000000000') {
    errors.push('NFT 合约地址未配置')
  }
  
  if (config.contracts.mockUSDC === '0x0000000000000000000000000000000000000000') {
    errors.push('USDC 合约地址未配置')
  }
  
  if (errors.length > 0) {
    console.warn('⚠️ 配置警告:', errors)
    
    if (config.debug) {
      console.log('📝 当前配置:', {
        chainId: config.chainId,
        rpcUrl: config.rpcUrl,
        contracts: config.contracts,
      })
    }
  }
  
  return errors.length === 0
}

// 开发环境下验证配置
if (import.meta.env.DEV) {
  validateConfig()
}

/**
 * 导出给其他模块使用
 */
export const CONTRACTS = config.contracts
export const API_BASE_URL = config.apiBaseUrl
export const CHAIN_ID = config.chainId
export const RPC_URL = config.rpcUrl

