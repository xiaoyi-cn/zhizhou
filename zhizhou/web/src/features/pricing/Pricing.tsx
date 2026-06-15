import { useState, useEffect } from 'react';
import api from '../../lib/api';

interface FeatureInfo {
  tier: string;
  features: { [key: string]: boolean };
}

const plans = [
  {
    name: '免费版',
    tier: 'free',
    price: '永久免费',
    features: [
      { name: '链接采集 + AI 摘要/标签', free: true },
      { name: '审核确认 → 归档', free: true },
      { name: '全文搜索 + 标签筛选', free: true },
      { name: 'API Key 自配', free: true },
      { name: '云端存储', free: '10GB' },
      { name: '本地存储', free: '完整功能' },
      { name: 'MCP Server', free: false },
      { name: '开放 API', free: false },
      { name: '自动化规则引擎', free: false },
    ],
  },
  {
    name: '专业版',
    tier: 'pro',
    price: '订阅制',
    recommended: true,
    features: [
      { name: '全部免费版功能', free: true },
      { name: '云端存储', free: '无限制' },
      { name: 'MCP Server', free: true },
      { name: '开放 API', free: true },
      { name: '自动化规则引擎', free: true },
      { name: '优先技术支持', free: true },
    ],
  },
];

export default function Pricing() {
  const [features, setFeatures] = useState<FeatureInfo | null>(null);

  useEffect(() => {
    api.get('/features').then((res: any) => setFeatures(res)).catch(() => {});
  }, []);

  return (
    <div>
      <h2 className="page-title">版本对比</h2>

      {features && (
        <div className="card" style={{ marginBottom: 24, textAlign: 'center' }}>
          <div style={{ fontSize: 14, color: 'var(--text-secondary)' }}>
            当前版本：<strong>{features.tier === 'pro' ? '专业版' : '免费版'}</strong>
          </div>
        </div>
      )}

      <div className="grid" style={{ gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))' }}>
        {plans.map(plan => (
          <div key={plan.tier} className="card" style={{
            border: plan.recommended ? '2px solid var(--primary)' : '1px solid var(--border)',
            position: 'relative',
          }}>
            {plan.recommended && (
              <div style={{
                position: 'absolute', top: -12, left: '50%', transform: 'translateX(-50%)',
                background: 'var(--primary)', color: 'white', padding: '4px 16px',
                borderRadius: 20, fontSize: 12, fontWeight: 600,
              }}>
                推荐
              </div>
            )}

            <div style={{ textAlign: 'center', marginBottom: 24, marginTop: 8 }}>
              <h3 style={{ fontSize: 20, fontWeight: 700 }}>{plan.name}</h3>
              <div style={{ fontSize: 28, fontWeight: 700, color: 'var(--primary)', marginTop: 8 }}>
                {plan.price}
              </div>
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
              {plan.features.map(f => (
                <div key={f.name} style={{ display: 'flex', alignItems: 'center', gap: 8, fontSize: 14 }}>
                  {f.free === true ? (
                    <span style={{ color: 'var(--success)', fontWeight: 600 }}>✓</span>
                  ) : f.free === false ? (
                    <span style={{ color: 'var(--text-muted)' }}>✗</span>
                  ) : (
                    <span style={{ color: 'var(--primary)', fontWeight: 500 }}>{f.free}</span>
                  )}
                  <span style={{ color: f.free === false ? 'var(--text-muted)' : 'var(--text)' }}>{f.name}</span>
                </div>
              ))}
            </div>

            {plan.recommended && (
              <button className="btn btn-primary" style={{ width: '100%', marginTop: 24 }}>升级专业版</button>
            )}
          </div>
        ))}
      </div>

      <div className="card" style={{ marginTop: 24 }}>
        <h3 style={{ fontSize: 16, fontWeight: 600, marginBottom: 12 }}>设计原则</h3>
        <ul style={{ fontSize: 14, color: 'var(--text-secondary)', lineHeight: 2, paddingLeft: 20 }}>
          <li>不歧视免费用户：免费版拥有完整的核心体验</li>
          <li>专业版 = 服务器成本 + 开发者工具</li>
          <li>同步更新，不分先后</li>
          <li>本地部署永远免费且完整</li>
          <li>不做商业焦虑</li>
        </ul>
      </div>
    </div>
  );
}