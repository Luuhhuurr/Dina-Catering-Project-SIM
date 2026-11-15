-- SIM Dina Catering — query sqlc
-- Konvensi parameter filter:
--   status/s sourced '' berarti "semua"
--   dari/sampai memakai tanggal batas (1900-01-01 / 2999-12-31)

-- =========================================================
-- Auth
-- =========================================================

-- name: GetUserByUsername :one
SELECT id, username, password_hash, role, created_at
FROM users
WHERE username = $1;

-- name: CreateUser :exec
INSERT INTO users (username, password_hash, role)
VALUES ($1, $2, $3);

-- =========================================================
-- Pengaturan
-- =========================================================

-- name: GetPengaturan :one
SELECT kunci, nilai, updated_at
FROM pengaturan
WHERE kunci = $1;

-- name: ListPengaturan :many
SELECT kunci, nilai, updated_at
FROM pengaturan
ORDER BY kunci;

-- name: UpsertPengaturan :exec
INSERT INTO pengaturan (kunci, nilai, updated_at)
VALUES ($1, $2, now())
ON CONFLICT (kunci) DO UPDATE SET nilai = EXCLUDED.nilai, updated_at = now();

-- =========================================================
-- Pelanggan
-- =========================================================

-- name: ListPelanggan :many
SELECT p.id, p.nama, p.kontak, p.alamat, p.catatan, p.created_at,
       COALESCE(c.jumlah_pesanan, 0)::bigint AS jumlah_pesanan
FROM pelanggan p
LEFT JOIN (
    SELECT pelanggan_id, COUNT(*) AS jumlah_pesanan FROM pesanan GROUP BY pelanggan_id
) c ON c.pelanggan_id = p.id
WHERE sqlc.narg('cari')::text IS NULL
   OR p.nama ILIKE '%' || sqlc.narg('cari')::text || '%'
   OR p.kontak ILIKE '%' || sqlc.narg('cari')::text || '%'
ORDER BY p.nama;

-- name: GetPelanggan :one
SELECT id, nama, kontak, alamat, catatan, created_at
FROM pelanggan
WHERE id = $1;

-- name: CreatePelanggan :one
INSERT INTO pelanggan (nama, kontak, alamat, catatan)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: UpdatePelanggan :exec
UPDATE pelanggan
SET nama = $2, kontak = $3, alamat = $4, catatan = $5
WHERE id = $1;

-- name: DeletePelanggan :exec
DELETE FROM pelanggan WHERE id = $1;

-- name: ListPesananByPelanggan :many
SELECT p.id, p.pelanggan_id, p.paket_menu_id, p.jumlah_tamu, p.tanggal_acara, p.lokasi,
       p.sumber_pesanan, p.total_harga, p.status, p.catatan, p.masak_di_lokasi,
       p.tgl_jatuh_tempo_dp, p.tgl_jatuh_tempo_lunas, p.created_at,
       pk.nama AS paket_nama
FROM pesanan p
JOIN paket_menu pk ON pk.id = p.paket_menu_id
WHERE p.pelanggan_id = $1
ORDER BY p.tanggal_acara DESC;

-- =========================================================
-- Bahan baku
-- =========================================================

-- name: ListBahanBaku :many
SELECT id, nama, satuan, stok_saat_ini, harga_terakhir, created_at
FROM bahan_baku
ORDER BY nama;

-- name: GetBahanBaku :one
SELECT id, nama, satuan, stok_saat_ini, harga_terakhir, created_at
FROM bahan_baku
WHERE id = $1;

-- name: CreateBahanBaku :one
INSERT INTO bahan_baku (nama, satuan, stok_saat_ini, harga_terakhir)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: UpdateBahanBaku :exec
UPDATE bahan_baku
SET nama = $2, satuan = $3, stok_saat_ini = $4, harga_terakhir = $5
WHERE id = $1;

-- name: DeleteBahanBaku :exec
DELETE FROM bahan_baku WHERE id = $1;

-- name: SetStokBahan :exec
UPDATE bahan_baku SET stok_saat_ini = sqlc.arg('stok') WHERE id = $1;

-- name: KurangiStokBahan :exec
UPDATE bahan_baku SET stok_saat_ini = stok_saat_ini - sqlc.arg('jumlah') WHERE id = $1;

-- =========================================================
-- Riwayat harga bahan
-- =========================================================

