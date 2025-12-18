package cmd

import (
	"simple-golang/config"
	"simple-golang/internal/database/seed"

	"github.com/labstack/gommon/log"

	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed initial data into the database",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.NewConfig()

		db, err := cfg.ConnectionPostgres()
		if err != nil {
			log.Fatalf("[RunSeed-1] failed to connect to DB Gorm: %v", err)
		}

		seed.RunAll(db.DB)

		log.Infof("Database seeding completed successfully.")
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
}
