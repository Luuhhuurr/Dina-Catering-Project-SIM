# Rancangan Antarmuka (Prototipe)

Dokumen ini menjelaskan rancangan antarmuka aplikasi SIM Dina Catering. Prototipe
low-fidelity tersedia di folder [`prototipe/`](prototipe/) dan dapat dibuka langsung di
browser tanpa server.

---

## 1. Prinsip Rancangan

| Aspek | Keputusan | Alasan |
|---|---|---|
| **Layout** | Sidebar tetap di kiri, topbar di atas, konten di tengah | Owner memakai aplikasi sambil berdiri di dapur; navigasi harus selalu terlihat. |
| **Grup menu** | 5 grup: Utama, Pelanggan & Menu, Bahan & Produksi, Penjualan, Analitik & Sistem | Menu dikelompokkan sesuai urutan kerja pemilik, bukan sesuai tabel basis data. |
| **Teknologi** | Server-rendered `html/template` + HTMX 2, tanpa framework frontend | Agar aplikasi tetap berjalan tanpa JavaScript saat presentasi di jaringan tanpa internet. |
| **Ikon** | Sprite SVG internal | Menghindari permintaan ke CDN tambahan; aplikasi tetap utuh saat offline. |
| **Peringatan margin** | Badge warna pada dashboard, daftar paket, dan laporan | Persoalan utama owner: "berapa yang benar-benar untung". |
| **Dark mode** | Mengikuti `prefers-color-scheme` sistem | Aplikasi menyesuaikan tema gelap atau terang secara otomatis, tanpa tombol mode yang perlu dipelajari pengguna. |

---

## 2. Peta Navigasi

```
Login /login
 └─ (berhasil) → Dashboard /

Sidebar
├── Utama
│ └── Dashboard /
├── Pelanggan & Menu
│ ├── Pelanggan /pelanggan
│ └── Menu & Resep /menu
├── Bahan & Produksi
│ ├── Bahan Baku /bahan
│ ├── Pembelian /pembelian
│ └── Produksi /produksi
├── Penjualan
│ ├── Paket & HPP /paket
│ ├── Kalkulator HPP /hpp
│ └── Pesanan /pesanan
└── Analitik & Sistem
 ├── Laporan /laporan
 ├── Biaya /biaya
 └── Pengaturan /pengaturan

Topbar (selalu di atas konten)
├── Breadcrumb: Dina Catering / <judul halaman>
├── Identitas pengguna: inisial + username + role
└── Tombol Keluar → POST /logout
```

**Alur navigasi antar layar yang paling sering dipakai:**

```
Pelanggan → Pesanan (buat pesanan) → Detail Pesanan (catat DP)
 ↓
 Produksi (papan dapur)
 ↓
 Laporan ( laba rugi)
```

---

## 3. Ringkasan Layar

| No | Layar | Route | Judul Halaman | Aksi Utama | Prototipe |
|---|---|---|---|---|---|
| 1 | Login | `GET/POST /login` | Masuk | Isi username dan password, tekan Masuk | [01-login.html](prototipe/01-login.html) |
| 2 | Dashboard | `GET /` | Dashboard | Lihat ringkasan, buka pesanan yang menunggu bayar | [02-dashboard.html](prototipe/02-dashboard.html) |
| 3 | Pelanggan | `GET/POST /pelanggan` | Pelanggan | Tambah, cari, hapus pelanggan | [03-pelanggan.html](prototipe/03-pelanggan.html) |
| 4 | Bahan Baku | `GET/POST /bahan`, `/bahan/{id}/harga` | Bahan Baku & Riwayat Harga | Tambah bahan, ubah harga, hapus bahan | [04-bahan-baku.html](prototipe/04-bahan-baku.html) |
| 5 | Pembelian | `GET/POST /pembelian` | Pembelian Bahan di Pasar | Catat pembelian dari pasar | [05-pembelian.html](prototipe/05-pembelian.html) |
| 6 | Menu & Resep | `GET/POST /menu`, `/menu/{id}/resep` | Menu & Resep | Tambah menu, atur bahan per porsi | [06-menu-resep.html](prototipe/06-menu-resep.html) |
| 7 | Paket & Margin | `GET/POST /paket`, `/paket/{id}/isi` | Paket Menu, HPP & Margin | Tambah paket, isi menu, lihat badge margin | [07-paket-margin.html](prototipe/07-paket-margin.html) |
| 8 | Kalkulator HPP | `GET /hpp`, `/hpp/hasil` | Kalkulator HPP | Pilih paket/menu + tanggal, hitung HPP | [08-kalkulator-hpp.html](prototipe/08-kalkulator-hpp.html) |
| 9 | Pesanan | `GET/POST /pesanan`, `/pesanan/{id}` | Pesanan | Buat pesanan, ubah status, catat pembayaran | [09-pesanan.html](prototipe/09-pesanan.html), [10-pesanan-detail.html](prototipe/10-pesanan-detail.html) |
| 10 | Produksi | `GET /produksi`, `/produksi/tugas/{id}/pj` | Jadwal Produksi | Set PJ, ubah status tugas | [11-produksi.html](prototipe/11-produksi.html) |
| 11 | Laporan | `GET /laporan` | Laporan | Pilih bulan, baca laba rugi dan margin | [12-laporan.html](prototipe/12-laporan.html) |
| 12 | Biaya | `GET/POST /biaya` | Biaya Operasional | Catat biaya operasional | [13-biaya.html](prototipe/13-biaya.html) |
| 13 | Pengaturan | `GET/POST /pengaturan` | Pengaturan | Ubah persentase biaya operasional | [14-pengaturan.html](prototipe/14-pengaturan.html) |

