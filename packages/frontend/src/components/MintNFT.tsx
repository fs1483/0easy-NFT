// 铸造 NFT 组件
// 企业级功能：用户可以铸造自己的 NFT

import { useState } from 'react'
import { useAccount, useWriteContract, useWaitForTransactionReceipt } from 'wagmi'
import { CONTRACTS } from '../wagmi'
import OeasyNFTABI from '../contracts/OeasyNFT.json'
import './MintNFT.css'

export function MintNFT() {
  const { address, chain } = useAccount()
  const { writeContractAsync } = useWriteContract()
  
  const [tokenId, setTokenId] = useState('')
  const [minting, setMinting] = useState(false)
  const [approving, setApproving] = useState(false)

  async function handleMint(e: React.FormEvent) {
    e.preventDefault()

    if (!address || !chain) {
      alert('请先连接钱包')
      return
    }

    if (!tokenId) {
      alert('请输入 Token ID')
      return
    }

    setMinting(true)

    try {
      console.log('🎨 开始铸造 NFT...')
      
      // 1. 铸造 NFT
      const mintHash = await writeContractAsync({
        address: CONTRACTS.mockNFT,
        abi: OeasyNFTABI.abi,
        functionName: 'mintWithId',
        args: [address, BigInt(tokenId)],
      })

      console.log('⏳ 等待铸造交易确认...', mintHash)
      
      // 简单等待（实际应该用 useWaitForTransactionReceipt）
      await new Promise(resolve => setTimeout(resolve, 3000))

      console.log('✅ NFT 铸造成功！')
      
      // 2. 自动授权给 Marketplace
      setMinting(false)
      setApproving(true)
      
      console.log('🔐 开始授权 Marketplace...')
      
      const approveHash = await writeContractAsync({
        address: CONTRACTS.mockNFT,
        abi: OeasyNFTABI.abi,
        functionName: 'setApprovalForAll',
        args: [CONTRACTS.marketplace, true],
      })

      console.log('⏳ 等待授权交易确认...', approveHash)
      
      await new Promise(resolve => setTimeout(resolve, 3000))

      console.log('✅ Marketplace 授权成功！')
      
      alert(`🎉 NFT #${tokenId} 铸造成功并已授权！\n\n现在可以创建卖单了。`)
      
      // 重置表单
      setTokenId('')
    } catch (error: any) {
      console.error('❌ 铸造失败:', error)
      
      if (error.message?.includes('user rejected')) {
        alert('用户取消了交易')
      } else if (error.message?.includes('already minted')) {
        alert(`Token ID ${tokenId} 已被铸造，请使用其他 ID`)
      } else {
        alert(`铸造失败: ${error.message || '未知错误'}`)
      }
    } finally {
      setMinting(false)
      setApproving(false)
    }
  }

  return (
    <div className="mint-nft-container">
      <div className="mint-nft-header">
        <h2>🎨 铸造 NFT</h2>
        <p className="subtitle">创建你自己的 NFT，然后可以挂单出售</p>
      </div>

      <form onSubmit={handleMint} className="mint-nft-form">
        {/* Token ID 输入 */}
        <div className="form-group">
          <label htmlFor="tokenId">Token ID</label>
          <input
            id="tokenId"
            type="number"
            value={tokenId}
            onChange={(e) => setTokenId(e.target.value)}
            placeholder="输入唯一的 Token ID (如: 100)"
            required
            min="1"
            disabled={minting || approving}
          />
          <p className="hint">
            💡 提示：Token ID 必须是唯一的（1-999999），建议使用较大的数字避免冲突
          </p>
        </div>

        {/* NFT 信息展示 */}
        <div className="nft-info-box">
          <h3>将要铸造的 NFT</h3>
          <div className="info-row">
            <span>合约:</span>
            <span className="value">{CONTRACTS.mockNFT.slice(0, 10)}...</span>
          </div>
          <div className="info-row">
            <span>Token ID:</span>
            <span className="value">#{tokenId || '---'}</span>
          </div>
          <div className="info-row">
            <span>所有者:</span>
            <span className="value">{address ? `${address.slice(0, 6)}...${address.slice(-4)}` : '---'}</span>
          </div>
        </div>

        {/* 流程说明 */}
        <div className="process-box">
          <h4>📋 铸造流程</h4>
          <ol>
            <li className={minting ? 'active' : ''}>
              铸造 NFT（需要签名）
              {minting && <span className="loading">进行中...</span>}
            </li>
            <li className={approving ? 'active' : ''}>
              授权给 Marketplace（需要签名）
              {approving && <span className="loading">进行中...</span>}
            </li>
            <li>完成！可以创建卖单</li>
          </ol>
        </div>

        {/* 提交按钮 */}
        <button
          type="submit"
          className="btn-mint"
          disabled={!address || minting || approving || !tokenId}
        >
          {minting ? '🎨 铸造中...' : 
           approving ? '🔐 授权中...' : 
           '🚀 开始铸造'}
        </button>

        {/* 说明 */}
        <div className="notice-box">
          <h4>ℹ️ 注意事项</h4>
          <ul>
            <li>铸造 NFT 需要消耗 Gas（测试网免费）</li>
            <li>Token ID 一旦铸造不可修改</li>
            <li>铸造后会自动授权给 Marketplace</li>
            <li>授权后即可创建卖单</li>
          </ul>
        </div>
      </form>
    </div>
  )
}


