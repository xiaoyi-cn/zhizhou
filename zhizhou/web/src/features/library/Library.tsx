import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import api from '../../lib/api';

interface Content {
  id: string;
  url: string;
  title: string;
  summary: string;
  category: string;
  tags: string[];
  created_at: string;
}

export default function Library() {
  const [contents, setContents] = useState<Content[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [selectedTag, setSelectedTag] = useState('');
  const [category, setCategory] = useState('');

  const fetchContents = async () => {
    setLoading(true);
    try {
      const params: any = { page, limit: 20 };
      if (selectedTag) params.tags = [selectedTag];
      if (category) params.category = category;
      const res: any = await api.get('/contents', { params });
      setContents(res.contents || []);
      setTotal(res.total || 0);
    } catch { /* ignore */ } finally {
      setLoading(false);
    }
  };

  useEffect(() => { fetchContents(); }, [page, selectedTag, category]);

  // 收集所有标签
  const allTags = [...new Set(contents.flatMap(c => c.tags || []))];

  return (
    <div>
      <h2 className="page-title">知识库</h2>

      <div style={{ display: 'flex', gap: 12, marginBottom: 24, flexWrap: 'wrap' }}>
        <select className="input" style={{ width: 'auto' }} value={category} onChange={e => { setCategory(e.target.value); setPage(1); }}>
          <option value="">全部分类</option>
          {[...new Set(contents.map(c => c.category).filter(Boolean))].map(cat => (
            <option key={cat} value={cat}>{cat}</option>
          ))}
        </select>
        <select className="input" style={{ width: 'auto' }} value={selectedTag} onChange={e => { setSelectedTag(e.target.value); setPage(1); }}>
          <option value="">全部标签</option>
          {allTags.map(tag => (
            <option key={tag} value={tag}>{tag}</option>
          ))}
        </select>
      </div>

      {loading ? (
        <div className="loading">加载中...</div>
      ) : contents.length === 0 ? (
        <div className="empty-state">
          <div style={{ fontSize: 48, marginBottom: 16 }}>📚</div>
          <h3>知识库还是空的</h3>
          <p>采集并审核内容后，它们会出现在这里</p>
        </div>
      ) : (
        <>
          <div className="grid">
            {contents.map(c => (
              <Link key={c.id} to={`/detail/${c.id}`} style={{ color: 'inherit' }}>
                <div className="card" style={{ cursor: 'pointer', transition: 'box-shadow 0.2s' }}
                  onMouseEnter={e => (e.currentTarget.style.boxShadow = 'var(--shadow-md)')}
                  onMouseLeave={e => (e.currentTarget.style.boxShadow = 'var(--shadow)')}
                >
                  <div style={{ fontWeight: 600, marginBottom: 8 }}>{c.title}</div>
                  {c.summary && (
                    <p style={{ fontSize: 14, color: 'var(--text-secondary)', marginBottom: 12, display: '-webkit-box', WebkitLineClamp: 2, WebkitBoxOrient: 'vertical', overflow: 'hidden' }}>
                      {c.summary}
                    </p>
                  )}
                  <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap' }}>
                    {c.category && <span className="badge badge-approved">{c.category}</span>}
                    {c.tags?.slice(0, 3).map(tag => (
                      <span key={tag} className="tag">{tag}</span>
                    ))}
                  </div>
                </div>
              </Link>
            ))}
          </div>

          {total > 20 && (
            <div style={{ display: 'flex', gap: 8, justifyContent: 'center', marginTop: 24 }}>
              <button className="btn btn-secondary btn-sm" disabled={page === 1} onClick={() => setPage(p => p - 1)}>上一页</button>
              <span style={{ lineHeight: '32px', fontSize: 14 }}>{page} / {Math.ceil(total / 20)}</span>
              <button className="btn btn-secondary btn-sm" disabled={page >= Math.ceil(total / 20)} onClick={() => setPage(p => p + 1)}>下一页</button>
            </div>
          )}
        </>
      )}
    </div>
  );
}