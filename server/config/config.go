package config

import (
	"os"
	"path/filepath"
)

var (
	HomePath     = os.Getenv("GEO_NOTIFY_SERVER_HOME")
	ResourcePath string
)

func init() {
	ResourcePath = filepath.Join(HomePath, "./resource")
}
