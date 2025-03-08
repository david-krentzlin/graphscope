package cmd

import (
	"fmt"
	"log"

	"github.com/david-krentzlin/graphscope/internal/analyzer"
	"github.com/spf13/cobra"
)

var resolverCmd = &cobra.Command{
	Use:   "resolver",
	Short: "Query resolvers in your GraphQL schema",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := analyzePath(schemaDir)
		if verbose {
			fmt.Printf("Analyzed schema in %s\n", schemaDir)
		}

		if err != nil {
			log.Fatalf("Error analyzing schema: %v", err)
		}
	},
}

func init() {
	resolverCmd.Flags().StringP("name", "n", "", "Name of the resolver to query")
	resolverCmd.Flags().StringP("type", "t", "", "Type of the resolver to query")
	resolverCmd.Flags().StringP("path", "p", "", "Path of the resolver in case of REST resolvers")
	rootCmd.AddCommand(resolverCmd)
}

func analyzePath(schemaDir string) (*analyzer.Analyzer, error) {
	analyzerInstance, err := analyzer.New(schemaDir, nil)
	if err != nil {
		log.Fatalf("Error creating analyzer: %v", err)
		return nil, err
	}

	analyzerInstance.Ingest()
	return analyzerInstance, nil
}
