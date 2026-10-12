import React, { useState } from 'react'
import { useAuth } from '../context/AuthContext'
import { useToast } from '../components/Toast'
import { Link } from 'react-router-dom'

const Login: React.FC = () => {
  const { login } = useAuth()
  const { addToast } = useToast()
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [isLoading, setIsLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)

    try {
      await login(username, password)
      addToast('Successfully logged in!', 'success')
    } catch (err: any) {
      const message = err.message || err.details || 'Login failed. Please check your credentials.'
      addToast(message, 'error')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="max-w-md mx-auto mt-12">
      <div className="bg-[#1a1f2b] rounded-lg p-8 border border-[#373e4a]">
        <h1 className="text-2xl font-bold mb-6 text-center">Welcome Back</h1>

        <form data-testid="auth-login-form" onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-[#8b929a] mb-1">
              Username
            </label>
            <input
              data-testid="auth-username"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full px-4 py-2 bg-[#0f1419] border border-[#373e4a] rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-500 text-[#e7e9ea]"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-[#8b929a] mb-1">
              Password
            </label>
            <input
              data-testid="auth-password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-4 py-2 bg-[#0f1419] border border-[#373e4a] rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-500 text-[#e7e9ea]"
              required
            />
          </div>

          <button
            data-testid="auth-submit"
            type="submit"
            disabled={isLoading}
            className="w-full bg-orange-500 hover:bg-orange-600 disabled:bg-orange-500/50 text-white py-2 rounded-lg font-medium transition-colors"
          >
            {isLoading ? 'Logging in...' : 'Login'}
          </button>
        </form>

        <p className="mt-6 text-center text-[#8b929a]">
          Don't have an account?{' '}
          <Link to="/register" className="text-orange-400 hover:text-orange-300">
            Register with invite code
          </Link>
        </p>
      </div>
    </div>
  )
}

export default Login
