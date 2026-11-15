-- Data awal SIM Dina Catering.
-- Dijalankan hanya saat tabel users masih kosong.
-- Placeholder {PASSWORD_HASH} diganti hash bcrypt saat runtime.

-- =========================================================
-- Pengguna & pengaturan
-- =========================================================
INSERT INTO users (username, password_hash, role)
VALUES ('admin', '{PASSWORD_HASH}', 'admin');

INSERT INTO pengaturan (kunci, nilai) VALUES
    ('biaya_operasional_persen', '10');

-- =========================================================
-- Bahan baku (harga Rupiah per satuan, tahun 2026)
-- =========================================================
INSERT INTO bahan_baku (nama, satuan, stok_saat_ini, harga_terakhir) VALUES
    ('Beras Premium',   'kg',    40, 16000),
    ('Daging Sapi',      'kg',    15, 135000),
    ('Ayam Potong',     'kg',    30, 38000),
    ('Telur Ayam',      'kg',    12, 28000),
    ('Sayur Sop',       'kg',    25, 8000),
    ('Minyak Goreng',   'liter', 20, 18000),
    ('Bumbu Racik',     'kg',    8,  12000),
    ('Gas LPG 3 kg',    'kg',    14, 22000),
    ('Buah Potong',     'kg',    18, 15000);

-- Riwayat harga awal: berlaku sejak 1 Mei 2026.
INSERT INTO harga_bahan_history (bahan_baku_id, harga_satuan, tanggal_berlaku)
SELECT b.id, b.harga_terakhir, DATE '2026-05-01'
FROM bahan_baku b;

-- =========================================================
-- Menu & resep (takaran per porsi)
-- =========================================================
INSERT INTO menu (nama, deskripsi) VALUES
    ('Nasi Putih',     'Nasi putih steamer yang selalu tersedia di setiap paket.'),
    ('Ayam Bakar',     'Ayam potong bakar kecap dengan sambal doorstep.'),
    ('Rendang Daging', 'Rendang daging sapi dengan santan dan bumbu racik.'),
    ('Sayur Sop',      'Sayur bening bening khas dapur Dina Catering.'),
    ('Buah Potong',    'Potongan buah segar sebagai penutup jamuan.');

INSERT INTO resep_item (menu_id, bahan_baku_id, jumlah_per_porsi, satuan)
SELECT m.id, b.id, r.jumlah, r.satuan
FROM (VALUES
    ('Nasi Putih',     'Beras Premium',  0.150, 'kg'),
    ('Nasi Putih',     'Minyak Goreng',  0.015, 'liter'),
    ('Ayam Bakar',     'Ayam Potong',    0.250, 'kg'),
    ('Ayam Bakar',     'Bumbu Racik',    0.040, 'kg'),
    ('Ayam Bakar',     'Minyak Goreng',  0.015, 'liter'),
    ('Rendang Daging', 'Daging Sapi',    0.150, 'kg'),
    ('Rendang Daging', 'Bumbu Racik',    0.030, 'kg'),
    ('Rendang Daging', 'Minyak Goreng',  0.020, 'liter'),
    ('Sayur Sop',      'Sayur Sop',      0.200, 'kg'),
    ('Sayur Sop',      'Bumbu Racik',    0.020, 'kg'),
    ('Sayur Sop',      'Minyak Goreng',  0.015, 'liter'),
    ('Buah Potong',    'Buah Potong',    0.150, 'kg')
) AS r(menu_nama, bahan_nama, jumlah, satuan)
JOIN menu m ON m.nama = r.menu_nama
JOIN bahan_baku b ON b.nama = r.bahan_nama;

-- =========================================================
-- Paket menu (tier harga per tahun)
-- =========================================================
INSERT INTO paket_menu (nama, harga_per_porsi, harga_per_porsi_weekend, tahun_berlaku, deskripsi) VALUES
    ('Hemat',  30000, 32000, 2026, 'Nasi, Sayur Sop, dan Ayam Bakar.'),
    ('Standar', 35000, 38000, 2026, 'Nasi, Sayur Sop, Ayam Bakar, dan Buah Potong.'),
    ('Premium', 40000, 43000, 2026, 'Nasi, Rendang Daging, Sayur Sop, dan Buah Potong.');

