package service

import (
	"context"
	"fmt"

	"dinacatering/internal/store"
)

// kurangiStokUntukProduksi mengurangi stok bahan sebesar resep x jumlah tamu.
//
// Asumsi satuan: takaran pada resep memakai satuan yang sama dengan satuan
// stok bahan (kg untuk beras, liter untuk minyak, dan seterusnya).
func (s *Service) kurangiStokUntukProduksi(ctx context.Context, q store.Querier, pesanan store.GetPesananRow) error {
	itemPaket, err := q.ListPaketItem(ctx, pesanan.PaketMenuID)
	if err != nil {
		return fmt.Errorf("ambil isi paket: %w", err)
	}

	for _, it := range itemPaket {
		resep, err := q.ListResepByMenu(ctx, it.MenuID)
		if err != nil {
			return fmt.Errorf("ambil resep menu %d: %w", it.MenuID, err)
		}
		for _, r := range resep {
			dimakan := r.JumlahPerPorsi * float64(pesanan.JumlahTamu)
			if dimakan <= 0 {
				continue
			}
			if err := q.KurangiStokBahan(ctx, store.KurangiStokBahanParams{
				ID:     r.BahanBakuID,
				Jumlah: dimakan,
			}); err != nil {
				return fmt.Errorf("kurangi stok bahan %d: %w", r.BahanBakuID, err)
			}
		}
	}
	return nil
}
