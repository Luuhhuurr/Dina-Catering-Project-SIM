# Panduan Presentasi — SIM Dina Catering

> **Untuk apa dokumen ini?** Ini naskah presentasi Anda saat mempresentasikan Final Proyek
> di depan dosen, plus lembar contekan dan bank pertanyaan-jawaban. Dokumen ini **bukan**
> dokumentasi produk — untuk itu sudah ada `README.md` dan `docs/DOCKER.md`.
>
> **Aturan pakai:** Bagian 1 adalah naskah linear (baca/ucapkan) dan **sudah dibagi
> ke 4 pembicara**. Bagian 2 adalah **feature tour aplikasi: aplikasi dibuka, lalu tiap
> pembicara menjelaskan layar-layar jadi bagiannya** (satu segmen per orang).
> Bagian 3–5 adalah contekan saat dosen bertanya. Bagian 6–7 adalah netsaibi sebelum
> presentasi.

**Total durasi: ±35 menit** — slide 23 menit + demo langsung 15 menit
(tanya jawab menambah 10–15 menit). **Kalau quota Anda 20 menit**, pakai jalur cepat
di akhir Bagian 2: slide 3–4 cukup 2 menit, slide 8–9 cukup 1 menit, dan demo dipangkas
jadi 7 menit dengan cara melewati semua langkah bertanda ★.

**Jumlah slide: 10. Demo langsung: 4 segmen / 32 langkah (±15 menit). Bank Q&A: 30 pertanyaan.**
**Kelompok: 4 orang.** Pembagian peran ada di Bagian 0, skrip demo lengkap ada di Bagian 2.

---

## Peta Dokumen

| Bagian | Isi | Kapan dipakai |
|---|---|---|
| 0 | Identitas project, ringkasan, tech stack, **pembagian peran 4 orang** | Slide 1 |
| 1 | Naskah presentasi linear 10 slide (ada durasi + nama pemateri) | Presentasi inti |
| 2 | **Feature tour: 13 layar dibagi 4 segmen, satu per pembicara** | Slide 7 |
| 3 | Angka acuan HPP & margin (terverifikasi) | Slide 5 & 6 |
| 4 | Peta kode "kalau ditanya, buka ini" | Saat dosen minta buka kode |
| 5 | Bank Q&A 30 pertanyaan | Tanya jawab |
| 6 | Checklist sebelum presentasi & troubleshooting | H-1 dan H-0 |
| 7 | Lampiran: diagram, wireframe, screenshot | Slide 4 & lampiran laporan |

---

# BAGIAN 0 — IDENTITAS & RINGKASAN

| | |
|---|---|
| **Nama aplikasi** | SIM Dina Catering (Sistem Informasi Manajemen) |
| **Jenis** | Final Proyek kuliah — aplikasi web manajemen usaha katering |
| **Obyek** | Usaha katering **Dina Catering**, Ungaran, Kabupaten Semarang |
| **Repo** | `https://github.com/Luuhhuurr/Dina-Catering-Project-SIM` |
| **Aktor utama** | Owner / Admin — satu-satunya pihak yang masuk ke aplikasi |
| **Bahasa** | Go 1.26 · `html/template` + HTMX 2 · vanilla CSS (814 baris, tulis tangan) |
| **Database** | PostgreSQL 16 · akses via **sqlc (pgx/v5)** · 14 tabel |
| **Migrasi** | golang-migrate v4, file migrasi **di-embed ke dalam binary** |
| **Runtime** | Docker + docker-compose (2 service: `app`, `postgres`) |
| **Dependency** | **Hanya 4**: chi, pgx, golang-migrate, golang.org/x/crypto |
| **Skala kode** | 40 file Go · 6.206 baris Go · 67 query sqlc · 49 rute · 13 halaman · 16 partial |

### Ringkasan 3 kalimat (untuk slide pembuka)

> "SIM Dina Catering adalah aplikasi web untuk mengelola usaha katering dari bahan baku sampai
> laporan laba rugi. Masalah utamanya: owner mencatat stok di buku, harga di Excel, jadwal
> dapur di papan tulis — sehingga harga bahan naik, harga paket tetap, dan dia tidak tahu
> berapa yang benar-benar untung. Solusinya: kalkulator HPP berbasis resep, di mana begitu
> harga bahan berubah, HPP dan margin semua paket langsung dihitung ulang, dan paket yang
> turun di bawah margin 30% otomatis diberi badge peringatan."

### ⚠️ Satu koreksi yang WAJIB Anda akimati sebelum presentasi

Project ini **bukan** aplikasi e-commerce/restaurant untuk pelanggan. Pelanggan di sini adalah
**data master yang di-input owner**, bukan pengguna aplikasi. Tidak ada keranjang, checkout,
katalog publik, atau pendaftaran akun. Yang ada adalah **13 layar operasional** untuk owner.

Kalau Anda keliru dan berkata "pelanggan browse menu lalu add to cart", dosen akan langsung
menyadari sistem ini tidak seperti yang Anda gambarkan — dan seluruh presentasi kehilangan
kredibilitas. Gunakan istilah: **owner / admin** untuk yang memakai aplikasi, **pelanggan**
untuk data pesanan yang dicatat owner.

---

## Pembagian peran 4 orang

> **Prinsip:** setiap orang memegang satu bidang penuh — slide-nya *dan* demo-nya.
> Kalau ada yang ditanya dosen, orang itulah yang menjawab. Jadi tidak perlu hafal semua
> naskah; hafal **bidangmu**, dan kalau ditanya bagian lain, teruskan ke teammate.

| Pembicara | Nama | Slide | Isi Presentasi | Segmen Demo | Durasi |
|---|---|---|---|---|---|
| **P1** | *(isi nama)* | 1, 2 | Pembuka, masalah, solusi & ruang lingkup 7 modul | **Segmen 1** — Login → Dashboard → Pelanggan → Pengaturan | 4 mnt + 3 mnt |
| **P2** | *(isi nama)* | 3, 4 | Arsitektur berlapis, alasan teknologi, model data | **Segmen 2** — Bahan Baku → Menu & Resep → Paket & HPP → Kalkulator HPP | 7 mnt + 5 mnt |
| **P3** | *(isi nama)* | 5, 6 | Fitur inti kalkulator HPP, alur operasional pesanan | **Segmen 3** — Pembelian → Pesanan → Detail pesanan → Produksi | 7 mnt + 4 mnt |
| **P4** | *(isi nama)* | 8, 9, 10 | Pengujian, keputusan desain & batasan, penutup | **Segmen 4** — Laporan → Biaya Operasional → bukti pengujian | 5 mnt + 3 mnt |

> Durasi = durasi slide + durasi segmen demo-nya. Total: 23 menit slide + 15 menit demo.

**Pembagian peran backup (kalau ada yang berhalangan):**

| Yang berhalangan | Diganti oleh | Konsekuensi |
|---|---|---|
| P1 (pembuka) | P2 | P2 memegang slide 1–2 **dan** segmen 1 (dashboard). Segmen 2 jadi lebih padat. |
| P2 (arsitektur) | P4 | P4 memegang slide 3–4 **dan** segmen 2 (HPP). Segmen 4 jadi lebih padat. |
| P3 (operasional) | P4 | P4 memegang slide 5–6 **dan** segmen 3. Segmen 4 jadi lebih padat. |
| P4 (penutup) | P1 | P1 memegang slide 8–10 **dan** segmen 4, lalu menutup sendiri presentasi. |

> **Aturan serah terima antar pembicara:** bagian terakhir tiap segmen selalu ditutup dengan
> kalimat eksplisit "sekarang saya serahkan ke **[nama]**" supaya tidak ada jeda bingung.
> Kalimat siap pakai ada di akhir tiap segmen (Bagian 2).

---

# BAGIAN 1 — NASKAH PRESENTASI

> **Cara pakai:** Setiap slide punya **Durasi**, **Yang diucapkan** (siap dibacakan, pakai
> bahasa sendiri kalau perlu), **Yang dibuka** (file + nomor baris, untuk kasus dosen minta
> "tunjukkan kodenya"), dan **Penekan** (satu kalimat penutup yang mengena).

---

## SLIDE 1 — Pembuka & masalah (2 menit)

> **Pemateri: P1** · Durasi 2 menit · Bagian dari naskah P1

**Yang diucapkan:**

> "Selamat pagi. Perkenalkan, saya [nama]. Project saya adalah Sistem Informasi Manajemen
> untuk usaha katering Dina Catering di Ungaran — ini Final Proyek saya.
>
> Saya mau mulai dari masalahnya. Hj. Dina, pemilik usaha ini, bekerja seperti begini: stok
> bahan dicatat di buku tulis, harga dihitung di Excel, jadwal dapur ditulis di papan tulis
> di dapur. Akibatnya, ketika harga bahan naik di pasar — dan di Indonesia harga bahan
> katering naik hampir setiap bulan — harga paket yang sudah terlanjur dipasang ke
> pelanggan tidak ikut naik. Dan karena tidak ada hitungan HPP yang benar, pemilik tidak tahu
> paket mana yang masih untung dan paket mana yang diam-diam subsidari pelanggan.
>
> Pertanyaan yang jadi inti project saya: **bagaimana supaya begitu harga bahan berubah,
> sistem langsung bisa menghitung ulang berapa HPP dan berapa margin setiap paket — dan
> memberi tahu kalau ada paket yang sudah tidak layak dijual?**"

**Penekan:** "Ini bukan aplikasi kasir. Ini aplikasi yang menjawab: paket mana yang
sekarang embarrassing untuk dijual."

**Yang dibuka (kalau diminta):** `Diagram,dll/PNG/01-proses-usulan-manual.png` (alur manual)
dibandingkan `02-proses-usulan-sistem-baru.png` (alur usulan).

---

## SLIDE 2 — Solusi & ruang lingkup (2 menit)

> **Pemateri: P1** · Durasi 2 menit · Bagian dari naskah P1

**Yang diucapkan:**

> "Dari masalah itu saya bangun satu aplikasi dengan tujuh modul.
>
> Yang paling penting, modul A, adalah **Kalkulator HPP dan Resep**. Setiap menu punya resep:
> bahan apa saja dan takaran berapa per porsi. Setiap bahan punya riwayat harga bertanggal,
> jadi HPP historis bisa dihitung ulang. Dari situ HPP per porsi, margin, dan harga jual
> minimum untuk margin 30% semuanya turun otomatis.
>
> Modul lain: pelanggan; pesanan dengan alur status dan tenggat DP H-30 dan pelunasan H-7;
> pembelian dan stok; jadwal produksi menggantikan papan tulis di dapur; laporan; dan login.
>
> Saya juga perlu menyebut batasannya di sini, karena saya ingin presentasi ini jujur.
> **Sistem ini belum multi-user.** Kolom `role` sudah saya siapkan di tabel `users`, tapi
> saat ini semua fitur hanya dibatasi 'sudah login' atau 'belum login'. Ini keputusan yang
> saya ambil dengan sadar: saya memilih menyelesaikan satu aplikasi dengan utuh, bukan
> setengah jadi dengan lebih banyak modul — dan saya jelaskan di slide 9."

**Ringkas ruang lingkup (7 modul, `README.md` §3):**

