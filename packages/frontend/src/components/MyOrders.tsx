// 我的订单组件
// 显示当前用户创建的所有订单

import { useEffect, useState } from 'react'
import { useAccount, useSignTypedData } from 'wagmi'
import { useOrderAPI, type OrderResponse } from '../hooks/useOrderAPI'
import { formatPrice, getEIP712Domain, EIP712_ORDER_TYPES } from '../utils/eip712'
import { CONTRACTS } from '../wagmi'
import './MyOrders.css'

export function MyOrders() {
  const { address, chain } = useAccount()
  const { signTypedDataAsync } = useSignTypedData()
  const { loading, error, fetchOrders, cancelOrder } = useOrderAPI()
  const [orders, setOrders] = useState<OrderResponse[]>([])
  const [statusFilter, setStatusFilter] = useState<'all' | 'active' | 'filled' | 'cancelled'>('all')
  const [cancelling, setCancelling] = useState<number | null>(null)

  useEffect(() => {
    if (address) {
      loadMyOrders()
      // 每 10 秒刷新
      const interval = setInterval(loadMyOrders, 10000)
      return () => clearInterval(interval)
    }
  }, [address, statusFilter])

  async function loadMyOrders() {
    if (!address) return

    try {
      // 获取所有状态的订单（active, filled, cancelled）
      const activeOrders = await fetchOrders({ status: 'active' })
      const filledOrders = await fetchOrders({ status: 'filled' })
      const cancelledOrders = await fetchOrders({ status: 'cancelled' })
      
      const allOrders = [...activeOrders, ...filledOrders, ...cancelledOrders]
      
      // 筛选当前用户的订单
      const myOrders = allOrders.filter(order => 
        order.maker.toLowerCase() === address.toLowerCase()
      )
      
      // 按状态筛选
      const filtered = statusFilter === 'all' 
        ? myOrders 
        : myOrders.filter(order => order.status === statusFilter)
      
      setOrders(filtered)
    } catch (err) {
      console.error('加载我的订单失败:', err)
    }
  }
  
  // 处理取消订单
  async function handleCancelOrder(order: OrderResponse) {
    if (!address || !chain) {
      alert('请先连接钱包')
      return
    }
    
    if (!confirm(`确定要取消订单 #${order.id} 吗？`)) {
      return
    }
    
    setCancelling(order.id)
    
    try {
      // 1. 构建取消消息（简化版，使用 maker + nonce）
      const cancelMessage = {
        maker: address.toLowerCase(),
        nonce: BigInt(order.nonce),
      }
      
      // 2. 签名
      const domain = getEIP712Domain(chain.id, CONTRACTS.marketplace)
      
      const cancelSignature = await signTypedDataAsync({
        domain,
        types: {
          Cancel: [
            { name: 'maker', type: 'address' },
            { name: 'nonce', type: 'uint256' },
          ],
        },
        primaryType: 'Cancel',
        message: cancelMessage,
      })
      
      console.log('🔒 取消订单签名:', cancelSignature)
      
      // 3. 调用 API
      const success = await cancelOrder(
        order.id,
        address,
        order.nonce,
        cancelSignature
      )
      
      if (success) {
        alert(`订单 #${order.id} 已取消`)
        loadMyOrders()  // 刷新列表
      } else {
        alert('取消订单失败')
      }
    } catch (err) {
      console.error('取消订单错误:', err)
      alert(err instanceof Error ? err.message : '取消失败')
    } finally {
      setCancelling(null)
    }
  }

  // 统计数据
  const stats = {
    total: orders.length,
    active: orders.filter(o => o.status === 'active').length,
    filled: orders.filter(o => o.status === 'filled').length,
    cancelled: orders.filter(o => o.status === 'cancelled').length,
  }

  return (
    <div className="my-orders-container">
      <div className="my-orders-header">
        <div>
          <h2>我的订单</h2>
          <p className="address-info">
            地址: {address?.slice(0, 6)}...{address?.slice(-4)}
          </p>
        </div>
        
        {/* 统计卡片 */}
        <div className="stats-cards">
          <div className="stat-card">
            <div className="stat-value">{stats.total}</div>
            <div className="stat-label">总订单</div>
          </div>
          <div className="stat-card">
            <div className="stat-value active">{stats.active}</div>
            <div className="stat-label">活跃</div>
          </div>
          <div className="stat-card">
            <div className="stat-value filled">{stats.filled}</div>
            <div className="stat-label">已成交</div>
          </div>
          <div className="stat-card">
            <div className="stat-value cancelled">{stats.cancelled}</div>
            <div className="stat-label">已取消</div>
          </div>
        </div>
      </div>

      {/* 状态筛选 */}
      <div className="filter-tabs">
        <button
          className={statusFilter === 'all' ? 'filter-tab active' : 'filter-tab'}
          onClick={() => setStatusFilter('all')}
        >
          全部 ({stats.total})
        </button>
        <button
          className={statusFilter === 'active' ? 'filter-tab active' : 'filter-tab'}
          onClick={() => setStatusFilter('active')}
        >
          活跃 ({stats.active})
        </button>
        <button
          className={statusFilter === 'filled' ? 'filter-tab active' : 'filter-tab'}
          onClick={() => setStatusFilter('filled')}
        >
          已成交 ({stats.filled})
        </button>
        <button
          className={statusFilter === 'cancelled' ? 'filter-tab active' : 'filter-tab'}
          onClick={() => setStatusFilter('cancelled')}
        >
          已取消 ({stats.cancelled})
        </button>
      </div>

      {loading && <div className="loading">加载中...</div>}
      {error && <div className="error">错误: {error}</div>}

      {!loading && orders.length === 0 && (
        <div className="empty-state">
          <p>暂无订单</p>
          <p className="hint">
            {statusFilter === 'all' 
              ? '创建你的第一个订单！' 
              : `暂无${statusFilter === 'active' ? '活跃' : statusFilter === 'filled' ? '已成交' : '已取消'}的订单`}
          </p>
        </div>
      )}

      {!loading && orders.length > 0 && (
        <div className="orders-list">
          {orders.map(order => (
            <div key={order.id} className={`order-card status-${order.status}`}>
              <div className="order-header">
                <div className="order-type-badge">
                  <span className={order.side === 'ask' ? 'badge-ask' : 'badge-bid'}>
                    {order.side === 'ask' ? '卖单' : '买单'}
                  </span>
                  <span className={`status-badge status-${order.status}`}>
                    {order.status === 'active' ? '活跃' : 
                     order.status === 'filled' ? '已成交' : '已取消'}
                  </span>
                </div>
                <span className="order-id">#{order.id}</span>
              </div>

              <div className="order-body">
                <div className="info-row">
                  <span className="label">NFT:</span>
                  <span className="value">OeasyNFT #{order.tokenId}</span>
                </div>

                <div className="info-row">
                  <span className="label">合约:</span>
                  <span className="value truncate" title={order.nftAddress}>
                    {order.nftAddress.slice(0, 10)}...
                  </span>
                </div>

                <div className="info-row">
                  <span className="label">价格:</span>
                  <span className="value price">
                    {formatPrice(BigInt(order.price))} USDC
                  </span>
                </div>

                <div className="info-row">
                  <span className="label">创建时间:</span>
                  <span className="value">
                    {new Date(order.createdAt).toLocaleString('zh-CN')}
                  </span>
                </div>

                <div className="info-row">
                  <span className="label">过期时间:</span>
                  <span className="value">
                    {new Date(order.expiry).toLocaleDateString('zh-CN')}
                  </span>
                </div>
              </div>

              {order.status === 'active' && (
                <div className="order-actions">
                  <button 
                    className="btn-cancel" 
                    onClick={() => handleCancelOrder(order)}
                    disabled={cancelling === order.id}
                  >
                    {cancelling === order.id ? '取消中...' : '取消订单'}
                  </button>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

