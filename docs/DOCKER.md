# Panduan Docker untuk Proyek SIM Dina Catering

Dokumen ini ditulis untuk teman-teman yang **belum pernah memakai Docker**.
Kalau ini pertama kalinya, baca Bagian 1 sampai Bagian 3 dulu — setelah itu
menjalankan proyek ini tinggal satu perintah.

---

## 1. Apa itu Docker? (versi sangat singkat)

Bayangkan Anda memindahkan program dari laptop sendiri ke laptop teman.

Normally you would have to install Go, PostgreSQL, dan semua dependensinya
di laptop teman, lalu berharap versinya sama. **Docker** membuat ini jadi satu
file: Anda meng-*package* program beserta seluruh kebutuhannya ke dalam satu
**image**, lalu image itu bisa dijalankan di mana saja tanpa instalasi apa pun.

Tiga istilah yang perlu Anda kenal:

| Istilah | Analogi | Arti sebenarnya |
|---|---|---|
| **Image** | Cetak biru (blueprint) | Paket berisi program + seluruh dependensinya. Dibuat sekali, bisa dipakai berkali-kali. Tidak bisa "dijalankan" langsung. |
| **Container** | Rumah yang dibangun dari cetak biru | Instance yang sedang berjalan dari sebuah image. Anda bisa punya banyak container dari satu image. |
| **Volume** | Lemari arsip | Penyimpanan data yang **tidak ikut hilang** ketika container dihapus. Ini yang membuat data PostgreSQL Anda tetap ada setelah `docker compose down`. |

Analoginya begini:

```
Image  = Cetak biru rumah  ──►  Container  =  Rumah yang dibangun
                                                 (bisa dihapus, lalu dibangun lagi)
Volume =  Lemari arsip di dalam rumah
          (tidak ikut hilang bersama rumah)
```

> **Yang paling sering bikin bingung:** menghapus container **tidak** menghapus
> data PostgreSQL, karena datanya disimpan di volume, bukan di dalam container.

---

## 2. Yang perlu disiapkan

Cuma satu: **Docker Desktop**.

1. Unduh di <https://www.docker.com/products/docker-desktop/>.
2. Install, lalu **start** Docker Desktop dan tunggu ikonnya menunjukkan
   "Engine running" (di Windows, pastikan **WSL 2** sudah aktif).
3. Cek apakah sudah jalan:

   ```bash
   docker --version
   docker compose version
   ```

Kalau dua perintah itu mengeluarkan versinya, Docker sudah siap.

> Perintah yang dipakai di dokumen ini memakai `docker compose` (v2, spasi).
> Docker versi lama memakai `docker-compose` (v1, dengan tanda hubung).
> Kalau `docker compose` tidak dikenali, coba yang pakai tanda hubung.

---

## 3. Menjalankan proyek ini (3 perintah)

```bash
git clone https://github.com/Luuhhuurr/Dina-Catering-Project-SIM.git
cd Dina-Catering-Project-SIM
docker compose up --build
```

Lalu buka **http://localhost:8080** dan masuk dengan `admin` / `admin123`.

Yang terjadi di balik layar:

1. Docker membaca `Dockerfile`, lalu **build** image aplikasi Go.
2. Docker menarik image `postgres:16-alpine` (kalau belum ada di komputer).
3. Docker membuat volume bernama `dina-catering_pgdata`.
4. PostgreSQL dijalankan, lalu menunggu sampai benar-benar siap (bukan sekadar
   container-nya hidup — lihat `healthcheck` di Bagian 5).
5. Aplikasi Go dijalankan, otomatis membuat tabel (migrasi) dan mengisi data
   awal (seed) kalau database masih kosong.
6. Port `8080` di komputer diteruskan ke port `8080` di dalam container.

### Perintah Docker Compose yang paling sering dipakai

| Perintah | Fungsi |
|---|---|
| `docker compose up --build` | Build image lalu jalankan semua layanan. Ini perintah utama. |
| `docker compose up -d` | Sama, tapi jalankan di latar belakang (tidak ada log yang ikut tampil). |
| `docker compose logs -f app` | Lihat log aplikasi secara realtime (Ctrl+C untuk berhenti, container tetap jalan). |
| `docker compose ps` | Lihat status tiap layanan. |
| `docker compose down` | Matikan semua container. **Data di volume tetap aman.** |
| `docker compose down -v` | Matikan container **dan hapus volume** → database kembali ke kondisi awal saat dijalankan ulang. |
| `docker compose restart app` | Restart hanya satu layanan. |
| `docker compose build --no-cache app` | Build ulang dari nol (dipakai kalau ada masalah cache build yang aneh). |

---

## 4. Membaca `Dockerfile`

`Dockerfile` berisi **dua tahap** (*multi-stage build*). Inilah yang membuat
hasil akhirnya ringan.

### Tahap 1 — `build` (menggantung image `golang:1.26-alpine`)

```dockerfile
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server
```

Penjelasan per baris:

