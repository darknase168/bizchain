import React, { useState } from 'react'
import {
  GraduationCap, Users, BookOpen, CreditCard, Award, Plus,
  Search, CheckCircle, Clock, AlertCircle, FileText, BadgeCheck,
  Calendar, DollarSign, TrendingUp, X
} from 'lucide-react'
import { Student, Course, TuitionPayment, Certificate, Grade, formatCoin, formatDate } from './types'

// Mock data
const mockStudents: Student[] = [
  { id: 1, nis: '2025001', name: 'Ahmad Rizki', class: 'XII IPA 1', major: 'IPA', parent_name: 'Bapak Rizki', parent_phone: '081234567890', address: 'Jakarta Selatan', wallet_address: 'rupiah1stu1...', status: 'active', enrolled_at: '2023-07-15T00:00:00Z' },
  { id: 2, nis: '2025002', name: 'Siti Nurhaliza', class: 'XII IPA 2', major: 'IPA', parent_name: 'Ibu Siti', parent_phone: '081234567891', address: 'Jakarta Timur', wallet_address: 'rupiah1stu2...', status: 'active', enrolled_at: '2023-07-15T00:00:00Z' },
  { id: 3, nis: '2025003', name: 'Budi Hartono', class: 'XI IPS 1', major: 'IPS', parent_name: 'Bapak Hartono', parent_phone: '081234567892', address: 'Depok', wallet_address: 'rupiah1stu3...', status: 'active', enrolled_at: '2024-07-15T00:00:00Z' },
  { id: 4, nis: '2025004', name: 'Dewi Lestari', class: 'X Bahasa 1', major: 'Bahasa', parent_name: 'Ibu Lestari', parent_phone: '081234567893', address: 'Bekasi', wallet_address: 'rupiah1stu4...', status: 'active', enrolled_at: '2025-07-15T00:00:00Z' },
  { id: 5, nis: '2024015', name: 'Eko Prasetyo', class: 'Alumni', major: 'IPA', parent_name: 'Bapak Prasetyo', parent_phone: '081234567894', address: 'Tangerang', wallet_address: 'rupiah1stu5...', status: 'graduated', enrolled_at: '2022-07-15T00:00:00Z' },
]

const mockCourses: Course[] = [
  { id: 1, code: 'MTK-101', name: 'Matematika Wajib', description: 'Matematika dasar dan lanjutan', credits: 4, teacher: 'Ibu Susi, S.Pd', schedule: 'Senin & Rabu 07:00-08:30', semester: 'Ganjil 2026/2027', max_students: 36, enrolled_students: 32, fee: '500000', status: 'open' },
  { id: 2, code: 'FIS-201', name: 'Fisika', description: 'Mekanika, listrik, dan magnet', credits: 3, teacher: 'Pak Joko, M.Si', schedule: 'Selasa & Kamis 08:30-10:00', semester: 'Ganjil 2026/2027', max_students: 30, enrolled_students: 28, fee: '450000', status: 'open' },
  { id: 3, code: 'ING-101', name: 'Bahasa Inggris', description: 'English for academic purposes', credits: 3, teacher: 'Mrs. Diana, S.S', schedule: 'Rabu & Jumat 10:00-11:30', semester: 'Ganjil 2026/2027', max_students: 40, enrolled_students: 40, fee: '400000', status: 'closed' },
  { id: 4, code: 'KOM-301', name: 'Informatika', description: 'Pemrograman dan blockchain dasar', credits: 2, teacher: 'Pak Andi, M.Kom', schedule: 'Kamis 13:00-15:00', semester: 'Ganjil 2026/2027', max_students: 25, enrolled_students: 20, fee: '600000', status: 'open' },
]