| | Modul | Inti |
|---|---|---|
| A | **Kalkulator HPP & Resep** | Resep per porsi, riwayat harga, HPP, margin, badge <30% |
| B | **Pelanggan** | Nama, kontak WA, alamat, catatan, jumlah pesanan |
| C | **Pesanan & Pembayaran** | Alur status 7 tahap, DP H-30 / lunas H-7, harga akhir pekan |
| D | **Pembelian & Stok** | Pembelian di pasar → stok naik + riwayat harga + HPP ikut berubah |
| E | **Jadwal Produksi** | 4 tugas otomatis per pesanan, PJ, status; stok berkurang saat status `produksi` |
| F | **Laporan** | Laba rugi bulanan, rekap per status/sumber, margin per pesanan, posisi stok |
| G | **Login** | Cookie session bertanda tangan HMAC-SHA256, password bcrypt |

**Penekan:** "Tujuh modul, tapi semuanya bermuara pada satu angka: HPP."

---

## SLIDE 3 — Arsitektur & teknologi (4 menit)

> **Pemateri: P2** · Durasi 4 menit · Bagian dari naskah P2

> Ini slide teknis. Kalau dosen tidak teknis, bacakan **Bagian A** saja dan lewati ke Bagian B.

### Bagian A — Arsitektur (baca)

> "Arsitektur saya berlapis, dan aturannya saya tegakkan ketat.
>
> **Aturan pertama: tidak ada SQL di dalam handler.** Handler HTTP tidak boleh menulis query
> sama sekali. Dia hanya menerima request, memanggil service, lalu merender template.
>
> **Aturan kedua: tidak ada ORM.** Saya pakai **sqlc** — query saya tulis di
> `db/queries.sql` dengan anotasi, lalu sqlc meng-*generate* kode Go-nya ke `internal/store`.
> Jadi saya tetap menulis SQL — pakai sintaks SQL asli, bukan DSL — tapi hasilnya
> type-safe: kalau saya ubah skema, sqlc dan compiler langsung memberi tahu.
>
> **Aturan ketiga: semua perhitungan bisnis ada di `internal/service`** dan semuanya punya
> unit test. Handler cumaorkestration, service yang mengambil keputusan.
>
> Alurnya: request masuk ke router chi, melewati middleware (logging, recover, auth), masuk
> handler, handler panggil service, service panggil store (sqlc), store bicara dengan PostgreSQL.
> Result-nya balik ke template `html/template`, dirender jadi HTML."

### Bagian B —-Alasan pemilihan teknologi (baca)

> "Kenapa begini? Karena saya memang embracing kesederhanaan.
>
> Total dependency saya cuma **empat**: chi, pgx, golang-migrate, dan x/crypto. Nggak ada
> ORM, nggak ada framework frontend, nggak ada template engine, nggak ada build step untuk
> JavaScript. CSS-nya 814 baris vanilla, saya tulis sendiri.
>
> Dependensi sedikit itu deliberate — kalau ada yang rusak atau perlu di-upgrade, saya
> tahu persis apa yang terpengaruh.
>
> Untuk frontend saya pakai **HTMX**. Alasannya: aplikasi ini memang server-side rendering
> dari awal, dan HTMX memberi saya partial refresh tanpa harus membangun SPA dan tanpa
> framework JavaScript. Dan ada satu keputusan yang saya sengaja ambil: **HTMX saya muat
> dari CDN, tapi setiap form tetap punya `method` dan `action` biasa** — artinya kalau
> JavaScript-nya gagal dimuat, aplikasinya **tetap berfungsi penuh**. Ini penting karena
> saat presentasi, kalau wifi kampus mati, demo saya tidak ikut mati."

**Yang dibuka (kalau diminta):**

| Yang ditanyakan | Buka file |
|---|---|
| "Tunjukkan bootstrap aplikasi" | `cmd/server/main.go:32-87` |
| "Tunjukkan rutenya" | `internal/handler/handler.go:58-163` |
| "Kenapa tanpa ORM?" | `go.mod`, `sqlc.yaml`, `internal/store/querier.go:12-125` |
| "Tunjukkan session" | `internal/middleware/session.go:42-134` |

**Penekan:** "Arsitektur saya dijaga oleh tiga aturan yang bisa diaudit: tidak ada SQL di
handler, tidak ada ORM, semua perhitungan ada di service."

---

## SLIDE 4 — Model data (3 menit)

> **Pemateri: P2** · Durasi 3 menit · Bagian dari naskah P2

**Yang diucapkan:**

> "Basis datanya 14 tabel, PostgreSQL 16. Saya tidak perlu tampilkan semua kolomnya, tapi
> saya mau jelaskan **empat keputusan desain** yang membentuk hampir semua logika sistem.
>
> **Pertama, `harga_bahan_history`.** Harga bahan tidak disimpan di satu kolom saja, tapi
> sebagai riwayat bertanggal berlaku. Ini yang membuat HPP historis bisa dihitung ulang —
> dulu saya bisa tanya 'berapa HPP paket Premium bulan Mei?', dan jawabannya tetap sama
> walaupun harga daging hari ini sudah naik.
>
> **Kedua, `paket_menu.harga_per_porsi_weekend` itu nullable.** Kalau NULL, berarti pakai
> harga biasa. Saya choosing kolom nullable daripada tabel terpisah, karena tidak setiap
> paket perlu tarif akhir pekan, dan menambah tabel berarti menambah join tanpa manfaat.
>
> **Ketiga, `paket_menu.tahun_berlaku`.** Harga paket itu historis juga. Baris (Premium,
> 2026) dan (Premium, 2027) bisa hidup berdampingan. Satu pesanan terkunci ke satu baris
> harga — jadi kalau owner menaikkan harga di tengah jalan, pesanan lama tidak ikut berubah.
>
> **Keempat, tabel `pengaturan` yang berisi key-value.** Persentase biaya operasional
> disimpan di database, bukan environment variable, supaya bisa diubah dari UI oleh owner
> yang bukan pakai terminal."

**Tabel inti (ringkas 14 tabel → 9 yang paling penting):**

| Tabel | Peran |
|---|---|
| `users` | Login. `role` sudah ada, belum dipakai gating |
| `pengaturan` | Key-value config (`biaya_operasional_persen` = 10) |
| `pelanggan` | Data master pelanggan (6 data awal) |
| `bahan_baku` | Master bahan + stok saat ini + harga terakhir (9 bahan) |
| `harga_bahan_history` | **Riwayat harga bertanggal berlaku** |
| `pembelian_bahan` | Pembelian di pasar (6 transaksi awal) |
| `menu` | Master menu (5 menu) |
| `resep_item` | **Resep**: bahan + takaran per porsi (12 baris) |
| `paket_menu` | **Harga paket per tahun**, incl. harga weekend nullable (3 paket) |
| `paket_menu_item` | Isi menu dalam paket (many-to-many) |
| `pesanan` | Tabel pusat: pelanggan, paket, tamu, tanggal, total, status, tenggat DP/lunas |
| `pembayaran` | DP & pelunasan |
| `tugas_produksi` | Papan dapur: tugas, PJ, deadline, status |
| `biaya_operasional` | Biaya aktual (gas, listrik, mobilitas, perlengkapan) |

**Yang dibuka:** `Diagram,dll/PNG/09-erd.png` (ERD lengkap) atau `Diagram,dll/04-ERD.md`
(ERD + kamus data + daftar index).

**Penekan:** "Empat keputusan di atas yang membuat seluruh fitur bisa jadi — sisanya
justru konsekuensi."

---

## SLIDE 5 — Fitur inti: kalkulator HPP (4 menit)

> **Pemateri: P3** · Durasi 4 menit · Bagian dari naskah P3

> Ini slide paling penting. Jangan buru-buru.

**Yang diucapkan:**

> "Sekarang saya masuk ke inti sistem: kalkulator HPP.
>
> Rumusnya begini. **HPP bahan** per porsi adalah jumlah dari semua bahan di resep, masing-masing
> dikali harga bahan yang berlaku pada tanggal itu. Jadi kalau harga daging naik, HPP
> otomatis naik — karena harga yang dipakai bukan harga hardcode, tapi harga dari riwayat
> yang berlaku pada tanggal tersebut.
>
> **HPP per porsi** = HPP bahan × (1 + persentase biaya operasional). Default-nya 10%, dan
> owner bisa ubah dari halaman Pengaturan.
>
> **Margin** = (harga jual − HPP) / harga jual. Dan **harga jual minimum untuk margin 30%**
> = HPP dibagi 0,7.
>
> Perhatikan urutannya: biaya operasional ditambahkan **proporsional ke HPP bahan**, baru
> margin dihitung dari HPP total. Kalau urutannya dibalik, rumus HPP/0,7 jadi tidak
> konsisten. Ini detail kecil tapi penting, dan saya sudah dokumentasikan alasannya."

### Angka nyata dari data awal (boleh dibuka saat menjelaskan)

| Paket | Komponen HPP bahan/porsi | HPP bahan | +10% biaya ops | **HPP/por** | Harga jual | **Margin** | Min. harga 30% |
|---|---|---|---|---|---|---|---|
| **Hemat** | Nasi 2.670 + Sayur 2.110 + Ayam 10.250 | 15.030 | +1.503 | **16.533** | 30.000 | **44.9%** | 23.619 |
| **Standar** | Nasi 2.670 + Sayur 2.110 + Ayam 10.250 + Buah 2.250 | 17.280 | +1.728 | **19.008** | 35.000 | **45.7%** | 27.155 |
| **Premium** | Nasi 2.670 + Rendang 20.970 + Sayur 2.110 + Buah 2.250 | 28.000 | +2.800 | **30.800** | 40.000 | **23.0%** ⚠️ | 44.000 |

> **Contoh penjelasan Paket Hemat:** Nasi Putih = beras 0,150 kg × Rp16.000 = Rp2.400, plus
> minyak goreng 0,015 liter × Rp18.000 = Rp270 → Rp2.670. Ayam Bakar = ayam 0,250 kg ×
> Rp38.000 = Rp9.500, bumbu racik 0,040 kg × Rp12.000 = Rp480, minyak Rp270 → Rp10.250.
> Sayur Sop = Rp1.600 + Rp240 + Rp270 → Rp2.110. Total bahan Rp15.030, tambah 10% biaya
> operasional jadi **Rp16.533 per porsi**. Harga jual Rp30.000, jadi margin
> (30.000 − 16.533) / 30.000 = **44.9%**.

> **Poin yang harus Anda sampaikan eksplisit:** Paket **Premium** sengaja saya buat margin-nya
> di bawah 30% — hasilnya 23%. Ini **bukan bug**, ini fitur: begitu harga daging naik lagi
> atau harga paket turun, badge peringatan langsung muncul di dashboard, di daftar paket,
> dan di laporan. Angka Premium ini sudah saya jadikan unit test, jadi regression-nya ketahuan
> kalau ada perubahan kode.

**Yang dibuka (kalau diminta):** `internal/service/hpp.go:42-79` — empat fungsi inti:
`BiayaKomponen`, `TotalHPPPerPorsi`, `MarginPaket`, `HargaJualMinimum30Persen`.
Testnya: `internal/service/hpp_test.go:75-107`.

**Penekan:** "Satu angka mengubah segalanya. Harga bahan berubah → HPP berubah → margin
berubah → sistem memberi tahu paket mana yang tidak lagi layak."

---

## SLIDE 6 — Alur operasional (3 menit)

> **Pemateri: P3** · Durasi 3 menit · Bagian dari naskah P3