| Baris | Arti |
|---|---|
| `FROM golang:1.26-alpine AS build` | Mulai dari image Go. `AS build` memberi nama tahap ini. |
| `WORKDIR /src` | Semua perintah berikutnya dijalankan di folder `/src` di dalam container. |
| `COPY go.mod go.sum ./` | **Hanya** salin daftar dependensi dulu. |
| `RUN go mod download` | Unduh dependensi. Karena baris ini lapisan (*layer*) tersendiri, Docker menyimpannya di cache. |
| `COPY . .` | Baru salin seluruh source code. |
| `RUN go build ...` | Kompilasi jadi satu file executable `/out/server`. |

**Kenapa `go.mod` disalin duluan?** Supaya setiap kali Anda hanya mengubah
kode Go, Docker tidak perlu mengunduh ulang semua dependensi. Build kedua
jadi jauh lebih cepat.

**Kenapa `CGO_ENABLED=0`?** Go secara default butuh kompiler C untuk
fitur tertentu. Karena image runtime-nya `alpine` (yang tidak punya C
compiler), kita matikan saja fiturnya — proyek ini tidak butuh. Hasilnya
binary statis yang bisa jalan di image Alpine yang sangat kecil.

**Apa arti `-trimpath` dan `-ldflags="-s -w"`?**
`-trimpath` menghapus path absolut dari komputer Anda dari hasil build
(jadi binary lebih reproducible). `-s -w` membuang tabel simbol debug,
sehingga binary lebih kecil.

### Tahap 2 — runtime (image `alpine:3.20`)

```dockerfile
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -u 10001 appuser
ENV TZ=Asia/Jakarta
ENV PORT=8080
WORKDIR /app
COPY --from=build /out/server /app/server
USER appuser
EXPOSE 8080
ENTRYPOINT ["/app/server"]
```

| Baris | Arti |
|---|---|
| `FROM alpine:3.20` | Image runtime baru yang **tidak** memuat Go sama sekali — hanya Linux minimal (~7 MB). |
| `apk add --no-cache ca-certificates tzdata` | Pasang sertifikat CA (untuk koneksi HTTPS keluar) dan basis data zona waktu. `--no-cache` supaya tidak menyimpan cache paket. |
| `adduser -D -u 10001 appuser` | Buat user non-root. Jalankan aplikasi **bukan** sebagai root — praktik keamanan dasar di container. |
| `ENV TZ=Asia/Jakarta` | Zona waktu container. Penting untuk app ini: perhitungan H-30/H-7 harus mengikuti kalender toko di Ungaran, bukan UTC. |
| `COPY --from=build /out/server /app/server` | Ambil **hanya** binary hasil build dari tahap sebelumnya. Compiler Go, source code, dan dependensi build tidak ikut ke image akhir. |
| `USER appuser` | Jalankan sebagai user tadi. |
| `EXPOSE 8080` | Dokumentasi: aplikasi Listen di port 8080. |
| `ENTRYPOINT ["/app/server"]` | Perintah yang dijalankan saat container start. |

**Inti multi-stage build:** image `golang:1.26-alpine` itu besar (±800 MB)
karena berisi compiler Go. Kalau kita用它 langsung untuk runtime, image
akhirnya ikut besar. Dengan dua tahap, yang dikirim ke server hanya
`alpine` + satu file binary — totalnya **sekitar 20–30 MB**.

---

## 5. Membaca `docker-compose.yml`

`docker-compose.yml` adalah "resep": beberapa container yang dijalankan
bersama sebagai satu paket.

### Layanan `postgres`

```yaml
services:
  postgres:
    image: postgres:16-alpine
    restart: unless-stopped
    environment:
      POSTGRES_USER: ${POSTGRES_USER:-dina}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:-dina}
      POSTGRES_DB: ${POSTGRES_DB:-dinacatering}
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER:-dina} -d ${POSTGRES_DB:-dinacatering}"]
      interval: 5s
      timeout: 5s
      retries: 12
      start_period: 10s
```

| Bagian | Fungsi |
|---|---|
| `image: postgres:16-alpine` | Paket PostgreSQL 16 versi Alpine. Tidak perlu `build:` karena image ini sudah ada resmi di Docker Hub. |
| `restart: unless-stopped` | Kalau container mati (mis. laptop tidur lalu bangun), Docker otomatis menyalakannya lagi. |
| `environment` | Nilai diambil dari file `.env` kalau ada. Tanda `${NAMA:-default}` artinya: pakai nilai dari `.env`; kalau tidak ada, pakai nilai setelah `:-`. |
| `volumes: pgdata:/var/lib/postgresql/data` | Folder data PostgreSQL **di dalam container** dipetakan ke volume `pgdata`. Inilah alasan data tidak hilang saat container dihapus. |
| `healthcheck` | Docker menjalankan `pg_isready` tiap 5 detik untuk mengecek apakah PostgreSQL **benar-benar siap menerima koneksi** (bukan cuma container-nya hidup). Maksimum 12 percobaan dengan jeda 5 detik. |

### Layanan `app`

