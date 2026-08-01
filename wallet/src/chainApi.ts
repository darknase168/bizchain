// BizChain Chain API
// Queries the live chain through the Cosmos SDK REST gateway (grpc-gateway),
// which proxies directly to each module's gRPC Query service.
// Endpoint mapping:
//   agen     GET /bizchain/agen/agen          -> Query/AgenAll
//   oleholeh GET /bizchain/oleholeh/product   -> Query/ProductAll
//            GET /bizchain/oleholeh/order     -> Query/OrderAll
//   asuransi GET /bizchain/asuransi/asuransi  -> Query/AsuransiAll
//            GET /bizchain/asuransi/claim     -> Query/ClaimAll
//   dao      GET /bizchain/dao/proposal       -> Query/ProposalAll
//   IBC      GET /ibc/core/channel/v1/channels            -> ibc.channel Query/Channels
//            GET /ibc/core/connection/v1/connections      -> ibc.connection Query/Connections
//            GET /ibc/core/client/v1/client_states        -> ibc.client Query/ClientStates

import { CHAIN_CONFIG } from './types'
import type {
  Agen, Complaint, AgentPerformance,
  OlehOlehProduct, OlehOlehOrder,
  Asuransi, AsuransiClaim,
  DaoProposal, Vote,
  IBCChannel, IBCCounterpartyChain,
  Jamaah, DocumentHash, VaccinationRecord,
  Paket, Review,
  Pembayaran, Installment, EscrowStage,
  Visa, Hotel,
  Ticket as TicketType,
  Referral, Reward, RewardBalance,
  AuditLog, Keberangkatan, BaggageItem,
  PosUnit, PriceLevel, BundleComponent, PosProduct,
  Account, JournalLine, JournalEntry, AccountBalance,
  TrialBalance, IncomeStatement, BalanceSheet, Ledger, LedgerLine,
  PriceLevelReportItem,
  Transaction, TransactionItem,
} from './types'

const REST_URL = CHAIN_CONFIG.restUrl

class ChainApiError extends Error {
  status: number
  constructor(message: string, status: number) {
    super(message)
    this.status = status
  }
}

async function get<T>(path: string): Promise<T> {
  let res: Response
  try {
    res = await fetch(`${REST_URL}${path}`)
  } catch {
    throw new ChainApiError(
      `Tidak dapat terhubung ke node (${REST_URL}). Pastikan \`bizchaind start\` berjalan.`,
      0,
    )
  }
  if (!res.ok) {
    throw new ChainApiError(`REST ${path} → HTTP ${res.status}`, res.status)
  }
  return res.json() as Promise<T>
}

// proto3 JSON encodes uint64 as strings; normalize to number to match TS types.
const num = (v: string | number | null | undefined): number => {
  if (v === null || v === undefined || v === '') return 0
  const n = Number(v)
  return Number.isFinite(n) ? n : 0
}
const str = (v: string | null | undefined): string => v ?? ''
// Cast a raw status string to the TS literal union, falling back to a default.
const status = <T extends string>(v: string | null | undefined, fallback: T): T => {
  const s = str(v)
  return (s || fallback) as T
}

// ------------------- New multi-cabang helpers -------------------

interface RawPosProductWithBranch {
  id?: string
  name?: string
  description?: string
  price?: string
  cost_price?: string
  sku?: string
  category?: string
  image_url?: string
  stock?: string
  owner?: string
  active?: boolean
  created_at?: string
  updated_at?: string
  base_unit_id?: string
  price_levels?: RawPriceLevel[]
  is_bundle?: boolean
  components?: RawBundleComponent[]
  barcode?: string
  min_stock?: string
  branch_id?: string
}

interface RawTransactionWithBranch {
  id?: string
  seller?: string
  customer_address?: string
  items?: RawTransactionItem[]
  total?: string
  discount?: string
  tax?: string
  grand_total?: string
  payment_method?: string
  status?: string
  note?: string
  created_at?: string
  updated_at?: string
  branch_id?: string
}

interface RawTransactionItem {
  product_id?: string
  quantity?: string
  unit_id?: string
  price?: string
  subtotal?: string
  cost?: string
}

// ==================== Normalizers (proto JSON -> TS types) ====================

interface RawComplaint {
  id?: string
  reporter?: string
  reason?: string
  status?: string
  resolution?: string
  created_at?: string
  resolved_at?: string
}

interface RawAgentPerformance {
  agent?: string
  period?: string
  sales?: string
  volume?: string
  rating_avg?: string
  score?: string
  created_at?: string
}

interface RawAgen {
  id?: string
  address?: string
  name?: string
  parent_id?: string
  level?: string
  status?: string
  commission_rate?: string
  score?: string
  rating_avg?: string
  total_downline?: string
  total_sales?: string
  total_volume?: string
  complaints?: RawComplaint[]
  performance?: RawAgentPerformance[]
  creator?: string
  created_at?: string
  updated_at?: string
}

