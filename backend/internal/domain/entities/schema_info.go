package entities

type ColumnInfo struct {
	Name             string  `json:"name"`
	DataType         string  `json:"data_type"`
	IsNullable       bool    `json:"is_nullable"`
	DefaultValue     *string `json:"default_value,omitempty"`
	MaxLength        *int    `json:"max_length,omitempty"`
	NumericPrecision *int    `json:"numeric_precision,omitempty"`
	NumericScale     *int    `json:"numeric_scale,omitempty"`
	OrdinalPosition  int     `json:"ordinal_position"`
}

type ForeignKeyInfo struct {
	ConstraintName    string   `json:"constraint_name"`
	ColumnNames       []string `json:"column_names"`
	ReferencedTable   string   `json:"referenced_table"`
	ReferencedColumns []string `json:"referenced_columns"`
	UpdateRule        string   `json:"update_rule"`
	DeleteRule        string   `json:"delete_rule"`
}

type IndexInfo struct {
	Name        string   `json:"name"`
	ColumnNames []string `json:"column_names"`
	IsUnique    bool     `json:"is_unique"`
	IsPrimary   bool     `json:"is_primary"`
	IndexType   string   `json:"index_type"`
}

type TableSchemaInfo struct {
	Columns     []ColumnInfo     `json:"columns"`
	PrimaryKeys []string         `json:"primary_keys"`
	ForeignKeys []ForeignKeyInfo `json:"foreign_keys"`
	Indexes     []IndexInfo      `json:"indexes"`
}
