# Diagram & Dokumentasi Rancangan — SIM Dina Catering

Folder ini berisi materi lampiran untuk laporan tugas akhir Sistem Informasi Manajemen
Dina Catering: diagram proses usulan, Data Flow Diagram, diagram use case, ERD, dan
rancangan antarmuka berupa wireframe prototipe.

Semua materi menggambarkan **sistem yang benar-benar dibangun** di repo ini — nama tabel,
kolom, rute, status, dan aturan mainannya diambil langsung dari
[`db/schema.sql`](../db/schema.sql), [`internal/model/model.go`](../internal/model/model.go),
[`internal/handler/handler.go`](../internal/handler/handler.go), dan
[`internal/service/`](../internal/service/).

---

## 1. Daftar Isi

| Berkas | Isi | Gambar |
|---|---|---|
| [01-Diagram-Proses-Usulan.md](01-Diagram-Proses-Usulan.md) | Diagram alur proses manual (situasi lama), alur proses usulan (sistem baru), perhitungan HPP, pengadaan bahan dan stok, serta pelaporan keuangan. | `01`–`05` |
| [02-DFD-Level-0-dan-Level-1.md](02-DFD-Level-0-dan-Level-1.md) | Context diagram (Level 0) dan DFD Level 1 dengan 6 proses dan 9 penyimpanan data, dilengkapi kamus data. | `06`–`07` |
| [03-Diagram-Use-Case.md](03-Diagram-Use-Case.md) | Use case diagram dengan 4 aktor dan 18 use case, relasi `include`, deskripsi use case, dan matriks aktor. | `08` |
| [04-ERD.md](04-ERD.md) | ERD lengkap 14 tabel beserta kolom, constraint, tabel relasi, daftar index, dan tabel desain basis data. | `09` |
| [05-Rancangan-Antarmuka.md](05-Rancangan-Antarmuka.md) | Peta navigasi, ringkasan layar, dan rancangan tiap layar (field, aksi, aturan mainannya). | — |
| [prototipe/](prototipe/) | 14 wireframe low-fidelity dalam HTML yang dapat dibuka langsung di browser. | — |
| [PNG/](PNG/) | Kumpulan 9 diagram PNG dalam satu folder, siap disisipkan ke laporan. | — |
| [mermaid/](mermaid/) | Sumber Mermaid tiap diagram dalam format `.mmd` untuk diedit atau dirender ulang. | — |
| [viewer.html](viewer.html) | Halaman web yang merender seluruh diagram Mermaid sekaligus, untuk dilihat tanpa tooling tambahan. | — |

---

## 2. Diagram Proses Usulan

Bagian 1 menjelaskan masalah yang dijawab sistem: owner mencatat stok di buku, harga
di Excel, dan jadwal di papan tulis, sehingga tidak tahu berapa yang benar-benar untung.
Diagram alur manual menunjukkan alur lama; diagram alur usulan menunjukkan alur baru dengan
penandaan otomatis: total harga, tenggat DP H-30, tenggat pelunasan H-7, pengurangan stok
otomatis, tugas produksi otomatis, dan badge peringatan margin di bawah 30%.

## 3. Data Flow Diagram

Context diagram menampilkan satu proses tunggal dengan 8 aliran data masuk dan 14 aliran data
keluar. Level 1 menguraikannya menjadi enam proses:

| Proses | Cakupan |
|---|---|
| 1.0 Manajemen Data Master | Pelanggan, bahan baku, menu, resep, paket menu, pengaturan |
| 2.0 Pemrosesan Pesanan | Pesanan baru, harga efektif, tenggat pembayaran, perubahan status |
| 3.0 Pemrosesan Pembelian dan Stok | Pembelian di pasar, stok, riwayat harga, pengurangan stok |
| 4.0 Pemrosesan Pembayaran | DP, pelunasan, sisa tagihan, status pembayaran, jatuh tempo |
| 5.0 Penjadwalan Produksi | Tugas produksi default, PJ, status pengerjaan |
| 6.0 Pelaporan | Laba rugi, rekap pesanan, margin per pesanan, posisi stok |

