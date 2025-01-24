package database

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//	interface DB {
//		StoreType(t *GraphQLType) error
//		StoreField(f *Field) error
//		StoreFieldArgument(a *FieldArgument) error
//		StoreDirective(d *Directive) error
//	}

func NewDB() (*gorm.DB, error) {
	//db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	// delete test.db if it exists
	os.Remove("test.db")

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&Type{}, &Field{}, &Directive{}, &TypeDirective{}, &FieldDirective{}, &TypeDirectiveArgument{}, &FieldDirectiveArgument{}); err != nil {

		log.Fatalf("failed to migrate schema: %v", err)
	}

	return db, nil
}