> Ini slide yang menjelaskan *sistem hidup*, bukan statis.

### Alur 1 — Pesanan sampai lunas

> "Alur utama: owner menerima order lewat Instagram, WhatsApp, atau hotel, lalu mencatatnya.
>
> Waktu pesanan dibuat, sistem menghitung sendiri tiga hal: **total harga** — pakai harga
> akhir pekan kalau acaranya Sabtu atau Minggu; **tenggat DP** = H-30; dan **tenggat pelunasan**
> = H-7. Owner tidak perlu menghitung apa pun.
>
> Status pesanan punya alur: baru → final meeting → survei lokasi → produksi → pengiriman →
> selesai, atau batal.
>
> Dan ini bagian yang paling ingin saya tunjukkan: **saat DP dicatat, sistem otomatis membuat
> 4 tugas produksi** untuk hari acara — persiapan bahan H-1, masak H, plating H, pengiriman H.
> Alasannya bisnis: begitu DP 50% masuk, pesanan itu sudah pasti jalan, jadi jadwal dapur
> harus langsung tersedia. Pembuatannya idempoten, jadi tugas tidak akan pernah dobel."

### Alur 2 — Pembelian bahan (satu transaksi, tiga efek)

> "Kalau owner mencatat pembelian di pasar, dalam **satu transaksi** terjadi tiga hal:
> stok bahan bertambah, harga terakhir diperbarui, dan satu baris baru masuk ke riwayat harga.
> Efek ketiga itu yang penting — artinya HPP semua paket ikut terhitung ulang, dan tidak ada
> satu pun langkah yang bisa terlewat."

### Alur 3 — Produksi (stok berkurang otomatis)

> "Saat status pesanan diubah menjadi `produksi`, stok bahan otomatis berkurang sebesar
> resep × jumlah tamu, dalam satu transaksi database — jadi tidak mungkin berkurang dobel.
> Saya memilih pengurangan terjadi di titik ini, bukan saat pesanan dibuat, karena kalau
> pesanan dibatalkan setelah DP masuk, stoknya harus utuh."

**Penekan:** "Tiga alur ini saling berantai: mencatat pesanan, mencatat pembelian, lalu
memasuki produksi — dan tiap langkah otomatis memperbarui angka di langkah berikutnya."

**Yang dibuka (kalau diminta):**

| Yang ditanyakan | Buka file |
|---|---|
| "Di mana pesanan dibuat?" | `internal/service/pesanan.go:27-77` |
| "Di mana DP memicu tugas produksi?" | `internal/service/pembayaran.go:115-143` |
| "Di mana stok berkurang?" | `internal/service/pesanan.go:99-106` → `internal/service/stok.go:14-39` |
| "Di mana tugas default dibuat?" | `internal/service/pesanan.go:113-147` |
| "Di mana laba rugi?" | `internal/service/laporan.go:23-57` |

---

## SLIDE 7 — Demo langsung: feature tour aplikasi (15 menit)

> **Pemateri: P1 membuka, lalu P2 → P3 → P4 masing-masing segmen.**
> Naskah lengkap tiap segmen ada di **Bagian 2**.

> Lihat **Bagian 2** untuk langkah persis. Yang diucapkan saat demo cukup sedikit —
> biarkan layarnya yang menjelaskan. Yang penting: **tiap orang menjelaskan layar-layar
> yang jadi bagiannya**, jadi dosen tahu aplikasi ini punya berapa layar dan apa saja.

**Kalimat pembuka sebelum demo (dibaca P1):**

> "Saya akan tunjukkan sistemnya bekerja. Fokusnya satu: **ubah harga satu bahan, lihat
> HPP dan margin semua paket ikut berubah.** Itu inti dari seluruh aplikasi."
>
> "Aplikasi ini punya **13 halaman**. Saya mulai dari dashboard dan data master, lalu
> **[nama P2]** lanjut ke bagian HPP, **[nama P3]** ke operasional pesanan, dan
> **[nama P4]** menutup dengan laporan."

**Urutan segmen (jangan diacak, ujung ke ujung) — feature tour, sesuai urutan menu:**

`P1 · login → dashboard → pelanggan → pengaturan` → `P2 · bahan baku → menu & resep → paket & HPP → kalkulator HPP`
→ `P3 · pembelian → pesanan → detail pesanan → produksi` → `P4 · laporan → biaya operasional → bukti tes`

---

## SLIDE 8 — Pengujian & kualitas (2 menit)

> **Pemateri: P4** · Durasi 2 menit · Bagian dari naskah P4

**Yang diucapkan:**

> "Untuk soal pengujian, saya punya tiga lapis.
>
> **Lapis pertama, unit test** — 13 fungsi di `internal/service` yang menguji HPP, margin,
> harga jual minimum, jadwal DP dan pelunasan, dan status pembayaran. Plus 2 test di
> `internal/handler` yang memastikan semua 13 template halaman bisa di-parse dan
> dirender. Jadi 15 fungsi test. Saya sengaja menaruh angka HPP Paket Premium sebagai
> test case — kalau logika HPP berubah, test ini langsung gagal.
>
> **Lapis kedua, analisis statis** — `go vet` dan `go build`.
>
> **Lapis ketiga, smoke test end-to-end** — `scripts/smoke-test.sh`, sebuah skrip bash yang
> menguji **18 hal** secara berurutan lewat HTTP: redirect tanpa session, login, dashboard,
> badge margin, kalkulator HPP termasuk fallback tanpa JavaScript, perubahan harga bahan
> sampai HPP berubah, pembuatan pesanan dengan cek tenggat H-30 dan H-7, pencatatan DP,
> papan produksi, pembelian sampai stok dan riwayat harga, laporan laba rugi dengan
> rekonsiliasi, sampai logout.
>
> Saya juga tahu apa yang **belum** diuji: belum ada unit test untuk lapisan handler HTTP
> dan middleware session. Untuk sekarang lapisan itu tertutup oleh smoke test."

**Yang dibuka:**
```
go vet ./...
go build ./...
go test ./...
```
↳ Output: `ok` di **2 paket** — `internal/service` dan `internal/handler`
(sisanya `[no test files]`).

**Penekan:** "HPP punya unit test, dan seluruh alur punya smoke test."

---

## SLIDE 9 — Keputusan desain & batasan (2 menit)

> **Pemateri: P4** · Durasi 2 menit · Bagian dari naskah P4

**Yang diucapkan:**

> "Saya mau tutup dengan beberapa keputusan desain yang perlu saya belakan, dan beberapa
> batasan yang masih ada."

### Enam keputusan desain yang paling layak dijelaskan

| # | Keputusan | Alasan |
|---|---|---|
| 1 | Migrasi **di-embed ke binary** (`//go:embed *.sql`) | Tidak butuh executable `migrate` terpisah di image Docker. Migrasi ikut versi binary — selalu cocok. |
| 2 | Harga akhir pekan = **kolom nullable**, bukan tabel terpisah | Tidak semua paket butuh tarif akhir pekan. Kolom nullable menghindari join yang tidak perlu. |
| 3 | Harga paket per **tahun berlaku** (`tahun_berlaku`) | Harga juga historis. Pesanan lama terkunci ke baris harga lamanya, jadi kenaikan harga tidak mengubah pesanan yang sudah jalan. |
| 4 | **Session dibuat sendiri** — cookie HMAC-SHA256 | Tidak butuh tabel session, tidak butuh Redis, tidak menambah dependency. Cookie `HttpOnly` + `SameSite=Lax`, TTL 12 jam, perbandingan `hmac.Equal` (constant-time). |
| 5 | **HTMX dari CDN, tapi form tetap punya `action` biasa** | Aplikasi tetap berfungsi penuh tanpa JavaScript. Aman untuk demo di jaringan tanpa internet. |
| 6 | **Stok hanya berkurang saat status `produksi`** | Kalau pesanan dibatalkan setelah DP masuk, stok harus utuh. Pengurangan dalam satu transaksi, jadi tidak dobel. |

### Batasan yang saya akui sendiri

- **Belum multi-user.** Kolom `role` ada, tapi belum jadi middleware pembatasan akses.
- **Semua actor domain lain tanpa akun.** Pelanggan, supplier, dan tim dapur diwakili
  oleh data (`sumber_pesanan`, `pasar`, `pj`) — sesuai kenyataan usaha katering.
- **Belum ada notifikasi otomatis.** Pengingat tenggat tampil di dashboard, tapi belum
  dikirim via WhatsApp/email.
- **Nilai uang pakai `BIGINT` (Rupiah bulat), kuantitas `DOUBLE PRECISION`.** Ini
  konsekuensi pilihan "tidak pakai decimal" — pembulatan dilakukan per komponen resep
  dengan pembulatan ke rupiah terdekat.
- **Cakupan pengujian:** logika bisnis dan template teruji; lapisan HTTP dan middleware
  belum punya unit test sendiri.
- **Belum ada migrasi incremental.** Hanya `0001_init` — wajar untuk tahap awal.

**Penekan:** "Yang belum dikerjakan saya sebut sendiri, bukan saya tunggu ditanya."

---

## SLIDE 10 — Penutup (1 menit)

> **Pemateri: P4** · Durasi 1 menit · Bagian dari naskah P4

**Yang diucapkan:**

> "Jadi ringkasannya: dari satu masalah nyata — owner tidak tahu paketnya untung atau tidak —
> saya bangun satu aplikasi yang bikin jawaban itu otomatis muncul, lengkap dengan peringatan
> kalau ada paket yang perlu harga dinaikkan.
>
> Yang membuat ini mudah dipelajari dan dirawat: saya menulis SQL sungguhan tanpa ORM,
> semua perhitungan bisnis ada di satu layer yang punya test, dan aplikasinya tetap
> jalan meski JavaScript-nya gagal dimuat.
>
> Untuk tahap berikutnya, tiga hal yang paling masuk akal: mengaktifkan pembatasan akses
> berbasis role, menambah notifikasi WhatsApp untuk tenggat DP dan pelunasan, dan memindahkan
> perhitungan ke struktur historis yang lebih kaya untuk analisis harga.
>
> Terima kasih. Silakan kalau ada pertanyaan."

**Penekan (kalau ditanya "apa intinya?"):** "Satu masalah nyata, satu angka yang jadi
jawaban, dan satu sistem yang menghitung ulang otomatis."

---

# BAGIAN 2 — SKENARIO DEMO: FEATURE TOUR 4 PEMBICARA (±15 menit)

> **Cara demo:** aplikasinya dibuka di layar, lalu **tiap orang menjelaskan layar-layar
> yang jadi bagiannya**. Setiap orang pegang mouse selama segmennya, lalu serah terima ke
> orang berikutnya. Urutannya mengikuti urutan menu di sidebar, jadi navigasinya mulus
> dan tidak ada lompat-lompat.
>
> Prinsipnya: **pembicara demo modul yang sama dengan yang ia jelaskan di slide.**
> P1 → masalah & data master · P2 → HPP (inti) · P3 → operasional · P4 → laporan & mutu.
>
> Demo **mengubah data** (menambah pesanan, pembelian, riwayat harga). Itu memang
> bagian dari ceritanya — dan kalau mau presentasi ulang dengan kondisi bersih:
> `docker compose down -v` lalu `docker compose up --build`.

## Persiapan (kerjakan H-1, bukan H-0)

