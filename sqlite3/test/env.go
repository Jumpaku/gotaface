package test

import "os"

const (
	EnvTestSQLite3 = "GOTAFACE_TEST_SQLITE3"
)

func GetEnvSQLite3() (database string) {
	return os.Getenv(EnvTestSQLite3)
}