INSERT INTO paket_menu_item (paket_menu_id, menu_id)
SELECT pk.id, m.id
FROM (VALUES
    ('Hemat',   'Nasi Putih'),
    ('Hemat',   'Sayur Sop'),
    ('Hemat',   'Ayam Bakar'),
    ('Standar', 'Nasi Putih'),
    ('Standar', 'Sayur Sop'),
    ('Standar', 'Ayam Bakar'),
    ('Standar', 'Buah Potong'),
    ('Premium', 'Nasi Putih'),
    ('Premium', 'Rendang Daging'),
    ('Premium', 'Sayur Sop'),
    ('Premium', 'Buah Potong')
) AS p(paket_nama, menu_nama)
JOIN paket_menu pk ON pk.nama = p.paket_nama AND pk.tahun_berlaku = 2026
JOIN menu m ON m.nama = p.menu_nama;

-- =========================================================
-- Pelanggan
-- =========================================================
INSERT INTO pelanggan (nama, kontak, alamat, catatan) VALUES
    ('Ratna Wulandari', '0812-3345-6781', 'Jl. Sudirman No. 12, Ungaran',      'Pemesanan lewat Instagram, biasa ambil paket Standar.'),
    ('Budi Santoso',   '0813-2211-4455', 'Jl. Ahmad Yani No. 45, Ungaran',    'Kenalan dari tetangga, sering pesan untuk acara keluarga.'),
    ('Siti Aminah',    '0857-9988-7766', 'Villa Surya Bakti, Tembalang',       'Paket Premium untuk acara pernikahan putrinya.'),
    ('Agus Prasetyo',   '0819-5544-3322', 'Hotel Sindang, Jl.KESsauf No. 3',    'Pasokan hotel rutin dengan tagihan bulanan.'),
    ('Dewi Lestari',   '0822-7788-9900', 'Jl. Diponegoro No. 8, Ungaran',    'Pelanggan lama, sudah beberapa kali pesan.'),
    ('Rizky Maulana',  '0815-6677-1122', 'Jl. Sudirman No. 45, Ungaran',    'Pesanan pernah dibatalkan karena berubah rencana.');

-- =========================================================
-- Pesanan contoh
-- hari_ke dihitung relatif terhadap CURRENT_DATE agar demo selalu relevan.
-- =========================================================
WITH pesanan_seed(pelanggan_nama, paket_nama, jumlah_tamu, hari_ke, lokasi, sumber, status, catatan, masak_di_lokasi) AS (
    VALUES
        ('Ratna Wulandari', 'Standar',  80,  25, 'Gedung Serbaguna Ungaran', 'instagram',      'baru',          'Prasmanan, minta sambal terpisah.', false),
        ('Budi Santoso',   'Hemat',   50,  40, 'Rumah Pak Budi, Ungaran',      'mulut_ke_mulut', 'final_meeting', 'DP sudah masuk, pelunasan H-7.', false),
        ('Siti Aminah',    'Premium', 120,  18, 'Villa Surya Bakti, Tembalang', 'instagram',      'survei_lokasi', 'Pernikahan putri, makanan dimasak di lokasi.', true),
        ('Agus Prasetyo',  'Standar', 150,  5, 'Hotel Sindang, Ungaran',       'hotel',          'produksi',      'Pasokan hotel, tempo pembayaran 30 hari.', false),
        ('Dewi Lestari',   'Premium', 200, -10, 'Gedung Serbaguna Ungaran',     'instagram',      'selesai',       'Selesai lunas, terima kasih banyak.', false),
        ('Rizky Maulana',  'Hemat',   60,  60, 'Rumah Rizky, Ungaran',         'lainnya',        'batal',         'Batal karena berubah rencana.', false)
)
INSERT INTO pesanan (pelanggan_id, paket_menu_id, jumlah_tamu, tanggal_acara, lokasi,
                     sumber_pesanan, total_harga, status, catatan, masak_di_lokasi,
                     tgl_jatuh_tempo_dp, tgl_jatuh_tempo_lunas)
SELECT pl.id,
       pk.id,
       s.jumlah_tamu,
       (CURRENT_DATE + s.hari_ke)::date,
       s.lokasi,
       s.sumber,
       s.jumlah_tamu * CASE
           WHEN EXTRACT(ISODOW FROM (CURRENT_DATE + s.hari_ke)::date) IN (6, 7)
               THEN COALESCE(pk.harga_per_porsi_weekend, pk.harga_per_porsi)
           ELSE pk.harga_per_porsi
       END,
       s.status,
       s.catatan,
       s.masak_di_lokasi,
       (CURRENT_DATE + s.hari_ke - 30)::date,
       (CURRENT_DATE + s.hari_ke - 7)::date
