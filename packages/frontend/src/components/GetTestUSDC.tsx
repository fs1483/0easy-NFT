// 获取测试 USDC 组件
// 企业级功能：用户可以铸造测试代币

import { useState } from 'react'
import { useAccount, useWriteContract, useReadContract } from 'wagmi'
import { CONTRACTS } from '../wagmi'
import { formatUnits } from 'viem'
import MockUSDCABI from '../contracts/MockUSDC.json'

export function GetTestUSDC() {
  const { address } = useAccount()
  const { writeContractAsync } = useWriteContract()
  const [minting, setMinting] = useState(false)

  // 读取当前 USDC 余额
  const { data: balance, refetch: refetchBalance } = useReadContract({
    address: CONTRACTS.mockUSDC,
    abi: MockUSDCABI.abi,
    functionName: 'balanceOf',
    args: address ? [address] : undefined,
  })

  // 读取授权额度
  const { data: allowance, refetch: refetchAllowance } = useReadContract({
    address: CONTRACTS.mockUSDC,
    abi: MockUSDCABI.abi,
    functionName: 'allowance',
    args: address ? [address, CONTRACTS.marketplace] : undefined,
  })

  async function handleGetUSDC() {
    if (!address) {
      alert('请先连接钱包')
      return
    }

    setMinting(true)

    try {
      console.log('💰 开始铸造 USDC...')
      
      // 1. 铸造 1000 USDC
      const amount = 1000n * 1000000n  // 1000 USDC (6 decimals)
      
      const mintHash = await writeContractAsync({
        address: CONTRACTS.mockUSDC,
        abi: MockUSDCABI.abi,
        functionName: 'mint',
        args: [address, amount],
      })

      console.log('⏳ 等待铸造交易确认...', mintHash)
      
      await new Promise(resolve => setTimeout(resolve, 3000))
      await refetchBalance()

      console.log('✅ USDC 铸造成功！')
      
      alert('🎉 获得 1000 测试 USDC！\n\n💡 提示：创建买单时系统会自动检查并提示授权。')
    } catch (error: any) {
      console.error('❌ 获取 USDC 失败:', error)
      
      if (error.message?.includes('user rejected')) {
        alert('用户取消了交易')
      } else {
        alert(`获取 USDC 失败: ${error.message || '未知错误'}`)
      }
    } finally {
      setMinting(false)
    }
  }

  const balanceFormatted = balance ? formatUnits(balance as bigint, 6) : '0'
  const allowanceFormatted = allowance ? formatUnits(allowance as bigint, 6) : '0'

  return (
    <div className="get-usdc-container">
      <div className="get-usdc-header">
        <h2>💰 获取测试 USDC</h2>
        <p className="subtitle">铸造测试代币用于购买 NFT（仅测试网）</p>
        <div className="testnet-badge">
          🧪 测试网专用功能
        </div>
      </div>

      <div className="usdc-info">
        {/* 余额显示 */}
        <div className="balance-card">
          <div className="balance-label">当前余额</div>
          <div className="balance-value">{balanceFormatted} USDC</div>
        </div>

        {/* 授权额度 */}
        <div className="balance-card">
          <div className="balance-label">Marketplace 授权额度</div>
          <div className="balance-value">
            {allowanceFormatted === '0' ? '未授权' : 
             Number(allowanceFormatted) > 1000000 ? '已授权 (∞)' :
             `${allowanceFormatted} USDC`}
          </div>
        </div>
      </div>

      {/* 获取按钮 */}
      <button
        onClick={handleGetUSDC}
        className="btn-get-usdc"
        disabled={!address || minting}
      >
        {minting ? '💰 铸造中...' : '🎁 免费获取 1000 USDC'}
      </button>

      {/* 说明 */}
      <div className="info-box">
        <h4>📝 使用说明</h4>
        <ul>
          <li>点击按钮后会铸造 1000 测试 USDC 到你的钱包</li>
          <li>💡 <strong>授权在创建买单时自动处理</strong>（企业级标准）</li>
          <li>每次点击都会增加 1000 USDC</li>
          <li>⚠️ 仅限测试网使用，无真实价值</li>
          <li>🌐 主网环境：用户需要从交易所购买真实 USDC</li>
        </ul>
      </div>

      <style jsx>{`
        .get-usdc-container {
          max-width: 600px;
          margin: 0 auto;
          padding: 20px;
        }

        .get-usdc-header {
          text-align: center;
          margin-bottom: 30px;
        }

        .get-usdc-header h2 {
          margin: 0 0 8px 0;
          font-size: 32px;
          font-weight: 700;
        }

        .subtitle {
          margin: 0 0 12px 0;
          font-size: 16px;
          color: #6b7280;
        }
        
        .testnet-badge {
          display: inline-block;
          padding: 6px 12px;
          background: #fef3c7;
          color: #92400e;
          border-radius: 6px;
          font-size: 13px;
          font-weight: 600;
        }

        .usdc-info {
          display: grid;
          grid-template-columns: 1fr 1fr;
          gap: 16px;
          margin-bottom: 24px;
        }

        .balance-card {
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          padding: 24px;
          border-radius: 16px;
          color: white;
        }

        .balance-label {
          font-size: 14px;
          opacity: 0.9;
          margin-bottom: 8px;
        }

        .balance-value {
          font-size: 28px;
          font-weight: 800;
        }

        .btn-get-usdc {
          width: 100%;
          padding: 20px;
          background: linear-gradient(135deg, #10b981 0%, #059669 100%);
          color: white;
          border: none;
          border-radius: 12px;
          font-size: 20px;
          font-weight: 700;
          cursor: pointer;
          transition: all 0.3s;
          margin-bottom: 24px;
        }

        .btn-get-usdc:hover:not(:disabled) {
          transform: translateY(-2px);
          box-shadow: 0 8px 24px rgba(16, 185, 129, 0.4);
        }

        .btn-get-usdc:disabled {
          background: #9ca3af;
          cursor: not-allowed;
          transform: none;
        }

        .info-box {
          background: #eff6ff;
          padding: 20px;
          border-left: 4px solid #3b82f6;
          border-radius: 8px;
        }

        .info-box h4 {
          margin: 0 0 12px 0;
          font-size: 16px;
          color: #1e40af;
        }

        .info-box ul {
          margin: 0;
          padding-left: 20px;
        }

        .info-box li {
          margin-bottom: 8px;
          font-size: 14px;
          color: #1e40af;
          line-height: 1.6;
        }

        @media (max-width: 768px) {
          .usdc-info {
            grid-template-columns: 1fr;
          }
        }
      `}</style>
    </div>
  )
}


