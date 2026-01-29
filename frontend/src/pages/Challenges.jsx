import { useState, useEffect, useCallback } from 'react';
import { useWalletContext } from '../components/WalletConnect';
import axios from 'axios';
import './Challenges.css';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

const getDefaultChallengeBg = (behaviorType) => {
  switch (behaviorType) {
    case 'waste_sorting':
      return 'linear-gradient(135deg, rgba(34,197,94,0.20), rgba(20,184,166,0.10))';
    case 'tree_planting':
      return 'linear-gradient(135deg, rgba(16,185,129,0.22), rgba(34,197,94,0.10))';
    case 'low_carbon_travel':
      return 'linear-gradient(135deg, rgba(59,130,246,0.22), rgba(56,189,248,0.10))';
    default:
      return 'linear-gradient(135deg, rgba(148,163,184,0.18), rgba(226,232,240,0.10))';
  }
};

// 根据状态和时间判断显示文本
const getStatusText = (status, startTime, endTime) => {
  const now = Date.now();
  const start = startTime ? new Date(startTime).getTime() : 0;
  const end = endTime ? new Date(endTime).getTime() : 0;

  // 挑战时间已结束，则统一展示“已结束”
  if (end > 0 && now >= end) {
    return '⏱ 已结束';
  }

  // 如果状态是草稿，但开始时间还没到，显示"即将开始"
  if (status === 'draft' && start > 0 && now < start) {
    return '⏰ 即将开始';
  }
  // 如果状态是草稿，但已经开始，显示"进行中"
  if (status === 'draft' && start > 0 && now >= start && now < end) {
    return '🔥 进行中';
  }
  // 如果状态是active，但还没开始，显示"即将开始"
  if (status === 'active' && start > 0 && now < start) {
    return '⏰ 即将开始';
  }

  const statusMap = {
    draft: '📝 草稿',
    active: '🔥 进行中',
    completed: '✅ 已完成',
    cancelled: '❌ 已取消',
  };
  return statusMap[status] || status || '-';
};

// 视觉状态对应的样式（避免“进行中”文字相同但底色不同）
const getStatusClass = (status, startTime, endTime) => {
  const now = Date.now();
  const start = startTime ? new Date(startTime).getTime() : 0;
  const end = endTime ? new Date(endTime).getTime() : 0;

  if (end > 0 && now >= end || status === 'completed') return 'status-completed';
  if (status === 'cancelled') return 'status-cancelled';

  // 时间在范围内，一律按“进行中”视觉处理
  if (start > 0 && now >= start && (end === 0 || now < end)) {
    return 'status-active';
  }

  // 其余都按“即将开始/草稿”的灰色
  return 'status-draft';
};

const getBehaviorTypeText = (type) => {
  const typeMap = {
    waste_sorting: '🗑️ 垃圾分类',
    tree_planting: '🌳 植树造林',
    low_carbon_travel: '🚴 低碳出行',
  };
  return typeMap[type] || type || '-';
};

const getChallengeTypeText = (type) => {
  const typeMap = {
    individual: '👤 个人挑战',
    team: '👥 团队挑战',
    community: '🏘️ 社区挑战',
  };
  return typeMap[type] || type || '-';
};

