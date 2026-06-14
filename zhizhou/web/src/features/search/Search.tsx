import { useState } from 'react';
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

export default function Search() {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<Content[]>([]);
  const [loading, setLoading] = useState(false);
  const [searched, setSearched] = useState(false);
  const [total, setTotal] = useState(0);

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!query.trim()) return;

    setLoading(true);
    setSearched(true);
    try {
      const res: any = await api.get('/search', { params: { q: query } });
      setResults(res.contents || []);
      setTotal(res.total || 0);
    } catch { /* ignore */ } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h2 className="page-title">搜索</h2>

      <form onSubmit={handleSearch} style={{ marginBottom: 24 }}>
        <div style={{ display: 'flex', gap: 8 }}>
          <input
            className="input"
            type="text"
            placeholder="搜索标题或摘要..."
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            style={{ fontSize: 16, padding: '12px 16px' }}
          />
          <button className="btn btn-primary" type="submit" disabled={loading}>
            {loading ? '搜索中...' : '搜索'}
          </button>
        </div>
      </form>

      {loading && <div className="loading">搜索中...</div>}

      {searched && !loading && (
        <div style={{ marginBottom: 16, fontSize: 14, color: 'var(--text-secondary)' }}>
          找到 {total} 条结果
        </div>
      )}

      {searched && !loading && results.length === 0 && (
        <div className="empty-state">
          <div style={{ fontSize: 48, marginBottom: 16 }}>🔍</div>
          <h3>没有找到相关结果</h3>
          <p>试试其他关键词</p>
        </div>
      )}

      <div className="grid">
        {results.map(c => (
          <Link key={c.id} to={`/detail/${c.id}`} style={{ color: 'inherit' }}>
            <div className="card">
              <div style={{ fontWeight: 600, marginBottom: 8 }}>{c.title}</div>
              {c.summary && (
                <p style={{ fontSize: 14, color: 'var(--text-secondary)', marginBottom: 12 }}>
                  {c.summary}
                </p>
              )}
              <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap' }}>
                {c.category && <span className="badge badge-approved">{c.category}</span>}
                {c.tags?.map(tag => (
                  <span key={tag} className="tag">{tag}</span>
                ))}
              </div>
            </div>
          </Link>
        ))}
      </div>
    </div>
  );
}