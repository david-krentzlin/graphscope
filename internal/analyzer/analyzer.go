package analyzer

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
)

type ProgressUpdate interface {
	isProgressUpdate()
}

type ParseSchemaFile struct {
	TotalFiles  int
	CurrentFile string
	FilesParsed int
}

func (ParseSchemaFile) isProgressUpdate() {}

type ParseSchemaComplete struct{}

func (ParseSchemaComplete) isProgressUpdate() {}

type ExtractresolversComplete struct{}
type ExtractResolverDefinition struct {
	Name           string
	DefinitionName string
}
type ExtractResolverUsage struct {
	ResolverName string
	FieldPath    string
}

func (ExtractresolversComplete) isProgressUpdate()  {}
func (ExtractResolverDefinition) isProgressUpdate() {}
func (ExtractResolverUsage) isProgressUpdate()      {}

type Resolver struct {
	Definition *ast.Definition
	References []*ast.FieldDefinition
}

type ResolverFilter struct {
	name *string
	tpe  *string
	path *string
}

type Analyzer struct {
	schemaDef     *ast.SchemaDocument
	updateChan    chan<- ProgressUpdate
	BaseDir       string
	SourcePaths   []string
	TotalFiles    int
	FilesParsed   int
	ResolverIndex map[string]*Resolver
}

func New(baseDir string, updateChan chan<- ProgressUpdate) (*Analyzer, error) {
	sourcePaths, err := collectSourcePaths(baseDir)
	if err != nil {
		return nil, err
	}
	return &Analyzer{FilesParsed: 0, updateChan: updateChan, BaseDir: baseDir, SourcePaths: sourcePaths, TotalFiles: len(sourcePaths)}, nil
}

func (a *Analyzer) Ingest() error {
	if err := a.Parse(); err != nil {
		return err
	}
	if err := a.IndexResolvers(); err != nil {
		return err
	}
	return nil
}

func (a *Analyzer) Parse() error {
	schema, err := a.parseSchemaDefinitions()
	if err != nil {
		return err
	}
	a.schemaDef = schema
	a.sendUpdate(ParseSchemaComplete{})
	return nil
}

func (a *Analyzer) IndexResolvers() error {
	a.ResolverIndex = make(map[string]*Resolver)
	var indexable bool

	for _, def := range a.schemaDef.Definitions {
		indexable = false

		if def.Kind == ast.InputObject {
			for _, directive := range def.Directives {
				switch directive.Name {
				case "httpGet":
					indexable = true
					a.ResolverIndex[strings.ToLower(def.Name)] = &Resolver{Definition: def}
				case "httpGetBatched":
					indexable = true
					a.ResolverIndex[strings.ToLower(def.Name)] = &Resolver{Definition: def}
				case "httpPost":
					indexable = true
					a.ResolverIndex[strings.ToLower(def.Name)] = &Resolver{Definition: def}
				case "httpPut":
					indexable = true
					a.ResolverIndex[strings.ToLower(def.Name)] = &Resolver{Definition: def}
				case "httpPatch":
					indexable = true
					a.ResolverIndex[strings.ToLower(def.Name)] = &Resolver{Definition: def}
				case "httpDelete":
					indexable = true
					a.ResolverIndex[strings.ToLower(def.Name)] = &Resolver{Definition: def}
				default:
					indexable = false
				}

				if indexable {
					a.sendUpdate(ExtractResolverDefinition{Name: def.Name, DefinitionName: directive.Name})
				}
			}
		}
	}

	for _, def := range a.schemaDef.Definitions {
		if def.Kind == ast.Object {
			for _, field := range def.Fields {
				for _, directive := range field.Directives {
					// check if the directive is a resolver directive
					if idx, found := a.ResolverIndex[strings.ToLower(directive.Name)]; found {
						idx.References = append(idx.References, field)
						a.sendUpdate(ExtractResolverUsage{ResolverName: directive.Name, FieldPath: fmt.Sprintf("%s.%s", def.Name, field.Name)})
					}
				}
			}
		}
	}

	a.sendUpdate(ExtractresolversComplete{})
	return nil
}

func collectSourcePaths(baseDir string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".graphql" {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return paths, nil
}

func (a *Analyzer) parseSchemaDefinitions() (*ast.SchemaDocument, error) {

	var schema *ast.SchemaDocument = &ast.SchemaDocument{}
	for _, source := range a.SourcePaths {
		newSchema, err := a.parseFile(source)
		if err != nil {
			return nil, fmt.Errorf("failed to parse file %s: %w", source, err)
		}
		a.FilesParsed++
		relativeSource := strings.TrimPrefix(source, a.BaseDir)
		a.sendUpdate(ParseSchemaFile{TotalFiles: a.TotalFiles, CurrentFile: relativeSource, FilesParsed: a.FilesParsed})
		schema.Merge(newSchema)
	}

	return schema, nil
}

func (a *Analyzer) parseFile(path string) (*ast.SchemaDocument, error) {
	content, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return parser.ParseSchema(&ast.Source{
		Name:    path,
		Input:   string(content),
		BuiltIn: true,
	})
}

func (a *Analyzer) sendUpdate(update ProgressUpdate) {
	if a.updateChan == nil {
		return
	}
	a.updateChan <- update
}
