import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { authService } from '../services/auth'

interface User {
  id: string
  username: string
  email: string
  role: 'admin' | 'moderator' | 'member'
  created_at: string
}

interface AuthContextType {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (username: string, password: string) => Promise<void>
  register: (username: string, email: string, password: string, code: string) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const token = authService.getToken()
    if (token) {
      authService.getCurrentUser().then(setUser).catch(() => {
        authService.removeToken()
        setUser(null)
      })
    }
    setIsLoading(false)
  }, [])

  const login = async (username: string, password: string) => {
    const response = await authService.login({ username, password })
    setUser(response.user)
  }

  const register = async (username: string, email: string, password: string, code: string) => {
    const response = await authService.register({ username, email, password, code })
    setUser(response.user)
  }

  const logout = () => {
    authService.removeToken()
    setUser(null)
  }

  return (
    <AuthContext.Provider value={{
      user,
      isAuthenticated: !!user,
      isLoading,
      login,
      register,
      logout,
    }}>
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
