package analyzer

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/david-krentzlin/graphscope/internal/database"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
	"gorm.io/gorm"
)

type Analyzer struct {
	db          *gorm.DB
	tpeMap      map[string]database.Type
	FilesLoaded int
}

func New(db *gorm.DB) *Analyzer {
	return &Analyzer{db: db, tpeMap: make(map[string]database.Type)}
}

func (a *Analyzer) IngestFiles(dir string) error {
	schema, err := a.loadSchema(dir)
	if err != nil {
		return err
	}
	err = a.ingestTypes(schema)
	if err != nil {
		return err
	}

	// clean up the map
	for k := range a.tpeMap {
		delete(a.tpeMap, k)
	}
	return nil
}

func (a *Analyzer) loadSchema(dir string) (*ast.SchemaDocument, error) {
	var schema *ast.SchemaDocument = &ast.SchemaDocument{}

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(path) == ".graphql" {
			newSchema, err := a.parseFile(path)
			if err != nil {
				return fmt.Errorf("failed to parse file %s: %w", path, err)
			}
			schema.Merge(newSchema)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return schema, nil
}

func (a *Analyzer) parseFile(path string) (*ast.SchemaDocument, error) {
	content, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	a.FilesLoaded++

	return parser.ParseSchema(&ast.Source{
		Name:    path,
		Input:   string(content),
		BuiltIn: true,
	})
}

func (a *Analyzer) ingestNativeTypes() error {
	for _, scalar := range []string{"PhoneNumberE164", "String", "Int", "Long", "Markdown", "Email", "UploadAuthToken", "UploadId", "Float", "Boolean", "ID", "Date", "URL", "JSON", "URN", "GlobalID", "SlugOrID", "Expression", "UUID"} {
		tpe := database.Type{
			Name: scalar,
			Kind: "SCALAR",
		}

		if _, ok := a.tpeMap[scalar]; !ok {
			result := a.db.Create(&tpe)
			if result.Error != nil {
				return result.Error
			}
			a.tpeMap[tpe.Name] = tpe
		}
	}

	return nil
}

func (a *Analyzer) ingestTypes(schema *ast.SchemaDocument) error {
	// first make the types known
	ingenstErr := a.ingestNativeTypes()

	if ingenstErr != nil {
		return ingenstErr
	}

	for _, def := range schema.Definitions {
		tpe := database.Type{
			Name: def.Name,
			Kind: fmt.Sprint(def.Kind),
		}

		if _, ok := a.tpeMap[tpe.Name]; !ok {
			result := a.db.Create(&tpe)
			if result.Error != nil {
				return result.Error
			}
			a.tpeMap[def.Name] = tpe
		}
	}

	knownDirectives := make(map[string]uint)
	var directiveId uint

	for _, def := range schema.Definitions {
		tpe := a.tpeMap[def.Name]

		for _, directive := range def.Directives {
			if foundDirective, ok := knownDirectives[directive.Name]; !ok {
				d := database.Directive{
					Name: directive.Name,
				}
				a.db.Create(&d)
				directiveId = d.ID
				knownDirectives[directive.Name] = directiveId
			} else {
				directiveId = foundDirective
			}

			assoc := database.TypeDirective{
				TypeID:      tpe.ID,
				DirectiveID: directiveId,
			}

			result := a.db.Create(&assoc)
			if result.Error != nil {
				return result.Error
			}

			var arguments []database.TypeDirectiveArgument
			for _, arg := range directive.Arguments {
				arguments = append(arguments, database.TypeDirectiveArgument{
					Name:  arg.Name,
					Value: arg.Value.Raw,
				})
			}
		}

		for _, field := range def.Fields {
			fieldTpe, ok := a.tpeMap[field.Type.Name()]

			if !ok {
				fmt.Println("Field type not found. Likely scala defined type. Will skip for now", field.Type.Name())
				continue
			}

			fld := database.Field{
				Name:   field.Name,
				TypeID: fieldTpe.ID,
			}
			result := a.db.Create(&fld)
			if result.Error != nil {
				return result.Error
			}
		}
	}

	return nil
}
