import { useState, useEffect, useCallback } from 'react';
import { useWalletContext } from '../components/WalletConnect';
import axios from 'axios';
import './Marketing.css';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

export const Marketing = () => {
  const { publicKey } = useWalletContext();
  const [activities, setActivities] = useState([]);
  const [loading, setLoading] = useState(false);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [selected, setSelected] = useState(null);
  const [joinedMap, setJoinedMap] = useState(() => ({}));

  const loadActivities = useCallback(async () => {
    setLoading(true);
    try {
      // 后端会返回 active 或 scheduled 的活动；前端不应只筛 active，否则新创建的 scheduled 活动永远看不到
      const res = await axios.get(`${API_BASE_URL}/api/v1/marketing/activities`);
      setActivities(res.data.data || []);
    } catch (err) {
      console.error('加载活动失败:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadActivities();
  }, [loadActivities]);

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
      setJoinedMap((m) => ({ ...m, [activityId]: true }));
      alert('参与成功！');
      loadActivities();
    } catch (err) {
      alert(err.response?.data?.error || '参与失败');
    }
  };

  const handleClaim = async (activityId) => {
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
      const res = await axios.post(
        `${API_BASE_URL}/api/v1/marketing/activities/${activityId}/claim`,
        {},
        { headers: { Authorization: `Bearer ${token}` } }
      );
      alert(`领取成功！交易哈希：${res.data.tx_hash || ''}`);
      loadActivities();
    } catch (err) {
      alert(err.response?.data?.error || '领取失败');
    }
  };

  return (
    <div className="marketing">
      <div className="marketing-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h2>🎁 营销活动</h2>
        {publicKey && (
          <button
            className="btn btn-primary"
            onClick={() => setShowCreateModal(true)}
          >
            + 发起活动
          </button>
        )}
      </div>
      
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
                <div style={{ display: 'flex', gap: 10 }}>
                  <button
                    className="btn btn-secondary"
                    onClick={() => setSelected(activity)}
                  >
                    查看详情
                  </button>
                  <button
                    className="btn btn-primary"
                    onClick={() => handleJoin(activity.id)}
                  >
                    立即参与
                  </button>
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      {showCreateModal && (
        <CreateMarketingModal
          onClose={() => setShowCreateModal(false)}
          onSuccess={(created) => {
            setShowCreateModal(false);
            if (created?.id) {
              setActivities((prev) => [created, ...prev.filter((a) => a.id !== created.id)]);
            }
            loadActivities();
          }}
        />
      )}

      {selected && (
        <ActivityDetailModal
          activity={selected}
          joined={Boolean(joinedMap[selected.id])}
          onJoin={() => handleJoin(selected.id)}
          onClaim={() => handleClaim(selected.id)}
          onClose={() => setSelected(null)}
        />
      )}
    </div>
  );
};

const ActivityDetailModal = ({ activity, joined, onJoin, onClaim, onClose }) => {
  const isSignIn = activity.activity_type === 'sign_in';
  const today = new Date();
  const days = Array.from({ length: 7 }).map((_, i) => {
    const d = new Date();
    d.setDate(today.getDate() - (6 - i));
    return d;
  });

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <h3>{activity.title}</h3>
        <p style={{ color: '#666' }}>{activity.description}</p>

        <div className="activity-info" style={{ marginTop: 12 }}>
          <div className="info-item">
            <span>类型:</span>
            <span>{activity.activity_type}</span>
          </div>
          <div className="info-item">
            <span>状态:</span>
            <span>{activity.status}</span>
          </div>
          <div className="info-item">
            <span>时间:</span>
            <span>
              {new Date(activity.start_time).toLocaleString('zh-CN')} - {new Date(activity.end_time).toLocaleString('zh-CN')}
            </span>
          </div>
          <div className="info-item">
            <span>奖励:</span>
            <span className="reward">{activity.reward_amount} SOLGREEN</span>
          </div>
        </div>

        {isSignIn ? (
          <div style={{ marginTop: 16 }}>
            <h4>签到日历（近 7 天）</h4>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(7, 1fr)', gap: 8 }}>
              {days.map((d, idx) => {
                const isToday = d.toDateString() === today.toDateString();
                return (
                  <div
                    key={idx}
                    style={{
                      border: '1px solid #e5e7eb',
                      borderRadius: 8,
                      padding: 10,
                      textAlign: 'center',
                      background: isToday ? '#e8f5e9' : '#fff',
                      fontWeight: isToday ? 700 : 400,
                    }}
                  >
                    <div style={{ fontSize: 12, color: '#64748b' }}>
                      {d.getMonth() + 1}/{d.getDate()}
                    </div>
                    <div style={{ marginTop: 6, fontSize: 14 }}>{isToday ? '今天' : ' '}</div>
                  </div>
                );
              })}
            </div>
            <p style={{ marginTop: 10, color: '#64748b', fontSize: 13 }}>
              参考文档效果：今日高亮；已签到显示 ✓（当前为简化演示）。
            </p>
          </div>
        ) : null}

        {activity.rules ? (
          <div style={{ marginTop: 16 }}>
            <h4>活动规则</h4>
            <p style={{ whiteSpace: 'pre-wrap' }}>{activity.rules}</p>
          </div>
        ) : null}

        <div style={{ display: 'flex', gap: 10, justifyContent: 'flex-end', marginTop: 20 }}>
          <button className="btn btn-secondary" onClick={onClose}>关闭</button>
          <button className="btn btn-secondary" onClick={onClaim} disabled={!joined}>
            领取奖励
          </button>
          <button className="btn btn-primary" onClick={onJoin}>
            {joined ? '已参与' : (isSignIn ? '立即签到' : '立即参与')}
          </button>
        </div>
      </div>
    </div>
  );
};

