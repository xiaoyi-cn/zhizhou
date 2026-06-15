import { useState } from 'react';
import api from '../../lib/api';

export default function Ingest() {
  const [url, setUrl] = useState('');
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<{ id: string; status: string } | null>(null);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!url.trim()) return;

    setLoading(true);
    setError('');
    setResult(null);
    try {
      const res: any = await api.post('/contents/ingest', { url });
      setResult(res);
      setUrl('');
    } catch {
      setError('提交失败，请检查链接是否有效');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h2 className="page-title">内容采集</h2>
      <p style={{ color: 'var(--text-secondary)', marginBottom: 24 }}>
        粘贴文章链接，知舟会自动抓取内容并生成 AI 摘要
      </p>

      <div className="card" style={{ marginBottom: 24 }}>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label className="label">文章链接</label>
            <input
              className="input"
              type="url"
              placeholder="https://..."
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              required
              style={{ fontSize: 16, padding: '12px 16px' }}
            />
          </div>
          {error && (
            <div style={{ color: 'var(--danger)', marginBottom: 16, fontSize: 14 }}>{error}</div>
          )}
          <button className="btn btn-primary" type="submit" disabled={loading}>
            {loading ? '处理中...' : '开始采集'}
          </button>
        </form>
      </div>

      {result && (
        <div className="card" style={{ borderLeft: '4px solid var(--success)' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <span style={{ fontSize: 20 }}>✅</span>
            <div>
              <div style={{ fontWeight: 600 }}>提交成功</div>
              <div style={{ color: 'var(--text-secondary)', fontSize: 14 }}>
                内容已进入待审核列表，AI 正在自动生成摘要和标签...
              </div>
            </div>
          </div>
        </div>
      )}

      <div style={{ marginTop: 40 }}>
        <h3 style={{ fontSize: 16, marginBottom: 12 }}>支持的采集方式</h3>
        <div className="grid" style={{ gridTemplateColumns: 'repeat(auto-fill, minmax(200px, 1fr))' }}>
          {['粘贴链接', '浏览器扩展', '移动端分享', '截图 OCR'].map(item => (
            <div key={item} className="card" style={{ textAlign: 'center' }}>
              <div style={{ fontSize: 24, marginBottom: 8 }}>📌</div>
              <div style={{ fontSize: 14, fontWeight: 500 }}>{item}</div>
              <div style={{ fontSize: 12, color: 'var(--text-muted)', marginTop: 4 }}>
                {item === '粘贴链接' ? 'Phase 1 已支持' : 'Phase 2-3 规划中'}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}