-- name: InsertHargaBahan :one
INSERT INTO harga_bahan_history (bahan_baku_id, harga_satuan, tanggal_berlaku)
VALUES ($1, $2, $3)
RETURNING id;

-- name: UpdateHargaTerakhirBahanDariRiwayat :exec
-- harga_terakhir selalu mengikuti baris riwayat terbaru (tanggal berlaku
-- paling akhir). Jadi mencatat harga atau pembelian dengan tanggal lampau
-- tidak menurunkan harga yang sedang berlaku.
UPDATE bahan_baku b
SET harga_terakhir = (
    SELECT h.harga_satuan
    FROM harga_bahan_history h
    WHERE h.bahan_baku_id = b.id
    ORDER BY h.tanggal_berlaku DESC, h.id DESC
    LIMIT 1
)
WHERE b.id = $1;

-- name: ListHargaBahanByBahan :many
SELECT id, bahan_baku_id, harga_satuan, tanggal_berlaku, created_at
FROM harga_bahan_history
WHERE bahan_baku_id = $1
ORDER BY tanggal_berlaku DESC, id DESC;

-- name: GetHargaBahanEfektif :one
SELECT id, bahan_baku_id, harga_satuan, tanggal_berlaku, created_at
FROM harga_bahan_history
WHERE bahan_baku_id = $1 AND tanggal_berlaku <= $2
ORDER BY tanggal_berlaku DESC, id DESC
LIMIT 1;

-- name: ListHargaEfektif :many
-- Harga terbaru yang berlaku pada atau sebelum tanggal yang diminta, untuk semua bahan.
SELECT hbh.bahan_baku_id, hbh.harga_satuan
FROM harga_bahan_history hbh
WHERE hbh.tanggal_berlaku <= sqlc.arg('tanggal')::date
  AND hbh.id = (
    SELECT h2.id
    FROM harga_bahan_history h2
    WHERE h2.bahan_baku_id = hbh.bahan_baku_id
      AND h2.tanggal_berlaku <= sqlc.arg('tanggal')::date
    ORDER BY h2.tanggal_berlaku DESC, h2.id DESC
    LIMIT 1
  );

-- =========================================================
-- Pembelian bahan
-- =========================================================

-- name: CreatePembelian :one
INSERT INTO pembelian_bahan (bahan_baku_id, jumlah, harga_satuan, total, pasar, tanggal)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: ListPembelian :many
SELECT pb.id, pb.bahan_baku_id, bb.nama AS bahan_nama, bb.satuan, pb.jumlah, pb.harga_satuan,
       pb.total, pb.pasar, pb.tanggal, pb.created_at
FROM pembelian_bahan pb
JOIN bahan_baku bb ON bb.id = pb.bahan_baku_id
WHERE pb.tanggal >= sqlc.arg('dari')::date AND pb.tanggal <= sqlc.arg('sampai')::date
ORDER BY pb.tanggal DESC, pb.id DESC;

-- name: ListPembelianByBahan :many
SELECT pb.id, pb.bahan_baku_id, bb.nama AS bahan_nama, bb.satuan, pb.jumlah, pb.harga_satuan,
       pb.total, pb.pasar, pb.tanggal, pb.created_at
FROM pembelian_bahan pb
JOIN bahan_baku bb ON bb.id = pb.bahan_baku_id
WHERE pb.bahan_baku_id = $1
ORDER BY pb.tanggal DESC, pb.id DESC;

-- name: TotalPembelianRentang :one
SELECT COALESCE(SUM(total), 0)::bigint AS total
FROM pembelian_bahan
WHERE tanggal >= sqlc.arg('dari')::date AND tanggal <= sqlc.arg('sampai')::date;

-- =========================================================
-- Menu & resep
-- =========================================================

-- name: ListMenu :many
SELECT m.id, m.nama, m.deskripsi, m.created_at,
       COALESCE(r.jumlah_resep, 0)::bigint AS jumlah_resep
FROM menu m
LEFT JOIN (
    SELECT menu_id, COUNT(*) AS jumlah_resep FROM resep_item GROUP BY menu_id
) r ON r.menu_id = m.id
ORDER BY m.nama;

-- name: GetMenu :one
SELECT id, nama, deskripsi, created_at FROM menu WHERE id = $1;

-- name: CreateMenu :one
INSERT INTO menu (nama, deskripsi) VALUES ($1, $2) RETURNING id;

