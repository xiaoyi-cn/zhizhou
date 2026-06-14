import { Outlet, NavLink, useNavigate } from 'react-router-dom';
import { useAuthStore } from '../store/auth';

const navItems = [
  { to: '/ingest', label: '采集' },
  { to: '/review', label: '审核' },
  { to: '/library', label: '知识库' },
  { to: '/search', label: '搜索' },
  { to: '/settings', label: '设置' },
];

export default function Layout() {
  const navigate = useNavigate();
  const { logout } = useAuthStore();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div style={{ minHeight: '100vh' }}>
      <header style={{
        background: 'var(--bg-white)',
        borderBottom: '1px solid var(--border)',
        padding: '0 24px',
        height: 56,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        position: 'sticky',
        top: 0,
        zIndex: 100,
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 32 }}>
          <h1 style={{ fontSize: 20, fontWeight: 700, color: 'var(--primary)' }}>知舟</h1>
          <nav style={{ display: 'flex', gap: 4 }}>
            {navItems.map(item => (
              <NavLink
                key={item.to}
                to={item.to}
                style={({ isActive }) => ({
                  padding: '8px 16px',
                  borderRadius: 'var(--radius)',
                  fontSize: 14,
                  fontWeight: 500,
                  color: isActive ? 'var(--primary)' : 'var(--text-secondary)',
                  background: isActive ? 'var(--primary-light)' : 'transparent',
                })}
              >
                {item.label}
              </NavLink>
            ))}
          </nav>
        </div>
        <div style={{ display: 'flex', gap: 12, alignItems: 'center' }}>
          <NavLink to="/pricing" style={{ fontSize: 13, color: 'var(--text-secondary)' }}>定价</NavLink>
          <NavLink to="/account" style={{ fontSize: 13, color: 'var(--text-secondary)' }}>账户</NavLink>
          <button onClick={handleLogout} className="btn btn-secondary btn-sm">退出</button>
        </div>
      </header>
      <main className="container">
        <Outlet />
      </main>
    </div>
  );
}