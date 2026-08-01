# Multi Branch (Multi Cabang)

## Overview

RetailChain mendukung pengelolaan banyak cabang (Branch) dalam satu Group (Toko). Seluruh cabang menggunakan blockchain yang sama sehingga data selalu tersinkronisasi dan memiliki riwayat transaksi yang dapat diaudit.

```
Group (TOKO ABC)
│
├── Cabang Jakarta
├── Cabang Bandung
├── Cabang Surabaya
└── Cabang Semarang
```

---

## Konsep

Satu **Group** merepresentasikan satu perusahaan atau toko.

Di dalam Group dapat dibuat satu atau lebih **Branch (Cabang)**.

Setiap Branch memiliki:

- Nama Cabang
- Kode Cabang
- Gudang
- Kasir
- Admin Cabang
- Stok Barang
- Transaksi Penjualan

Semua data tetap berada dalam blockchain yang sama, tetapi dipisahkan berdasarkan `branch_id`.

---

## Struktur Data

```text
Group
└── Branch
    ├── branch_id
    ├── name
    ├── address
    ├── manager
    ├── warehouse
    └── status
```

---

## Contoh Transaksi

### Tambah Stok

```json
{
  "group_id": "abc",
  "branch_id": "SMG",
  "type": "ADD_STOCK",
  "product_id": "P001",
  "qty": 100
}
```

### Penjualan

```json
{
  "group_id": "abc",
  "branch_id": "SMG",
  "type": "SALE",
  "invoice": "INV-001",
  "items": [
    {
      "product_id": "P001",
      "qty": 2
    }
  ]
}
```

---

## Laporan

Laporan dapat difilter berdasarkan:

- Seluruh Group
- Cabang tertentu
- Periode waktu

Contoh:

```
TOKO ABC

Total Penjualan
Rp120.000.000

Cabang Semarang
Rp40.000.000

Cabang Jakarta
Rp35.000.000

Cabang Bandung
Rp25.000.000

Cabang Surabaya
Rp20.000.000
```

---

## Keuntungan

- Satu blockchain melayani banyak cabang.
- Setiap cabang memiliki stok dan transaksi sendiri.
- Owner dapat melihat laporan seluruh cabang secara real-time.
- Audit transaksi tetap terjaga karena seluruh aktivitas tercatat di blockchain.
- Mendukung transfer stok antar cabang dengan riwayat yang transparan.