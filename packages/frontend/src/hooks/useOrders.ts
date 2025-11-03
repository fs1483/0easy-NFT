// 订单管理 Hook - 封装与后端 API 的交互逻辑
import { useState, useEffect } from 'react'
import type { Order } from '@oeasy-nft/common'
import { API_BASE_URL } from '../wagmi'

interface UseOrdersResult {
  orders: Order[]
  loading: boolean
  error: string | null
  refetch: () => void
}

/**
 * useOrders Hook 用于获取市场活跃订单列表
 * @param side - 订单类型筛选 ('ask' | 'bid' | undefined)
 * @param collection - NFT 合约地址筛选
 */
export function useOrders(side?: 'ask' | 'bid', collection?: string): UseOrdersResult {
  const [orders, setOrders] = useState<Order[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchOrders = async () => {
    try {
      setLoading(true)
      setError(null)

      const params = new URLSearchParams()
      if (side) params.append('side', side)
      if (collection) params.append('collection', collection)

      const url = `${API_BASE_URL}/api/orders${params.toString() ? '?' + params.toString() : ''}`
      const response = await fetch(url)

      if (!response.ok) {
        throw new Error(`获取订单失败: ${response.statusText}`)
      }

      const data = await response.json()
      setOrders(data.orders || [])
    } catch (err) {
      setError(err instanceof Error ? err.message : '未知错误')
      setOrders([])
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchOrders()
  }, [side, collection])

  return {
    orders,
    loading,
    error,
    refetch: fetchOrders,
  }
}