// 创建营销活动模态框
const CreateMarketingModal = ({ onClose, onSuccess }) => {
  const { publicKey } = useWalletContext();
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    activity_type: 'sign_in',
    start_time: '',
    end_time: '',
    reward_amount: 1000,
    max_participants: 0,
    rules: '',
    image_url: '',
    config: '',
  });
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!publicKey) {
      setError('请先连接钱包');
      return;
    }

    const token = localStorage.getItem('token');
    if (!token) {
      setError('请先登录');
      return;
    }

    // 基础表单校验（避免后端兜底 + 减少无效请求）
    if (!formData.title.trim() || !formData.description.trim()) {
      setError('请填写活动标题与描述');
      return;
    }
    if (!formData.start_time || !formData.end_time) {
      setError('请选择开始/结束时间');
      return;
    }
    const start = new Date(formData.start_time);
    const end = new Date(formData.end_time);
    if (Number.isNaN(start.getTime()) || Number.isNaN(end.getTime())) {
      setError('时间格式不正确');
      return;
    }
    if (end <= start) {
      setError('结束时间必须晚于开始时间');
      return;
    }
    if (!formData.reward_amount || formData.reward_amount < 1) {
      setError('奖励数量必须大于 0');
      return;
    }

    try {
      setSubmitting(true);
      setError('');
      const res = await axios.post(
        `${API_BASE_URL}/api/v1/marketing/activities`,
        formData,
        { headers: { Authorization: `Bearer ${token}` } }
      );
      onSuccess(res.data?.data);
    } catch (err) {
      setError(err.response?.data?.detail || err.response?.data?.error || '创建失败');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <h3>发起营销活动</h3>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>活动标题 *</label>
            <input
              type="text"
              value={formData.title}
              onChange={(e) => setFormData({ ...formData, title: e.target.value })}
              required
            />
          </div>

          <div className="form-group">
            <label>活动描述 *</label>
            <textarea
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              required
            />
          </div>

          <div className="form-group">
            <label>活动类型 *</label>
            <select
              value={formData.activity_type}
              onChange={(e) => setFormData({ ...formData, activity_type: e.target.value })}
            >
              <option value="sign_in">每日签到</option>
              <option value="invite">邀请好友</option>
              <option value="daily_task">每日任务</option>
              <option value="lucky_draw">幸运抽奖</option>
              <option value="flash_sale">限时抢购</option>
              <option value="festival">节日活动</option>
            </select>
          </div>

          <div className="form-row">
            <div className="form-group">
              <label>开始时间 *</label>
              <input
                type="datetime-local"
                value={formData.start_time}
                onChange={(e) => setFormData({ ...formData, start_time: e.target.value })}
                required
              />
            </div>
            <div className="form-group">
              <label>结束时间 *</label>
              <input
                type="datetime-local"
                value={formData.end_time}
                onChange={(e) => setFormData({ ...formData, end_time: e.target.value })}
                required
              />
            </div>
          </div>

          <div className="form-row">
            <div className="form-group">
              <label>奖励数量 *</label>
              <input
                type="number"
                min="1"
                value={formData.reward_amount}
                onChange={(e) => setFormData({ ...formData, reward_amount: parseInt(e.target.value || '0', 10) })}
                required
              />
            </div>
            <div className="form-group">
              <label>最大参与人数（0=不限）</label>
              <input
                type="number"
                min="0"
                value={formData.max_participants}
                onChange={(e) => setFormData({ ...formData, max_participants: parseInt(e.target.value || '0', 10) })}
              />
            </div>
          </div>

          <div className="form-group">
            <label>规则（可选）</label>
            <textarea
              value={formData.rules}
              onChange={(e) => setFormData({ ...formData, rules: e.target.value })}
            />
          </div>

          {error ? <div className="error-message">{error}</div> : null}

          <div className="modal-actions">
            <button type="button" onClick={onClose} className="btn btn-secondary">
              取消
            </button>
            <button type="submit" disabled={submitting} className="btn btn-primary">
              {submitting ? '创建中...' : '创建活动'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
