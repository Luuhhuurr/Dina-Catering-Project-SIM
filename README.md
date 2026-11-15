# SIM Dina Catering

Sistem Informasi Manajemen (SIM) untuk usaha katering **Dina Catering**, Ungaran,
Kabupaten Semarang. Aplikasi web untuk mengelola bahan baku, resep, paket menu,
HPP, pesanan, pembayaran, jadwal produksi, dan laporan keuangan.

Dibangun sebagai Final Proyek dengan **Go 1.26 + chi + sqlc (pgx) + PostgreSQL 16 +
HTMX**. Tanpa ORM, tanpa framework frontend.

---

## 1. Menjalankan aplikasi

```bash
git clone <repo>
cd dina-catering
cp .env.example .env      # opsional, docker compose sudah punya nilai bawaan
docker compose up --build
```

Buka **http://localhost:8080**.

| Akun | Password |
|------|----------|
| `admin` | `admin123` |

Database, migrasi, dan data awal (seed) dibuat otomatis saat container pertama kali
dibangun. Menghapus volume `pgdata` akan mengembalikan database ke kondisi awal.

Untuk menghentikan:

```bash
docker compose down          # matikan, data tetap tersimpan
docker compose down -v       # matikan dan hapus data
```

### Menjalankan di luar Docker

```bash
go run ./cmd/server
```

Butuh PostgreSQL yang sudah berjalan. Set variabel lingkungan:

```bash
export DATABASE_URL="postgres://dina:dina@localhost:5432/dinacatering?sslmode=disable"
export SESSION_SECRET="secret-minimal-untuk-pengembangan"
export PORT=8080
```

---

## 2. Fitur

### A. Kalkulator HPP & Resep (fitur prioritas)

Ini adalah inti sistem. Owner complained: *"harga bahan naik, tapi harga paket tetap —
saya harus tahu HPP per porsi dan untung berapa."*

* Setiap menu punya **resep** (bahan + takaran per porsi).
* Setiap bahan punya **riwayat harga** bertanggal berlaku, jadi HPP historis bisa
  dihitung ulang.
* **HPP per porsi** = Σ(takaran × harga efektif bahan) × (1 + % biaya operasional).
* **Margin** = (harga jual − HPP) / harga jual.
* **Harga jual minimum 30%** = HPP / 0,7.
* Paket dengan margin di bawah 30% diberi **badge peringatan** di dashboard, daftar
  paket, dan laporan.
* Mengubah harga satu bahan langsung mengubah HPP dan margin seluruh paket.

### B. Pelanggan

Data pelanggan (nama, kontak WhatsApp, alamat, catatan) plus jumlah pesanan.

### C. Pesanan & Pembayaran

* Sumber pesanan: Instagram, Mulut ke Mulut, Hotel/Wedding, Lainnya.
* Alur status: `baru` → `final_meeting` → `survei_lokasi` → `produksi` →
  `pengiriman` → `selesai` (atau `batal`).
* Tanggal acara pada hari Sabtu/Minggu memakai **harga akhir pekan**.
* **Tenggat DP (H-30)** dan **tenggat pelunasan (H-7)** dihitung otomatis.
* DP 50% dicatat lewat form pembayaran; pesanan yang belum lunas muncul di daftar
  **"menunggu pembayaran"** dengan keterangan jatuh tempo/terlambat.
* Flag **masak di lokasi**.

### D. Pembelian & Stok

Pencatatan pembelian di pasar tradisional (nama pasar, jumlah, harga satuan).
Satu transaksi irrevocably: stok bahan bertambah, harga terakhir diperbarui, dan satu
baris riwayat harga ditambahkan — sehingga HPP ikut terhitung ulang.

### E. Jadwal Produksi

Papan produksi menggantikan papan tulis di dapur. Saat DP dicatat, pesanan naik
status dari `baru` ke `final_meeting` dan **tugas produksi default dibuat untuk
hari acara** (Persiapan bahan H-1, Masak H, Plating H, Pengiriman H). Saat status
pesanan diubah menjadi `produksi`, stok bahan otomatis berkurang sesuai resep ×
jumlah tamu. Pembuatan tugas bersifat idempoten, jadi tugas tidak pernah
digandakan.
Setiap tugas punya **penanggung jawab (PJ)** dan status pengerjaan.