```bash
cd dina-catering
docker compose down -v             # LANGKAH WAJIB: volume lama dibuang
docker compose up --build          # build + migrate + seed otomatis
```

> **Kenapa `down -v` wajib?** Seed hanya jalan kalau tabel `users` masih kosong. Kalau
> volumenya pernah dipakai untuk latihan (menambah pesanan, mengubah harga bahan), data
> lamanya **tidak hilang** dan angkanya jadi tidak lagi cocok dengan slide. Perintah di
> atas membuang volume itu dan membangun ulang dari nol — database jadi bersih dan
> semua angka di dokumen ini **pasti cocok**. Butuh ±2 menit.

- Buka **http://localhost:8080**
- Login: **`admin`** / **`admin123`**
- Pastikan Docker Desktop sudah running sebelum presentasi.
- **Cek cepat:** di `/bahan`, harga terakhir **Ayam Potong harus 38.000** (bukan 12.345)
  dan **Daging Sapi 135.000**. Kalau tidak, berarti volume-nya belum bersih — ulangi
  `down -v` lalu `up --build`.

## ⚠️ Dua hal yang angkanya BERGANTUNG pada tanggal demo

Ini yang paling mudah membuat salah kutip di depan dosen. Dua hal di bawah ini sudah
**diverifikasi langsung di aplikasi yang berjalan** — please believe ini.

**1. Margin paket di `/paket` dan dashboard memakai harga jual HARI INI.**
Kalau hari ini hari kerja (Senin–Jumat), harga jual yang dipakai Rp40.000; kalau hari ini
Sabtu/Minggu, harga jual akhir pekan Rp43.000 dipakai. Dampaknya:

| Yang tampil | Hari kerja (Sen–Jum) | Akhir pekan (Sab–Min) |
|---|---|---|
| Margin Paket **Hemat** | 44.9% | 48.3% |
| Margin Paket **Standar** | 45.7% | 50.0% |
| Margin Paket **Premium** | **23.0%** ⚠️ | **28.4%** ⚠️ |
| Harga jual minimum Premium | 44.000 | 44.000 |

> **Aturan main:** di dashboard (Segmen 1) dan di tabel paket (Segmen 2), **jangan sebut
> angka margin secara spesifik** — sebut saja "di bawah 30 persen". Itu selalu benar.
> Angka spesifiknya sudah ada di slide 5, dan bisa dibacakan di sana.

**2. Di Kalkulator HPP, "Tanggal acuan" menentukan angka yang tampil.**
Default-nya adalah hari ini. Supaya angkanya pasti dan bisa dihafal, **selalu isi tanggal
acuan dengan HARI KERJA (Senin–Jumat)**:

| Kalau hari ini | Isi tanggal acuan dengan | Alasan |
|---|---|---|
| Senin–Jumat | **hari ini** | Hari kerja, harga jual Rp40.000 |
| Sabtu | **Senin minggu depan** | Tetap hari kerja, dan masih ≥ hari ini |
| Minggu | **Senin besok** | Sama seperti di atas |

> Kenapa harus hari kerja? Karena hanya di hari kerja harga jual Paket Premium = Rp40.000,
> sehingga marginnya **23.0%** dan minimum harganya **44.000** — angka yang sama persis
> dengan yang tertulis di slide 5 dan Bagian 3. Kalau tanggal acuan jatuh akhir pekan,
> angkanya jadi 28.4% dan dosen akan berpikir "kok beda dari slide?".

## Aturan main bersama (baca sebelum hari-H)

| Aturan | Detail |
|---|---|
| **Satu browser, satu sesi login** | Jangan logout di tengah demo. Session bertahan 12 jam, jadi aman. Kalau memang terlanjur logout, tinggal login lagi dengan `admin`/`admin123`. |
| **Satu orang pegang mouse** | Pembicara yang sedang segmen yang memegang mouse. Yang lain menjelaskan atau mengangguk. Empat orang berebut mouse = dosen bingung. |
| **Sebut nama** | "Sekarang saya **[nama]**, menunjukkan ..." — dosen tahu harus bertanya siapa kalau ada pertanyaan. |
| **Jeda maksimal 3 detik** | Kalau ada yang butuh waktu klik, pembicara lain langsung bantu jelaskan. Jangan diam — diam terlihat seperti macet. |
| **Papan kontingensi** | `docs/screenshots/` (7 PNG) sudah dibuka di tab terpisah sebelum presentasi. Kalau browser/proyektor mati, langsung pindah ke screenshot. |
| **Reset kalau perlu** | Kalau data jadi kacau di tengah jalan: `docker compose down -v` → `docker compose up --build` → ulangi segmen 1 (login). |

## Peta segmen: feature tour per pembicara (±15 menit)

| Segmen | Pembicara | Layar yang ditunjukkan | Slide terkait | Durasi |
|---|---|---|---|---|
| **1** | **P1** | Login → **Dashboard** → **Pelanggan** → **Pengaturan** | Slide 1–2 | 3 mnt |
| **2** | **P2** | **Bahan Baku** → **Menu & Resep** → **Paket & HPP** → **Kalkulator HPP** | Slide 3–5 | 5 mnt |
| **3** | **P3** | **Pembelian** → **Pesanan** → Detail pesanan → **Produksi** | Slide 6 | 4 mnt |
| **4** | **P4** | **Laporan** → **Biaya Operasional** → bukti pengujian di terminal | Slide 8–10 | 3 mnt |

> Langkah bertanda **★** itu opsional. Kalau waktu mepet, lewati yang ★ — inti demo tetap utuh.
> Semua layar yang disebut di atas bisa diklik dari **sidebar kiri**; tidak ada menu yang
> harus diketik manual.

### 13 layar aplikasi, dibagi ke 4 pembicara

| # | Layar | URL | Pembicara |
|---|---|---|---|
| 1 | Masuk | `/login` | **P1** |
| 2 | Dashboard | `/` | **P1** |
| 3 | Pelanggan | `/pelanggan` | **P1** |
| 4 | Pengaturan | `/pengaturan` | **P1** |
| 5 | Bahan Baku | `/bahan` | **P2** |
| 6 | Menu & Resep | `/menu` | **P2** |
| 7 | Paket & HPP | `/paket` | **P2** |
| 8 | Kalkulator HPP | `/hpp` | **P2** |
| 9 | Pembelian | `/pembelian` | **P3** |
| 10 | Pesanan | `/pesanan` | **P3** |
| 11 | Produksi | `/produksi` | **P3** |
| 12 | Laporan | `/laporan` | **P4** |
| 13 | Biaya Operasional | `/biaya` | **P4** |

> Tepat **13 halaman** — cocok dengan angka "13 halaman" di slide 4. Semuanya ada di
> sidebar kiri dengan label Bahasa Indonesia, jadi tidak ada istilah teknis yang perlu
> diterjemahkan saat demo. Satu layar tambahan yang bukan halaman utuh: **Detail Pesanan**
> (`/pesanan/{id}`), dibuka P3 dengan klik nama pelanggan.

---

## SEGMENT 1 — PEMBICARA 1: Masuk, Dashboard, Pelanggan, Pengaturan (±3 menit)

**Layar yang dibahas:** `/login` → `/` (Dashboard) → `/pelanggan` → `/pengaturan`

| # | Layar & aksi | **Yang harus ditunjukkan ke dosen** |
|---|---|---|
| 1 | **Login** — buka `http://localhost:8080`, isi `admin` / `admin123`, klik **Masuk** | Halaman login: logo "SIM Dina Catering", kartu **Masuk**, catatan "Akun demo: `admin` · `admin123`". Session dibuat server-side (tanpa tabel session), TTL 12 jam. |
| 2 | **Dashboard** (otomatis terbuka) | 4 kartu statistik: Pelanggan · Menunggu pembayaran · Pemasukan bulan ini · Keuntungan bulan ini. |
| 3 | Tunjuk **banner merah** di paling atas | "**Perhatian:** 1 paket punya margin di bawah 30%. Update harga jual atau turunkan HPP." |
| 4 | Tunjuk tabel **"Paket dengan margin di bawah 30%"** di paling bawah | **Premium**: HPP/porsi **30.800** · margin **di bawah 30%** (23.0% hari kerja / 28.4% akhir pekan) · min. harga **44.000**. |
| 5 | Tunjuk kartu **"Menunggu pembayaran"** | Ada label tenggat: "DP terlambat", "Pelunasan jatuh tempo" — dihitung dari H-30/H-7, bukan diisi manual. |
| 6 | Tunjuk **"Agenda acara 7 hari ke depan"**, **"Pesanan per status"**, **"Pesanan per sumber"** | Seed pakai `CURRENT_DATE + hari_ke` → agenda selalu relevan kapan pun Anda demo. |
| 7 | Sidebar → **Pelanggan** | Form "Tambah pelanggan" (nama, kontak WhatsApp, alamat, catatan) + tabel 6 pelanggan + jumlah pesanan per pelanggan. |
| 8 | Sidebar → **Pengaturan** | "Biaya operasional" **10%** — disimpan di **database**, bisa diubah owner dari UI (bukan edit env var). Tabel "Seluruh nilai pengaturan" (key-value). |

**Yang diucapkan (P1):**

> "Selamat datang. Ini aplikasi untuk owner — satu-satunya pihak yang masuk ke aplikasi
> adalah owner sendiri. Pelanggan di sini adalah **data master**, bukan pengguna aplikasi.
>
> Ini dashboard-nya. Semua angka di sini dihitung ulang oleh sistem setiap kali datanya
> berubah — tidak ada satu pun yang saya hitung manual. Yang paling penting justru banner
> merahnya: ada satu paket, Premium, yang margin-nya di bawah 30 persen, dan sistem yang
> memberi tahu. Badge yang sama muncul di daftar paket dan di laporan.
>
> Di halaman Pelanggan, owner menyimpan data pelanggan. Dan di Pengaturan, persentase biaya
> operasional 10% disimpan di database, jadi bisa diubah dari layar — tidak perlu terminal."

**Serah terima:**

> "Untuk masuk ke inti aplikasinya, saya serahkan ke **[nama P2]** — dari bahan baku sampai
> angka HPP."

---

## SEGMENT 2 — PEMBICARA 2: Bahan Baku, Resep, Paket, Kalkulator HPP (±5 menit)

**Layar yang dibahas:** `/bahan` → `/menu` → `/paket` → `/hpp`

> Ini **segmen paling penting** — dan inti dari seluruh aplikasi. Ajarkan urutannya:
> bahan → resep → paket → HPP.