FROM pesanan_seed s
JOIN pelanggan pl ON pl.nama = s.pelanggan_nama
JOIN LATERAL (
    SELECT pm.* FROM paket_menu pm
    WHERE pm.nama = s.paket_nama
      AND pm.tahun_berlaku <= EXTRACT(YEAR FROM (CURRENT_DATE + s.hari_ke)::date)::int
    ORDER BY pm.tahun_berlaku DESC
    LIMIT 1
) pk ON true;

-- Pembayaran: DP 50% untuk pesanan yang sudah lewat final meeting / produksi.
INSERT INTO pembayaran (pesanan_id, jumlah, tipe, tanggal, metode)
SELECT p.id, p.total_harga / 2, 'dp', p.tgl_jatuh_tempo_dp, 'transfer'
FROM pesanan p
JOIN pelanggan pl ON pl.id = p.pelanggan_id
WHERE pl.nama IN ('Budi Santoso', 'Siti Aminah', 'Agus Prasetyo');

-- Pelunasan untuk pesanan yang sudah selesai.
INSERT INTO pembayaran (pesanan_id, jumlah, tipe, tanggal, metode)
SELECT p.id, p.total_harga - (p.total_harga / 2), 'pelunasan', p.tgl_jatuh_tempo_lunas, 'transfer'
FROM pesanan p
JOIN pelanggan pl ON pl.id = p.pelanggan_id
WHERE pl.nama = 'Dewi Lestari';

-- =========================================================
-- Jadwal produksi untuk pesanan yang sudah masuk produksi
-- =========================================================
INSERT INTO tugas_produksi (pesanan_id, nama_tugas, pj, deadline, status)
SELECT p.id, t.nama_tugas, t.pj, (p.tanggal_acara + t.hari_ke)::date, t.status
FROM pesanan p
JOIN pelanggan pl ON pl.id = p.pelanggan_id
CROSS JOIN (VALUES
    ('Persiapan bahan & alat', 'Sari',  -1, 'selesai'),
    ('Masak seluruh menu',     'Dedi',   0, 'jalan'),
    ('Plating & cek porsi',    'Sari',   0, 'belum'),
    ('Pengiriman ke lokasi',   'Rizky',  0, 'belum')
) AS t(nama_tugas, pj, hari_ke, status)
WHERE pl.nama = 'Agus Prasetyo';

-- =========================================================
-- Pembelian bahan di pasar tradisional
-- Tanggal tersebar di dua bulan supaya laporan laba rugi bulanan punya
-- pembelanjaan bahan di bulan berjalan maupun bulan sebelumnya.
-- =========================================================
INSERT INTO pembelian_bahan (bahan_baku_id, jumlah, harga_satuan, total, pasar, tanggal)
SELECT b.id, p.jumlah, p.harga_satuan, (p.jumlah * p.harga_satuan)::bigint, p.pasar,
       (CURRENT_DATE - p.hari_kemudian)::date
FROM (VALUES
    ('Beras Premium',  25, 16000, 'Pasar Ungaran',        5),
    ('Ayam Potong',   20, 38000, 'Pasar Ungaran',       35),
    ('Sayur Sop',     18,  8000, 'Pasar Sultan',          2),
    ('Minyak Goreng', 12, 18000, 'Pasar Ungaran',       32),
    ('Buah Potong',   15, 15000, 'Pasar Sultan',          1),
    ('Daging Sapi',   10,135000, 'Pasar Hewan Ungaran',  30)
) AS p(bahan_nama, jumlah, harga_satuan, pasar, hari_kemudian)
JOIN bahan_baku b ON b.nama = p.bahan_nama;

-- =========================================================
-- Biaya operasional
-- =========================================================
INSERT INTO biaya_operasional (kategori, deskripsi, jumlah, tanggal)
SELECT b.kategori, b.deskripsi, b.jumlah, (CURRENT_DATE - b.hari_kemudian)::date
FROM (VALUES
    ('Mobilitas',   'Bensin motor pengantar',            75000, 2),
    ('Mobilitas',   'Sewa kendaraan angkut pesanan',    350000, 3),
    ('Gas',         'Isi ulang tabung LPG 3 kg',          44000, 1),
    ('Listrik',     'Token listrik dapur',               150000, 4),
    ('Perlengkapan','Peralatan dapur dan kemasan',       275000, 6)
) AS b(kategori, deskripsi, jumlah, hari_kemudian);
