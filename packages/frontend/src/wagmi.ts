// ============================================
// Wagmi 配置文件 - Web3 钱包连接和链配置
// ============================================
// 企业级最佳实践：使用环境变量而非硬编码

import { http, createConfig } from 'wagmi'
import { sepolia, type Chain } from 'wagmi/chains'
import { injected } from 'wagmi/connectors'

/**
 * 从环境变量获取配置
 */
const RPC_URL = import.meta.env.VITE_RPC_URL || 'http://localhost:8545'

/**
 * 自定义 Anvil 本地链配置
 * 解决 "unsupported chain" 警告
 */
const anvilChain: Chain = {
  id: 31337,
  name: 'Anvil Local',
  nativeCurrency: {
    decimals: 18,
    name: 'Ethereum',
    symbol: 'ETH',
  },
  rpcUrls: {
    default: { http: [RPC_URL] },
    public: { http: [RPC_URL] },
  },
  blockExplorers: undefined,
  contracts: undefined,
  testnet: true,
}

/**
 * Wagmi 配置
 * 支持多个链：Anvil 本地链 + Sepolia 测试网
 * 用户可以在 MetaMask 中切换
 */
export const config = createConfig({
  chains: [anvilChain, sepolia],
  connectors: [
    injected(), // MetaMask 等注入式钱包
  ],
  transports: {
    [anvilChain.id]: http(RPC_URL),
    [sepolia.id]: http(), // 使用公共 RPC
  },
})

/**
 * 合约地址配置（从环境变量读取）
 * 部署后在 .env.local 中配置，无需修改代码
 */
export const CONTRACTS = {
  marketplace: (import.meta.env.VITE_MARKETPLACE_ADDRESS || 
    '0x5FbDB2315678afecb367f032d93F642f64180aa3') as `0x${string}`,
  mockNFT: (import.meta.env.VITE_NFT_ADDRESS || 
    '0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512') as `0x${string}`,
  mockUSDC: (import.meta.env.VITE_USDC_ADDRESS || 
    '0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0') as `0x${string}`,
}

/**
 * API 端点配置（从环境变量读取）
 */
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081'

