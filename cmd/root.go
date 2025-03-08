package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	verbose   bool
	schemaDir string
)

var rootCmd = &cobra.Command{
	Use:   "graphscope <schema-directory>",
	Short: "Graphscope gives insights into your GraphQL schema",
}

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVarP(&schemaDir, "dir", "d", ".", "Specify the working directory")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
