import React, { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { apiService, Node } from '../services/api'

const Home: React.FC = () => {
  const [nodes, setNodes] = useState<Node[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    apiService.getNodes()
      .then(setNodes)
      .finally(() => setLoading(false))
  }, [])

  return (
    <div className="space-y-6">
      {/* Hero Section */}
      <div data-testid="home-hero" className="text-center py-12 bg-gradient-to-b from-[#1a1f2b] to-transparent rounded-2xl">
        <h2 data-testid="home-hero-title" className="text-4xl md:text-5xl font-bold mb-4">
          Your Community, <span className="text-orange-500">Your Rules</span>
        </h2>
        <p data-testid="home-hero-description" className="text-[#8b929a] text-lg max-w-2xl mx-auto">
          SubNode is a self-hosted, invite-only community platform. 
          Join your community's node to start posting and discussing.
        </p>
      </div>

      {/* Nodes Grid */}
      <div data-testid="home-nodes-list">
        <h2 className="text-2xl font-bold mb-4">Available Nodes</h2>
        {loading ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {[1, 2, 3].map(i => (
              <div key={i} className="bg-[#1a1f2b] rounded-lg p-6 animate-pulse">
                <div className="h-6 bg-[#242a38] rounded w-3/4 mb-4"></div>
                <div className="h-4 bg-[#242a38] rounded w-full mb-2"></div>
                <div className="h-4 bg-[#242a38] rounded w-2/3"></div>
              </div>
            ))}
          </div>
        ) : nodes.length === 0 ? (
          <div data-testid="home-node-card" className="text-center py-12 text-[#8b929a]">
            <p>No nodes available yet.</p>
            <p className="mt-2">Check back later or contact your host.</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {nodes.map(node => (
              <Link
                key={node.id}
                to={`/node/${node.slug}`}
                data-testid="home-node-card"
                className="bg-[#1a1f2b] hover:bg-[#242a38] rounded-lg p-6 transition-colors border border-transparent hover:border-[#373e4a]"
              >
                <h3 data-testid="home-node-card-title" className="text-xl font-bold mb-2 text-orange-400">
                  {node.name}
                </h3>
                <p className="text-[#8b929a] text-sm mb-4 line-clamp-2">
                  {node.description || 'No description'}
                </p>
                <div className="flex items-center justify-between text-sm text-[#8b929a]">
                  <span>{node.rules ? `${node.rules.split('\n').length} rules` : 'No rules'}</span>
                  <span>→</span>
                </div>
              </Link>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

export default Home
