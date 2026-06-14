package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// RunMigrations reads and executes all SQL migration files in order
func RunMigrations(ctx context.Context, pool *pgxpool.Pool, migrationsDir string, logger *zap.Logger) error {
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory %s: %w", migrationsDir, err)
	}

	// Filter and sort .sql files
	var sqlFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			sqlFiles = append(sqlFiles, f.Name())
		}
	}
	sort.Strings(sqlFiles)

	if len(sqlFiles) == 0 {
		logger.Warn("No migration files found", zap.String("dir", migrationsDir))
		return nil
	}

	for _, fileName := range sqlFiles {
		path := filepath.Join(migrationsDir, fileName)
		sql, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", fileName, err)
		}

		_, err = pool.Exec(ctx, string(sql))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", fileName, err)
		}

		logger.Info("Migration applied successfully", zap.String("file", fileName))
	}

	return nil
}