const mapComplaint = (c: RawComplaint): Complaint => ({
  id: num(c.id),
  reporter: str(c.reporter),
  reason: str(c.reason),
  status: status(c.status, 'open'),
  resolution: str(c.resolution),
  created_at: str(c.created_at),
  resolved_at: str(c.resolved_at),
})

const mapPerformance = (p: RawAgentPerformance): AgentPerformance => ({
  agent: str(p.agent),
  period: str(p.period),
  sales: num(p.sales),
  volume: str(p.volume),
  rating_avg: str(p.rating_avg),
  score: str(p.score),
  created_at: str(p.created_at),
})

const mapAgen = (a: RawAgen): Agen => ({
  id: num(a.id),
  address: str(a.address),
  name: str(a.name),
  parent_id: str(a.parent_id) || '0',
  level: status(a.level, 'subagen'),
  status: status(a.status, 'active'),
  commission_rate: str(a.commission_rate),
  score: str(a.score) || '0',
  rating_avg: str(a.rating_avg) || '0',
  total_downline: num(a.total_downline),
  total_sales: num(a.total_sales),
  total_volume: str(a.total_volume) || '0',
  complaints: (a.complaints ?? []).map(mapComplaint),
  performance: (a.performance ?? []).map(mapPerformance),
  creator: str(a.creator),
  created_at: str(a.created_at),
  updated_at: str(a.updated_at),
})

interface RawProduct {
  id?: string
  name?: string
  description?: string
  price?: string
  image_url?: string
  stock?: string
  seller?: string
  category?: string
  status?: string
  creator?: string
  created_at?: string
  updated_at?: string
}

const mapProduct = (p: RawProduct): OlehOlehProduct => ({
  id: num(p.id),
  name: str(p.name),
  description: str(p.description),
  price: str(p.price) || '0',
  image_url: str(p.image_url),
  stock: num(p.stock),
  seller: str(p.seller),
  category: str(p.category),
  status: status(p.status, 'active'),
  creator: str(p.creator),
  created_at: str(p.created_at),
  updated_at: str(p.updated_at),
})

interface RawOrder {
  id?: string
  product_id?: string
  product_name?: string
  jamaah?: string
  quantity?: string
  total?: string
  status?: string
  shipping_address?: string
  creator?: string
  created_at?: string
  updated_at?: string
}

const mapOrder = (o: RawOrder): OlehOlehOrder => ({
  id: num(o.id),
  product_id: num(o.product_id),
  product_name: str(o.product_name),
  jamaah: str(o.jamaah),
  quantity: num(o.quantity),
  total: str(o.total) || '0',
  status: status(o.status, 'pending'),
  shipping_address: str(o.shipping_address),
  creator: str(o.creator),
  created_at: str(o.created_at),
  updated_at: str(o.updated_at),
})

interface RawAsuransi {
  id?: string
  jamaah?: string
  policy_type?: string
  premium?: string
  coverage?: string
  start_date?: string
  end_date?: string
  status?: string
  document_hash?: string
  provider?: string
  creator?: string
  created_at?: string
  updated_at?: string
}

const mapAsuransi = (a: RawAsuransi): Asuransi => ({
  id: num(a.id),
  jamaah: str(a.jamaah),
  policy_type: str(a.policy_type),
  premium: str(a.premium) || '0',
  coverage: str(a.coverage) || '0',
  start_date: str(a.start_date),
  end_date: str(a.end_date),
  status: status(a.status, 'active'),
  document_hash: str(a.document_hash),
  provider: str(a.provider),
  creator: str(a.creator),
  created_at: str(a.created_at),
  updated_at: str(a.updated_at),
})

interface RawClaim {
  id?: string
  asuransi_id?: string
  jamaah?: string
  reason?: string
  amount?: string
  status?: string
  evidence_hash?: string
  decision_by?: string
  decision_note?: string
  submitted_at?: string
  decided_at?: string
}

const mapClaim = (c: RawClaim): AsuransiClaim => ({
  id: num(c.id),
  asuransi_id: num(c.asuransi_id),
  jamaah: str(c.jamaah),
  reason: str(c.reason),
  amount: str(c.amount) || '0',
  status: status(c.status, 'submitted'),
  evidence_hash: str(c.evidence_hash),
  decision_by: str(c.decision_by),
  decision_note: str(c.decision_note),
  submitted_at: str(c.submitted_at),
  decided_at: str(c.decided_at),
})

interface RawVote {
  voter?: string
  option?: string
  weight?: string
  voted_at?: string
}

const mapVote = (v: RawVote): Vote => ({
  voter: str(v.voter),
  option: str(v.option),
  weight: str(v.weight) || '1',
  voted_at: str(v.voted_at),
})

interface RawProposal {
  id?: string
  title?: string
  description?: string
  options?: string[]
  votes?: RawVote[]
  deadline?: string
  status?: string
  result_option?: string
  creator?: string
  created_at?: string
  closed_at?: string
}

