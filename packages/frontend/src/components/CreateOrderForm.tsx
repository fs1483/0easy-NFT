// 创建订单表单组件
// 支持创建 Ask (卖单) 和 Bid (买单)

import { useState } from 'react'
import { useAccount, useSignTypedData, useReadContract } from 'wagmi'
import { type Address } from 'viem'
import { useOrderAPI } from '../hooks/useOrderAPI'
import { ApprovalModal } from './ApprovalModal'
import {
  createAskOrder,
  createBidOrder,
  getEIP712Domain,
  EIP712_ORDER_TYPES,
  parsePrice,
  OrderSide,
  type Order,
} from '../utils/eip712'
import { CONTRACTS } from '../wagmi'
import MockUSDCABI from '../contracts/MockUSDC.json'
import OeasyNFTABI from '../contracts/OeasyNFT.json'
import './CreateOrderForm.css'

interface CreateOrderFormProps {
  nftAddress?: Address
  tokenId?: string
  defaultSide?: 'ask' | 'bid'
  defaultPrice?: string  // 新增：预填充价格
  onSuccess?: () => void
  onCancel?: () => void
}

export function CreateOrderForm({
  nftAddress: propNftAddress,
  tokenId: propTokenId,
  defaultSide = 'ask',
  defaultPrice,  // 新增
  onSuccess,
  onCancel,
}: CreateOrderFormProps) {
  const { address, chain } = useAccount()
  const { signTypedDataAsync } = useSignTypedData()
  const { createOrder, loading, error } = useOrderAPI()

  const [orderSide, setOrderSide] = useState<'ask' | 'bid'>(defaultSide)
  const [nftAddress, setNftAddress] = useState(propNftAddress || CONTRACTS.mockNFT)
  const [tokenId, setTokenId] = useState(propTokenId || '')
  const [price, setPrice] = useState(defaultPrice || '')
  const [expiryDays, setExpiryDays] = useState('7')
  const [submitting, setSubmitting] = useState(false)
  const [showApprovalModal, setShowApprovalModal] = useState(false)
  const [approvalType, setApprovalType] = useState<'NFT' | 'USDC'>('NFT')
  const [pendingOrder, setPendingOrder] = useState<Order | null>(null)
  
  // 读取 USDC 授权额度
  const { data: usdcAllowance, refetch: refetchAllowance } = useReadContract({
    address: CONTRACTS.mockUSDC,
    abi: MockUSDCABI.abi,
    functionName: 'allowance',
    args: address ? [address, CONTRACTS.marketplace] : undefined,
  })
  
  // 读取 NFT 授权状态
  const { data: nftApproved, refetch: refetchNFTApproval } = useReadContract({
    address: nftAddress,
    abi: OeasyNFTABI.abi,
    functionName: 'isApprovedForAll',
    args: address ? [address, CONTRACTS.marketplace] : undefined,
  })

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()

    if (!address || !chain) {
      alert('请先连接钱包')
      return
    }

    if (!tokenId || !price) {
      alert('请填写完整信息')
      return
    }

    // 1. 构建订单数据
    const priceWei = parsePrice(price, 6)
    const order: Order = orderSide === 'ask'
      ? createAskOrder(address, nftAddress, BigInt(tokenId), CONTRACTS.mockUSDC, priceWei, parseInt(expiryDays))
      : createBidOrder(address, nftAddress, BigInt(tokenId), CONTRACTS.mockUSDC, priceWei, parseInt(expiryDays))

    // 2. 检查授权（企业级标准流程）
    if (orderSide === 'ask') {
      // 卖单：检查 NFT 授权
      if (!nftApproved) {
        setPendingOrder(order)
        setApprovalType('NFT')
        setShowApprovalModal(true)
        return  // 等待授权完成
      }
    } else {
      // 买单：检查 USDC 授权
      const allowance = (usdcAllowance as bigint) || 0n
      if (allowance < priceWei) {
        setPendingOrder(order)
        setApprovalType('USDC')
        setShowApprovalModal(true)
        return  // 等待授权完成
      }
    }

    // 3. 授权充足，直接创建订单
    await submitOrder(order)
  }
  
  async function submitOrder(order: Order) {
    if (!chain) return
    
    setSubmitting(true)

    try {
      console.log('创建订单:', order)

      // EIP-712 签名
      const domain = getEIP712Domain(chain.id, CONTRACTS.marketplace)
      
      const signature = await signTypedDataAsync({
        domain,
        types: EIP712_ORDER_TYPES,
        primaryType: 'Order',
        message: order,
      })

      console.log('签名成功:', signature)

      // 提交到后端
      const result = await createOrder(order, signature)

      if (result) {
        alert(`订单创建成功! ID: ${result.id}`)
        setTokenId('')
        setPrice('')
        onSuccess?.()
      } else {
        alert('订单创建失败，请查看控制台')
      }
    } catch (err) {
      console.error('创建订单错误:', err)
      alert(err instanceof Error ? err.message : '创建订单失败')
    } finally {
      setSubmitting(false)
    }
  }
  
  // 授权完成后的回调
  async function handleApprovalComplete() {
    setShowApprovalModal(false)
    
    // 刷新授权状态
    await refetchAllowance()
    await refetchNFTApproval()
    
    // 继续创建订单
    if (pendingOrder) {
      await submitOrder(pendingOrder)
      setPendingOrder(null)
    }
  }

  return (
    <div className="create-order-form">
      <h2>创建订单</h2>
      
      {/* 授权弹窗 */}
      {showApprovalModal && (
        <ApprovalModal
          type={approvalType}
          nftAddress={nftAddress}
          requiredAmount={parsePrice(price || '0', 6)}
          onApproved={handleApprovalComplete}
          onCancel={() => {
            setShowApprovalModal(false)
            setPendingOrder(null)
          }}
        />
      )}

      <form onSubmit={handleSubmit}>
        {/* 订单类型选择 */}
        <div className="form-group">
          <label>订单类型</label>
          <div className="order-type-selector">
            <button
              type="button"
              className={orderSide === 'ask' ? 'type-btn active ask' : 'type-btn ask'}
              onClick={() => setOrderSide('ask')}
            >
              卖单 (Ask)
            </button>
            <button
              type="button"
              className={orderSide === 'bid' ? 'type-btn active bid' : 'type-btn bid'}
              onClick={() => setOrderSide('bid')}
            >
              买单 (Bid)
            </button>
          </div>
          <p className="hint">
            {orderSide === 'ask' 
              ? '卖单：你希望以什么价格出售这个 NFT' 
              : '买单：你愿意出价多少购买这个 NFT'}
          </p>
        </div>

        {/* NFT 合约地址 */}
        <div className="form-group">
          <label htmlFor="nftAddress">NFT 合约地址</label>
          <input
            id="nftAddress"
            type="text"
            value={nftAddress}
            onChange={(e) => setNftAddress(e.target.value as Address)}
            placeholder="0x..."
            required
            disabled={!!propNftAddress}
          />
        </div>

        {/* Token ID */}
        <div className="form-group">
          <label htmlFor="tokenId">Token ID</label>
          <input
            id="tokenId"
            type="number"
            value={tokenId}
            onChange={(e) => setTokenId(e.target.value)}
            placeholder="1"
            required
            disabled={!!propTokenId}
          />
        </div>

        {/* 价格 */}
        <div className="form-group">
          <label htmlFor="price">价格 (USDC)</label>
          <input
            id="price"
            type="number"
            step="0.000001"
            value={price}
            onChange={(e) => setPrice(e.target.value)}
            placeholder="100.00"
            required
          />
          <p className="hint">输入 USDC 价格（例如：100.5）</p>
        </div>

        {/* 有效期 */}
        <div className="form-group">
          <label htmlFor="expiry">有效期（天）</label>
          <select
            id="expiry"
            value={expiryDays}
            onChange={(e) => setExpiryDays(e.target.value)}
          >
            <option value="1">1 天</option>
            <option value="3">3 天</option>
            <option value="7">7 天</option>
            <option value="14">14 天</option>
            <option value="30">30 天</option>
          </select>
        </div>

        {/* 错误提示 */}
        {error && (
          <div className="error-message">
            {error}
          </div>
        )}

        {/* 提交按钮 */}
        <div className="form-actions">
          {onCancel && (
            <button type="button" onClick={onCancel} className="btn-cancel">
              取消
            </button>
          )}
          <button
            type="submit"
            className="btn-submit"
            disabled={loading || submitting || !address}
          >
            {submitting ? '签名中...' : loading ? '提交中...' : '创建订单'}
          </button>
        </div>

        {/* 提示信息 */}
        <div className="info-box">
          <h4>📝 注意事项</h4>
          <ul>
            <li>创建订单需要签名，但<strong>不消耗 Gas</strong></li>
            <li>订单会在链下存储，等待撮合</li>
            <li>
              {orderSide === 'ask'
                ? '请确保你拥有该 NFT，并已授权给市场合约'
                : '请确保你有足够的 USDC，并已授权给市场合约'}
            </li>
            <li>订单一旦撮合成功，会自动在链上执行</li>
          </ul>
        </div>
      </form>
    </div>
  )
}

