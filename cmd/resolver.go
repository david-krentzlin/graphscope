package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/david-krentzlin/graphscope/internal/analyzer"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
)

var (
	flagName       bool
	flagType       bool
	flagUrl        bool
	flagEdit       bool
	flagReferences bool
	flagSourcePath bool
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

		// analyzerInstance.DumpResolverIndex()
		// return
		resolvers := analyzerInstance.FindResolvers()
		if len(resolvers) == 0 {
			fmt.Println("No resolvers found")
			os.Exit(0)
		}

		maxNameLen := 0
		maxTypeLen := 0
		maxUrlLen := 0

		for _, resolver := range resolvers {
			maxNameLen = max(maxNameLen, len(resolver.Name()))
			maxTypeLen = max(maxTypeLen, len(resolver.Type()))
			maxUrlLen = max(maxUrlLen, len(resolver.Path()))
		}

		idx, err := fuzzyfinder.Find(resolvers, func(i int) string {
			return resolverOutputLine(resolvers[i], schemaDir, maxNameLen, maxTypeLen, maxUrlLen)
		})

		if idx == -1 {
			return
		}

		selectedResolver := resolvers[idx]
		resolverName := selectedResolver.Name()
		def := analyzerInstance.ResolverDefinition(resolverName)
		if def != nil {
			if flagReferences {
				if len(def.References) == 0 {
					fmt.Println("Resolver has no references")
					fmt.Println("At the moment indirect references through @router or @onAuthenticationState are not supported")
					return
				}

				// find the references rather than the definition
				idx, err = fuzzyfinder.Find(def.References, func(i int) string {
					return def.References[i].Path()
				})
				if idx == -1 {
					return
				}
				selectedReference := def.References[idx]
				lessCmd := exec.Command("less", fmt.Sprintf("+%d", selectedReference.Field.Position.Line), selectedReference.Field.Position.Src.Name)
				lessCmd.Stdout = os.Stdout
				lessCmd.Stderr = os.Stderr
				lessCmd.Run()

			} else {
				lessCmd := exec.Command("less", fmt.Sprintf("+%d", def.Definition.Position.Line), def.Definition.Position.Src.Name)
				lessCmd.Stdout = os.Stdout
				lessCmd.Stderr = os.Stderr
				lessCmd.Run()
			}
		}
	},
}

func resolverOutputLine(resolver *analyzer.Resolver, baseDir string, maxNameLen, maxTypeLen, maxUrlLen int) string {
	var output string

	if flagName {
		output += fmt.Sprintf("%s%*s", resolver.Name(), maxNameLen-len(resolver.Name()), "")
		output += " "
	}
	if flagType {
		output += fmt.Sprintf("%s%*s", resolver.Type(), maxTypeLen-len(resolver.Type()), "")
		output += " "
	}
	if flagUrl {
		output += fmt.Sprintf("%s%*s", resolver.Path(), maxUrlLen-len(resolver.Path()), "")
		output += " "
	}

	if flagSourcePath {
		output += resolver.RelativeSourcePath(baseDir)
		output += " "
	}
	return output
}

func init() {

	resolverCmd.Flags().BoolVarP(&flagName, "name", "n", true, "Show the name of the resolver")
	resolverCmd.Flags().BoolVarP(&flagType, "type", "t", false, "Show the type of the resolver")
	resolverCmd.Flags().BoolVarP(&flagUrl, "url", "u", false, "Show the URL of the resolver")
	resolverCmd.Flags().BoolVarP(&flagEdit, "edit", "e", false, "Edit the resolver")
	resolverCmd.Flags().BoolVarP(&flagReferences, "references", "r", false, "Find references rather than definitions")
	resolverCmd.Flags().BoolVarP(&flagSourcePath, "source-path", "s", false, "Show the source path of the resolver")

	rootCmd.AddCommand(resolverCmd)
}
