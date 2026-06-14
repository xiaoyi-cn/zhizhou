import { useState, useEffect } from 'react';
import api from '../../lib/api';

interface Content {
  id: string;
  url: string;
  title: string;
  summary: string;
  category: string;
  tags: string[];
  status: string;
  created_at: string;
}

export default function Review() {
  const [contents, setContents] = useState<Content[]>([]);
  const [loading, setLoading] = useState(true);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editForm, setEditForm] = useState({ title: '', summary: '', category: '', tags: '' });

  const fetchPending = async () => {
    try {
      const res: any = await api.get('/contents/pending');
      setContents(res.contents || []);
    } catch {
      // ignore
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { fetchPending(); }, []);

  const approve = async (id: string) => {
    try {
      await api.post(`/contents/${id}/approve`);
      setContents(prev => prev.filter(c => c.id !== id));
    } catch { /* ignore */ }
  };

  const skip = async (id: string) => {
    try {
      await api.post(`/contents/${id}/skip`);
      setContents(prev => prev.filter(c => c.id !== id));
    } catch { /* ignore */ }
  };

  const startEdit = (c: Content) => {
    setEditingId(c.id);
    setEditForm({
      title: c.title,
      summary: c.summary,
      category: c.category,
      tags: c.tags.join(', '),
    });
  };

  const saveEdit = async (id: string) => {
    try {
      await api.put(`/contents/${id}`, {
        title: editForm.title,
        summary: editForm.summary,
        category: editForm.category,
        tags: editForm.tags.split(',').map(t => t.trim()).filter(Boolean),
      });
      setEditingId(null);
      fetchPending();
    } catch { /* ignore */ }
  };

  return (
    <div>
      <div className="flex-between" style={{ marginBottom: 24 }}>
        <h2 className="page-title" style={{ marginBottom: 0 }}>审核</h2>
        <span className="badge badge-pending">{contents.length} 条待审核</span>
      </div>

      {loading ? (
        <div className="loading">加载中...</div>
      ) : contents.length === 0 ? (
        <div className="empty-state">
          <div style={{ fontSize: 48, marginBottom: 16 }}>🎉</div>
          <h3>没有待审核内容</h3>
          <p>去采集页添加新内容吧</p>
        </div>
      ) : (
        <div className="grid" style={{ gap: 12 }}>
          {contents.map(c => (
            <div key={c.id} className="card">
              {editingId === c.id ? (
                <div>
                  <div className="form-group">
                    <label className="label">标题</label>
                    <input className="input" value={editForm.title} onChange={e => setEditForm({ ...editForm, title: e.target.value })} />
                  </div>
                  <div className="form-group">
                    <label className="label">摘要</label>
                    <textarea className="input" rows={3} value={editForm.summary} onChange={e => setEditForm({ ...editForm, summary: e.target.value })} />
                  </div>
                  <div className="form-group">
                    <label className="label">分类</label>
                    <input className="input" value={editForm.category} onChange={e => setEditForm({ ...editForm, category: e.target.value })} />
                  </div>
                  <div className="form-group">
                    <label className="label">标签（逗号分隔）</label>
                    <input className="input" value={editForm.tags} onChange={e => setEditForm({ ...editForm, tags: e.target.value })} />
                  </div>
                  <div style={{ display: 'flex', gap: 8 }}>
                    <button className="btn btn-primary btn-sm" onClick={() => saveEdit(c.id)}>保存</button>
                    <button className="btn btn-secondary btn-sm" onClick={() => setEditingId(null)}>取消</button>
                  </div>
                </div>
              ) : (
                <div>
                  <div style={{ marginBottom: 12 }}>
                    <div style={{ fontWeight: 600, fontSize: 16, marginBottom: 4 }}>{c.title}</div>
                    {c.url && (
                      <a href={c.url} target="_blank" rel="noopener noreferrer" style={{ fontSize: 13, wordBreak: 'break-all' }}>
                        {c.url}
                      </a>
                    )}
                  </div>

                  {c.summary && (
                    <p style={{ fontSize: 14, color: 'var(--text-secondary)', marginBottom: 12, lineHeight: 1.7 }}>
                      {c.summary}
                    </p>
                  )}

                  <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap', marginBottom: 16 }}>
                    {c.category && <span className="badge badge-approved">{c.category}</span>}
                    {c.tags?.map(tag => (
                      <span key={tag} className="tag">{tag}</span>
                    ))}
                  </div>

                  <div style={{ display: 'flex', gap: 8 }}>
                    <button className="btn btn-success btn-sm" onClick={() => approve(c.id)}>确认归档</button>
                    <button className="btn btn-secondary btn-sm" onClick={() => skip(c.id)}>跳过</button>
                    <button className="btn btn-secondary btn-sm" onClick={() => startEdit(c)}>编辑</button>
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}