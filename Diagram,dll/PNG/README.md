# Kumpulan Diagram PNG — SIM Dina Catering

Folder ini berisi seluruh diagram yang sudah dirender menjadi PNG, dikumpulkan dalam
satu tempat agar mudah dicari saat menyisipkan ke laporan tugas akhir.

| # | Berkas PNG | Isi diagram | Sumber |
|---|---|---|---|
| 1 | `01-proses-usulan-manual.png` | Alur proses manual (situasi lama) | `01-Diagram-Proses-Usulan.md` |
| 2 | `02-proses-usulan-sistem-baru.png` | Alur proses usulan (sistem baru) | `01-Diagram-Proses-Usulan.md` |
| 3 | `03-perhitungan-hpp.png` | Alur perhitungan HPP | `01-Diagram-Proses-Usulan.md` |
| 4 | `04-pengadaan-stok.png` | Alur pengadaan bahan dan stok | `01-Diagram-Proses-Usulan.md` |
| 5 | `05-pelaporan-keuangan.png` | Alur pelaporan keuangan | `01-Diagram-Proses-Usulan.md` |
| 6 | `06-dfd-level-0.png` | Context diagram (DFD Level 0) | `02-DFD-Level-0-dan-Level-1.md` |
| 7 | `07-dfd-level-1.png` | DFD Level 1 (6 proses, 9 data store) | `02-DFD-Level-0-dan-Level-1.md` |
| 8 | `08-use-case.png` | Use case diagram (4 aktor, 18 use case) | `03-Diagram-Use-Case.md` |
| 9 | `09-erd.png` | ERD (14 tabel, 10 relasi FK) | `04-ERD.md` |

Berkas asli (bersama versi SVG) ada di [`../gambar/`](../gambar/). Sumber editnya
berada di [`../mermaid/`](../mermaid/) dalam format `.mmd`.

## Ukuran gambar

| Berkas | Resolusi |
|---|---|
| `01-proses-usulan-manual.png` | 798 × 3666 |
| `02-proses-usulan-sistem-baru.png` | 1104 × 8892 |
| `03-perhitungan-hpp.png` | 900 × 4882 |
| `04-pengadaan-stok.png` | 1006 × 3374 |
| `05-pelaporan-keuangan.png` | 1142 × 3792 |
| `06-dfd-level-0.png` | 2128 × 1250 |
| `07-dfd-level-1.png` | 7102 × 4006 |
| `08-use-case.png` | 3286 × 4016 |
| `09-erd.png` | 4694 × 2902 |

Diagram proses (1–5) sengaja digambar vertikal memanjang agar mudah dibaca di halaman
cetak A4. Diagram relasional (6–9) digambar mendatar.

## Render ulang

Bila diagram diedit, render ulang dengan:

```bash
node scripts/render-diagram.mjs
```

Lalu salin kembali hasil PNG dari `gambar/` ke folder ini.
