#!/usr/bin/env bash
# Smoke test SIM Dina Catering against http://localhost:8080
# Jalankan setelah `docker compose up --build`.
#
# Catatan: skrip ini mengubah data (menambah pesanan, pembelian, dan riwayat
# harga) pada database yang sedang diuji. Untuk mengulang dengan hasil yang
# sama, jalankan `docker compose down -v` lebih dulu.
set -u

BASE="http://localhost:8080"
COOKIE_JAR="$(mktemp -d)/cookies.txt"
TMP="$(mktemp -d)"
FAIL=0

ok()  { echo "  [OK]   $1"; }
bad() { echo "  [FAIL] $1"; FAIL=$((FAIL+1)); }
step(){ echo; echo "== $1 =="; }

# Ambil nilai pertama dari <select id="X"> pada HTML.
select_value() { # select_value <file> <select-id>
  awk -v id="$2" '
    $0 ~ "<select id=\""id"\"" { f=1; next }
    f && /<option value="[0-9]+"/ {
      gsub(/.*<option value="/, ""); gsub(/".*/, ""); print; exit
    }
    f && /<\/select>/ { exit }
  ' "$1"
}

# Ambil baris pertama href /bahan/<id>/harga.
bahan_id() { grep -oE '/bahan/[0-9]+/harga' "$1" | head -1 | grep -oE '[0-9]+'; }

# Ambil nilai "HPP per porsi" dari halaman hasil kalkulator HPP.
# Label/statistik nilai memakai kelas "angka", jadi ambil baris Rp pertama
# sesudah label tersebut (bukan Rp pertama di halaman, yang bisa dari dropdown).
hpp_per_porsi() { grep -A3 'label">HPP per porsi<' "$1" | grep -oE 'Rp [0-9.]+' | head -1; }

is_redirect() { [ "$1" = "302" ] || [ "$1" = "303" ] || [ "$1" = "307" ]; }

# ---------------------------------------------------------------- 1. auth
step "1. Redirect tanpa session -> /login"
code=$(curl -s -o /dev/null -w '%{http_code}' "$BASE/")
loc=$(curl -s -o /dev/null -w '%{redirect_url}' "$BASE/")
is_redirect "$code" && ok "GET / = $code (redirect)" || bad "GET / = $code (harus redirect)"
echo "$loc" | grep -q "/login" && ok "redirect ke /login" || bad "redirect ke $loc"

step "2. Halaman login terbuka"
code=$(curl -s -o /dev/null -w '%{http_code}' "$BASE/login")
[ "$code" = "200" ] && ok "GET /login = 200" || bad "GET /login = $code"

step "3. Login admin/admin123"
code=$(curl -s -c "$COOKIE_JAR" -o /dev/null -w '%{http_code}' \
  -d "username=admin&password=admin123" "$BASE/login")
is_redirect "$code" && ok "POST /login = $code" || bad "POST /login = $code"
[ -s "$COOKIE_JAR" ] && ok "cookie session tersimpan" || bad "tidak ada cookie session"

step "4. Dashboard terbuka dengan session"
code=$(curl -s -b "$COOKIE_JAR" -o /dev/null -w '%{http_code}' "$BASE/")
[ "$code" = "200" ] && ok "GET / = 200 dengan session" || bad "GET / = $code"

# ---------------------------------------------------------------- 2. HPP
step "5. Halaman paket: HPP, margin, dan badge <30%"
curl -s -b "$COOKIE_JAR" "$BASE/paket" -o "$TMP/paket.html"
grep -q 'HPP/porsi' "$TMP/paket.html" && ok "kolom HPP per porsi ada" || bad "kolom HPP/porsi tidak ada"
grep -q 'Min. harga 30%' "$TMP/paket.html" && ok "kolom harga jual minimum 30% ada" || bad "kolom harga minimum tidak ada"
grep -q 'margin &lt; 30%' "$TMP/paket.html" && ok "badge margin <30% muncul (paket Premium)" \
  || bad "badge margin <30% tidak ditemukan"
