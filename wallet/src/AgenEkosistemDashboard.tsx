import React, { useState, useEffect, useCallback } from 'react'
import {
  Users, Package, ShieldCheck, Vote, Network, Plus, Search, Star,
  CheckCircle2, AlertTriangle, ArrowUpRight, Globe, KeyRound, Layers, RefreshCw, Loader2
} from 'lucide-react'
import {
  Agen, OlehOlehProduct, OlehOlehOrder, Asuransi, AsuransiClaim, DaoProposal,
  IBCChannel, IBCCounterpartyChain, formatRupiah, formatDate
} from './types'
import {
  fetchAgents, fetchOlehOlehProducts, fetchOlehOlehOrders,
  fetchAsuransi, fetchAsuransiClaims, fetchDaoProposals, fetchIBCNetworkStatus,
} from './chainApi'

const statusColors: Record<string, string> = {
  active: 'badge-success', inactive: 'badge-warning', suspended: 'badge-danger',
  open: 'badge-warning', resolved: 'badge-success',
  pending: 'badge-warning', paid: 'badge-success', shipped: 'badge-info', delivered: 'badge-success', cancelled: 'badge-danger',
  expired: 'badge-danger', claimed: 'badge-info',
  submitted: 'badge-warning', approved: 'badge-success', rejected: 'badge-danger',
  passed: 'badge-success', closed: 'badge-info',
  Open: 'badge-success', OPEN: 'badge-success', CLOSED: 'badge-danger', Uninitialized: 'badge-warning', UNINITIALIZED: 'badge-warning', INIT: 'badge-info', TRYOPEN: 'badge-info',
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

const AgenEkosistemDashboard: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'agen' | 'oleholeh' | 'asuransi' | 'dao' | 'ibc'>('agen')
  const [searchTerm, setSearchTerm] = useState('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [refreshing, setRefreshing] = useState(false)

  // Live chain state
  const [agents, setAgents] = useState<Agen[]>([])
  const [products, setProducts] = useState<OlehOlehProduct[]>([])
  const [orders, setOrders] = useState<OlehOlehOrder[]>([])
  const [asuransi, setAsuransi] = useState<Asuransi[]>([])
  const [claims, setClaims] = useState<AsuransiClaim[]>([])
  const [proposals, setProposals] = useState<DaoProposal[]>([])
  const [channels, setChannels] = useState<IBCChannel[]>([])
  const [counterparties, setCounterparties] = useState<IBCCounterpartyChain[]>([])

  const load = useCallback(async () => {
    // Fail-soft: load every dataset independently so a single failing
    // endpoint (e.g. IBC without a relayer) never blanks the whole dashboard.
    const settled = await Promise.allSettled([
      fetchAgents(), fetchOlehOlehProducts(), fetchOlehOlehOrders(),
      fetchAsuransi(), fetchAsuransiClaims(), fetchDaoProposals(),
      fetchIBCNetworkStatus(),
    ])

    const [ag, pr, or, as, cl, pp, ibc] = settled
    if (ag.status === 'fulfilled') setAgents(ag.value)
    if (pr.status === 'fulfilled') setProducts(pr.value)
    if (or.status === 'fulfilled') setOrders(or.value)
    if (as.status === 'fulfilled') setAsuransi(as.value)
    if (cl.status === 'fulfilled') setClaims(cl.value)
    if (pp.status === 'fulfilled') setProposals(pp.value)
    if (ibc.status === 'fulfilled') {
      setChannels(ibc.value.channels)
      setCounterparties(ibc.value.counterparties)
    }

    // All datasets failed -> node unreachable. Partial failure -> keep data.
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

  const pusat = agents.find(a => a.level === 'pusat')
  const totalVolume = agents.reduce((s, a) => s + parseInt(a.total_volume || '0'), 0).toString()
  const avgScore = agents.length > 0
    ? (agents.reduce((s, a) => s + parseFloat(a.score || '0'), 0) / agents.length).toFixed(1)
    : '0.0'
  const openComplaints = agents.reduce((s, a) => s + a.complaints.filter(c => c.status === 'open').length, 0)

  const tabs = [
    { key: 'agen', label: 'Multi-Agen', icon: <Users size={16} /> },
    { key: 'oleholeh', label: 'Oleh-oleh', icon: <Package size={16} /> },
    { key: 'asuransi', label: 'Asuransi', icon: <ShieldCheck size={16} /> },
    { key: 'dao', label: 'DAO Voting', icon: <Vote size={16} /> },
    { key: 'ibc', label: 'IBC Bridge', icon: <Network size={16} /> },
  ] as const

  return (
    <div className="space-y-6 animate-fade-in">
      <div className="flex items-center justify-between flex-wrap gap-4">
        <div>
          <h1 className="text-2xl font-bold text-white">Ekosistem Agen Haji & Umroh</h1>
          <p className="text-surface-400 mt-1">Data live dari chain • multi-agen, oleh-oleh, asuransi, DAO & IBC</p>
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
            <AlertTriangle size={16} />
            {error} — pastikan `bizchaind start` berjalan dan REST aktif di port 1317.
          </p>
          <button className="btn-primary text-sm" onClick={handleRefresh} disabled={refreshing}>
            <RefreshCw size={14} className={refreshing ? 'animate-spin' : ''} /> Coba Lagi
          </button>
        </div>
      )}

      {/* ============ MULTI-AGEN ============ */}
      {activeTab === 'agen' && !error && (
        <div className="space-y-6">
          {loading ? (
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
              {[0, 1, 2, 3].map(i => <SkeletonCard key={i} />)}
            </div>
          ) : agents.length === 0 ? (
            <EmptyState
              icon={<Users size={28} className="text-emerald-400" />}
              title="Belum ada agen terdaftar"
              hint="Daftarkan agen pusat/cabang/subagen via `bizchaind tx agen create-agen`"
            />
          ) : (
            <>
              {/* Hierarchy Overview */}
              <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                <div className="glass-card p-5">
                  <div className="flex items-center gap-3 mb-3">
                    <div className="w-10 h-10 rounded-xl bg-emerald-500/10 flex items-center justify-center"><Users size={20} className="text-emerald-400" /></div>
                    <span className="text-surface-400 text-sm">Total Agen</span>
                  </div>
                  <p className="text-2xl font-bold text-white">{agents.length}</p>
                  <p className="text-xs text-surface-500 mt-2">{agents.filter(a => a.level === 'pusat').length} pusat • {agents.filter(a => a.level === 'cabang').length} cabang • {agents.filter(a => a.level === 'subagen').length} subagen</p>
                </div>
                <div className="glass-card p-5">
                  <div className="flex items-center gap-3 mb-3">
                    <div className="w-10 h-10 rounded-xl bg-blue-500/10 flex items-center justify-center"><ArrowUpRight size={20} className="text-blue-400" /></div>
                    <span className="text-surface-400 text-sm">Total Volume</span>
                  </div>
                  <p className="text-2xl font-bold text-white">{formatRupiah(totalVolume)}</p>
                  <p className="text-xs text-surface-500 mt-2">Seluruh jaringan</p>
                </div>
                <div className="glass-card p-5">
                  <div className="flex items-center gap-3 mb-3">
                    <div className="w-10 h-10 rounded-xl bg-purple-500/10 flex items-center justify-center"><Star size={20} className="text-purple-400" /></div>
                    <span className="text-surface-400 text-sm">Rata-rata Skor</span>
                  </div>
                  <p className="text-2xl font-bold text-white">{avgScore}</p>
                  <p className="text-xs text-surface-500 mt-2">Rekam jejak (0-100)</p>
                </div>
                <div className="glass-card p-5">
                  <div className="flex items-center gap-3 mb-3">
                    <div className="w-10 h-10 rounded-xl bg-amber-500/10 flex items-center justify-center"><AlertTriangle size={20} className="text-amber-400" /></div>
                    <span className="text-surface-400 text-sm">Komplain Terbuka</span>
                  </div>
                  <p className="text-2xl font-bold text-white">{openComplaints}</p>
                  <p className="text-xs text-surface-500 mt-2">Perlu resolusi</p>
                </div>
              </div>

              {/* Tree */}
              {pusat && (
                <div className="glass-card p-6">
                  <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2"><Layers size={18} className="text-emerald-400" />Hierarki Agen (Pusat → Cabang → Subagen)</h3>
                  <div className="flex flex-col items-center gap-4">
                    <div className="p-4 rounded-2xl bg-emerald-500/10 border border-emerald-500/30 text-center w-64">
                      <p className="font-semibold text-white">{pusat.name}</p>
                      <p className="text-xs text-emerald-400 mt-1">{pusat.level.toUpperCase()} • Skor {pusat.score}</p>
                    </div>
                    <div className="flex gap-6 flex-wrap justify-center">
                      {agents.filter(a => a.parent_id === String(pusat.id)).map(cabang => (
                        <div key={cabang.id} className="flex flex-col items-center gap-3">
                          <div className="p-3 rounded-xl bg-blue-500/10 border border-blue-500/30 text-center w-52">
                            <p className="text-sm font-semibold text-white">{cabang.name}</p>
                            <p className="text-xs text-blue-400 mt-1">CABANG • Skor {cabang.score}</p>
                          </div>
                          <div className="flex gap-2 flex-wrap justify-center">
                            {agents.filter(s => s.parent_id === String(cabang.id)).map(sub => (
                              <div key={sub.id} className={`p-2.5 rounded-lg text-center w-44 ${sub.status === 'suspended' ? 'bg-red-500/10 border border-red-500/30' : 'bg-white/5 border border-white/10'}`}>
                                <p className="text-xs font-medium text-white">{sub.name}</p>
                                <p className="text-[10px] text-surface-500 mt-0.5">SUBAGEN • {sub.score}</p>
                              </div>
                            ))}
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                </div>
              )}

              {/* Agent Cards */}
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {agents.map(agen => (
                  <div key={agen.id} className="glass-card p-6">
                    <div className="flex items-start justify-between mb-3">
                      <div className="flex items-center gap-3">
                        <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-emerald-500/10 to-teal-500/10 flex items-center justify-center">
                          <Users size={22} className="text-emerald-400" />
                        </div>
                        <div>
                          <h3 className="text-base font-semibold text-white">{agen.name}</h3>
                          <p className="text-xs text-surface-500">{agen.level} • {agen.address.slice(0, 12)}...</p>
                        </div>
                      </div>
                      <span className={`badge ${statusColors[agen.status] || 'badge-warning'}`}>{agen.status}</span>
                    </div>
                    <div className="flex items-center gap-1 mb-3">
                      <span className="text-2xl font-bold text-emerald-400">{agen.score}</span>
                      <span className="text-xs text-surface-500">skor</span>
                      <div className="flex-1 bg-white/5 rounded-full h-1.5 mx-2">
                        <div className="bg-gradient-to-r from-emerald-500 to-teal-500 h-1.5 rounded-full" style={{ width: `${Math.min(parseFloat(agen.score || '0'), 100)}%` }} />
                      </div>
                      <span className="text-xs text-yellow-400">★ {agen.rating_avg}</span>
                    </div>
                    <div className="space-y-1.5 text-xs text-surface-500 mb-4">
                      <div className="flex justify-between"><span>Komisi</span><span className="text-surface-300">{agen.commission_rate}%</span></div>
                      <div className="flex justify-between"><span>Downline</span><span className="text-surface-300">{agen.total_downline}</span></div>
                      <div className="flex justify-between"><span>Penjualan</span><span className="text-surface-300">{agen.total_sales}</span></div>
                      <div className="flex justify-between"><span>Volume</span><span className="text-surface-300">{formatRupiah(agen.total_volume)}</span></div>
                    </div>
                    {agen.complaints.length > 0 && (
                      <div className="space-y-2">
                        {agen.complaints.map(c => (
                          <div key={c.id} className={`p-2.5 rounded-lg text-xs ${c.status === 'open' ? 'bg-amber-500/10 border border-amber-500/20' : 'bg-white/5'}`}>
                            <div className="flex items-center justify-between mb-1">
                              <span className="text-white">{c.reason}</span>
                              <span className={`badge ${statusColors[c.status] || 'badge-warning'}`}>{c.status}</span>
                            </div>
                            {c.resolution && <p className="text-surface-500">{c.resolution}</p>}
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </>
          )}
        </div>
      )}

      {/* ============ OLEH-OLEH ============ */}
      {activeTab === 'oleholeh' && !error && (
        <div className="space-y-4">
          <div className="flex justify-between items-center flex-wrap gap-3">
            <div className="relative flex-1 max-w-md">
              <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-surface-400" />
              <input className="input-field pl-10" placeholder="Cari produk oleh-oleh..." value={searchTerm} onChange={(e) => setSearchTerm(e.target.value)} />
            </div>
            <button className="btn-primary"><Plus size={18} />Tambah Produk</button>
          </div>
          {loading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
              {[0, 1, 2, 3].map(i => <SkeletonCard key={i} />)}
            </div>
          ) : products.length === 0 ? (
            <EmptyState
              icon={<Package size={28} className="text-emerald-400" />}
              title="Belum ada produk oleh-oleh"
              hint="Daftarkan produk via `bizchaind tx oleholeh create-product`"
            />
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
              {products.filter(p => p.name.toLowerCase().includes(searchTerm.toLowerCase())).map(p => (
                <div key={p.id} className="glass-card p-5">
                  <div className="flex items-start justify-between mb-3">
                    <div className="w-14 h-14 rounded-xl bg-emerald-500/10 flex items-center justify-center">
                      <Package size={24} className="text-emerald-400" />
                    </div>
                    <span className={`badge ${statusColors[p.status] || 'badge-warning'}`}>{p.status}</span>
                  </div>
                  <h3 className="text-sm font-semibold text-white mb-1">{p.name}</h3>
                  <p className="text-xs text-surface-500 mb-2 line-clamp-2">{p.description}</p>
                  <p className="text-lg font-bold text-emerald-400 mb-2">{formatRupiah(p.price)}</p>
                  <div className="flex items-center justify-between text-xs mb-3">
                    <span className="text-surface-500">{p.category}</span>
                    <span className={`${p.stock === 0 ? 'text-red-400' : 'text-surface-300'}`}>Stok: {p.stock}</span>
                  </div>
                  <button className="btn-secondary w-full justify-center text-sm" disabled={p.stock === 0}>Pre-order</button>
                </div>
              ))}
            </div>
          )}

          {/* Orders */}
          <div className="glass-card overflow-hidden">
            <div className="px-4 py-3 border-b border-white/5 flex items-center justify-between">
              <h3 className="text-sm font-semibold text-white">Pre-order Masuk</h3>
              <span className="text-xs text-surface-500">{orders.length} pesanan</span>
            </div>
            {loading ? (
              <div className="p-6 space-y-3 animate-pulse">
                {[0, 1, 2].map(i => <div key={i} className="h-10 bg-white/5 rounded-lg" />)}
              </div>
            ) : orders.length === 0 ? (
              <div className="p-6 text-center text-sm text-surface-500">Belum ada pre-order. Buat via `bizchaind tx oleholeh order-product`.</div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b border-white/5">
                      <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">#</th>
                      <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Produk</th>
                      <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Jamaah</th>
                      <th className="text-center text-xs text-surface-400 font-medium px-4 py-3">Qty</th>
                      <th className="text-right text-xs text-surface-400 font-medium px-4 py-3">Total</th>
                      <th className="text-center text-xs text-surface-400 font-medium px-4 py-3">Status</th>
                    </tr>
                  </thead>
                  <tbody>
                    {orders.map(o => (
                      <tr key={o.id} className="border-b border-white/5 hover:bg-white/5 transition-colors">
                        <td className="px-4 py-3 text-sm text-surface-400">#{o.id}</td>
                        <td className="px-4 py-3 text-sm text-white">{o.product_name}</td>
                        <td className="px-4 py-3 text-sm text-surface-300 font-mono">{o.jamaah}</td>
                        <td className="px-4 py-3 text-sm text-surface-300 text-center">{o.quantity}</td>
                        <td className="px-4 py-3 text-sm text-emerald-400 text-right font-medium">{formatRupiah(o.total)}</td>
                        <td className="px-4 py-3 text-center"><span className={`badge ${statusColors[o.status] || 'badge-warning'}`}>{o.status}</span></td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        </div>
      )}

      {/* ============ ASURANSI ============ */}
      {activeTab === 'asuransi' && !error && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-semibold text-white">Polis Asuransi Digital</h3>
            <button className="btn-primary"><Plus size={18} />Terbitkan Polis</button>
          </div>
          {loading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {[0, 1].map(i => <SkeletonCard key={i} />)}
            </div>
          ) : asuransi.length === 0 ? (
            <EmptyState
              icon={<ShieldCheck size={28} className="text-emerald-400" />}
              title="Belum ada polis asuransi"
              hint="Terbitkan polis via `bizchaind tx asuransi create-asuransi`"
            />
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {asuransi.map(a => (
                <div key={a.id} className="glass-card p-6">
                  <div className="flex items-start justify-between mb-4">
                    <div className="flex items-center gap-3">
                      <div className="w-12 h-12 rounded-xl bg-emerald-500/10 flex items-center justify-center"><ShieldCheck size={22} className="text-emerald-400" /></div>
                      <div>
                        <h3 className="text-base font-semibold text-white">{a.policy_type} insurance</h3>
                        <p className="text-xs text-surface-500">{a.provider} • {a.jamaah}</p>
                      </div>
                    </div>
                    <span className={`badge ${statusColors[a.status] || 'badge-warning'}`}>{a.status}</span>
                  </div>
                  <div className="grid grid-cols-2 gap-3 text-xs mb-4">
                    <div className="p-3 rounded-lg bg-white/5">
                      <p className="text-surface-500 mb-1">Premi</p>
                      <p className="text-white font-medium">{formatRupiah(a.premium)}</p>
                    </div>
                    <div className="p-3 rounded-lg bg-white/5">
                      <p className="text-surface-500 mb-1">Coverage</p>
                      <p className="text-emerald-400 font-medium">{formatRupiah(a.coverage)}</p>
                    </div>
                  </div>
                  <div className="flex items-center justify-between text-xs text-surface-500 mb-4">
                    <span>Mulai: {a.start_date}</span>
                    <span>Berakhir: {a.end_date}</span>
                  </div>
                  <div className="flex gap-2">
                    <button className="btn-secondary flex-1 justify-center text-sm">Klaim</button>
                    <button className="btn-primary flex-1 justify-center text-sm">Detail Polis</button>
                  </div>
                </div>
              ))}
            </div>
          )}

          {/* Claims */}
          <div className="glass-card p-6">
            <h3 className="text-sm font-semibold text-white mb-3">Proses Klaim (Transparan)</h3>
            {loading ? (
              <div className="space-y-3 animate-pulse">
                {[0, 1].map(i => <div key={i} className="h-16 bg-white/5 rounded-xl" />)}
              </div>
            ) : claims.length === 0 ? (
              <p className="text-sm text-surface-500">Belum ada klaim. Ajukan via `bizchaind tx asuransi submit-claim`.</p>
            ) : (
              <div className="space-y-3">
                {claims.map(c => {
                  const asuransiItem = asuransi.find(a => a.id === c.asuransi_id)
                  return (
                    <div key={c.id} className="p-4 rounded-xl bg-white/5">
                      <div className="flex items-start justify-between mb-2 flex-wrap gap-2">
                        <div>
                          <p className="text-sm text-white">{c.reason}</p>
                          <p className="text-xs text-surface-500">Polis #{c.asuransi_id} ({asuransiItem?.policy_type}) • {c.jamaah}</p>
                        </div>
                        <div className="flex items-center gap-2">
                          <span className="text-sm text-emerald-400 font-medium">{formatRupiah(c.amount)}</span>
                          <span className={`badge ${statusColors[c.status] || 'badge-warning'}`}>{c.status}</span>
                        </div>
                      </div>
                      {c.decision_note && <p className="text-xs text-surface-400">Keputusan: {c.decision_note}</p>}
                      <p className="text-[10px] text-surface-500 mt-2">{c.submitted_at ? formatDate(c.submitted_at) : '—'}</p>
                    </div>
                  )
                })}
              </div>
            )}
          </div>
        </div>
      )}

      {/* ============ DAO VOTING ============ */}
      {activeTab === 'dao' && !error && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-semibold text-white">Voting DAO Asosiasi Agen</h3>
            <button className="btn-primary"><Plus size={18} />Buat Proposal</button>
          </div>
          {loading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {[0, 1].map(i => <SkeletonCard key={i} />)}
            </div>
          ) : proposals.length === 0 ? (
            <EmptyState
              icon={<Vote size={28} className="text-emerald-400" />}
              title="Belum ada proposal DAO"
              hint="Buat proposal via `bizchaind tx dao create-proposal`"
            />
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {proposals.map(p => {
                const totalVotes = p.votes.length
                const counts: Record<string, number> = {}
                p.options.forEach(o => counts[o] = p.votes.filter(v => v.option === o).length)
                return (
                  <div key={p.id} className="glass-card p-6">
                    <div className="flex items-start justify-between mb-3">
                      <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-xl bg-emerald-500/10 flex items-center justify-center"><Vote size={20} className="text-emerald-400" /></div>
                        <div>
                          <h3 className="text-base font-semibold text-white">{p.title}</h3>
                          <p className="text-xs text-surface-500">{p.description}</p>
                        </div>
                      </div>
                      <span className={`badge ${statusColors[p.status] || 'badge-warning'}`}>{p.status}</span>
                    </div>
                    <div className="space-y-2 mb-3">
                      {p.options.map(opt => (
                        <div key={opt}>
                          <div className="flex justify-between text-xs mb-1">
                            <span className="text-surface-300">{opt}</span>
                            <span className="text-surface-500">{counts[opt] || 0} suara ({totalVotes > 0 ? Math.round(((counts[opt] || 0) / totalVotes) * 100) : 0}%)</span>
                          </div>
                          <div className="bg-white/5 rounded-full h-2">
                            <div className="bg-gradient-to-r from-emerald-500 to-teal-500 h-2 rounded-full" style={{ width: `${totalVotes > 0 ? ((counts[opt] || 0) / totalVotes) * 100 : 0}%` }} />
                          </div>
                        </div>
                      ))}
                    </div>
                    <div className="flex items-center justify-between text-xs text-surface-500">
                      <span>{totalVotes} suara • Deadline {p.deadline}</span>
                      {p.result_option && <span className="text-emerald-400 font-medium">Hasil: {p.result_option}</span>}
                    </div>
                  </div>
                )
              })}
            </div>
          )}
        </div>
      )}

      {/* ============ IBC ============ */}
      {activeTab === 'ibc' && !error && (
        <div className="space-y-4">
          <div className="glass-card p-6">
            <div className="flex items-center gap-3 mb-4">
              <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-emerald-500/10 to-blue-500/10 flex items-center justify-center"><Network size={24} className="text-emerald-400" /></div>
              <div>
                <h3 className="text-lg font-semibold text-white">Integrasi IBC (Inter-Blockchain Communication)</h3>
                <p className="text-xs text-surface-500">Transfer aset & interoperabilitas dengan blockchain Cosmos lain</p>
              </div>
            </div>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-3 text-xs">
              <div className="p-3 rounded-lg bg-white/5">
                <p className="text-surface-500 mb-1 flex items-center gap-1"><Globe size={12} />Client</p>
                <p className="text-white font-medium">{counterparties.length} tendermint</p>
              </div>
              <div className="p-3 rounded-lg bg-white/5">
                <p className="text-surface-500 mb-1 flex items-center gap-1"><KeyRound size={12} />Koneksi</p>
                <p className="text-white font-medium">{counterparties.filter(c => c.status === 'Open').length} aktif</p>
              </div>
              <div className="p-3 rounded-lg bg-white/5">
                <p className="text-surface-500 mb-1 flex items-center gap-1"><Layers size={12} />Channel</p>
                <p className="text-white font-medium">{channels.length} transfer</p>
              </div>
              <div className="p-3 rounded-lg bg-white/5">
                <p className="text-surface-500 mb-1 flex items-center gap-1"><CheckCircle2 size={12} />Status</p>
                <p className="text-emerald-400 font-medium">{channels.some(c => c.state === 'OPEN') ? 'Ready' : 'Setup needed'}</p>
              </div>
            </div>
          </div>

          {loading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {[0, 1].map(i => <SkeletonCard key={i} />)}
            </div>
          ) : channels.length === 0 && counterparties.length === 0 ? (
            <EmptyState
              icon={<Network size={28} className="text-emerald-400" />}
              title="Belum ada koneksi IBC"
              hint="Setup client & channel via CLI (hermes / relayer) atau `bizchaind tx ibc`"
            />
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="glass-card p-6">
                <h3 className="text-sm font-semibold text-white mb-3">Channel Transfer</h3>
                {channels.length === 0 ? (
                  <p className="text-sm text-surface-500">Belum ada channel IBC terdaftar.</p>
                ) : (
                  <div className="space-y-2">
                    {channels.map(ch => (
                      <div key={`${ch.port_id}-${ch.channel_id}`} className="p-3 rounded-lg bg-white/5 flex items-center justify-between">
                        <div>
                          <p className="text-sm text-white font-mono">{ch.channel_id} ({ch.port_id})</p>
                          <p className="text-xs text-surface-500">↔ {ch.counterparty_channel} • {ch.connection_id}</p>
                        </div>
                        <span className={`badge ${statusColors[ch.state] || 'badge-warning'}`}>{ch.state}</span>
                      </div>
                    ))}
                  </div>
                )}
              </div>
              <div className="glass-card p-6">
                <h3 className="text-sm font-semibold text-white mb-3">Blockchain Counterparty</h3>
                {counterparties.length === 0 ? (
                  <p className="text-sm text-surface-500">Belum ada client/connection IBC.</p>
                ) : (
                  <div className="space-y-2">
                    {counterparties.map(chain => (
                      <div key={chain.client_id} className="p-3 rounded-lg bg-white/5 flex items-center justify-between">
                        <div>
                          <p className="text-sm text-white">{chain.chain_id}</p>
                          <p className="text-xs text-surface-500">{chain.client_id} • {chain.connection_id}</p>
                        </div>
                        <span className={`badge ${statusColors[chain.status] || 'badge-warning'}`}>{chain.status}</span>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>
          )}

          <div className="glass-card p-6">
            <h3 className="text-sm font-semibold text-white mb-3">Transfer Aset (IBC Transfer)</h3>
            <div className="flex flex-col md:flex-row gap-3">
              <div className="flex-1">
                <p className="text-xs text-surface-500 mb-2">Token</p>
                <input className="input-field" placeholder="mis. 1000000upoint" />
              </div>
              <div className="flex-1">
                <p className="text-xs text-surface-500 mb-2">Chain Tujuan</p>
                <select className="input-field">
                  {counterparties.length > 0
                    ? counterparties.map(c => <option key={c.chain_id}>{c.chain_id}</option>)
                    : <option>Belum ada chain</option>}
                </select>
              </div>
              <div className="flex-1">
                <p className="text-xs text-surface-500 mb-2">Alamat Tujuan</p>
                <input className="input-field" placeholder="cosmos1..." />
              </div>
              <div className="flex items-end">
                <button className="btn-primary" disabled={channels.length === 0}><Network size={16} />Kirim via IBC</button>
              </div>
            </div>
          </div>
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

export default AgenEkosistemDashboard