const mapProposal = (p: RawProposal): DaoProposal => ({
  id: num(p.id),
  title: str(p.title),
  description: str(p.description),
  options: p.options ?? [],
  votes: (p.votes ?? []).map(mapVote),
  deadline: str(p.deadline),
  status: status(p.status, 'active'),
  result_option: str(p.result_option),
  creator: str(p.creator),
  created_at: str(p.created_at),
  closed_at: str(p.closed_at),
})

// ---- IBC raw shapes ----

interface RawIBCChannel {
  channel_id?: string
  channel?: {
    state?: string
    ordering?: string
    counterparty?: { port_id?: string; channel_id?: string }
    connection_hops?: string[]
    version?: string
    port_id?: string
  }
}

interface RawIBCConnection {
  id?: string
  client_id?: string
  state?: string
}

interface RawIBCClientState {
  client_id?: string
  client_state?: {
    chain_id?: string
    '@type'?: string
  }
}

const mapChannel = (c: RawIBCChannel): IBCChannel => ({
  channel_id: str(c.channel_id),
  port_id: str(c.channel?.port_id),
  counterparty_channel: str(c.channel?.counterparty?.channel_id),
  counterparty_port: str(c.channel?.counterparty?.port_id),
  connection_id: c.channel?.connection_hops?.[0] ?? '',
  state: (c.channel?.state ?? 'STATE_UNINITIALIZED').replace('STATE_', ''),
})

// ==================== Public query API ====================

/** Query/AgenAll — all registered agents (pusat/cabang/subagen) */
export async function fetchAgents(): Promise<Agen[]> {
  const data = await get<{ agen?: RawAgen[] }>('/bizchain/agen/agen')
  return (data.agen ?? []).map(mapAgen)
}

/** Query/ProductAll — all oleh-oleh marketplace products */
export async function fetchOlehOlehProducts(): Promise<OlehOlehProduct[]> {
  const data = await get<{ product?: RawProduct[] }>('/bizchain/oleholeh/product')
  return (data.product ?? []).map(mapProduct)
}

/** Query/OrderAll — all oleh-oleh pre-orders */
export async function fetchOlehOlehOrders(): Promise<OlehOlehOrder[]> {
  const data = await get<{ order?: RawOrder[] }>('/bizchain/oleholeh/order')
  return (data.order ?? []).map(mapOrder)
}

/** Query/AsuransiAll — all insurance policies */
export async function fetchAsuransi(): Promise<Asuransi[]> {
  const data = await get<{ asuransi?: RawAsuransi[] }>('/bizchain/asuransi/asuransi')
  return (data.asuransi ?? []).map(mapAsuransi)
}

/** Query/ClaimAll — all insurance claims */
export async function fetchAsuransiClaims(): Promise<AsuransiClaim[]> {
  const data = await get<{ claim?: RawClaim[] }>('/bizchain/asuransi/claim')
  return (data.claim ?? []).map(mapClaim)
}

/** Query/ProposalAll — all DAO proposals */
export async function fetchDaoProposals(): Promise<DaoProposal[]> {
  const data = await get<{ proposal?: RawProposal[] }>('/bizchain/dao/proposal')
  return (data.proposal ?? []).map(mapProposal)
}

/** ibc.core.channel Query/Channels — all IBC channels */
export async function fetchIBCChannels(): Promise<IBCChannel[]> {
  const data = await get<{ channels?: RawIBCChannel[] }>('/ibc/core/channel/v1/channels')
  return (data.channels ?? []).map(mapChannel)
}

export interface IBCNetworkStatus {
  channels: IBCChannel[]
  counterparties: IBCCounterpartyChain[]
}

/**
 * IBC counterparty status, built by joining ibc.core.connection Query/Connections
 * with ibc.core.client Query/ClientStates (client -> connection, chain-id from
 * the client state's chain_id field).
 */
export async function fetchIBCNetworkStatus(): Promise<IBCNetworkStatus> {
  const [channels, conns, clients] = await Promise.all([
    fetchIBCChannels(),
    get<{ connections?: RawIBCConnection[] }>('/ibc/core/connection/v1/connections'),
    get<{ client_states?: RawIBCClientState[] }>('/ibc/core/client/v1/client_states'),
  ])

  const counterparties: IBCCounterpartyChain[] = (clients.client_states ?? []).map((cs) => {
    const conn = (conns.connections ?? []).find((c) => c.client_id === cs.client_id)
    const state = str(conn?.state).replace('STATE_', '')
    return {
      chain_id: str(cs.client_state?.chain_id),
      client_id: str(cs.client_id),
      connection_id: str(conn?.id),
      status: state === 'OPEN' ? 'Open' : state === 'UNINITIALIZED' ? 'Uninitialized' : state || 'Uninitialized',
    }
  })

  return { channels, counterparties }
}

// ==================== Haji & Umroh raw shapes ====================

