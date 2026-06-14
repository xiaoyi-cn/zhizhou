import { useState, useEffect } from 'react';
import api from '../../lib/api';

interface APIKey {
  id: string;
  provider: string;
  base_url: string;
  model: string;
  is_active: boolean;
  created_at: string;
}

interface Category {
  id: string;
  name: string;
  parent_id: string | null;
}

export default function Settings() {
  const [activeTab, setActiveTab] = useState<'keys' | 'categories'>('keys');

  // API Key
  const [keys, setKeys] = useState<APIKey[]>([]);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ provider: 'deepseek', key: '', base_url: '', model: 'deepseek-chat' });
  const [testing, setTesting] = useState(false);
  const [testResult, setTestResult] = useState('');

  // Categories
  const [categories, setCategories] = useState<Category[]>([]);
  const [newCategory, setNewCategory] = useState('');

  useEffect(() => {
    if (activeTab === 'keys') fetchKeys();
    if (activeTab === 'categories') fetchCategories();
  }, [activeTab]);

  const fetchKeys = async () => {
    try {
      const res: any = await api.get('/keys');
      setKeys(res.keys || []);
    } catch { /* ignore */ }
  };

  const fetchCategories = async () => {
    try {
      const res: any = await api.get('/categories');
      setCategories(res.categories || []);
    } catch { /* ignore */ }
  };

  const saveKey = async () => {
    try {
      await api.post('/keys', form);
      setShowForm(false);
      setForm({ provider: 'deepseek', key: '', base_url: '', model: 'deepseek-chat' });
      fetchKeys();
    } catch { /* ignore */ }
  };

  const deleteKey = async (id: string) => {
    try {
      await api.delete(`/keys/${id}`);
      fetchKeys();
    } catch { /* ignore */ }
  };

  const testKey = async () => {
    setTesting(true);
    setTestResult('');
    try {
      await api.post('/keys/:id/test', form);
      setTestResult('连接成功');
    } catch {
      setTestResult('连接失败');
    } finally {
      setTesting(false);
    }
  };

  const addCategory = async () => {
    if (!newCategory.trim()) return;
    try {
      await api.post('/categories', { name: newCategory });
      setNewCategory('');
      fetchCategories();
    } catch { /* ignore */ }
  };

  const deleteCategory = async (id: string) => {
    try {
      await api.delete(`/categories/${id}`);
      fetchCategories();
    } catch { /* ignore */ }
  };

  return (
    <div>
      <h2 className="page-title">设置</h2>

      <div style={{ display: 'flex', gap: 4, marginBottom: 24, borderBottom: '1px solid var(--border)' }}>
        {[
          { key: 'keys' as const, label: 'API Key' },
          { key: 'categories' as const, label: '分类管理' },
        ].map(tab => (
          <button
            key={tab.key}
            onClick={() => setActiveTab(tab.key)}
            style={{
              padding: '10px 20px',
              border: 'none',
              background: 'none',
              fontSize: 14,
              fontWeight: 500,
              color: activeTab === tab.key ? 'var(--primary)' : 'var(--text-secondary)',
              borderBottom: activeTab === tab.key ? '2px solid var(--primary)' : '2px solid transparent',
              marginBottom: -1,
            }}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {activeTab === 'keys' && (
        <div>
          <p style={{ color: 'var(--text-secondary)', marginBottom: 16, fontSize: 14 }}>
            绑定你的 LLM API Key，知舟不提供 AI 服务，所有 AI 处理走你的账户。
          </p>

          {keys.map(k => (
            <div key={k.id} className="card" style={{ marginBottom: 12 }}>
              <div className="flex-between">
                <div>
                  <div style={{ fontWeight: 600 }}>{k.provider}</div>
                  <div style={{ fontSize: 13, color: 'var(--text-muted)' }}>模型: {k.model}</div>
                  {k.base_url && <div style={{ fontSize: 13, color: 'var(--text-muted)' }}>{k.base_url}</div>}
                </div>
                <button className="btn btn-danger btn-sm" onClick={() => deleteKey(k.id)}>删除</button>
              </div>
            </div>
          ))}

          {!showForm ? (
            <button className="btn btn-primary" onClick={() => setShowForm(true)}>添加 API Key</button>
          ) : (
            <div className="card" style={{ marginTop: 16 }}>
              <div className="form-group">
                <label className="label">服务商</label>
                <select className="input" value={form.provider} onChange={e => setForm({ ...form, provider: e.target.value })}>
                  <option value="deepseek">DeepSeek</option>
                  <option value="openai">OpenAI</option>
                  <option value="hunyuan">混元</option>
                  <option value="custom">自定义</option>
                </select>
              </div>
              <div className="form-group">
                <label className="label">API Key</label>
                <input className="input" type="password" value={form.key} onChange={e => setForm({ ...form, key: e.target.value })} placeholder="sk-..." />
              </div>
              <div className="form-group">
                <label className="label">Base URL（可选）</label>
                <input className="input" value={form.base_url} onChange={e => setForm({ ...form, base_url: e.target.value })} placeholder="https://api.deepseek.com/v1" />
              </div>
              <div className="form-group">
                <label className="label">模型</label>
                <input className="input" value={form.model} onChange={e => setForm({ ...form, model: e.target.value })} placeholder="deepseek-chat" />
              </div>
              {testResult && (
                <div style={{ marginBottom: 12, fontSize: 14, color: testResult === '连接成功' ? 'var(--success)' : 'var(--danger)' }}>
                  {testResult}
                </div>
              )}
              <div style={{ display: 'flex', gap: 8 }}>
                <button className="btn btn-primary btn-sm" onClick={saveKey} disabled={!form.key}>保存</button>
                <button className="btn btn-secondary btn-sm" onClick={testKey} disabled={testing || !form.key}>
                  {testing ? '测试中...' : '测试连接'}
                </button>
                <button className="btn btn-secondary btn-sm" onClick={() => setShowForm(false)}>取消</button>
              </div>
            </div>
          )}
        </div>
      )}

      {activeTab === 'categories' && (
        <div>
          <div style={{ display: 'flex', gap: 8, marginBottom: 16 }}>
            <input
              className="input"
              placeholder="新分类名，如 技术/后端"
              value={newCategory}
              onChange={e => setNewCategory(e.target.value)}
              onKeyDown={e => e.key === 'Enter' && addCategory()}
            />
            <button className="btn btn-primary" onClick={addCategory}>添加</button>
          </div>

          {categories.length === 0 ? (
            <div className="empty-state">
              <h3>还没有分类</h3>
              <p>添加分类后，AI 会自动将内容归入对应分类</p>
            </div>
          ) : (
            <div className="grid">
              {categories.map(c => (
                <div key={c.id} className="card">
                  <div className="flex-between">
                    <span style={{ fontWeight: 500 }}>{c.name}</span>
                    <button className="btn btn-danger btn-sm" onClick={() => deleteCategory(c.id)}>删除</button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  );
}