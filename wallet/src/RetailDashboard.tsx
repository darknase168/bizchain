import React, { useState, useEffect, useCallback } from 'react'
import {
  Package, Truck, Percent, Star, BarChart3, ArrowUpRight, ArrowDownRight,
  Plus, Search, AlertTriangle, TrendingUp, ShoppingBag, Users, Tag,
  ArrowUp, ArrowDown, RefreshCw, CheckCircle, Calculator, Ruler, Layers, Boxes,
  BookOpen, Scale, Landmark, Loader2, Map
} from 'lucide-react'
import {
  Supplier, PurchaseOrder, Discount, LoyaltyMember, StockMovement, SalesReport,
  PosUnit, PosProduct, Account, JournalEntry, TrialBalance, IncomeStatement,
  BalanceSheet, Ledger, PriceLevelReportItem,
  formatCoin, formatDate, formatRupiah, formatNumber
} from './types'
import {
  fetchPosUnits, fetchPosProducts, fetchPosAccounts, fetchJournalEntries,
  fetchTrialBalance, fetchIncomeStatement, fetchBalanceSheet, fetchLedger,
  fetchPriceLevelReport,
} from './chainApi'

// Mock data (existing tabs — inventory/suppliers/discounts/loyalty/reports)
const mockSuppliers: Supplier[] = [
  { id: 1, name: 'PT Sumber Makmur', contact_person: 'Bapak Hendra', phone: '021-5551234', email: 'hendra@sumbermakmur.id', address: 'Jakarta Barat', wallet_address: 'rupiah1sup1...', products_supplied: 45, total_orders: 120, rating: 4.8, status: 'active' },
  { id: 2, name: 'CV Berkah Jaya', contact_person: 'Ibu Ratna', phone: '021-5555678', email: 'ratna@berkahjaya.id', address: 'Tangerang', wallet_address: 'rupiah1sup2...', products_supplied: 32, total_orders: 85, rating: 4.5, status: 'active' },
  { id: 3, name: 'UD Sentosa', contact_person: 'Bapak Sentosa', phone: '021-5559012', email: 'sentosa@udsentosa.id', address: 'Bekasi', wallet_address: 'rupiah1sup3...', products_supplied: 18, total_orders: 42, rating: 4.2, status: 'active' },
]

const mockPurchaseOrders: PurchaseOrder[] = [
  { id: 1, supplier_id: 1, supplier_name: 'PT Sumber Makmur', items: [{ product_id: 1, product_name: 'Indomie Goreng', quantity: 100, unit_price: '2800', total: '280000' }], total: '280000', status: 'received', ordered_at: '2026-07-25T00:00:00Z', received_at: '2026-07-27T00:00:00Z', tx_hash: 'PO123ABC...' },
  { id: 2, supplier_id: 2, supplier_name: 'CV Berkah Jaya', items: [{ product_id: 2, product_name: 'Aqua 600ml', quantity: 200, unit_price: '2500', total: '500000' }], total: '500000', status: 'ordered', ordered_at: '2026-07-29T00:00:00Z', received_at: null, tx_hash: 'PO456DEF...' },
  { id: 3, supplier_id: 3, supplier_name: 'UD Sentosa', items: [{ product_id: 3, product_name: 'Teh Botol Sosro', quantity: 150, unit_price: '3500', total: '525000' }], total: '525000', status: 'draft', ordered_at: '2026-07-30T00:00:00Z', received_at: null, tx_hash: '' },
]

const mockDiscounts: Discount[] = [
  { id: 1, code: 'MERDEKA17', name: 'Promo Kemerdekaan', type: 'percentage', value: 17, min_purchase: '100000', max_discount: '50000', start_date: '2026-08-01T00:00:00Z', end_date: '2026-08-31T00:00:00Z', usage_limit: 1000, used_count: 0, status: 'active' },
  { id: 2, code: 'NEWBIE10', name: 'Diskon Member Baru', type: 'percentage', value: 10, min_purchase: '50000', max_discount: '25000', start_date: '2026-07-01T00:00:00Z', end_date: '2026-12-31T00:00:00Z', usage_limit: 500, used_count: 234, status: 'active' },
  { id: 3, code: 'FLASH5K', name: 'Flash Sale 5K', type: 'fixed', value: 5000, min_purchase: '25000', max_discount: '5000', start_date: '2026-07-28T00:00:00Z', end_date: '2026-07-30T00:00:00Z', usage_limit: 200, used_count: 200, status: 'expired' },
]

const mockLoyaltyMembers: LoyaltyMember[] = [
  { id: 1, customer_address: 'rupiah1cust1...', name: 'Ibu Sari', phone: '0811111111', tier: 'Platinum', points: 15420, total_spent: '25000000', total_transactions: 156, joined_at: '2025-01-10T00:00:00Z', last_transaction: '2026-07-30T00:00:00Z' },
  { id: 2, customer_address: 'rupiah1cust2...', name: 'Bapak Dedi', phone: '0822222222', tier: 'Gold', points: 8750, total_spent: '12000000', total_transactions: 89, joined_at: '2025-03-15T00:00:00Z', last_transaction: '2026-07-29T00:00:00Z' },
  { id: 3, customer_address: 'rupiah1cust3...', name: 'Ibu Maya', phone: '0833333333', tier: 'Silver', points: 3200, total_spent: '5500000', total_transactions: 42, joined_at: '2025-06-20T00:00:00Z', last_transaction: '2026-07-28T00:00:00Z' },
  { id: 4, customer_address: 'rupiah1cust4...', name: 'Bapak Rudi', phone: '0844444444', tier: 'Bronze', points: 850, total_spent: '1200000', total_transactions: 12, joined_at: '2026-05-01T00:00:00Z', last_transaction: '2026-07-25T00:00:00Z' },
]

const mockStockMovements: StockMovement[] = [
  { id: 1, product_id: 1, product_name: 'Indomie Goreng', type: 'in', quantity: 100, before_stock: 50, after_stock: 150, reference: 'PO-001', note: 'Restock from supplier', created_at: '2026-07-27T10:00:00Z' },
  { id: 2, product_id: 1, product_name: 'Indomie Goreng', type: 'out', quantity: 25, before_stock: 150, after_stock: 125, reference: 'TX-101', note: 'Sales transaction', created_at: '2026-07-28T14:30:00Z' },
  { id: 3, product_id: 2, product_name: 'Aqua 600ml', type: 'out', quantity: 48, before_stock: 200, after_stock: 152, reference: 'TX-102', note: 'Sales transaction', created_at: '2026-07-29T09:15:00Z' },
  { id: 4, product_id: 3, product_name: 'Teh Botol Sosro', type: 'adjustment', quantity: -5, before_stock: 80, after_stock: 75, reference: 'ADJ-001', note: 'Damaged goods', created_at: '2026-07-29T16:00:00Z' },
  { id: 5, product_id: 1, product_name: 'Indomie Goreng', type: 'return', quantity: 3, before_stock: 125, after_stock: 128, reference: 'RET-001', note: 'Customer return', created_at: '2026-07-30T11:00:00Z' },
]

