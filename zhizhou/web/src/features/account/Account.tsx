import { useState, useEffect } from 'react';
import api from '../../lib/api';

interface SubscriptionInfo {
  tier: string;
  storage_used: number;
  storage_limit: number;
  pro_expires_at: string | null;
  features: { [key: string]: boolean };
}

export default function Account() {
  const [info, setInfo] = useState<SubscriptionInfo | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get('/subscription')
      .then((res: any) => setInfo(res))
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '无限制';
    const gb = bytes / (1024 * 1024 * 1024);
    return gb >= 1 ? `${gb.toFixed(1)} GB` : `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  };

  if (loading) return <div className="loading">加载中...</div>;

  return (
    <div>
      <h2 className="page-title">账户</h2>

      {info && (
        <>
          <div className="card" style={{ marginBottom: 24 }}>
            <h3 style={{ fontSize: 16, fontWeight: 600, marginBottom: 16 }}>订阅状态</h3>
            <div className="grid" style={{ gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))' }}>
              <div>
                <div style={{ fontSize: 13, color: 'var(--text-muted)' }}>当前版本</div>
                <div style={{ fontSize: 18, fontWeight: 700, color: info.tier === 'pro' ? 'var(--primary)' : 'var(--text)' }}>
                  {info.tier === 'pro' ? '专业版' : '免费版'}
                </div>
              </div>
              <div>
                <div style={{ fontSize: 13, color: 'var(--text-muted)' }}>存储用量</div>
                <div style={{ fontSize: 18, fontWeight: 700 }}>
                  {formatBytes(info.storage_used)}
                  {info.storage_limit > 0 && (
                    <span style={{ fontSize: 14, fontWeight: 400, color: 'var(--text-muted)' }}> / {formatBytes(info.storage_limit)}</span>
                  )}
                </div>
              </div>
              {info.pro_expires_at && (
                <div>
                  <div style={{ fontSize: 13, color: 'var(--text-muted)' }}>到期时间</div>
                  <div style={{ fontSize: 18, fontWeight: 700 }}>
                    {new Date(info.pro_expires_at).toLocaleDateString('zh-CN')}
                  </div>
                </div>
              )}
            </div>
          </div>

          <div className="card" style={{ marginBottom: 24 }}>
            <h3 style={{ fontSize: 16, fontWeight: 600, marginBottom: 16 }}>功能权限</h3>
            <div className="grid" style={{ gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))' }}>
              {Object.entries(info.features).map(([key, enabled]) => (
                <div key={key} style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                  <span style={{
                    width: 8, height: 8, borderRadius: '50%',
                    background: enabled ? 'var(--success)' : 'var(--text-muted)',
                  }} />
                  <span style={{ fontSize: 14 }}>
                    {key === 'mcp_server' ? 'MCP Server' :
                     key === 'open_api' ? '开放 API' :
                     key === 'auto_rules' ? '自动化规则' : key}
                  </span>
                  <span style={{ fontSize: 12, color: enabled ? 'var(--success)' : 'var(--text-muted)' }}>
                    {enabled ? '已开启' : '专业版'}
                  </span>
                </div>
              ))}
            </div>
          </div>

          <div className="card">
            <h3 style={{ fontSize: 16, fontWeight: 600, marginBottom: 16 }}>操作</h3>
            <div style={{ display: 'flex', gap: 8 }}>
              <button className="btn btn-primary btn-sm">升级专业版</button>
              <button className="btn btn-secondary btn-sm">导出数据</button>
            </div>
          </div>
        </>
      )}
    </div>
  );
}