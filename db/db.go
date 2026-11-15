// Package db memegang file SQL yang dibutuhkan aplikasi: skema (untuk sqlc)
// dan data awal (seed).
package db

import _ "embed"

//go:embed seed.sql
var SeedSQL []byte
