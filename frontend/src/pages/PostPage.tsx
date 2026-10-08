import React, { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { apiService, Post } from '../services/api'

const PostPage: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const [post, setPost] = useState<Post | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    if (!id) return
    
    apiService.getPost(id)
      .then(setPost)
      .catch((err: any) => {
        setError(err.response?.data?.error || 'Failed to load post')
      })
      .finally(() => setLoading(false))
  }, [id])

  if (loading) {
    return (
      <div className="bg-[#1a1f2b] rounded-lg p-8 animate-pulse">
        <div className="h-8 bg-[#242a38] rounded w-3/4 mb-4"></div>
        <div className="h-4 bg-[#242a38] rounded w-full mb-2"></div>
        <div className="h-4 bg-[#242a38] rounded w-2/3"></div>
      </div>
    )
  }

  if (error || !post) {
    return (
      <div data-testid="common-not-found" className="text-center py-12">
        <h2 className="text-2xl font-bold text-red-400 mb-2">Post Not Found</h2>
        <Link to="/" className="text-orange-400 hover:text-orange-300">
          ← Back to Home
        </Link>
      </div>
    )
  }

  return (
    <div className="max-w-4xl mx-auto">
      <Link to={`/`} data-testid="post-back-to-node" className="inline-block text-orange-400 hover:text-orange-300 mb-4">
        ← Back to Node
      </Link>

      <article data-testid="post-article" className="bg-[#1a1f2b] rounded-lg p-8 border border-[#373e4a]">
        <div className="flex items-start gap-4 mb-6">
          <div className="flex flex-col items-center">
            <button data-testid="post-vote-up" className="text-[#8b929a] hover:text-orange-400 transition-colors">
              <svg className="w-6 h-6" fill="currentColor" viewBox="0 0 20 20">
                <path d="M10 1l2.5 6.5H19l-5.5 4 2.5 6.5L10 13l-6 4 2.5-6.5L1 7.5h6.5z" />
              </svg>
            </button>
            <span className="font-bold text-lg my-1">{post.vote_count}</span>
            <button data-testid="post-vote-down" className="text-[#8b929a] hover:text-blue-400 transition-colors">
              <svg className="w-6 h-6" fill="currentColor" viewBox="0 0 20 20">
                <path d="M10 19l-2.5-6.5H1l5.5-4L10 2l6 4.5L13.5 12.5 19 19h-6.5z" />
              </svg>
            </button>
          </div>

          <div className="flex-1">
            <h1 data-testid="post-title" className="text-3xl font-bold mb-4">{post.title}</h1>
            <div className="flex items-center gap-4 text-sm text-[#8b929a] mb-6">
              <span>Posted by <span className="text-orange-400">{post.author.username}</span></span>
              <span>•</span>
              <span>{new Date(post.created_at).toLocaleDateString()}</span>
            </div>
          </div>
        </div>

        <div className="prose prose-invert max-w-none">
          <p className="text-[#e7e9ea] whitespace-pre-wrap">{post.body || 'No content'}</p>
          {post.url && (
            <a href={post.url} target="_blank" rel="noopener noreferrer" className="text-orange-400 hover:text-orange-300 break-all">
              {post.url}
            </a>
          )}
        </div>

        <div data-testid="post-comments" className="mt-8 pt-6 border-t border-[#373e4a]">
          <h3 className="text-xl font-bold mb-4">Comments</h3>
          <div className="text-[#8b929a] text-center py-8">
            <p>Comments coming soon...</p>
          </div>
        </div>
      </article>
    </div>
  )
}

export default PostPage