const mockPayments: TuitionPayment[] = [
  { id: 1, student_id: 1, student_name: 'Ahmad Rizki', type: 'spp', amount: '500000', semester: 'Ganjil', academic_year: '2026/2027', status: 'paid', due_date: '2026-07-10T00:00:00Z', paid_at: '2026-07-08T10:30:00Z', tx_hash: 'ABC123DEF456...' },
  { id: 2, student_id: 2, student_name: 'Siti Nurhaliza', type: 'spp', amount: '500000', semester: 'Ganjil', academic_year: '2026/2027', status: 'paid', due_date: '2026-07-10T00:00:00Z', paid_at: '2026-07-09T14:20:00Z', tx_hash: 'DEF456GHI789...' },
  { id: 3, student_id: 3, student_name: 'Budi Hartono', type: 'spp', amount: '500000', semester: 'Ganjil', academic_year: '2026/2027', status: 'pending', due_date: '2026-08-10T00:00:00Z', paid_at: null, tx_hash: '' },
  { id: 4, student_id: 4, student_name: 'Dewi Lestari', type: 'uang_gedung', amount: '5000000', semester: 'Ganjil', academic_year: '2026/2027', status: 'paid', due_date: '2026-07-20T00:00:00Z', paid_at: '2026-07-15T09:00:00Z', tx_hash: 'GHI789JKL012...' },
  { id: 5, student_id: 3, student_name: 'Budi Hartono', type: 'praktikum', amount: '750000', semester: 'Ganjil', academic_year: '2026/2027', status: 'overdue', due_date: '2026-07-25T00:00:00Z', paid_at: null, tx_hash: '' },
]

const mockCertificates: Certificate[] = [
  { id: 1, student_id: 5, student_name: 'Eko Prasetyo', type: 'diploma', title: 'Ijazah SMA - Jurusan IPA', description: 'Telah menyelesaikan pendidikan SMA dengan predikat Sangat Memuaskan', issued_at: '2025-06-20T00:00:00Z', tx_hash: 'JKL012MNO345...', verified: true, ipfs_hash: 'QmX7Y8Z9...' },
  { id: 2, student_id: 1, student_name: 'Ahmad Rizki', type: 'achievement', title: 'Juara 1 Olimpiade Matematika', description: 'Olimpiade Matematika Tingkat Provinsi 2026', issued_at: '2026-05-15T00:00:00Z', tx_hash: 'MNO345PQR678...', verified: true, ipfs_hash: 'QmA1B2C3...' },
  { id: 3, student_id: 2, student_name: 'Siti Nurhaliza', type: 'certificate', title: 'Sertifikat English Proficiency', description: 'TOEFL ITP Score 580', issued_at: '2026-04-10T00:00:00Z', tx_hash: 'PQR678STU901...', verified: true, ipfs_hash: 'QmD4E5F6...' },
]

const mockGrades: Grade[] = [
  { id: 1, student_id: 1, course_id: 1, course_name: 'Matematika Wajib', semester: 'Ganjil 2025/2026', score: 92, grade: 'A', teacher: 'Ibu Susi, S.Pd', recorded_at: '2026-01-15T00:00:00Z', tx_hash: 'STU901VWX234...' },
  { id: 2, student_id: 1, course_id: 2, course_name: 'Fisika', semester: 'Ganjil 2025/2026', score: 88, grade: 'A-', teacher: 'Pak Joko, M.Si', recorded_at: '2026-01-15T00:00:00Z', tx_hash: 'VWX234YZA567...' },
  { id: 3, student_id: 2, course_id: 1, course_name: 'Matematika Wajib', semester: 'Ganjil 2025/2026', score: 95, grade: 'A', teacher: 'Ibu Susi, S.Pd', recorded_at: '2026-01-15T00:00:00Z', tx_hash: 'YZA567BCD890...' },
]

const paymentTypeLabels: Record<string, string> = {
  spp: 'SPP Bulanan',
  uang_gedung: 'Uang Gedung',
  praktikum: 'Biaya Praktikum',
  wisuda: 'Biaya Wisuda',
  lainnya: 'Lainnya',
}

const certTypeLabels: Record<string, string> = {
  diploma: 'Ijazah',
  certificate: 'Sertifikat',
  transcript: 'Transkrip',
  achievement: 'Prestasi',
}