const mockSalesReport: SalesReport = {
  period: 'July 2026',
  total_sales: '45600000',
  total_transactions: 342,
  total_items_sold: 1250,
  average_transaction: '133333',
  top_products: [
    { name: 'Indomie Goreng', quantity: 450, revenue: '1350000' },
    { name: 'Aqua 600ml', quantity: 380, revenue: '1140000' },
    { name: 'Teh Botol Sosro', quantity: 220, revenue: '880000' },
  ],
  sales_by_category: [
    { category: 'Makanan', total: '18500000' },
    { category: 'Minuman', total: '15200000' },
    { category: 'Snack', total: '8900000' },
    { category: 'Lainnya', total: '3000000' },
  ],
  sales_by_hour: [
    { hour: 8, total: '1200000' }, { hour: 10, total: '2800000' }, { hour: 12, total: '4500000' },
    { hour: 14, total: '3200000' }, { hour: 16, total: '2800000' }, { hour: 18, total: '5100000' },
    { hour: 20, total: '3800000' },
  ],
}

const tierColors: Record<string, string> = {
  Bronze: 'text-amber-600 bg-amber-600/10',
  Silver: 'text-gray-300 bg-gray-300/10',
  Gold: 'text-yellow-400 bg-yellow-400/10',
  Platinum: 'text-cyan-400 bg-cyan-400/10',
}

