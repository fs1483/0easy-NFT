// 订单列表组件
// 显示市场上的活跃订单，支持按类型筛选

import { useEffect, useState } from 'react'
import { useOrderAPI, type OrderResponse } from '../hooks/useOrderAPI'
import { formatPrice, formatExpiry } from '../utils/eip712'
import './OrderList.css'

interface OrderListProps {
  filterSide?: 'ask' | 'bid'
  filterCollection?: string
  onSelectOrder?: (order: OrderResponse) => void
  onMakeOffer?: (order: OrderResponse) => void  // 新增：出价回调
}

export function OrderList({ filterSide, filterCollection, onSelectOrder, onMakeOffer }: OrderListProps) {
  const { loading, error, fetchOrders } = useOrderAPI()
  const [orders, setOrders] = useState<OrderResponse[]>([])
  const [activeFilter, setActiveFilter] = useState<'all' | 'ask' | 'bid'>(filterSide || 'all')

  // 加载订单列表
  useEffect(() => {
    loadOrders()
    // 每 10 秒刷新一次
    const interval = setInterval(loadOrders, 10000)
    return () => clearInterval(interval)
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [filterSide, filterCollection])

  async function loadOrders() {
    const filters = {
      side: filterSide,
      collection: filterCollection as any,
      status: 'active' as const,
    }
    const data = await fetchOrders(filters)
    setOrders(data)
  }

  // 筛选订单
  const filteredOrders = orders.filter(order => {
    if (activeFilter === 'all') return true
    return order.side === activeFilter
  })

  // 格式化订单类型显示
  function getOrderTypeLabel(side: 'ask' | 'bid'): string {
    return side === 'ask' ? '卖单' : '买单'
  }

  // 格式化订单类型样式
  function getOrderTypeClass(side: 'ask' | 'bid'): string {
    return side === 'ask' ? 'order-type-ask' : 'order-type-bid'
  }

  return (
    <div className="order-list-container">
      <div className="order-list-header">
        <h2>市场订单</h2>
        <div className="filter-buttons">
          <button
            className={activeFilter === 'all' ? 'filter-btn active' : 'filter-btn'}
            onClick={() => setActiveFilter('all')}
          >
            全部订单
          </button>
          <button
            className={activeFilter === 'ask' ? 'filter-btn active' : 'filter-btn'}
            onClick={() => setActiveFilter('ask')}
          >
            卖单 (Ask)
          </button>
          <button
            className={activeFilter === 'bid' ? 'filter-btn active' : 'filter-btn'}
            onClick={() => setActiveFilter('bid')}
          >
            买单 (Bid)
          </button>
        </div>
      </div>

      {loading && <div className="loading">加载中...</div>}
      {error && <div className="error">错误: {error}</div>}

      {!loading && filteredOrders.length === 0 && (
        <div className="empty-state">
          <p>暂无活跃订单</p>
          <p className="hint">成为第一个创建订单的人！</p>
        </div>
      )}

      {!loading && filteredOrders.length > 0 && (
        <div className="orders-grid">
          {filteredOrders.map(order => (
            <div
              key={order.id}
              className="order-card"
              onClick={() => onSelectOrder?.(order)}
            >
              <div className="order-header">
                <span className={`order-type ${getOrderTypeClass(order.side)}`}>
                  {getOrderTypeLabel(order.side)}
                </span>
                <span className="order-id">#{order.id}</span>
              </div>

              <div className="order-body">
                <div className="order-info-row">
                  <span className="label">NFT:</span>
                  <span className="value truncate" title={order.nftAddress}>
                    OeasyNFT #{order.tokenId}
                  </span>
                </div>

                <div className="order-info-row">
                  <span className="label">合约:</span>
                  <span className="value truncate" title={order.nftAddress}>
                    {order.nftAddress.slice(0, 10)}...
                  </span>
                </div>

                <div className="order-info-row">
                  <span className="label">价格:</span>
                  <span className="value price">
                    {formatPrice(BigInt(order.price))} USDC
                  </span>
                </div>

                <div className="order-info-row">
                  <span className="label">制作者:</span>
                  <span className="value truncate">
                    {order.maker.slice(0, 6)}...{order.maker.slice(-4)}
                  </span>
                </div>

                <div className="order-info-row">
                  <span className="label">过期:</span>
                  <span className="value expiry">
                    {new Date(order.expiry).toLocaleDateString('zh-CN')}
                  </span>
                </div>
              </div>

              <div className="order-footer">
                <div className="order-meta">
                  <span className="status">{order.status === 'active' ? '活跃' : order.status}</span>
                  <span className="created-at">
                    {new Date(order.createdAt).toLocaleString('zh-CN')}
                  </span>
                </div>
                
                {/* 如果是卖单，显示"出价"按钮 */}
                {order.side === 'ask' && onMakeOffer && (
                  <button 
                    className="btn-make-offer"
                    onClick={(e) => {
                      e.stopPropagation()  // 阻止卡片点击事件
                      onMakeOffer(order)
                    }}
                  >
                    💰 出价
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

