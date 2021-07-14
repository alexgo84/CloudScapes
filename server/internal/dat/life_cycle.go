package dat

import (
	"CloudScapes/pkg/logger"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	dbMap *sqlx.DB
}

func NewDB(ctx context.Context) (*DB, error) {
	user := getEnv("POSTGRES_USER", "cloudscapes")
	pass := getEnv("POSTGRES_DB", "cloudscapes")
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	dbName := getEnv("POSTGRES_DB_NAME", "cloudscapes")

	connectionURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, dbName)

	var err error
	dbMap, err := sqlx.ConnectContext(ctx, "pgx", connectionURL)
	if err != nil {
		return nil, err
	}
	logger.Log(logger.INFO, fmt.Sprintf("Connected user '%s' to DB '%s' @ %s:%s", user, dbName, host, port))

	return &DB{dbMap: dbMap}, nil
}

func (db *DB) PingDB(ctx context.Context) error {
	err := db.dbMap.PingContext(ctx)
	return err
}

func (db *DB) Close() error {
	err := db.dbMap.Close()
	return err
}

func (db *DB) GetNewTransaction(ctx context.Context) (*sqlx.Tx, error) {
	txn, err := db.dbMap.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return txn, nil
}

func (db *DB) RunMigrations(ctx context.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	migrationsDirectory := cwd + "/../../server/internal/dat/migrations"
	files, err := ioutil.ReadDir(migrationsDirectory)
	if err != nil {
		return err
	}

	// sort migration files so they are run in correct order
	sort.Slice(files, func(i, j int) bool {
		return strings.Compare(files[i].Name(), files[j].Name()) < 0
	})

	allMigrations := `
	CREATE TABLE IF NOT EXISTS migrations(
		id serial PRIMARY KEY,
		name text NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW() 
	);`

	for _, f := range files {
		fileName := f.Name()
		if path.Ext(fileName) != ".pgsql" {
			logger.Log(logger.WARN, "skipping file in migration folder since it doesn't have 'pgsql' extention", logger.Str("filename", fileName))
			continue
		}
		migrationName := strings.ToLower(f.Name())
		migrationName = migrationName[0 : len(migrationName)-len(".pgsql")]

		migration, err := ioutil.ReadFile(filepath.Join(migrationsDirectory, fileName))
		if err != nil {
			return err
		}

		decoratedMigration := fmt.Sprintf(`
			DO $$
			DECLARE
				migration_name text := '%s';
			BEGIN
				IF EXISTS (
					SELECT
						1
					FROM
						migrations
					WHERE
						name = migration_name) THEN
				RETURN;
			END IF;

			%s

			INSERT INTO migrations (name)
			VALUES (migration_name);
			END
			$$;`, migrationName, string(migration))

		allMigrations += decoratedMigration
	}
	//fmt.Println(allMigrations)

	_, err = db.dbMap.ExecContext(ctx, allMigrations)
	return err
}

func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
