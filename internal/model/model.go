// Package model berisi konstanta domain yang dipakai bersama oleh service dan handler.
package model

// Status alur pesanan (aturan bisnis No.5).
const (
	StatusBaru         = "baru"
	StatusFinalMeeting = "final_meeting"
	StatusSurveiLokasi = "survei_lokasi"
	StatusProduksi     = "produksi"
	StatusPengiriman   = "pengiriman"
	StatusSelesai      = "selesai"
	StatusBatal        = "batal"
)

// StatusPesananList adalah urutan alur kerja pesanan.
var StatusPesananList = []string{
	StatusBaru, StatusFinalMeeting, StatusSurveiLokasi,
	StatusProduksi, StatusPengiriman, StatusSelesai, StatusBatal,
}

// LabelStatusPesanan memakai nama yang mudah dibaca owner.
var LabelStatusPesanan = map[string]string{
	StatusBaru:         "Baru",
	StatusFinalMeeting: "Final Meeting",
	StatusSurveiLokasi: "Survei Lokasi",
	StatusProduksi:     "Produksi",
	StatusPengiriman:   "Pengiriman",
	StatusSelesai:      "Selesai",
	StatusBatal:        "Batal",
}

// Sumber pesanan (aturan bisnis No.1).
const (
	SumberInstagram  = "instagram"
	SumberMulutMulut = "mulut_ke_mulut"
	SumberHotel      = "hotel"
	SumberLainnya    = "lainnya"
)

// SumberPesananList adalah kanal pemesanan yang dipakai owner.
var SumberPesananList = []string{SumberInstagram, SumberMulutMulut, SumberHotel, SumberLainnya}

// LabelSumberPesanan memakai nama kanal yang lazim dipakai owner.
var LabelSumberPesanan = map[string]string{
	SumberInstagram:  "Instagram",
	SumberMulutMulut: "Mulut ke Mulut",
	SumberHotel:      "Hotel / Wedding",
	SumberLainnya:    "Lainnya",
}

// Tipe pembayaran.
const (
	TipeDP        = "dp"
	TipePelunasan = "pelunasan"
)

// Status pembayaran turunan dari total pembayaran vs total pesanan.
const (
	StatusBelumDP    = "belum_dp"
	StatusDPTerbayar = "dp_terbayar"
	StatusLunas      = "lunas"
)

// Metode pembayaran.
var MetodePembayaranList = []string{"transfer", "tunai", "kas"}

const (
	// KunciPengaturanBiayaOperasional adalah persentase biaya operasional yang
	// ditambahkan ke HPP bahan saat menghitung HPP per porsi.
	KunciPengaturanBiayaOperasional = "biaya_operasional_persen"

	// DefaultPersenBiayaOperasional dipakai bila pengaturan belum diisi.
	DefaultPersenBiayaOperasional = 10.0

	// MarginMinimumPersen adalah target untung minimum owner (aturan bisnis No.8).
	MarginMinimumPersen = 30.0
)
