#!/bin/sh

export GOTAFACE_TEST_SPANNER_PROJECT=gotaface
export GOTAFACE_TEST_SPANNER_INSTANCE=test
go test ./spanner/schema/...