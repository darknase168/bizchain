import React, { useState, useEffect, useCallback } from 'react'
import { Routes, Route, Link, useLocation } from 'react-router-dom'
import { Wallet, ShoppingCart, Package, BarChart3, Settings, Menu, X, Zap, Globe, Shield, Network, GraduationCap, Store, Plane, Users } from 'lucide-react'
import WalletDashboard from './Wallet'
import PosDashboard from './PosDashboard'
import MlmDashboard from './MlmDashboard'
import EducationDashboard from './EducationDashboard'
import RetailDashboard from './RetailDashboard'
import HajiUmrohDashboard from './HajiUmrohDashboard'
import AgenEkosistemDashboard from './AgenEkosistemDashboard'

const NavLink: React.FC<{ to: string; icon: React.ReactNode; label: string; active?: boolean }> = ({ to, icon, label, active }) => (
  <Link
    to={to}
    className={`flex items-center gap-3 px-4 py-3 rounded-xl transition-all duration-200 ${
      active
        ? 'bg-primary-500/10 text-primary-400 border border-primary-500/20'
        : 'text-surface-400 hover:text-white hover:bg-white/5'
    }`}
  >
    {icon}
    <span className="font-medium">{label}</span>
  </Link>
)

const App: React.FC = () => {
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const [connected, setConnected] = useState(false)
  const location = useLocation()

  const navItems = [
    { to: '/', icon: <Wallet size={20} />, label: 'Wallet' },
    { to: '/pos', icon: <ShoppingCart size={20} />, label: 'POS Dashboard' },
    { to: '/retail', icon: <Store size={20} />, label: 'Retail' },
    { to: '/hajiumroh', icon: <Plane size={20} />, label: 'Haji & Umroh' },
    { to: '/ekosistem', icon: <Users size={20} />, label: 'Ekosistem Agen' },
    { to: '/mlm', icon: <Network size={20} />, label: 'MLM Network' },
    { to: '/education', icon: <GraduationCap size={20} />, label: 'Education' },
    { to: '/products', icon: <Package size={20} />, label: 'Products' },
    { to: '/analytics', icon: <BarChart3 size={20} />, label: 'Analytics' },
    { to: '/settings', icon: <Settings size={20} />, label: 'Settings' },
  ]

  const isActive = (path: string) => {
    if (path === '/') return location.pathname === '/'
    return location.pathname.startsWith(path)
  }

  return (
    <div className="min-h-screen flex">
      {/* Sidebar */}
      <aside className={`
        fixed lg:static inset-y-0 left-0 z-50 w-72
        bg-surface-900/50 backdrop-blur-2xl border-r border-white/5
        transform transition-transform duration-300 lg:transform-none
        ${sidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'}
      `}>
        {/* Logo */}
        <div className="p-6 border-b border-white/5">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-primary-500 to-accent-600 flex items-center justify-center">
              <Zap size={22} className="text-white" />
            </div>
            <div>
              <h1 className="text-lg font-bold text-white">BizChain</h1>
              <p className="text-xs text-surface-400">Feeless Blockchain</p>
            </div>
          </div>

          {/* Network Status */}
          <div className="mt-4 flex items-center gap-2 px-3 py-2 rounded-lg bg-white/5">
            <span className={`status-dot ${connected ? 'active' : 'inactive'}`} />
            <span className="text-xs text-surface-400">
              {connected ? 'Connected to BizChain' : 'Not connected'}
            </span>
          </div>
        </div>

        {/* Navigation */}
        <nav className="p-4 space-y-1">
          {navItems.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              icon={item.icon}
              label={item.label}
              active={isActive(item.to)}
            />
          ))}
        </nav>

        {/* Footer */}
        <div className="absolute bottom-0 left-0 right-0 p-6 border-t border-white/5">
          <div className="flex items-center gap-2 text-xs text-surface-500">
            <Shield size={12} />
            <span>Secured by Cosmos SDK</span>
          </div>
        </div>
      </aside>

      {/* Overlay */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 bg-black/50 z-40 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* Main Content */}
      <main className="flex-1 min-h-screen">
        {/* Top Bar */}
        <header className="sticky top-0 z-30 bg-surface-900/80 backdrop-blur-xl border-b border-white/5">
          <div className="flex items-center justify-between px-6 py-4">
            <button
              className="lg:hidden p-2 rounded-lg hover:bg-white/5"
              onClick={() => setSidebarOpen(!sidebarOpen)}
            >
              {sidebarOpen ? <X size={20} /> : <Menu size={20} />}
            </button>

            <div className="flex items-center gap-4">
              <div className="flex items-center gap-2 text-sm text-surface-400">
                <Globe size={14} />
                <span>bizchain-1</span>
              </div>
              <div className="h-4 w-px bg-white/10" />
              <div className="flex items-center gap-2">
                <span className="text-xs text-surface-500">0 FEE</span>
                <span className="badge-success">Active</span>
              </div>
            </div>
          </div>
        </header>

        {/* Page Content */}
        <div className="p-6">
          <Routes>
            <Route path="/" element={<WalletDashboard onConnectChange={setConnected} />} />
            <Route path="/pos" element={<PosDashboard />} />
            <Route path="/retail" element={<RetailDashboard />} />
            <Route path="/hajiumroh" element={<HajiUmrohDashboard />} />
            <Route path="/ekosistem" element={<AgenEkosistemDashboard />} />
            <Route path="/mlm" element={<MlmDashboard />} />
            <Route path="/education" element={<EducationDashboard />} />
            <Route path="/products" element={<PosDashboard />} />
            <Route path="*" element={
              <div className="flex items-center justify-center h-96">
                <div className="text-center">
                  <Zap size={48} className="mx-auto text-primary-500 mb-4" />
                  <h2 className="text-2xl font-bold text-white mb-2">Coming Soon</h2>
                  <p className="text-surface-400">This feature is under development</p>
                </div>
              </div>
            } />
          </Routes>
        </div>
      </main>
    </div>
  )
}

export default App