interface RawDocumentHash {
  doc_type?: string
  hash?: string
  storage_ref?: string
  uploaded_at?: string
}

interface RawVaccination {
  vaccine_name?: string
  date?: string
  issuer?: string
  batch?: string
}

interface RawJamaah {
  id?: string
  name?: string
  phone?: string
  email?: string
  address?: string
  passport_number?: string
  photo_hash?: string
  status?: string
  did?: string
  creator?: string
  documents?: RawDocumentHash[]
  vaccinations?: RawVaccination[]
  created_at?: string
  updated_at?: string
}

interface RawReview {
  reviewer?: string
  rating?: number
  comment?: string
  created_at?: string
}

interface RawPaket {
  id?: string
  name?: string
  price?: string
  schedule?: string
  quota?: string
  booked?: string
  hotel?: string
  airline?: string
  muthawif?: string
  status?: string
  departure_date?: string
  return_date?: string
  category?: string
  description?: string
  image_url?: string
  reviews?: RawReview[]
  creator?: string
  created_at?: string
  updated_at?: string
}

interface RawInstallment {
  id?: string
  amount?: string
  due_date?: string
  paid_at?: string
  paid?: boolean
  late_fee?: string
}

interface RawEscrowStage {
  name?: string
  amount?: string
  released?: boolean
  released_at?: string
}

interface RawPembayaran {
  id?: string
  jamaah?: string
  paket_id?: string
  total?: string
  down_payment?: string
  paid?: string
  remaining?: string
  status?: string
  payment_method?: string
  installments?: RawInstallment[]
  escrow_stages?: RawEscrowStage[]
  creator?: string
  created_at?: string
  updated_at?: string
}

interface RawVisa {
  id?: string
  jamaah?: string
  paket_id?: string
  status?: string
  visa_number?: string
  issue_date?: string
  expiry_date?: string
  document_hash?: string
  notes?: string
  creator?: string
  created_at?: string
  updated_at?: string
}

interface RawHotel {
  id?: string
  name?: string
  city?: string
  address?: string
  star_rating?: string
  price_per_night?: string
  room_type?: string
  available_rooms?: string
  distance_haram?: string
  status?: string
  creator?: string
  created_at?: string
  updated_at?: string
}

interface RawTicket {
  id?: string
  jamaah?: string
  paket_id?: string
  airline?: string
  flight_number?: string
  seat?: string
  schedule?: string
  qr_code?: string
  nft_id?: string
  status?: string
  document_hash?: string
  creator?: string
  created_at?: string
}

interface RawReferral {
  id?: string
  agent?: string
  referred_jamaah?: string
  paket_id?: string
  commission_rate?: string
  commission?: string
  status?: string
  paid_at?: string
  creator?: string
  created_at?: string
}

interface RawReward {
  id?: string
  jamaah?: string
  points?: string
  reward_type?: string
  reason?: string
  status?: string
  creator?: string
  created_at?: string
}

interface RawAuditLog {
  id?: string
  module?: string
  action?: string
  actor?: string
  target_id?: string
  data_hash?: string
  metadata?: string
  created_at?: string
}

interface RawBaggageItem {
  id?: string
  tag?: string
  weight?: string
  status?: string
  updated_at?: string
}

interface RawKeberangkatan {
  id?: string
  jamaah?: string
  paket_id?: string
  pembayaran_id?: string
  stage?: string
  status_label?: string
  departure_date?: string
  return_date?: string
  manasik_date?: string
  hotel_confirm?: string
  airline_confirm?: string
  baggage?: RawBaggageItem[]
  creator?: string
  created_at?: string
  updated_at?: string
}

// ---- Haji & Umroh mappers ----

const mapDocumentHash = (d: RawDocumentHash): DocumentHash => ({
  doc_type: str(d.doc_type),
  hash: str(d.hash),
  storage_ref: str(d.storage_ref),
  uploaded_at: str(d.uploaded_at),
})

const mapVaccination = (v: RawVaccination): VaccinationRecord => ({
  vaccine_name: str(v.vaccine_name),
  date: str(v.date),
  issuer: str(v.issuer),
  batch: str(v.batch),
})

const mapJamaah = (j: RawJamaah): Jamaah => ({
  id: num(j.id),
  name: str(j.name),
  phone: str(j.phone),
  email: str(j.email),
  address: str(j.address),
  passport_number: str(j.passport_number),
  photo_hash: str(j.photo_hash),
  status: status(j.status, 'active'),
  did: str(j.did),
  creator: str(j.creator),
  documents: (j.documents ?? []).map(mapDocumentHash),
  vaccinations: (j.vaccinations ?? []).map(mapVaccination),
  created_at: str(j.created_at),
  updated_at: str(j.updated_at),
})

const mapReview = (r: RawReview): Review => ({
  reviewer: str(r.reviewer),
  rating: num(r.rating),
  comment: str(r.comment),
  created_at: str(r.created_at),
})

