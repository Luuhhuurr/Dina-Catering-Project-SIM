-- SIM Dina Catering — skema database
-- Catatan: uang (harga/total) memakai BIGINT (Rupiah bulat), kuantitas DOUBLE PRECISION.

CREATE TABLE users (
    id            SERIAL PRIMARY KEY,
    username      TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role          TEXT NOT NULL DEFAULT 'admin',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Konfigurasi runtime yang bisa diubah owner (persentase biaya operasional, dll)
CREATE TABLE pengaturan (
    kunci      TEXT PRIMARY KEY,
    nilai      TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE pelanggan (
    id         SERIAL PRIMARY KEY,
    nama       TEXT NOT NULL,
    kontak     TEXT NOT NULL DEFAULT '',
    alamat     TEXT NOT NULL DEFAULT '',
    catatan    TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE bahan_baku (
    id             SERIAL PRIMARY KEY,
    nama           TEXT NOT NULL UNIQUE,
    satuan         TEXT NOT NULL DEFAULT 'kg',
    stok_saat_ini  DOUBLE PRECISION NOT NULL DEFAULT 0,
    harga_terakhir BIGINT NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Riwayat harga bahan: HPP historis dihitung dari harga yang berlaku pada tanggal tertentu
CREATE TABLE harga_bahan_history (
    id             SERIAL PRIMARY KEY,
    bahan_baku_id  INTEGER NOT NULL REFERENCES bahan_baku(id) ON DELETE CASCADE,
    harga_satuan   BIGINT NOT NULL,
    tanggal_berlaku DATE NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_harga_bahan_riwayat ON harga_bahan_history (bahan_baku_id, tanggal_berlaku DESC);

CREATE TABLE pembelian_bahan (
    id            SERIAL PRIMARY KEY,
    bahan_baku_id INTEGER NOT NULL REFERENCES bahan_baku(id),
    jumlah        DOUBLE PRECISION NOT NULL,
    harga_satuan  BIGINT NOT NULL,
    total         BIGINT NOT NULL,
    pasar         TEXT NOT NULL DEFAULT '',
    tanggal       DATE NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_pembelian_tanggal ON pembelian_bahan (tanggal);

CREATE TABLE menu (
    id          SERIAL PRIMARY KEY,
    nama        TEXT NOT NULL UNIQUE,
    deskripsi   TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE resep_item (
    id               SERIAL PRIMARY KEY,
    menu_id          INTEGER NOT NULL REFERENCES menu(id) ON DELETE CASCADE,
    bahan_baku_id    INTEGER NOT NULL REFERENCES bahan_baku(id),
    jumlah_per_porsi DOUBLE PRECISION NOT NULL,
    satuan           TEXT NOT NULL DEFAULT ''
);
CREATE INDEX idx_resep_menu ON resep_item (menu_id);

CREATE TABLE paket_menu (
    id                       SERIAL PRIMARY KEY,
    nama                     TEXT NOT NULL,
    harga_per_porsi          BIGINT NOT NULL,
    -- NULL = harga akhir pekan sama dengan harga biasa
    harga_per_porsi_weekend  BIGINT,
    tahun_berlaku            INTEGER NOT NULL,
    deskripsi                TEXT NOT NULL DEFAULT '',
    created_at               TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (nama, tahun_berlaku)
);

CREATE TABLE paket_menu_item (
    id             SERIAL PRIMARY KEY,
    paket_menu_id  INTEGER NOT NULL REFERENCES paket_menu(id) ON DELETE CASCADE,
    menu_id        INTEGER NOT NULL REFERENCES menu(id) ON DELETE CASCADE
);
CREATE INDEX idx_paket_item_paket ON paket_menu_item (paket_menu_id);

CREATE TABLE pesanan (
    id                    SERIAL PRIMARY KEY,
    pelanggan_id          INTEGER NOT NULL REFERENCES pelanggan(id),
    paket_menu_id         INTEGER NOT NULL REFERENCES paket_menu(id),
    jumlah_tamu           INTEGER NOT NULL,
    tanggal_acara         DATE NOT NULL,
    lokasi                TEXT NOT NULL DEFAULT '',
    sumber_pesanan        TEXT NOT NULL CHECK (sumber_pesanan IN ('instagram', 'mulut_ke_mulut', 'hotel', 'lainnya')),
    total_harga           BIGINT NOT NULL DEFAULT 0,
    status                TEXT NOT NULL DEFAULT 'baru' CHECK (status IN ('baru', 'final_meeting', 'survei_lokasi', 'produksi', 'pengiriman', 'selesai', 'batal')),
    catatan               TEXT NOT NULL DEFAULT '',
    masak_di_lokasi       BOOLEAN NOT NULL DEFAULT false,
    tgl_jatuh_tempo_dp    DATE NOT NULL,
    tgl_jatuh_tempo_lunas DATE NOT NULL,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_pesanan_tanggal ON pesanan (tanggal_acara);
CREATE INDEX idx_pesanan_status ON pesanan (status);

CREATE TABLE pembayaran (
    id         SERIAL PRIMARY KEY,
    pesanan_id INTEGER NOT NULL REFERENCES pesanan(id) ON DELETE CASCADE,
    jumlah     BIGINT NOT NULL,
    tipe       TEXT NOT NULL CHECK (tipe IN ('dp', 'pelunasan')),
    tanggal    DATE NOT NULL,
    metode     TEXT NOT NULL DEFAULT 'transfer',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_pembayaran_pesanan ON pembayaran (pesanan_id);

CREATE TABLE tugas_produksi (
    id         SERIAL PRIMARY KEY,
    pesanan_id INTEGER NOT NULL REFERENCES pesanan(id) ON DELETE CASCADE,
    nama_tugas TEXT NOT NULL,
    pj         TEXT NOT NULL DEFAULT '',
    deadline   DATE,
    status     TEXT NOT NULL DEFAULT 'belum' CHECK (status IN ('belum', 'jalan', 'selesai')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_tugas_pesanan ON tugas_produksi (pesanan_id);
CREATE INDEX idx_tugas_deadline ON tugas_produksi (deadline);

CREATE TABLE biaya_operasional (
    id         SERIAL PRIMARY KEY,
    kategori   TEXT NOT NULL,
    deskripsi  TEXT NOT NULL DEFAULT '',
    jumlah     BIGINT NOT NULL,
    tanggal    DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_biaya_tanggal ON biaya_operasional (tanggal);
