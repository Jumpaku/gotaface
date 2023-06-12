#!/bin/sh

export GOTAFACE_TEST_SPANNER_PROJECT=gotaface
export GOTAFACE_TEST_SPANNER_INSTANCE=test
go test ./spanner/ddl/...
go test ./spanner/dml/delete/...
go test ./spanner/dml/dump/...
go test ./spanner/dml/insert/...