const mapPaket = (p: RawPaket): Paket => ({
  id: num(p.id),
  name: str(p.name),
  price: str(p.price) || '0',
  schedule: str(p.schedule),
  quota: num(p.quota),
  booked: num(p.booked),
  hotel: str(p.hotel),
  airline: str(p.airline),
  muthawif: str(p.muthawif),
  status: status(p.status, 'open'),
  departure_date: str(p.departure_date),
  return_date: str(p.return_date),
  category: status(p.category, 'umroh'),
  description: str(p.description),
  image_url: str(p.image_url),
  reviews: (p.reviews ?? []).map(mapReview),
  creator: str(p.creator),
  created_at: str(p.created_at),
  updated_at: str(p.updated_at),
})

const mapInstallment = (i: RawInstallment): Installment => ({
  id: num(i.id),
  amount: str(i.amount) || '0',
  due_date: str(i.due_date),
  paid_at: str(i.paid_at),
  paid: !!i.paid,
  late_fee: str(i.late_fee) || '0',
})

const mapEscrowStage = (e: RawEscrowStage): EscrowStage => ({
  name: str(e.name),
  amount: str(e.amount) || '0',
  released: !!e.released,
  released_at: str(e.released_at),
})

const mapPembayaran = (p: RawPembayaran): Pembayaran => ({
  id: num(p.id),
  jamaah: str(p.jamaah),
  paket_id: num(p.paket_id),
  total: str(p.total) || '0',
  down_payment: str(p.down_payment) || '0',
  paid: str(p.paid) || '0',
  remaining: str(p.remaining) || '0',
  status: status(p.status, 'pending'),
  payment_method: str(p.payment_method),
  installments: (p.installments ?? []).map(mapInstallment),
  escrow_stages: (p.escrow_stages ?? []).map(mapEscrowStage),
  creator: str(p.creator),
  created_at: str(p.created_at),
  updated_at: str(p.updated_at),
})

const mapVisa = (v: RawVisa): Visa => ({
  id: num(v.id),
  jamaah: str(v.jamaah),
  paket_id: num(v.paket_id),
  status: status(v.status, 'processing'),
  visa_number: str(v.visa_number),
  issue_date: str(v.issue_date),
  expiry_date: str(v.expiry_date),
  document_hash: str(v.document_hash),
  notes: str(v.notes),
  creator: str(v.creator),
  created_at: str(v.created_at),
  updated_at: str(v.updated_at),
})

const mapHotel = (h: RawHotel): Hotel => ({
  id: num(h.id),
  name: str(h.name),
  city: str(h.city),
  address: str(h.address),
  star_rating: str(h.star_rating),
  price_per_night: str(h.price_per_night) || '0',
  room_type: str(h.room_type),
  available_rooms: num(h.available_rooms),
  distance_haram: str(h.distance_haram),
  status: status(h.status, 'active'),
  creator: str(h.creator),
  created_at: str(h.created_at),
  updated_at: str(h.updated_at),
})

const mapTicket = (t: RawTicket): TicketType => ({
  id: num(t.id),
  jamaah: str(t.jamaah),
  paket_id: num(t.paket_id),
  airline: str(t.airline),
  flight_number: str(t.flight_number),
  seat: str(t.seat),
  schedule: str(t.schedule),
  qr_code: str(t.qr_code),
  nft_id: str(t.nft_id),
  status: status(t.status, 'issued'),
  document_hash: str(t.document_hash),
  creator: str(t.creator),
  created_at: str(t.created_at),
})

const mapReferral = (r: RawReferral): Referral => ({
  id: num(r.id),
  agent: str(r.agent),
  referred_jamaah: str(r.referred_jamaah),
  paket_id: num(r.paket_id),
  commission_rate: str(r.commission_rate) || '0',
  commission: str(r.commission) || '0',
  status: status(r.status, 'pending'),
  paid_at: str(r.paid_at),
  creator: str(r.creator),
  created_at: str(r.created_at),
})

const mapReward = (r: RawReward): Reward => ({
  id: num(r.id),
  jamaah: str(r.jamaah),
  points: str(r.points) || '0',
  reward_type: status(r.reward_type, 'cashback'),
  reason: str(r.reason),
  status: status(r.status, 'awarded'),
  creator: str(r.creator),
  created_at: str(r.created_at),
})

const mapAuditLog = (l: RawAuditLog): AuditLog => ({
  id: num(l.id),
  module: str(l.module),
  action: str(l.action),
  actor: str(l.actor),
  target_id: str(l.target_id),
  data_hash: str(l.data_hash),
  metadata: str(l.metadata),
  created_at: str(l.created_at),
})

const mapBaggageItem = (b: RawBaggageItem): BaggageItem => ({
  id: num(b.id),
  tag: str(b.tag),
  weight: str(b.weight),
  status: status(b.status, 'checked_in'),
  updated_at: str(b.updated_at),
})

