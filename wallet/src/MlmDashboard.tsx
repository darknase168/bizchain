import React, { useState } from 'react'
import {
  Users, TrendingUp, Award, DollarSign, GitBranch, ChevronRight,
  Crown, Star, Target, ArrowUpRight, ArrowDownRight, Copy, Check,
  UserPlus, Network, Gift, BarChart3, Wallet, Eye, EyeOff
} from 'lucide-react'
import { MLMMember, MLMCommission, MLMRank, CommissionType, formatCoin, formatAddress, formatDate } from './types'

// Mock data
const mockMembers: MLMMember[] = [
  { id: 1, address: 'rupiah1abc...', name: 'Andi Wijaya', sponsor_id: null, sponsor_name: '-', rank: 'Crown', personal_volume: '500000000', group_volume: '5000000000', total_downline: 156, direct_downline: 12, status: 'active', joined_at: '2025-01-15T00:00:00Z', last_active: '2026-07-30T00:00:00Z' },
  { id: 2, address: 'rupiah1def...', name: 'Budi Santoso', sponsor_id: 1, sponsor_name: 'Andi Wijaya', rank: 'Diamond', personal_volume: '200000000', group_volume: '1500000000', total_downline: 45, direct_downline: 8, status: 'active', joined_at: '2025-02-20T00:00:00Z', last_active: '2026-07-29T00:00:00Z' },
  { id: 3, address: 'rupiah1ghi...', name: 'Citra Dewi', sponsor_id: 1, sponsor_name: 'Andi Wijaya', rank: 'Platinum', personal_volume: '150000000', group_volume: '800000000', total_downline: 32, direct_downline: 6, status: 'active', joined_at: '2025-03-10T00:00:00Z', last_active: '2026-07-30T00:00:00Z' },
  { id: 4, address: 'rupiah1jkl...', name: 'Deni Pratama', sponsor_id: 2, sponsor_name: 'Budi Santoso', rank: 'Gold', personal_volume: '80000000', group_volume: '300000000', total_downline: 15, direct_downline: 5, status: 'active', joined_at: '2025-04-05T00:00:00Z', last_active: '2026-07-28T00:00:00Z' },
  { id: 5, address: 'rupiah1mno...', name: 'Eka Putri', sponsor_id: 2, sponsor_name: 'Budi Santoso', rank: 'Silver', personal_volume: '50000000', group_volume: '120000000', total_downline: 8, direct_downline: 3, status: 'active', joined_at: '2025-05-12T00:00:00Z', last_active: '2026-07-30T00:00:00Z' },
  { id: 6, address: 'rupiah1pqr...', name: 'Fajar Ramadhan', sponsor_id: 3, sponsor_name: 'Citra Dewi', rank: 'Bronze', personal_volume: '20000000', group_volume: '45000000', total_downline: 4, direct_downline: 2, status: 'active', joined_at: '2025-06-18T00:00:00Z', last_active: '2026-07-25T00:00:00Z' },
]

const mockCommissions: MLMCommission[] = [
  { id: 1, member_id: 1, member_name: 'Andi Wijaya', from_member: 'Budi Santoso', type: 'sponsor_bonus', amount: '5000000', level: 1, transaction_id: 101, status: 'paid', created_at: '2026-07-28T10:30:00Z' },
  { id: 2, member_id: 1, member_name: 'Andi Wijaya', from_member: 'Citra Dewi', type: 'pairing_bonus', amount: '3000000', level: 1, transaction_id: 102, status: 'paid', created_at: '2026-07-28T14:20:00Z' },
  { id: 3, member_id: 2, member_name: 'Budi Santoso', from_member: 'Deni Pratama', type: 'level_bonus', amount: '1500000', level: 2, transaction_id: 103, status: 'paid', created_at: '2026-07-29T09:15:00Z' },
  { id: 4, member_id: 1, member_name: 'Andi Wijaya', from_member: 'System', type: 'rank_bonus', amount: '10000000', level: 0, transaction_id: 104, status: 'pending', created_at: '2026-07-30T08:00:00Z' },
  { id: 5, member_id: 3, member_name: 'Citra Dewi', from_member: 'Fajar Ramadhan', type: 'matching_bonus', amount: '800000', level: 3, transaction_id: 105, status: 'pending', created_at: '2026-07-30T11:45:00Z' },
]

const rankColors: Record<MLMRank, string> = {
  Bronze: 'text-amber-600 bg-amber-600/10',
  Silver: 'text-gray-300 bg-gray-300/10',
  Gold: 'text-yellow-400 bg-yellow-400/10',
  Platinum: 'text-cyan-400 bg-cyan-400/10',
  Diamond: 'text-blue-400 bg-blue-400/10',
  Crown: 'text-purple-400 bg-purple-400/10',
}

