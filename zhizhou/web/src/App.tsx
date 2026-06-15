import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Layout from './components/Layout';
import Login from './features/auth/Login';
import Ingest from './features/ingest/Ingest';
import Review from './features/review/Review';
import Library from './features/library/Library';
import Search from './features/search/Search';
import Detail from './features/detail/Detail';
import Settings from './features/settings/Settings';
import Pricing from './features/pricing/Pricing';
import Account from './features/account/Account';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/" element={<Layout />}>
          <Route index element={<Navigate to="/ingest" replace />} />
          <Route path="ingest" element={<Ingest />} />
          <Route path="review" element={<Review />} />
          <Route path="library" element={<Library />} />
          <Route path="search" element={<Search />} />
          <Route path="detail/:id" element={<Detail />} />
          <Route path="settings" element={<Settings />} />
          <Route path="pricing" element={<Pricing />} />
          <Route path="account" element={<Account />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;