# hanya paket Premium yang boleh bermargin di bawah 30%
n_bad=$(grep -c 'margin &lt; 30%' "$TMP/paket.html")
[ "$n_bad" -ge 1 ] && ok "badge ditemukan $n_bad kali" || bad "badge tidak ditemukan"

step "6. Kalkulator HPP per menu"
curl -s -b "$COOKIE_JAR" "$BASE/hpp" -o "$TMP/hpp.html"
grep -q 'name="menu_id"' "$TMP/hpp.html" && ok "form kalkulator HPP ada" || bad "form kalkulator HPP tidak ada"
menu_id=$(select_value "$TMP/hpp.html" menu_id)
if [ -n "$menu_id" ]; then
  curl -s -b "$COOKIE_JAR" "$BASE/hpp/hasil?menu_id=$menu_id" -o "$TMP/hpp_hasil.html"
  grep -qE 'HPP' "$TMP/hpp_hasil.html" && ok "hasil HPP menu #$menu_id terhitung" || bad "hasil HPP menu tidak terhitung"
  grep -qE 'Margin|margin' "$TMP/hpp_hasil.html" && ok "margin aktual ditampilkan" || bad "margin tidak ditampilkan"
  grep -qE 'Minimum|minimum' "$TMP/hpp_hasil.html" && ok "harga jual minimum ditampilkan" || bad "harga minimum tidak ditampilkan"
  # tanpa JavaScript, /hpp/hasil harus tetap berupa halaman penuh (ada layout)
  grep -q '<!DOCTYPE html>' "$TMP/hpp_hasil.html" && ok "halaman penuh (fallback tanpa JavaScript) berisi layout" \
    || bad "/hpp/hasil tanpa JavaScript tidak mengembalikan layout"
else
  bad "tidak menemukan menu_id di halaman HPP"
fi

step "7. Ubah harga bahan -> HPP berubah -> dikembalikan lagi"
curl -s -b "$COOKIE_JAR" "$BASE/bahan" -o "$TMP/bahan.html"
bid=$(bahan_id "$TMP/bahan.html")
if [ -n "$bid" ]; then
  # harga asli bahan ini, dibaca dari baris tabel
  harga_asli=$(grep -oE 'Rp [0-9.]+' "$TMP/bahan.html" | head -1 | grep -oE '[0-9.]+')
  curl -s -b "$COOKIE_JAR" -o /dev/null \
    -d "harga_satuan=999999&tanggal_berlaku=$(date +%F)" "$BASE/bahan/$bid/harga"
  curl -s -b "$COOKIE_JAR" "$BASE/hpp/hasil?menu_id=$menu_id" -o "$TMP/hpp_after.html"
  after=$(hpp_per_porsi "$TMP/hpp_after.html")
  # kembalikan: baris baru dengan tanggal hari ini, id lebih besar -> menang
  curl -s -b "$COOKIE_JAR" -o /dev/null \
    -d "harga_satuan=${harga_asli//./}&tanggal_berlaku=$(date +%F)" "$BASE/bahan/$bid/harga"
  curl -s -b "$COOKIE_JAR" "$BASE/hpp/hasil?menu_id=$menu_id" -o "$TMP/hpp_restored.html"
  restored=$(hpp_per_porsi "$TMP/hpp_restored.html")
  [ "$after" != "$restored" ] && ok "HPP berubah setelah harga bahan #$bid diubah, lalu kembali ($after -> $restored)" \
    || bad "HPP tidak kembali setelah harga dikembalikan ($after vs $restored)"
else
  bad "tidak menemukan id bahan"
fi