const EducationDashboard: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'students' | 'courses' | 'payments' | 'certificates' | 'grades'>('students')
  const [searchTerm, setSearchTerm] = useState('')
  const [showAddStudent, setShowAddStudent] = useState(false)
  const [showPayModal, setShowPayModal] = useState(false)
  const [selectedPayment, setSelectedPayment] = useState<TuitionPayment | null>(null)
  const [newStudent, setNewStudent] = useState({ name: '', class: '', major: '', parent_name: '', parent_phone: '', address: '' })
  const [paySuccess, setPaySuccess] = useState(false)

  const filteredStudents = mockStudents.filter(s =>
    s.name.toLowerCase().includes(searchTerm.toLowerCase()) || s.nis.includes(searchTerm)
  )

  const totalPaid = mockPayments.filter(p => p.status === 'paid').reduce((sum, p) => sum + parseInt(p.amount), 0)
  const totalPending = mockPayments.filter(p => p.status === 'pending' || p.status === 'overdue').reduce((sum, p) => sum + parseInt(p.amount), 0)

  const handleAddStudent = () => {
    if (!newStudent.name || !newStudent.class) return
    setShowAddStudent(false)
    setNewStudent({ name: '', class: '', major: '', parent_name: '', parent_phone: '', address: '' })
  }

  const handlePay = () => {
    setPaySuccess(true)
    setTimeout(() => {
      setPaySuccess(false)
      setShowPayModal(false)
      setSelectedPayment(null)
    }, 2000)
  }

  return (
    <div className="space-y-6 animate-fade-in">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-white">Education Dashboard</h1>
          <p className="text-surface-400 mt-1">School & Course Management on Blockchain</p>
        </div>
        <div className="flex gap-2 flex-wrap">
          {(['students', 'courses', 'payments', 'certificates', 'grades'] as const).map(tab => (
            <button key={tab} className={`px-4 py-2 rounded-xl text-sm font-medium transition-all capitalize ${activeTab === tab ? 'bg-primary-500/20 text-primary-400 border border-primary-500/30' : 'bg-white/5 text-surface-400 hover:text-white'}`} onClick={() => setActiveTab(tab)}>
              {tab === 'students' && <Users size={16} className="inline mr-2" />}
              {tab === 'courses' && <BookOpen size={16} className="inline mr-2" />}
              {tab === 'payments' && <CreditCard size={16} className="inline mr-2" />}
              {tab === 'certificates' && <Award size={16} className="inline mr-2" />}
              {tab === 'grades' && <FileText size={16} className="inline mr-2" />}
              {tab}
            </button>
          ))}
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="glass-card p-5">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-10 h-10 rounded-xl bg-primary-500/10 flex items-center justify-center"><Users size={20} className="text-primary-400" /></div>
            <span className="text-surface-400 text-sm">Active Students</span>
          </div>
          <p className="text-2xl font-bold text-white">{mockStudents.filter(s => s.status === 'active').length}</p>
        </div>
        <div className="glass-card p-5">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-10 h-10 rounded-xl bg-accent-500/10 flex items-center justify-center"><BookOpen size={20} className="text-accent-400" /></div>
            <span className="text-surface-400 text-sm">Open Courses</span>
          </div>
          <p className="text-2xl font-bold text-white">{mockCourses.filter(c => c.status === 'open').length}</p>
        </div>
        <div className="glass-card p-5">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-10 h-10 rounded-xl bg-green-500/10 flex items-center justify-center"><TrendingUp size={20} className="text-green-400" /></div>
            <span className="text-surface-400 text-sm">Revenue (Paid)</span>
          </div>
          <p className="text-2xl font-bold text-white">{formatCoin(totalPaid.toString(), 'uidr')}</p>
        </div>
        <div className="glass-card p-5">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-10 h-10 rounded-xl bg-yellow-500/10 flex items-center justify-center"><AlertCircle size={20} className="text-yellow-400" /></div>
            <span className="text-surface-400 text-sm">Outstanding</span>
          </div>
          <p className="text-2xl font-bold text-white">{formatCoin(totalPending.toString(), 'uidr')}</p>
        </div>
      </div>

      {paySuccess && (
        <div className="glass-card p-4 border-primary-500/20 bg-primary-500/10 flex items-center gap-3 animate-slide-down">
          <CheckCircle size={20} className="text-primary-400" />
          <span className="text-primary-300">Payment recorded on blockchain successfully! 0 fees charged.</span>
        </div>
      )}

      {activeTab === 'students' && (
        <div className="space-y-4">
          <div className="flex gap-3">
            <div className="relative flex-1">
              <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-surface-400" />
              <input className="input-field pl-10" placeholder="Search by name or NIS..." value={searchTerm} onChange={(e) => setSearchTerm(e.target.value)} />
            </div>
            <button className="btn-primary" onClick={() => setShowAddStudent(true)}><Plus size={18} />Add Student</button>
          </div>

          {showAddStudent && (
            <div className="glass-card p-6 animate-slide-down">
              <h3 className="text-lg font-semibold text-white mb-4">Register New Student</h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="block text-sm text-surface-400 mb-2">Full Name *</label>
                  <input className="input-field" placeholder="Student name" value={newStudent.name} onChange={(e) => setNewStudent(p => ({ ...p, name: e.target.value }))} />
                </div>
                <div>
                  <label className="block text-sm text-surface-400 mb-2">Class *</label>
                  <input className="input-field" placeholder="X IPA 1" value={newStudent.class} onChange={(e) => setNewStudent(p => ({ ...p, class: e.target.value }))} />
                </div>
                <div>
                  <label className="block text-sm text-surface-400 mb-2">Major</label>
                  <select className="input-field" value={newStudent.major} onChange={(e) => setNewStudent(p => ({ ...p, major: e.target.value }))}>
                    <option value="">Select major</option>
                    <option value="IPA">IPA</option>
                    <option value="IPS">IPS</option>
                    <option value="Bahasa">Bahasa</option>
                  </select>
                </div>
                <div>
                  <label className="block text-sm text-surface-400 mb-2">Parent Name</label>
                  <input className="input-field" placeholder="Parent/Guardian name" value={newStudent.parent_name} onChange={(e) => setNewStudent(p => ({ ...p, parent_name: e.target.value }))} />
                </div>
                <div>
                  <label className="block text-sm text-surface-400 mb-2">Parent Phone</label>
                  <input className="input-field" placeholder="08xxxxxxxxxx" value={newStudent.parent_phone} onChange={(e) => setNewStudent(p => ({ ...p, parent_phone: e.target.value }))} />
                </div>
                <div>
                  <label className="block text-sm text-surface-400 mb-2">Address</label>
                  <input className="input-field" placeholder="Home address" value={newStudent.address} onChange={(e) => setNewStudent(p => ({ ...p, address: e.target.value }))} />
                </div>
              </div>
              <div className="flex gap-3">
                <button className="btn-primary" onClick={handleAddStudent}>Register Student</button>
                <button className="btn-secondary" onClick={() => setShowAddStudent(false)}>Cancel</button>
              </div>
            </div>
          )}

          <div className="glass-card overflow-hidden">
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-white/5">
                    <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">NIS</th>
                    <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Name</th>
                    <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Class</th>
                    <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Major</th>
                    <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Parent</th>
                    <th className="text-center text-xs text-surface-400 font-medium px-4 py-3">Status</th>
                  </tr>
                </thead>
                <tbody>
                  {filteredStudents.map(student => (
                    <tr key={student.id} className="border-b border-white/5 hover:bg-white/5 transition-colors">
                      <td className="px-4 py-3 text-sm text-surface-400 font-mono">{student.nis}</td>
                      <td className="px-4 py-3">
                        <div className="flex items-center gap-2">
                          <div className="w-8 h-8 rounded-lg bg-primary-500/10 flex items-center justify-center"><GraduationCap size={14} className="text-primary-400" /></div>
                          <p className="text-sm text-white">{student.name}</p>
                        </div>
                      </td>
                      <td className="px-4 py-3 text-sm text-surface-300">{student.class}</td>
                      <td className="px-4 py-3"><span className="badge-info">{student.major}</span></td>
                      <td className="px-4 py-3 text-sm text-surface-300">{student.parent_name}</td>
                      <td className="px-4 py-3 text-center">
                        <span className={`badge ${student.status === 'active' ? 'badge-success' : student.status === 'graduated' ? 'badge-info' : ''}`}>{student.status}</span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}

      {activeTab === 'courses' && (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {mockCourses.map(course => (
            <div key={course.id} className="glass-card p-6">
              <div className="flex items-start justify-between mb-3">
                <div>
                  <span className="text-xs text-surface-500 font-mono">{course.code}</span>
                  <h3 className="text-lg font-semibold text-white">{course.name}</h3>
                  <p className="text-sm text-surface-400">{course.description}</p>
                </div>
                <span className={`badge ${course.status === 'open' ? 'badge-success' : course.status === 'closed' ? '' : 'badge-info'}`}>{course.status}</span>
              </div>
              <div className="space-y-2 text-sm mb-4">
                <div className="flex items-center gap-2 text-surface-400"><Users size={14} /><span>{course.teacher}</span></div>
                <div className="flex items-center gap-2 text-surface-400"><Calendar size={14} /><span>{course.schedule}</span></div>
                <div className="flex items-center gap-2 text-surface-400"><BookOpen size={14} /><span>{course.credits} credits • {course.semester}</span></div>
                <div className="flex items-center gap-2 text-surface-400"><DollarSign size={14} /><span>{formatCoin(course.fee, 'uidr')}</span></div>
              </div>
              <div className="flex items-center justify-between">
                <div className="flex-1 mr-4">
                  <div className="flex justify-between text-xs text-surface-500 mb-1">
                    <span>Enrollment</span>
                    <span>{course.enrolled_students}/{course.max_students}</span>
                  </div>
                  <div className="w-full bg-white/5 rounded-full h-2">
                    <div className="bg-gradient-to-r from-primary-500 to-accent-500 h-2 rounded-full" style={{ width: `${(course.enrolled_students / course.max_students) * 100}%` }} />
                  </div>
                </div>
                {course.status === 'open' && <button className="btn-primary text-sm">Enroll</button>}
              </div>
            </div>
          ))}
        </div>
      )}

      {activeTab === 'payments' && (
        <div className="space-y-4">
          <div className="glass-card overflow-hidden">
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-white/5">
                    <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Student</th>
                    <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Type</th>
                    <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Period</th>
                    <th className="text-right text-xs text-surface-400 font-medium px-4 py-3">Amount</th>
                    <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Due Date</th>
                    <th className="text-center text-xs text-surface-400 font-medium px-4 py-3">Status</th>
                    <th className="text-center text-xs text-surface-400 font-medium px-4 py-3">Action</th>
                  </tr>
                </thead>
                <tbody>
                  {mockPayments.map(payment => (
                    <tr key={payment.id} className="border-b border-white/5 hover:bg-white/5 transition-colors">
                      <td className="px-4 py-3 text-sm text-white">{payment.student_name}</td>
                      <td className="px-4 py-3 text-sm text-surface-300">{paymentTypeLabels[payment.type]}</td>
                      <td className="px-4 py-3 text-sm text-surface-300">{payment.semester} {payment.academic_year}</td>
                      <td className="px-4 py-3 text-sm text-primary-400 text-right font-medium">{formatCoin(payment.amount, 'uidr')}</td>
                      <td className="px-4 py-3 text-sm text-surface-400">{formatDate(payment.due_date)}</td>
                      <td className="px-4 py-3 text-center">
                        <span className={`badge ${payment.status === 'paid' ? 'badge-success' : payment.status === 'overdue' ? 'bg-red-500/10 text-red-400' : 'badge-info'}`}>
                          {payment.status === 'paid' && <CheckCircle size={12} className="inline mr-1" />}
                          {payment.status === 'overdue' && <AlertCircle size={12} className="inline mr-1" />}
                          {payment.status === 'pending' && <Clock size={12} className="inline mr-1" />}
                          {payment.status}
                        </span>
                      </td>
                      <td className="px-4 py-3 text-center">
                        {payment.status !== 'paid' && (
                          <button className="btn-primary text-xs px-3 py-1.5" onClick={() => { setSelectedPayment(payment); setShowPayModal(true) }}>Pay Now</button>
                        )}
                        {payment.status === 'paid' && payment.tx_hash && (
                          <span className="text-xs text-surface-500 font-mono" title={payment.tx_hash}>TX: {payment.tx_hash.slice(0, 8)}...</span>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          {showPayModal && selectedPayment && (
            <div className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center p-4" onClick={() => setShowPayModal(false)}>
              <div className="glass-card p-6 max-w-md w-full animate-slide-down" onClick={(e) => e.stopPropagation()}>
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-lg font-semibold text-white">Pay Tuition</h3>
                  <button onClick={() => setShowPayModal(false)} className="p-1 rounded-lg hover:bg-white/10"><X size={18} className="text-surface-400" /></button>
                </div>
                <div className="space-y-3 mb-6">
                  <div className="flex justify-between text-sm"><span className="text-surface-400">Student</span><span className="text-white">{selectedPayment.student_name}</span></div>
                  <div className="flex justify-between text-sm"><span className="text-surface-400">Type</span><span className="text-white">{paymentTypeLabels[selectedPayment.type]}</span></div>
                  <div className="flex justify-between text-sm"><span className="text-surface-400">Period</span><span className="text-white">{selectedPayment.semester} {selectedPayment.academic_year}</span></div>
                  <div className="flex justify-between text-sm"><span className="text-surface-400">Fee</span><span className="text-primary-400">0 IDR (Feeless)</span></div>
                  <div className="border-t border-white/10 pt-3 flex justify-between"><span className="text-white font-medium">Total</span><span className="text-primary-400 font-bold text-lg">{formatCoin(selectedPayment.amount, 'uidr')}</span></div>
                </div>
                <button className="btn-primary w-full justify-center" onClick={handlePay}>
                  <CreditCard size={18} />Pay with RUPIAH (0 Fee)
                </button>
              </div>
            </div>
          )}
        </div>
      )}

      {activeTab === 'certificates' && (
        <div className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {mockCertificates.map(cert => (
              <div key={cert.id} className="glass-card p-6">
                <div className="flex items-start justify-between mb-3">
                  <div className="flex items-center gap-3">
                    <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-primary-500/10 to-accent-500/10 flex items-center justify-center">
                      <Award size={22} className="text-primary-400" />
                    </div>
                    <div>
                      <span className="text-xs text-surface-500">{certTypeLabels[cert.type]}</span>
                      <h3 className="text-base font-semibold text-white">{cert.title}</h3>
                    </div>
                  </div>
                  {cert.verified && (
                    <span className="badge-success flex items-center gap-1"><BadgeCheck size={12} />Verified</span>
                  )}
                </div>
                <p className="text-sm text-surface-400 mb-3">{cert.description}</p>
                <div className="space-y-1.5 text-xs text-surface-500">
                  <div className="flex justify-between"><span>Student</span><span className="text-surface-300">{cert.student_name}</span></div>
                  <div className="flex justify-between"><span>Issued</span><span className="text-surface-300">{formatDate(cert.issued_at)}</span></div>
                  <div className="flex justify-between"><span>TX Hash</span><span className="text-surface-300 font-mono">{cert.tx_hash}</span></div>
                  <div className="flex justify-between"><span>IPFS</span><span className="text-surface-300 font-mono">{cert.ipfs_hash}</span></div>
                </div>
                <button className="btn-secondary w-full justify-center mt-4 text-sm">
                  <BadgeCheck size={16} />Verify on Blockchain
                </button>
              </div>
            ))}
          </div>
        </div>
      )}

      {activeTab === 'grades' && (
        <div className="glass-card overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-white/5">
                  <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Student</th>
                  <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Course</th>
                  <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Semester</th>
                  <th className="text-center text-xs text-surface-400 font-medium px-4 py-3">Score</th>
                  <th className="text-center text-xs text-surface-400 font-medium px-4 py-3">Grade</th>
                  <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">Teacher</th>
                  <th className="text-left text-xs text-surface-400 font-medium px-4 py-3">TX Hash</th>
                </tr>
              </thead>
              <tbody>
                {mockGrades.map(grade => (
                  <tr key={grade.id} className="border-b border-white/5 hover:bg-white/5 transition-colors">
                    <td className="px-4 py-3 text-sm text-white">{mockStudents.find(s => s.id === grade.student_id)?.name}</td>
                    <td className="px-4 py-3 text-sm text-surface-300">{grade.course_name}</td>
                    <td className="px-4 py-3 text-sm text-surface-300">{grade.semester}</td>
                    <td className="px-4 py-3 text-sm text-white text-center font-medium">{grade.score}</td>
                    <td className="px-4 py-3 text-center">
                      <span className={`badge ${grade.grade.startsWith('A') ? 'badge-success' : 'badge-info'}`}>{grade.grade}</span>
                    </td>
                    <td className="px-4 py-3 text-sm text-surface-300">{grade.teacher}</td>
                    <td className="px-4 py-3 text-xs text-surface-500 font-mono">{grade.tx_hash}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  )
}

export default EducationDashboard
