// 授权弹窗组件
// 企业级标准：创建订单前检查并请求授权

import { useState } from 'react'
import { useWriteContract } from 'wagmi'
import { formatUnits } from 'viem'
import { CONTRACTS } from '../wagmi'
import MockUSDCABI from '../contracts/MockUSDC.json'
import OeasyNFTABI from '../contracts/OeasyNFT.json'
import './ApprovalModal.css'

interface ApprovalModalProps {
  type: 'NFT' | 'USDC'
  nftAddress?: string
  requiredAmount?: bigint
  onApproved: () => void
  onCancel: () => void
}

export function ApprovalModal({ 
  type, 
  nftAddress,
  requiredAmount, 
  onApproved, 
  onCancel 
}: ApprovalModalProps) {
  const { writeContractAsync } = useWriteContract()
  const [approving, setApproving] = useState(false)
  const [approvalType, setApprovalType] = useState<'unlimited' | 'exact'>('unlimited')

  async function handleApprove() {
    setApproving(true)

    try {
      if (type === 'NFT') {
        // 授权 NFT（只能是全部授权）
        await writeContractAsync({
          address: (nftAddress || CONTRACTS.mockNFT) as `0x${string}`,
          abi: OeasyNFTABI.abi,
          functionName: 'setApprovalForAll',
          args: [CONTRACTS.marketplace, true],
        })
      } else {
        // 授权 USDC
        const amount = approvalType === 'unlimited'
          ? 2n ** 256n - 1n  // 无限授权
          : requiredAmount || 0n  // 精确授权

        await writeContractAsync({
          address: CONTRACTS.mockUSDC,
          abi: MockUSDCABI.abi,
          functionName: 'approve',
          args: [CONTRACTS.marketplace, amount],
        })
      }

      // 等待确认
      await new Promise(resolve => setTimeout(resolve, 2000))

      onApproved()
    } catch (error: any) {
      console.error('授权失败:', error)
      if (!error.message?.includes('user rejected')) {
        alert(`授权失败: ${error.message || '未知错误'}`)
      }
    } finally {
      setApproving(false)
    }
  }

  const requiredAmountFormatted = requiredAmount 
    ? formatUnits(requiredAmount, 6)
    : '0'

  return (
    <div className="modal-overlay" onClick={onCancel}>
      <div className="approval-modal" onClick={e => e.stopPropagation()}>
        <div className="modal-header">
          <h3>🔐 需要授权</h3>
        </div>

        <div className="modal-body">
          {type === 'NFT' ? (
            <>
              <p className="description">
                要创建卖单，需要先授权 Marketplace 合约操作你的 NFT。
              </p>
              <div className="info-box">
                <div className="info-row">
                  <span>授权对象:</span>
                  <span className="value">Oeasy Marketplace</span>
                </div>
                <div className="info-row">
                  <span>授权范围:</span>
                  <span className="value">所有 NFT</span>
                </div>
              </div>
              <div className="notice">
                ℹ️ 授权后，Marketplace 可以转移你的 NFT（仅在成交时）
              </div>
            </>
          ) : (
            <>
              <p className="description">
                要创建买单，需要先授权 Marketplace 合约使用你的 USDC 进行支付。
              </p>
              
              <div className="approval-options">
                <label className="radio-option">
                  <input
                    type="radio"
                    value="exact"
                    checked={approvalType === 'exact'}
                    onChange={() => setApprovalType('exact')}
                  />
                  <div className="option-content">
                    <div className="option-title">仅本次交易</div>
                    <div className="option-desc">
                      授权 {requiredAmountFormatted} USDC
                    </div>
                    <div className="option-note">
                      ✅ 最安全 ⚠️ 每次都要授权
                    </div>
                  </div>
                </label>

                <label className="radio-option recommended">
                  <input
                    type="radio"
                    value="unlimited"
                    checked={approvalType === 'unlimited'}
                    onChange={() => setApprovalType('unlimited')}
                  />
                  <div className="option-content">
                    <div className="option-title">
                      无限授权 <span className="badge">推荐</span>
                    </div>
                    <div className="option-desc">
                      授权无限额度（一次授权永久有效）
                    </div>
                    <div className="option-note">
                      ✅ 方便 ✅ 节省 Gas
                    </div>
                  </div>
                </label>
              </div>

              <div className="info-box">
                <div className="info-row">
                  <span>代币:</span>
                  <span className="value">USDC</span>
                </div>
                <div className="info-row">
                  <span>授权对象:</span>
                  <span className="value">Oeasy Marketplace</span>
                </div>
                <div className="info-row">
                  <span>当前交易需要:</span>
                  <span className="value">{requiredAmountFormatted} USDC</span>
                </div>
              </div>

              <div className="notice">
                ℹ️ 授权不会转移你的资金，只是允许 Marketplace 在交易成交时使用
              </div>
            </>
          )}
        </div>

        <div className="modal-footer">
          <button 
            className="btn-cancel" 
            onClick={onCancel}
            disabled={approving}
          >
            取消
          </button>
          <button 
            className="btn-approve" 
            onClick={handleApprove}
            disabled={approving}
          >
            {approving ? '授权中...' : '授权'}
          </button>
        </div>
      </div>
    </div>
  )
}