-- name: UpdateMenu :exec
UPDATE menu SET nama = $2, deskripsi = $3 WHERE id = $1;

-- name: DeleteMenu :exec
DELETE FROM menu WHERE id = $1;

-- name: ListResepByMenu :many
SELECT ri.id, ri.menu_id, ri.bahan_baku_id, ri.jumlah_per_porsi, ri.satuan,
       bb.nama AS bahan_nama, bb.satuan AS bahan_satuan
FROM resep_item ri
JOIN bahan_baku bb ON bb.id = ri.bahan_baku_id
WHERE ri.menu_id = $1
ORDER BY bb.nama;

-- name: InsertResepItem :one
INSERT INTO resep_item (menu_id, bahan_baku_id, jumlah_per_porsi, satuan)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: DeleteResepItem :exec
DELETE FROM resep_item WHERE id = $1;

-- =========================================================
-- Paket menu
-- =========================================================

-- name: ListPaket :many
SELECT pm.id, pm.nama, pm.harga_per_porsi, pm.harga_per_porsi_weekend, pm.tahun_berlaku,
       pm.deskripsi, pm.created_at,
       COALESCE(pmi.jumlah_menu, 0)::bigint AS jumlah_menu
FROM paket_menu pm
LEFT JOIN (
    SELECT paket_menu_id, COUNT(*) AS jumlah_menu FROM paket_menu_item GROUP BY paket_menu_id
) pmi ON pmi.paket_menu_id = pm.id
ORDER BY pm.tahun_berlaku DESC, pm.harga_per_porsi;

-- name: ListPaketPerTahun :many
SELECT DISTINCT ON (nama) id, nama, harga_per_porsi, harga_per_porsi_weekend, tahun_berlaku,
       deskripsi, created_at
FROM paket_menu
WHERE tahun_berlaku <= sqlc.arg('tahun')::int
ORDER BY nama, tahun_berlaku DESC;

-- name: GetPaket :one
SELECT id, nama, harga_per_porsi, harga_per_porsi_weekend, tahun_berlaku, deskripsi, created_at
FROM paket_menu
WHERE id = $1;

-- name: GetPaketEfektif :one
-- Baris harga paket yang berlaku untuk tahun acara (tahun terdekat <= tahun acara).
SELECT id, nama, harga_per_porsi, harga_per_porsi_weekend, tahun_berlaku, deskripsi, created_at
FROM paket_menu
WHERE nama = $1 AND tahun_berlaku <= sqlc.arg('tahun')::int
ORDER BY tahun_berlaku DESC
LIMIT 1;

-- name: CreatePaket :one
INSERT INTO paket_menu (nama, harga_per_porsi, harga_per_porsi_weekend, tahun_berlaku, deskripsi)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: UpdatePaket :exec
UPDATE paket_menu
SET nama = $2, harga_per_porsi = $3, harga_per_porsi_weekend = $4, tahun_berlaku = $5, deskripsi = $6
WHERE id = $1;

-- name: DeletePaket :exec
DELETE FROM paket_menu WHERE id = $1;

-- name: ListPaketItem :many
SELECT pmi.id, pmi.paket_menu_id, pmi.menu_id, m.nama AS menu_nama, m.deskripsi AS menu_deskripsi
FROM paket_menu_item pmi
JOIN menu m ON m.id = pmi.menu_id
WHERE pmi.paket_menu_id = $1
ORDER BY m.nama;

-- name: InsertPaketItem :one
INSERT INTO paket_menu_item (paket_menu_id, menu_id) VALUES ($1, $2) RETURNING id;

-- name: DeletePaketItem :exec
DELETE FROM paket_menu_item WHERE id = $1;

-- =========================================================
-- Pesanan
-- =========================================================

-- name: CreatePesanan :one
INSERT INTO pesanan (pelanggan_id, paket_menu_id, jumlah_tamu, tanggal_acara, lokasi,
                     sumber_pesanan, total_harga, status, catatan, masak_di_lokasi,
                     tgl_jatuh_tempo_dp, tgl_jatuh_tempo_lunas)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING id;

-- name: GetPesanan :one
SELECT p.id, p.pelanggan_id, pl.nama AS pelanggan_nama, pl.kontak AS pelanggan_kontak,
       p.paket_menu_id, pk.nama AS paket_nama, p.jumlah_tamu, p.tanggal_acara, p.lokasi,
       p.sumber_pesanan, p.total_harga, p.status, p.catatan, p.masak_di_lokasi,
       p.tgl_jatuh_tempo_dp, p.tgl_jatuh_tempo_lunas, p.created_at,
       COALESCE(pay.total_bayar, 0)::bigint AS total_bayar