---

## 4. Rancangan Per Layar

### 4.1 Login — `01-login.html`

**Tujuan:** memastikan hanya pemilik usaha yang masuk.

**Komponen:** kartu login di tengah layar (tanpa sidebar), kolom username, kolom password,
tombol Masuk, pesan galat.

**Aturan:**
- Password ditampilkan dengan atribut `type="password"`.
- Bila gagal, pesan "Username atau password salah" muncul di atas tombol, dan field username
 tetap terisi.
- Redirect: berhasil → `303` ke `/`; sudah punya sesi → `303` ke `/`.

---

### 4.2 Dashboard — `02-dashboard.html`

**Tujuan:** memberi gambaran keadaan usaha hari ini dalam satu layar.

**Komponen:**

| Bagian | Isi |
|---|---|
| Banner peringatan | Muncul hanya bila ada paket bermargin di bawah 30%: "N paket punya margin di bawah 30%. Update harga jual atau turunkan HPP." |
| Empat kartu statistik | Jumlah pelanggan · Pesanan menunggu pembayaran · Pemasukan bulan ini · Keuntungan bulan ini |
| Tabel "Menunggu pembayaran" | Pelanggan, tanggal acara, sisa tagihan, label jatuh tempo (DP terlambat / Pelunasan terlambat / DP jatuh tempo / Pelunasan jatuh tempo) |
| Tabel "Agenda acara 7 hari ke depan" | Tanggal, pelanggan, paket, jumlah tamu, total harga |
| Tabel "Pesanan per status" | Lencana status + jumlah |
| Tabel "Pesanan per sumber" | Instagram / Mulut ke Mulut / Hotel / Lainnya + jumlah |
| Tabel "Paket dengan margin di bawah 30%" | Harga jual, HPP per porsi, margin, harga jual minimum |

**Interaksi:** setiap nama pelanggan dan tanggal acara pada tabel tertaut ke detail pesanan.

---

### 4.3 Pelanggan — `03-pelanggan.html`

**Tujuan:** menyimpan data pelanggan agar setiap pesanan terhubung dengan pemesan yang jelas.

**Komponen:** form tambah pelanggan (nama, kontak WhatsApp, alamat, catatan) dan tabel
seluruh pelanggan (nama, kontak, alamat, jumlah pesanan, aksi hapus).

**Aturan:**
- Nama wajib diisi; kontak, alamat, dan catatan boleh dikosongkan.
- Nama pelanggan tidak dibatasi unik — satu orang boleh punya lebih dari satu entri bila
 misalnya berbeda nomor telepon.
- Kolom jumlah pesanan dihitung otomatis dari tabel pesanan, bukan disimpan.
- Penghapusan pelanggan ditolak bila masih dipakai oleh pesanan.

---

### 4.4 Bahan Baku — `04-bahan-baku.html`

**Tujuan:** mengelola master bahan baku beserta riwayat harganya.

