// Package migrations membungkus file SQL migrasi agar bisa dijalankan dari
// dalam binary (tidak perlu executable migrate terpisah di image Docker).
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