Sembilan penyimpanan data (D1–D9) beserta neraca aliran datanya dijelaskan pada dokumen.

## 4. Diagram Use Case

Aktor utama adalah **Owner / Admin** — satu-satunya pihak yang masuk ke aplikasi.
Pelanggan, Pasar/Supplier, dan Tim Dapur adalah aktor sekunder yang hanya diwakili oleh
data yang diinput owner (`sumber_pesanan`, `pasar`, dan `pj`). Kolom `users.role` sudah
tersedia tetapi belum dipakai membatasi akses.

## 5. ERD

ERD memuat 14 tabel dengan seluruh kolom beserta tipe dan constraint-nya, 10 relasi
foreign key, 10 index eksplisit, dan 4 constraint `UNIQUE` (plus `pengaturan.kunci` sebagai
primary key).

> **Catatan penting:** ERD lama di `README.md` §6 masih menuliskan relasi
> `users — pesanan`. DDL pada `db/schema.sql` **tidak** memiliki kolom foreign key dari
> `pesanan` ke `users`, sehingga relasi tersebut tidak digambar di ERD folder ini. Entitas
> `users` dan `pengaturan` berdiri sendiri.

## 6. Prototipe Antarmuka

Folder [`prototipe/`](prototipe/) berisi 15 berkas HTML — satu wireframe untuk tiap layar
aplikasi, ditambah halaman indeks:

```
prototipe/
├── index.html daftar layar, peta navigasi, legenda simbol
├── assets/app.css gaya low-fidelity
├── 01-login.html 02-dashboard.html
├── 03-pelanggan.html 04-bahan-baku.html
├── 05-pembelian.html 06-menu-resep.html
├── 07-paket-margin.html 08-kalkulator-hpp.html
├── 09-pesanan.html 10-pesanan-detail.html
├── 11-produksi.html 12-laporan.html
├── 13-biaya.html 14-pengaturan.html
```

**Cara membuka:** klik dua kali `prototipe/index.html`, atau jalankan di VS Code dengan
Live Server. Tidak membutuhkan server, basis data, atau koneksi internet.

Wireframe sengaja dibuat low-fidelity (monokrom) karena yang diuji pada tahap perancangan
adalah tata letak dan alur, bukan tampilan akhir. Setiap wireframe memuat blok
**"Keterangan rancangan"** di bagian bawah berisi daftar field, aksi, dan aturan mainannya,
sehingga wireframe dapat berdiri sendiri sebagai lampiran.

---

## 7. Melihat dan Merender Diagram

### Cara 1 — Viewer HTML (paling mudah)

Buka [`viewer.html`](viewer.html) di browser. Halaman ini memuat seluruh diagram Mermaid
dari dokumen-dokumen di folder ini dan merendernya sekaligus. Memerlukan koneksi internet
untuk memuat Mermaid dari CDN.

### Cara 2 — Pratinjau di VS Code

Buka salah satu berkas `.md` dan tekan `Ctrl+Shift+V`. Blok Mermaid langsung dirender oleh
pratinjau bawaan VS Code — tanpa tambahan apa pun.

### Cara 3 — Render menjadi gambar (untuk disisipkan ke laporan)

Berkas sumber Mermaid tersedia di folder [`mermaid/`](mermaid/) dalam format `.mmd`.
Render menjadi PNG atau SVG memakai `mermaid-cli`:

```bash
npm install @mermaid-js/mermaid-cli
mmdc -i mermaid/09-erd.mmd -o gambar/09-erd.png -b white -s 2
```

Atau pakai skrip yang sudah disiapkan di repo, yang mengekstrak seluruh blok Mermaid dari
berkas `.md` lalu merendernya sekaligus:

```bash
node scripts/render-diagram.mjs
```