const rankIcons: Record<MLMRank, React.ReactNode> = {
  Bronze: <Award size={14} />,
  Silver: <Award size={14} />,
  Gold: <Star size={14} />,
  Platinum: <Star size={14} />,
  Diamond: <Crown size={14} />,
  Crown: <Crown size={14} />,
}

const commissionLabels: Record<CommissionType, string> = {
  sponsor_bonus: 'Sponsor Bonus',
  pairing_bonus: 'Pairing Bonus',
  matching_bonus: 'Matching Bonus',
  level_bonus: 'Level Bonus',
  rank_bonus: 'Rank Bonus',
  retail_profit: 'Retail Profit',
}

const MlmDashboard: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'overview' | 'network' | 'commissions' | 'ranks'>('overview')
  const [copied, setCopied] = useState<string | null>(null)
  const [showDownline, setShowDownline] = useState(true)
  const [selectedMember, setSelectedMember] = useState<MLMMember | null>(null)

  const currentMember = mockMembers[0] // Simulate logged in member
  const totalEarnings = mockCommissions.filter(c => c.member_id === currentMember.id && c.status === 'paid').reduce((sum, c) => sum + parseInt(c.amount), 0)
  const pendingEarnings = mockCommissions.filter(c => c.member_id === currentMember.id && c.status === 'pending').reduce((sum, c) => sum + parseInt(c.amount), 0)

  const copyToClipboard = async (text: string, label: string) => {
    try {
      await navigator.clipboard.writeText(text)
      setCopied(label)
      setTimeout(() => setCopied(null), 2000)
    } catch { /* fallback */ }
  }

  const referralLink = `https://bizchain.id/ref/${currentMember.address}`

  return (
    <div className="space-y-6 animate-fade-in">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-white">MLM Dashboard</h1>
          <p className="text-surface-400 mt-1">Multi-Level Marketing Management</p>
        </div>
        <div className="flex gap-2">
          {(['overview', 'network', 'commissions', 'ranks'] as const).map(tab => (
            <button key={tab} className={`px-4 py-2 rounded-xl text-sm font-medium transition-all capitalize ${activeTab === tab ? 'bg-primary-500/20 text-primary-400 border border-primary-500/30' : 'bg-white/5 text-surface-400 hover:text-white'}`} onClick={() => setActiveTab(tab)}>
              {tab === 'overview' && <BarChart3 size={16} className="inline mr-2" />}
              {tab === 'network' && <Network size={16} className="inline mr-2" />}
              {tab === 'commissions' && <DollarSign size={16} className="inline mr-2" />}
              {tab === 'ranks' && <Award size={16} className="inline mr-2" />}
              {tab}
            </button>
          ))}
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="glass-card p-5">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-10 h-10 rounded-xl bg-primary-500/10 flex items-center justify-center"><DollarSign size={20} className="text-primary-400" /></div>
            <span className="text-surface-400 text-sm">Total Earnings</span>
          </div>
          <p className="text-2xl font-bold text-white">{formatCoin(totalEarnings.toString(), 'uidr')}</p>
          <div className="flex items-center gap-1 mt-2 text-xs text-primary-400"><ArrowUpRight size={12} /><span>+12.5% this month</span></div>
        </div>
        <div className="glass-card p-5">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-10 h-10 rounded-xl bg-yellow-500/10 flex items-center justify-center"><Wallet size={20} className="text-yellow-400" /></div>
            <span className="text-surface-400 text-sm">Pending</span>
          </div>
          <p className="text-2xl font-bold text-white">{formatCoin(pendingEarnings.toString(), 'uidr')}</p>
          <p className="text-xs text-surface-500 mt-2">Awaiting distribution</p>
        </div>
        <div className="glass-card p-5">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-10 h-10 rounded-xl bg-accent-500/10 flex items-center justify-center"><Users size={20} className="text-accent-400" /></div>
            <span className="text-surface-400 text-sm">Total Downline</span>
          </div>
          <p className="text-2xl font-bold text-white">{currentMember.total_downline}</p>
          <p className="text-xs text-surface-500 mt-2">{currentMember.direct_downline} direct referrals</p>
        </div>
        <div className="glass-card p-5">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-10 h-10 rounded-xl bg-purple-500/10 flex items-center justify-center"><Crown size={20} className="text-purple-400" /></div>
            <span className="text-surface-400 text-sm">Current Rank</span>
          </div>
          <p className="text-2xl font-bold text-white">{currentMember.rank}</p>
          <span className={`inline-flex items-center gap-1 mt-2 px-2 py-1 rounded-lg text-xs font-medium ${rankColors[currentMember.rank]}`}>
            {rankIcons[currentMember.rank]} {currentMember.rank}
          </span>
        </div>
      </div>

      {/* Referral Link */}
      <div className="glass-card p-5">
        <div className="flex items-center justify-between">
          <div>
            <h3 className="text-sm font-medium text-surface-400 mb-1">Your Referral Link</h3>
            <p className="text-white font-mono text-sm">{referralLink}</p>
          </div>
          <button onClick={() => copyToClipboard(referralLink, 'ref')} className="btn-secondary text-sm">
            {copied === 'ref' ? <Check size={16} className="text-primary-400" /> : <Copy size={16} />}
            {copied === 'ref' ? 'Copied!' : 'Copy Link'}
          </button>
        </div>
      </div>

      {activeTab === 'overview' && (
        <>
          {/* Volume Progress */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="glass-card p-6">
              <h3 className="text-lg font-semibold text-white mb-4">Personal Volume</h3>
              <p className="text-3xl font-bold text-primary-400 mb-2">{formatCoin(currentMember.personal_volume, 'uidr')}</p>
              <div className="w-full bg-white/5 rounded-full h-2 mb-2">
                <div className="bg-gradient-to-r from-primary-500 to-accent-500 h-2 rounded-full" style={{ width: '75%' }} />
              </div>
              <p className="text-xs text-surface-500">75% towards next rank requirement</p>
            </div>
            <div className="glass-card p-6">
              <h3 className="text-lg font-semibold text-white mb-4">Group Volume</h3>
              <p className="text-3xl font-bold text-accent-400 mb-2">{formatCoin(currentMember.group_volume, 'uidr')}</p>
              <div className="w-full bg-white/5 rounded-full h-2 mb-2">
                <div className="bg-gradient-to-r from-accent-500 to-purple-500 h-2 rounded-full" style={{ width: '60%' }} />
              </div>
              <p className="text-xs text-surface-500">60% towards next rank requirement</p>
            </div>
          </div>

          {/* Recent Commissions */}
          <div className="glass-card p-6">
            <h3 className="text-lg font-semibold text-white mb-4">Recent Commissions</h3>
            <div className="space-y-3">
              {mockCommissions.slice(0, 5).map(comm => (
                <div key={comm.id} className="flex items-center justify-between p-3 rounded-xl bg-white/5">
                  <div className="flex items-center gap-3">
                    <div className={`w-8 h-8 rounded-lg flex items-center justify-center ${comm.status === 'paid' ? 'bg-primary-500/10' : 'bg-yellow-500/10'}`}>
                      <DollarSign size={14} className={comm.status === 'paid' ? 'text-primary-400' : 'text-yellow-400'} />
                    </div>
                    <div>
                      <p className="text-sm text-white">{commissionLabels[comm.type]}</p>
                      <p className="text-xs text-surface-500">from {comm.from_member} • Level {comm.level}</p>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className="text-sm font-medium text-primary-400">{formatCoin(comm.amount, 'uidr')}</p>
                    <span className={`text-xs ${comm.status === 'paid' ? 'text-primary-400' : 'text-yellow-400'}`}>{comm.status}</span>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </>
      )}

      {activeTab === 'network' && (
        <div className="glass-card p-6">
          <div className="flex items-center justify-between mb-6">
            <h3 className="text-lg font-semibold text-white">Network Tree</h3>
            <button onClick={() => setShowDownline(!showDownline)} className="flex items-center gap-2 text-sm text-surface-400 hover:text-white transition-colors">
              {showDownline ? <EyeOff size={16} /> : <Eye size={16} />}
              {showDownline ? 'Hide' : 'Show'} Downline
            </button>
          </div>
          <div className="space-y-3">
            {mockMembers.map(member => (
              <div key={member.id} className={`p-4 rounded-xl border transition-all cursor-pointer hover:bg-white/5 ${selectedMember?.id === member.id ? 'border-primary-500/30 bg-primary-500/5' : 'border-white/5 bg-white/5'}`} onClick={() => setSelectedMember(member)}>
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-primary-500/10 to-accent-500/10 flex items-center justify-center">
                      <Users size={18} className="text-primary-400" />
                    </div>
                    <div>
                      <p className="text-sm font-medium text-white">{member.name}</p>
                      <p className="text-xs text-surface-500">{formatAddress(member.address)} • Sponsor: {member.sponsor_name}</p>
                    </div>
                  </div>
                  <div className="flex items-center gap-4">
                    <div className="text-right">
                      <p className="text-sm text-white">{member.total_downline} downline</p>
                      <p className="text-xs text-surface-500">GV: {formatCoin(member.group_volume, 'uidr')}</p>
                    </div>
                    <span className={`inline-flex items-center gap-1 px-2 py-1 rounded-lg text-xs font-medium ${rankColors[member.rank]}`}>
                      {rankIcons[member.rank]} {member.rank}
                    </span>
                    <ChevronRight size={16} className="text-surface-500" />
                  </div>
                </div>
                {showDownline && selectedMember?.id === member.id && (
                  <div className="mt-4 pt-4 border-t border-white/5 grid grid-cols-3 gap-4 animate-slide-down">
                    <div className="text-center p-3 rounded-lg bg-white/5">
                      <p className="text-xs text-surface-500">Personal Volume</p>
                      <p className="text-sm font-medium text-white">{formatCoin(member.personal_volume, 'uidr')}</p>
                    </div>
                    <div className="text-center p-3 rounded-lg bg-white/5">
                      <p className="text-xs text-surface-500">Direct Downline</p>
                      <p className="text-sm font-medium text-white">{member.direct_downline}</p>
                    </div>
                    <div className="text-center p-3 rounded-lg bg-white/5">
                      <p className="text-xs text-surface-500">Status</p>
                      <p className={`text-sm font-medium ${member.status === 'active' ? 'text-primary-400' : 'text-red-400'}`}>{member.status}</p>
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      {activeTab === 'commissions' && (
        <div className="glass-card p-6">
          <h3 className="text-lg font-semibold text-white mb-4">Commission History</h3>
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-white/5">
                  <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">ID</th>
                  <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Type</th>
                  <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">From</th>
                  <th className="text-center text-xs text-surface-400 font-medium px-4 py-3">Level</th>
                  <th className="text-right text-xs text-surface-400 font-medium px-4 py-3">Amount</th>
                  <th className="text-center text-xs text-surface-400 font-medium px-4 py-3">Status</th>
                  <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Date</th>
                </tr>
              </thead>
              <tbody>
                {mockCommissions.map(comm => (
                  <tr key={comm.id} className="border-b border-white/5 hover:bg-white/5 transition-colors">
                    <td className="px-4 py-3 text-sm text-surface-300">#{comm.id}</td>
                    <td className="px-4 py-3 text-sm text-white">{commissionLabels[comm.type]}</td>
                    <td className="px-4 py-3 text-sm text-surface-300">{comm.from_member}</td>
                    <td className="px-4 py-3 text-sm text-surface-300 text-center">{comm.level}</td>
                    <td className="px-4 py-3 text-sm text-primary-400 text-right font-medium">{formatCoin(comm.amount, 'uidr')}</td>
                    <td className="px-4 py-3 text-center">
                      <span className={`badge ${comm.status === 'paid' ? 'badge-success' : comm.status === 'pending' ? 'badge-info' : ''}`}>{comm.status}</span>
                    </td>
                    <td className="px-4 py-3 text-sm text-surface-400">{formatDate(comm.created_at)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {activeTab === 'ranks' && (
        <div className="space-y-4">
          <div className="glass-card p-6">
            <h3 className="text-lg font-semibold text-white mb-4">Rank Requirements & Rewards</h3>
            <div className="space-y-3">
              {(['Bronze', 'Silver', 'Gold', 'Platinum', 'Diamond', 'Crown'] as MLMRank[]).map((rank, i) => {
                const requirements = [
                  { pv: '10,000,000', gv: '50,000,000', dd: 2, reward: '500,000' },
                  { pv: '25,000,000', gv: '150,000,000', dd: 4, reward: '1,500,000' },
                  { pv: '50,000,000', gv: '500,000,000', dd: 6, reward: '5,000,000' },
                  { pv: '100,000,000', gv: '1,500,000,000', dd: 8, reward: '15,000,000' },
                  { pv: '250,000,000', gv: '5,000,000,000', dd: 10, reward: '50,000,000' },
                  { pv: '500,000,000', gv: '15,000,000,000', dd: 12, reward: '150,000,000' },
                ]
                const req = requirements[i]
                const isCurrentRank = rank === currentMember.rank
                return (
                  <div key={rank} className={`p-4 rounded-xl border ${isCurrentRank ? 'border-primary-500/30 bg-primary-500/5' : 'border-white/5 bg-white/5'}`}>
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-3">
                        <span className={`inline-flex items-center gap-1 px-3 py-1.5 rounded-lg text-sm font-medium ${rankColors[rank]}`}>
                          {rankIcons[rank]} {rank}
                        </span>
                        {isCurrentRank && <span className="text-xs text-primary-400 font-medium">Current Rank</span>}
                      </div>
                      <div className="flex items-center gap-6 text-sm">
                        <div className="text-right"><p className="text-surface-500 text-xs">PV Required</p><p className="text-white">{req.pv} IDR</p></div>
                        <div className="text-right"><p className="text-surface-500 text-xs">GV Required</p><p className="text-white">{req.gv} IDR</p></div>
                        <div className="text-right"><p className="text-surface-500 text-xs">Direct</p><p className="text-white">{req.dd}</p></div>
                        <div className="text-right"><p className="text-surface-500 text-xs">Reward</p><p className="text-primary-400 font-medium">{req.reward} IDR</p></div>
                      </div>
                    </div>
                  </div>
                )
              })}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default MlmDashboard