const mapKeberangkatan = (k: RawKeberangkatan): Keberangkatan => ({
  id: num(k.id),
  jamaah: str(k.jamaah),
  paket_id: num(k.paket_id),
  pembayaran_id: num(k.pembayaran_id),
  stage: num(k.stage),
  status_label: str(k.status_label),
  departure_date: str(k.departure_date),
  return_date: str(k.return_date),
  manasik_date: str(k.manasik_date),
  hotel_confirm: str(k.hotel_confirm),
  airline_confirm: str(k.airline_confirm),
  baggage: (k.baggage ?? []).map(mapBaggageItem),
  creator: str(k.creator),
  created_at: str(k.created_at),
  updated_at: str(k.updated_at),
})

// ==================== Haji & Umroh query API ====================

/** Query/JamaahAll — all registered pilgrims */
export async function fetchJamaah(): Promise<Jamaah[]> {
  const data = await get<{ jamaah?: RawJamaah[] }>('/bizchain/jamaah/jamaah')
  return (data.jamaah ?? []).map(mapJamaah)
}

/** Query/PaketAll — all hajj/umroh packages */
export async function fetchPaket(): Promise<Paket[]> {
  const data = await get<{ paket?: RawPaket[] }>('/bizchain/paket/paket')
  return (data.paket ?? []).map(mapPaket)
}

/** Query/PembayaranAll — all escrow payments */
export async function fetchPembayaran(): Promise<Pembayaran[]> {
  const data = await get<{ pembayaran?: RawPembayaran[] }>('/bizchain/pembayaran/pembayaran')
  return (data.pembayaran ?? []).map(mapPembayaran)
}

/** Query/VisaAll — all visa applications */
export async function fetchVisa(): Promise<Visa[]> {
  const data = await get<{ visa?: RawVisa[] }>('/bizchain/visa/visa')
  return (data.visa ?? []).map(mapVisa)
}

/** Query/HotelAll — all hotels */
export async function fetchHotel(): Promise<Hotel[]> {
  const data = await get<{ hotel?: RawHotel[] }>('/bizchain/hotel/hotel')
  return (data.hotel ?? []).map(mapHotel)
}

/** Query/TicketAll — all NFT flight tickets */
export async function fetchTickets(): Promise<TicketType[]> {
  const data = await get<{ ticket?: RawTicket[] }>('/bizchain/ticket/ticket')
  return (data.ticket ?? []).map(mapTicket)
}

/** Query/ReferralAll — all referrals */
export async function fetchReferrals(): Promise<Referral[]> {
  const data = await get<{ referral?: RawReferral[] }>('/bizchain/referral/referral')
  return (data.referral ?? []).map(mapReferral)
}

/** Query/RewardAll — all loyalty rewards */
export async function fetchRewards(): Promise<Reward[]> {
  const data = await get<{ reward?: RawReward[] }>('/bizchain/reward/reward')
  return (data.reward ?? []).map(mapReward)
}

/** Query/AuditLogAll — all audit logs */
export async function fetchAuditLogs(): Promise<AuditLog[]> {
  const data = await get<{ audit_log?: RawAuditLog[] }>('/bizchain/audit/audit_log')
  return (data.audit_log ?? []).map(mapAuditLog)
}

/** Query/KeberangkatanAll — all departure journeys */
export async function fetchKeberangkatan(): Promise<Keberangkatan[]> {
  const data = await get<{ keberangkatan?: RawKeberangkatan[] }>('/bizchain/keberangkatan/keberangkatan')
  return (data.keberangkatan ?? []).map(mapKeberangkatan)
}

/** Aggregate per-jamaah reward balances from the reward ledger (earned/redeemed). */
export function deriveRewardBalances(rewards: Reward[]): RewardBalance[] {
  const map = new Map<string, { earned: number; redeemed: number }>()
  for (const r of rewards) {
    if (r.status === 'expired') continue // expired points no longer count
    const cur = map.get(r.jamaah) ?? { earned: 0, redeemed: 0 }
    const pts = parseInt(r.points || '0', 10) || 0
    if (r.status === 'redeemed') cur.redeemed += pts
    else cur.earned += pts
    map.set(r.jamaah, cur)
  }
  return Array.from(map.entries()).map(([jamaah, v]) => ({
    jamaah,
    earned: String(v.earned),
    redeemed: String(v.redeemed),
    balance: String(v.earned - v.redeemed),
  }))
}

// ==================== POS / Retail raw shapes (akuntansi, satuan, harga level, bundle) ====================

interface RawPosUnit {
  id?: string
  name?: string
  symbol?: string
  conversion_factor?: string
  is_base?: boolean
}

interface RawPriceLevel {
  level?: string
  price?: string
  min_quantity?: string
}

interface RawBundleComponent {
  product_id?: string
  quantity?: string
}

interface RawAccount {
  id?: string
  code?: string
  name?: string
  type?: string
  description?: string
  created_at?: string
}

interface RawJournalLine {
  account_id?: string
  debit?: string
  credit?: string
}