| # | Layar & aksi | **Yang harus ditunjukkan ke dosen** |
|---|---|---|
| 1 | Sidebar → **Bahan Baku** | Tabel 9 bahan: nama, satuan, **stok saat ini**, **harga terakhir**. Tombol **Harga** di setiap baris. |
| 2 | Di baris **Daging Sapi** klik tombol **Harga** | Panel "Riwayat harga · Daging Sapi": satu baris, **135.000**, berlaku mulai 1 Mei 2026. Ini tabel `harga_bahan_history` — harga **bertanggal berlaku**, bukan satu kolom. |
| 3 | Sidebar → **Menu & Resep** | Daftar 5 menu. Klik tombol **Resep** pada **Nasi Putih** → takaran per porsi: Beras Premium **0.15 kg**, Minyak Goreng **0.015 liter**. Ini tabel `resep_item`. |
| 4 | Sidebar → **Paket & HPP** | Tabel 3 paket: harga biasa, **harga weekend**, **HPP/porsi**, **margin**, **min. harga 30%**. Premium punya badge merah **"margin < 30%"**, min. harga **44.000**. |
| 5 | Sidebar → **Kalkulator HPP** → dropdown **Paket = Premium** → **Tanggal acuan = HARI KERJA** (lihat catatan di atas) → klik **Hitung HPP** | Panel hasil: HPP per porsi **30.800** · harga jual **40.000** · margin **23.0%** · minimum harga 30% **44.000** · tabel rincian biaya tiap bahan: Nasi Putih 2.670, Rendang Daging 20.970, Sayur Sop 2.110, Buah Potong 2.250. |
| 6 | **★ Uji harga bahan** — kembali ke **Bahan Baku** → Daging Sapi → **Harga** → isi **150.000** → **Simpan harga baru** → kembali **Kalkulator HPP** → **Hitung HPP** lagi | HPP **30.800 → 33.275**, margin **23.0% → 16.8%**, min. harga **44.000 → 47.536**. Tidak ada satu pun field kalkulator yang disentuh. |
| 7 | **★ Bukti historis** — di Kalkulator HPP ubah **Tanggal acuan** ke **1 Juni 2026** → **Hitung HPP** | HPP **kembali 30.800**, margin **23.0%** — harga daging hari ini tidak merusak angka bulan lalu. |
| 8 | **★ Kembalikan data** — **Bahan Baku** → Daging Sapi → **Harga** → isi **135.000** → **Simpan harga baru** | Data demo bersih lagi (supaya kalau demo diulang, angkanya sama). |

**Yang diucapkan (P2):**

> "Di layar-layar ini ada dua keputusan desain yang membentuk hampir semua logika sistem.
> Pertama, harga bahan disimpan sebagai **riwayat bertanggal berlaku**, bukan satu kolom —
> itu yang membuat HPP historis bisa dihitung ulang. Kedua, tiap menu punya **resep**:
> bahan apa saja dan takaran berapa per porsi.
>
> Dari situ HPP dihitung otomatis. Untuk paket Premium: HPP bahan 28.000, ditambah 10% biaya
> operasional, jadi **30.800** per porsi. Harga jual 40.000, jadi margin 23 persen. Dan
> harga jual minimum untuk margin 30 persen adalah 44.000.
>
> Dan ini buktinya. Saya naikkan harga daging dari 135.000 menjadi 150.000 — tidak ada satu
> pun field di kalkulator HPP yang saya sentuh. HPP Premium naik dari 30.800 ke 33.275, dan
> margin turun dari 23 persen ke 16,8 persen. Badge peringatan di halaman tadi juga makin merah.
>
> Terakhir, kalau saya hitung ulang dengan tanggal acuan bulan lalu, angkanya kembali ke 30.800.
> Harga bulan lalu tidak ikut berubah karena harga hari ini naik — itu gunanya harga disimpan
> sebagai riwayat."

**Serah terima:**

> "Sekarang saya serahkan ke **[nama P3]** — bagaimana angka-angka itu dipakai untuk pesanan
> nyata."

---

## SEGMENT 3 — PEMBICARA 3: Pembelian, Pesanan, dan Papan Produksi (±4 menit)

**Layar yang dibahas:** `/pembelian` → `/pesanan` → `/pesanan/{id}` → `/produksi`

| # | Layar & aksi | **Yang harus ditunjukkan ke dosen** |
|---|---|---|
| 1 | Sidebar → **Pembelian** | Form "Catat pembelian": bahan, jumlah, harga satuan, pasar, tanggal + tabel "Riwayat pembelian". Satu transaksi → **tiga efek**: stok naik, harga terakhir diperbarui, riwayat harga nambah. |
| 2 | Isi form: **Ayam Potong** · jumlah **20** · harga satuan **38.000** · Pasar **Pasar Klewer** · tanggal hari ini → **Simpan pembelian** | Riwayat bertambah satu baris, total **Rp760.000**. |
| 3 | Sidebar → **Pesanan** → form "Pesanan baru": Pelanggan **Ratna Wulandari** · Paket **Standar** · **60** tamu · tanggal acara **Sabtu/Minggu minggu depan** · sumber **Instagram** → **Simpan pesanan** | Toast: "Pesanan tersimpan, tenggat DP H-30 dan pelunasan H-7 dihitung otomatis". |
| 4 | Tunjuk baris pesanan itu di tabel bawah | Total **Rp2.280.000** — bukan 60 × 35.000, karena acaranya akhir pekan sehingga harga **Rp38.000** dipakai. Owner tidak perlu hitung sendiri. |
| 5 | Klik nama **Ratna Wulandari** → **detail pesanan** | Kartu ringkasan: tanggal acara, total, status pembayaran "Belum DP", margin untung **50.0%**. Tabel bawah: **Tenggat DP (H-30)** dan **Tenggat pelunasan (H-7)** terisi tanggal otomatis. |
| 6 | Kartu "Catat pembayaran": Jenis **DP 50%** (nominal terisi sendiri **Rp1.140.000**) · metode **Transfer** → **Simpan pembayaran** | Status jadi **"DP terbayar"**, status pesanan naik `baru` → **Final Meeting**, dan tabel "Tugas produksi" otomatis berisi **4 baris**. |
| 7 | Sidebar → **Produksi** | **Papan produksi**: tugas dikelompokkan per tanggal — Persiapan bahan & alat (H-1) · Masak seluruh menu (H) · Plating & cek porsi (H) · Pengiriman ke lokasi (H). Tiap kartu ada kolom **Penanggung jawab** + **Status tugas**. |
| 8 | Di satu kartu: isi **PJ** = `Sari` → **Simpan**; lalu **Status tugas** `Belum` → `Dalam proses` → **Simpan** | Kartu berubah: ada PJ, badge status berubah. "Ini menggantikan papan tulis di dapur." |
| 9 | **★ Stok berkurang** — detail pesanan → "Ubah status pesanan" → **Produksi** → simpan; lalu cek **Bahan Baku** | Stok **Ayam Potong** turun **15 kg** (resep 0,25 kg × 60 tamu), **Beras Premium** turun 9 kg. Dikurangi dalam satu transaksi database, jadi tidak mungkin dobel. |

**Yang diucapkan (P3):**

> "Alurnya seperti di dunia nyata. Owner mencatat pembelian di pasar — dalam satu transaksi,
> stok naik, harga terakhir diperbarui, dan satu baris riwayat harga masuk. Efek terakhir itu
> yang penting: artinya HPP semua paket ikut terhitung ulang, dan tidak ada satu pun langkah
> yang bisa terlewat.
>
> Sekarang owner mencatat pesanan. Waktu pesanan dibuat, sistem menghitung sendiri tiga hal:
> total harga — pakai harga akhir pekan karena acaranya hari Sabtu — tenggat DP H-30, dan
> tenggat pelunasan H-7. Owner tidak perlu menghitung apa pun.
>
> Dan bagian yang paling menarik: begitu DP 50% dicatat, **empat tugas produksi otomatis
> dibuat** untuk hari acara. Alasannya bisnis — begitu DP 50% masuk, pesanan itu sudah pasti
> jalan, jadi jadwal dapur harus langsung tersedia. Pembuatannya idempoten, jadi tugas tidak
> akan pernah dobel. Di halaman Produksi, tugas itu dikelompokkan per tanggal, lengkap dengan
> penanggung jawab dan status — menggantikan papan tulis di dapur.
>
> Terakhir, kalau status pesanan diubah ke Produksi, stok bahan otomatis berkurang: 0,25 kg
> ayam per porsi dikali 60 tamu, jadi 15 kg, dalam satu transaksi."

**Serah terima:**

> "Terima kasih **[nama P2]**. Saya serahkan ke **[nama P4]** — menutup dari laporan
> laba rugi sampai bukti pengujian."

---

## SEGMENT 4 — PEMBICARA 4: Laporan, Biaya Operasional, Bukti Mutu (±3 menit)

**Layar yang dibahas:** `/laporan` → `/biaya` → terminal

| # | Layar & aksi | **Yang harus ditunjukkan ke dosen** |
|---|---|---|
| 1 | Sidebar → **Laporan** → pilih **bulan ini** → **Terapkan** | 4 kartu: **Pemasukan** · **Pengeluaran bahan** · **Biaya operasional** · **Keuntungan**. Di bawahnya: rekap pesanan per status & per sumber, tabel **Margin per pesanan**, tabel **Paket dengan margin di bawah 30%**, dan **Stok bahan saat ini**. |
| 2 | Tunjuk angka **Keuntungan** | "Keuntungan = Pemasukan − Pengeluaran bahan − Biaya operasional." Bandingkan dengan kartu "Keuntungan bulan ini" di dashboard — angkanya sama, jadi konsisten. |
| 3 | Tunjuk tabel **Margin per pesanan** | Setiap pesanan punya margin; pesanan yang baru dibuat (Standar, tanggal akhir pekan) **50.0%** — lebih sehat dari 45.7% di hari kerja, karena harga jual akhir pekan lebih tinggi. |
| 4 | Sidebar → **Biaya** | Form catat biaya aktual (gas, listrik, mobilitas, perlengkapan) + tabel riwayatnya. Penting: laporan memakai biaya **aktual**, bukan estimasi 10% — jadi tidak tumpang tindih dengan kalkulator HPP. |
| 5 | **★ Terminal:** `bash scripts/smoke-test.sh` | 18 langkah `[OK]` bergantian. "Ini bukti saya tidak cuma demo manual." |
| 6 | **★ Terminal:** `go test ./...` | `ok` di 2 paket: `internal/service` dan `internal/handler`. |
| 7 | **★ Tanpa JavaScript** — matikan JS di browser, buka `/hpp/hasil` | Tetap **halaman penuh yang berfungsi**, bukan error. Bukti keputusan desain nomor 5. |

**Yang diucapkan (P4):**

> "Ini rekap bisnisnya. Laba rugi bulanan, rekap pesanan, margin per pesanan, dan posisi stok
> bahan. Keuntungan = pemasukan dikurangi pengeluaran bahan dikurangi biaya operasional —
> dan angkanya cocok dengan yang muncul di dashboard, jadi konsisten.
>
> Perhatikan satu hal: laporan ini memakai biaya operasional **aktual** yang dicatat, bukan
> estimasi 10%. Jadi dua perhitungan itu tidak tumpang tindih.
>
> Terakhir, saya jalankan smoke test: 18 pemeriksaan end-to-end lewat HTTP, semuanya `[OK]`.
> Ini yang membuat saya percaya sistem ini benar-benar jalan, bukan cuma terlihat jalan."

**Penutup (lalu ke Slide 10):**

> "That's the whole application running live. Sekarang kita kembali ke slide terakhir untuk penutup."

---

## Jalur cepat (kalau waktu tinggal 5 menit)

| Prioritas | Yang DIPERTAHANKAN | Yang DILEWATI |
|---|---|---|
| **P1** | Login → banner merah → tabel margin Premium → Pelanggan | Daftar menunggu bayar, agenda, Pengaturan |
| **P2** | Paket & HPP → Kalkulator HPP → harga daging 150.000 → HPP 33.275 | Menu & Resep, bukti historis, mengembalikan harga |
| **P3** | Pesanan → total 2.280.000 → DP → 4 tugas | Pembelian, set PJ, status Produksi |
| **P4** | Laporan → kartu keuntungan → tabel margin | Biaya, smoke test, `go test` |