**Komponen:** form tambah bahan (nama, satuan, stok awal, harga terakhir), tabel bahan
(nama, satuan, stok, harga terakhir, aksi), dan panel riwayat harga yang terbuka di sisi
kiri atau bawah saat memilih salah satu bahan.

**Aturan:**
- Nama bahan unik; bila nama sudah dipakai, sistem menolak dan menampilkan pesan.
- Harga disimpan dalam Rupiah penuh tanpa desimal.
- Menambah harga baru membuka tanggal berlaku baru dan `harga_terakhir` diperbarui mengikuti
 baris terbaru.

---

### 4.5 Pembelian — `05-pembelian.html`

**Tujuan:** mencatat pembelian bahan di pasar dalam satu langkah, sekaligus memperbarui stok
dan riwayat harga.

**Komponen:** form catat pembelian (bahan, jumlah, harga satuan, pasar, tanggal pembelian) dan
panel riwayat pembelian pada periode terpilih (tanggal, bahan, jumlah, harga, total, pasar).

**Aturan:**
- Jumlah boleh pecahan (mis. 2,5 kg); satuan mengikuti satuan bahan yang dipilih.
- Total transaksi dihitung otomatis dari jumlah x harga satuan.
- Menyimpan pembelian menambah `stok_saat_ini` bahan, memperbarui `harga_terakhir`, dan
 menulis baris riwayat harga baru.
- Bila tanggal pembelian lebih lama daripada riwayat harga terbaru, riwayat harga lama tetap
 dipertahankan sebagai arsip.
- Kolom pasar bebas diisi, hanya untuk jejak belanja.

---

### 4.6 Menu & Resep — `06-menu-resep.html`

**Tujuan:** mendefinisikan menu dan bahannya per porsi.

**Komponen:** form tambah menu (nama, deskripsi), tabel menu (nama, jumlah bahan, aksi),
panel resep per menu (tabel bahan, takaran per porsi, satuan), dan form tambah bahan ke resep.

**Aturan:**
- Resep memakai satuan yang sama dengan satuan bahan.
- Menghapus menu ikut menghapus seluruh resepnya (`ON DELETE CASCADE`).

---

### 4.7 Paket & Margin — `07-paket-margin.html`

**Tujuan:** menetapkan harga jual paket dan memantau marginnya.

**Komponen:** form tambah paket (nama, harga per porsi, harga akhir pekan opsional, tahun
berlaku, deskripsi), tabel paket (nama, harga, HPP per porsi, margin, badge peringatan),
dan panel isi paket.

**Aturan:**
- Kombinasi nama paket dan tahun berlaku harus unik. Paket dengan nama sama dapat dibuat
 ulang untuk tahun berikutnya dengan harga yang berbeda.
- Harga akhir pekan boleh dikosongkan; bila kosong berarti sama dengan harga biasa.
- Margin di bawah 30% diberi badge merah dan disertai harga jual minimum `HPP / 0,7`.

---

### 4.8 Kalkulator HPP — `08-kalkulator-hpp.html`

**Tujuan:** menjawab pertanyaan inti owner, berapa HPP dan berapa untung.

**Komponen:** form (pilih paket atau pilih satu menu, tanggal acuan harga, tombol Hitung HPP),
dan panel hasil.

**Isi panel hasil paket:** rincian per menu (bahan, takaran, harga satuan, biaya), HPP bahan,
HPP per porsi, HPP paket, harga jual efektif, margin, dan harga jual minimum 30%.

**Isi panel hasil menu:** rincian per bahan, HPP bahan, HPP per porsi setelah biaya
operasional.

**Aturan:**
- Persentase biaya operasional yang ditambahkan ditampilkan di bawah form agar owner tahu
 angka yang dipakai (bawaan 10%, dapat diubah di Pengaturan).
- Tanggal acuan menentukan harga bahan yang dipakai, sehingga HPP masa lalu dapat
 dicek ulang.

---

### 4.9 Pesanan — `09-pesanan.html`

**Tujuan:** mencatat pesanan masuk dan memantau pembayarannya.

**Komponen (dua kolom):**

| Kolom kiri | Kolom kanan |
|---|---|
| Form pesanan baru: pelanggan, paket, jumlah tamu, tanggal acara, sumber pesanan, lokasi, catatan, centang "Masak di lokasi", tombol Simpan | Tabel pesanan yang menunggu pembayaran: pelanggan, paket, sisa tagihan, label jatuh tempo |

