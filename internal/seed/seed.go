// Package seed mengisi data awal aplikasi saat database masih kosong.
package seed

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"dinacatering/db"
)

// Run mengisi data awal bila tabel users masih kosong. Password admin di-hash
// saat runtime sehingga tidak ada hash yang tertulis di dalam kode.
func Run(ctx context.Context, pool *pgxpool.Pool, adminPassword string) error {
	var jumlahUser int
	if err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&jumlahUser); err != nil {
		return fmt.Errorf("cek data awal: %w", err)
	}
	if jumlahUser > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password admin: %w", err)
	}

	script := strings.ReplaceAll(string(db.SeedSQL), "{PASSWORD_HASH}", string(hash))
	if _, err := pool.Exec(ctx, script); err != nil {
		return fmt.Errorf("jalankan data awal: %w", err)
	}
	return nil
}