const movementIcons: Record<string, React.ReactNode> = {
  in: <ArrowDown size={14} className="text-green-400" />,
  out: <ArrowUp size={14} className="text-red-400" />,
  adjustment: <RefreshCw size={14} className="text-yellow-400" />,
  return: <ArrowDown size={14} className="text-blue-400" />,
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

type RetailTab = 'inventory' | 'suppliers' | 'discounts' | 'loyalty' | 'reports' | 'akuntansi' | 'satuan' | 'hargalevel' | 'gabungan' | 'cabang'

const RetailDashboard: React.FC = () => {
  const [activeTab, setActiveTab] = useState<RetailTab>('inventory')
  const [selectedBranch, setSelectedBranch] = useState<string>('all')
  const [searchTerm, setSearchTerm] = useState('')
  const [showAddDiscount, setShowAddDiscount] = useState(false)
  const [newDiscount, setNewDiscount] = useState({ code: '', name: '', type: 'percentage', value: '', min_purchase: '', max_discount: '' })

  // Live chain state (POS module)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [refreshing, setRefreshing] = useState(false)
  const [units, setUnits] = useState<PosUnit[]>([])
  const [products, setProducts] = useState<PosProduct[]>([])
  const [accounts, setAccounts] = useState<Account[]>([])
  const [journalEntries, setJournalEntries] = useState<JournalEntry[]>([])
  const [trialBalance, setTrialBalance] = useState<TrialBalance | null>(null)
  const [incomeStatement, setIncomeStatement] = useState<IncomeStatement | null>(null)
  const [balanceSheet, setBalanceSheet] = useState<BalanceSheet | null>(null)
  const [priceLevels, setPriceLevels] = useState<PriceLevelReportItem[]>([])

  // Accounting sub-tab + ledger selector
  const [accTab, setAccTab] = useState<'trial' | 'income' | 'balance' | 'ledger'>('trial')
  const [selectedAccount, setSelectedAccount] = useState(0)
  const [ledger, setLedger] = useState<Ledger | null>(null)

  const load = useCallback(async () => {
    // Fail-soft: load every dataset independently so a single failing
    // endpoint never blanks the whole dashboard.
    const settled = await Promise.allSettled([
      fetchPosUnits(), fetchPosProducts(), fetchPosAccounts(), fetchJournalEntries(),
      fetchTrialBalance(), fetchIncomeStatement(), fetchBalanceSheet(), fetchPriceLevelReport(),
    ])
    const [un, pr, ac, jl, tb, is, bs, pl] = settled
    if (un.status === 'fulfilled') setUnits(un.value)
    if (pr.status === 'fulfilled') setProducts(pr.value)
    if (ac.status === 'fulfilled') setAccounts(ac.value)
    if (jl.status === 'fulfilled') setJournalEntries(jl.value)
    if (tb.status === 'fulfilled') setTrialBalance(tb.value)
    if (is.status === 'fulfilled') setIncomeStatement(is.value)
    if (bs.status === 'fulfilled') setBalanceSheet(bs.value)
    if (pl.status === 'fulfilled') setPriceLevels(pl.value)

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

  // Load ledger for the selected account
  useEffect(() => {
    if (accTab !== 'ledger' || selectedAccount === 0) return
    setLedger(null)
    fetchLedger(selectedAccount).then(setLedger).catch(() => setLedger(null))
  }, [accTab, selectedAccount])

  const handleRefresh = () => {
    setRefreshing(true)
    load()
  }

  // Live low-stock from chain when connected; mock fallback only when offline
  const lowStockItems = products.filter(p => p.stock <= p.min_stock && p.min_stock > 0).slice(0, 3)
  const displayLowStock = error
    ? [{ name: 'Indomie Goreng', stock: 15, min_stock: 20 }, { name: 'Kopi Kapal Api', stock: 8, min_stock: 15 }, { name: 'Gula Pasir 1kg', stock: 5, min_stock: 10 }]
    : lowStockItems.map(p => ({ name: p.name, stock: p.stock, min_stock: p.min_stock }))

  const bundleProducts = products.filter(p => p.is_bundle)
  const productName = (id: number) => products.find(p => p.id === id)?.name ?? `Produk #${id}`

  const tabs: { key: RetailTab; label: string; icon: React.ReactNode }[] = [
    { key: 'inventory', label: 'Inventory', icon: <Package size={16} /> },
    { key: 'suppliers', label: 'Suppliers', icon: <Truck size={16} /> },
    { key: 'discounts', label: 'Discounts', icon: <Percent size={16} /> },
    { key: 'loyalty', label: 'Loyalty', icon: <Star size={16} /> },
    { key: 'reports', label: 'Reports', icon: <BarChart3 size={16} /> },
    { key: 'akuntansi', label: 'Akuntansi', icon: <Calculator size={16} /> },
    { key: 'satuan', label: 'Multi-Satuan', icon: <Ruler size={16} /> },
    { key: 'hargalevel', label: 'Harga Level', icon: <Layers size={16} /> },
    { key: 'gabungan', label: 'Barang Gabungan', icon: <Boxes size={16} /> },
    { key: 'cabang', label: 'Multi-Cabang', icon: <Map size={16} /> },
  ]

  return (
    <div className="space-y-6 animate-fade-in">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-white">Retail Dashboard</h1>
          <p className="text-surface-400 mt-1">
            {selectedBranch === 'all' 
              ? 'All Branches - Inventory, Suppliers, Discounts, Loyalty, Akuntansi & Laporan' 
              : `Branch ${selectedBranch} - Inventory, Suppliers, Discounts, Loyalty, Akuntansi & Laporan`}
          </p>
        </div>
        <div className="flex items-center gap-3 flex-wrap">
          <span className={`flex items-center gap-1.5 text-xs px-2.5 py-1.5 rounded-lg border ${error ? 'bg-red-500/10 text-red-400 border-red-500/20' : 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20'}`}>
            <span className={`status-dot ${error ? 'inactive' : 'active'}`} /> {error ? 'Offline' : 'Live'}
          </span>
          <select 
            className="input-field text-sm"
            value={selectedBranch}
            onChange={(e) => setSelectedBranch(e.target.value)}
          >
            <option value="all">All Branches</option>
            <option value="JKT">Jakarta</option>
            <option value="BDG">Bandung</option>
            <option value="SMG">Semarang</option>
            <option value="SUR">Surabaya</option>
          </select>
          <button
            className="px-3 py-2 rounded-xl text-sm font-medium bg-white/5 text-surface-300 hover:text-white hover:bg-white/10 transition-all flex items-center gap-2"
            onClick={handleRefresh}
            disabled={loading || refreshing}
          >
            <RefreshCw size={14} className={refreshing ? 'animate-spin' : ''} />
            Refresh
          </button>
          <div className="flex gap-2 flex-wrap">
            {tabs.map(tab => (
              <button
                key={tab.key}
                className={`px-3 py-2 rounded-xl text-sm font-medium transition-all flex items-center gap-2 capitalize ${activeTab === tab.key ? 'bg-primary-500/20 text-primary-400 border border-primary-500/30' : 'bg-white/5 text-surface-400 hover:text-white'}`}
                onClick={() => setActiveTab(tab.key)}
              >
                {tab.icon}
                {tab.label}
              </button>
            ))}
          </div>
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

      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="glass-card p-5">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-10 h-10 rounded-xl bg-primary-500/10 flex items-center justify-center"><ShoppingBag size={20} className="text-primary-400" /></div>
            <span className="text-surface-400 text-sm">Monthly Sales</span>
          </div>
          <p className="text-2xl font-bold text-white">{formatCoin(mockSalesReport.total_sales, 'uidr')}</p>
          <div className="flex items-center gap-1 mt-2 text-xs text-primary-400"><ArrowUpRight size={12} /><span>+8.2% vs last month</span></div>
        </div>
        <div className="glass-card p-5">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-10 h-10 rounded-xl bg-accent-500/10 flex items-center justify-center"><TrendingUp size={20} className="text-accent-400" /></div>
            <span className="text-surface-400 text-sm">Transactions</span>
          </div>
          <p className="text-2xl font-bold text-white">{mockSalesReport.total_transactions}</p>
          <p className="text-xs text-surface-500 mt-2">Avg: {formatCoin(mockSalesReport.average_transaction, 'uidr')}</p>
        </div>
        <div className="glass-card p-5">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-10 h-10 rounded-xl bg-yellow-500/10 flex items-center justify-center"><AlertTriangle size={20} className="text-yellow-400" /></div>
            <span className="text-surface-400 text-sm">Low Stock</span>
          </div>
          <p className="text-2xl font-bold text-white">{displayLowStock.length}</p>
          <p className="text-xs text-yellow-400 mt-2">Items need restock</p>
        </div>
        <div className="glass-card p-5">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-10 h-10 rounded-xl bg-purple-500/10 flex items-center justify-center"><Users size={20} className="text-purple-400" /></div>
            <span className="text-surface-400 text-sm">Loyalty Members</span>
          </div>
          <p className="text-2xl font-bold text-white">{mockLoyaltyMembers.length}</p>
          <p className="text-xs text-surface-500 mt-2">Active members</p>
        </div>
      </div>

      {activeTab === 'inventory' && (
        <div className="space-y-4">
          {/* Low Stock Alert */}
          <div className="glass-card p-4 border-yellow-500/20 bg-yellow-500/5">
            <div className="flex items-center gap-2 mb-3">
              <AlertTriangle size={18} className="text-yellow-400" />
              <h3 className="text-sm font-semibold text-yellow-400">Low Stock Alert</h3>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
              {displayLowStock.map((item, i) => (
                <div key={i} className="flex items-center justify-between p-3 rounded-lg bg-white/5">
                  <div>
                    <p className="text-sm text-white">{item.name}</p>
                    <p className="text-xs text-surface-500">Min: {item.min_stock}</p>
                  </div>
                  <span className="text-lg font-bold text-yellow-400">{item.stock}</span>
                </div>
              ))}
            </div>
          </div>

          {/* Stock Movements */}
          <div className="glass-card p-6">
            <h3 className="text-lg font-semibold text-white mb-4">Stock Movements</h3>
            <div className="space-y-3">
              {mockStockMovements.map(mov => (
                <div key={mov.id} className="flex items-center justify-between p-3 rounded-xl bg-white/5">
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center">{movementIcons[mov.type]}</div>
                    <div>
                      <p className="text-sm text-white">{mov.product_name}</p>
                      <p className="text-xs text-surface-500">{mov.reference} • {mov.note}</p>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className={`text-sm font-medium ${mov.type === 'in' || mov.type === 'return' ? 'text-green-400' : 'text-red-400'}`}>
                      {mov.type === 'in' || mov.type === 'return' ? '+' : ''}{mov.quantity}
                    </p>
                    <p className="text-xs text-surface-500">{mov.before_stock} → {mov.after_stock}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Purchase Orders */}
          <div className="glass-card p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold text-white">Purchase Orders</h3>
              <button className="btn-primary text-sm"><Plus size={16} />New PO</button>
            </div>
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-white/5">
                    <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">PO #</th>
                    <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Supplier</th>
                    <th className="text-right text-xs text-surface-400 font-medium px-4 py-3">Total</th>
                    <th className="text-center text-xs text-surface-400 font-medium px-4 py-3">Status</th>
                    <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Ordered</th>
                  </tr>
                </thead>
                <tbody>
                  {mockPurchaseOrders.map(po => (
                    <tr key={po.id} className="border-b border-white/5 hover:bg-white/5 transition-colors">
                      <td className="px-4 py-3 text-sm text-surface-300">PO-{String(po.id).padStart(3, '0')}</td>
                      <td className="px-4 py-3 text-sm text-white">{po.supplier_name}</td>
                      <td className="px-4 py-3 text-sm text-primary-400 text-right font-medium">{formatCoin(po.total, 'uidr')}</td>
                      <td className="px-4 py-3 text-center">
                        <span className={`badge ${po.status === 'received' ? 'badge-success' : po.status === 'ordered' ? 'badge-info' : ''}`}>{po.status}</span>
                      </td>
                      <td className="px-4 py-3 text-sm text-surface-400">{formatDate(po.ordered_at)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}

      {activeTab === 'suppliers' && (
        <div className="space-y-4">
          <div className="flex gap-3">
            <div className="relative flex-1">
              <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-surface-400" />
              <input className="input-field pl-10" placeholder="Search suppliers..." value={searchTerm} onChange={(e) => setSearchTerm(e.target.value)} />
            </div>
            <button className="btn-primary"><Plus size={18} />Add Supplier</button>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {mockSuppliers.map(supplier => (
              <div key={supplier.id} className="glass-card p-6">
                <div className="flex items-start justify-between mb-3">
                  <div className="flex items-center gap-3">
                    <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-primary-500/10 to-accent-500/10 flex items-center justify-center">
                      <Truck size={22} className="text-primary-400" />
                    </div>
                    <div>
                      <h3 className="text-base font-semibold text-white">{supplier.name}</h3>
                      <p className="text-xs text-surface-500">{supplier.contact_person}</p>
                    </div>
                  </div>
                  <div className="flex items-center gap-1">
                    <Star size={14} className="text-yellow-400 fill-yellow-400" />
                    <span className="text-sm text-white">{supplier.rating}</span>
                  </div>
                </div>
                <div className="space-y-1.5 text-xs text-surface-500 mb-4">
                  <div className="flex justify-between"><span>Phone</span><span className="text-surface-300">{supplier.phone}</span></div>
                  <div className="flex justify-between"><span>Email</span><span className="text-surface-300">{supplier.email}</span></div>
                  <div className="flex justify-between"><span>Products</span><span className="text-surface-300">{supplier.products_supplied}</span></div>
                  <div className="flex justify-between"><span>Total Orders</span><span className="text-surface-300">{supplier.total_orders}</span></div>
                </div>
                <div className="flex gap-2">
                  <button className="btn-primary flex-1 text-sm justify-center">Order</button>
                  <button className="btn-secondary flex-1 text-sm justify-center">Details</button>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {activeTab === 'discounts' && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-semibold text-white">Discount & Promotions</h3>
            <button className="btn-primary" onClick={() => setShowAddDiscount(true)}><Plus size={18} />Create Discount</button>
          </div>

          {showAddDiscount && (
            <div className="glass-card p-6 animate-slide-down">
              <h3 className="text-lg font-semibold text-white mb-4">Create New Discount</h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="block text-sm text-surface-400 mb-2">Code *</label>
                  <input className="input-field" placeholder="PROMO10" value={newDiscount.code} onChange={(e) => setNewDiscount(p => ({ ...p, code: e.target.value.toUpperCase() }))} />
                </div>
                <div>
                  <label className="block text-sm text-surface-400 mb-2">Name *</label>
                  <input className="input-field" placeholder="Promo name" value={newDiscount.name} onChange={(e) => setNewDiscount(p => ({ ...p, name: e.target.value }))} />
                </div>
                <div>
                  <label className="block text-sm text-surface-400 mb-2">Type</label>
                  <select className="input-field" value={newDiscount.type} onChange={(e) => setNewDiscount(p => ({ ...p, type: e.target.value }))}>
                    <option value="percentage">Percentage (%)</option>
                    <option value="fixed">Fixed Amount (IDR)</option>
                    <option value="buy_x_get_y">Buy X Get Y</option>
                  </select>
                </div>
                <div>
                  <label className="block text-sm text-surface-400 mb-2">Value</label>
                  <input className="input-field" type="number" placeholder={newDiscount.type === 'percentage' ? '10' : '10000'} value={newDiscount.value} onChange={(e) => setNewDiscount(p => ({ ...p, value: e.target.value }))} />
                </div>
                <div>
                  <label className="block text-sm text-surface-400 mb-2">Min Purchase (IDR)</label>
                  <input className="input-field" type="number" placeholder="50000" value={newDiscount.min_purchase} onChange={(e) => setNewDiscount(p => ({ ...p, min_purchase: e.target.value }))} />
                </div>
                <div>
                  <label className="block text-sm text-surface-400 mb-2">Max Discount (IDR)</label>
                  <input className="input-field" type="number" placeholder="25000" value={newDiscount.max_discount} onChange={(e) => setNewDiscount(p => ({ ...p, max_discount: e.target.value }))} />
                </div>
              </div>
              <div className="flex gap-3">
                <button className="btn-primary" onClick={() => setShowAddDiscount(false)}>Create Discount</button>
                <button className="btn-secondary" onClick={() => setShowAddDiscount(false)}>Cancel</button>
              </div>
            </div>
          )}

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {mockDiscounts.map(discount => (
              <div key={discount.id} className={`glass-card p-6 ${discount.status === 'expired' ? 'opacity-50' : ''}`}>
                <div className="flex items-start justify-between mb-3">
                  <div>
                    <span className="text-xs font-mono text-primary-400 bg-primary-500/10 px-2 py-1 rounded">{discount.code}</span>
                    <h3 className="text-base font-semibold text-white mt-2">{discount.name}</h3>
                  </div>
                  <span className={`badge ${discount.status === 'active' ? 'badge-success' : ''}`}>{discount.status}</span>
                </div>
                <div className="text-center py-4">
                  <p className="text-3xl font-bold text-primary-400">
                    {discount.type === 'percentage' ? `${discount.value}%` : formatCoin(discount.value.toString(), 'uidr')}
                  </p>
                  <p className="text-xs text-surface-500 mt-1">
                    {discount.type === 'percentage' ? 'discount' : 'off'}
                  </p>
                </div>
                <div className="space-y-1.5 text-xs text-surface-500">
                  <div className="flex justify-between"><span>Min Purchase</span><span className="text-surface-300">{formatCoin(discount.min_purchase, 'uidr')}</span></div>
                  <div className="flex justify-between"><span>Max Discount</span><span className="text-surface-300">{formatCoin(discount.max_discount, 'uidr')}</span></div>
                  <div className="flex justify-between"><span>Usage</span><span className="text-surface-300">{discount.used_count}/{discount.usage_limit}</span></div>
                  <div className="flex justify-between"><span>Valid Until</span><span className="text-surface-300">{formatDate(discount.end_date)}</span></div>
                </div>
                <div className="w-full bg-white/5 rounded-full h-1.5 mt-3">
                  <div className="bg-gradient-to-r from-primary-500 to-accent-500 h-1.5 rounded-full" style={{ width: `${(discount.used_count / discount.usage_limit) * 100}%` }} />
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {activeTab === 'loyalty' && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-semibold text-white">Loyalty Program Members</h3>
            <button className="btn-primary"><Plus size={18} />Add Member</button>
          </div>
          <div className="glass-card overflow-hidden">
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-white/5">
                    <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Member</th>
                    <th className="text-center text-xs text-surface-400 font-medium px-4 py-3">Tier</th>
                    <th className="text-right text-xs text-surface-400 font-medium px-4 py-3">Points</th>
                    <th className="text-right text-xs text-surface-400 font-medium px-4 py-3">Total Spent</th>
                    <th className="text-center text-xs text-surface-400 font-medium px-4 py-3">Transactions</th>
                    <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Last Active</th>
                  </tr>
                </thead>
                <tbody>
                  {mockLoyaltyMembers.map(member => (
                    <tr key={member.id} className="border-b border-white/5 hover:bg-white/5 transition-colors">
                      <td className="px-4 py-3">
                        <div className="flex items-center gap-2">
                          <div className="w-8 h-8 rounded-lg bg-primary-500/10 flex items-center justify-center"><Users size={14} className="text-primary-400" /></div>
                          <div>
                            <p className="text-sm text-white">{member.name}</p>
                            <p className="text-xs text-surface-500">{member.phone}</p>
                          </div>
                        </div>
                      </td>
                      <td className="px-4 py-3 text-center">
                        <span className={`inline-flex items-center gap-1 px-2 py-1 rounded-lg text-xs font-medium ${tierColors[member.tier]}`}>
                          <Star size={12} />{member.tier}
                        </span>
                      </td>
                      <td className="px-4 py-3 text-sm text-primary-400 text-right font-medium">{member.points.toLocaleString()}</td>
                      <td className="px-4 py-3 text-sm text-white text-right">{formatCoin(member.total_spent, 'uidr')}</td>
                      <td className="px-4 py-3 text-sm text-surface-300 text-center">{member.total_transactions}</td>
                      <td className="px-4 py-3 text-sm text-surface-400">{formatDate(member.last_transaction)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}

      {activeTab === 'reports' && (
        <div className="space-y-4">
          {/* Top Products */}
          <div className="glass-card p-6">
            <h3 className="text-lg font-semibold text-white mb-4">Top Selling Products</h3>
            <div className="space-y-3">
              {mockSalesReport.top_products.map((product, i) => (
                <div key={i} className="flex items-center justify-between p-3 rounded-xl bg-white/5">
                  <div className="flex items-center gap-3">
                    <span className="w-8 h-8 rounded-lg bg-primary-500/10 flex items-center justify-center text-sm font-bold text-primary-400">#{i + 1}</span>
                    <div>
                      <p className="text-sm text-white">{product.name}</p>
                      <p className="text-xs text-surface-500">{product.quantity} sold</p>
                    </div>
                  </div>
                  <p className="text-sm font-medium text-primary-400">{formatCoin(product.revenue, 'uidr')}</p>
                </div>
              ))}
            </div>
          </div>

          {/* Sales by Category */}
          <div className="glass-card p-6">
            <h3 className="text-lg font-semibold text-white mb-4">Sales by Category</h3>
            <div className="space-y-3">
              {mockSalesReport.sales_by_category.map((cat, i) => {
                const total = parseInt(mockSalesReport.total_sales)
                const value = parseInt(cat.total)
                const percent = ((value / total) * 100).toFixed(1)
                return (
                  <div key={i}>
                    <div className="flex justify-between text-sm mb-1">
                      <span className="text-surface-300">{cat.category}</span>
                      <span className="text-white">{formatCoin(cat.total, 'uidr')} ({percent}%)</span>
                    </div>
                    <div className="w-full bg-white/5 rounded-full h-2">
                      <div className="bg-gradient-to-r from-primary-500 to-accent-500 h-2 rounded-full" style={{ width: `${percent}%` }} />
                    </div>
                  </div>
                )
              })}
            </div>
          </div>

          {/* Sales by Hour */}
          <div className="glass-card p-6">
            <h3 className="text-lg font-semibold text-white mb-4">Sales by Hour</h3>
            <div className="flex items-end gap-2 h-40">
              {mockSalesReport.sales_by_hour.map((hour, i) => {
                const maxVal = Math.max(...mockSalesReport.sales_by_hour.map(h => parseInt(h.total)))
                const height = (parseInt(hour.total) / maxVal) * 100
                return (
                  <div key={i} className="flex-1 flex flex-col items-center gap-1">
                    <span className="text-xs text-surface-500">{formatCoin(hour.total, 'uidr').split(' ')[0]}</span>
                    <div className="w-full bg-gradient-to-t from-primary-500 to-accent-500 rounded-t-lg transition-all hover:opacity-80" style={{ height: `${height}%` }} />
                    <span className="text-xs text-surface-400">{hour.hour}:00</span>
                  </div>
                )
              })}
            </div>
          </div>
        </div>
      )}

      {/* ============ AKUNTANSI (live from chain) ============ */}
      {activeTab === 'akuntansi' && !error && (
        <div className="space-y-4">
          <div className="flex flex-wrap gap-2">
            {([
              { key: 'trial', label: 'Neraca Saldo', icon: <Scale size={14} /> },
              { key: 'income', label: 'Laba Rugi', icon: <TrendingUp size={14} /> },
              { key: 'balance', label: 'Neraca', icon: <Landmark size={14} /> },
              { key: 'ledger', label: 'Buku Besar', icon: <BookOpen size={14} /> },
            ] as const).map(t => (
              <button
                key={t.key}
                className={`px-3 py-2 rounded-xl text-sm font-medium transition-all flex items-center gap-2 ${accTab === t.key ? 'bg-primary-500/20 text-primary-400 border border-primary-500/30' : 'bg-white/5 text-surface-400 hover:text-white'}`}
                onClick={() => setAccTab(t.key)}
              >
                {t.icon}
                {t.label}
              </button>
            ))}
          </div>

          {loading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {[0, 1].map(i => <SkeletonCard key={i} />)}
            </div>
          ) : accTab === 'trial' ? (
            trialBalance && trialBalance.accounts.length > 0 ? (
              <div className="glass-card p-6">
                <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2"><Scale size={18} className="text-primary-400" />Neraca Saldo</h3>
                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead>
                      <tr className="border-b border-white/5">
                        <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Kode</th>
                        <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Akun</th>
                        <th className="text-center text-xs text-surface-400 font-medium px-4 py-3">Tipe</th>
                        <th className="text-right text-xs text-surface-400 font-medium px-4 py-3">Debit</th>
                        <th className="text-right text-xs text-surface-400 font-medium px-4 py-3">Kredit</th>
                        <th className="text-right text-xs text-surface-400 font-medium px-4 py-3">Saldo</th>
                      </tr>
                    </thead>
                    <tbody>
                      {trialBalance.accounts.map(a => (
                        <tr key={a.account_id} className="border-b border-white/5 hover:bg-white/5 transition-colors">
                          <td className="px-4 py-3 text-sm text-surface-400 font-mono">{a.code}</td>
                          <td className="px-4 py-3 text-sm text-white">{a.name}</td>
                          <td className="px-4 py-3 text-center"><span className="text-xs text-surface-400 capitalize">{a.type}</span></td>
                          <td className="px-4 py-3 text-sm text-emerald-400 text-right font-mono">{formatNumber(a.debit)}</td>
                          <td className="px-4 py-3 text-sm text-red-400 text-right font-mono">{formatNumber(a.credit)}</td>
                          <td className="px-4 py-3 text-sm text-white text-right font-medium">{formatRupiah(a.balance)}</td>
                        </tr>
                      ))}
                    </tbody>
                    <tfoot>
                      <tr className="border-t border-white/10">
                        <td colSpan={3} className="px-4 py-3 text-sm text-surface-300 font-medium">Total</td>
                        <td className="px-4 py-3 text-sm text-emerald-400 text-right font-semibold">{formatRupiah(trialBalance.total_debit)}</td>
                        <td className="px-4 py-3 text-sm text-red-400 text-right font-semibold">{formatRupiah(trialBalance.total_credit)}</td>
                        <td />
                      </tr>
                    </tfoot>
                  </table>
                </div>
                {trialBalance.total_debit === trialBalance.total_credit ? (
                  <p className="text-xs text-emerald-400 mt-3 flex items-center gap-1"><CheckCircle size={12} />Balance OK — debit = kredit</p>
                ) : (
                  <p className="text-xs text-yellow-400 mt-3 flex items-center gap-1"><AlertTriangle size={12} />Saldo tidak balance — selisih {formatRupiah((parseInt(trialBalance.total_debit || '0') - parseInt(trialBalance.total_credit || '0')).toString())}</p>
                )}
              </div>
            ) : (
              <EmptyState icon={<Scale size={28} className="text-primary-400" />} title="Belum ada jurnal akuntansi" hint="Buat jurnal via `bizchaind tx pos create-journal-entry` atau transaksi penjualan" />
            )
          ) : accTab === 'income' ? (
            incomeStatement && (incomeStatement.revenues.length > 0 || incomeStatement.expenses.length > 0) ? (
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
                <div className="glass-card p-6">
                  <h3 className="text-sm font-semibold text-white mb-3 flex items-center gap-2"><ArrowUpRight size={16} className="text-emerald-400" />Pendapatan</h3>
                  <div className="space-y-2">
                    {incomeStatement.revenues.map(r => (
                      <div key={r.account_id} className="flex justify-between text-sm">
                        <span className="text-surface-300">{r.code} • {r.name}</span>
                        <span className="text-emerald-400 font-medium">{formatRupiah(r.balance)}</span>
                      </div>
                    ))}
                  </div>
                  <div className="flex justify-between border-t border-white/10 mt-3 pt-3">
                    <span className="text-surface-300 font-medium">Total Pendapatan</span>
                    <span className="text-emerald-400 font-bold">{formatRupiah(incomeStatement.total_revenue)}</span>
                  </div>
                </div>
                <div className="glass-card p-6">
                  <h3 className="text-sm font-semibold text-white mb-3 flex items-center gap-2"><ArrowDownRight size={16} className="text-red-400" />Beban</h3>
                  <div className="space-y-2">
                    {incomeStatement.expenses.map(e => (
                      <div key={e.account_id} className="flex justify-between text-sm">
                        <span className="text-surface-300">{e.code} • {e.name}</span>
                        <span className="text-red-400 font-medium">{formatRupiah(e.balance)}</span>
                      </div>
                    ))}
                  </div>
                  <div className="flex justify-between border-t border-white/10 mt-3 pt-3">
                    <span className="text-surface-300 font-medium">Total Beban</span>
                    <span className="text-red-400 font-bold">{formatRupiah(incomeStatement.total_expense)}</span>
                  </div>
                </div>
                <div className="glass-card p-6 lg:col-span-2">
                  <div className="flex items-center justify-between">
                    <span className="text-surface-300 font-medium">Laba / Rugi Bersih</span>
                    <span className={`text-2xl font-bold ${parseInt(incomeStatement.net_income || '0') >= 0 ? 'text-emerald-400' : 'text-red-400'}`}>{formatRupiah(incomeStatement.net_income)}</span>
                  </div>
                </div>
              </div>
            ) : (
              <EmptyState icon={<TrendingUp size={28} className="text-primary-400" />} title="Belum ada data laba rugi" hint="Laporan muncul setelah ada jurnal pendapatan/beban" />
            )
          ) : accTab === 'balance' ? (
            balanceSheet && (balanceSheet.assets.length > 0 || balanceSheet.liabilities.length > 0 || balanceSheet.equities.length > 0) ? (
              <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
                {([
                  { title: 'Aset', list: balanceSheet.assets, total: balanceSheet.total_assets, color: 'text-emerald-400', icon: <Landmark size={16} className="text-emerald-400" /> },
                  { title: 'Liabilitas', list: balanceSheet.liabilities, total: balanceSheet.total_liabilities, color: 'text-red-400', icon: <ArrowDownRight size={16} className="text-red-400" /> },
                  { title: 'Ekuitas', list: balanceSheet.equities, total: balanceSheet.total_equity, color: 'text-blue-400', icon: <Users size={16} className="text-blue-400" /> },
                ] as const).map(section => (
                  <div key={section.title} className="glass-card p-6">
                    <h3 className="text-sm font-semibold text-white mb-3 flex items-center gap-2">{section.icon}{section.title}</h3>
                    <div className="space-y-2 min-h-[80px]">
                      {section.list.map(a => (
                        <div key={a.account_id} className="flex justify-between text-sm">
                          <span className="text-surface-300">{a.code} • {a.name}</span>
                          <span className={`${section.color} font-medium`}>{formatRupiah(a.balance)}</span>
                        </div>
                      ))}
                      {section.list.length === 0 && <p className="text-xs text-surface-500">Tidak ada akun</p>}
                    </div>
                    <div className="flex justify-between border-t border-white/10 mt-3 pt-3">
                      <span className="text-surface-300 font-medium">Total {section.title}</span>
                      <span className={`${section.color} font-bold`}>{formatRupiah(section.total)}</span>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <EmptyState icon={<Landmark size={28} className="text-primary-400" />} title="Belum ada data neraca" hint="Neraca muncul setelah ada jurnal aset/liabilitas/ekuitas" />
            )
          ) : (
            <div className="space-y-4">
              <div className="glass-card p-6">
                <h3 className="text-sm font-semibold text-white mb-3">Pilih Akun</h3>
                <select className="input-field" value={selectedAccount} onChange={(e) => setSelectedAccount(Number(e.target.value))}>
                  <option value={0}>— Pilih akun —</option>
                  {accounts.map(a => (
                    <option key={a.id} value={a.id}>{a.code} — {a.name}</option>
                  ))}
                </select>
              </div>
              {selectedAccount === 0 && (
                <EmptyState icon={<BookOpen size={28} className="text-primary-400" />} title="Pilih akun untuk melihat buku besar" hint="Detail debit/kredit & saldo berjalan per akun" />
              )}
              {selectedAccount > 0 && ledger && (
                <div className="glass-card p-6">
                  <h3 className="text-lg font-semibold text-white mb-1">{ledger.account.code} — {ledger.account.name}</h3>
                  <p className="text-xs text-surface-500 mb-4 capitalize">{ledger.account.type} • Saldo akhir: <span className="text-white font-medium">{formatRupiah(ledger.ending_balance)}</span></p>
                  <div className="overflow-x-auto">
                    <table className="w-full">
                      <thead>
                        <tr className="border-b border-white/5">
                          <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">#</th>
                          <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Deskripsi</th>
                          <th className="text-right text-xs text-surface-400 font-medium px-4 py-3">Debit</th>
                          <th className="text-right text-xs text-surface-400 font-medium px-4 py-3">Kredit</th>
                          <th className="text-right text-xs text-surface-400 font-medium px-4 py-3">Saldo</th>
                        </tr>
                      </thead>
                      <tbody>
                        {ledger.lines.map((l, li) => (
                          <tr key={li} className="border-b border-white/5 hover:bg-white/5 transition-colors">
                            <td className="px-4 py-3 text-sm text-surface-400">#{l.journal_entry_id}</td>
                            <td className="px-4 py-3 text-sm text-white">{l.description}</td>
                            <td className="px-4 py-3 text-sm text-emerald-400 text-right font-mono">{formatNumber(l.debit)}</td>
                            <td className="px-4 py-3 text-sm text-red-400 text-right font-mono">{formatNumber(l.credit)}</td>
                            <td className="px-4 py-3 text-sm text-white text-right font-medium">{formatRupiah(l.balance)}</td>
                          </tr>
                        ))}
                        {ledger.lines.length === 0 && (
                          <tr><td colSpan={5} className="px-4 py-6 text-center text-sm text-surface-500">Belum ada transaksi untuk akun ini</td></tr>
                        )}
                      </tbody>
                    </table>
                  </div>
                </div>
              )}
              {selectedAccount > 0 && !ledger && (
                <div className="flex items-center justify-center gap-2 text-surface-500 text-sm py-4">
                  <Loader2 size={16} className="animate-spin" /> Memuat buku besar...
                </div>
              )}
            </div>
          )}
        </div>
      )}

      {/* ============ MULTI-SATUAN (live from chain) ============ */}
      {activeTab === 'satuan' && !error && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-semibold text-white">Multi Satuan Produk</h3>
            <button className="btn-primary"><Plus size={18} />Tambah Satuan</button>
          </div>
          {loading ? (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {[0, 1, 2].map(i => <SkeletonCard key={i} />)}
            </div>
          ) : units.length === 0 ? (
            <EmptyState icon={<Ruler size={28} className="text-primary-400" />} title="Belum ada satuan terdaftar" hint="Buat satuan via `bizchaind tx pos create-unit`" />
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {units.map(u => (
                <div key={u.id} className="glass-card p-6">
                  <div className="flex items-start justify-between mb-3">
                    <div className="flex items-center gap-3">
                      <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-primary-500/10 to-accent-500/10 flex items-center justify-center">
                        <Ruler size={22} className="text-primary-400" />
                      </div>
                      <div>
                        <h3 className="text-base font-semibold text-white">{u.name}</h3>
                        <p className="text-xs text-surface-500 font-mono">{u.symbol}</p>
                      </div>
                    </div>
                    {u.is_base && <span className="badge badge-success">Satuan Dasar</span>}
                  </div>
                  <div className="space-y-1.5 text-xs text-surface-500">
                    <div className="flex justify-between"><span>Faktor Konversi</span><span className="text-surface-300">1 {u.symbol} = {u.conversion_factor} base</span></div>
                  </div>
                </div>
              ))}
            </div>
          )}

          {/* Products with units */}
          <div className="glass-card p-6">
            <h3 className="text-sm font-semibold text-white mb-3">Produk & Satuan Dasarnya</h3>
            {loading ? (
              <div className="space-y-3 animate-pulse">
                {[0, 1, 2].map(i => <div key={i} className="h-12 bg-white/5 rounded-lg" />)}
              </div>
            ) : products.length === 0 ? (
              <p className="text-sm text-surface-500">Belum ada produk. Buat via `bizchaind tx pos create-product`.</p>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b border-white/5">
                      <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">SKU</th>
                      <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Produk</th>
                      <th className="text-right text-xs text-surface-400 font-medium px-4 py-3">Stok</th>
                      <th className="text-right text-xs text-surface-400 font-medium px-4 py-3">Harga</th>
                      <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Satuan Dasar</th>
                    </tr>
                  </thead>
                  <tbody>
                    {products.map(p => {
                      const baseUnit = units.find(u => u.id === p.base_unit_id)
                      return (
                        <tr key={p.id} className="border-b border-white/5 hover:bg-white/5 transition-colors">
                          <td className="px-4 py-3 text-sm text-surface-400 font-mono">{p.sku}</td>
                          <td className="px-4 py-3 text-sm text-white">{p.name}</td>
                          <td className="px-4 py-3 text-sm text-surface-300 text-right">{p.stock}</td>
                          <td className="px-4 py-3 text-sm text-primary-400 text-right font-medium">{formatRupiah(p.price)}</td>
                          <td className="px-4 py-3 text-sm text-surface-300">{baseUnit ? `${baseUnit.name} (${baseUnit.symbol})` : '-'}</td>
                        </tr>
                      )
                    })}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        </div>
      )}

      {/* ============ HARGA LEVEL (live from chain) ============ */}
      {activeTab === 'hargalevel' && !error && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-semibold text-white">Laporan Harga per Level</h3>
            <button className="btn-primary"><Plus size={18} />Atur Harga Level</button>
          </div>
          {loading ? (
            <div className="space-y-3 animate-pulse">
              {[0, 1, 2].map(i => <div key={i} className="h-16 bg-white/5 rounded-lg" />)}
            </div>
          ) : priceLevels.length === 0 ? (
            <EmptyState icon={<Layers size={28} className="text-primary-400" />} title="Belum ada data harga level" hint="Set harga level saat `bizchaind tx pos create-product` (ecer/grosir/member)" />
          ) : (
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
              {priceLevels.map(item => (
                <div key={item.product_id} className="glass-card p-6">
                  <div className="flex items-start justify-between mb-3">
                    <div>
                      <h3 className="text-base font-semibold text-white">{item.product_name}</h3>
                      <p className="text-xs text-surface-500 font-mono">{item.sku}</p>
                    </div>
                    <span className="text-sm text-surface-300">Harga Dasar: <span className="text-white font-medium">{formatRupiah(item.base_price)}</span>/{item.base_unit}</span>
                  </div>
                  <div className="space-y-2">
                    {item.price_levels.map(pl => (
                      <div key={pl.level} className="flex items-center justify-between p-3 rounded-lg bg-white/5">
                        <div className="flex items-center gap-2">
                          <Tag size={14} className="text-primary-400" />
                          <span className="text-sm text-white capitalize">{pl.level}</span>
                          <span className="text-xs text-surface-500">min {pl.min_quantity}</span>
                        </div>
                        <span className="text-sm text-primary-400 font-medium">{formatRupiah(pl.price)}</span>
                      </div>
                    ))}
                    {item.price_levels.length === 0 && <p className="text-xs text-surface-500">Belum ada level harga untuk produk ini</p>}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {/* ============ BARANG GABUNGAN (live from chain) ============ */}
      {activeTab === 'gabungan' && !error && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-semibold text-white">Barang Gabungan (Bundle)</h3>
            <button className="btn-primary"><Plus size={18} />Buat Bundle</button>
          </div>
          {loading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {[0, 1].map(i => <SkeletonCard key={i} />)}
            </div>
          ) : bundleProducts.length === 0 ? (
            <EmptyState icon={<Boxes size={28} className="text-primary-400" />} title="Belum ada barang gabungan" hint="Buat bundle via `bizchaind tx pos create-product` dengan is_bundle=true" />
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {bundleProducts.map(b => (
                <div key={b.id} className="glass-card p-6">
                  <div className="flex items-start justify-between mb-3">
                    <div className="flex items-center gap-3">
                      <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-primary-500/10 to-accent-500/10 flex items-center justify-center">
                        <Boxes size={22} className="text-primary-400" />
                      </div>
                      <div>
                        <h3 className="text-base font-semibold text-white">{b.name}</h3>
                        <p className="text-xs text-surface-500 font-mono">{b.sku} • Stok: {b.stock}</p>
                      </div>
                    </div>
                    <span className="badge badge-info">Bundle</span>
                  </div>
                  <p className="text-xl font-bold text-primary-400 mb-3">{formatRupiah(b.price)}</p>
                  <div className="p-3 rounded-lg bg-white/5">
                    <p className="text-xs text-surface-400 mb-2 flex items-center gap-1"><Layers size={12} />Komponen Penyusun</p>
                    <div className="space-y-1.5">
                      {b.components.map(c => (
                        <div key={c.product_id} className="flex items-center justify-between text-sm">
                          <span className="text-surface-300">{productName(c.product_id)}</span>
                          <span className="text-white font-mono">× {c.quantity}</span>
                        </div>
                      ))}
                      {b.components.length === 0 && <p className="text-xs text-surface-500">Tidak ada komponen</p>}
                    </div>
                  </div>
                  {b.price_levels.length > 0 && (
                    <div className="flex gap-2 mt-3 flex-wrap">
                      {b.price_levels.map(pl => (
                        <span key={pl.level} className="text-[10px] px-2 py-1 rounded-lg bg-primary-500/10 text-primary-400 capitalize">{pl.level}: {formatRupiah(pl.price)}</span>
                      ))}
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {/* ============ MULTI-CABANG ============ */}
      {activeTab === 'cabang' && !error && (
        <div className="space-y-6">
          <div className="glass-card p-6">
            <h3 className="text-lg font-semibold text-white mb-4">Multi-Cabang Dashboard</h3>
            <p className="text-sm text-surface-400 mb-4">
              {selectedBranch === 'all'
                ? 'Menampilkan data seluruh cabang. Gunakan filter untuk melihat detail per cabang.'
                : `Menampilkan data cabang ${selectedBranch}. Gunakan dropdown di atas untuk mengganti cabang.`}
            </p>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="glass-card p-5">
                <div className="flex items-center gap-3 mb-3">
                  <div className="w-10 h-10 rounded-xl bg-primary-500/10 flex items-center justify-center">
                    <ShoppingBag size={20} className="text-primary-400" />
                  </div>
                  <span className="text-surface-400 text-sm">Total Produk</span>
                </div>
                <p className="text-2xl font-bold text-white">
                  {selectedBranch === 'all' ? products.length : products.filter(p => p.branch_id === selectedBranch).length}
                </p>
                <p className="text-xs text-surface-500 mt-2">Produk aktif</p>
              </div>
              <div className="glass-card p-5">
                <div className="flex items-center gap-3 mb-3">
                  <div className="w-10 h-10 rounded-xl bg-accent-500/10 flex items-center justify-center">
                    <TrendingUp size={20} className="text-accent-400" />
                  </div>
                  <span className="text-surface-400 text-sm">Total Transaksi</span>
                </div>
                <p className="text-2xl font-bold text-white">0</p>
                <p className="text-xs text-surface-500 mt-2">Segera tersedia</p>
              </div>
              <div className="glass-card p-5">
                <div className="flex items-center gap-3 mb-3">
                  <div className="w-10 h-10 rounded-xl bg-yellow-500/10 flex items-center justify-center">
                    <Map size={20} className="text-yellow-400" />
                  </div>
                  <span className="text-surface-400 text-sm">Cabang Terdaftar</span>
                </div>
                <p className="text-2xl font-bold text-white">4</p>
                <p className="text-xs text-surface-500 mt-2">Jakarta, Bandung, Semarang, Surabaya</p>
              </div>
            </div>
          </div>

          {selectedBranch === 'all' && (
            <div className="glass-card p-6">
              <h3 className="text-lg font-semibold text-white mb-4">Ringkasan Penjualan per Cabang</h3>
              <div className="space-y-3">
                {[
                  { id: 'JKT', name: 'Jakarta', sales: '45000000', transactions: 120 },
                  { id: 'BDG', name: 'Bandung', sales: '32000000', transactions: 85 },
                  { id: 'SMG', name: 'Semarang', sales: '28000000', transactions: 95 },
                  { id: 'SUR', name: 'Surabaya', sales: '35000000', transactions: 110 },
                ].map((branch, i) => (
                  <div key={branch.id} className="flex items-center justify-between p-4 rounded-xl bg-white/5">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 rounded-lg bg-primary-500/10 flex items-center justify-center">
                        <Map size={18} className="text-primary-400" />
                      </div>
                      <div>
                        <p className="text-sm font-medium text-white">{branch.name}</p>
                        <p className="text-xs text-surface-500">Cabang {branch.id}</p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="text-sm font-medium text-primary-400">{formatRupiah(branch.sales)}</p>
                      <p className="text-xs text-surface-500">{branch.transactions} transaksi</p>
                    </div>
                  </div>
                ))}
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

export default RetailDashboard
