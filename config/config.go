package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func ConfigInit(cfgFile string) {

	// Set some defaults
	viper.SetDefault("log_level", "DEBUG")
	viper.SetDefault("server.bind_address", "127.0.0.1")
	viper.SetDefault("server.bind_port", "8080")
	viper.SetDefault("server.access_log", true)
	viper.SetDefault("backend.db_driver", "mysql")
	viper.SetDefault("poller.interval", 60)

	// Configuring and pulling overrides from environmental variables
	viper.SetEnvPrefix("VSUMMARY")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// set default config name and paths to look for it
	viper.SetConfigType("yaml")
	viper.SetConfigName("vsummary-config")
	viper.AddConfigPath("./")
	goPath := os.Getenv("GOPATH")
	if goPath != "" {
		viper.AddConfigPath(fmt.Sprintf("%s/src/github.com/gbolo/vsummary/testdata/sampleconfig", goPath))
	}

	// if the user provides a config file in a flag, lets use it
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	// If a config file is found, read it in.
	err := viper.ReadInConfig()

	// Kick-off the logging module
	loggingInit(viper.GetString("log_level"))

	if err == nil {
		log.Infof("using config file: %s", viper.ConfigFileUsed())
	} else {
		log.Warning(("no config file found: using environment variables and hard-coded defaults"))
	}

	// Print config in debug
	printConfigSummary()

	// Sanity checks
	sanityChecks()

	return

}

func ConfigInitPoller(cfgFile string) {

	// Set some defaults
	viper.SetDefault("log_level", "CRITICAL")
	viper.SetDefault("poller.interval", 60)

	// Configuring and pulling overrides from environmental variables
	viper.SetEnvPrefix("VSUMMARY")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// set default config name and paths to look for it
	viper.SetConfigType("yaml")
	viper.SetConfigName("vsummary-poller")
	viper.AddConfigPath("./")
	goPath := os.Getenv("GOPATH")
	if goPath != "" {
		viper.AddConfigPath(fmt.Sprintf("%s/src/github.com/gbolo/vsummary/testdata/sampleconfig", goPath))
	}

	// if the user provides a config file in a flag, lets use it
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	// If a config file is found, read it in.
	err := viper.ReadInConfig()

	// Kick-off the logging module
	loggingInit(viper.GetString("log_level"))

	if err == nil {
		log.Infof("using config file: %s", viper.ConfigFileUsed())
	} else {
		log.Warning(("no config file found: using environment variables and hard-coded defaults"))
	}

	return
}

func printConfigSummary() {

	log.Debugf("Configuration:\n")
	for _, c := range []string{
		"log_level",
		"server.bind_address",
		"server.bind_port",
		"backend.db_driver",
		"backend.db_dsn",
		"server.static_files_dir",
		"server.templates_dir",
	} {

		log.Debugf("%s: %s\n", c, viper.GetString(c))
	}
}

func sanityChecks() {

	if viper.GetString("backend.db_driver") != "mysql" && viper.GetString("backend.db_driver") != "" {

		log.Fatalf("only mysql is supported for backend_db_driver (value: %s)", viper.GetString("backend.db_driver"))
	}

	if viper.GetString("backend.db_dsn") == "" {

		log.Fatal("backend_db_dsn must be set")
	}

}