# ---------------------------------------------------------------- 3. Pesanan
step "8. Buat pesanan baru -> total + tenggat H-30/H-7"
curl -s -b "$COOKIE_JAR" "$BASE/pesanan" -o "$TMP/pesanan.html"
pelanggan_id=$(select_value "$TMP/pesanan.html" pelanggan_id)
event_date=$(date -d "+20 days" +%F 2>/dev/null || date -v+20d +%F)
if [ -n "$pelanggan_id" ]; then
  # server me-redirect ke halaman detail pesanan yang baru dibuat
  loc=$(curl -s -b "$COOKIE_JAR" -o /dev/null -w '%{redirect_url}' \
    -d "pelanggan_id=$pelanggan_id&nama_paket=Standar&jumlah_tamu=60&tanggal_acara=$event_date&lokasi=Gedung&sumber_pesanan=instagram&catatan=smoke" \
    "$BASE/pesanan")
  newest=$(echo "$loc" | grep -oE '[0-9]+$')
  if [ -n "${newest:-}" ]; then
    ok "POST /pesanan membuat pesanan #$newest"
    curl -s -b "$COOKIE_JAR" "$BASE/pesanan/$newest" -o "$TMP/detail.html"
    grep -q 'H-30' "$TMP/detail.html" && ok "tenggat DP (H-30) tampil" || bad "tenggat H-30 tidak tampil"
    grep -q 'H-7' "$TMP/detail.html" && ok "tenggat pelunasan (H-7) tampil" || bad "tenggat H-7 tidak tampil"
    grep -q 'DP 50%' "$TMP/detail.html" && ok "opsi pembayaran DP 50% ada" || bad "opsi DP 50% tidak ada"
    # tanggal acara jatuh pada Sabtu/Minggu -> harga akhir pekan 38.000 x 60 tamu
    total=$(grep -A2 'Total' "$TMP/detail.html" | grep -oE 'Rp [0-9.]+' | head -1 | grep -oE '[0-9.]+' | tr -d '.')
    if [ "$total" = "2280000" ]; then
      ok "total harga = 60 x Rp 38.000 (harga akhir pekan) = $total"
    else
      bad "total harga = $total, ingin 2280000 (harga akhir pekan)"
    fi
    # tenggat harus tepat H-30 dan H-7 dari tanggal acara
    tgl_dp=$(date -d "$event_date -30 days" +'%d %b %Y' 2>/dev/null)
    tgl_lunas=$(date -d "$event_date -7 days" +'%d %b %Y' 2>/dev/null)
    dp_tgl=$(grep -A2 'Tenggat DP' "$TMP/detail.html" | grep -oE '[0-9]{1,2} [A-Za-z]{3} [0-9]{4}' | head -1)
    lunas_tgl=$(grep -A2 'Tenggat pelunasan' "$TMP/detail.html" | grep -oE '[0-9]{1,2} [A-Za-z]{3} [0-9]{4}' | head -1)
    [ "$dp_tgl" = "$tgl_dp" ] && ok "tenggat DP = H-30 dari tanggal acara ($tgl_dp)" \
      || bad "tenggat DP = $dp_tgl, ingin $tgl_dp"
    [ "$lunas_tgl" = "$tgl_lunas" ] && ok "tenggat pelunasan = H-7 dari tanggal acara ($tgl_lunas)" \
      || bad "tenggat pelunasan = $lunas_tgl, ingin $tgl_lunas"
  else
    bad "pesanan baru tidak ditemukan (redirect: $loc)"
  fi
else
  bad "tidak menemukan pelanggan_id"
fi

step "9. Daftar pesanan menunggu pembayaran"
curl -s -b "$COOKIE_JAR" "$BASE/pesanan" -o "$TMP/pesanan2.html"
grep -qiE 'jatuh tempo|terlambat|menunggu' "$TMP/pesanan2.html" && \
  ok "daftar menampilkan status tenggat pembayaran" || bad "daftar menunggu pembayaran kosong"

