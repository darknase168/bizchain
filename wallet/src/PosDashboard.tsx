import React, { useState } from 'react'
import {
  ShoppingCart, Package, Plus, Search, Minus, 
  Receipt, TrendingUp, CheckCircle,
  BarChart3
} from 'lucide-react'
import { Product, Transaction, CartItem, CreateProductForm, formatCoin, formatDate } from './types'

const initialProducts: Product[] = [
  { id: 1, name: 'Indomie Goreng', description: 'Instant noodles', price: '3500000', sku: 'IDM-001', category: 'Makanan', image_url: '', stock: 150, owner: 'rupiah1...', active: true, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
  { id: 2, name: 'Aqua 600ml', description: 'Mineral water', price: '1500000', sku: 'AQU-001', category: 'Minuman', image_url: '', stock: 200, owner: 'rupiah1...', active: true, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
  { id: 3, name: 'Roti Tawar', description: 'White bread', price: '5000000', sku: 'ROT-001', category: 'Makanan', image_url: '', stock: 50, owner: 'rupiah1...', active: true, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
  { id: 4, name: 'Susu UHT 1L', description: 'UHT milk', price: '12000000', sku: 'SUS-001', category: 'Minuman', image_url: '', stock: 80, owner: 'rupiah1...', active: true, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
  { id: 5, name: 'Buku Tulis', description: 'Notebook 40 pages', price: '2000000', sku: 'BUK-001', category: 'Pendidikan', image_url: '', stock: 300, owner: 'rupiah1...', active: true, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
  { id: 6, name: 'Pensil 2B', description: 'Drawing pencil', price: '1000000', sku: 'PEN-001', category: 'Pendidikan', image_url: '', stock: 500, owner: 'rupiah1...', active: true, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
]

const PosDashboard: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'pos' | 'products'>('pos')
  const [products, setProducts] = useState<Product[]>(initialProducts)
  const [cart, setCart] = useState<CartItem[]>([])
  const [searchTerm, setSearchTerm] = useState('')
  const [categoryFilter, setCategoryFilter] = useState<string>('all')
  const [showCreateProduct, setShowCreateProduct] = useState(false)
  const [transactions, setTransactions] = useState<Transaction[]>([])
  const [customerAddress, setCustomerAddress] = useState('')
  const [note, setNote] = useState('')
  const [showSuccess, setShowSuccess] = useState(false)
  const [newProduct, setNewProduct] = useState<CreateProductForm>({ name: '', description: '', price: '', sku: '', category: '', stock: '' })

  const filteredProducts = products.filter(p => {
    const matchesSearch = p.name.toLowerCase().includes(searchTerm.toLowerCase()) || p.sku.toLowerCase().includes(searchTerm.toLowerCase())
    const matchesCategory = categoryFilter === 'all' || p.category === categoryFilter
    return matchesSearch && matchesCategory && p.active
  })

  const categories = [...new Set(products.map(p => p.category))]

  const formatTotal = (price: string, qty: number): string => {
    const priceNum = parseInt(price) || 0
    return (priceNum * qty).toString()
  }

  const addToCart = (product: Product) => {
    setCart(prev => {
      const existing = prev.find(item => item.productId === product.id)
      if (existing) {
        return prev.map(item =>
          item.productId === product.id
            ? { ...item, quantity: item.quantity + 1, total: formatTotal(item.price, item.quantity + 1) }
            : item
        )
      }
      return [...prev, { productId: product.id, productName: product.name, quantity: 1, price: product.price, total: formatTotal(product.price, 1) }]
    })
  }

  const removeFromCart = (productId: number) => {
    setCart(prev => {
      const existing = prev.find(item => item.productId === productId)
      if (existing && existing.quantity > 1) {
        return prev.map(item =>
          item.productId === productId
            ? { ...item, quantity: item.quantity - 1, total: formatTotal(item.price, item.quantity - 1) }
            : item
        )
      }
      return prev.filter(item => item.productId !== productId)
    })
  }

  const cartTotal = cart.reduce((sum, item) => sum + (parseInt(item.total) || 0), 0)
  const clearCart = () => setCart([])

  const createTransaction = () => {
    if (cart.length === 0) return
    const newTx: Transaction = {
      id: transactions.length + 1,
      seller: 'rupiah1...',
      customer_address: customerAddress || 'Anonymous',
      items: cart.map(item => ({ product_id: item.productId, quantity: item.quantity, price: item.price })),
      total: cartTotal.toString(),
      status: 'completed',
      note: note,
      created_at: new Date().toISOString(),
      branch_id: 'main'
    }
    setTransactions(prev => [newTx, ...prev])
    clearCart()
    setCustomerAddress('')
    setNote('')
    setShowSuccess(true)
    setTimeout(() => setShowSuccess(false), 3000)
  }

  const handleCreateProduct = () => {
    if (!newProduct.name || !newProduct.price) return
    const newId = products.length + 1
    const product: Product = {
      id: newId,
      name: newProduct.name,
      description: newProduct.description,
      price: newProduct.price,
      sku: newProduct.sku || `PRD-${String(newId).padStart(3, '0')}`,
      category: newProduct.category || 'General',
      image_url: '',
      stock: parseInt(newProduct.stock) || 0,
      owner: 'rupiah1...',
      active: true,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    }
    setProducts(prev => [...prev, product])
    setShowCreateProduct(false)
    setNewProduct({ name: '', description: '', price: '', sku: '', category: '', stock: '' })
  }

  const totalSales = transactions.reduce((sum, tx) => sum + parseInt(tx.total), 0)
  const todaySales = transactions.filter(tx => new Date(tx.created_at).toDateString() === new Date().toDateString()).reduce((sum, tx) => sum + parseInt(tx.total), 0)

  return (
    <div className="space-y-6 animate-fade-in">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-white">POS Dashboard</h1>
          <p className="text-surface-400 mt-1">Point of Sales - Retail Management</p>
        </div>
        <div className="flex gap-2">
          <button className={`px-4 py-2 rounded-xl text-sm font-medium transition-all ${activeTab === 'pos' ? 'bg-primary-500/20 text-primary-400 border border-primary-500/30' : 'bg-white/5 text-surface-400 hover:text-white'}`} onClick={() => setActiveTab('pos')}>
            <ShoppingCart size={16} className="inline mr-2" />POS
          </button>
          <button className={`px-4 py-2 rounded-xl text-sm font-medium transition-all ${activeTab === 'products' ? 'bg-primary-500/20 text-primary-400 border border-primary-500/30' : 'bg-white/5 text-surface-400 hover:text-white'}`} onClick={() => setActiveTab('products')}>
            <Package size={16} className="inline mr-2" />Products
          </button>
        </div>
      </div>

      {showSuccess && (
        <div className="glass-card p-4 border-primary-500/20 bg-primary-500/10 flex items-center gap-3 animate-slide-down">
          <CheckCircle size={20} className="text-primary-400" />
          <span className="text-primary-300">Transaction completed successfully! 0 fees charged.</span>
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="glass-card p-5">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-10 h-10 rounded-xl bg-primary-500/10 flex items-center justify-center"><TrendingUp size={20} className="text-primary-400" /></div>
            <span className="text-surface-400 text-sm">Total Sales</span>
          </div>
          <p className="text-2xl font-bold text-white">{formatCoin(totalSales.toString(), 'uidr')}</p>
        </div>
        <div className="glass-card p-5">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-10 h-10 rounded-xl bg-accent-500/10 flex items-center justify-center"><BarChart3 size={20} className="text-accent-400" /></div>
            <span className="text-surface-400 text-sm">Today's Sales</span>
          </div>
          <p className="text-2xl font-bold text-white">{formatCoin(todaySales.toString(), 'uidr')}</p>
        </div>
        <div className="glass-card p-5">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-10 h-10 rounded-xl bg-blue-500/10 flex items-center justify-center"><ShoppingCart size={20} className="text-blue-400" /></div>
            <span className="text-surface-400 text-sm">Transactions</span>
          </div>
          <p className="text-2xl font-bold text-white">{transactions.length}</p>
        </div>
      </div>

      {activeTab === 'pos' ? (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2 space-y-4">
            <div className="flex gap-3">
              <div className="relative flex-1">
                <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-surface-400" />
                <input className="input-field pl-10" placeholder="Search products..." value={searchTerm} onChange={(e) => setSearchTerm(e.target.value)} />
              </div>
              <select className="input-field w-auto" value={categoryFilter} onChange={(e) => setCategoryFilter(e.target.value)}>
                <option value="all">All Categories</option>
                {categories.map(cat => (<option key={cat} value={cat}>{cat}</option>))}
              </select>
            </div>
            <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
              {filteredProducts.map(product => (
                <button key={product.id} onClick={() => addToCart(product)} className="glass-card-hover p-4 text-left">
                  <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-primary-500/10 to-accent-500/10 flex items-center justify-center mb-3">
                    <Package size={22} className="text-primary-400" />
                  </div>
                  <h3 className="font-medium text-white text-sm mb-1">{product.name}</h3>
                  <p className="text-primary-400 font-semibold text-sm mb-2">{formatCoin(product.price, 'uidr')}</p>
                  <span className="text-xs text-surface-500">Stock: {product.stock}</span>
                </button>
              ))}
            </div>
          </div>
          <div className="glass-card p-4 h-fit sticky top-24">
            <div className="flex items-center justify-between mb-4">
              <h3 className="font-semibold text-white flex items-center gap-2"><ShoppingCart size={18} className="text-primary-400" />Cart ({cart.length})</h3>
              {cart.length > 0 && (<button onClick={clearCart} className="text-xs text-red-400 hover:text-red-300 transition-colors">Clear</button>)}
            </div>
            {cart.length === 0 ? (
              <div className="text-center py-8">
                <ShoppingCart size={32} className="mx-auto text-surface-500 mb-2" />
                <p className="text-surface-500 text-sm">Cart is empty</p>
                <p className="text-surface-600 text-xs mt-1">Click products to add</p>
              </div>
            ) : (
              <>
                <div className="space-y-2 mb-4 max-h-64 overflow-y-auto">
                  {cart.map(item => (
                    <div key={item.productId} className="flex items-center justify-between p-2 rounded-xl bg-white/5">
                      <div className="flex-1 min-w-0">
                        <p className="text-sm text-white truncate">{item.productName}</p>
                        <p className="text-xs text-primary-400">{formatCoin(item.total, 'uidr')}</p>
                      </div>
                      <div className="flex items-center gap-1 ml-2">
                        <button onClick={() => removeFromCart(item.productId)} className="p-1 rounded-lg hover:bg-white/10 transition-colors"><Minus size={14} className="text-surface-400" /></button>
                        <span className="text-sm text-white w-6 text-center">{item.quantity}</span>
                        <button onClick={() => addToCart(products.find(p => p.id === item.productId)!)} className="p-1 rounded-lg hover:bg-white/10 transition-colors"><Plus size={14} className="text-surface-400" /></button>
                      </div>
                    </div>
                  ))}
                </div>
                <input className="input-field text-sm mb-2" placeholder="Customer address (optional)" value={customerAddress} onChange={(e) => setCustomerAddress(e.target.value)} />
                <input className="input-field text-sm mb-4" placeholder="Transaction note (optional)" value={note} onChange={(e) => setNote(e.target.value)} />
                <div className="flex items-center justify-between p-3 rounded-xl bg-primary-500/10 mb-4">
                  <span className="text-sm text-surface-300">Total</span>
                  <span className="text-lg font-bold text-primary-400">{formatCoin(cartTotal.toString(), 'uidr')}</span>
                </div>
                <button onClick={createTransaction} className="btn-primary w-full justify-center" disabled={cart.length === 0}>
                  <Receipt size={18} />Complete Transaction (0 Fee)
                </button>
              </>
            )}
          </div>
        </div>
      ) : (
        <div className="space-y-4">
          <div className="flex justify-end">
            <button className="btn-primary" onClick={() => setShowCreateProduct(true)}><Plus size={18} />Add Product</button>
          </div>
          {showCreateProduct && (
            <div className="glass-card p-6 animate-slide-down">
              <h3 className="text-lg font-semibold text-white mb-4">Add New Product</h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="block text-sm text-surface-400 mb-2">Product Name *</label>
                  <input className="input-field" placeholder="Product name" value={newProduct.name} onChange={(e) => setNewProduct(prev => ({ ...prev, name: e.target.value }))} />
                </div>
                <div>
                  <label className="block text-sm text-surface-400 mb-2">SKU</label>
                  <input className="input-field" placeholder="SKU-001" value={newProduct.sku} onChange={(e) => setNewProduct(prev => ({ ...prev, sku: e.target.value }))} />
                </div>
                <div>
                  <label className="block text-sm text-surface-400 mb-2">Price (in uidr) *</label>
                  <input className="input-field" type="number" placeholder="1000000" value={newProduct.price} onChange={(e) => setNewProduct(prev => ({ ...prev, price: e.target.value }))} />
                </div>
                <div>
                  <label className="block text-sm text-surface-400 mb-2">Category</label>
                  <input className="input-field" placeholder="Makanan, Minuman, Pendidikan" value={newProduct.category} onChange={(e) => setNewProduct(prev => ({ ...prev, category: e.target.value }))} />
                </div>
                <div>
                  <label className="block text-sm text-surface-400 mb-2">Initial Stock</label>
                  <input className="input-field" type="number" placeholder="0" value={newProduct.stock} onChange={(e) => setNewProduct(prev => ({ ...prev, stock: e.target.value }))} />
                </div>
                <div className="md:col-span-2">
                  <label className="block text-sm text-surface-400 mb-2">Description</label>
                  <textarea className="input-field min-h-[80px] resize-none" placeholder="Product description" value={newProduct.description} onChange={(e) => setNewProduct(prev => ({ ...prev, description: e.target.value }))} />
                </div>
              </div>
              <div className="flex gap-3">
                <button className="btn-primary" onClick={handleCreateProduct}>Add Product</button>
                <button className="btn-secondary" onClick={() => setShowCreateProduct(false)}>Cancel</button>
              </div>
            </div>
          )}
          <div className="glass-card overflow-hidden">
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-white/5">
                    <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">ID</th>
                    <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Name</th>
                    <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">SKU</th>
                    <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Category</th>
                    <th className="text-right text-xs text-surface-400 font-medium px-4 py-3">Price</th>
                    <th className="text-right text-xs text-surface-400 font-medium px-4 py-3">Stock</th>
                    <th className="text-center text-xs text-surface-400 font-medium px-4 py-3">Status</th>
                  </tr>
                </thead>
                <tbody>
                  {products.map(product => (
                    <tr key={product.id} className="border-b border-white/5 hover:bg-white/5 transition-colors">
                      <td className="px-4 py-3 text-sm text-surface-300">{product.id}</td>
                      <td className="px-4 py-3">
                        <div className="flex items-center gap-2">
                          <div className="w-8 h-8 rounded-lg bg-primary-500/10 flex items-center justify-center"><Package size={14} className="text-primary-400" /></div>
                          <div><p className="text-sm text-white">{product.name}</p><p className="text-xs text-surface-500">{product.description}</p></div>
                        </div>
                      </td>
                      <td className="px-4 py-3 text-sm text-surface-400 font-mono">{product.sku}</td>
                      <td className="px-4 py-3"><span className="badge-info">{product.category}</span></td>
                      <td className="px-4 py-3 text-sm text-primary-400 text-right font-medium">{formatCoin(product.price, 'uidr')}</td>
                      <td className="px-4 py-3 text-sm text-surface-300 text-right">{product.stock}</td>
                      <td className="px-4 py-3 text-center"><span className={`badge ${product.active ? 'badge-success' : ''}`}>{product.active ? 'Active' : 'Inactive'}</span></td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}

      {transactions.length > 0 && activeTab === 'pos' && (
        <div className="glass-card p-6">
          <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2"><Receipt size={18} className="text-primary-400" />Recent Transactions</h3>
          <div className="space-y-2">
            {transactions.slice(0, 5).map(tx => (
              <div key={tx.id} className="flex items-center justify-between p-3 rounded-xl bg-white/5">
                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 rounded-lg bg-primary-500/10 flex items-center justify-center"><Receipt size={14} className="text-primary-400" /></div>
                  <div><p className="text-sm text-white">Transaction #{tx.id}</p><p className="text-xs text-surface-500">{tx.items.length} items • {formatDate(tx.created_at)}</p></div>
                </div>
                <div className="text-right">
                  <p className="text-sm text-primary-400 font-medium">{formatCoin(tx.total, 'uidr')}</p>
                  <span className="badge-success text-xs">Completed</span>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}

export default PosDashboard
