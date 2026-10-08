import React, { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { apiService, Node, Post } from '../services/api'
import { useAuth } from '../context/AuthContext'

const NodePage: React.FC = () => {
  const { slug } = useParams<{ slug: string }>()
  const { isAuthenticated } = useAuth()
  const [node, setNode] = useState<Node | null>(null)
  const [posts, setPosts] = useState<Post[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    if (!slug) return
    
    const fetchData = async () => {
      try {
        const [nodeData, postsData] = await Promise.all([
          apiService.getNode(slug),
          apiService.getPosts(slug),
        ])
        setNode(nodeData)
        setPosts(postsData)
      } catch (err: any) {
        setError(err.response?.data?.error || 'Failed to load node')
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [slug])

  if (loading) {
    return (
      <div className="space-y-4">
        {[1, 2, 3].map(i => (
          <div key={i} className="bg-[#1a1f2b] rounded-lg p-6 animate-pulse">
            <div className="h-6 bg-[#242a38] rounded w-3/4 mb-4"></div>
            <div className="h-4 bg-[#242a38] rounded w-full mb-2"></div>
          </div>
        ))}
      </div>
    )
  }

  if (error || !node) {
    return (
      <div data-testid="common-not-found" className="text-center py-12">
        <h2 className="text-2xl font-bold text-red-400 mb-2">Node Not Found</h2>
        <p className="text-[#8b929a]">{error || 'The node you\'re looking for doesn\'t exist.'}</p>
        <Link to="/" className="inline-block mt-4 text-orange-400 hover:text-orange-300">
          ← Back to Home
        </Link>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Node Header */}
      <div data-testid="node-header" className="bg-[#1a1f2b] rounded-lg p-6 border border-[#373e4a]">
        <h1 data-testid="node-page-title" className="text-3xl font-bold text-orange-400 mb-2">{node.name}</h1>
        {node.description && (
          <p className="text-[#8b929a] mb-4">{node.description}</p>
        )}
        {node.rules && (
          <div className="bg-[#0f1419] rounded p-4">
            <h3 className="font-medium mb-2">Rules:</h3>
            <ul className="list-disc list-inside text-[#8b929a] text-sm space-y-1">
              {node.rules.split('\n').map((rule, i) => (
                <li key={i}>{rule}</li>
              ))}
            </ul>
          </div>
        )}
      </div>

      {/* Posts */}
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-bold">Posts</h2>
        {isAuthenticated && (
          <Link
            to={`/node/${slug}/new-post`}
            data-testid="node-new-post-btn"
            className="px-4 py-2 bg-orange-500 hover:bg-orange-600 text-white rounded-lg transition-colors"
          >
            New Post
          </Link>
        )}
      </div>

      {posts.length === 0 ? (
        <div className="text-center py-12 text-[#8b929a]">
          <p>No posts yet. Be the first to post!</p>
        </div>
      ) : (
        <div data-testid="node-posts-list" className="space-y-3">
          {posts.map(post => (
            <Link
              key={post.id}
              to={`/post/${post.id}`}
              className="block bg-[#1a1f2b] hover:bg-[#242a38] rounded-lg p-4 transition-colors border border-transparent hover:border-[#373e4a]"
            >
              <div className="flex items-start gap-3">
                <div className="flex flex-col items-center text-[#8b929a] min-w-[40px]">
                  <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                    <path d="M10 1l2.5 6.5H19l-5.5 4 2.5 6.5L10 13l-6 4 2.5-6.5L1 7.5h6.5z" />
                  </svg>
                  <span className="font-medium">{post.vote_count}</span>
                </div>
                <div className="flex-1">
                  <h3 data-testid="home-node-card-title" className="text-lg font-medium text-[#e7e9ea] hover:text-orange-400">
                    {post.title}
                  </h3>
                  <p className="text-sm text-[#8b929a] mt-1">
                    Posted by <span className="text-orange-400">{post.author.username}</span>
                  </p>
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}

      <Link to="/" data-testid="node-back-home" className="inline-block text-orange-400 hover:text-orange-300">
        ← Back to Home
      </Link>
    </div>
  )
}

export default NodePage
