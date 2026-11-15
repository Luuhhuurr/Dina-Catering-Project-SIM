package handler

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"

	"dinacatering/internal/model"
)

// templateFuncs berisi helper yang dipakai seluruh template.
func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"rupiah":          rupiah,
		"tanggal":         func(t time.Time) string { return t.Format("02 Jan 2006") },
		"tanggalPanjang":  func(t time.Time) string { return t.Format("02 January 2006") },
		"tanggalPendek":   func(t time.Time) string { return t.Format("02/01/2006") },
		"bulan":           func(t time.Time) string { return t.Format("January 2006") },
		"hari":            func(t time.Time) string { return t.Format("Monday") },
		"numerik":         func(n float64) string { return strconv.FormatFloat(n, 'f', 2, 64) },
		"tanpaDesimal":    func(n float64) string { return strconv.FormatFloat(n, 'f', -1, 64) },
		"persen":          func(n float64) string { return fmt.Sprintf("%.1f%%", n) },
		"labelStatus":     func(s string) string { return model.LabelStatusPesanan[s] },
		"labelSumber":     func(s string) string { return model.LabelSumberPesanan[s] },
		"labelJatuhTempo": func(s string) string { return labelJatuhTempo(s) },
		"tambah":          func(a, b int) int { return a + b },
		"kurang":          func(a, b int) int { return a - b },
		"sub":             func(a, b int64) int64 { return a - b },
		"today":           hariIni,
		"year":            func(t time.Time) int { return t.Year() },
		"eq":              func(a, b string) bool { return a == b },
		"orDefault":       orDefault,
	}
}

// rupiah memformat angka menjadi Rupiah dengan pemisah ribuan.
func rupiah(n int64) string {
	negatif := n < 0
	if negatif {
		n = -n
	}
	teks := strconv.FormatInt(n, 10)
	var bagian []string
	for len(teks) > 3 {
		bagian = append([]string{teks[len(teks)-3:]}, bagian...)
		teks = teks[:len(teks)-3]
	}
	bagian = append([]string{teks}, bagian...)
	hasil := strings.Join(bagian, ".")
	if negatif {
		hasil = "-" + hasil
	}
	return "Rp " + hasil
}

// labelJatuhTempo memberi keterangan singkat untuk daftar menunggu pembayaran.
func labelJatuhTempo(status string) string {
	switch status {
	case "dp_jatuh_tempo":
		return "DP jatuh tempo"
	case "dp_terlambat":
		return "DP terlambat"
	case "pelunasan_jatuh_tempo":
		return "Pelunasan jatuh tempo"
	case "pelunasan_terlambat":
		return "Pelunasan terlambat"
	default:
		return "Lunas"
	}
}

// orDefault mengembalikan nilai cadangan bila string kosong.
func orDefault(nilai, cadangan string) string {
	if nilai == "" {
		return cadangan
	}
	return nilai
}
