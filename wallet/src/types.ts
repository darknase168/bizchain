// BizChain Wallet Types

export const CHAIN_CONFIG = {
  chainId: 'bizchain-1',
  chainName: 'BizChain',
  rpcUrl: 'http://localhost:26657',
  restUrl: 'http://localhost:1317',
  bech32Prefix: 'rupiah',
  coinDenom: 'RUPIAH',
  coinMinimalDenom: 'uidr',
  coinDecimals: 6,
} as const;

export interface WalletState {
  mnemonic: string;
  address: string;
  publicKey: Uint8Array | null;
  balance: Coin | null;
  connected: boolean;
}

export interface Coin {
  denom: string;
  amount: string;
}

export interface Product {
  id: number;
  name: string;
  description: string;
  price: string;
  sku: string;
  category: string;
  image_url: string;
  stock: number;
  owner: string;
  active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Transaction {
  id: number;
  seller: string;
  customer_address: string;
  items: TransactionItem[];
  total: string;
  status: string;
  note: string;
  created_at: string;
  branch_id: string;
}

export interface TransactionItem {
  product_id: number;
  quantity: number;
  price: string;
}

export interface POSState {
  products: Product[];
  transactions: Transaction[];
  selectedProduct: Product | null;
  cart: CartItem[];
  loading: boolean;
  error: string | null;
}

export interface CartItem {
  productId: number;
  productName: string;
  quantity: number;
  price: string;
  total: string;
}

export interface CreateProductForm {
  name: string;
  description: string;
  price: string;
  sku: string;
  category: string;
  stock: string;
}

export interface CreateTransactionForm {
  customerAddress: string;
  note: string;
}

// ==================== MLM Types ====================

export interface MLMMember {
  id: number;
  address: string;
  name: string;
  sponsor_id: number | null;
  sponsor_name: string;
  rank: MLMRank;
  personal_volume: string;
  group_volume: string;
  total_downline: number;
  direct_downline: number;
  status: 'active' | 'inactive' | 'suspended';
  joined_at: string;
  last_active: string;
}

export type MLMRank = 'Bronze' | 'Silver' | 'Gold' | 'Platinum' | 'Diamond' | 'Crown';

export interface MLMCommission {
  id: number;
  member_id: number;
  member_name: string;
  from_member: string;
  type: CommissionType;
  amount: string;
  level: number;
  transaction_id: number;
  status: 'pending' | 'paid' | 'cancelled';
  created_at: string;
}

export type CommissionType = 'sponsor_bonus' | 'pairing_bonus' | 'matching_bonus' | 'level_bonus' | 'rank_bonus' | 'retail_profit';

export interface MLMTreeNode {
  member: MLMMember;
  children: MLMTreeNode[];
  level: number;
  position: 'left' | 'right' | 'center';
}

export interface MLMSettings {
  sponsor_bonus_percent: number;
  pairing_bonus_percent: number;
  matching_bonus_levels: number[];
  level_bonus_percents: number[];
  rank_requirements: RankRequirement[];
  max_withdrawal: string;
  min_withdrawal: string;
}

export interface RankRequirement {
  rank: MLMRank;
  personal_volume: string;
  group_volume: string;
  direct_downline: number;
  reward: string;
}

// ==================== Education Types ====================

export interface Student {
  id: number;
  nis: string;
  name: string;
  class: string;
  major: string;
  parent_name: string;
  parent_phone: string;
  address: string;
  wallet_address: string;
  status: 'active' | 'graduated' | 'transferred' | 'dropped';
  enrolled_at: string;
}

export interface Course {
  id: number;
  code: string;
  name: string;
  description: string;
  credits: number;
  teacher: string;
  schedule: string;
  semester: string;
  max_students: number;
  enrolled_students: number;
  fee: string;
  status: 'open' | 'closed' | 'completed';
}

export interface TuitionPayment {
  id: number;
  student_id: number;
  student_name: string;
  type: 'spp' | 'uang_gedung' | 'praktikum' | 'wisuda' | 'lainnya';
  amount: string;
  semester: string;
  academic_year: string;
  status: 'pending' | 'paid' | 'overdue' | 'cancelled';
  due_date: string;
  paid_at: string | null;
  tx_hash: string;
}

export interface Certificate {
  id: number;
  student_id: number;
  student_name: string;
  type: 'diploma' | 'certificate' | 'transcript' | 'achievement';
  title: string;
  description: string;
  issued_at: string;
  tx_hash: string;
  verified: boolean;
  ipfs_hash: string;
}

export interface Grade {
  id: number;
  student_id: number;
  course_id: number;
  course_name: string;
  semester: string;
  score: number;
  grade: string;
  teacher: string;
  recorded_at: string;
  tx_hash: string;
}

// ==================== Retail Types ====================

export interface Supplier {
  id: number;
  name: string;
  contact_person: string;
  phone: string;
  email: string;
  address: string;
  wallet_address: string;
  products_supplied: number;
  total_orders: number;
  rating: number;
  status: 'active' | 'inactive';
}

export interface PurchaseOrder {
  id: number;
  supplier_id: number;
  supplier_name: string;
  items: PurchaseOrderItem[];
  total: string;
  status: 'draft' | 'ordered' | 'received' | 'cancelled';
  ordered_at: string;
  received_at: string | null;
  tx_hash: string;
}

export interface PurchaseOrderItem {
  product_id: number;
  product_name: string;
  quantity: number;
  unit_price: string;
  total: string;
}

export interface Discount {
  id: number;
  code: string;
  name: string;
  type: 'percentage' | 'fixed' | 'buy_x_get_y';
  value: number;
  min_purchase: string;
  max_discount: string;
  start_date: string;
  end_date: string;
  usage_limit: number;
  used_count: number;
  status: 'active' | 'expired' | 'disabled';
}

export interface LoyaltyMember {
  id: number;
  customer_address: string;
  name: string;
  phone: string;
  tier: 'Bronze' | 'Silver' | 'Gold' | 'Platinum';
  points: number;
  total_spent: string;
  total_transactions: number;
  joined_at: string;
  last_transaction: string;
}

export interface StockMovement {
  id: number;
  product_id: number;
  product_name: string;
  type: 'in' | 'out' | 'adjustment' | 'return';
  quantity: number;
  before_stock: number;
  after_stock: number;
  reference: string;
  note: string;
  created_at: string;
}

export interface SalesReport {
  period: string;
  total_sales: string;
  total_transactions: number;
  total_items_sold: number;
  average_transaction: string;
  top_products: { name: string; quantity: number; revenue: string }[];
  sales_by_category: { category: string; total: string }[];
  sales_by_hour: { hour: number; total: string }[];
}

// ==================== Agen (Multi-Agen + Rekam Jejak) Types ====================

export interface Agen {
  id: number;
  address: string;
  name: string;
  parent_id: string;
  level: 'pusat' | 'cabang' | 'subagen';
  status: 'active' | 'inactive' | 'suspended';
  commission_rate: string;
  score: string;
  rating_avg: string;
  total_downline: number;
  total_sales: number;
  total_volume: string;
  complaints: Complaint[];
  performance: AgentPerformance[];
  creator: string;
  created_at: string;
  updated_at: string;
}

export interface Complaint {
  id: number;
  reporter: string;
  reason: string;
  status: 'open' | 'resolved';
  resolution: string;
  created_at: string;
  resolved_at: string;
}

export interface AgentPerformance {
  agent: string;
  period: string;
  sales: number;
  volume: string;
  rating_avg: string;
  score: string;
  created_at: string;
}

// ==================== Oleh-oleh Marketplace Types ====================

export interface OlehOlehProduct {
  id: number;
  name: string;
  description: string;
  price: string;
  image_url: string;
  stock: number;
  seller: string;
  category: string;
  status: 'active' | 'inactive';
  creator: string;
  created_at: string;
  updated_at: string;
}

export interface OlehOlehOrder {
  id: number;
  product_id: number;
  product_name: string;
  jamaah: string;
  quantity: number;
  total: string;
  status: 'pending' | 'paid' | 'shipped' | 'delivered' | 'cancelled';
  shipping_address: string;
  creator: string;
  created_at: string;
  updated_at: string;
}

// ==================== Asuransi Digital Types ====================

export interface Asuransi {
  id: number;
  jamaah: string;
  policy_type: string;
  premium: string;
  coverage: string;
  start_date: string;
  end_date: string;
  status: 'active' | 'expired' | 'cancelled' | 'claimed';
  document_hash: string;
  provider: string;
  creator: string;
  created_at: string;
  updated_at: string;
}

export interface AsuransiClaim {
  id: number;
  asuransi_id: number;
  jamaah: string;
  reason: string;
  amount: string;
  status: 'submitted' | 'approved' | 'rejected' | 'paid';
  evidence_hash: string;
  decision_by: string;
  decision_note: string;
  submitted_at: string;
  decided_at: string;
}

// ==================== DAO Voting Types ====================

export interface DaoProposal {
  id: number;
  title: string;
  description: string;
  options: string[];
  votes: Vote[];
  deadline: string;
  status: 'active' | 'closed' | 'passed' | 'rejected';
  result_option: string;
  creator: string;
  created_at: string;
  closed_at: string;
}

export interface Vote {
  voter: string;
  option: string;
  weight: string;
  voted_at: string;
}

// ==================== IBC Types ====================

export interface IBCChannel {
  channel_id: string;
  port_id: string;
  counterparty_channel: string;
  counterparty_port: string;
  connection_id: string;
  state: string;
}

export interface IBCCounterpartyChain {
  chain_id: string;
  client_id: string;
  connection_id: string;
  status: string;
}

// ==================== POS / Retail — Akuntansi, Multi-Satuan, Harga Level, Barang Gabungan ====================

export interface PosUnit {
  id: number;
  name: string;
  symbol: string;
  conversion_factor: number;
  is_base: boolean;
}

export interface PriceLevel {
  level: string;
  price: string;
  min_quantity: number;
}

export interface BundleComponent {
  product_id: number;
  quantity: number;
}

export interface PosProduct {
  id: number;
  name: string;
  description: string;
  price: string;
  cost_price: string;
  sku: string;
  category: string;
  image_url: string;
  stock: number;
  owner: string;
  active: boolean;
  created_at: string;
  updated_at: string;
  base_unit_id: number;
  price_levels: PriceLevel[];
  is_bundle: boolean;
  components: BundleComponent[];
  barcode: string;
  min_stock: number;
  branch_id: string;
}

export interface Account {
  id: number;
  code: string;
  name: string;
  type: 'asset' | 'liability' | 'equity' | 'revenue' | 'expense';
  description: string;
  created_at: string;
}

export interface JournalLine {
  account_id: number;
  debit: string;
  credit: string;
}

export interface JournalEntry {
  id: number;
  reference_type: string;
  reference_id: number;
  description: string;
  lines: JournalLine[];
  created_at: string;
  creator: string;
}

export interface AccountBalance {
  account_id: number;
  code: string;
  name: string;
  type: string;
  debit: string;
  credit: string;
  balance: string;
}

export interface TrialBalance {
  accounts: AccountBalance[];
  total_debit: string;
  total_credit: string;
}

export interface IncomeStatement {
  revenues: AccountBalance[];
  expenses: AccountBalance[];
  total_revenue: string;
  total_expense: string;
  net_income: string;
}

export interface BalanceSheet {
  assets: AccountBalance[];
  liabilities: AccountBalance[];
  equities: AccountBalance[];
  total_assets: string;
  total_liabilities: string;
  total_equity: string;
}

export interface LedgerLine {
  journal_entry_id: number;
  reference_type: string;
  reference_id: number;
  description: string;
  debit: string;
  credit: string;
  balance: string;
  created_at: string;
}

export interface Ledger {
  account: Account;
  lines: LedgerLine[];
  ending_balance: string;
}

export interface PriceLevelReportItem {
  product_id: number;
  product_name: string;
  sku: string;
  base_price: string;
  base_unit: string;
  price_levels: PriceLevel[];
}

// ==================== Haji / Umroh Types ====================

export interface Jamaah {
  id: number;
  name: string;
  phone: string;
  email: string;
  address: string;
  passport_number: string;
  photo_hash: string;
  status: 'active' | 'inactive' | 'blocked';
  did: string;
  creator: string;
  documents: DocumentHash[];
  vaccinations: VaccinationRecord[];
  created_at: string;
  updated_at: string;
}

export interface DocumentHash {
  doc_type: string;
  hash: string;
  storage_ref: string;
  uploaded_at: string;
}

export interface VaccinationRecord {
  vaccine_name: string;
  date: string;
  issuer: string;
  batch: string;
}

export interface Paket {
  id: number;
  name: string;
  price: string;
  schedule: string;
  quota: number;
  booked: number;
  hotel: string;
  airline: string;
  muthawif: string;
  status: 'open' | 'full' | 'closed' | 'departed' | 'completed';
  departure_date: string;
  return_date: string;
  category: 'haji' | 'umroh' | 'haji_plus';
  description: string;
  image_url: string;
  reviews: Review[];
  creator: string;
  created_at: string;
  updated_at: string;
}

export interface Review {
  reviewer: string;
  rating: number;
  comment: string;
  created_at: string;
}

export interface Pembayaran {
  id: number;
  jamaah: string;
  paket_id: number;
  total: string;
  down_payment: string;
  paid: string;
  remaining: string;
  status: 'pending' | 'dp_paid' | 'active' | 'completed' | 'refunded' | 'cancelled';
  payment_method: string;
  installments: Installment[];
  escrow_stages: EscrowStage[];
  creator: string;
  created_at: string;
  updated_at: string;
}

export interface Installment {
  id: number;
  amount: string;
  due_date: string;
  paid_at: string;
  paid: boolean;
  late_fee: string;
}

export interface EscrowStage {
  name: string;
  amount: string;
  released: boolean;
  released_at: string;
}

export interface Visa {
  id: number;
  jamaah: string;
  paket_id: number;
  status: 'processing' | 'issued' | 'rejected' | 'expired';
  visa_number: string;
  issue_date: string;
  expiry_date: string;
  document_hash: string;
  notes: string;
  creator: string;
  created_at: string;
  updated_at: string;
}

export interface Hotel {
  id: number;
  name: string;
  city: string;
  address: string;
  star_rating: string;
  price_per_night: string;
  room_type: string;
  available_rooms: number;
  distance_haram: string;
  status: 'active' | 'inactive';
  creator: string;
  created_at: string;
  updated_at: string;
}

export interface Ticket {
  id: number;
  jamaah: string;
  paket_id: number;
  airline: string;
  flight_number: string;
  seat: string;
  schedule: string;
  qr_code: string;
  nft_id: string;
  status: 'issued' | 'checked_in' | 'boarded' | 'used' | 'void';
  document_hash: string;
  creator: string;
  created_at: string;
}

export interface Referral {
  id: number;
  agent: string;
  referred_jamaah: string;
  paket_id: number;
  commission_rate: string;
  commission: string;
  status: 'pending' | 'paid' | 'cancelled';
  paid_at: string;
  creator: string;
  created_at: string;
}

export interface RewardBalance {
  jamaah: string;
  balance: string;
  earned: string;
  redeemed: string;
}

export interface Reward {
  id: number;
  jamaah: string;
  points: string;
  reward_type: 'cashback' | 'discount' | 'referral_bonus' | 'redeem';
  reason: string;
  status: 'awarded' | 'redeemed' | 'expired';
  creator: string;
  created_at: string;
}

export interface AuditLog {
  id: number;
  module: string;
  action: string;
  actor: string;
  target_id: string;
  data_hash: string;
  metadata: string;
  created_at: string;
}

export interface Keberangkatan {
  id: number;
  jamaah: string;
  paket_id: number;
  pembayaran_id: number;
  stage: number;
  status_label: string;
  departure_date: string;
  return_date: string;
  manasik_date: string;
  hotel_confirm: string;
  airline_confirm: string;
  baggage: BaggageItem[];
  creator: string;
  created_at: string;
  updated_at: string;
}

export interface BaggageItem {
  id: number;
  tag: string;
  weight: string;
  status: 'checked_in' | 'in_transit' | 'arrived' | 'delivered';
  updated_at: string;
}

export const JOURNEY_STAGES = [
  'Daftar',
  'DP Dibayar',
  'Visa Diproses',
  'Visa Terbit',
  'Tiket Terbit',
  'Hotel Confirm',
  'Manasik',
  'Berangkat',
  'Pulang',
]

// Helper functions
export function formatCoin(amount: string, denom: string): string {
  const denomination = denom === 'uidr' ? 'IDR' : denom.toUpperCase();
  const value = Number(amount) / 1000000;
  return `${value.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 6 })} ${denomination}`;
}

export function formatRupiah(amount: string): string {
  const value = Number(amount || 0);
  if (isNaN(value)) return `Rp ${amount}`;
  return `Rp ${value.toLocaleString('id-ID')}`;
}

export function formatNumber(amount: string): string {
  const value = Number(amount || 0);
  if (isNaN(value)) return amount;
  return value.toLocaleString('id-ID');
}

export function formatAddress(address: string): string {
  if (!address) return '';
  return `${address.slice(0, 8)}...${address.slice(-6)}`;
}

export function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}
