// 交易历史组件
// 显示已成交的订单和链上交易事件

import { useEffect, useState } from 'react'
import { API_BASE_URL } from '../wagmi'
import { formatPrice } from '../utils/eip712'
import './TradeHistory.css'

interface TradeEvent {
  id: number
  transactionHash: string
  blockNumber: number
  maker: string
  taker: string
  nftAddress: string
  tokenId: string
  paymentToken: string
  price: string
  side: number
  fee: string
  createdAt: string
}

export function TradeHistory() {
  const [trades, setTrades] = useState<TradeEvent[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    loadTrades()
    // 每 10 秒刷新
    const interval = setInterval(loadTrades, 10000)
    return () => clearInterval(interval)
  }, [])

  async function loadTrades() {
    setLoading(true)
    setError(null)

    try {
      // 查询已成交的订单
      const response = await fetch(`${API_BASE_URL}/api/orders?status=filled`)
      
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}`)
      }

      const data = await response.json()
      const orders = data.orders || data || []
      
      // 转换为交易历史格式
      const tradeList = orders.map((order: any) => ({
        id: order.id,
        transactionHash: order.hash,
        blockNumber: 0,
        maker: order.maker,
        taker: '未知',
        nftAddress: order.nftAddress,
        tokenId: order.tokenId,
        paymentToken: order.paymentToken,
        price: order.price,
        side: order.side === 'ask' ? 0 : 1,
        fee: '0',
        createdAt: order.updatedAt || order.createdAt,
      }))
      
      setTrades(tradeList)
    } catch (err) {
      setError(err instanceof Error ? err.message : '加载失败')
      console.error('加载交易历史失败:', err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="trade-history-container">
      <div className="trade-history-header">
        <h2>交易历史</h2>
        <p className="subtitle">已成交的订单记录</p>
      </div>

      {loading && <div className="loading">加载中...</div>}
      {error && <div className="error">错误: {error}</div>}

      {!loading && trades.length === 0 && (
        <div className="empty-state">
          <p>暂无成交记录</p>
          <p className="hint">创建订单并等待撮合成功！</p>
        </div>
      )}

      {!loading && trades.length > 0 && (
        <div className="trades-list">
          {trades.map(trade => (
            <div key={trade.id} className="trade-card">
              <div className="trade-header">
                <span className="trade-type">
                  {trade.side === 0 ? '🔴 卖单成交' : '🔵 买单成交'}
                </span>
                <span className="trade-time">
                  {new Date(trade.createdAt).toLocaleString('zh-CN')}
                </span>
              </div>

              <div className="trade-body">
                <div className="trade-info">
                  <div className="info-row">
                    <span className="label">NFT:</span>
                    <span className="value">
                      OeasyNFT #{trade.tokenId}
                    </span>
                  </div>

                  <div className="info-row">
                    <span className="label">合约:</span>
                    <span className="value truncate" title={trade.nftAddress}>
                      {trade.nftAddress.slice(0, 10)}...{trade.nftAddress.slice(-6)}
                    </span>
                  </div>

                  <div className="info-row">
                    <span className="label">成交价:</span>
                    <span className="value price">
                      {formatPrice(BigInt(trade.price))} USDC
                    </span>
                  </div>

                  <div className="info-row">
                    <span className="label">卖方:</span>
                    <span className="value truncate">
                      {trade.maker.slice(0, 6)}...{trade.maker.slice(-4)}
                    </span>
                  </div>

                  {trade.taker !== '未知' && (
                    <div className="info-row">
                      <span className="label">买方:</span>
                      <span className="value truncate">
                        {trade.taker.slice(0, 6)}...{trade.taker.slice(-4)}
                      </span>
                    </div>
                  )}
                </div>
              </div>

              <div className="trade-footer">
                <span className="status-badge">✅ 已成交</span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

