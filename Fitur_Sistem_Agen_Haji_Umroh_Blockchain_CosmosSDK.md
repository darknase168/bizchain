# Sistem Agen Haji & Umroh Berbasis Blockchain (Cosmos SDK)

## 1. Identitas Jamaah (DID)

-   ID jamaah berbasis blockchain.
-   Riwayat keberangkatan.
-   Riwayat vaksin.
-   Riwayat pembayaran.
-   Riwayat visa.
-   Hash dokumen untuk verifikasi.

## 2. Smart Contract Paket

-   Nama paket
-   Harga
-   Jadwal
-   Kuota
-   Hotel
-   Maskapai
-   Muthawif
-   Status keberangkatan
-   Penutupan kuota otomatis.

## 3. Escrow Pembayaran

Dana disimpan smart contract dan dicairkan bertahap sesuai progres
(visa, tiket, hotel).

## 4. NFT Tiket

NFT berisi data paket, kursi, hotel, jadwal, dan QR Code.

## 5. Token Loyalitas

Token reward untuk cashback, diskon, dan referral.

## 6. Referral On-chain

Komisi dihitung dan dibagikan otomatis melalui smart contract.

## 7. Marketplace Paket

Menampilkan paket, harga, rating, review, jadwal, dan kuota.

## 8. Verifikasi Dokumen

Hash blockchain untuk paspor, visa, vaksin, dan surat kesehatan.

## 9. Tracking Proses

Status: 1. Daftar 2. DP Dibayar 3. Visa Diproses 4. Visa Terbit 5. Tiket
Terbit 6. Hotel Confirm 7. Manasik 8. Berangkat 9. Pulang

## 10. Wallet Jamaah

Saldo, token, NFT, bukti pembayaran, dan cicilan.

## 11. Pembayaran Cicilan

DP, cicilan berkala, pengingat, dan aturan keterlambatan.

## 12. Multi Agen

Mendukung agen pusat, cabang, dan subagen dalam satu jaringan.

## 13. Audit Transparan

Laporan transaksi, refund, pembatalan, dan pendapatan.

## 14. Tracking Bagasi

Status bagasi melalui QR/NFC.

## 15. Marketplace Oleh-oleh

Pre-order oleh-oleh dengan pembayaran melalui wallet.

## 16. Asuransi Digital

Polis digital dan proses klaim transparan.

## 17. Voting DAO

Voting untuk keputusan organisasi atau asosiasi agen.

## 18. Rekam Jejak Agen

Skor berdasarkan performa, rating, dan komplain.

## 19. Integrasi Cosmos IBC

Transfer aset dan interoperabilitas dengan blockchain Cosmos lain.

## 20. Modul Cosmos SDK

-   x-jamaah
-   x-paket
-   x-pembayaran
-   x-visa
-   x-hotel
-   x-ticket
-   x-referral
-   x-reward
-   x-audit
-   x-keberangkatan

## Arsitektur

``` text
Mobile Jamaah
      |
 REST / gRPC
      |
+----------------------+
| Cosmos SDK Blockchain|
| x-jamaah             |
| x-paket              |
| x-pembayaran         |
| x-visa               |
| x-hotel              |
| x-ticket             |
| x-referral           |
| x-reward             |
| x-audit              |
+----------------------+
      |      |      |
 Hotel API Airline Payment Gateway
```

## Rekomendasi

-   Simpan dokumen pribadi secara off-chain (object storage/IPFS
    terenkripsi).
-   Simpan hash dan metadata di blockchain.
-   Gunakan permissioned validator untuk agen resmi dan regulator.
