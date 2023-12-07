#!/bin/sh

set -eux

GOTAFACE_SPANNER_DATABASE=projects/gotaface/instances/example/databases/example

go run github.com/Jumpaku/spanner/cmd/gaf-spanner-fetch-schema \
    -format=txt.tpl \
    -input-txt-tpl=/work/examples/cmd/gaf-spanner-fetch-schema/dml/dml.sql.tpl \
    -output=/work/examples/cmd/gaf-spanner-fetch-schema/dml/dml.sql \
    "${GOTAFACE_SPANNER_DATABASE}" \
    A B_1 B_2 B_3 B_4 C_1 C_2 C_3 C_4 C_5 D_1 E_1 E_2 F_1 F_2 F_3 G
