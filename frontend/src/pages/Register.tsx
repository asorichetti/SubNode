import React, { useState } from 'react'
import { useAuth } from '../context/AuthContext'
import { Link } from 'react-router-dom'

const Register: React.FC = () => {
  const { register } = useAuth()
  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [code, setCode] = useState('')
  const [error, setError] = useState('')
  const [isLoading, setIsLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setIsLoading(true)

    try {
      await register(username, email, password, code)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Registration failed')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="max-w-md mx-auto mt-12">
      <div className="bg-[#1a1f2b] rounded-lg p-8 border border-[#373e4a]">
        <h1 className="text-2xl font-bold mb-2 text-center">Join SubNode</h1>
        <p className="text-[#8b929a] text-sm text-center mb-6">
          Enter your invite code to create an account
        </p>
        
        {error && (
          <div className="bg-red-900/20 border border-red-500/50 text-red-400 px-4 py-3 rounded mb-4">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-[#8b929a] mb-1">
              Invite Code
            </label>
            <input
              type="text"
              value={code}
              onChange={(e) => setCode(e.target.value.toUpperCase())}
              className="w-full px-4 py-2 bg-[#0f1419] border border-[#373e4a] rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-500 text-[#e7e9ea] font-mono tracking-wider"
              placeholder="XXXX-XXXX"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-[#8b929a] mb-1">
              Username
            </label>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full px-4 py-2 bg-[#0f1419] border border-[#373e4a] rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-500 text-[#e7e9ea]"
              minLength={3}
              maxLength={50}
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-[#8b929a] mb-1">
              Email
            </label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full px-4 py-2 bg-[#0f1419] border border-[#373e4a] rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-500 text-[#e7e9ea]"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-[#8b929a] mb-1">
              Password
            </label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-4 py-2 bg-[#0f1419] border border-[#373e4a] rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-500 text-[#e7e9ea]"
              minLength={8}
              required
            />
          </div>

          <button
            type="submit"
            disabled={isLoading}
            className="w-full bg-orange-500 hover:bg-orange-600 disabled:bg-orange-500/50 text-white py-2 rounded-lg font-medium transition-colors"
          >
            {isLoading ? 'Creating account...' : 'Register'}
          </button>
        </form>

        <p className="mt-6 text-center text-[#8b929a]">
          Already have an account?{' '}
          <Link to="/login" className="text-orange-400 hover:text-orange-300">
            Login
          </Link>
        </p>
      </div>
    </div>
  )
}

export default Register