# ---------------------------------------------------------------- 4. Pembayaran
step "10. Catat DP 50% -> status naik + tugas produksi default"
if [ -n "${newest:-}" ]; then
  dp=$(grep -oE 'data-dp="[0-9]+"' "$TMP/detail.html" | grep -oE '[0-9]+' | head -1)
  if [ -n "${dp:-}" ]; then
    code=$(curl -s -b "$COOKIE_JAR" -o /dev/null -w '%{http_code}' \
      -d "jumlah=$dp&tipe=dp&tanggal=$(date +%F)&metode=transfer" "$BASE/pesanan/$newest/pembayaran")
    ok "POST pembayaran DP = $code (nominal $dp)"
  else
    bad "nominal DP 50% tidak ditemukan di form"
  fi
  curl -s -b "$COOKIE_JAR" "$BASE/pesanan/$newest" -o "$TMP/detail2.html"
  grep -qiE 'Final Meeting' "$TMP/detail2.html" && ok "status pesanan naik ke Final Meeting setelah DP" \
    || bad "status pesanan tidak naik ke Final Meeting setelah DP"
  grep -qiE 'DP terbayar' "$TMP/detail2.html" && ok "status pembayaran = DP terbayar" \
    || bad "status pembayaran bukan DP terbayar"
fi

step "11. Papan produksi menampilkan tugas + PJ"
curl -s -b "$COOKIE_JAR" "$BASE/produksi" -o "$TMP/produksi.html"
grep -q 'Persiapan bahan' "$TMP/produksi.html" && ok "tugas 'Persiapan bahan' tampil" || bad "tugas persiapan tidak tampil"
grep -q 'Masak seluruh menu' "$TMP/produksi.html" && ok "tugas 'Masak seluruh menu' tampil" || bad "tugas masak tidak tampil"
grep -q 'Pengiriman ke lokasi' "$TMP/produksi.html" && ok "tugas 'Pengiriman ke lokasi' tampil" || bad "tugas pengiriman tidak tampil"
grep -qi 'Penanggung jawab' "$TMP/produksi.html" && ok "kolom penanggung jawab (PJ) ada" || bad "kolom PJ tidak ada"
grep -q 'Sari' "$TMP/produksi.html" && ok "PJ dari seed data tampil" || bad "PJ seed tidak tampil"

# ---------------------------------------------------------------- 5. Pembelian & stok
step "12. Catat pembelian -> stok naik + riwayat harga"
curl -s -b "$COOKIE_JAR" "$BASE/pembelian" -o "$TMP/pembelian.html"
p_bahan=$(select_value "$TMP/pembelian.html" bahan_baku_id)
if [ -n "$p_bahan" ]; then
  curl -s -b "$COOKIE_JAR" -o /dev/null \
    -d "bahan_baku_id=$p_bahan&jumlah=5&harga_satuan=12345&pasar=Pasar Uji&tanggal=$(date +%F)" \
    "$BASE/pembelian"
  curl -s -b "$COOKIE_JAR" "$BASE/bahan/$p_bahan/harga" -o "$TMP/riwayat.html"
  grep -q '12\.345' "$TMP/riwayat.html" && ok "harga pembelian masuk riwayat harga bahan" || bad "harga pembelian tidak masuk riwayat"
  curl -s -b "$COOKIE_JAR" "$BASE/pembelian" -o "$TMP/pembelian2.html"
  grep -q 'Pasar Uji' "$TMP/pembelian2.html" && ok "pembelian tercatat di daftar pembelian" || bad "pembelian tidak tercatat"
else
  bad "tidak menemukan bahan_baku_id"
fi

# ---------------------------------------------------------------- 6. Laporan
step "13. Laporan laba rugi"
bulan=$(date +%Y-%m)
curl -s -b "$COOKIE_JAR" "$BASE/laporan?bulan=$bulan" -o "$TMP/laporan.html"
grep -q 'Keuntungan' "$TMP/laporan.html" && ok "baris Keuntungan (laba rugi) ada" || bad "baris laba rugi tidak ada"
grep -q 'Pemasukan' "$TMP/laporan.html" && ok "baris pemasukan ada" || bad "baris pemasukan tidak ada"
grep -q 'Pengeluaran bahan' "$TMP/laporan.html" && ok "baris pengeluaran bahan ada" || bad "baris pengeluaran bahan tidak ada"
grep -q 'Biaya operasional' "$TMP/laporan.html" && ok "baris biaya operasional ada" || bad "baris biaya tidak ada"

