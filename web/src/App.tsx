import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import Layout from './components/layout/Layout';
import HomePage from './pages/HomePage';
import OwnersPage from './pages/OwnersPage';
import OwnerProfilePage from './pages/OwnerProfilePage';
import SeasonsPage from './pages/SeasonsPage';
import SeasonDetailPage from './pages/SeasonDetailPage';
import BeltPage from './pages/BeltPage';
import AwardsPage from './pages/AwardsPage';
import NotFoundPage from './pages/NotFoundPage';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000,
      retry: 1,
    },
  },
});

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route element={<Layout />}>
            <Route index element={<HomePage />} />
            <Route path="owners" element={<OwnersPage />} />
            <Route path="owners/:id" element={<OwnerProfilePage />} />
            <Route path="seasons" element={<SeasonsPage />} />
            <Route path="seasons/:year" element={<SeasonDetailPage />} />
            <Route path="belt" element={<BeltPage />} />
            <Route path="awards" element={<AwardsPage />} />
            <Route path="*" element={<NotFoundPage />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  );
}
