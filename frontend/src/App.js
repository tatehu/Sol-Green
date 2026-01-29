import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import { WalletContextProvider, WalletConnect } from './components/WalletConnect';
import { GreenSubmit } from './pages/GreenSubmit';
import { BehaviorStatus } from './pages/BehaviorStatus';
import { Challenges } from './pages/Challenges';
import { Marketing } from './pages/Marketing';
import './App.css';

function App() {
  return (
    <WalletContextProvider>
      <Router>
        <div className="App">
          <nav className="navbar">
            <div className="container">
              <Link to="/" className="logo">
                <h1>🌱 Sol-Green</h1>
              </Link>
              <div className="nav-links">
                <Link to="/">首页</Link>
                <Link to="/submit">提交行为</Link>
                <Link to="/challenges">挑战活动</Link>
                <Link to="/marketing">营销活动</Link>
                <Link to="/status">查询状态</Link>
                <WalletConnect />
              </div>
            </div>
          </nav>

          <main className="main-content">
            <div className="container">
              <Routes>
                <Route path="/" element={<Home />} />
                <Route path="/submit" element={<GreenSubmit />} />
                <Route path="/challenges" element={<Challenges />} />
                <Route path="/marketing" element={<Marketing />} />
                <Route path="/status" element={<BehaviorStatus />} />
              </Routes>
            </div>
          </main>
        </div>
      </Router>
    </WalletContextProvider>
  );
}

function Home() {
  return (
    <div className="home">
      <div className="hero">
        <h2>欢迎使用 Sol-Green 环保奖励平台</h2>
        <p>通过区块链技术记录和奖励您的环保行为</p>
        <div className="features">
          <div className="feature-card">
            <h3>🌿 垃圾分类</h3>
            <p>记录您的垃圾分类行为，获得代币奖励</p>
          </div>
          <div className="feature-card">
            <h3>🌳 植树造林</h3>
            <p>上传植树证明，获得丰厚奖励</p>
          </div>
          <div className="feature-card">
            <h3>🚴 低碳出行</h3>
            <p>记录低碳出行方式，为环保贡献力量</p>
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