FROM pesanan p
JOIN pelanggan pl ON pl.id = p.pelanggan_id
JOIN paket_menu pk ON pk.id = p.paket_menu_id
LEFT JOIN (
    SELECT pesanan_id, SUM(jumlah) AS total_bayar FROM pembayaran GROUP BY pesanan_id
) pay ON pay.pesanan_id = p.id
WHERE p.id = $1;

-- name: ListPesanan :many
SELECT p.id, p.pelanggan_id, pl.nama AS pelanggan_nama, p.paket_menu_id, pk.nama AS paket_nama,
       p.jumlah_tamu, p.tanggal_acara, p.lokasi, p.sumber_pesanan, p.total_harga, p.status,
       p.catatan, p.masak_di_lokasi, p.tgl_jatuh_tempo_dp, p.tgl_jatuh_tempo_lunas, p.created_at,
       COALESCE(pay.total_bayar, 0)::bigint AS total_bayar
FROM pesanan p
JOIN pelanggan pl ON pl.id = p.pelanggan_id
JOIN paket_menu pk ON pk.id = p.paket_menu_id
LEFT JOIN (
    SELECT pesanan_id, SUM(jumlah) AS total_bayar FROM pembayaran GROUP BY pesanan_id
) pay ON pay.pesanan_id = p.id
WHERE (sqlc.narg('status')::text IS NULL OR p.status = sqlc.narg('status')::text)
  AND (sqlc.narg('sumber')::text IS NULL OR p.sumber_pesanan = sqlc.narg('sumber')::text)
  AND (p.tanggal_acara >= sqlc.arg('dari')::date)
  AND (p.tanggal_acara <= sqlc.arg('sampai')::date)
ORDER BY p.tanggal_acara DESC, p.id DESC;

-- name: UpdateStatusPesanan :exec
UPDATE pesanan SET status = $2 WHERE id = $1;

-- name: UpdatePesanan :exec
UPDATE pesanan
SET pelanggan_id = $2, paket_menu_id = $3, jumlah_tamu = $4, tanggal_acara = $5, lokasi = $6,
    sumber_pesanan = $7, total_harga = $8, catatan = $9, masak_di_lokasi = $10,
    tgl_jatuh_tempo_dp = $11, tgl_jatuh_tempo_lunas = $12
WHERE id = $1;

-- name: ListPesananJatuhTempoDP :many
SELECT p.id, p.pelanggan_id, pl.nama AS pelanggan_nama, p.paket_menu_id, pk.nama AS paket_nama,
       p.jumlah_tamu, p.tanggal_acara, p.lokasi, p.sumber_pesanan, p.total_harga, p.status,
       p.catatan, p.masak_di_lokasi, p.tgl_jatuh_tempo_dp, p.tgl_jatuh_tempo_lunas, p.created_at,
       COALESCE(pay.total_bayar, 0)::bigint AS total_bayar
FROM pesanan p
JOIN pelanggan pl ON pl.id = p.pelanggan_id
JOIN paket_menu pk ON pk.id = p.paket_menu_id
LEFT JOIN (
    SELECT pesanan_id, SUM(jumlah) AS total_bayar FROM pembayaran GROUP BY pesanan_id
) pay ON pay.pesanan_id = p.id
WHERE p.status NOT IN ('batal', 'selesai')
  AND p.tgl_jatuh_tempo_dp <= sqlc.arg('sampai')::date
  AND COALESCE(pay.total_bayar, 0) < p.total_harga
ORDER BY p.tgl_jatuh_tempo_dp;

-- name: ListPesananRentangTanggal :many
SELECT p.id, p.pelanggan_id, pl.nama AS pelanggan_nama, p.paket_menu_id, pk.nama AS paket_nama,
       p.jumlah_tamu, p.tanggal_acara, p.lokasi, p.sumber_pesanan, p.total_harga, p.status,
       p.catatan, p.masak_di_lokasi, p.tgl_jatuh_tempo_dp, p.tgl_jatuh_tempo_lunas, p.created_at,
       COALESCE(pay.total_bayar, 0)::bigint AS total_bayar