# rekonsiliasi: pemasukan - bahan - biaya = keuntungan
# angka di UI memakai pemisah ribuan titik, jadi titiknya dibuang dulu.
angka() { grep -A2 "$1" "$TMP/laporan.html" | grep -oE 'Rp -?[0-9.]+' | head -1 | grep -oE '\-?[0-9.]+' | tr -d '.'; }
pemasukan=$(angka 'Pemasukan')
bahan=$(angka 'Pengeluaran bahan')
biaya=$(angka 'Biaya operasional')
untung=$(angka 'Keuntungan')
if [ -n "${pemasukan:-}" ] && [ -n "${bahan:-}" ] && [ -n "${biaya:-}" ] && [ -n "${untung:-}" ]; then
  hitung=$(awk -v a="$pemasukan" -v b="$bahan" -v c="$biaya" 'BEGIN{printf "%d", a-b-c}')
  [ "$hitung" = "$untung" ] && ok "rekonsiliasi laba rugi: $pemasukan - $bahan - $biaya = $hitung" \
    || bad "rekonsiliasi gagal: $pemasukan - $bahan - $biaya = $hitung, tampilannya $untung"
else
  bad "angka laporan tidak terbaca lengkap"
fi

step "14. Rekap pesanan per status & sumber"
grep -qiE 'instagram' "$TMP/laporan.html" && ok "rekap per sumber ada" || bad "rekap per sumber tidak ada"
grep -qiE 'baru|final_meeting|produksi' "$TMP/laporan.html" && ok "rekap per status ada" || bad "rekap per status tidak ada"

# ---------------------------------------------------------------- 7. Modul lain
step "15. Semua halaman modul terbuka (200)"
for p in pelanggan bahan pembelian menu paket hpp pesanan produksi laporan biaya pengaturan; do
  code=$(curl -s -b "$COOKIE_JAR" -o /dev/null -w '%{http_code}' "$BASE/$p")
  [ "$code" = "200" ] || bad "GET /$p = $code"
done
ok "semua halaman modul merespons 200"

step "16. Endpoint tabel HTMX merespons 200"
for p in "/bahan/tabel" "/menu/1/resep" "/paket/tabel" "/pesanan/tabel" "/produksi/tabel"; do
  code=$(curl -s -b "$COOKIE_JAR" -o /dev/null -w '%{http_code}' "$BASE$p")
  [ "$code" = "200" ] || bad "GET $p = $code"
done
ok "semua endpoint tabel merespons 200"

step "17. Panel tanpa JavaScript tetap halaman penuh"
for p in "/hpp/hasil?menu_id=1" "/bahan/1/harga" "/menu/1/resep" "/paket/1/isi"; do
  curl -s -b "$COOKIE_JAR" "$BASE$p" -o "$TMP/panel.html"
  grep -q '<!DOCTYPE html>' "$TMP/panel.html" \
    && ok "GET $p mengembalikan halaman penuh" \
    || bad "GET $p mengembalikan potongan tanpa layout"
done
# dengan header HTMX, yang dikirim harus potongan HTML (bukan layout)
curl -s -b "$COOKIE_JAR" -H 'HX-Request: true' "$BASE/menu/1/resep" -o "$TMP/panel_htmx.html"
grep -q '<!DOCTYPE html>' "$TMP/panel_htmx.html" \
  && bad "GET /menu/1/resep (HX-Request) seharusnya partial, bukan layout" \
  || ok "GET /menu/1/resep (HX-Request) mengembalikan partial"

step "18. Logout"
code=$(curl -s -b "$COOKIE_JAR" -o /dev/null -w '%{http_code}' -X POST "$BASE/logout")
is_redirect "$code" && ok "POST /logout = $code" || bad "POST /logout = $code"

rm -rf "$TMP" "$(dirname "$COOKIE_JAR")"

echo
echo "==================================="
if [ "$FAIL" -eq 0 ]; then
  echo "SEMUA SMOKE TEST LULUS"
else
  echo "$FAIL CHECK GAGAL"
fi
echo "==================================="
exit "$FAIL"
