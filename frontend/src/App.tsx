import React from 'react'
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from './context/AuthContext'
import Header from './components/Header'
import Home from './pages/Home'
import Login from './pages/Login'
import Register from './pages/Register'
import NodePage from './pages/NodePage'
import PostPage from './pages/PostPage'
import Dashboard from './pages/Dashboard'
import NotFound from './pages/NotFound'
import BrokenPage from './pages/BrokenPage'

const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated } = useAuth()
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" />
}

const AppRoutes: React.FC = () => {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />
      <Route path="/node/:slug" element={<NodePage />} />
      <Route path="/post/:id" element={<PostPage />} />
      <Route path="/dashboard" element={
        <ProtectedRoute>
          <Dashboard />
        </ProtectedRoute>
      } />
      <Route path="/broken" element={<BrokenPage />} />
      <Route path="*" element={<NotFound />} />
    </Routes>
  )
}

const App: React.FC = () => {
  return (
    <AuthProvider>
      <Router>
        <div className="min-h-screen bg-[#0f1419]">
          <Header />
          <main data-testid="page-content" className="container mx-auto px-4 py-6 max-w-6xl">
            <AppRoutes />
          </main>
        </div>
      </Router>
    </AuthProvider>
  )
}

export default App
