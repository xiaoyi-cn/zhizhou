import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../../lib/api';
import { useAuthStore } from '../../store/auth';

export default function Login() {
  const [phone, setPhone] = useState('');
  const [code, setCode] = useState('');
  const [step, setStep] = useState<'phone' | 'code'>('phone');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const navigate = useNavigate();
  const { setAuth } = useAuthStore();

  const sendCode = async () => {
    if (!/^1\d{10}$/.test(phone)) {
      setError('请输入正确的手机号');
      return;
    }
    setLoading(true);
    setError('');
    try {
      await api.post('/auth/send-code', { phone });
      setStep('code');
    } catch {
      setError('发送验证码失败');
    } finally {
      setLoading(false);
    }
  };

  const verifyCode = async () => {
    if (!code) {
      setError('请输入验证码');
      return;
    }
    setLoading(true);
    setError('');
    try {
      const res: any = await api.post('/auth/verify-code', { phone, code });
      setAuth(res.token, res.user_id);
      navigate('/ingest');
    } catch {
      setError('验证失败，请检查验证码');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{
      minHeight: '100vh',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    }}>
      <div className="card" style={{ width: 400, padding: 40 }}>
        <div style={{ textAlign: 'center', marginBottom: 32 }}>
          <h1 style={{ fontSize: 28, fontWeight: 700, color: 'var(--primary)' }}>知舟</h1>
          <p style={{ color: 'var(--text-secondary)', marginTop: 8 }}>个人知识管理工具</p>
        </div>

        {error && (
          <div style={{
            padding: '10px 14px',
            background: '#FEF2F2',
            color: 'var(--danger)',
            borderRadius: 'var(--radius)',
            marginBottom: 16,
            fontSize: 14,
          }}>
            {error}
          </div>
        )}

        {step === 'phone' ? (
          <>
            <div className="form-group">
              <label className="label">手机号</label>
              <input
                className="input"
                type="tel"
                placeholder="请输入手机号"
                value={phone}
                onChange={(e) => setPhone(e.target.value)}
                maxLength={11}
              />
            </div>
            <button className="btn btn-primary" style={{ width: '100%' }} onClick={sendCode} disabled={loading}>
              {loading ? '发送中...' : '获取验证码'}
            </button>
          </>
        ) : (
          <>
            <div className="form-group">
              <label className="label">验证码</label>
              <input
                className="input"
                type="text"
                placeholder="请输入验证码"
                value={code}
                onChange={(e) => setCode(e.target.value)}
                maxLength={6}
              />
            </div>
            <button className="btn btn-primary" style={{ width: '100%' }} onClick={verifyCode} disabled={loading}>
              {loading ? '验证中...' : '登录'}
            </button>
            <button
              className="btn btn-secondary"
              style={{ width: '100%', marginTop: 12 }}
              onClick={() => { setStep('phone'); setError(''); }}
            >
              返回
            </button>
          </>
        )}

        <p style={{ textAlign: 'center', marginTop: 24, fontSize: 12, color: 'var(--text-muted)' }}>
          登录即表示同意服务条款和隐私政策
        </p>
      </div>
    </div>
  );
}