> **Jangan pernah** melewati langkah "ubah harga bahan → HPP berubah" (P2 langkah 6).
> Itu inti dari seluruh aplikasi, dan itu yang paling sering jadi bahan pertanyaan dosen.

## Papan kontingensi (kalau ada masalah saat demo)

| Gejala | Yang harus dilakukan |
|---|---|
| Halaman putih / 502 | `docker compose logs -f app` → `docker compose restart app`. Cookie session masih valid setelah restart. |
| "connection refused" | `docker compose ps` → tunggu `postgres` healthy (app sudah retry 30×). |
| Login ditolak | `docker compose down -v` → `up --build`. Cookie browser ikut tidak berlaku karena session ikut hilang. |
| HTMX gagal (toast tidak muncul, panel tak refresh) | **Jangan panik.** Form tetap jalan via full-page reload. Katakan: "Ini justru fallback yang saya rancang." Lanjut demo. |
| Proyektor mati / browser crash | Buka tab `docs/screenshots/` (7 PNG) — 2-dashboard.png, 3-kalkulator-hpp.png, 4-paket-margin.png, 6-produksi.png, 7-laporan.png. |
| WiFi kampus mati | Tidak masalah — HTMX dari CDN, tapi semua form punya `action` biasa, jadi aplikasi tetap berfungsi penuh. |
| Angka HPP tidak sesuai Bagian 3 | Kemungkinan ada yang **mengubah harga bahan** tanpa mengembalikannya. Reset: `docker compose down -v` → `up --build`. |

**Catatan demo — siapkan sebelum presentasi:** Dataset sudah berisi 6 pesanan dengan 6 status
berbeda dan tanggal tersebar dari −10 sampai +60 hari, jadi Agenda 7 hari dan daftar
"menunggu pembayaran" **tidak pernah kosong** — bagus untuk demo.

---

# BAGIAN 3 — ANGKA ACUAN (verifikasi silang saat ditanya)

HPP dihitung dari `db/seed.sql`; angka Premium diverifikasi oleh unit test
`internal/service/hpp_test.go:75-107`.

### Komponen HPP bahan per porsi

| Menu | Komponen | Total |
|---|---|---|
| Nasi Putih | Beras 0,150 kg × 16.000 = 2.400 · Minyak 0,015 ltr × 18.000 = 270 | **2.670** |
| Ayam Bakar | Ayam 0,250 kg × 38.000 = 9.500 · Bumbu 0,040 kg × 12.000 = 480 · Minyak 0,015 × 18.000 = 270 | **10.250** |
| Rendang Daging | Daging 0,150 kg × 135.000 = 20.250 · Bumbu 0,030 × 12.000 = 360 · Minyak 0,020 × 18.000 = 360 | **20.970** |
| Sayur Sop | Sayur 0,200 kg × 8.000 = 1.600 · Bumbu 0,020 × 12.000 = 240 · Minyak 0,015 × 18.000 = 270 | **2.110** |
| Buah Potong | Buah 0,150 kg × 15.000 = 2.250 | **2.250** |

### Rekap paket

| Paket | Isi | HPP bahan | +10% | HPP/por | Harga | Weekend | Margin | Min. harga |
|---|---|---|---|---|---|---|---|---|
| Hemat | Nasi, Sayur, Ayam | 15.030 | 16.533 | **16.533** | 30.000 | 32.000 | 44.9% | 23.619 |
| Standar | + Buah | 17.280 | 19.008 | **19.008** | 35.000 | 38.000 | 45.7% | 27.155 |
| Premium | Nasi, Rendang, Sayur, Buah | 28.000 | 30.800 | **30.800** | 40.000 | 43.000 | **23.0%** ⚠️ | 44.000 |

### Angka lain yang mungkin ditanya

| | Nilai |
|---|---|
| Margin minimum (aturan owner) | **30%** — `internal/model/model.go:76` |
| Biaya operasional default | **10%** — dapat diubah di `/pengaturan` |
| DP | **50%** dari total |
| Tenggat DP / pelunasan | **H-30** / **H-7** dari tanggal acara |
| TTL session | **12 jam** |
| Margin weekend Premium | (43.000 − 30.800) / 43.000 = **28.4%** — juga di bawah 30%, badge tetap muncul |

### Catatan format angka (supaya tidak salah baca di layar)

Aplikasi memakai **titik** sebagai pemisah ribuan dan desimal (hasil `fmt.Sprintf` Go),
jadi yang tampil di layar adalah `30.800`, `44.9%`, `0.25 kg`. Dokumen ini memakai
**koma** untuk desimal saat menulis takaran (`0,250 kg`) karena itu kebiasaan tulis
di Indonesia — angkanya sama, hanya formatnya berbeda. **Tolocalisasi ke bentuk di layar
sebelum membacakan angka ke dosen.**

---

# BAGIAN 4 — PETA KODE: "KALAU DITANYA, BUKA INI"

| Pertanyaan dosen | Buka file:baris |
|---|---|
| "Tunjukkan cara aplikasinya start" | `cmd/server/main.go:32-87` |
| "Semua rutenya di mana?" | `internal/handler/handler.go:58-163` |
| "Kenapa tanpa ORM? Bagaimana SQL-nya?" | `db/queries.sql` → di-*generate* ke `internal/store/queries.sql.go`; interface di `internal/store/querier.go:12-125`; config `sqlc.yaml` |
| "Di mana logika HPP?" | `internal/service/hpp.go:42-79` (rumus), `:84-103` (harga efektif), `:203-238` (daftar paket) |
| "Di mana margin & harga minimum?" | `internal/service/hpp.go:57-79` |
| "Kenapa session tanpa tabel?" | `internal/middleware/session.go:42-134` |
| "Bagaimana middleware auth?" | `internal/middleware/auth.go:23-45` |
| "Di mana password di-hash?" | `internal/handler/auth.go:26` + `internal/seed/seed.go` |
| "Di mana stok berkurang?" | `internal/service/pesanan.go:99-106` → `internal/service/stok.go:14-39` |
| "Di mana transaksi database?" | `internal/service/service.go:66` (`withTx`) |
| "Di mana tugas produksi dibuat?" | `internal/service/pesanan.go:113-147` |
| "Di mana DP memicu-naik status?" | `internal/service/pembayaran.go:131-141` |
| "Di mana jadwal H-30/H-7?" | `internal/service/pembayaran.go:15-18` |
| "Di mana harga akhir pekan dipilih?" | `internal/service/harga.go:24-36` |
| "Di mana laba rugi?" | `internal/service/laporan.go:23-57` |
| "Bagaimana HTMX vs full page dibedakan?" | `internal/handler/handler.go:176-215` + `internal/handler/pelanggan.go:94` (`isHTMX`) |
| "Bagaimana template diparse?" | `internal/handler/handler.go:35-55` |
| "Di mana skema database?" | `db/schema.sql` (DDL) + `db/migrations/0001_init.up.sql` |
| "Di mana data demo?" | `db/seed.sql` (199 baris). Demo dibersihkan dengan `docker compose down -v` lalu `up --build` — seed hanya jalan kalau tabel `users` kosong. |
| "Di mana migrasi dijalankan?" | `internal/database/database.go:50` |
| "Di mana config dibaca?" | `internal/config/config.go:20-32` |
| "Di mana test HPP?" | `internal/service/hpp_test.go` |
| "Di mana test pembayaran?" | `internal/service/pembayaran_test.go` |
| "Di mana test template?" | `internal/handler/template_test.go` |

---

# BAGIAN 5 — BANK Q&A

> 30 pertanyaan yang paling sering muncul saat presentasi tugas akhir, dengan jawaban siap
> pakai. **Jawab dengan jujur** kalau memang belum dikerjakan — itu dinilai lebih tinggi
> daripada dikarang.

## A. Tentang pemilihan teknologi

**1. "Kenapa pakai Go? Kenapa bukan PHP/JS/Python yang lebih umum?"**
Go karena tiga hal: compile-time type checking (banyak bug tertangkap sebelum runtime),
satu binary tanpa runtime dependency (cocok untuk Docker — image kecil, start cepat), dan
performanya cukup untuk aplikasi CRUD seperti ini. Saya juga butuh yang mendukung
*static typing* secara ketat supaya arsitektur berlapis bisa dijaga.

**2. "Kenapa tidak pakai ORM seperti GORM?"**
Ini keputusan sadar, bukan karena belum tahu ORM. Dengan sqlc saya: (1) tetap menulis SQL
asli sehingga performa dan query bisa dikontrol penuh, (2) tapi hasilnya type-safe — compiler
yang-menegakkan, dan (3) tidak ada *magic* — saya tahu persis query apa yang dikirim ke
database. Untuk sistem berbasis laporan dan agregasi seperti ini, SQL eksplisit lebih mudah
di-audit. Konsekuensinya: saya harus benar-benar paham SQL, dan tidak bisa membuat CRUD
sederhana tanpa effort.

**3. "Kenapa HTMX, bukan React/Vue?"**
Karena aplikasinya memang server-side rendering dari awal, dan problem-nya nyata: 13 halaman
form + tabel. Dengan HTMX saya dapat partial refresh tanpa membangun SPA, tanpa build step,
dan tanpa menambah dependency frontend. Yang lebih penting: saya bisa **fallback ke
full-page reload** untuk setiap interaksi, jadi aplikasi tetap utuh kalau HTMX gagal dimuat.

**4. "Kenapa server-side rendering, bukan SPA?"**
Datanya sudah dimuat server; tidak ada kebutuhan state management di client. SSR juga
membuat halaman pertama cepat dan SEO-friendly, dan yang paling terasa: **tidak ada build
step**, jadi tidak ada toolchain frontend yang harus dijaga. Untuk aplikasi internal
seperti ini, pendekatan ini paling hemat dan paling mudah dirawat.

**5. "Kenapa CSS-nya ditulis sendiri, tidak pakai Bootstrap/Tailwind?"**
Saya memang butuh sistem desain yang konsisten, tapi dependensi frontend saya tetap 1
(HTMX). 814 baris CSS vanilla dengan CSS custom properties — termasuk dukungan dark mode
otomatis lewat `prefers-color-scheme` — sudah cukup untuk 13 layar, dan saya paham setiap
barisnya karena saya yang menulisnya.

**6. "Kenapa CSS-nya sampai 814 baris? Bukankah itu banyak?"**
Ya, tapi itu **satu file** dan menggantikan third-party. 814 baris untuk 13 layar dengan
komponen kartu, tabel, badge, form, toast, papan produksi, dan dark mode — itu rasio yang
wajar kalau tidak pakai framework. Alternatifnya menarik 200KB Bootstrap yang 80%-nya
tidak saya pakai.

## B. Tentang arsitektur & kualitas kode

**7. "Bagaimana Anda memastikan layering-nya benar, tidak bocor?"**
Tiga aturan yang dijaga dan bisa diaudit: (1) tidak ada SQL di handler — semua query lewat
interface `store.Querier` yang di-*generate* sqlc; (2) semua perhitungan numerik hanya di
`internal/service`; (3) handler tidak pernah mengubah `service` secara langsung, lewat
method yang sudah deterministik. Saya juga menulis unit test di layer service supaya
perilaku bisnisnya terkunci di satu tempat.

