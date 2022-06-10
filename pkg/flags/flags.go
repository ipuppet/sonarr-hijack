package flags

import (
	"flag"
)

var (
	ConfigPath string
	LogPath    string
)

const (
	configPathDefault = "./config"
	configPathUsage   = "Set config file path."

	logPathDefault = "/go/log"
	logPathUsage   = "Set log file path."
)

func init() {
	flag.StringVar(&ConfigPath, "config", configPathDefault, configPathUsage)

	flag.StringVar(&LogPath, "log", logPathDefault, logPathUsage)

	flag.Parse()
}
