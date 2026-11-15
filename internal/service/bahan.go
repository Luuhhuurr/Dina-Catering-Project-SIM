package service

import (
	"context"
	"fmt"
	"time"

	"dinacatering/internal/store"
)

// CatatHargaBahan menambah baris riwayat harga dan memperbarui harga
// terakhir bahan. Inilah yang membuat HPP paket langsung ikut berubah ketika
// harga bahan naik di pasar.
func (s *Service) CatatHargaBahan(ctx context.Context, bahanID int32, harga int64, tanggalBerlaku time.Time) error {
	return s.withTx(ctx, func(q store.Querier) error {
		if _, err := q.InsertHargaBahan(ctx, store.InsertHargaBahanParams{
			BahanBakuID:    bahanID,
			HargaSatuan:    harga,
			TanggalBerlaku: TanggalNormalize(tanggalBerlaku),
		}); err != nil {
			return fmt.Errorf("simpan riwayat harga: %w", err)
		}
		if err := q.UpdateHargaTerakhirBahanDariRiwayat(ctx, bahanID); err != nil {
			return fmt.Errorf("update harga terakhir: %w", err)
		}
		return nil
	})
}

// CatatPembelian menyimpan pembelian bahan di pasar: stok bertambah, harga
// terakhir diperbarui, dan riwayat harga bertambah satu baris.
func (s *Service) CatatPembelian(ctx context.Context, bahanID int32, jumlah float64, hargaSatuan int64, pasar string, tanggal time.Time) error {
	if jumlah <= 0 {
		return fmt.Errorf("jumlah pembelian harus lebih dari nol")
	}
	total := int64(jumlah * float64(hargaSatuan))

	return s.withTx(ctx, func(q store.Querier) error {
		bahan, err := q.GetBahanBaku(ctx, bahanID)
		if err != nil {
			return fmt.Errorf("ambil bahan: %w", err)
		}

		if _, err := q.CreatePembelian(ctx, store.CreatePembelianParams{
			BahanBakuID: bahanID,
			Jumlah:      jumlah,
			HargaSatuan: hargaSatuan,
			Total:       total,
			Pasar:       pasar,
			Tanggal:     TanggalNormalize(tanggal),
		}); err != nil {
			return fmt.Errorf("simpan pembelian: %w", err)
		}

		stokBaru := bahan.StokSaatIni + jumlah
		if err := q.SetStokBahan(ctx, store.SetStokBahanParams{
			ID:   bahanID,
			Stok: stokBaru,
		}); err != nil {
			return fmt.Errorf("update stok bahan: %w", err)
		}

		if _, err := q.InsertHargaBahan(ctx, store.InsertHargaBahanParams{
			BahanBakuID:    bahanID,
			HargaSatuan:    hargaSatuan,
			TanggalBerlaku: TanggalNormalize(tanggal),
		}); err != nil {
			return fmt.Errorf("simpan riwayat harga: %w", err)
		}

		if err := q.UpdateHargaTerakhirBahanDariRiwayat(ctx, bahanID); err != nil {
			return fmt.Errorf("update harga terakhir: %w", err)
		}
		return nil
	})
}