**8. "Bagaimana kalau ada perubahan skema di masa depan?"**
sqlc di-*generate* dari `db/schema.sql` — saya ubah skema, jalankan `sqlc generate`, dan
compiler akan menunjuk di mana kode saya harus disesuaikan. Ditambah migrasi
golang-migrate yang sudah tertanam di binary, jadi migrasi selalu sesuai versi binary-nya.

**9. "Kenapa hanya ada satu file migrasi?"**
Karena project ini belum pernah di-deploy ke environment yang sudah ada datanya — semua
skema dibangun dalam satu kali. Jadi begitu ada iterasi skema di masa depan, file
`0002_...` tinggal ditambah. Embedding-nya sudah disiapkan, jadi menambah migrasi baru
tidak mengubah arsitektur.

**10. "Apa saja yang belum Anda kerjakan?"**
Ini jawaban jujur yang selalu dihargai:
- **Belum multi-user.** Kolom `role` ada tapi belum jadi middleware pembatasan akses.
- **Belum ada notifikasi otomatis** (WhatsApp/email) untuk tenggat DP dan pelunasan.
- **Belum ada unit test untuk handler HTTP dan middleware** — lapisan itu sekarang
  ditutup oleh smoke test end-to-end.
- **Belum ada integrasi pembayaran digital** — pembayaran dicatat manual.
- **Belum ada dashboard analitik lanjutan** (grafik tren harga, forecast).

## C. Tentang keamanan

**11. "Bagaimana password disimpan?"**
bcrypt, dan **hash-nya dihitung saat runtime** — tidak ada hash yang ditulis di dalam kode
maupun di `seed.sql`. Yang ada di `db/seed.sql` cuma placeholder `{PASSWORD_HASH}` yang
diganti dengan hash bcrypt dari environment variable `ADMIN_PASSWORD` saat boot.

**12. "Kenapa session-nya dibuat sendiri, tidak pakai library?"**
Karena kebutuhannya sempit: user ID, username, role, dan waktu kedaluwarsa. Saya encode
payload base64, tandatangani dengan **HMAC-SHA256**, bandingkan dengan `hmac.Equal` (constant-time,
mencegah timing attack). Cookie-nya `HttpOnly` + `SameSite=Lax`, TTL 12 jam. Dengan begitu
tidak ada tabel session, tidak ada Redis, tidak ada dependency tambahan. Trade-off-nya
saya akui: cookie ini **tidak bisa dicabut server-side** (tidak ada "logout semua perangkat")
— kalau logout, token-nya masih valid sampai kedaluwarsa.

**13. "Apakah aman dari SQL injection?"**
Ya, karena sqlc meng-*generate* kode dengan **prepared statement parameterized** —
nilai selalu dikirim sebagai parameter terpisah, tidak pernah digabungkan ke string SQL.
Tidak ada satu pun tempat di kode saya yang membangun query dengan penggabungan string.
Tambahan: `users.username` punya constraint `UNIQUE`, dan password dibandingkan dengan
`bcrypt.CompareHashAndPassword` yang aman terhadap timing attack.

**14. "Bagaimana dengan CSRF?"**
Ini **limitasi yang saya akui**: cookie `SameSite=Lax` memberi perlindungan sebagian
terhadap POST lintas-situs, dan `HttpOnly` mencegah cookie dibaca JavaScript — tapi saya
belum menambahkan token CSRF eksplisit. Ini termasuk rencana perbaikan. Untuk konteks:
aplikasi ini internal dan hanya satu owner yang menggunakan, dengan deploy di jaringan internal,
sehingga risikonya relatif rendah — tapi saya tahu ini bukan alasan yang baik untuk produksi.

**15. "Kenapa hanya satu akun admin? Bagaimana kalau ada banyak owner?"**
Kebutuhan real-nya memang satu owner. Skema sudah punya kolom `role` dan tabel `users`
yang mendukung banyak akun — yang belum ada adalah UI dan middleware untuk mengelola
akun dan pembatasan akses per peran. Itu salah satu item prioritas tahap berikutnya.

**16. "Kenapa `SESSION_SECRET` wajib diisi? Kenapa tidak default?"**
Karena default yang dipublikasikan di repository berarti siapa pun yang tahu bisa
memalsukan cookie session. Jadi aplikasi **menolak start** kalau `SESSION_SECRET` kosong —
fail fast, bukan gagal diam-diam di tengah jalan. Di `.env.example` nilainya juga sudah
ditandai harus diganti.

## D. Tentang domain bisnis & perhitungan

**17. "Kenapa biaya operasional ditambahkan proporsional ke HPP bahan, bukan ditambahkan di akhir?"**
Supaya rumus `harga jual minimum = HPP / 0,7` tetap konsisten. Kalau biaya operasional
ditambahkan setelah margin dihitung, maka harga minimum yang dihitung tidak akan
menghasilkan margin 30% yang sesungguhnya. Urutan "biaya → HPP total → margin" menjaga
kedua angka itu konsisten satu sama lain. Ini saya dokumentasikan di `README.md` §4.
Catatan penting: laporan laba rugi memakai **biaya operasional aktual yang dicatat**
(bukan estimasi persentase), supaya dua perhitungan itu tidak tumpang tindih.

**18. "Kenapa harga bahan disimpan sebagai riwayat, bukan satu kolom?"**
Karena HPP harus bisa dihitung ulang secara historis. Kalau cuma ada kolom `harga_terakhir`,
maka HPP bulan Mei akan ikut berubah ketika harga daging naik hari ini — dan laporan
keuntungan bulan Mei jadi bohong. Dengan `harga_bahan_history` bertanggal berlaku,
`hargaEfektif(bahan, t)` mengambil harga terakhir yang berlaku pada tanggal t, jadi
historis tetap konsisten.

**19. "Kenapa harga paket disimpan per tahun berlaku?"**
Alasan yang sama: harga juga historis. Kalau owner menaikkan harga paket di tengah tahun,
pesanan yang sudah dibuat bulan lalu harus tetap memakai harga yang berlaku saat itu.
Kolom `tahun_berlaku` + `UNIQUE(nama, tahun_berlaku)` membuat satu pesanan terkunci ke
satu baris harga secara alami.

**20. "Kenapa `BIGINT` untuk uang dan `DOUBLE` untuk kuantitas? Bukankah tidak konsisten?"**
Sengaja, dan memang ada konsekuensinya. Uang dalam Rupiah bulat tidak pernah pecahan —
jadi `BIGINT` menghindari masalah pembulatan floating point pada nominal. Kuantitas bahan
(0,150 kg) memang pecahan, jadi `DOUBLE PRECISION` dipakai. Konsekuensinya: pembulatan
terjadi di level per komponen resep (`math.Round` di `BiayaKomponen`), bukan di akhir.
Ini keputusan yang bisa diperdebatkan — kalau presisi lebih tinggi dibutuhkan, `NUMERIC`
adalah pilihan yang lebih tepat.

**21. "Kenapa stok hanya berkurang saat status jadi `produksi`, bukan saat pesanan dibuat?"**
Karena pesanan bisa dibatalkan setelah DP masuk. Kalau stok dipotong saat pesanan dibuat,
pesanan yang dibatalkan akan menyebabkan stok berkurang padahal bahan belum dimasak.
Titik pengurangan saya pilih saat bahan benar-benar dipakai, dan seluruhnya dalam satu
transaksi database supaya tidak mungkin berkurang dobel.

**22. "Bagaimana kalau tugas produksi terduplikasi?"**
Pembuatannya idempoten — `buatTugasDefault` mengecek dulu apakah sudah ada tugas untuk
pesanan itu (`internal/service/pesanan.go:114-120`), dan kalau ada, dia keluar tanpa
menambah apa pun. Jadi tidak masalah baik tugas dibuat saat DP dicatat maupun saat status
diubah ke `produksi` — keduanya aman.

**23. "Kenapa tanggal acara dihitung relatif `CURRENT_DATE + hari_ke` di seed?"**
Supaya demo selalu relevan: ada pesanan yang sudah lewat (selesai), ada yang sedang
berjalan (produksi), dan ada yang akan datang (agenda 7 hari ke depan) — kapan pun Anda
presentasikan. Kalau tanggalnya hardcode, dati demo bisa jadi basi.

## E. Tentang pengujian & operasional

**24. "Bagaimana cara Anda memastikan sistemnya benar?"**
Tiga lapis, dan saya bisa menunjukkan secara langsung:
1. **15 fungsi unit test** — HPP, margin, harga minimum, jadwal pembayaran, status
   pembayaran (13 di `internal/service`) + 2 test template di `internal/handler`.
   Angka HPP Paket Premium saya jadikan test case supaya regresi terdeteksi.
2. **`go vet` + `go build`** untuk analisis statis.
3. **`scripts/smoke-test.sh`** — 18 pemeriksaan end-to-end lewat HTTP, termasuk yang
   paling penting: harga bahan diubah → HPP berubah → dikembalikan.

Sisanya saya akui belum ada: unit test untuk handler HTTP dan middleware session.

**25. "Bagaimana cara deploy-nya?"**
`docker compose up --build`. Dua service: `postgres:16-alpine` dengan healthcheck
`pg_isready`, dan `app` yang build dari `Dockerfile` multi-stage (build `golang:1.26-alpine`
→ runtime `alpine:3.20`, non-root user `appuser`, `TZ=Asia/Jakarta`). Migrasi dan seed
otomatis jalan saat container pertama kali start. Database, kode aplikasi, template, dan
statis semuanya di-*embed* ke dalam satu binary — jadi tidak ada volume kode yang harus
di-mount.

**26. "Kalau dipresentasikan tanpa internet, HTMX-nya gagal dimuat — bagaimana?"**
Aplikasi tetap berfungsi penuh, dan ini memang saya rencanakan. Setiap form punya pasangan
`method`/`action` biasa, dan setiap endpoint panel mengembalikan **partial** kalau request-nya
dari HTMX (`HX-Request: true`), atau **halaman penuh** kalau request biasa. Saya bisa
membuktikannya secara live: matikan JavaScript, buka `/hpp/hasil` — tetap jalan.

**27. "Kenapa tidak pakai Redis atau message queue?"**
Kebutuhannya belum ada. Trafik satu owner dengan satu container; session sudah cukup
ditangani cookie bertanda tangan tanpa state server. Menambah Redis berarti menambah satu
service lagi untuk masalah yang belum ada.

## F. Pertanyaan reflektif yang sering muncul di penutup

**28. "Kalau diberi waktu lebih, apa yang akan Anda kerjakan lebih dulu?"**
Tiga hal, urut: **(1) pembatasan akses berbasis role** — fondasi agar aplikasi bisa dipakai
lebih dari satu orang, dan ini juga yang akan jadi prasyarat fitur **(2) notifikasi
tenggat via WhatsApp**, karena di katering nyata keterlambatan bayar DP itu kerugian langsung;
**(3) pelaporan historis** — grafik tren harga bahan dan margin per bulan, supaya owner
bisa memprediksi kapan harus menaikkan harga paket, bukan baru tahu setelah kejadian.

