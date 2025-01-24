package analyzer

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/david-krentzlin/graphscope/internal/database"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
)

type Analyzer struct {
	db          *database.DB
	FilesLoaded int
}

func New(db *database.DB) *Analyzer {
	return &Analyzer{db: db}
}

func (a *Analyzer) ProcessDirectory(dir string) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(path) == ".graphql" {
			return a.processFile(path)
		}

		return nil
	})
}

func (a *Analyzer) processFile(path string) error {
	content, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	schemaDoc, err := parser.ParseSchema(&ast.Source{
		Name:    path,
		Input:   string(content),
		BuiltIn: false,
	})
	if err != nil {
		return err
	}

	a.FilesLoaded++
	return a.storeSchema(schemaDoc)
}

func (a *Analyzer) storeSchema(schema *ast.SchemaDocument) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Store types
	for _, def := range schema.Definitions {
		if err := a.storeType(tx, def); err != nil {
			return err
		}
	}

	// Store directives
	for _, dir := range schema.Directives {
		if err := a.storeDirective(tx, dir); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (a *Analyzer) storeType(tx *sql.Tx, def *ast.Definition) error {
	// Implementation for storing types
	fmt.Printf("Storing type %s\n", def.Name)
	return nil
}

func (a *Analyzer) storeDirective(tx *sql.Tx, dir *ast.DirectiveDefinition) error {
	// Implementation for storing directives
	fmt.Printf("Storing directive %s\n", dir.Name)

	return nil
}