**Di bawahnya:** tabel seluruh pesanan (nomor, pelanggan, paket, tanggal, total, status, aksi).

**Aturan:**
- Jumlah tamu minimal 1.
- Tanggal acara wajib diisi.
- Total harga, tenggat DP (H-30), dan tenggat pelunasan (H-7) dihitung otomatis dan tidak
 dapat diisi manual.
- Bijak: setelah pesanan tersimpan, owner langsung diarahkan ke halaman detail pesanan agar
 bisa langsung mencatat DP.

---

### 4.10 Detail Pesanan — `10-pesanan-detail.html`

**Tujuan:** tempat kendali satu pesanan dari awal sampai lunas.

**Komponen:**

| Bagian | Isi |
|---|---|
| Header | Nomor pesanan, nama pelanggan, paket, jumlah tamu, tanggal acara, lokasi, sumber pesanan, status |
| Ringkasan pembayaran | Total harga, total dibayar, sisa tagihan, status pembayaran, tanggal jatuh tempo DP dan pelunasan |
| Alur status | Tombol/daftar status berikutnya sesuai alur `baru` → `final_meeting` → `survei_lokasi` → `produksi` → `pengiriman` → `selesai`, serta opsi `batal` |
| Form pembayaran | Tipe (DP / Pelunasan), jumlah, tanggal, metode; jumlah DP terisi otomatis 50% |
| Riwayat pembayaran | Tabel tanggal, tipe, jumlah, metode |
| Daftar tugas produksi | Tugas, PJ, deadline, status |
| Margin pesanan | Total HPP dan margin persen pesanan ini |

**Aturan:**
- Total dibayar dan sisa tagihan dihitung otomatis dari tabel pembayaran.
- Status `produksi` mengurangi stok bahan sesuai resep dikali jumlah tamu dalam satu transaksi.
- Status `batal` tidak mengurangi stok.

---

### 4.11 Produksi — `11-produksi.html`

**Tujuan:** menggantikan papan tulis di dapur.

**Komponen:** papan tugas dikelompokkan per tanggal acara, tiap tugas memuat nama tugas,
nama pelanggan, paket, jumlah tamu, penanggung jawab, deadline, dan status.

**Aksi:** atur PJ (input teks), ubah status tugas (`belum` → `jalan` → `selesai`).

**Aturan:**
- Tugas produksi default dibuat otomatis saat DP dicatat: Persiapan bahan & alat (H-1),
 Masak seluruh menu (H), Plating & cek porsi (H), Pengiriman ke lokasi (H).
- Pembuatan tugas bersifat idempoten sehingga tugas tidak pernah digandakan.

---

### 4.12 Laporan — `12-laporan.html`

**Tujuan:** menjadi sumber kebenaran untung-rugi usaha.

**Komponen:** filter bulan, empat kartu statistik (pemasukan, pengeluaran bahan, biaya
operasional, keuntungan), tabel rekap pesanan per status, tabel rekap pesanan per sumber,
tabel margin per pesanan, tabel peringatan margin, dan tabel posisi stok bahan.

**Aturan:**
- Laba rugi memakai biaya operasional **aktual** yang dicatat, bukan estimasi persentase.
 Karena itu persentase pada kalkulator HPP tidak ditambahkan lagi di sini.
- Tanggal acuan default adalah bulan berjalan.

---

### 4.13 Biaya — `13-biaya.html`

**Tujuan:** mencatat biaya operasional yang tidak termasuk harga bahan, untuk dipakai pada
laporan laba rugi.

**Komponen:** form catat biaya (kategori, deskripsi, nominal, tanggal) dan tabel daftar biaya
pada periode terpilih beserta baris total.

**Aturan:**
- Kategori dipilih dari daftar tetap: gas, listrik, mobilitas, tenaga kerja, perlengkapan,
 lainnya.
- Nominal disimpan dalam Rupiah penuh (BIGINT), tanpa pecahan.
- Biaya yang dicatat di sini **bukan** bagian dari HPP bahan; hubungannya dengan HPP diatur
 lewat persentase pada halaman Pengaturan.
- Total pada kaki tabel adalah jumlah seluruh biaya pada periode yang sedang ditampilkan.

---

### 4.14 Pengaturan — `14-pengaturan.html`

**Tujuan:** menetapkan parameter perhitungan yang dipakai seluruh sistem tanpa mengubah kode.

