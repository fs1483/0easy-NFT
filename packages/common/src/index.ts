// 导出合约 ABI
export { default as OeasyMarketplaceABI } from './contracts/OeasyMarketplace.json'
export { default as OeasyNFTABI } from './contracts/OeasyNFT.json'
export { default as MockUSDCABI } from './contracts/MockUSDC.json'

// 合约地址类型
export interface ContractAddresses {
  marketplace: `0x${string}`
  mockNFT: `0x${string}`
  mockUSDC: `0x${string}`
}

// 订单类型定义（与后端保持一致）
export interface Order {
  id: number
  maker: string
  nftAddress: string
  tokenId: string
  paymentToken: string
  price: string
  expiry: string
  nonce: string
  side: 'ask' | 'bid'
  status: 'active' | 'cancelled' | 'filled'
  signature: string
  hash: string
  createdAt: string
  updatedAt: string
}

// EIP-712 订单结构（用于签名）
export interface OrderStruct {
  maker: `0x${string}`
  nft: `0x${string}`
  tokenId: bigint
  paymentToken: `0x${string}`
  price: bigint
  expiry: bigint
  nonce: bigint
  side: 0 | 1 // 0 = Ask, 1 = Bid
}