```yaml
  app:
    build: .
    restart: unless-stopped
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      DATABASE_URL: postgres://${POSTGRES_USER:-dina}:${POSTGRES_PASSWORD:-dina}@postgres:5432/${POSTGRES_DB:-dinacatering}?sslmode=disable
      PORT: "8080"
      SESSION_SECRET: ${SESSION_SECRET:-ganti-secret-ini-sebelum-dipakai-produksi}
      ADMIN_PASSWORD: ${ADMIN_PASSWORD:-admin123}
      SEED_ON_BOOT: ${SEED_ON_BOOT:-true}
    ports:
      - "${PORT:-8080}:8080"
```

| Bagian | Fungsi |
|---|---|
| `build: .` | Bangun image dari `Dockerfile` yang ada di folder ini. |
| `depends_on: condition: service_healthy` | **Bagian paling penting.** Aplikasi baru dijalankan setelah PostgreSQL benar-benar *healthy*. Tanpa ini, aplikasi bisa start duluan lalu gagal karena database belum siap. |
| `DATABASE_URL` | Alamat database. Perhatikan host-nya **`postgres`**, bukan `localhost`! Di dalam jaringan Docker, tiap container punya hostname sendiri. `localhost` di dalam container `app` berarti container `app` itu sendiri. |
| `ports: "8080:8080"` | Format `port_komputer:port_container`. Browser Anda mengakses `localhost:8080` di komputer, lalu Docker meneruskannya ke port 8080 di dalam container. |
| `SESSION_SECRET` | Kunci penanda tangan cookie session. **Wajib diganti** untuk pemakaian nyata. |

### Volume

```yaml
volumes:
  pgdata:
```

Deklarasi volume bernama. Docker akan membuatnya otomatis dengan nama lengkap
`dina-catering_pgdata` (nama project + nama volume).

---

## 6. Variabel lingkungan (`.env`)

`docker-compose.yml` membaca nilai dari file `.env` di folder yang sama.
Salin dari `.env.example`:

```bash
cp .env.example .env
```

Isinya:

```ini
POSTGRES_USER=dina
POSTGRES_PASSWORD=dina
POSTGRES_DB=dinacatering
SESSION_SECRET=ganti-secret-ini-sebelum-dipakai-produksi
ADMIN_PASSWORD=admin123
SEED_ON_BOOT=true
PORT=8080
DATABASE_URL=postgres://dina:dina@localhost:5432/dinacatering?sslmode=disable
```

`.env` **tidak pernah** ikut ter-commit ke Git (sudah masuk `.gitignore`),
karena bisa berisi password produksi.

> `DATABASE_URL` di `.env` hanya dipakai kalau Anda menjalankan aplikasi
> **di luar** Docker (`go run ./cmd/server`). Saat memakai `docker compose`,
> nilainya overridden oleh `environment` di `docker-compose.yml`.

---

## 7. Database: di mana datanya tersimpan?

Data PostgreSQL disimpan di **volume Docker**, bukan di file di dalam project.
Lihat volumenya:

```bash
docker volume ls | grep pgdata
# dina-catering_pgdata
```

Kalau mau masuk ke psql (melalui `docker compose exec` — perintah dijalankan
**di dalam** container yang sedang berjalan):

```bash
docker compose exec postgres psql -U dina -d dinacatering
```

Lalu di prompt `psql` ketik:

```sql
SELECT nama, satuan, stok_saat_ini FROM bahan_baku ORDER BY id;
\q          -- keluar dari psql
```

Perintah yang sering dipakai:

| Tujuan | Perintah |
|---|---|
| Lihat semua tabel | `\dt` |
| Lihat isi satu tabel | `SELECT * FROM pesanan;` |
| Keluar | `\q` |

---

## 8. Kalau ada masalah

| Gejala | Penyebab & Solusi |
|---|---|
| `docker compose up` error: port 8080 sudah dipakai | Ada aplikasi lain di port itu. Ubah `PORT` di `.env` (mis. `8081`) atau matikan program yang memakainya. |
| `connection refused` ke database | Database belum siap. Lihat `docker compose logs postgres`. |
| Perubahan kode tidak terlihat | Build ulang: `docker compose up --build --force-recreate app` |
| Ingin mulai dari data awal lagi | `docker compose down -v` lalu `docker compose up --build` |
| Lihat log aplikasi | `docker compose logs -f app` |
| Docker Desktop lambat / error WSL | Restart Docker Desktop, atau pastikan WSL 2 aktif (`wsl --status`). |

---

## 9. Ringkasan perintah harian

```bash
docker compose up --build      # mulai (perintah pertama, dan setiap ada perubahan kode)
docker compose logs -f app     # lihat log
docker compose ps              # cek status
docker compose down            # stop (data aman)
docker compose down -v         # stop + reset database
```

---

## Referensi

* Docker Compose file reference — <https://docs.docker.com/reference/compose-file/>
* Dockerfile reference — <https://docs.docker.com/reference/dockerfile/>
* PostgreSQL Docker official image — <https://hub.docker.com/_/postgres>
* Docker glossary — <https://docs.docker.com/get-started/glossary/>
