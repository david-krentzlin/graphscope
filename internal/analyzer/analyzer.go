package analyzer

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

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

type Analyzer struct {
	schemaDef   *ast.SchemaDocument
	updateChan  chan<- ProgressUpdate
	SourcePaths []string
	TotalFiles  int
	FilesParsed int
}

func New(baseDir string, updateChan chan<- ProgressUpdate) (*Analyzer, error) {
	sourcePaths, err := collectSourcePaths(baseDir)
	if err != nil {
		return nil, err
	}
	return &Analyzer{FilesParsed: 0, updateChan: updateChan, SourcePaths: sourcePaths, TotalFiles: len(sourcePaths)}, nil
}

func (a *Analyzer) Parse() error {
	schema, err := a.parseSchemaDefinitions()
	if err != nil {
		return err
	}
	a.schemaDef = schema
	a.updateChan <- ParseSchemaComplete{}
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
		a.updateChan <- ParseSchemaFile{TotalFiles: a.TotalFiles, CurrentFile: source, FilesParsed: a.FilesParsed}
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