### F. Laporan

* Laba rugi per bulan: pemasukan − pembelian bahan − biaya operasional.
* Rekap pesanan per status dan per sumber.
* Margin per pesanan.
* Stok bahan terkini.

### G. Login

Login berbasis cookie session bertanda tangan (HMAC-SHA256), bcrypt untuk password.
Kolom `role` sudah disiapkan untuk pembatasan hak akses di masa depan.

---

## 3. Rumus perhitungan

Semua rumus ada di `internal/service` dan tersedia unit test.

```
hargaEfektif(bahan, t) = harga terakhir pada harga_bahan_history dengan tanggal_berlaku <= t
HPPBahan(menu, t)       = Σ (resep_item.jumlah_per_porsi × hargaEfektif(bahan, t))
HPPPerPorsi(menu, t)    = HPPBahan × (1 + persenOperasional / 100)
HPPPaket(paket, t)      = Σ HPPPerPorsi(menu di paket)
Margin(paket, t)        = (hargaJualEfektif − HPPPaket) / hargaJualEfektif × 100
HargaJualMinimum30(HPP) = HPP / 0,7
HargaJualEfektif(paket) = hari Sabtu/Minggu dan ada harga weekend ? harga weekend : harga biasa
TotalPesanan            = jumlah_tamu × HargaJualEfektif(paket, tanggal_acara)
JadwalPembayaran(acara) = (acara − 30 hari, acara − 7 hari)
StatusPembayaran        = belum_dp | dp_terbayar | lunas
```

**Asumsi biaya operasional.** Persentase default 10% (bisa diubah di halaman
Pengaturan) ditambahkan **proporsional ke HPP bahan**, baru margin dihitung dari HPP
total. Dengan cara ini rumus `HPP / 0,7` tetap konsisten. Laba rugi di laporan memakai
**biaya operasional aktual yang dicatat**, bukan estimasi persentase, sehingga kedua
perhitungan tidak tumpang tindih.

---

## 4. Struktur proyek

```
dina-catering/
├── cmd/
│   └── server/main.go            # titik masuk: config → connect → migrate → seed → router
├── internal/
│   ├── config/                   # baca environment variable
│   ├── database/                 # koneksi pool, migrasi (golang-migrate, embedded)
│   ├── handler/                  # HTTP handler (chi) + template per halaman
│   ├── middleware/               # logging, recover, auth, session (cookie bertanda tangan)
│   ├── model/                    # konstanta domain (status, sumber, margin minimum)
│   ├── seed/                     # data awal, hash password dihitung saat runtime
│   ├── service/                  # ⭐ logika bisnis: HPP, margin, jadwal bayar, stok, laporan
│   └── store/                    # kode hasil generate sqlc (tidak ditulis manual)
├── db/
│   ├── schema.sql                # DDL (dipakai sqlc)
│   ├── queries.sql               # query beranotasi (dipakai sqlc)
│   ├── seed.sql                  # data awal
│   └── migrations/               # migrasi golang-migrate (di-embed ke binary)
├── web/
│   ├── templates/
│   │   ├── base.layout.html      # layout + navbar
│   │   ├── pages/                # 13 halaman
│   │   └── partials/             # potongan yang di-refresh HTMX
│   └── static/                   # app.css, app.js
├── sqlc.yaml
├── Dockerfile
├── docker-compose.yml
├── .env.example
├── scripts/smoke-test.sh             # uji otomatis seluruh alur aplikasi
├── go.mod
└── README.md
```

Aturan yang dijaga:

* **Tidak ada SQL di dalam handler.** Handler hanya memakai `store` (hasil generate
  sqlc) dan `service`.
* **Tidak ada ORM.** Seluruh akses database lewat kode yang dihasilkan sqlc.
* **Semua perhitungan berada di `internal/service`** dan diuji dengan unit test.

### Regenerasi kode sqlc

```bash
sqlc generate
```

### Menjalankan pengujian

```bash
go vet ./...
go build ./...
go test ./...
```