FROM pesanan p
JOIN pelanggan pl ON pl.id = p.pelanggan_id
JOIN paket_menu pk ON pk.id = p.paket_menu_id
LEFT JOIN (
    SELECT pesanan_id, SUM(jumlah) AS total_bayar FROM pembayaran GROUP BY pesanan_id
) pay ON pay.pesanan_id = p.id
WHERE p.tanggal_acara >= sqlc.arg('dari')::date
  AND p.tanggal_acara <= sqlc.arg('sampai')::date
  AND p.status <> 'batal'
ORDER BY p.tanggal_acara, p.jumlah_tamu DESC;

-- name: CountPesananByStatus :many
SELECT status, COUNT(*)::bigint AS jumlah
FROM pesanan
GROUP BY status
ORDER BY status;

-- name: CountPesananBySumber :many
SELECT sumber_pesanan, COUNT(*)::bigint AS jumlah
FROM pesanan
GROUP BY sumber_pesanan
ORDER BY jumlah DESC;

-- =========================================================
-- Pembayaran
-- =========================================================

-- name: CreatePembayaran :one
INSERT INTO pembayaran (pesanan_id, jumlah, tipe, tanggal, metode)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: ListPembayaranByPesanan :many
SELECT id, pesanan_id, jumlah, tipe, tanggal, metode, created_at
FROM pembayaran
WHERE pesanan_id = $1
ORDER BY tanggal DESC, id DESC;

-- name: TotalPembayaranPesanan :one
SELECT COALESCE(SUM(jumlah), 0)::bigint AS total
FROM pembayaran
WHERE pesanan_id = $1;

-- name: TotalPenerimaanRentang :one
SELECT COALESCE(SUM(pb.jumlah), 0)::bigint AS total
FROM pembayaran pb
JOIN pesanan p ON p.id = pb.pesanan_id
WHERE pb.tanggal >= sqlc.arg('dari')::date
  AND pb.tanggal <= sqlc.arg('sampai')::date
  AND p.status <> 'batal';

-- =========================================================
-- Tugas produksi
-- =========================================================

-- name: InsertTugas :one
INSERT INTO tugas_produksi (pesanan_id, nama_tugas, pj, deadline, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: ListTugasByPesanan :many
SELECT id, pesanan_id, nama_tugas, pj, deadline, status, created_at
FROM tugas_produksi
WHERE pesanan_id = $1
ORDER BY deadline NULLS LAST, id;

-- name: ListTugasRentangTanggal :many
SELECT t.id, t.pesanan_id, t.nama_tugas, t.pj, t.deadline, t.status, t.created_at,
       p.tanggal_acara, pl.nama AS pelanggan_nama, pk.nama AS paket_nama,
       p.jumlah_tamu, p.masak_di_lokasi
FROM tugas_produksi t
JOIN pesanan p ON p.id = t.pesanan_id
JOIN pelanggan pl ON pl.id = p.pelanggan_id
JOIN paket_menu pk ON pk.id = p.paket_menu_id
WHERE t.deadline >= sqlc.arg('dari')::date
  AND t.deadline <= sqlc.arg('sampai')::date
  AND p.status <> 'batal'
ORDER BY t.deadline, t.id;

-- name: UpdateTugasPJ :exec
UPDATE tugas_produksi SET pj = $2 WHERE id = $1;

-- name: UpdateTugasStatus :exec
UPDATE tugas_produksi SET status = $2 WHERE id = $1;

-- name: ListTugasBelumAdaPJ :many
SELECT id, pesanan_id, nama_tugas, pj, deadline, status, created_at
FROM tugas_produksi
WHERE pj = '' AND status <> 'selesai'
ORDER BY deadline NULLS LAST;

-- =========================================================
-- Biaya operasional
-- =========================================================

-- name: CreateBiaya :one
INSERT INTO biaya_operasional (kategori, deskripsi, jumlah, tanggal)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: ListBiaya :many
SELECT id, kategori, deskripsi, jumlah, tanggal, created_at
FROM biaya_operasional
WHERE tanggal >= sqlc.arg('dari')::date AND tanggal <= sqlc.arg('sampai')::date
ORDER BY tanggal DESC, id DESC;

-- name: TotalBiayaRentang :one
SELECT COALESCE(SUM(jumlah), 0)::bigint AS total
FROM biaya_operasional
WHERE tanggal >= sqlc.arg('dari')::date AND tanggal <= sqlc.arg('sampai')::date;
