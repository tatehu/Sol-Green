import { useState, useEffect } from 'react';
import { useWalletContext } from '../components/WalletConnect';
import axios from 'axios';
import './Marketing.css';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

export const Marketing = () => {
  const { publicKey } = useWalletContext();
  const [activities, setActivities] = useState([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    loadActivities();
  }, []);

  const loadActivities = async () => {
    setLoading(true);
    try {
      const res = await axios.get(`${API_BASE_URL}/api/v1/marketing/activities?status=active`);
      setActivities(res.data.data || []);
    } catch (err) {
      console.error('加载活动失败:', err);
    } finally {
      setLoading(false);
    }
  };

  const getActivityTypeText = (type) => {
    const typeMap = {
      'sign_in': '📅 每日签到',
      'invite': '👥 邀请好友',
      'daily_task': '✅ 每日任务',
      'lucky_draw': '🎁 幸运抽奖',
      'flash_sale': '⚡ 限时抢购',
      'festival': '🎉 节日活动',
    };
    return typeMap[type] || type;
  };

  const handleJoin = async (activityId) => {
    if (!publicKey) {
      alert('请先连接钱包');
      return;
    }

    const token = localStorage.getItem('token');
    if (!token) {
      alert('请先登录');
      return;
    }

    try {
      await axios.post(
        `${API_BASE_URL}/api/v1/marketing/activities/${activityId}/join`,
        {},
        { headers: { Authorization: `Bearer ${token}` } }
      );
      alert('参与成功！');
      loadActivities();
    } catch (err) {
      alert(err.response?.data?.error || '参与失败');
    }
  };

  return (
    <div className="marketing">
      <h2>🎁 营销活动</h2>
      
      {loading ? (
        <div className="loading">加载中...</div>
      ) : activities.length === 0 ? (
        <div className="empty-state">暂无进行中的活动</div>
      ) : (
        <div className="activities-grid">
          {activities.map((activity) => (
            <div key={activity.id} className="activity-card">
              <div className="activity-header">
                <h3>{activity.title}</h3>
                <span className="activity-type">{getActivityTypeText(activity.activity_type)}</span>
              </div>
              <p className="activity-description">{activity.description}</p>
              <div className="activity-info">
                <div className="info-item">
                  <span>奖励:</span>
                  <span className="reward">{activity.reward_amount} SOLGREEN</span>
                </div>
                <div className="info-item">
                  <span>参与人数:</span>
                  <span>{activity.current_participants}</span>
                </div>
              </div>
              {publicKey && (
                <button 
                  className="btn btn-primary"
                  onClick={() => handleJoin(activity.id)}
                >
                  立即参与
                </button>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
