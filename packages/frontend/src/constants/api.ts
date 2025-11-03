// API 端点常量配置

/**
 * API 基础 URL 配置
 */
export const API_BASE_URLS = {
  development: 'http://localhost:8081',
  staging: 'https://api-staging.oeasy-nft.com',
  production: 'https://api.oeasy-nft.com',
} as const

/**
 * 获取当前环境的 API URL
 */
export function getApiBaseUrl(): string {
  const env = import.meta.env.MODE || 'development'
  
  if (env === 'production') {
    return API_BASE_URLS.production
  } else if (env === 'staging') {
    return API_BASE_URLS.staging
  }
  
  return API_BASE_URLS.development
}

/**
 * API 端点路径
 */
export const API_ENDPOINTS = {
  // 订单相关
  ORDERS: '/api/orders',
  ORDER_BY_ID: (id: number) => `/api/orders/${id}`,
  CANCEL_ORDER: (id: number) => `/api/orders/${id}/cancel`,
  
  // 健康检查
  HEALTH: '/health',
} as const

/**
 * API 请求超时时间（毫秒）
 */
export const API_TIMEOUT = {
  DEFAULT: 10000,      // 10 秒
  SIGNATURE: 30000,    // 30 秒（签名操作）
  LONG_POLL: 60000,    // 60 秒（长轮询）
} as const

/**
 * HTTP 状态码
 */
export const HTTP_STATUS = {
  OK: 200,
  CREATED: 201,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  NOT_FOUND: 404,
  CONFLICT: 409,
  INTERNAL_ERROR: 500,
} as const

