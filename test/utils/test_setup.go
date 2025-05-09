package utils

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupPostgres(ctx context.Context, dbName string) (*sql.DB, func(), error) {
	tc, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:15-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env:          map[string]string{"POSTGRES_USER": "user", "POSTGRES_PASSWORD": "pass", "POSTGRES_DB": dbName},
			WaitingFor:   wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		return nil, nil, err
	}
	host, _ := tc.Host(ctx)
	port, _ := tc.MappedPort(ctx, "5432/tcp")
	dsn := fmt.Sprintf("postgres://user:pass@%s:%s/%s?sslmode=disable", host, port.Port(), dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		tc.Terminate(ctx)
		return nil, nil, err
	}
	for i := 0; i < 10; i++ {
		if err := db.PingContext(ctx); err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	wd, err := os.Getwd()
	if err != nil {
		db.Close()
		tc.Terminate(ctx)
		return nil, nil, err
	}
	baseDir := filepath.ToSlash(filepath.Clean(filepath.Join(wd, "..", "..", "..", "..", "scripts", "migrations")))
	m, err := migrate.New("file://"+filepath.ToSlash(baseDir), dsn)
	if err != nil {
		db.Close()
		tc.Terminate(ctx)
		return nil, nil, err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		db.Close()
		tc.Terminate(ctx)
		return nil, nil, err
	}
	testDir := filepath.Join(baseDir, "test")
	files, _ := os.ReadDir(testDir)
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".up.sql") {
			sqlBytes, _ := os.ReadFile(filepath.Join(testDir, f.Name()))
			db.Exec(string(sqlBytes))
		}
	}
	cleanup := func() { db.Close(); tc.Terminate(ctx) }
	return db, cleanup, nil
}