Unit test yang tersedia: `internal/service/hpp_test.go` (HPP, margin, harga jual
minimum), `internal/service/pembayaran_test.go` (jadwal DP/pelunasan, status
pembayaran), `internal/handler/template_test.go` (semua halaman template valid).

### Uji otomatis alur aplikasi (smoke test)

Dengan aplikasi sudah jalan di `localhost:8080`:

```bash
bash scripts/smoke-test.sh
```

Skrip ini menguji 18 hal secara berurutan: redirect tanpa session, halaman
login, login + cookie session, dashboard, badge margin <30%, kalkulator HPP
(termasuk fallback halaman penuh tanpa JavaScript), perubahan harga bahan →
HPP berubah, pembuatan pesanan (total + tenggat H-30/H-7), daftar menunggu
pembayaran, pencatatan DP → status naik ke Final Meeting + tugas produksi, papan
produksi + kolom PJ, pembelian → stok dan riwayat harga, laporan laba rugi,
rekap per status dan sumber, seluruh halaman modul, endpoint tabel, panel
HTMX versus halaman penuh, serta logout.

---

## 5. ERD

```mermaid
erDiagram
    users ||--o{ pesanan : "login admin"
    pelanggan ||--o{ pesanan : "memesan"
    paket_menu ||--o{ pesanan : "dipilih"
    paket_menu ||--o{ paket_menu_item : "berisi"
    menu ||--o{ paket_menu_item : "tercantum"
    menu ||--o{ resep_item : "diresep"
    bahan_baku ||--o{ resep_item : "dipakai"
    bahan_baku ||--o{ harga_bahan_history : "riwayat harga"
    bahan_baku ||--o{ pembelian_bahan : "dibeli"
    pesanan ||--o{ pembayaran : "dibayar"
    pesanan ||--o{ tugas_produksi : "dikerjakan"

    users {
        int id PK
        text username UK
        text password_hash
        text role
    }
    pengaturan {
        text kunci PK
        text nilai
    }
    pelanggan {
        int id PK
        text nama
        text kontak
        text alamat
        text catatan
    }
    bahan_baku {
        int id PK
        text nama UK
        text satuan
        double stok_saat_ini
        bigint harga_terakhir
    }
    harga_bahan_history {
        int id PK
        int bahan_baku_id FK
        bigint harga_satuan
        date tanggal_berlaku
    }
    pembelian_bahan {
        int id PK
        int bahan_baku_id FK
        double jumlah
        bigint harga_satuan
        bigint total
        text pasar
        date tanggal
    }
    menu {
        int id PK
        text nama UK
        text deskripsi
    }
    resep_item {
        int id PK
        int menu_id FK
        int bahan_baku_id FK
        double jumlah_per_porsi
        text satuan
    }
    paket_menu {
        int id PK
        text nama
        bigint harga_per_porsi
        bigint harga_per_porsi_weekend
        int tahun_berlaku
    }
    paket_menu_item {
        int id PK
        int paket_menu_id FK
        int menu_id FK
    }
    pesanan {
        int id PK
        int pelanggan_id FK
        int paket_menu_id FK
        int jumlah_tamu
        date tanggal_acara
        text sumber_pesanan
        bigint total_harga
        text status
        boolean masak_di_lokasi
        date tgl_jatuh_tempo_dp
        date tgl_jatuh_tempo_lunas
    }
    pembayaran {
        int id PK
        int pesanan_id FK
        bigint jumlah
        text tipe
        date tanggal
        text metode
    }
    tugas_produksi {
        int id PK
        int pesanan_id FK
        text nama_tugas
        text pj
        date deadline
        text status
    }
    biaya_operasional {
        int id PK
        text kategori
        text deskripsi
        bigint jumlah
        date tanggal
    }
```

---

## 6. Catatan keputusan desain

Hal-hal yang tidak sepenuhnya ditentukan oleh spesifikasi, dan pilihan yang dipakai:

1. **Migrasi memakai golang-migrate** dengan file migrasi di-embed ke dalam binary
   (`db/migrations`). Tidak perlu executable `migrate` terpisah di image Docker.
