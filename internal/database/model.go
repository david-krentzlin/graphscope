package database

type Type struct {
	ID     uint   `gorm:"primaryKey"`
	Name   string `gorm:"unique;not null"`
	Kind   string `gorm:"not null;check:kind IN ('OBJECT', 'INPUT_OBJECT', 'SCALAR', 'ENUM', 'INTERFACE', 'UNION')"`
	Fields []Field
}

type Field struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	TypeID      uint `gorm:"not null"` // Foreign key to Type
	Type        Type `gorm:"foreignKey:TypeID"`
	//	IsNullable  bool        `gorm:"default:true"`
	// IsList      bool            `gorm:"default:false"`
	//	Arguments []FieldArgument `gorm:"foreignKey:FieldID"`
}

type FieldArgument struct {
	ID      uint   `gorm:"primaryKey"`
	Name    string `gorm:"not null"`
	FieldID uint   `gorm:"not null"` // Foreign key to Field
	TypeID  uint   `gorm:"not null"` // Foreign key to GraphQLType
	// DefaultValue string
	Description string

	Field Field `gorm:"foreignKey:FieldID"`
	Type  Type  `gorm:"foreignKey:TypeID"`
}

type Directive struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}

type TypeDirective struct {
	ID                 uint `gorm:"primaryKey"`
	TypeID             uint `gorm:"not null"` // Foreign key to Field
	DirectiveID        uint `gorm:"not null"` // Foreign key to Directive
	DirectiveArguments []TypeDirectiveArgument
}

type FieldDirective struct {
	ID                 uint `gorm:"primaryKey"`
	FieldID            uint `gorm:"not null"` // Foreign key to Field
	DirectiveID        uint `gorm:"not null"` // Foreign key to Directive
	DirectiveArguments []FieldDirectiveArgument
}

type FieldDirectiveArgument struct {
	ID               uint   `gorm:"primaryKey"`
	Name             string `gorm:"not null"`
	FieldDirectiveID uint   `gorm:"not null"` // Foreign key to Directive
	Value            string
}

type TypeDirectiveArgument struct {
	ID              uint   `gorm:"primaryKey"`
	Name            string `gorm:"not null"`
	TypeDirectiveID uint   `gorm:"not null"` // Foreign key to Directive
	Value           string
}
