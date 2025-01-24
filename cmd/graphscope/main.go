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
		Short: "GraphScope analyzes GraphQL schemas with focus on directives",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			schemaDir := args[0]

			// Initialize database
			db, err := database.NewDB()
			if err != nil {
				log.Fatalf("Failed to initialize database: %v", err)
			}
			defer db.Close()

			// Initialize analyzer
			a := analyzer.New(db)

			// Process schema directory
			if err := a.ProcessDirectory(schemaDir); err != nil {
				log.Fatalf("Failed to process schema directory: %v", err)
			}

			fmt.Println("Schema analysis complete. Use queries to explore the schema.")
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
