// 订单 API 调用 hooks
// 封装与后端订单服务的交互逻辑

import { useState, useCallback } from 'react'
import { type Address, type Hex } from 'viem'
import { API_BASE_URL } from '../wagmi'
import type { Order } from '../utils/eip712'

/**
 * 后端订单响应数据结构
 */
export interface OrderResponse {
  id: number
  maker: Address
  nftAddress: Address
  tokenId: string
  paymentToken: Address
  price: string
  expiry: string // ISO 时间字符串
  nonce: string
  side: 'ask' | 'bid'
  status: 'active' | 'filled' | 'cancelled'
  signature: Hex
  hash: Hex
  createdAt: string
  updatedAt: string
}

/**
 * 创建订单请求参数
 */
export interface CreateOrderRequest {
  maker: string
  nft: string
  nftAddress: string  // 后端期望的字段名
  tokenId: string
  paymentToken: string
  price: string
  expiry: number
  nonce: string
  side: string
  signature: string
}

/**
 * 订单查询过滤参数
 */
export interface OrderFilters {
  side?: 'ask' | 'bid'
  collection?: Address
  status?: 'active' | 'filled' | 'cancelled'
}

/**
 * 订单 API Hook
 * 提供创建、查询、取消订单的功能
 */
export function useOrderAPI() {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  /**
   * 创建订单
   */
  const createOrder = useCallback(async (
    order: Order,
    signature: Hex
  ): Promise<OrderResponse | null> => {
    setLoading(true)
    setError(null)

    try {
      const request = {
        maker: order.maker.toLowerCase(),
        nftAddress: order.nft.toLowerCase(),
        tokenId: order.tokenId.toString(),
        paymentToken: order.paymentToken.toLowerCase(),
        price: order.price.toString(),
        expiry: Number(order.expiry),
        nonce: order.nonce.toString(),
        side: order.side === 0 ? 'ask' : 'bid',
        signature: signature.toLowerCase(),
      }
      
      console.log('📤 创建订单请求:', request)

      const response = await fetch(`${API_BASE_URL}/api/orders`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ error: '创建订单失败' }))
        throw new Error(errorData.error || `HTTP ${response.status}`)
      }

      const data = await response.json()
      return data as OrderResponse
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : '未知错误'
      setError(errorMsg)
      console.error('创建订单失败:', err)
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  /**
   * 查询订单列表
   */
  const fetchOrders = useCallback(async (
    filters?: OrderFilters
  ): Promise<OrderResponse[]> => {
    setLoading(true)
    setError(null)

    try {
      const params = new URLSearchParams()
      if (filters?.side) params.append('side', filters.side)
      if (filters?.collection) params.append('collection', filters.collection)
      if (filters?.status) params.append('status', filters.status || 'active')

      const url = `${API_BASE_URL}/api/orders${params.toString() ? '?' + params.toString() : ''}`
      
      console.log('正在请求订单:', url)
      
      const response = await fetch(url, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      })

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`)
      }

      const data = await response.json()
      console.log('订单数据:', data)
      
      // API 返回格式可能是 {orders: [...]} 或直接 [...]
      const orders = data.orders || data
      return Array.isArray(orders) ? orders : []
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : '未知错误'
      setError(`获取订单失败: ${errorMsg}`)
      console.error('查询订单失败:', err)
      console.error('API_BASE_URL:', API_BASE_URL)
      return []
    } finally {
      setLoading(false)
    }
  }, [])

  /**
   * 取消订单
   */
  const cancelOrder = useCallback(async (
    orderId: number,
    maker: string,
    nonce: string,
    cancelSignature: Hex
  ): Promise<boolean> => {
    setLoading(true)
    setError(null)

    try {
      const request = {
        maker: maker.toLowerCase(),
        nonce: nonce,
        signature: cancelSignature.toLowerCase(),
      }
      
      console.log('🗑️ 取消订单请求:', request)
      
      const response = await fetch(`${API_BASE_URL}/api/orders/${orderId}/cancel`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ error: '取消订单失败' }))
        throw new Error(errorData.error || `HTTP ${response.status}`)
      }

      return true
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : '未知错误'
      setError(errorMsg)
      console.error('取消订单失败:', err)
      return false
    } finally {
      setLoading(false)
    }
  }, [])

  return {
    loading,
    error,
    createOrder,
    fetchOrders,
    cancelOrder,
  }
}

