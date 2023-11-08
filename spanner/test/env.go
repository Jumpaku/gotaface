package test

import "os"

const (
	EnvTestSpannerProject  = "GOTAFACE_TEST_SPANNER_PROJECT"
	EnvTestSpannerInstance = "GOTAFACE_TEST_SPANNER_INSTANCE"
)

func GetEnvSpanner() (project string, instance string) {
	return os.Getenv(EnvTestSpannerProject), os.Getenv(EnvTestSpannerInstance)
}
