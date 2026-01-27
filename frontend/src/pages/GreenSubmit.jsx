import { useState } from 'react';
import { useWalletContext } from '../components/WalletConnect';
import axios from 'axios';
import './GreenSubmit.css';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

export const GreenSubmit = () => {
  const { publicKey, signMessage } = useWalletContext();
  const [behaviorType, setBehaviorType] = useState('waste_sorting');
  const [mediaFiles, setMediaFiles] = useState([]);
  const [location, setLocation] = useState('');
  const [loading, setLoading] = useState(false);
  const [msg, setMsg] = useState('');
  const [result, setResult] = useState(null);

  // 处理文件选择
  const handleFileChange = (e) => {
    const files = Array.from(e.target.files);
    setMediaFiles(files);
  };

  // 上传文件到服务器（简化版，实际应上传到 OSS/IPFS）
  const uploadFiles = async (files) => {
    // 实际项目中应调用文件上传接口
    // 这里返回模拟 URL
    return files.map((file, idx) => 
      `https://sol-green.oss-cn-shanghai.aliyuncs.com/${publicKey?.toString()}-${Date.now()}-${idx}.${file.name.split('.').pop()}`
    );
  };

  // 钱包登录
  const handleWalletLogin = async () => {
    if (!publicKey || !signMessage) {
      setMsg('请先连接钱包');
      return;
    }

    try {
      const message = `Sol-Green 登录验证\n钱包地址: ${publicKey.toString()}\n时间戳: ${Date.now()}`;
      const signature = await signMessage(new TextEncoder().encode(message));
      
      const res = await axios.post(`${API_BASE_URL}/api/v1/auth/wallet`, {
        wallet_addr: publicKey.toString(),
        signature: Array.from(signature).map(b => b.toString(16).padStart(2, '0')).join(''),
        message: message,
      });

      localStorage.setItem('token', res.data.token);
      setMsg('登录成功');
    } catch (err) {
      setMsg(err.response?.data?.error || '登录失败');
      console.error(err);
    }
  };

  // 提交行为申报
  const handleSubmit = async () => {
    if (!publicKey) {
      setMsg('请先连接钱包');
      return;
    }
    if (mediaFiles.length === 0) {
      setMsg('请上传至少一张图片/视频');
      return;
    }

    // 检查是否已登录
    const token = localStorage.getItem('token');
    if (!token) {
      await handleWalletLogin();
      return;
    }

    setLoading(true);
    setMsg('');
    setResult(null);

    try {
      // 1. 上传文件（实际应调用上传接口）
      const mediaURLs = await uploadFiles(mediaFiles);
      
      // 2. 调用后端接口
      const res = await axios.post(
        `${API_BASE_URL}/api/v1/green/behavior/submit`,
        {
          behavior_type: behaviorType,
          media_urls: mediaURLs,
          location: location || '未指定',
        },
        {
          headers: { Authorization: `Bearer ${token}` }
        }
      );

      setResult(res.data);
      setMsg(res.data.msg || '提交成功');
      
      // 清空表单
      setMediaFiles([]);
      setLocation('');
    } catch (err) {
      setMsg(err.response?.data?.error || '提交失败');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="green-submit">
      <div className="card">
        <h2>🌱 环保行为申报 / Environmental Behavior Submission</h2>
        
        {!publicKey && (
          <div className="alert">
            <p>请先连接钱包 / Please connect wallet first</p>
          </div>
        )}

        <div className="form-group">
          <label className="label">行为类型 / Behavior Type</label>
          <select 
            value={behaviorType} 
            onChange={(e) => setBehaviorType(e.target.value)}
            className="input"
          >
            <option value="waste_sorting">🗑️ 垃圾分类 / Waste Sorting</option>
            <option value="tree_planting">🌳 植树造林 / Tree Planting</option>
            <option value="low_carbon_travel">🚴 低碳出行 / Low-Carbon Travel</option>
          </select>
        </div>

        <div className="form-group">
          <label className="label">行为地点（可选）/ Location (Optional)</label>
          <input
            type="text"
            value={location}
            onChange={(e) => setLocation(e.target.value)}
            placeholder="例如：上海张江软件园 / e.g., Shanghai Zhangjiang Software Park"
            className="input"
          />
        </div>

        <div className="form-group">
          <label className="label">上传证明（图片/视频）/ Upload Proof (Images/Videos)</label>
          <input 
            type="file" 
            multiple 
            accept="image/*,video/*"
            onChange={handleFileChange}
            className="input"
          />
          {mediaFiles.length > 0 && (
            <p className="file-info">已选择 {mediaFiles.length} 个文件</p>
          )}
        </div>

        <button 
          onClick={handleSubmit}
          disabled={loading || !publicKey || mediaFiles.length === 0}
          className="btn btn-primary"
        >
          {loading ? '提交中... / Submitting...' : '提交申报 / Submit'}
        </button>

        {msg && (
          <div className={`message ${msg.includes('成功') ? 'success' : 'error'}`}>
            {msg}
          </div>
        )}

        {result && (
          <div className="result">
            <h3>提交结果</h3>
            <p><strong>行为ID:</strong> {result.id}</p>
            <p><strong>状态:</strong> {result.status}</p>
            {result.tx_hash && (
              <p>
                <strong>交易哈希:</strong>{' '}
                <a 
                  href={`https://explorer.solana.com/tx/${result.tx_hash}?cluster=devnet`}
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  {result.tx_hash}
                </a>
              </p>
            )}
            {result.proof && (
              <p><strong>存证哈希:</strong> {result.proof}</p>
            )}
          </div>
        )}
      </div>
    </div>
  );
};
