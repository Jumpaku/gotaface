package dml

import "reflect"

// DBValue represents values supported by a database.
type DBValue interface {
	// GoType returns type of value that DBValue represents.
	GoType() reflect.Type
	// Set sets val to DBValue. Set may panic if val is wrong.
	Set(val any)
	// Get assigns a value that DBValue represents to the destination referenced by ptr. Get may panic if ptr is wrong or incompatible.
	Get(ptr any)
}

type Row map[string]DBValue
type Rows []Row
