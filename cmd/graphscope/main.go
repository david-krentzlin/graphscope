package main

import (
	"fmt"
	"log"
	"os"

	"github.com/david-krentzlin/graphscope/internal/analyzer"
	"github.com/david-krentzlin/graphscope/internal/database"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "graphscope <schema-director>",
		Short: "Graphscope gives insights into your GraphQL schema",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			schemaDir := args[0]

			// Initialize database
			db, err := database.NewDB()
			if err != nil {
				log.Fatalf("Failed to initialize database: %v", err)
			}

			a := analyzer.New(db)
			if err := a.IngestFiles(schemaDir); err != nil {
				log.Fatalf("Failed to process schema directory: %v", err)
			}

			fmt.Printf("Schema analysis complete. Loaded %d files", a.FilesLoaded)
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