interface RawJournalEntry {
  id?: string
  reference_type?: string
  reference_id?: string
  description?: string
  lines?: RawJournalLine[]
  created_at?: string
  creator?: string
}

interface RawAccountBalance {
  account_id?: string
  code?: string
  name?: string
  type?: string
  debit?: string
  credit?: string
  balance?: string
}

interface RawTrialBalance {
  accounts?: RawAccountBalance[]
  total_debit?: string
  total_credit?: string
}

interface RawIncomeStatement {
  revenues?: RawAccountBalance[]
  expenses?: RawAccountBalance[]
  total_revenue?: string
  total_expense?: string
  net_income?: string
}

interface RawBalanceSheet {
  assets?: RawAccountBalance[]
  liabilities?: RawAccountBalance[]
  equities?: RawAccountBalance[]
  total_assets?: string
  total_liabilities?: string
  total_equity?: string
}

interface RawLedgerLine {
  journal_entry_id?: string
  reference_type?: string
  reference_id?: string
  description?: string
  debit?: string
  credit?: string
  balance?: string
  created_at?: string
}

interface RawLedger {
  account?: RawAccount
  lines?: RawLedgerLine[]
  ending_balance?: string
}

interface RawPriceLevelReportItem {
  product_id?: string
  product_name?: string
  sku?: string
  base_price?: string
  base_unit?: string
  price_levels?: RawPriceLevel[]
}

// ---- POS / Retail mappers ----

const mapPriceLevel = (p: RawPriceLevel): PriceLevel => ({
  level: str(p.level),
  price: str(p.price) || '0',
  min_quantity: num(p.min_quantity),
})

const mapBundleComponent = (c: RawBundleComponent): BundleComponent => ({
  product_id: num(c.product_id),
  quantity: num(c.quantity),
})

const mapPosUnit = (u: RawPosUnit): PosUnit => ({
  id: num(u.id),
  name: str(u.name),
  symbol: str(u.symbol),
  conversion_factor: num(u.conversion_factor),
  is_base: !!u.is_base,
})

const mapPosProduct = (p: RawPosProductWithBranch): PosProduct => ({
  id: num(p.id),
  name: str(p.name),
  description: str(p.description),
  price: str(p.price) || '0',
  cost_price: str(p.cost_price) || '0',
  sku: str(p.sku),
  category: str(p.category),
  image_url: str(p.image_url),
  stock: num(p.stock),
  owner: str(p.owner),
  active: p.active !== false,
  created_at: str(p.created_at),
  updated_at: str(p.updated_at),
  base_unit_id: num(p.base_unit_id),
  price_levels: (p.price_levels ?? []).map(mapPriceLevel),
  is_bundle: !!p.is_bundle,
  components: (p.components ?? []).map(mapBundleComponent),
  barcode: str(p.barcode),
  min_stock: num(p.min_stock),
  branch_id: str(p.branch_id),
})

const mapTransactionItem = (i: RawTransactionItem): TransactionItem => ({
  product_id: num(i.product_id),
  quantity: num(i.quantity),
  price: str(i.price) || '0',
})

const mapTransaction = (t: RawTransactionWithBranch): Transaction => ({
  id: num(t.id),
  seller: str(t.seller),
  customer_address: str(t.customer_address),
  items: (t.items ?? []).map(mapTransactionItem),
  total: str(t.total) || '0',
  status: str(t.status),
  note: str(t.note),
  created_at: str(t.created_at),
  branch_id: str(t.branch_id),
})

const mapAccount = (a: RawAccount): Account => ({
  id: num(a.id),
  code: str(a.code),
  name: str(a.name),
  type: status(a.type, 'asset'),
  description: str(a.description),
  created_at: str(a.created_at),
})

const mapJournalLine = (l: RawJournalLine): JournalLine => ({
  account_id: num(l.account_id),
  debit: str(l.debit) || '0',
  credit: str(l.credit) || '0',
})

const mapJournalEntry = (j: RawJournalEntry): JournalEntry => ({
  id: num(j.id),
  reference_type: str(j.reference_type),
  reference_id: num(j.reference_id),
  description: str(j.description),
  lines: (j.lines ?? []).map(mapJournalLine),
  created_at: str(j.created_at),
  creator: str(j.creator),
})

const mapAccountBalance = (b: RawAccountBalance): AccountBalance => ({
  account_id: num(b.account_id),
  code: str(b.code),
  name: str(b.name),
  type: str(b.type),
  debit: str(b.debit) || '0',
  credit: str(b.credit) || '0',
  balance: str(b.balance) || '0',
})

const mapLedgerLine = (l: RawLedgerLine): LedgerLine => ({
  journal_entry_id: num(l.journal_entry_id),
  reference_type: str(l.reference_type),
  reference_id: num(l.reference_id),
  description: str(l.description),
  debit: str(l.debit) || '0',
  credit: str(l.credit) || '0',
  balance: str(l.balance) || '0',
  created_at: str(l.created_at),
})