2. **Harga akhir pekan** dimodelkan sebagai kolom nullable
   `paket_menu.harga_per_porsi_weekend` (NULL = harga biasa), bukan tabel terpisah.
3. **Harga paket per tahun**: `pesanan.paket_menu_id` menunjuk baris harga spesifik.
   Form pesanan memilih *nama paket* + tanggal acara; service mencari baris dengan
   `tahun_berlaku <= tahun acara` yang paling dekat. Satu pesanan terkunci ke satu
   baris harga.
4. **Session** dibuat sendiri: cookie berisi HMAC-SHA256 dari `userID|username|role|expiry`,
   dibandingkan dengan `hmac.Equal` (constant-time). Tidak ada tabel session dan tidak
   ada dependensi tambahan.
5. **Tabel `pengaturan`** (key/value) dipakai untuk persentase biaya operasional supaya
   bisa diubah dari UI, bukan dari environment variable.
6. **Tidak ada kolom `kode` pesanan.** UI menampilkan `#<id>` supaya tidak ada transaksi
   ekstra untuk membuat nomor urut.
7. **Data awal** sedikit lebih banyak dari minimum yang diminta, agar semua demo
   Acceptance Criteria terlihat tanpa perlu input manual: 9 bahan baku (bertambah
   *Buah Potong* untuk menu buah), 5 menu (bertambah *Rendang Daging* supaya paket
   Premium masuk akal), 6 pelanggan, 6 pesanan dengan status dan sumber berbeda.
   Paket **Premium** sengaja dibuat bermargin di bawah 30% agar badge peringatan
   langsung terlihat saat demo.
8. **HTMX dimuat dari CDN** (`unpkg.com`) dan setiap form juga punya `method="post"`
   + `action="..."` biasa, sehingga aplikasi **tetap berfungsi tanpa JavaScript** —
   termasuk saat presentasi/demo di jaringan tanpa internet.
9. **Uang** disimpan sebagai `BIGINT` (Rupiah bulat), kuantitas sebagai
   `DOUBLE PRECISION`, margin sebagai `float64` hanya di dalam Go.
10. **Pengurangan stok** hanya terjadi saat transisi status pesanan masuk `produksi`,
    dan dilakukan dalam satu transaksi sehingga tidak terjadi pengurangan ganda.
11. **Tugas produksi** dibuat begitu DP dicatat (bukan menunggu status `produksi`),
    sesuai acceptance criteria: DP 50% yang masuk menandai pesanan sudah pasti
    jalan, sehingga jadwal produksi untuk hari acara langsung tersedia di papan
    dapur. Pembuatannya idempoten (dicek dulu apakah tugas sudah ada).
12. **Setelah pesanan dibuat, halaman dialihkan ke detail pesanan** (`/pesanan/<id>`),
    bukan ke daftar. Owner langsung melihat tenggat DP/pelunasan dan bisa mencatat DP
    tanpa mencari lagi pesanan yang baru dibuat di tabel.
13. **Panel yang dimuat HTMX juga bisa dibuka tanpa JavaScript.** Rute
    `/hpp/hasil`, `/bahan/{id}/harga`, `/menu/{id}/resep`, dan `/paket/{id}/isi`
    mengembalikan *partial* (potongan HTML) bila request berasal dari HTMX, dan
    **halaman penuh** (layout + panel yang sama) bila request biasa — mis. saat
    JavaScript dimatikan atau link dibuka langsung. Data panel dikirim pada key
    `Panel` supaya halaman dan partial memakai sumber data yang sama, dan `id`
    elemen panel (`hpp-hasil`, `riwayat-harga`, `resep-panel`, `paket-isi-panel`)
    tetap muncul tepat satu kali.

---

## 7. Catatan operasional

* Zona waktu container diset `Asia/Jakarta` agar perhitungan H-30 dan H-7 konsisten
  dengan kalender toko.
* `SESSION_SECRET` **wajib** diisi (aplikasi menolak start bila kosong). Untuk
  pemakaian nyata, ganti nilai bawaan di `.env`.
* Password admin di-seed dari `ADMIN_PASSWORD` (default `admin123`) dan di-hash dengan
  bcrypt **saat runtime** — tidak ada hash yang ditulis di dalam kode.
