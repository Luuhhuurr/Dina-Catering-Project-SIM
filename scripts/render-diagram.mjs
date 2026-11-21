// render-diagram.mjs — merender sumber Mermaid pada Diagram,dll/mermaid/*.mmd
// menjadi PNG pada Diagram,dll/gambar/.
//
// Pakai:  node scripts/render-diagram.mjs
//
// Memerlukan sebuah browser untuk merender. Script ini memakai Edge yang sudah ada di
// mesin Windows lewat puppeteerConfigFile; jalankan `npm i @mermaid-js/mermaid-cli`
// lebih dulu bila belum terpasang.

import { readFileSync, writeFileSync, mkdirSync, readdirSync, existsSync } from 'node:fs';
import { join, dirname, basename } from 'node:path';
import { fileURLToPath } from 'node:url';
import { execFileSync } from 'node:child_process';
import { tmpdir } from 'node:os';

const root = join(dirname(fileURLToPath(import.meta.url)), '..');
const srcDir = join(root, 'Diagram,dll');
const mmdDir = join(srcDir, 'mermaid');
const outDir = join(srcDir, 'gambar');
const tmpDir = join(tmpdir(), 'dina-mermaid');

const EDGE_CANDIDATES = [
  'C:/Program Files (x86)/Microsoft/Edge/Application/msedge.exe',
  'C:/Program Files/Microsoft/Edge/Application/msedge.exe',
  'C:/Program Files/Google/Chrome/Application/chrome.exe',
  'C:/Program Files (x86)/Google/Chrome/Application/chrome.exe',
];

// Lokasi binary mmdc. Default: node_modules/.bin di repo. Bila mermaid-cli
// dipasang di luar repo, set MMDC_PATH ke binary tersebut.
const mmdcCandidates = [
  process.env.MMDC_PATH,
  join(root, 'node_modules', '.bin', process.platform === 'win32' ? 'mmdc.cmd' : 'mmdc'),
  'B:/Latihan Golang/mermaid-tools/node_modules/.bin/mmdc.cmd',
  'B:/Latihan Golang/mermaid-tools/node_modules/.bin/mmdc',
].filter(Boolean);

const mmdc = mmdcCandidates.find((p) => existsSync(p));
if (!mmdc) {
  console.error(
    'Binary mmdc tidak ditemukan. Pasang @mermaid-js/mermaid-cli lalu set MMDC_PATH ke lokasinya.',
  );
  process.exit(1);
}

// Di Windows, file .cmd tidak bisa di-spawn langsung oleh execFileSync tanpa
// shell. Karena itu jalankan entry JavaScript miliknya melalui node.
const mmdcJs = mmdc.endsWith('.cmd')
  ? join(dirname(dirname(mmdc)), '@mermaid-js', 'mermaid-cli', 'src', 'cli.js')
  : null;

const command = mmdcJs && existsSync(mmdcJs) ? process.execPath : mmdc;
const commandArgs = mmdcJs && existsSync(mmdcJs) ? [mmdcJs] : [];

// Viewport harus lebih besar daripada diagram terlebar agar tidak terpotong.
const puppeteerConfig = {
  executablePath: EDGE_CANDIDATES.find((p) => existsSync(p)),
  defaultViewport: { width: 3600, height: 4800 },
};

mkdirSync(outDir, { recursive: true });
mkdirSync(tmpDir, { recursive: true });

writeFileSync(join(tmpDir, 'puppeteer.json'), JSON.stringify(puppeteerConfig, null, 2));

const sources = readdirSync(mmdDir)
  .filter((f) => f.endsWith('.mmd'))
  .sort();

let gagal = 0;
for (const file of sources) {
  const slug = basename(file, '.mmd');
  const pngFile = join(outDir, `${slug}.png`);
  const svgFile = join(outDir, `${slug}.svg`);
  try {
    // PNG untuk sisipan cepat; SVG untuk diagram lebar supaya tetap tajam
    // ketika diskalakan ke ukuran halaman A4.
    execFileSync(command, [
      ...commandArgs,
      '-i', join(mmdDir, file),
      '-o', pngFile,
      '-b', 'white',
      '-s', '2',
      '-p', join(tmpDir, 'puppeteer.json'),
    ], { stdio: 'pipe' });
    execFileSync(command, [
      ...commandArgs,
      '-i', join(mmdDir, file),
      '-o', svgFile,
      '-b', 'white',
      '-p', join(tmpDir, 'puppeteer.json'),
    ], { stdio: 'pipe' });
    console.log(`OK   ${file} -> gambar/${slug}.png + .svg`);
  } catch (e) {
    gagal++;
    console.error(`GAGAL ${file}: ${String(e.stderr || e.message).slice(0, 400)}`);
  }
}

console.log(`\nSelesai: ${sources.length - gagal} berhasil, ${gagal} gagal.`);
process.exit(gagal > 0 ? 1 : 0);
