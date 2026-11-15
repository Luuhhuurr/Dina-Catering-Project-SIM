// Package service berisi logika bisnis aplikasi: perhitungan HPP & margin,
// jadwal pembayaran, pengurangan stok, dan laporan keuangan.
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"dinacatering/internal/model"
	"dinacatering/internal/store"
)

// Service adalah pintu masuk seluruh logika bisnis.
//
// Queries mengekspos lapisan store hasil generate sqlc supaya handler bisa
// menjalankan query tanpa menulis SQL di dalam handler.
type Service struct {
	Queries store.Querier
	pool    *pgxpool.Pool
}

// New menyusun Service baru.
func New(queries store.Querier, pool *pgxpool.Pool) *Service {
	return &Service{Queries: queries, pool: pool}
}

// ErrPaketTidakDitemukan dipanggil saat baris harga paket tidak ada untuk tahun acara.
var ErrPaketTidakDitemukan = errors.New("harga paket untuk tahun acara belum tersedia")

// FindUserByUsername mengambil user untuk proses login.
func (s *Service) FindUserByUsername(ctx context.Context, username string) (store.User, error) {
	return s.Queries.GetUserByUsername(ctx, username)
}

// PersenBiayaOperasional membaca persentase biaya operasional dari tabel pengaturan.
func (s *Service) PersenBiayaOperasional(ctx context.Context) (float64, error) {
	pengaturan, err := s.Queries.GetPengaturan(ctx, model.KunciPengaturanBiayaOperasional)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.DefaultPersenBiayaOperasional, nil
	}
	if err != nil {
		return 0, fmt.Errorf("ambil pengaturan biaya operasional: %w", err)
	}
	var persen float64
	if _, err := fmt.Sscanf(pengaturan.Nilai, "%g", &persen); err != nil {
		return model.DefaultPersenBiayaOperasional, nil
	}
	return persen, nil
}

// SetPersenBiayaOperasional menyimpan persentase biaya operasional.
func (s *Service) SetPersenBiayaOperasional(ctx context.Context, persen float64) error {
	return s.Queries.UpsertPengaturan(ctx, store.UpsertPengaturanParams{
		Kunci: model.KunciPengaturanBiayaOperasional,
		Nilai: fmt.Sprintf("%g", persen),
	})
}

// di dalam transaksi, seluruh query diarahkan ke tx.
func (s *Service) withTx(ctx context.Context, fn func(q store.Querier) error) error {
	if s.pool == nil {
		return errors.New("transaksi butuh koneksi database")
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("mulai transaksi: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if err := fn(store.New(tx)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// pgDate mengubah tanggal biasa menjadi tipe pgtype.Date untuk kolom tanggal
// yang boleh kosong.
func pgDate(t time.Time) pgtype.Date {
	return pgtype.Date{Time: TanggalNormalize(t), Valid: true}
}