Skrip tersebut mencari `mmdc` di `node_modules/.bin` repo. Bila mermaid-cli dipasang di
folder lain, arahkan lewat variabel lingkungan:

```bash
MMDC_PATH=/path/ke/mmdc node scripts/render-diagram.mjs
```

> **Catatan:** pada mesin tanpa Chrome, tambahkan berkas konfigurasi Puppeteer yang
> menunjuk Edge yang sudah terpasang:
> ```json
> { "executablePath": "C:/Program Files (x86)/Microsoft/Edge/Application/msedge.exe" }
> ```
> lalu jalankan `mmdc... -p path/ke/puppeteer.json`.

### Cara 4 — mermaid.live

Buka <https://mermaid.live>, tempel isi berkas `.mmd` dari folder [`mermaid/`](mermaid/),
lalu unduh sebagai PNG atau SVG. Cara ini paling praktis bila tidak ingin memasang
apa pun.

> **Status folder `gambar/`:** render PNG **sudah dijalankan**. Folder ini berisi 9 berkas
> PNG (plus 9 SVG) yang sesuai dengan 9 sumber Mermaid di folder `mermaid/`.
> Setiap gambar sudah disisipkan ke dokumen `.md` terkait, tepat di bawah blok Mermaid
> yang sama. Bila diagram diedit ulang, render ulang dengan Cara 3 di atas.
>
> Nama berkas PNG mengikuti sumbernya:
>
> | Berkas sumber | Nama PNG hasil render |
> |---|---|
> | `mermaid/01-proses-usulan-manual.mmd` | `gambar/01-proses-usulan-manual.png` |
> | `mermaid/02-proses-usulan-sistem-baru.mmd` | `gambar/02-proses-usulan-sistem-baru.png` |
> | `mermaid/03-perhitungan-hpp.mmd` | `gambar/03-perhitungan-hpp.png` |
> | `mermaid/04-pengadaan-stok.mmd` | `gambar/04-pengadaan-stok.png` |
> | `mermaid/05-pelaporan-keuangan.mmd` | `gambar/05-pelaporan-keuangan.png` |
> | `mermaid/06-dfd-level-0.mmd` | `gambar/06-dfd-level-0.png` |
> | `mermaid/07-dfd-level-1.mmd` | `gambar/07-dfd-level-1.png` |
> | `mermaid/08-use-case.mmd` | `gambar/08-use-case.png` |
> | `mermaid/09-erd.mmd` | `gambar/09-erd.png` |

---

## 8. Konsistensi dengan Kode Aplikasi

Semua diagram dalam folder ini diverifikasi terhadap kode yang ada:

| Aspek | Nilai | Sumber |
|---|---|---|
| Jumlah tabel | 14 | `db/schema.sql` |
| Jumlah halaman | 13 | `internal/handler/handler.go` (`buildTemplates`) |
| Jumlah rute HTTP | 49 (25 GET + 24 POST) | `internal/handler/handler.go` (`Router`) |
| Status pesanan | 7 (`baru` … `selesai`, `batal`) | `db/schema.sql` |
| Sumber pesanan | 4 | `db/schema.sql` |
| Status tugas produksi | 3 (`belum`, `jalan`, `selesai`) | `db/schema.sql` |
| Margin minimum | 30% | `internal/model/model.go` |
| Biaya operasional bawaan | 10% | `internal/model/model.go` |
| Index eksplisit | 10 | `db/schema.sql` |
| Relasi foreign key | 10 | `db/schema.sql` |

---

## 9. Struktur Folder

```
Diagram, dll/
├── README.md indeks dan cara pakai
├── viewer.html viewer Mermaid
├── 01-Diagram-Proses-Usulan.md
├── 02-DFD-Level-0-dan-Level-1.md
├── 03-Diagram-Use-Case.md
├── 04-ERD.md
├── 05-Rancangan-Antarmuka.md
├── mermaid/ sumber.mmd per diagram
├── gambar/ hasil render PNG dan SVG
└── prototipe/ wireframe HTML
```
