import React, { useState, useEffect, useCallback, useMemo } from 'react'
import {
  Users, Package, CreditCard, Plane, MapPin, Ticket, Share2, Gift,
  Shield, Compass, Plus, Search, Star, ArrowUpRight, Clock, Wallet,
  CheckCircle2, Radio, BadgeCheck, Luggage, RefreshCw, Loader2
} from 'lucide-react'
import {
  Jamaah, Paket, Pembayaran, Visa, Hotel, Ticket as TicketType, Referral,
  Reward, RewardBalance, AuditLog, Keberangkatan, JOURNEY_STAGES,
  formatRupiah, formatNumber, formatDate
} from './types'
import {
  fetchJamaah, fetchPaket, fetchPembayaran, fetchVisa, fetchHotel,
  fetchTickets, fetchReferrals, fetchRewards, fetchAuditLogs,
  fetchKeberangkatan, deriveRewardBalances,
} from './chainApi'

const statusColors: Record<string, string> = {
  active: 'badge-success', inactive: 'badge-warning', blocked: 'badge-danger',
  open: 'badge-success', full: 'badge-warning', closed: 'badge-info', departed: 'badge-info', completed: 'badge-success',
  pending: 'badge-warning', dp_paid: 'badge-info', refunded: 'badge-info', cancelled: 'badge-danger',
  processing: 'badge-warning', issued: 'badge-success', rejected: 'badge-danger', expired: 'badge-danger',
  paid: 'badge-success', awarded: 'badge-success', redeemed: 'badge-info',
  checked_in: 'badge-info', in_transit: 'badge-warning', arrived: 'badge-success', delivered: 'badge-success',
}

const EmptyState: React.FC<{ icon: React.ReactNode; title: string; hint: string }> = ({ icon, title, hint }) => (
  <div className="glass-card p-10 text-center">
    <div className="w-14 h-14 rounded-2xl bg-white/5 flex items-center justify-center mx-auto mb-4">{icon}</div>
    <p className="text-white font-medium">{title}</p>
    <p className="text-xs text-surface-500 mt-1">{hint}</p>
  </div>
)

const SkeletonCard: React.FC = () => (
  <div className="glass-card p-6 animate-pulse">
    <div className="flex items-center gap-3 mb-4">
      <div className="w-10 h-10 rounded-xl bg-white/5" />
      <div className="flex-1 space-y-2">
        <div className="h-3 w-2/3 bg-white/5 rounded" />
        <div className="h-2.5 w-1/3 bg-white/5 rounded" />
      </div>
    </div>
    <div className="h-8 w-1/2 bg-white/5 rounded mb-3" />
    <div className="space-y-2">
      <div className="h-2.5 bg-white/5 rounded" />
      <div className="h-2.5 w-4/5 bg-white/5 rounded" />
    </div>
  </div>
)

