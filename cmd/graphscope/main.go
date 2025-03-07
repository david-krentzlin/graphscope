package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"

	//tea "github.com/charmbracelet/bubbletea"
	"github.com/david-krentzlin/graphscope/internal/analyzer"
	//"github.com/david-krentzlin/graphscope/internal/ui"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "graphscope <schema-directory>",
		Short: "Graphscope gives insights into your GraphQL schema",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			schemaDir := args[0]

			updateChan := make(chan analyzer.ProgressUpdate)
			analyzerInstance, err := analyzer.New(schemaDir, updateChan)
			if err != nil {
				log.Fatalf("Error creating analyzer: %v", err)
				return
			}
			bar := progressbar.Default(int64(analyzerInstance.TotalFiles), "Parsing schema files")
			go analyzerInstance.Parse()

			for {
				select {
				case update, more := <-updateChan:
					if !more {
						fmt.Println("✅ Parsing phase completed.")
						return
					}

					if update == nil {
						// Unexpected nil message, continue safely
						continue
					}

					switch update := update.(type) {
					case analyzer.ParseSchemaFile:
						// let's make the output relative to the schema directory

						relativePath, err := filepath.Rel(schemaDir, update.CurrentFile)
						bar.Add(1)
						bar.AddDetail(relativePath)
						if err != nil {
							log.Fatalf("Error getting relative path: %v", err)
						}
					case analyzer.ParseSchemaComplete:
						return

					}
				}
			}
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
