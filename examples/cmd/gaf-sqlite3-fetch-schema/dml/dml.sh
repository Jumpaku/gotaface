#!/bin/sh

set -eux

go run ./sqlite3/cmd/gaf-sqlite3-fetch-schema \
    -format=txt.tpl \
    -input-txt-tpl=./examples/cmd/gaf-sqlite3-fetch-schema/dml/dml.sql.tpl \
    -output=./examples/cmd/gaf-sqlite3-fetch-schema/dml/dml.sql \
    "example.db" \
    A C_1 C_2 C_3 C_4 C_5 D_1 E_1 E_2 F_1 F_2 F_3 G H I
