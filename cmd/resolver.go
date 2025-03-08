package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/david-krentzlin/graphscope/internal/analyzer"
	"github.com/spf13/cobra"
)

var (
	flagName bool
	flagType bool
	flagUrl  bool
	flagEdit bool
)

var resolverCmd = &cobra.Command{
	Use:   "resolver",
	Short: "Query resolvers in your GraphQL schema",
	Run: func(cmd *cobra.Command, args []string) {
		analyzerInstance, err := analyzer.New(schemaDir, nil)

		if err != nil {
			log.Fatalf("Error creating analyzer: %v", err)
			os.Exit(1)
		}

		err = analyzerInstance.Ingest()
		if err != nil {
			log.Fatalf("Error ingesting schema: %v", err)
			os.Exit(1)
		}

		resolvers := analyzerInstance.FindResolvers()
		if len(resolvers) == 0 {
			fmt.Println("No resolvers found")
			os.Exit(0)
		}

		candidates := []string{}
		for _, resolver := range resolvers {
			candidates = append(candidates, resolverOutputLine(resolver))
		}

		fzf := exec.Command("fzf")
		fzf.Stdin = strings.NewReader(strings.Join(candidates, "\n"))
		output, err := fzf.Output()
		if err != nil {
			fmt.Println("Error running fzf:", err)
			return
		}

		selected := strings.TrimSpace(string(output))
		if selected == "" {
			fmt.Println("No selection made")
			return
		} else {
			resolverName := strings.Split(selected, " ")[0]
			def := analyzerInstance.ResolverDefinition(resolverName)
			if def != nil {
				//use less to display the file
				lessCmd := exec.Command("less", fmt.Sprintf("+%d", def.Definition.Position.Line), def.Definition.Position.Src.Name)
				lessCmd.Stdout = os.Stdout
				lessCmd.Stderr = os.Stderr
				lessCmd.Run()
			}
		}
	},
}

func resolverOutputLine(resolver *analyzer.Resolver) string {
	var output string
	if flagName {
		output += resolver.Name()
		output += " "
	}
	if flagType {
		output += resolver.Type()
		output += " "
	}
	if flagUrl {
		output += resolver.Path()
		output += " "
	}
	return output
}

func init() {

	resolverCmd.Flags().BoolVarP(&flagName, "name", "n", true, "Show the name of the resolver")
	resolverCmd.Flags().BoolVarP(&flagType, "type", "t", false, "Show the type of the resolver")
	resolverCmd.Flags().BoolVarP(&flagUrl, "url", "u", false, "Show the URL of the resolver")
	resolverCmd.Flags().BoolVarP(&flagEdit, "edit", "e", false, "Edit the resolver")

	rootCmd.AddCommand(resolverCmd)
}
