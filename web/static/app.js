// Notifikasi singkat setelah aksi HTMX (HX-Trigger: showToast).
document.body.addEventListener("showToast", function (event) {
  var pesan = event.detail && event.detail.value;
  if (!pesan) return;

  var wadah = document.getElementById("toast");
  if (!wadah) return;

  var notif = document.createElement("div");
  notif.className = "notif";
  notif.innerHTML =
    '<svg class="i"><use href="#i-cek"></use></svg><span></span>';
  notif.querySelector("span").textContent = pesan;
  wadah.appendChild(notif);

  setTimeout(function () {
    notif.remove();
  }, 3500);
});

// Konfirmasi aksi yang merusak data (hapus).
document.body.addEventListener("htmx:confirm", function (event) {
  var pesan = event.detail.question;
  if (!pesan) return;
  if (!window.confirm(pesan)) {
    event.preventDefault();
  }
});

// Form pembayaran: nominal otomatis mengikuti jenis (DP 50% atau pelunasan).
document.body.addEventListener("change", function (event) {
  if (event.target.id !== "tipe") return;

  var nominal = document.getElementById("jumlah");
  if (!nominal) return;

  var nilai = event.target.value === "pelunasan"
    ? nominal.dataset.pelunasan
    : nominal.dataset.dp;

  if (nilai) nominal.value = nilai;
});

// Laci navigasi di layar kecil: tutup dengan tombol Escape atau setelah
// memilih menu. Controlled by a checkbox so it also works without JavaScript.
var navToggle = document.getElementById("nav-toggle");

document.addEventListener("keydown", function (event) {
  if (event.key === "Escape" && navToggle && navToggle.checked) {
    navToggle.checked = false;
  }
});

document.addEventListener("click", function (event) {
  if (!navToggle || !navToggle.checked) return;

  var sidebar = document.querySelector(".sidebar");
  var label = event.target.closest ? event.target.closest(".nav-btn, .nav-tirai") : null;
  if (label) return;

  // Klik di dalam sidebar (bukan pada linknya) tidak menutup laci.
  if (sidebar && event.target.closest && event.target.closest(".sidebar")) return;

  navToggle.checked = false;
});