**Komponen:** form biaya operasional (persentase) dan tabel seluruh nilai pengaturan beserta
waktu perubahan terakhir.

**Aturan:**
- Persentase boleh pecahan (step 0,5), misal 12,5%.
- Nilai bawaan bila belum pernah diatur adalah 10%.
- Perubahan persentase langsung berlaku pada perhitungan HPP, margin, dan harga minimum
 setiap paket.
- Nilai disimpan sebagai pasangan kunci dan nilai pada tabel `pengaturan`, sehingga jumlah
 jenis pengaturan dapat bertambah tanpa mengubah skema tabel.
- Margin minimum 30% bersifat tetap di dalam program dan tidak dapat diubah dari halaman ini.

---

---

## 5. Komponen Antarmuka Umum

| Komponen | Definisi | Perilaku |
|---|---|---|
| `kartu` | Wadah satu blok informasi | Latar putih, tepi tipis, radius 12px. |
| `statistik` | Kartu angka ringkasan | Empat varian warna: biru (data), hijau (uang), amber (perhatian), merah (bahaya). |
| `lencana` | Label status pendek | Varian netral, sukses, peringatan, bahaya. |
| `pesan` | Banner di atas konten | Varian sukses dan error. |
| `tombol` | Tombol aksi | Varian utama (warna brand), `secondary` (netral), `kecil` (ukuran kecil). |
| `field` | Satu baris form | Label di atas, input di bawah. |
| `form-grid` | Kisi form | Dua sampai tiga kolom pada layar lebar. |
| `filter-bar` | Baris filter laporan | Sepadan horizontal. |
| `tabel-bungkus` | Wadah tabel | Tabel dapat digulir mendatar bila lebar layar kecil. |
| `kosong` | Placeholder data kosong | Ikon centang atau kotak kosong + pesan singkat. |
| `pesan toast` | Notifikasi mengambang | Muncul setelah aksi HTMX, hilang otomatis. |
| `papan` / `tugas` | Papan produksi | Kartu tugas per tanggal acara. |

---

## 6. Kebersihan Antarmuka

Seluruh halaman memakai satu layout yang sama, sehingga tidak ada aturan tambahan per halaman:

1. Setiap halaman punya judul di kiri atas beserta keterangan satu kalimat di bawahnya.
2. Form pencarian atau filter selalu berada di atas, tabel di bawahnya.
3. Aksi yang menghapus data memakai konfirmasi (`hx-confirm`).
4. Setiap aksi HTMX punya padanan form biasa (`method` dan `action`), sehingga aplikasi tetap
 berfungsi tanpa JavaScript.
5. Field wajib ditandai dengan atribut `required` pada sisi HTML.

---

## 7. Daftar Berkas Prototipe

| Berkas | Isi |
|---|---|
| `prototipe/index.html` | Daftar seluruh layar, peta navigasi, dan legenda simbol. |
| `prototipe/assets/app.css` | Gaya low-fidelity (monokrom) untuk seluruh wireframe. |
| `prototipe/01-login.html` | Wireframe halaman login. |
| `prototipe/02-dashboard.html` | Wireframe dashboard. |
| `prototipe/03-pelanggan.html` | Wireframe data pelanggan. |
| `prototipe/04-bahan-baku.html` | Wireframe bahan baku dan riwayat harga. |
| `prototipe/05-pembelian.html` | Wireframe pembelian bahan di pasar. |
| `prototipe/06-menu-resep.html` | Wireframe menu dan resep. |
| `prototipe/07-paket-margin.html` | Wireframe paket, HPP, dan margin. |
| `prototipe/08-kalkulator-hpp.html` | Wireframe kalkulator HPP. |
| `prototipe/09-pesanan.html` | Wireframe daftar pesanan dan form pesanan baru. |
| `prototipe/10-pesanan-detail.html` | Wireframe detail pesanan. |
| `prototipe/11-produksi.html` | Wireframe papan produksi. |
| `prototipe/12-laporan.html` | Wireframe laporan keuangan. |
| `prototipe/13-biaya.html` | Wireframe biaya operasional. |
| `prototipe/14-pengaturan.html` | Wireframe pengaturan perhitungan. |

Setiap wireframe memuat blok **"Keterangan rancangan"** di bagian bawah yang merinci field,
aksi, dan aturan mainannya, sehingga wireframe dapat berdiri sendiri sebagai lampiran.
