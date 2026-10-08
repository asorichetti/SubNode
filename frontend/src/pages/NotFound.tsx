import React from 'react'
import { Link } from 'react-router-dom'

const NotFound: React.FC = () => {
  return (
    <div data-testid="common-not-found" className="text-center py-24">
      <div className="text-8xl font-bold text-orange-500 mb-4">404</div>
      <h1 className="text-3xl font-bold mb-4">Page Not Found</h1>
      <p className="text-[#8b929a] mb-8">
        The page you're looking for doesn't exist or has been moved.
      </p>
      <Link
        to="/"
        className="inline-block px-6 py-3 bg-orange-500 hover:bg-orange-600 text-white rounded-lg transition-colors"
      >
        Go Home
      </Link>
    </div>
  )
}

export default NotFound
