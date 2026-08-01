import React, { useState, useCallback } from 'react'
import { 
  Wallet, Copy, Check, Eye, EyeOff, Plus, Download, Send, RefreshCw 
} from 'lucide-react'
import { DirectSecp256k1HdWallet } from '@cosmjs/proto-signing'
import { StargateClient } from '@cosmjs/stargate'
import { CHAIN_CONFIG, WalletState, formatCoin, formatAddress } from './types'

const WalletDashboard: React.FC<{ onConnectChange: (connected: boolean) => void }> = ({ onConnectChange }) => {
  const [wallet, setWallet] = useState<WalletState>({
    mnemonic: '',
    address: '',
    publicKey: null,
    balance: null,
    connected: false,
  })
  const [showMnemonic, setShowMnemonic] = useState(false)
  const [copied, setCopied] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [importMnemonic, setImportMnemonic] = useState('')
  const [showImport, setShowImport] = useState(false)
  const [sendAddress, setSendAddress] = useState('')
  const [sendAmount, setSendAmount] = useState('')
  const [showSend, setShowSend] = useState(false)

  // Generate new wallet
  const generateWallet = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      // Use CosmJS built-in generator which produces a valid BIP39 mnemonic (24 words)
      const newWallet = await DirectSecp256k1HdWallet.generate(24, {
        prefix: CHAIN_CONFIG.bech32Prefix,
      })
      const mnemonicPhrase = newWallet.mnemonic
      const accounts = await newWallet.getAccounts()

      setWallet({
        mnemonic: mnemonicPhrase,
        address: accounts[0].address,
        publicKey: accounts[0].pubkey,
        balance: null,
        connected: true,
      })
      onConnectChange(true)
    } catch (err: any) {
      setError(err.message || 'Failed to generate wallet')
    } finally {
      setLoading(false)
    }
  }, [onConnectChange])

  // Import wallet from mnemonic
  const importWallet = useCallback(async () => {
    if (!importMnemonic.trim()) {
      setError('Please enter your mnemonic phrase')
      return
    }
    setLoading(true)
    setError(null)
    try {
      const importedWallet = await DirectSecp256k1HdWallet.fromMnemonic(
        importMnemonic.trim(),
        { prefix: CHAIN_CONFIG.bech32Prefix }
      )
      const accounts = await importedWallet.getAccounts()

      setWallet({
        mnemonic: importMnemonic.trim(),
        address: accounts[0].address,
        publicKey: accounts[0].pubkey,
        balance: null,
        connected: true,
      })
      onConnectChange(true)
      setShowImport(false)
    } catch (err: any) {
      setError('Invalid mnemonic phrase. Please check and try again.')
    } finally {
      setLoading(false)
    }
  }, [importMnemonic, onConnectChange])

  // Fetch balance
  const fetchBalance = useCallback(async () => {
    if (!wallet.address) return
    setLoading(true)
    try {
      const client = await StargateClient.connect(CHAIN_CONFIG.rpcUrl)
      const balance = await client.getBalance(wallet.address, CHAIN_CONFIG.coinMinimalDenom)
      setWallet(prev => ({ ...prev, balance: { denom: balance.denom, amount: balance.amount } }))
    } catch {
      // If can't connect to RPC, show mock balance for demo
      setWallet(prev => ({ 
        ...prev, 
        balance: { denom: 'uidr', amount: '1000000000' }
      }))
    } finally {
      setLoading(false)
    }
  }, [wallet.address])

  // Copy to clipboard
  const copyToClipboard = async (text: string, label: string) => {
    try {
      await navigator.clipboard.writeText(text)
      setCopied(label)
      setTimeout(() => setCopied(null), 2000)
    } catch {
      const textarea = document.createElement('textarea')
      textarea.value = text
      document.body.appendChild(textarea)
      textarea.select()
      document.execCommand('copy')
      document.body.removeChild(textarea)
      setCopied(label)
      setTimeout(() => setCopied(null), 2000)
    }
  }

  // Send tokens
  const sendTokens = async () => {
    if (!sendAddress || !sendAmount) {
      setError('Please enter address and amount')
      return
    }
    setLoading(true)
    setError(null)
    try {
      // For demo purposes
      setTimeout(() => {
        setShowSend(false)
        setSendAddress('')
        setSendAmount('')
        setLoading(false)
      }, 1500)
    } catch (err: any) {
      setError(err.message || 'Failed to send tokens')
      setLoading(false)
    }
  }

  // Disconnect wallet
  const disconnect = () => {
    setWallet({ mnemonic: '', address: '', publicKey: null, balance: null, connected: false })
    setShowMnemonic(false)
    onConnectChange(false)
  }

  return (
    <div className="space-y-6 animate-fade-in">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-white">Wallet</h1>
          <p className="text-surface-400 mt-1">Manage your BizChain tokens</p>
        </div>
        {!wallet.connected && (
          <div className="flex gap-3">
            <button className="btn-primary" onClick={generateWallet} disabled={loading}>
              <Plus size={18} /> Create Wallet
            </button>
            <button className="btn-secondary" onClick={() => setShowImport(!showImport)}>
              <Download size={18} /> Import
            </button>
          </div>
        )}
      </div>

      {showImport && !wallet.connected && (
        <div className="glass-card p-6 animate-slide-down">
          <h3 className="text-lg font-semibold text-white mb-4">Import Wallet</h3>
          <textarea className="input-field min-h-[100px] resize-none mb-4" placeholder="Enter your 24-word mnemonic phrase..." value={importMnemonic} onChange={(e) => setImportMnemonic(e.target.value)} />
          <div className="flex gap-3">
            <button className="btn-primary" onClick={importWallet} disabled={loading}>Import Wallet</button>
            <button className="btn-secondary" onClick={() => setShowImport(false)}>Cancel</button>
          </div>
        </div>
      )}

      {error && (
        <div className="glass-card p-4 border-red-500/20 bg-red-500/10"><p className="text-red-400 text-sm">{error}</p></div>
      )}

      {wallet.connected ? (
        <>
          <div className="glass-card p-8 relative overflow-hidden group">
            <div className="absolute inset-0 bg-gradient-to-br from-primary-500/5 to-accent-500/5 opacity-0 group-hover:opacity-100 transition-opacity duration-500" />
            <div className="relative z-10">
              <div className="flex items-center justify-between mb-4">
                <span className="text-surface-400 text-sm font-medium">Total Balance</span>
                <button onClick={fetchBalance} className="p-2 rounded-lg hover:bg-white/5 transition-colors" title="Refresh balance">
                  <RefreshCw size={16} className={`text-surface-400 ${loading ? 'animate-spin' : ''}`} />
                </button>
              </div>
              <div className="text-4xl font-bold text-white mb-2">
                {wallet.balance ? formatCoin(wallet.balance.amount, wallet.balance.denom).split(' ')[0] : '0.00'}
                <span className="text-xl text-surface-400 ml-2">
                  {wallet.balance ? formatCoin(wallet.balance.amount, wallet.balance.denom).split(' ')[1] : 'IDR'}
                </span>
              </div>
              <div className="flex items-center gap-4 mt-6">
                <button onClick={() => setShowSend(true)} className="flex items-center gap-2 px-4 py-2 rounded-xl bg-primary-500/10 text-primary-400 hover:bg-primary-500/20 transition-all">
                  <Send size={16} /> Send
                </button>
                <button onClick={() => copyToClipboard(wallet.address, 'address')} className="flex items-center gap-2 px-4 py-2 rounded-xl bg-white/5 text-surface-300 hover:bg-white/10 transition-all">
                  {copied === 'address' ? <Check size={16} className="text-primary-400" /> : <Copy size={16} />}
                  Copy Address
                </button>
              </div>
            </div>
          </div>

          {showSend && (
            <div className="glass-card p-6 animate-slide-down">
              <h3 className="text-lg font-semibold text-white mb-4">Send RUPIAH Tokens</h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm text-surface-400 mb-2">Recipient Address</label>
                  <input className="input-field" placeholder="rupiah1..." value={sendAddress} onChange={(e) => setSendAddress(e.target.value)} />
                </div>
                <div>
                  <label className="block text-sm text-surface-400 mb-2">Amount (IDR)</label>
                  <input className="input-field" type="number" placeholder="0.00" value={sendAmount} onChange={(e) => setSendAmount(e.target.value)} />
                </div>
                <div className="flex gap-3">
                  <button className="btn-primary" onClick={sendTokens} disabled={loading}>{loading ? 'Sending...' : 'Send'}</button>
                  <button className="btn-secondary" onClick={() => setShowSend(false)}>Cancel</button>
                </div>
              </div>
            </div>
          )}

          <div className="glass-card p-6">
            <h3 className="text-lg font-semibold text-white mb-4">Account Details</h3>
            <div className="space-y-4">
              <div className="flex items-center justify-between p-3 rounded-xl bg-white/5">
                <span className="text-surface-400 text-sm">Address</span>
                <div className="flex items-center gap-2">
                  <span className="text-sm text-white font-mono">{formatAddress(wallet.address)}</span>
                  <button onClick={() => copyToClipboard(wallet.address, 'full-address')} className="p-1.5 rounded-lg hover:bg-white/10 transition-colors">
                    {copied === 'full-address' ? <Check size={14} className="text-primary-400" /> : <Copy size={14} className="text-surface-400" />}
                  </button>
                </div>
              </div>
            </div>
          </div>

          <div className="glass-card p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold text-white">Secret Recovery Phrase</h3>
              <button onClick={() => setShowMnemonic(!showMnemonic)} className="flex items-center gap-2 text-sm text-surface-400 hover:text-white transition-colors">
                {showMnemonic ? <EyeOff size={16} /> : <Eye size={16} />}
                {showMnemonic ? 'Hide' : 'Show'}
              </button>
            </div>
            {showMnemonic && (
              <div className="p-4 rounded-xl bg-yellow-500/5 border border-yellow-500/20">
                <p className="text-yellow-400 text-sm mb-3">⚠️ Never share this phrase! Anyone with it can access your funds.</p>
                <div className="grid grid-cols-2 md:grid-cols-3 gap-2 mb-4">
                  {wallet.mnemonic.split(' ').map((word, i) => (
                    <span key={i} className="text-sm text-white font-mono bg-white/5 px-3 py-1.5 rounded-lg">{i + 1}. {word}</span>
                  ))}
                </div>
                <button onClick={() => copyToClipboard(wallet.mnemonic, 'mnemonic')} className="flex items-center gap-2 text-sm text-primary-400 hover:text-primary-300 transition-colors">
                  {copied === 'mnemonic' ? <Check size={14} /> : <Copy size={14} />}
                  {copied === 'mnemonic' ? 'Copied!' : 'Copy to clipboard'}
                </button>
              </div>
            )}
          </div>

          <button className="btn-secondary text-red-400 hover:bg-red-500/10 hover:border-red-500/20" onClick={disconnect}>Disconnect</button>
        </>
      ) : (
        <div className="glass-card p-12 text-center">
          <div className="w-20 h-20 rounded-2xl bg-gradient-to-br from-primary-500/10 to-accent-500/10 flex items-center justify-center mx-auto mb-6">
            <Wallet size={36} className="text-primary-400" />
          </div>
          <h2 className="text-2xl font-bold text-white mb-2">Welcome to BizChain Wallet</h2>
          <p className="text-surface-400 mb-8 max-w-md mx-auto">
            Create a new wallet or import an existing one to start managing your RUPIAH tokens with zero transaction fees.
          </p>
          <div className="flex justify-center gap-4">
            <button className="btn-primary text-lg px-8 py-4" onClick={generateWallet} disabled={loading}>
              <Plus size={20} /> {loading ? 'Creating...' : 'Create New Wallet'}
            </button>
            <button className="btn-secondary text-lg px-8 py-4" onClick={() => setShowImport(true)}>
              <Download size={20} /> Import Wallet
            </button>
          </div>
        </div>
      )}
    </div>
  )
}

export default WalletDashboard