const HajiUmrohDashboard: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'overview' | 'jamaah' | 'paket' | 'pembayaran' | 'visa' | 'hotel' | 'ticket' | 'referral' | 'reward' | 'keberangkatan' | 'audit'>('overview')
  const [searchTerm, setSearchTerm] = useState('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [refreshing, setRefreshing] = useState(false)

  // Live chain state
  const [jamaahs, setJamaahs] = useState<Jamaah[]>([])
  const [pakets, setPakets] = useState<Paket[]>([])
  const [pembayarans, setPembayarans] = useState<Pembayaran[]>([])
  const [visas, setVisas] = useState<Visa[]>([])
  const [hotels, setHotels] = useState<Hotel[]>([])
  const [tickets, setTickets] = useState<TicketType[]>([])
  const [referrals, setReferrals] = useState<Referral[]>([])
  const [rewards, setRewards] = useState<Reward[]>([])
  const [auditLogs, setAuditLogs] = useState<AuditLog[]>([])
  const [keberangkatans, setKeberangkatans] = useState<Keberangkatan[]>([])

  const rewardBalances: RewardBalance[] = useMemo(() => deriveRewardBalances(rewards), [rewards])

  const load = useCallback(async () => {
    // Fail-soft: load every dataset independently so a single failing
    // endpoint never blanks the whole dashboard.
    const settled = await Promise.allSettled([
      fetchJamaah(), fetchPaket(), fetchPembayaran(), fetchVisa(), fetchHotel(),
      fetchTickets(), fetchReferrals(), fetchRewards(), fetchAuditLogs(), fetchKeberangkatan(),
    ])

    const [ja, pk, pb, vs, ht, tk, rf, rw, au, kb] = settled
    if (ja.status === 'fulfilled') setJamaahs(ja.value)
    if (pk.status === 'fulfilled') setPakets(pk.value)
    if (pb.status === 'fulfilled') setPembayarans(pb.value)
    if (vs.status === 'fulfilled') setVisas(vs.value)
    if (ht.status === 'fulfilled') setHotels(ht.value)
    if (tk.status === 'fulfilled') setTickets(tk.value)
    if (rf.status === 'fulfilled') setReferrals(rf.value)
    if (rw.status === 'fulfilled') setRewards(rw.value)
    if (au.status === 'fulfilled') setAuditLogs(au.value)
    if (kb.status === 'fulfilled') setKeberangkatans(kb.value)

    if (settled.every(r => r.status === 'rejected')) {
      const reason = settled.find(r => r.status === 'rejected') as PromiseRejectedResult
      setError(reason?.reason?.message || 'Gagal terhubung ke node BizChain')
    } else {
      setError(null)
    }
    setLoading(false)
    setRefreshing(false)
  }, [])

  useEffect(() => { load() }, [load])

  const handleRefresh = () => {
    setRefreshing(true)
    load()
  }

  const totalRevenue = pembayarans.reduce((sum, p) => sum + parseInt(p.paid || '0'), 0).toString()
  const totalJamaah = jamaahs.filter(j => j.status === 'active').length
  const totalOpenQuota = pakets.reduce((sum, p) => sum + (p.quota - p.booked), 0)
  const totalQuota = pakets.reduce((s, p) => s + p.quota, 0)

  const tabs = [
    { key: 'overview', label: 'Overview', icon: <Compass size={16} /> },
    { key: 'jamaah', label: 'Jamaah', icon: <Users size={16} /> },
    { key: 'paket', label: 'Paket', icon: <Package size={16} /> },
    { key: 'pembayaran', label: 'Pembayaran', icon: <CreditCard size={16} /> },
    { key: 'keberangkatan', label: 'Tracking', icon: <Plane size={16} /> },
    { key: 'visa', label: 'Visa', icon: <Shield size={16} /> },
    { key: 'hotel', label: 'Hotel', icon: <MapPin size={16} /> },
    { key: 'ticket', label: 'Tiket', icon: <Ticket size={16} /> },
    { key: 'referral', label: 'Referral', icon: <Share2 size={16} /> },
    { key: 'reward', label: 'Reward', icon: <Gift size={16} /> },
    { key: 'audit', label: 'Audit', icon: <Shield size={16} /> },
  ] as const

  return (
    <div className="space-y-6 animate-fade-in">
      <div className="flex items-center justify-between flex-wrap gap-4">
        <div>
          <h1 className="text-2xl font-bold text-white">Haji & Umroh Dashboard</h1>
          <p className="text-surface-400 mt-1">Data live dari chain • Sistem Agen Haji & Umroh berbasis Blockchain</p>
        </div>
        <div className="flex items-center gap-2 flex-wrap">
          <span className={`flex items-center gap-1.5 text-xs px-2.5 py-1.5 rounded-lg border ${error ? 'bg-red-500/10 text-red-400 border-red-500/20' : 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20'}`}>
            <span className={`status-dot ${error ? 'inactive' : 'active'}`} /> {error ? 'Offline' : 'Live'}
          </span>
          <button
            className="px-3 py-2 rounded-xl text-sm font-medium bg-white/5 text-surface-300 hover:text-white hover:bg-white/10 transition-all flex items-center gap-2"
            onClick={handleRefresh}
            disabled={loading || refreshing}
          >
            <RefreshCw size={14} className={refreshing ? 'animate-spin' : ''} />
            Refresh
          </button>
          {tabs.map(tab => (
            <button
              key={tab.key}
              className={`px-3 py-2 rounded-xl text-sm font-medium transition-all flex items-center gap-2 ${activeTab === tab.key ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30' : 'bg-white/5 text-surface-400 hover:text-white'}`}
              onClick={() => setActiveTab(tab.key)}
            >
              {tab.icon}
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      {/* Connection error banner */}
      {error && (
        <div className="glass-card p-4 border-red-500/20 bg-red-500/10 flex items-center justify-between flex-wrap gap-3">
          <p className="text-red-400 text-sm flex items-center gap-2">
            <span className="status-dot inactive" />
            {error} — pastikan `bizchaind start` berjalan dan REST aktif di port 1317.
          </p>
          <button className="btn-primary text-sm" onClick={handleRefresh} disabled={refreshing}>
            <RefreshCw size={14} className={refreshing ? 'animate-spin' : ''} /> Coba Lagi
          </button>
        </div>
      )}

      {/* ============ OVERVIEW ============ */}
      {activeTab === 'overview' && !error && (
        <div className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div className="glass-card p-5">
              <div className="flex items-center gap-3 mb-3">
                <div className="w-10 h-10 rounded-xl bg-emerald-500/10 flex items-center justify-center"><Wallet size={20} className="text-emerald-400" /></div>
                <span className="text-surface-400 text-sm">Total Pembayaran</span>
              </div>
              <p className="text-2xl font-bold text-white">{loading ? '...' : formatRupiah(totalRevenue)}</p>
              <div className="flex items-center gap-1 mt-2 text-xs text-emerald-400"><ArrowUpRight size={12} /><span>Escrow aman di on-chain</span></div>
            </div>
            <div className="glass-card p-5">
              <div className="flex items-center gap-3 mb-3">
                <div className="w-10 h-10 rounded-xl bg-blue-500/10 flex items-center justify-center"><Users size={20} className="text-blue-400" /></div>
                <span className="text-surface-400 text-sm">Jamaah Aktif</span>
              </div>
              <p className="text-2xl font-bold text-white">{loading ? '...' : totalJamaah}</p>
              <p className="text-xs text-surface-500 mt-2">{jamaahs.length} total terdaftar di blockchain (DID)</p>
            </div>
            <div className="glass-card p-5">
              <div className="flex items-center gap-3 mb-3">
                <div className="w-10 h-10 rounded-xl bg-purple-500/10 flex items-center justify-center"><Package size={20} className="text-purple-400" /></div>
                <span className="text-surface-400 text-sm">Sisa Kuota</span>
              </div>
              <p className="text-2xl font-bold text-white">{loading ? '...' : totalOpenQuota}</p>
              <p className="text-xs text-surface-500 mt-2">Dari {totalQuota} total kuota</p>
            </div>
            <div className="glass-card p-5">
              <div className="flex items-center gap-3 mb-3">
                <div className="w-10 h-10 rounded-xl bg-amber-500/10 flex items-center justify-center"><Radio size={20} className="text-amber-400" /></div>
                <span className="text-surface-400 text-sm">Journey Berjalan</span>
              </div>
              <p className="text-2xl font-bold text-white">{loading ? '...' : keberangkatans.length}</p>
              <p className="text-xs text-surface-500 mt-2">Terlacak 9 tahap</p>
            </div>
          </div>

          {/* Journey Progress */}
          <div className="glass-card p-6">
            <h3 className="text-lg font-semibold text-white mb-4">Progres Perjalanan Jamaah</h3>
            {loading ? (
              <div className="space-y-4 animate-pulse">
                {[0, 1].map(i => <div key={i} className="h-12 bg-white/5 rounded-lg" />)}
              </div>
            ) : keberangkatans.length === 0 ? (
              <p className="text-sm text-surface-500">Belum ada journey. Buat tracking via `bizchaind tx keberangkatan create-keberangkatan`.</p>
            ) : (
              <div className="space-y-4">
                {keberangkatans.map(k => (
                  <div key={k.id}>
                    <div className="flex items-center justify-between mb-2">
                      <div className="flex items-center gap-2">
                        <Plane size={16} className="text-emerald-400" />
                        <span className="text-sm text-white">{k.jamaah}</span>
                        <span className="text-xs text-surface-500">Paket #{k.paket_id}</span>
                      </div>
                      <span className="text-xs text-emerald-400">Tahap {k.stage}/9 • {k.status_label}</span>
                    </div>
                    <div className="flex gap-1">
                      {JOURNEY_STAGES.map((stage, i) => (
                        <div key={i} className="flex-1 h-2 rounded-full transition-all" style={{ background: i < k.stage ? 'linear-gradient(to right, #10b981, #059669)' : 'rgba(255,255,255,0.08)' }} />
                      ))}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Active Packages */}
          <div>
            <h3 className="text-sm font-semibold text-white mb-3">Paket Tersedia</h3>
            {loading ? (
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                {[0, 1, 2].map(i => <SkeletonCard key={i} />)}
              </div>
            ) : pakets.length === 0 ? (
              <EmptyState
                icon={<Package size={28} className="text-emerald-400" />}
                title="Belum ada paket"
                hint="Buat paket via `bizchaind tx paket create-paket`"
              />
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                {pakets.map(paket => (
                  <div key={paket.id} className="glass-card p-6">
                    <div className="flex items-start justify-between mb-3">
                      <div>
                        <span className="text-xs text-emerald-400 bg-emerald-500/10 px-2 py-1 rounded-lg font-medium">{paket.category}</span>
                        <h3 className="text-base font-semibold text-white mt-2">{paket.name}</h3>
                      </div>
                      <span className={`badge ${statusColors[paket.status] || 'badge-warning'}`}>{paket.status}</span>
                    </div>
                    <div className="flex items-center gap-1 mb-2">
                      {[1, 2, 3, 4, 5].map(i => (
                        <Star key={i} size={14} className={paket.reviews.some(r => r.rating >= i) ? 'text-yellow-400 fill-yellow-400' : 'text-surface-600'} />
                      ))}
                      <span className="text-xs text-surface-500 ml-1">({paket.reviews.length} ulasan)</span>
                    </div>
                    <p className="text-2xl font-bold text-emerald-400 mb-4">{formatRupiah(paket.price)}</p>
                    <div className="space-y-2 text-xs text-surface-500 mb-4">
                      <div className="flex justify-between"><span>Kuota</span><span className="text-surface-300">{paket.booked}/{paket.quota} terisi</span></div>
                      <div className="flex justify-between"><span>Hotel</span><span className="text-surface-300">{paket.hotel}</span></div>
                      <div className="flex justify-between"><span>Maskapai</span><span className="text-surface-300">{paket.airline}</span></div>
                      <div className="flex justify-between"><span>Berangkat</span><span className="text-surface-300">{paket.departure_date}</span></div>
                    </div>
                    <div className="w-full bg-white/5 rounded-full h-2 mb-4">
                      <div className="bg-gradient-to-r from-emerald-500 to-teal-500 h-2 rounded-full" style={{ width: `${paket.quota > 0 ? Math.min((paket.booked / paket.quota) * 100, 100) : 0}%` }} />
                    </div>
                    <button className="btn-primary w-full justify-center"><Plus size={16} />Booking Paket</button>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      )}

      {/* ============ JAMAAH ============ */}
      {activeTab === 'jamaah' && !error && (
        <div className="space-y-4">
          <div className="flex gap-3 flex-wrap items-center justify-between">
            <div className="relative flex-1 max-w-md">
              <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-surface-400" />
              <input className="input-field pl-10" placeholder="Cari jamaah..." value={searchTerm} onChange={(e) => setSearchTerm(e.target.value)} />
            </div>
            <button className="btn-primary"><Plus size={18} />Daftarkan Jamaah</button>
          </div>
          {loading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {[0, 1, 2].map(i => <SkeletonCard key={i} />)}
            </div>
          ) : jamaahs.length === 0 ? (
            <EmptyState
              icon={<Users size={28} className="text-emerald-400" />}
              title="Belum ada jamaah terdaftar"
              hint="Daftarkan jamaah via `bizchaind tx jamaah create-jamaah`"
            />
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {jamaahs.filter(j => j.name.toLowerCase().includes(searchTerm.toLowerCase()) || j.passport_number.includes(searchTerm)).map(jamaah => (
                <div key={jamaah.id} className="glass-card p-6">
                  <div className="flex items-start justify-between mb-3">
                    <div className="flex items-center gap-3">
                      <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-emerald-500/10 to-teal-500/10 flex items-center justify-center">
                        <Users size={22} className="text-emerald-400" />
                      </div>
                      <div>
                        <h3 className="text-base font-semibold text-white flex items-center gap-1">{jamaah.name}<BadgeCheck size={14} className="text-emerald-400" /></h3>
                        <p className="text-xs text-surface-500">{jamaah.passport_number}</p>
                      </div>
                    </div>
                    <span className={`badge ${statusColors[jamaah.status] || 'badge-warning'}`}>{jamaah.status}</span>
                  </div>
                  <div className="space-y-1.5 text-xs text-surface-500 mb-4">
                    <div className="flex justify-between"><span>Phone</span><span className="text-surface-300">{jamaah.phone}</span></div>
                    <div className="flex justify-between"><span>Email</span><span className="text-surface-300">{jamaah.email}</span></div>
                    <div className="flex justify-between"><span>DID</span><span className="text-surface-300 font-mono">{jamaah.did}</span></div>
                    <div className="flex justify-between"><span>Vaksin</span><span className="text-surface-300">{jamaah.vaccinations.length} record</span></div>
                    <div className="flex justify-between"><span>Dokumen</span><span className="text-surface-300">{jamaah.documents.length} terverifikasi</span></div>
                  </div>
                  <div className="flex gap-2">
                    <button className="btn-primary flex-1 text-sm justify-center">Lihat Detail</button>
                    <button className="btn-secondary flex-1 text-sm justify-center">Tambah Dokumen</button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {/* ============ PAKET ============ */}
      {activeTab === 'paket' && !error && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-semibold text-white">Paket Haji & Umroh</h3>
            <button className="btn-primary"><Plus size={18} />Buat Paket</button>
          </div>
          {loading ? (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {[0, 1, 2].map(i => <SkeletonCard key={i} />)}
            </div>
          ) : pakets.length === 0 ? (
            <EmptyState
              icon={<Package size={28} className="text-emerald-400" />}
              title="Belum ada paket"
              hint="Buat paket via `bizchaind tx paket create-paket`"
            />
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {pakets.map(paket => (
                <div key={paket.id} className="glass-card p-6">
                  <div className="flex items-start justify-between mb-3">
                    <div>
                      <span className="text-xs text-emerald-400 bg-emerald-500/10 px-2 py-1 rounded-lg font-medium">{paket.category}</span>
                      <h3 className="text-base font-semibold text-white mt-2">{paket.name}</h3>
                    </div>
                    <span className={`badge ${statusColors[paket.status] || 'badge-warning'}`}>{paket.status}</span>
                  </div>
                  <p className="text-2xl font-bold text-emerald-400 mb-2">{formatRupiah(paket.price)}</p>
                  <p className="text-xs text-surface-500 mb-3 flex items-center gap-1"><Clock size={12} />{paket.schedule}</p>
                  <div className="space-y-1.5 text-xs text-surface-500 mb-4">
                    <div className="flex justify-between"><span>Muthawif</span><span className="text-surface-300">{paket.muthawif}</span></div>
                    <div className="flex justify-between"><span>Hotel</span><span className="text-surface-300">{paket.hotel}</span></div>
                    <div className="flex justify-between"><span>Maskapai</span><span className="text-surface-300">{paket.airline}</span></div>
                  </div>
                  <div className="w-full bg-white/5 rounded-full h-2 mb-1">
                    <div className="bg-gradient-to-r from-emerald-500 to-teal-500 h-2 rounded-full" style={{ width: `${paket.quota > 0 ? Math.min((paket.booked / paket.quota) * 100, 100) : 0}%` }} />
                  </div>
                  <p className="text-xs text-surface-500 text-right mb-4">{paket.booked}/{paket.quota} kursi</p>
                  <button className="btn-secondary w-full justify-center text-sm">Kelola Paket</button>
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {/* ============ PEMBAYARAN ============ */}
      {activeTab === 'pembayaran' && !error && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-semibold text-white">Pembayaran Escrow & Cicilan</h3>
            <button className="btn-primary"><Plus size={18} />Buat Pembayaran</button>
          </div>
          {loading ? (
            <div className="space-y-3 animate-pulse">
              {[0, 1, 2].map(i => <div key={i} className="h-16 bg-white/5 rounded-lg" />)}
            </div>
          ) : pembayarans.length === 0 ? (
            <EmptyState
              icon={<CreditCard size={28} className="text-emerald-400" />}
              title="Belum ada pembayaran"
              hint="Buat pembayaran escrow via `bizchaind tx pembayaran create-pembayaran`"
            />
          ) : (
            <div className="glass-card overflow-hidden">
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b border-white/5">
                      <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">#</th>
                      <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Jamaah</th>
                      <th className="text-right text-xs text-surface-400 font-medium px-4 py-3">Total</th>
                      <th className="text-right text-xs text-surface-400 font-medium px-4 py-3">Dibayar</th>
                      <th className="text-right text-xs text-surface-400 font-medium px-4 py-3">Sisa</th>
                      <th className="text-center text-xs text-surface-400 font-medium px-4 py-3">Status</th>
                      <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Escrow</th>
                    </tr>
                  </thead>
                  <tbody>
                    {pembayarans.map(p => (
                      <tr key={p.id} className="border-b border-white/5 hover:bg-white/5 transition-colors">
                        <td className="px-4 py-3 text-sm text-surface-400">#{p.id}</td>
                        <td className="px-4 py-3">
                          <div>
                            <p className="text-sm text-white font-mono">{p.jamaah}</p>
                            <p className="text-xs text-surface-500">Paket #{p.paket_id} • {p.payment_method}</p>
                          </div>
                        </td>
                        <td className="px-4 py-3 text-sm text-white text-right">{formatRupiah(p.total)}</td>
                        <td className="px-4 py-3 text-sm text-emerald-400 text-right font-medium">{formatRupiah(p.paid)}</td>
                        <td className="px-4 py-3 text-sm text-surface-300 text-right">{formatRupiah(p.remaining)}</td>
                        <td className="px-4 py-3 text-center"><span className={`badge ${statusColors[p.status] || 'badge-warning'}`}>{p.status}</span></td>
                        <td className="px-4 py-3">
                          <div className="flex gap-1">
                            {p.escrow_stages.map((stage, i) => (
                              <span key={i} title={`${stage.name}: ${formatRupiah(stage.amount)}`} className={`w-2 h-2 rounded-full ${stage.released ? 'bg-emerald-400' : 'bg-white/20'}`} />
                            ))}
                            {p.escrow_stages.length === 0 && <span className="text-xs text-surface-500">-</span>}
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          )}
        </div>
      )}

      {/* ============ KEBERANGKATAN (Tracking) ============ */}
      {activeTab === 'keberangkatan' && !error && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-semibold text-white">Tracking Keberangkatan (9 Tahap)</h3>
            <button className="btn-primary"><Plus size={18} />Buat Tracking</button>
          </div>
          {loading ? (
            <div className="space-y-4 animate-pulse">
              {[0, 1].map(i => <div key={i} className="h-40 bg-white/5 rounded-2xl" />)}
            </div>
          ) : keberangkatans.length === 0 ? (
            <EmptyState
              icon={<Plane size={28} className="text-emerald-400" />}
              title="Belum ada tracking keberangkatan"
              hint="Buat tracking via `bizchaind tx keberangkatan create-keberangkatan`"
            />
          ) : (
            <div className="space-y-4">
              {keberangkatans.map(k => (
                <div key={k.id} className="glass-card p-6">
                  <div className="flex items-center justify-between mb-4 flex-wrap gap-2">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 rounded-xl bg-emerald-500/10 flex items-center justify-center"><Plane size={20} className="text-emerald-400" /></div>
                      <div>
                        <p className="text-sm font-semibold text-white">{k.jamaah}</p>
                        <p className="text-xs text-surface-500">Paket #{k.paket_id} • Pembayaran #{k.pembayaran_id}</p>
                      </div>
                    </div>
                    <span className="badge badge-info">Tahap {k.stage}/9 • {k.status_label}</span>
                  </div>

                  <div className="flex items-center mb-4">
                    {JOURNEY_STAGES.map((stage, i) => (
                      <React.Fragment key={i}>
                        <div className="flex flex-col items-center gap-1 flex-1">
                          <div className={`w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold transition-all ${i < k.stage ? 'bg-emerald-500 text-white' : i === k.stage ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/50' : 'bg-white/5 text-surface-500'}`}>
                            {i < k.stage ? <CheckCircle2 size={14} /> : i + 1}
                          </div>
                          <span className={`text-[10px] text-center hidden md:block ${i === k.stage ? 'text-emerald-400' : 'text-surface-500'}`}>{stage}</span>
                        </div>
                        {i < JOURNEY_STAGES.length - 1 && (
                          <div className={`h-0.5 flex-1 mb-4 ${i < k.stage - 1 ? 'bg-emerald-500' : 'bg-white/10'}`} />
                        )}
                      </React.Fragment>
                    ))}
                  </div>

                  <div className="grid grid-cols-2 md:grid-cols-4 gap-3 text-xs">
                    <div className="p-3 rounded-lg bg-white/5">
                      <p className="text-surface-500 mb-1">Berangkat</p>
                      <p className="text-white font-medium">{k.departure_date || '-'}</p>
                    </div>
                    <div className="p-3 rounded-lg bg-white/5">
                      <p className="text-surface-500 mb-1">Pulang</p>
                      <p className="text-white font-medium">{k.return_date || '-'}</p>
                    </div>
                    <div className="p-3 rounded-lg bg-white/5">
                      <p className="text-surface-500 mb-1">Manasik</p>
                      <p className="text-white font-medium">{k.manasik_date || '-'}</p>
                    </div>
                    <div className="p-3 rounded-lg bg-white/5">
                      <p className="text-surface-500 mb-1">Bagasi</p>
                      <p className="text-white font-medium">{k.baggage.length} item terdaftar</p>
                    </div>
                  </div>

                  {k.baggage.length > 0 && (
                    <div className="mt-3">
                      <p className="text-xs text-surface-400 mb-2 flex items-center gap-1"><Luggage size={12} />Bagasi ter-track QR/NFC</p>
                      <div className="space-y-2">
                        {k.baggage.map(b => (
                          <div key={b.id} className="flex items-center justify-between p-2.5 rounded-lg bg-white/5">
                            <div className="flex items-center gap-2">
                              <Luggage size={14} className="text-emerald-400" />
                              <span className="text-xs text-white font-mono">{b.tag}</span>
                              <span className="text-xs text-surface-500">{b.weight}</span>
                            </div>
                            <span className={`badge ${statusColors[b.status] || 'badge-warning'}`}>{b.status}</span>
                          </div>
                        ))}
                      </div>
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {/* ============ VISA ============ */}
      {activeTab === 'visa' && !error && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-semibold text-white">Pengurusan Visa</h3>
            <button className="btn-primary"><Plus size={18} />Ajukan Visa</button>
          </div>
          {loading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {[0, 1].map(i => <SkeletonCard key={i} />)}
            </div>
          ) : visas.length === 0 ? (
            <EmptyState
              icon={<Shield size={28} className="text-emerald-400" />}
              title="Belum ada pengajuan visa"
              hint="Ajukan visa via `bizchaind tx visa create-visa`"
            />
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {visas.map(v => (
                <div key={v.id} className="glass-card p-6">
                  <div className="flex items-start justify-between mb-3">
                    <div>
                      <p className="text-sm font-semibold text-white">{v.jamaah}</p>
                      <p className="text-xs text-surface-500">Paket #{v.paket_id}</p>
                    </div>
                    <span className={`badge ${statusColors[v.status] || 'badge-warning'}`}>{v.status}</span>
                  </div>
                  <div className="space-y-1.5 text-xs text-surface-500">
                    <div className="flex justify-between"><span>No. Visa</span><span className="text-surface-300 font-mono">{v.visa_number || '-'}</span></div>
                    <div className="flex justify-between"><span>Terbit</span><span className="text-surface-300">{v.issue_date || '-'}</span></div>
                    <div className="flex justify-between"><span>Kadaluarsa</span><span className="text-surface-300">{v.expiry_date || '-'}</span></div>
                    <div className="flex justify-between"><span>Catatan</span><span className="text-surface-300">{v.notes}</span></div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {/* ============ HOTEL ============ */}
      {activeTab === 'hotel' && !error && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-semibold text-white">Hotel & Akomodasi</h3>
            <button className="btn-primary"><Plus size={18} />Tambah Hotel</button>
          </div>
          {loading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {[0, 1].map(i => <SkeletonCard key={i} />)}
            </div>
          ) : hotels.length === 0 ? (
            <EmptyState
              icon={<MapPin size={28} className="text-emerald-400" />}
              title="Belum ada hotel terdaftar"
              hint="Daftarkan hotel via `bizchaind tx hotel create-hotel`"
            />
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {hotels.map(h => (
                <div key={h.id} className="glass-card p-6">
                  <div className="flex items-start justify-between mb-3">
                    <div className="flex items-center gap-3">
                      <div className="w-12 h-12 rounded-xl bg-emerald-500/10 flex items-center justify-center"><MapPin size={22} className="text-emerald-400" /></div>
                      <div>
                        <h3 className="text-base font-semibold text-white">{h.name}</h3>
                        <p className="text-xs text-surface-500">{h.city} • {h.distance_haram}m dari Haram</p>
                      </div>
                    </div>
                    <span className="text-yellow-400 text-sm">{"★".repeat(parseInt(h.star_rating) || 0)}</span>
                  </div>
                  <p className="text-xl font-bold text-emerald-400 mb-3">{formatRupiah(h.price_per_night)}/malam</p>
                  <div className="space-y-1.5 text-xs text-surface-500 mb-4">
                    <div className="flex justify-between"><span>Tipe Kamar</span><span className="text-surface-300">{h.room_type}</span></div>
                    <div className="flex justify-between"><span>Kamar Tersedia</span><span className="text-surface-300">{h.available_rooms}</span></div>
                  </div>
                  <button className="btn-secondary w-full justify-center text-sm">Kelola Hotel</button>
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {/* ============ TICKET ============ */}
      {activeTab === 'ticket' && !error && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-semibold text-white">Tiket Pesawat (NFT)</h3>
            <button className="btn-primary"><Plus size={18} />Terbitkan Tiket</button>
          </div>
          {loading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {[0, 1].map(i => <SkeletonCard key={i} />)}
            </div>
          ) : tickets.length === 0 ? (
            <EmptyState
              icon={<Ticket size={28} className="text-emerald-400" />}
              title="Belum ada tiket terbit"
              hint="Terbitkan tiket via `bizchaind tx ticket issue-ticket`"
            />
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {tickets.map(t => (
                <div key={t.id} className="glass-card p-6 relative overflow-hidden">
                  <div className="absolute top-0 right-0 w-40 h-40 bg-emerald-500/5 rounded-full blur-2xl" />
                  <div className="flex items-start justify-between mb-4 relative">
                    <div>
                      <p className="text-xs text-surface-500 mb-1">Jamaah</p>
                      <p className="text-sm font-semibold text-white">{t.jamaah}</p>
                      <p className="text-xs text-surface-500 mt-1">Paket #{t.paket_id}</p>
                    </div>
                    <span className={`badge ${statusColors[t.status] || 'badge-warning'}`}>{t.status}</span>
                  </div>
                  <div className="flex items-center justify-between mb-4 relative">
                    <div>
                      <p className="text-xl font-bold text-white">{t.airline}</p>
                      <p className="text-xs text-surface-500">{t.flight_number}</p>
                    </div>
                    <div className="text-right">
                      <p className="text-sm text-surface-300">Seat <span className="text-emerald-400 font-bold">{t.seat}</span></p>
                      <p className="text-xs text-surface-500">{t.schedule}</p>
                    </div>
                  </div>
                  <div className="flex items-center justify-between p-3 rounded-lg bg-white/5 relative">
                    <div className="flex items-center gap-2">
                      <Ticket size={14} className="text-emerald-400" />
                      <span className="text-xs text-surface-400 font-mono">NFT: {t.nft_id}</span>
                    </div>
                    <span className="text-xs font-mono text-surface-500">{t.qr_code}</span>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {/* ============ REFERRAL ============ */}
      {activeTab === 'referral' && !error && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-semibold text-white">Program Referral Agen</h3>
            <button className="btn-primary"><Plus size={18} />Buat Referral</button>
          </div>
          {loading ? (
            <div className="space-y-3 animate-pulse">
              {[0, 1, 2].map(i => <div key={i} className="h-12 bg-white/5 rounded-lg" />)}
            </div>
          ) : referrals.length === 0 ? (
            <EmptyState
              icon={<Share2 size={28} className="text-emerald-400" />}
              title="Belum ada referral"
              hint="Catat referral via `bizchaind tx referral create-referral`"
            />
          ) : (
            <div className="glass-card overflow-hidden">
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b border-white/5">
                      <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">#</th>
                      <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Agen</th>
                      <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Jamaah</th>
                      <th className="text-center text-xs text-surface-400 font-medium px-4 py-3">Rate</th>
                      <th className="text-right text-xs text-surface-400 font-medium px-4 py-3">Komisi</th>
                      <th className="text-center text-xs text-surface-400 font-medium px-4 py-3">Status</th>
                    </tr>
                  </thead>
                  <tbody>
                    {referrals.map(r => (
                      <tr key={r.id} className="border-b border-white/5 hover:bg-white/5 transition-colors">
                        <td className="px-4 py-3 text-sm text-surface-400">#{r.id}</td>
                        <td className="px-4 py-3 text-sm text-white font-mono">{r.agent}</td>
                        <td className="px-4 py-3 text-sm text-surface-300 font-mono">{r.referred_jamaah}</td>
                        <td className="px-4 py-3 text-sm text-surface-300 text-center">{r.commission_rate}bp</td>
                        <td className="px-4 py-3 text-sm text-emerald-400 text-right font-medium">{formatRupiah(r.commission)}</td>
                        <td className="px-4 py-3 text-center"><span className={`badge ${statusColors[r.status] || 'badge-warning'}`}>{r.status}</span></td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          )}
        </div>
      )}

      {/* ============ REWARD ============ */}
      {activeTab === 'reward' && !error && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-semibold text-white">Loyalty Reward</h3>
            <button className="btn-primary"><Plus size={18} />Beri Reward</button>
          </div>
          {loading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {[0, 1].map(i => <SkeletonCard key={i} />)}
            </div>
          ) : rewards.length === 0 ? (
            <EmptyState
              icon={<Gift size={28} className="text-emerald-400" />}
              title="Belum ada reward"
              hint="Beri reward via `bizchaind tx reward award-reward`"
            />
          ) : (
            <>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {rewardBalances.map(b => (
                  <div key={b.jamaah} className="glass-card p-6">
                    <div className="flex items-center gap-3 mb-4">
                      <div className="w-10 h-10 rounded-xl bg-emerald-500/10 flex items-center justify-center"><Gift size={20} className="text-emerald-400" /></div>
                      <div>
                        <p className="text-sm font-semibold text-white">{b.jamaah}</p>
                        <p className="text-xs text-surface-500">Saldo reward</p>
                      </div>
                    </div>
                    <p className="text-3xl font-bold text-emerald-400 mb-4">{formatNumber(b.balance)}</p>
                    <div className="flex gap-3 text-xs">
                      <div className="flex-1 p-3 rounded-lg bg-white/5 text-center">
                        <p className="text-surface-500 mb-1">Earned</p>
                        <p className="text-white font-medium">{formatNumber(b.earned)}</p>
                      </div>
                      <div className="flex-1 p-3 rounded-lg bg-white/5 text-center">
                        <p className="text-surface-500 mb-1">Redeemed</p>
                        <p className="text-white font-medium">{formatNumber(b.redeemed)}</p>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
              <div className="glass-card p-6">
                <h3 className="text-sm font-semibold text-white mb-3">Riwayat Reward</h3>
                <div className="space-y-2">
                  {rewards.map(r => (
                    <div key={r.id} className="flex items-center justify-between p-3 rounded-lg bg-white/5">
                      <div className="flex items-center gap-2">
                        <Gift size={14} className="text-emerald-400" />
                        <span className="text-xs text-white">{r.reason}</span>
                        <span className="text-xs text-surface-500">({r.reward_type})</span>
                      </div>
                      <span className="text-xs text-emerald-400 font-medium">+{formatNumber(r.points)}</span>
                    </div>
                  ))}
                </div>
              </div>
            </>
          )}
        </div>
      )}

      {/* ============ AUDIT ============ */}
      {activeTab === 'audit' && !error && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-semibold text-white">Jejak Audit Immutable</h3>
            <button className="btn-primary"><Plus size={18} />Catat Aksi</button>
          </div>
          {loading ? (
            <div className="space-y-3 animate-pulse">
              {[0, 1, 2].map(i => <div key={i} className="h-12 bg-white/5 rounded-lg" />)}
            </div>
          ) : auditLogs.length === 0 ? (
            <EmptyState
              icon={<Shield size={28} className="text-emerald-400" />}
              title="Belum ada log audit"
              hint="Catat aksi via `bizchaind tx audit log-action`"
            />
          ) : (
            <div className="glass-card overflow-hidden">
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b border-white/5">
                      <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">#</th>
                      <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Modul</th>
                      <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Aksi</th>
                      <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Aktor</th>
                      <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Target</th>
                      <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Data Hash</th>
                      <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Waktu</th>
                    </tr>
                  </thead>
                  <tbody>
                    {auditLogs.map(l => (
                      <tr key={l.id} className="border-b border-white/5 hover:bg-white/5 transition-colors">
                        <td className="px-4 py-3 text-sm text-surface-400">#{l.id}</td>
                        <td className="px-4 py-3 text-sm text-white">{l.module}</td>
                        <td className="px-4 py-3 text-sm text-emerald-400">{l.action}</td>
                        <td className="px-4 py-3 text-sm text-surface-300 font-mono">{l.actor}</td>
                        <td className="px-4 py-3 text-sm text-surface-400">{l.target_id}</td>
                        <td className="px-4 py-3 text-xs text-surface-500 font-mono">{l.data_hash}</td>
                        <td className="px-4 py-3 text-sm text-surface-400">{l.created_at ? formatDate(l.created_at) : '—'}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          )}
        </div>
      )}

      {loading && (
        <div className="flex items-center justify-center gap-2 text-surface-500 text-sm py-2">
          <Loader2 size={16} className="animate-spin" />
          Memuat data dari chain...
        </div>
      )}
    </div>
  )
}

export default HajiUmrohDashboard
