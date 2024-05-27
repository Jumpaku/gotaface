package schema

type Schemas []Table

type Table struct {
	Name        string       `json:"name"`
	View        bool         `json:"view" spanner:"IsView"`
	Columns     []Column     `json:"columns"`
	PrimaryKey  []string     `json:"primary_key"`
	Parent      string       `json:"parent"`
	ForeignKeys []ForeignKey `json:"foreign_key"`
	Indexes     []Index      `json:"indexes"`
}

type Column struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Nullable  bool   `json:"nullable"`
	Generated bool   `json:"generated"`
}

type ForeignKey struct {
	Name            string   `json:"name"`
	ReferencedTable string   `json:"referenced_table"`
	ReferencedKey   []string `json:"referenced_key"`
	ReferencingKey  []string `json:"referencing_key"`
}

type Index struct {
	Name   string      `json:"name"`
	Key    []*IndexKey `json:"key"`
	Unique bool        `json:"unique" spanner:"IsUnique"`
}

type IndexKey struct {
	Name string `json:"name"`
	Desc bool   `json:"desc" spanner:"IsDesc"`
}
