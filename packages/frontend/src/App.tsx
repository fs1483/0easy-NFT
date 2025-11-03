import { useState } from 'react'
import { ConnectButton } from '@rainbow-me/rainbowkit'
import { useAccount } from 'wagmi'
import { OrderList } from './components/OrderList'
import { CreateOrderForm } from './components/CreateOrderForm'
import { TradeHistory } from './components/TradeHistory'
import { MyOrders } from './components/MyOrders'
import { MintNFT } from './components/MintNFT'
import { GetTestUSDC } from './components/GetTestUSDC'
import './App.css'

function App() {
  const { address, isConnected } = useAccount()
  const [activeTab, setActiveTab] = useState<'market' | 'create' | 'history' | 'myorders' | 'mint' | 'faucet'>('market')
  const [refreshKey, setRefreshKey] = useState(0)
  const [selectedOrder, setSelectedOrder] = useState<any>(null)

  // 订单创建成功后刷新列表
  const handleOrderCreated = () => {
    setRefreshKey(prev => prev + 1)
    setSelectedOrder(null)
    setActiveTab('market')
  }
  
  // 处理"出价"按钮点击
  const handleMakeOffer = (order: any) => {
    // 计算建议出价（等于卖价，用户可以修改）
    const suggestedPrice = (BigInt(order.price) / 1000000n).toString()
    
    setSelectedOrder({
      ...order,
      suggestedPrice
    })
    setActiveTab('create')
  }

  return (
    <div className="app-container">
      <header className="app-header">
        <div className="header-content">
          <div className="logo-section">
            <h1>🎨 Oeasy NFT</h1>
            <span className="subtitle">订单簿交易平台</span>
          </div>
          <ConnectButton />
        </div>
      </header>

      <main className="main-content">
        {!isConnected ? (
          <div className="connect-prompt">
            <h2>欢迎来到 Oeasy NFT 订单簿交易平台</h2>
            <p className="lead">链下挂单，链上结算 - 企业级 Web3 交易体验</p>
            <div className="features">
              <div className="feature-card">
                <div className="feature-icon">🎯</div>
                <h3>链下挂单</h3>
                <p>EIP-712 签名，无 Gas 费创建订单</p>
              </div>
              <div className="feature-card">
                <div className="feature-icon">⚡</div>
                <h3>快速撮合</h3>
                <p>企业级撮合引擎，实时匹配买卖订单</p>
              </div>
              <div className="feature-card">
                <div className="feature-icon">🔒</div>
                <h3>链上结算</h3>
                <p>智能合约原子化执行，去中心化保障</p>
              </div>
            </div>
            <div className="cta-section">
              <ConnectButton />
              <p className="cta-hint">连接钱包开始交易</p>
            </div>
          </div>
        ) : (
          <>
            <nav className="tabs">
              <button
                className={activeTab === 'market' ? 'tab active' : 'tab'}
                onClick={() => setActiveTab('market')}
              >
                <span className="tab-icon">📊</span>
                市场订单
              </button>
              <button
                className={activeTab === 'mint' ? 'tab active' : 'tab'}
                onClick={() => setActiveTab('mint')}
              >
                <span className="tab-icon">🎨</span>
                铸造 NFT
              </button>
              <button
                className={activeTab === 'faucet' ? 'tab active' : 'tab'}
                onClick={() => setActiveTab('faucet')}
              >
                <span className="tab-icon">💰</span>
                测试代币
              </button>
              <button
                className={activeTab === 'myorders' ? 'tab active' : 'tab'}
                onClick={() => setActiveTab('myorders')}
              >
                <span className="tab-icon">👤</span>
                我的订单
              </button>
              <button
                className={activeTab === 'history' ? 'tab active' : 'tab'}
                onClick={() => setActiveTab('history')}
              >
                <span className="tab-icon">✅</span>
                交易历史
              </button>
              <button
                className={activeTab === 'create' ? 'tab active' : 'tab'}
                onClick={() => setActiveTab('create')}
              >
                <span className="tab-icon">➕</span>
                创建订单
              </button>
            </nav>

            <div className="tab-content">
              {activeTab === 'market' ? (
                <OrderList 
                  key={refreshKey} 
                  onMakeOffer={handleMakeOffer}
                />
              ) : activeTab === 'mint' ? (
                <MintNFT />
              ) : activeTab === 'faucet' ? (
                <GetTestUSDC />
              ) : activeTab === 'myorders' ? (
                <MyOrders key={refreshKey} />
              ) : activeTab === 'history' ? (
                <TradeHistory key={refreshKey} />
              ) : (
                <CreateOrderForm 
                  nftAddress={selectedOrder?.nftAddress}
                  tokenId={selectedOrder?.tokenId}
                  defaultSide="bid"
                  defaultPrice={selectedOrder?.suggestedPrice}
                  onSuccess={handleOrderCreated}
                  onCancel={() => {
                    setSelectedOrder(null)
                    setActiveTab('market')
                  }}
                />
              )}
            </div>
          </>
        )}
      </main>

      <footer className="app-footer">
        <div className="footer-content">
          <div className="footer-info">
            <p className="footer-title">Oeasy NFT © 2025</p>
            <p className="footer-desc">准企业级订单簿交易平台 | 由 Web3 技术驱动</p>
          </div>
          <div className="footer-links">
            <a href="https://github.com/your-username/oeasy-nft" target="_blank" rel="noopener noreferrer">GitHub</a>
            <a href="#" onClick={(e) => { e.preventDefault(); alert('文档：查看 docs/ 目录') }}>文档</a>
          </div>
        </div>
      </footer>
    </div>
  )
}

export default App
