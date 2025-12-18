package cmd

import (
	"context"
	"simple-golang/config"

	"github.com/labstack/gommon/log"

	_ "simple-golang/internal/database/migration"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migrationas commands (up, down, status, etc.)",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Errorf("[RunMigrate-1] Please provide a goose command, e.g. migrate up")
		}

		cfg := config.NewConfig()

		// Ambil *sql.DB (bukan gorm.DB)
		sqlDB, err := cfg.ConnectionSqlDB()
		if err != nil {
			log.Errorf("[RunMigrate-2] failed to connect db: %v", err)
		}
		defer func() {
			if err := sqlDB.Close(); err != nil {
				log.Errorf("[RunMigrate-3] failed to close DB: %v", err)
			}
		}()

		// Jalankan goose command
		dir := "internal/database/migration"
		if err := goose.RunContext(context.Background(), args[0], sqlDB, dir, args[1:]...); err != nil {
			log.Errorf("[RunMigrate-4] goose %v: %v", args[0], err)
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
