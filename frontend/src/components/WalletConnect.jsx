import { WalletAdapterNetwork } from '@solana/wallet-adapter-base';
import { ConnectionProvider, WalletProvider, useWallet } from '@solana/wallet-adapter-react';
import { WalletModalProvider, WalletMultiButton } from '@solana/wallet-adapter-react-ui';
import { PhantomWalletAdapter, SolflareWalletAdapter } from '@solana/wallet-adapter-wallets';
import { clusterApiUrl } from '@solana/web3.js';
import axios from 'axios';
import { useEffect, useMemo, useRef, useState } from 'react';
import '@solana/wallet-adapter-react-ui/styles.css';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

function bytesToHex(bytes) {
  return Array.from(bytes)
    .map((b) => b.toString(16).padStart(2, '0'))
    .join('');
}

// 顶层钱包上下文 Provider，应该包裹整个应用
export const WalletContextProvider = ({ children }) => {
  const network = WalletAdapterNetwork.Devnet; // 测试网，上线改 Mainnet
  const endpoint = useMemo(() => clusterApiUrl(network), [network]);
  const wallets = useMemo(() => {
    const list = [new PhantomWalletAdapter(), new SolflareWalletAdapter()];
    // React 需要 key 唯一；以 adapter.name 去重，避免出现重复钱包名称导致的 warning
    const seen = new Set();
    return list.filter((w) => {
      const name = w?.name?.toString?.() || String(w?.name || '');
      if (!name) return true;
      if (seen.has(name)) return false;
      seen.add(name);
      return true;
    });
  }, []);

  return (
    <ConnectionProvider endpoint={endpoint}>
      <WalletProvider
        wallets={wallets}
        autoConnect={false}
        onError={(e) => {
          // 连接/授权/弹窗失败时，把错误打到控制台，便于定位“点了没反应”
          // eslint-disable-next-line no-console
          console.error('Wallet error:', e);
        }}
      >
        <WalletModalProvider>
          {children}
        </WalletModalProvider>
      </WalletProvider>
    </ConnectionProvider>
  );
};

// 只渲染按钮，本身不再创建 Provider
export const WalletConnect = () => {
  const { publicKey, connected, signMessage } = useWallet();
  const [authMsg, setAuthMsg] = useState('');
  const [authLoading, setAuthLoading] = useState(false);
  const lastAuthedWalletRef = useRef('');

  useEffect(() => {
    if (!connected) {
      setAuthMsg('');
      setAuthLoading(false);
      lastAuthedWalletRef.current = '';
    }
  }, [connected]);

  const walletAddr = publicKey?.toString() || '';
  const token = localStorage.getItem('token');
  const tokenWallet = localStorage.getItem('token_wallet_addr');
  const isAuthed = Boolean(token && tokenWallet === walletAddr);

  const handleLogin = async () => {
    if (!connected || !publicKey) {
      setAuthMsg('请先连接钱包');
      return;
    }
    if (!signMessage) {
      setAuthMsg('钱包不支持签名登录（signMessage 不可用）');
      return;
    }
    if (authLoading) return;

    try {
      setAuthLoading(true);
      setAuthMsg('请在钱包弹窗中确认签名...');

      const message = `Sol-Green 登录验证\n钱包地址: ${walletAddr}\n时间戳: ${Date.now()}`;
      const signatureBytes = await signMessage(new TextEncoder().encode(message));

      const res = await axios.post(`${API_BASE_URL}/api/v1/auth/wallet`, {
        wallet_addr: walletAddr,
        signature: bytesToHex(signatureBytes),
        message,
      });

      localStorage.setItem('token', res.data.token);
      localStorage.setItem('token_wallet_addr', walletAddr);
      lastAuthedWalletRef.current = walletAddr;
      setAuthMsg('已登录');
    } catch (e) {
      const raw = e?.response?.data?.error || e?.message || '登录失败';
      // 常见：Phantom/Solflare 扩展的 keyring 报错（未解锁/权限异常/请求被拒绝）
      if (String(raw).toLowerCase().includes('keyring request')) {
        setAuthMsg('钱包扩展 Keyring 异常：请打开钱包扩展解锁后重试；如仍失败，断开站点连接后重新连接。');
      } else {
        setAuthMsg(raw);
      }
    } finally {
      setAuthLoading(false);
    }
  };

  return (
    <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
      <WalletMultiButton />
      {connected && !isAuthed ? (
        <button
          onClick={handleLogin}
          disabled={authLoading}
          style={{
            padding: '8px 10px',
            borderRadius: 8,
            border: '1px solid rgba(255,255,255,0.2)',
            background: 'rgba(255,255,255,0.08)',
            color: 'inherit',
            cursor: authLoading ? 'not-allowed' : 'pointer',
            fontSize: 12,
          }}
        >
          {authLoading ? '登录中...' : '登录'}
        </button>
      ) : null}

      {connected && isAuthed ? <span style={{ fontSize: 12, opacity: 0.8 }}>已登录</span> : null}
      {authMsg ? <span style={{ fontSize: 12, opacity: 0.8 }}>{authMsg}</span> : null}
    </div>
  );
};

export const useWalletContext = () => {
  return useWallet();
};
