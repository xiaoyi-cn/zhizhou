import { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import api from '../../lib/api';

interface Content {
  id: string;
  url: string;
  title: string;
  summary: string;
  category: string;
  tags: string[];
  status: string;
  source_type: string;
  raw_content: string;
  created_at: string;
  updated_at: string;
}

export default function Detail() {
  const { id } = useParams<{ id: string }>();
  const [content, setContent] = useState<Content | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchDetail = async () => {
      try {
        const res: any = await api.get(`/contents/${id}`);
        setContent(res.content);
      } catch { /* ignore */ } finally {
        setLoading(false);
      }
    };
    fetchDetail();
  }, [id]);

  if (loading) return <div className="loading">加载中...</div>;
  if (!content) return <div className="empty-state"><h3>内容不存在</h3></div>;

  return (
    <div style={{ maxWidth: 800 }}>
      <Link to="/library" style={{ fontSize: 14, marginBottom: 16, display: 'inline-block' }}>← 返回知识库</Link>

      <h1 style={{ fontSize: 28, fontWeight: 700, marginBottom: 16 }}>{content.title}</h1>

      <div style={{ display: 'flex', gap: 16, marginBottom: 24, fontSize: 13, color: 'var(--text-muted)' }}>
        <span>收录于 {new Date(content.created_at).toLocaleDateString('zh-CN')}</span>
        {content.url && (
          <a href={content.url} target="_blank" rel="noopener noreferrer">查看原文 →</a>
        )}
      </div>

      <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap', marginBottom: 24 }}>
        {content.category && <span className="badge badge-approved">{content.category}</span>}
        {content.tags?.map(tag => (
          <span key={tag} className="tag">{tag}</span>
        ))}
      </div>

      {content.summary && (
        <div className="card" style={{ marginBottom: 24 }}>
          <h3 style={{ fontSize: 16, fontWeight: 600, marginBottom: 12 }}>AI 摘要</h3>
          <p style={{ fontSize: 15, lineHeight: 1.8, color: 'var(--text-secondary)' }}>{content.summary}</p>
        </div>
      )}

      {content.raw_content && (
        <div className="card">
          <h3 style={{ fontSize: 16, fontWeight: 600, marginBottom: 12 }}>原文内容</h3>
          <div style={{ fontSize: 14, lineHeight: 1.8, color: 'var(--text-secondary)', whiteSpace: 'pre-wrap', maxHeight: 600, overflow: 'auto' }}>
            {content.raw_content.substring(0, 5000)}
            {content.raw_content.length > 5000 && <p style={{ marginTop: 12, color: 'var(--text-muted)' }}>... 内容过长，已截断显示</p>}
          </div>
        </div>
      )}
    </div>
  );
}