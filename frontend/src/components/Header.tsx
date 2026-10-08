import React from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

const Header: React.FC = () => {
  const { user, isAuthenticated, logout } = useAuth()

  return (
    <header data-testid="common-header" className="bg-[#1a1f2b] border-b border-[#373e4a] sticky top-0 z-50">
      <div className="container mx-auto px-4 py-3 max-w-6xl">
        <div className="flex items-center justify-between">
          <Link to="/" data-testid="nav-home" className="flex items-center space-x-2">
            <div className="w-8 h-8 bg-orange-500 rounded-lg flex items-center justify-center">
              <span className="text-white font-bold text-lg">S</span>
            </div>
            <span className="text-xl font-bold text-[#e7e9ea]">SubNode</span>
          </Link>

          <nav data-testid="nav-main" className="flex items-center space-x-4" aria-label="Main navigation">
            {isAuthenticated ? (
              <>
                <Link
                  to="/dashboard"
                  data-testid="nav-dashboard"
                  className="text-[#8b929a] hover:text-[#e7e9ea] transition-colors"
                >
                  Dashboard
                </Link>
                <div className="flex items-center space-x-2">
                  <div className="w-8 h-8 bg-orange-500 rounded-full flex items-center justify-center">
                    <span className="text-white font-medium text-sm">
                      {user?.username.charAt(0).toUpperCase()}
                    </span>
                  </div>
                  <span className="text-[#e7e9ea] hidden sm:block">{user?.username}</span>
                </div>
                <button
                  data-testid="nav-logout"
                  onClick={logout}
                  className="px-4 py-2 text-[#8b929a] hover:text-[#e7e9ea] transition-colors"
                >
                  Logout
                </button>
              </>
            ) : (
              <>
                <Link
                  to="/login"
                  data-testid="nav-login"
                  className="px-4 py-2 text-[#8b929a] hover:text-[#e7e9ea] transition-colors"
                >
                  Login
                </Link>
                <Link
                  to="/register"
                  data-testid="nav-register"
                  className="px-4 py-2 bg-orange-500 hover:bg-orange-600 text-white rounded-lg transition-colors"
                >
                  Register
                </Link>
              </>
            )}
          </nav>
        </div>
      </div>
    </header>
  )
}

export default Header
