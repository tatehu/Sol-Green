import { useState } from 'react';
import { useWalletContext } from '../components/WalletConnect';
import axios from 'axios';
import './BehaviorStatus.css';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

export const BehaviorStatus = () => {
  const { publicKey } = useWalletContext();
  const [behaviorId, setBehaviorId] = useState('');
  const [loading, setLoading] = useState(false);
  const [behavior, setBehavior] = useState(null);
  const [error, setError] = useState('');

  const handleQuery = async () => {
    if (!behaviorId.trim()) {
      setError('请输入行为ID');
      return;
    }

    const token = localStorage.getItem('token');
    if (!token) {
      setError('请先登录');
      return;
    }

    setLoading(true);
    setError('');
    setBehavior(null);

    try {
      const res = await axios.get(
        `${API_BASE_URL}/api/v1/green/behavior/${behaviorId}`,
        {
          headers: { Authorization: `Bearer ${token}` }
        }
      );

      setBehavior(res.data.data);
    } catch (err) {
      setError(err.response?.data?.error || '查询失败');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const getStatusText = (status) => {
    const statusMap = {
      'pending_review': '⏳ 待审核',
      'approved': '✅ 已通过',
      'rejected': '❌ 已拒绝',
    };
    return statusMap[status] || status;
  };

  const getBehaviorTypeText = (type) => {
    const typeMap = {
      'waste_sorting': '🗑️ 垃圾分类',
      'tree_planting': '🌳 植树造林',
      'low_carbon_travel': '🚴 低碳出行',
    };
    return typeMap[type] || type;
  };

  return (
    <div className="behavior-status">
      <div className="card">
        <h2>📋 查询行为状态</h2>

        <div className="form-group">
          <label className="label">行为ID</label>
          <input
            type="text"
            value={behaviorId}
            onChange={(e) => setBehaviorId(e.target.value)}
            placeholder="请输入行为ID"
            className="input"
          />
        </div>

        <button 
          onClick={handleQuery}
          disabled={loading}
          className="btn btn-primary"
        >
          {loading ? '查询中...' : '查询'}
        </button>

        {error && (
          <div className="message error">
            {error}
          </div>
        )}

        {behavior && (
          <div className="behavior-detail">
            <h3>行为详情</h3>
            <div className="detail-item">
              <strong>行为ID:</strong> {behavior.id}
            </div>
            <div className="detail-item">
              <strong>钱包地址:</strong> {behavior.wallet_addr}
            </div>
            <div className="detail-item">
              <strong>行为类型:</strong> {getBehaviorTypeText(behavior.behavior_type)}
            </div>
            <div className="detail-item">
              <strong>状态:</strong> {getStatusText(behavior.status)}
            </div>
            <div className="detail-item">
              <strong>欺诈分数:</strong> {behavior.fraud_score?.toFixed(2) || 'N/A'}
            </div>
            <div className="detail-item">
              <strong>提交时间:</strong> {new Date(behavior.submit_time).toLocaleString('zh-CN')}
            </div>
            {behavior.approve_time && (
              <div className="detail-item">
                <strong>通过时间:</strong> {new Date(behavior.approve_time).toLocaleString('zh-CN')}
              </div>
            )}
            {behavior.location && (
              <div className="detail-item">
                <strong>地点:</strong> {behavior.location}
              </div>
            )}
            {behavior.tx_hash && (
              <div className="detail-item">
                <strong>交易哈希:</strong>{' '}
                <a 
                  href={`https://explorer.solana.com/tx/${behavior.tx_hash}?cluster=devnet`}
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  {behavior.tx_hash}
                </a>
              </div>
            )}
            {behavior.proof_hash && (
              <div className="detail-item">
                <strong>存证哈希:</strong> {behavior.proof_hash}
              </div>
            )}
            {behavior.media_urls && behavior.media_urls.length > 0 && (
              <div className="detail-item">
                <strong>媒体文件:</strong>
                <div className="media-list">
                  {behavior.media_urls.map((url, idx) => (
                    <a 
                      key={idx}
                      href={url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="media-link"
                    >
                      文件 {idx + 1}
                    </a>
                  ))}
                </div>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
};