const mapPriceLevelReportItem = (i: RawPriceLevelReportItem): PriceLevelReportItem => ({
  product_id: num(i.product_id),
  product_name: str(i.product_name),
  sku: str(i.sku),
  base_price: str(i.base_price) || '0',
  base_unit: str(i.base_unit),
  price_levels: (i.price_levels ?? []).map(mapPriceLevel),
})

// ==================== POS / Retail query API ====================

/** Query/UnitAll — all units of measure (multi satuan) */
export async function fetchPosUnits(): Promise<PosUnit[]> {
  const data = await get<{ unit?: RawPosUnit[] }>('/bizchain/pos/unit')
  return (data.unit ?? []).map(mapPosUnit)
}

/** Query/ProductAll — all POS products (incl. bundles & price levels) */
export async function fetchPosProducts(): Promise<PosProduct[]> {
  const data = await get<{ product?: RawPosProductWithBranch[] }>('/bizchain/pos/product')
  return (data.product ?? []).map(mapPosProduct)
}

/** Query/AccountAll — chart of accounts */
export async function fetchPosAccounts(): Promise<Account[]> {
  const data = await get<{ account?: RawAccount[] }>('/bizchain/pos/account')
  return (data.account ?? []).map(mapAccount)
}

/** Query/JournalEntryAll — all journal entries */
export async function fetchJournalEntries(): Promise<JournalEntry[]> {
  const data = await get<{ journal_entry?: RawJournalEntry[] }>('/bizchain/pos/journal_entry')
  return (data.journal_entry ?? []).map(mapJournalEntry)
}

/** Query/TrialBalance — neraca saldo (akuntansi) */
export async function fetchTrialBalance(): Promise<TrialBalance> {
  const data = await get<RawTrialBalance>('/bizchain/pos/trial_balance')
  return {
    accounts: (data.accounts ?? []).map(mapAccountBalance),
    total_debit: str(data.total_debit) || '0',
    total_credit: str(data.total_credit) || '0',
  }
}

/** Query/IncomeStatement — laba rugi (akuntansi) */
export async function fetchIncomeStatement(): Promise<IncomeStatement> {
  const data = await get<RawIncomeStatement>('/bizchain/pos/income_statement')
  return {
    revenues: (data.revenues ?? []).map(mapAccountBalance),
    expenses: (data.expenses ?? []).map(mapAccountBalance),
    total_revenue: str(data.total_revenue) || '0',
    total_expense: str(data.total_expense) || '0',
    net_income: str(data.net_income) || '0',
  }
}

/** Query/BalanceSheet — neraca (akuntansi) */
export async function fetchBalanceSheet(): Promise<BalanceSheet> {
  const data = await get<RawBalanceSheet>('/bizchain/pos/balance_sheet')
  return {
    assets: (data.assets ?? []).map(mapAccountBalance),
    liabilities: (data.liabilities ?? []).map(mapAccountBalance),
    equities: (data.equities ?? []).map(mapAccountBalance),
    total_assets: str(data.total_assets) || '0',
    total_liabilities: str(data.total_liabilities) || '0',
    total_equity: str(data.total_equity) || '0',
  }
}

/** Query/Ledger — buku besar untuk satu akun */
export async function fetchLedger(accountId: number): Promise<Ledger> {
  const data = await get<RawLedger>(`/bizchain/pos/ledger/${accountId}`)
  return {
    account: mapAccount(data.account ?? {}),
    lines: (data.lines ?? []).map(mapLedgerLine),
    ending_balance: str(data.ending_balance) || '0',
  }
}

/** Query/PriceLevelReport — laporan harga per level */
export async function fetchPriceLevelReport(): Promise<PriceLevelReportItem[]> {
  const data = await get<{ items?: RawPriceLevelReportItem[] }>('/bizchain/pos/price_level_report')
  return (data.items ?? []).map(mapPriceLevelReportItem)
}

/** Query/TransactionAll — all POS transactions */
export async function fetchPosTransactions(): Promise<Transaction[]> {
  const data = await get<{ transaction?: RawTransactionWithBranch[] }>('/bizchain/pos/transaction')
  return (data.transaction ?? []).map(mapTransaction)
}

/** Query/ProductByBranch — POS products filtered by branch */
export async function fetchPosProductsByBranch(branchId: string): Promise<PosProduct[]> {
  const data = await get<{ product?: RawPosProductWithBranch[] }>(`/bizchain/pos/product_by_branch/${branchId}`)
  return (data.product ?? []).map(mapPosProduct)
}

/** Query/TransactionByBranch — POS transactions filtered by branch */
export async function fetchPosTransactionsByBranch(branchId: string): Promise<Transaction[]> {
  const data = await get<{ transaction?: RawTransactionWithBranch[] }>(`/bizchain/pos/transaction_by_branch/${branchId}`)
  return (data.transaction ?? []).map(mapTransaction)
}

export { ChainApiError }