export const Challenges = () => {
  const { publicKey } = useWalletContext();
  const [challenges, setChallenges] = useState([]);
  const [loading, setLoading] = useState(false);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [selectedChallenge, setSelectedChallenge] = useState(null);
  const [joinChallengeId, setJoinChallengeId] = useState(null); // 要参与的挑战ID

  const loadChallenges = useCallback(async () => {
    setLoading(true);
    try {
      const res = await axios.get(`${API_BASE_URL}/api/v1/challenges`);
      setChallenges(res.data.data || []);
    } catch (err) {
      console.error('加载挑战失败:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadChallenges();
  }, [loadChallenges]);

  const formatDate = (dateStr) => {
    return new Date(dateStr).toLocaleString('zh-CN');
  };

  const calculateProgress = (current, target) => {
    return Math.min((current / target) * 100, 100);
  };

  return (
    <div className="challenges">
      <div className="challenges-header">
        <h2>🌱 环保挑战活动</h2>
        {publicKey && (
          <button 
            className="btn btn-primary"
            onClick={() => setShowCreateModal(true)}
          >
            + 发起挑战
          </button>
        )}
      </div>

      {loading ? (
        <div className="loading">加载中...</div>
      ) : challenges.length === 0 ? (
        <div className="empty-state">
          <p>暂无进行中的挑战活动</p>
        </div>
      ) : (
        <div className="challenges-grid">
          {challenges.map((challenge) => (
            <div
              key={challenge.id}
              className="challenge-card"
              style={{
                background: challenge.image_url ? undefined : getDefaultChallengeBg(challenge.behavior_type),
                backgroundImage: challenge.image_url ? `url(${challenge.image_url})` : undefined,
                backgroundSize: challenge.image_url ? 'cover' : undefined,
                backgroundPosition: challenge.image_url ? 'center' : undefined,
              }}
            >
              <div className="challenge-card-overlay" />
              <div className="challenge-card-content">
                <div className="challenge-header">
                  <h3>{challenge.title}</h3>
                  <span
                    className={`status-badge ${getStatusClass(challenge.status, challenge.start_time, challenge.end_time)}`}
                  >
                    {getStatusText(challenge.status, challenge.start_time, challenge.end_time)}
                  </span>
                </div>
              
                <p className="challenge-description">{challenge.description}</p>
              
                <div className="challenge-info">
                  <div className="info-item">
                    <span className="label">类型:</span>
                    <span>{getBehaviorTypeText(challenge.behavior_type)}</span>
                  </div>
                  <div className="info-item">
                    <span className="label">奖励:</span>
                    <span className="reward">{challenge.reward_amount} SOLGREEN</span>
                  </div>
                  <div className="info-item">
                    <span className="label">时间:</span>
                    <span>{formatDate(challenge.start_time)} - {formatDate(challenge.end_time)}</span>
                  </div>
                </div>

                <div className="challenge-progress">
                  <div className="progress-label">
                    <span>参与进度</span>
                    <span>{challenge.current_count} / {challenge.target_count}</span>
                  </div>
                  <div className="progress-bar">
                    <div 
                      className="progress-fill"
                      style={{ width: `${calculateProgress(challenge.current_count, challenge.target_count)}%` }}
                    />
                  </div>
                </div>

                <div className="challenge-actions">
                  <button 
                    className="btn btn-secondary"
                    onClick={() => setSelectedChallenge(challenge)}
                  >
                    查看详情
                  </button>
                  {(() => {
                    const now = Date.now();
                    const start = challenge.start_time ? new Date(challenge.start_time).getTime() : 0;
                    const end = challenge.end_time ? new Date(challenge.end_time).getTime() : 0;
                    const canJoin = (challenge.status === 'active' || challenge.status === 'draft') && 
                                    start > 0 && now >= start && now < end && publicKey;
                    return canJoin ? (
                      <button 
                        className="btn btn-primary"
                        onClick={() => setJoinChallengeId(challenge.id)}
                      >
                        参与挑战
                      </button>
                    ) : null;
                  })()}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {showCreateModal && (
        <CreateChallengeModal 
          onClose={() => setShowCreateModal(false)}
          onSuccess={(created) => {
            setShowCreateModal(false);
            if (created?.id) {
              // 先本地回显，避免用户觉得“没反应”
              setChallenges((prev) => {
                const next = [created, ...prev.filter((c) => c.id !== created.id)];
                return next;
              });
            }
            // 再从后端刷新一次，确保数据一致
            loadChallenges();
          }}
        />
      )}

      {selectedChallenge && (
        <ChallengeDetailModal
          challenge={selectedChallenge}
          onClose={() => setSelectedChallenge(null)}
          onJoin={() => {
            setSelectedChallenge(null);
            setJoinChallengeId(selectedChallenge.id);
          }}
        />
      )}

      {joinChallengeId && (
        <JoinChallengeModal
          challengeId={joinChallengeId}
          onClose={() => setJoinChallengeId(null)}
          onSuccess={() => {
            setJoinChallengeId(null);
            loadChallenges(); // 刷新列表
          }}
        />
      )}
    </div>
  );
};

// 创建挑战模态框
const CreateChallengeModal = ({ onClose, onSuccess }) => {
  const { publicKey } = useWalletContext();
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    challenge_type: 'individual',
    behavior_type: 'waste_sorting',
    reward_amount: 1000,
    target_count: 10,
    start_time: '',
    end_time: '',
    rules: '',
    image_url: '',
  });
  const [imageFile, setImageFile] = useState(null);
  const [imagePreview, setImagePreview] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const uploadImageIfNeeded = async (token) => {
    if (!imageFile) return '';
    const fd = new FormData();
    fd.append('file', imageFile);
    const res = await axios.post(`${API_BASE_URL}/api/v1/upload/image`, fd, {
      headers: { Authorization: `Bearer ${token}` },
    });
    const url = res.data?.url || '';
    if (!url) return '';
    // 后端返回相对路径时，补全为可直接访问的绝对地址（否则前端会去 3000 端口取图）
    if (typeof url === 'string' && url.startsWith('/')) return `${API_BASE_URL}${url}`;
    return url;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!publicKey) {
      setError('请先连接钱包');
      return;
    }

    setLoading(true);
    setError('');

    try {
      const token = localStorage.getItem('token');
      if (!token) {
        setError('请先登录');
        return;
      }

      // 基础表单校验（减少无效请求 + 更快反馈）
      if (!formData.title.trim() || !formData.description.trim()) {
        setError('请填写挑战标题与描述');
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

      let uploadedUrl = '';
      try {
        uploadedUrl = await uploadImageIfNeeded(token);
      } catch (e) {
        setError(e?.response?.data?.detail || e?.response?.data?.error || e?.message || '图片上传失败');
        return;
      }
      const res = await axios.post(
        `${API_BASE_URL}/api/v1/challenges`,
        { ...formData, image_url: uploadedUrl || formData.image_url },
        {
          headers: { Authorization: `Bearer ${token}` }
        }
      );

      onSuccess(res.data?.data);
    } catch (err) {
      setError(err.response?.data?.detail || err.response?.data?.error || '创建失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <h3>发起挑战</h3>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>挑战标题 *</label>
            <input
              type="text"
              value={formData.title}
              onChange={(e) => setFormData({...formData, title: e.target.value})}
              required
            />
          </div>
          <div className="form-group">
            <label>挑战描述 *</label>
            <textarea
              value={formData.description}
              onChange={(e) => setFormData({...formData, description: e.target.value})}
              required
            />
          </div>
          <div className="form-group">
            <label>挑战类型 *</label>
            <select
              value={formData.challenge_type}
              onChange={(e) => setFormData({...formData, challenge_type: e.target.value})}
            >
              <option value="individual">个人挑战</option>
              <option value="team">团队挑战</option>
              <option value="community">社区挑战</option>
            </select>
          </div>
          <div className="form-group">
            <label>行为类型 *</label>
            <select
              value={formData.behavior_type}
              onChange={(e) => setFormData({...formData, behavior_type: e.target.value})}
            >
              <option value="waste_sorting">垃圾分类</option>
              <option value="tree_planting">植树造林</option>
              <option value="low_carbon_travel">低碳出行</option>
            </select>
          </div>
          <div className="form-row">
            <div className="form-group">
              <label>奖励数量 *</label>
              <input
                type="number"
                value={formData.reward_amount}
                onChange={(e) => setFormData({...formData, reward_amount: parseInt(e.target.value)})}
                min="1"
                required
              />
            </div>
            <div className="form-group">
              <label>目标人数 *</label>
              <input
                type="number"
                value={formData.target_count}
                onChange={(e) => setFormData({...formData, target_count: parseInt(e.target.value)})}
                min="1"
                required
              />
            </div>
          </div>
          <div className="form-row">
            <div className="form-group">
              <label>开始时间 *</label>
              <input
                type="datetime-local"
                value={formData.start_time}
                onChange={(e) => setFormData({...formData, start_time: e.target.value})}
                required
              />
            </div>
            <div className="form-group">
              <label>结束时间 *</label>
              <input
                type="datetime-local"
                value={formData.end_time}
                onChange={(e) => setFormData({...formData, end_time: e.target.value})}
                required
              />
            </div>
          </div>

          <div className="form-group">
            <label>挑战规则（可选）</label>
            <textarea
              value={formData.rules}
              onChange={(e) => setFormData({ ...formData, rules: e.target.value })}
              placeholder="例如：完成指定次数的环保行为即可领取奖励"
            />
          </div>

          <div className="form-group">
            <label>背景图片（可选）</label>
            <input
              type="file"
              accept="image/*"
              onChange={(e) => {
                const file = e.target.files?.[0] || null;
                setImageFile(file);
                if (file) {
                  setImagePreview(URL.createObjectURL(file));
                } else {
                  setImagePreview('');
                }
              }}
            />
            <div style={{ marginTop: 6, fontSize: 12, color: '#64748b' }}>
              选择本地图片会自动上传；不选则使用默认背景。
            </div>
            {imagePreview && (
              <div style={{ marginTop: 8 }}>
                <img
                  src={imagePreview}
                  alt="预览"
                  style={{ width: '100%', maxHeight: 180, objectFit: 'cover', borderRadius: 8 }}
                />
              </div>
            )}
          </div>

          {error && <div className="error-message">{error}</div>}
          <div className="modal-actions">
            <button type="button" onClick={onClose} className="btn btn-secondary">
              取消
            </button>
            <button type="submit" disabled={loading} className="btn btn-primary">
              {loading ? '创建中...' : '创建挑战'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

// 挑战详情模态框
const formatRemaining = (endTime) => {
  if (!endTime) return '-';
  const diff = new Date(endTime).getTime() - Date.now();
  if (Number.isNaN(diff)) return '-';
  if (diff <= 0) return '已结束';
  const sec = Math.floor(diff / 1000);
  const d = Math.floor(sec / 86400);
  const h = Math.floor((sec % 86400) / 3600);
  const m = Math.floor((sec % 3600) / 60);
  return `${d}天 ${h}小时 ${m}分钟`;
};

const shortAddr = (addr) => {
  if (!addr) return '';
  if (addr.length <= 10) return addr;
  return `${addr.slice(0, 4)}...${addr.slice(-4)}`;
};

// 参与挑战模态框
const JoinChallengeModal = ({ challengeId, onClose, onSuccess }) => {
  const { publicKey } = useWalletContext();
  const [behaviors, setBehaviors] = useState([]);
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [selectedBehaviorId, setSelectedBehaviorId] = useState('');
  const [error, setError] = useState('');
  const [challengeInfo, setChallengeInfo] = useState(null);

  // 加载挑战信息和用户已通过审核的行为记录
  useEffect(() => {
    if (!challengeId || !publicKey) return;
    
    const loadData = async () => {
      setLoading(true);
      setError('');
      try {
        const token = localStorage.getItem('token');
        if (!token) {
          setError('请先登录');
          return;
        }

        // 获取挑战信息
        const challengeRes = await axios.get(`${API_BASE_URL}/api/v1/challenges/${challengeId}`);
        setChallengeInfo(challengeRes.data?.data);

        // 获取用户已通过审核的行为记录（匹配挑战的行为类型）
        const behaviorType = challengeRes.data?.data?.behavior_type;
        const behaviorsRes = await axios.get(
          `${API_BASE_URL}/api/v1/green/behaviors?status=approved${behaviorType ? `&behavior_type=${behaviorType}` : ''}`,
          { headers: { Authorization: `Bearer ${token}` } }
        );
        setBehaviors(behaviorsRes.data?.data || []);
        
        if (behaviorsRes.data?.data?.length === 0) {
          setError('您还没有已通过审核的环保行为记录，请先提交并等待审核通过');
        }
      } catch (err) {
        setError(err.response?.data?.error || err.response?.data?.detail || '加载失败');
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, [challengeId, publicKey]);

  const handleSubmit = async () => {
    if (!selectedBehaviorId) {
      setError('请选择一个环保行为记录');
      return;
    }

    setSubmitting(true);
    setError('');
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        setError('请先登录');
        return;
      }

      await axios.post(
        `${API_BASE_URL}/api/v1/challenges/${challengeId}/join`,
        { behavior_id: selectedBehaviorId },
        { headers: { Authorization: `Bearer ${token}` } }
      );

      onSuccess();
    } catch (err) {
      setError(err.response?.data?.error || err.response?.data?.detail || '参与失败');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()} style={{ maxWidth: '600px' }}>
        <h3>参与挑战</h3>
        
        {challengeInfo && (
          <div style={{ marginBottom: 16, padding: 12, background: '#f5f5f5', borderRadius: 4 }}>
            <p style={{ margin: 0, fontWeight: 'bold' }}>{challengeInfo.title}</p>
            <p style={{ margin: '4px 0 0 0', fontSize: '14px', color: '#666' }}>
              {challengeInfo.description}
            </p>
          </div>
        )}

        {loading ? (
          <div className="loading">加载中...</div>
        ) : error && !behaviors.length ? (
          <div className="error-message">{error}</div>
        ) : (
          <>
            <div className="form-group">
              <label>选择环保行为记录 *</label>
              <p style={{ fontSize: '12px', color: '#666', margin: '4px 0 8px 0' }}>
                请选择一条已通过审核的{challengeInfo?.behavior_type ? getBehaviorTypeText(challengeInfo.behavior_type) : '环保行为'}记录参与挑战
              </p>
              {behaviors.length === 0 ? (
                <div style={{ padding: 16, textAlign: 'center', color: '#999' }}>
                  暂无已通过审核的行为记录
                </div>
              ) : (
                <div style={{ maxHeight: '300px', overflowY: 'auto', border: '1px solid #ddd', borderRadius: 4 }}>
                  {behaviors.map((behavior) => (
                    <div
                      key={behavior.id}
                      onClick={() => setSelectedBehaviorId(behavior.id)}
                      style={{
                        padding: 12,
                        borderBottom: '1px solid #eee',
                        cursor: 'pointer',
                        background: selectedBehaviorId === behavior.id ? '#e3f2fd' : 'white',
                        transition: 'background 0.2s',
                      }}
                    >
                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <div>
                          <div style={{ fontWeight: 'bold' }}>
                            {getBehaviorTypeText(behavior.behavior_type)}
                          </div>
                          <div style={{ fontSize: '12px', color: '#666', marginTop: 4 }}>
                            {behavior.location || '未填写地点'} · {behavior.created_at ? new Date(behavior.created_at).toLocaleString('zh-CN') : '-'}
                          </div>
                          {behavior.reward_amount && (
                            <div style={{ fontSize: '12px', color: '#4caf50', marginTop: 4 }}>
                              奖励: {behavior.reward_amount} SOLGREEN
                            </div>
                          )}
                        </div>
                        {selectedBehaviorId === behavior.id && (
                          <span style={{ color: '#2196f3', fontSize: '20px' }}>✓</span>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>

            {error && <div className="error-message">{error}</div>}

            <div className="modal-actions">
              <button type="button" onClick={onClose} className="btn btn-secondary">
                取消
              </button>
              <button
                type="button"
                onClick={handleSubmit}
                disabled={submitting || !selectedBehaviorId || behaviors.length === 0}
                className="btn btn-primary"
              >
                {submitting ? '参与中...' : '确认参与'}
              </button>
            </div>
          </>
        )}
      </div>
    </div>
  );
};

const ChallengeDetailModal = ({ challenge, onClose, onJoin }) => {
  const [detailLoading, setDetailLoading] = useState(false);
  const [detailError, setDetailError] = useState('');
  const [detail, setDetail] = useState(null);

  useEffect(() => {
    if (!challenge?.id) return;
    let cancelled = false;
    (async () => {
      try {
        setDetailLoading(true);
        setDetailError('');
        const res = await axios.get(`${API_BASE_URL}/api/v1/challenges/${challenge.id}`);
        if (cancelled) return;
        setDetail(res.data);
      } catch (e) {
        if (cancelled) return;
        setDetailError(e?.response?.data?.error || e?.message || '加载详情失败');
      } finally {
        if (!cancelled) setDetailLoading(false);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [challenge?.id]);

  const data = detail?.data || challenge;
  const participantCount = detail?.participant_count;
  const participants = detail?.participants || [];

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content challenge-detail-modal" onClick={(e) => e.stopPropagation()}>
        {data?.image_url ? (
          <div className="challenge-detail-banner">
            <img
              src={data.image_url}
              alt={data.title}
              className="challenge-detail-banner-img"
            />
          </div>
        ) : null}

        <h3 className="challenge-detail-title">
          {data?.title || '挑战详情'}
        </h3>

        {detailLoading ? <div className="loading">加载中...</div> : null}
        {detailError ? <div className="error-message">{detailError}</div> : null}

        <p className="challenge-detail-desc">{data?.description}</p>

        <div className="challenge-info" style={{ marginTop: 12 }}>
          <div className="info-item">
            <span className="label">状态:</span>
            <span>{getStatusText(data?.status, data?.start_time, data?.end_time)}</span>
          </div>
          <div className="info-item">
            <span className="label">挑战类型:</span>
            <span>{getChallengeTypeText(data?.challenge_type)}</span>
          </div>
          <div className="info-item">
            <span className="label">类型:</span>
            <span>{getBehaviorTypeText(data?.behavior_type)}</span>
          </div>
          <div className="info-item">
            <span className="label">奖励:</span>
            <span className="reward">{data?.reward_amount} SOLGREEN</span>
          </div>
          <div className="info-item">
            <span className="label">倒计时:</span>
            <span>{formatRemaining(data?.end_time)}</span>
          </div>
          <div className="info-item">
            <span className="label">时间:</span>
            <span>
              {data?.start_time ? new Date(data.start_time).toLocaleString('zh-CN') : '-'}
              {' - '}
              {data?.end_time ? new Date(data.end_time).toLocaleString('zh-CN') : '-'}
            </span>
          </div>
          <div className="info-item">
            <span className="label">参与:</span>
            <span>
              {data?.current_count ?? 0} / {data?.target_count ?? 0}
              {typeof participantCount === 'number' ? `（参与者数：${participantCount}）` : ''}
            </span>
          </div>
        </div>

        {data?.rules ? (
          <div style={{ marginTop: 12 }}>
            <h4>挑战规则</h4>
            <p>{data.rules}</p>
          </div>
        ) : null}

        {participants.length > 0 ? (
          <div style={{ marginTop: 12 }}>
            <h4>最新参与者</h4>
            <ul style={{ paddingLeft: 18, margin: 0 }}>
              {participants.map((p) => (
                <li key={p.id}>
                  {shortAddr(p.wallet_addr)}（{p.joined_at ? new Date(p.joined_at).toLocaleString('zh-CN') : '-'}）
                </li>
              ))}
            </ul>
          </div>
        ) : null}

        <div style={{ display: 'flex', gap: 10, marginTop: 16 }}>
          <button
            type="button"
            className="btn btn-secondary"
            onClick={() => navigator.clipboard?.writeText?.(window.location.href)}
          >
            分享链接
          </button>
          {(() => {
            const now = Date.now();
            const start = data?.start_time ? new Date(data.start_time).getTime() : 0;
            const end = data?.end_time ? new Date(data.end_time).getTime() : 0;
            const canJoin = (data?.status === 'active' || data?.status === 'draft') && 
                           start > 0 && now >= start && now < end && onJoin;
            return canJoin ? (
              <button 
                className="btn btn-primary"
                onClick={onJoin}
              >
                参与挑战
              </button>
            ) : null;
          })()}
          <button onClick={onClose} className="btn btn-primary">关闭</button>
        </div>
      </div>
    </div>
  );
};