**29. "Apa kontribusi orisinal Anda?"**
Tiga hal: **(1)** model HPP historis — harga bahan dan harga paket sama-sama bertanggal
berlaku, sehingga margin masa lalu bisa diaudit ulang; **(2)** otomatisasi rantai
konsekuensi — satu pencatatan pembelian efeknya berantai sampai mengubah HPP, margin, dan
peringatan, dan satu pencatatan DP sampai membuat jadwal produksi; **(3)** jaminan
operasional — aplikasi tetap berfungsi tanpa JavaScript, dan seluruhnya jalan dengan satu
`docker compose up`.

**30. "Apa yang paling Anda pelajari dari project ini?"**
Bahwa **sebagian besar pekerjaan nyata justru ada di angka, bukan di layar.** Estetika
antarmuka itu pekerjaan beberapa hari; memastikan margin paket Premium benar-benar 23%,
bahwa harga jual minimum dihitung dari HPP yang sama dengan yang dipakai saat membuat
pesanan, dan bahwa pembelian di pasar benar-benar mengubah stok — itu pekerjaan berbulan-bulan
dan tidak terlihat kalau tidak ditest. Project ini juga mengajari saya bahwa **uji otomatis
itu bukan formalitas** — smoke test 18 langkah itu yang membuat saya percaya sistem saya
benar-benar jalan, bukan cuma terlihat jalan.

---

# BAGIAN 6 — CHECKLIST SEBELUM PRESENTASI & TROUBLESHOOTING

## Checklist H-1 (sehari sebelum)

- [ ] `docker compose down -v` lalu `docker compose up --build` berhasil — **bangun image
      sekarang**, jangan dekat presentasi. (`down -v` membuang volume lama supaya data
      demo bersih dan semua angka di dokumen ini pasti cocok.)
- [ ] `go vet ./... && go build ./... && go test ./...` → semua `ok`.
- [ ] `bash scripts/smoke-test.sh` → 18 `[OK]`, tanpa `[FAIL]`.
- [ ] Login `admin`/`admin123` berhasil.
- [ ] Buka **satu** browser untuk demo (Chrome/Edge). Nonaktifkan ekstensi yang mengganggu.
- [ ] Zoom browser 100% supaya angka rupiah terbaca dari kursi dosen.
- [ ] Siapkan **fallback offline**: `docs/screenshots/` (7 PNG) — login, dashboard, kalkulator
      HPP, paket & margin, resep, produksi, laporan. Buka dulu semua file-nya untuk
      memastikan bisa dibuka tanpa internet.
- [ ] Siapkan `Diagram,dll/PNG/` (9 diagram) kalau perlu menampilkan ERD/DFD/use case.
- [ ] Cetak / buka file ini di laptop kedua atau print — jangan hanya mengandalkan proyektor.
- [ ] **Cetak 4 bagian naskah ini** — satu untuk tiap anggota kelompok (Bagian 0 pembagian
      peran + Bagian 1 slide milikmu + Bagian 2 segmen milikmu). HP masing-masing juga
      sudah berisi file-nya, untuk berjaga-jaga kalau laptop gagal.
- [ ] **Dry run 15 menit sekali bersama**: P1 pegang mouse, jalankan segmen 1–2; P3
      isi form pesanan; P4 buka terminal. Pastikan tidak ada yang perlu cari tahu
      "form-nya di mana" saat demo.
- [ ] Sepakati siapa yang jawab kalau dosen bertanya di luar bidangnya — misalnya:
      "Itu bagian **[nama]**, tapi kalau saya boleh jawab…"

## Checklist H-0 (jam sebelum)

- [ ] Docker Desktop **running** dan `docker compose ps` menunjukkan 2 service healthy.
- [ ] Reset data supaya demo bersih: `docker compose down -v` lalu `docker compose up --build`.
- [ ] **Cek bahan setelah reset**: di `/bahan`, harga terakhir **Ayam Potong = 38.000** dan
      **Daging Sapi = 135.000**. Kalau angkanya lain, volumenya belum bersih — ulangi reset.
- [ ] **Cek tanggal di layar kalkulator HPP**: `Tanggal acuan` harus **hari kerja**
      (Senin–Jumat). Kalau hari ini akhir pekan, isi **Senin berikutnya** — biar HPP
      tetap **30.800** dan margin **23.0%** seperti di slide.
- [ ] Buka http://localhost:8080 sekali, login sekali, lalu **logout** — memastikan session
      cookie tidak tertinggal dari browser lain.
- [ ] Cek **HTMX termuat**: tekan F12 → Network → pastikan `htmx.min.js` status 200.
      Kalau gagal, demo tetap jalan (fallback), tapi toast tidak akan muncul — turunkan
      ekspektasi notifikasi.
- [ ] Buka tab kosong untuk: `README.md`, `Diagram,dll/04-ERD.md`,
      `internal/service/hpp.go`, `internal/handler/handler.go`.
- [ ] **Terminal sudah terbuka di folder project**, siap untuk `bash scripts/smoke-test.sh`
      (dipakai P4 di segmen 4).
- [ ] Matikan notifikasi Slack/WhatsApp/email di laptop.
- [ ] Siapkan air minum. Presentasi ±35 menit butuh banyak minum dan suara yang jelas.

## Troubleshooting cepat

| Gejala | Penyebab & solusi |
|---|---|
| `docker compose up` gagal, port 8080 dipakai | `docker compose down` lalu `netstat -ano \| findstr :8080` untuk cari proses, atau ubah `PORT` di `.env` |
| Halaman "`SESSION_SECRET wajib diisi`" | `.env` tidak ada / `SESSION_SECRET` kosong. Salin `.env.example` → `.env`, isi nilainya. |
| Login ditolak padahal benar | Database masih punya data lama dengan password berbeda. Reset: `docker compose down -v` lalu `up --build`. |
| Halaman kosong / template error | `go build` gagal. Jalankan `go test ./internal/handler/...` — test template memunculkan pesan error yang jauh lebih jelas. |
| "connection refused" dari app ke Postgres | `postgres` belum healthy. `docker compose logs -f postgres`. App sudah *retry* 30× — tunggu, atau `docker compose restart app`. |
| Table/panel HTMX tidak ter-refresh | HTMX gagal dimuat (CDN). Form tetap jalan via full reload — jelaskan itu sebagai fallback, bukan ganggu demo. |
| Toast tidak muncul | `HX-Trigger` butuh HTMX aktif. Disesuaikan: toast hanya pemanis, bukan inti demo. |
| Angka HPP tidak sesuai Bagian 3 | Kemungkinan ada yang **mengubah harga bahan** lewat `/bahan/{id}/harga` tanpa mengembalikannya. Reset database. |
| Tanggal demo terlihat "aneh" | Seed pakai `CURRENT_DATE + hari_ke` — disengaja, supaya agenda selalu relevan. |
| Ingin render ulang diagram | `npm i @mermaid-js/mermaid-cli` dulu. PNG sudah ter-commit, jadi **tidak wajib**. |

## Dua hal kecil yang belum dibersihkan (opsional, tapi rapi sebelum presentasi)

1. **Typo di seed:** deskripsi menu Sayur Sop tertulis *"Sayur bening bening khas dapur Dina
   Catering."* (`db/seed.sql:40`) — seharusnya "Sayur bening khas…". Tidak memengaruhi fungsi,
   tapi kalau dosen membuka detail menu, lebih baik sudah rapi.
2. **Karakter asing di `cara_ubah_database_manual.txt`** (di luar folder repo) — di baris 20
   dan 74 masih ada sisa draft berbahasa Mandarin (`，这才是…` dan `—改 langsung…`) yang
   menyisip di tengah kalimat Indonesia. Tidak memengaruhi kode, tapi file itu berada di
   folder yang sama dengan project dan bisa kebaca dosen.

---

# BAGIAN 7 — LAMPIRAN

## Materi visual siap pakai

| Kebutuhan | File |
|---|---|
| Alur proses manual (situasi lama) | `Diagram,dll/PNG/01-proses-usulan-manual.png` |
| Alur proses usulan (sistem baru) | `Diagram,dll/PNG/02-proses-usulan-sistem-baru.png` |
| Diagram perhitungan HPP | `Diagram,dll/PNG/03-perhitungan-hpp.png` |
| Diagram pengadaan bahan & stok | `Diagram,dll/PNG/04-pengadaan-stok.png` |
| Diagram pelaporan keuangan | `Diagram,dll/PNG/05-pelaporan-keuangan.png` |
| DFD Level 0 (context diagram) | `Diagram,dll/PNG/06-dfd-level-0.png` |
| DFD Level 1 (6 proses, 9 data store) | `Diagram,dll/PNG/07-dfd-level-1.png` |
| Use case diagram (4 aktor, 18 UC) | `Diagram,dll/PNG/08-use-case.png` |
| ERD 14 tabel | `Diagram,dll/PNG/09-erd.png` |
| Tampilkan semua diagram sekaligus | `Diagram,dll/viewer.html` (butuh internet) |
| Editable | `Diagram,dll/mermaid/*.mmd` (9 sumber) |
| Wireframe 14 layar | `Diagram,dll/prototipe/` (buka langsung di browser, tanpa server/DB) |
| Screenshot aplikasi | `docs/screenshots/` (7 PNG) |

## Dokumen pendukung

| Berkas | Isi |
|---|---|
| `README.md` | Dokumentasi produk lengkap (8 bab) |
| `docs/DOCKER.md` | Panduan Docker untuk pemula (9 bagian) |
| `Diagram,dll/01-Diagram-Proses-Usulan.md` | Alur manual vs usulan + alur HPP, stok, pelaporan |
| `Diagram,dll/02-DFD-Level-0-dan-Level-1.md` | Context diagram, DFD L1, kamus data |
| `Diagram,dll/03-Diagram-Use-Case.md` | 4 aktor, 18 use case, deskripsi, matriks aktor |
| `Diagram,dll/04-ERD.md` | ERD 14 tabel + kamus data + daftar index |
| `Diagram,dll/05-Rancangan-Antarmuka.md` | Peta navigasi + rancangan tiap layar |
| `../cara_ubah_database_manual.txt` | 3 cara mengubah database (UI / psql / pgAdmin) + tabel mana yang aman diedit manual |

## Angka ringkas untuk hafal

```
14 tabel · 49 rute · 13 halaman · 16 partial · 67 query sqlc
40 file Go · 6.206 baris Go · 15 fungsi test
7 status pesanan · 4 sumber pesanan · 3 status tugas produksi
margin minimum 30% · biaya operasional default 10% · DP 50%
tenggat DP H-30 · pelunasan H-7 · session TTL 12 jam
dependency: 4 (chi, pgx, golang-migrate, x/crypto)
```

## Urutan demo untuk dihafal (kalau nervous)

```
P1  login → dashboard → banner margin merah → pelanggan → pengaturan
P2  bahan baku → resep → paket → kalkulator HPP → harga daging 150.000 → HPP 33.275
P3  pembelian ayam 20 kg → pesanan Sabtu 60 tamu → total 2.280.000 → DP → 4 tugas → papan dapur
P4  laporan → keuntungan → margin per pesanan → biaya → smoke test
```

> Angka kunci yang **harus** keluar di mulut: **30.800 → 33.275**, **23% → 16,8%**,
> **2.280.000**, **H-30 / H-7**, **4 tugas produksi**, **15 kg ayam terpakai**.

---

*Dokumen ini dibuat sebagai panduan presentasi. Kalau ada perubahan kode setelahnya,
perbarui Bagian 3 (angka acuan) dan Bagian 4 (peta kode) — dua bagian itu yang paling
 cepat basi.*
