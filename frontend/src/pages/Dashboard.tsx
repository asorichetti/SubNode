import React, { useEffect, useState } from 'react'
import { apiService, InviteCode } from '../services/api'
import { useAuth } from '../context/AuthContext'
import { useToast } from '../components/Toast'

const Dashboard: React.FC = () => {
  const { user } = useAuth()
  const { addToast } = useToast()
  const [invites, setInvites] = useState<InviteCode[]>([])
  const [loading, setLoading] = useState(true)
  const [showCreateInvite, setShowCreateInvite] = useState(false)
  const [maxUses, setMaxUses] = useState(1)
  const [role, setRole] = useState<'member' | 'moderator'>('member')

  useEffect(() => {
    apiService.getInvites()
      .then(setInvites)
      .catch((err) => addToast(err.message || 'Failed to load invites', 'error'))
      .finally(() => setLoading(false))
  }, [])

  const handleCreateInvite = async (e: React.FormEvent) => {
    e.preventDefault()

    try {
      const newInvite = await apiService.createInvite({
        max_uses: maxUses,
        role,
      })
      setInvites([newInvite, ...invites])
      setShowCreateInvite(false)
      addToast(`Invite code created: ${newInvite.code}`, 'success')
    } catch (err: any) {
      addToast(err.message || 'Failed to create invite', 'error')
    }
  }

  const handleRevoke = async (id: string) => {
    try {
      await apiService.revokeInvite(id)
      setInvites(invites.filter(i => i.id !== id))
      addToast('Invite revoked', 'info')
    } catch (err: any) {
      addToast(err.message || 'Failed to revoke invite', 'error')
    }
  }

  const copyToClipboard = (code: string) => {
    navigator.clipboard.writeText(code)
    addToast('Copied to clipboard!', 'success')
  }

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <div className="bg-[#1a1f2b] rounded-lg p-6 border border-[#373e4a]">
        <h1 className="text-2xl font-bold mb-2">Dashboard</h1>
        <p className="text-[#8b929a]">
          Welcome, <span className="text-orange-400">{user?.username}</span> • Role: {user?.role}
        </p>
      </div>

      {/* Invite Codes Section */}
      <div className="bg-[#1a1f2b] rounded-lg p-6 border border-[#373e4a]">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-bold">Invite Codes</h2>
          <button
            onClick={() => setShowCreateInvite(!showCreateInvite)}
            className="px-4 py-2 bg-orange-500 hover:bg-orange-600 text-white rounded-lg transition-colors"
          >
            {showCreateInvite ? 'Cancel' : '+ New Invite'}
          </button>
        </div>

        {showCreateInvite && (
          <form onSubmit={handleCreateInvite} className="mb-6 bg-[#0f1419] rounded-lg p-4 space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-[#8b929a] mb-1">
                  Max Uses
                </label>
                <input
                  type="number"
                  value={maxUses}
                  onChange={(e) => setMaxUses(parseInt(e.target.value) || 1)}
                  className="w-full px-4 py-2 bg-[#1a1f2b] border border-[#373e4a] rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-500 text-[#e7e9ea]"
                  min={1}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-[#8b929a] mb-1">
                  Role
                </label>
                <select
                  value={role}
                  onChange={(e) => setRole(e.target.value as 'member' | 'moderator')}
                  className="w-full px-4 py-2 bg-[#1a1f2b] border border-[#373e4a] rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-500 text-[#e7e9ea]"
                >
                  <option value="member">Member</option>
                  <option value="moderator">Moderator</option>
                </select>
              </div>
            </div>
            <button
              type="submit"
              className="w-full bg-orange-500 hover:bg-orange-600 text-white py-2 rounded-lg transition-colors"
            >
              Generate Invite Code
            </button>
          </form>
        )}

        {loading ? (
          <div className="text-center py-4 text-[#8b929a]">Loading...</div>
        ) : invites.length === 0 ? (
          <div className="text-center py-4 text-[#8b929a]">
            No invite codes yet. Create one to invite others!
          </div>
        ) : (
          <div className="space-y-3">
            {invites.map(invite => (
              <div
                key={invite.id}
                className="flex items-center justify-between bg-[#0f1419] rounded-lg p-4"
              >
                <div>
                  <code className="text-lg font-mono text-orange-400 tracking-wider">
                    {invite.code}
                  </code>
                  <div className="text-sm text-[#8b929a] mt-1">
                    Uses: {invite.uses}/{invite.max_uses} • {invite.role}
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <button
                    onClick={() => copyToClipboard(invite.code)}
                    className="px-3 py-1 text-sm text-[#8b929a] hover:text-[#e7e9ea] transition-colors"
                  >
                    Copy
                  </button>
                  <button
                    onClick={() => handleRevoke(invite.id)}
                    className="px-3 py-1 text-sm text-red-400 hover:text-red-300 transition-colors"
                  >
                    Revoke
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-[#1a1f2b] rounded-lg p-6 border border-[#373e4a]">
          <h3 className="text-[#8b929a] text-sm mb-2">Total Invites</h3>
          <p className="text-3xl font-bold">{invites.length}</p>
        </div>
        <div className="bg-[#1a1f2b] rounded-lg p-6 border border-[#373e4a]">
          <h3 className="text-[#8b929a] text-sm mb-2">Active Invites</h3>
          <p className="text-3xl font-bold">
            {invites.filter(i => i.uses < i.max_uses).length}
          </p>
        </div>
        <div className="bg-[#1a1f2b] rounded-lg p-6 border border-[#373e4a]">
          <h3 className="text-[#8b929a] text-sm mb-2">Total Uses</h3>
          <p className="text-3xl font-bold">
            {invites.reduce((sum, i) => sum + i.uses, 0)}
          </p>
        </div>
      </div>
    </div>
  )
}

export default Dashboard
