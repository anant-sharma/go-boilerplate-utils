package config

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Load - Loads config
func Load(c interface{}) {
	// Set the path to look for the configurations file
	viper.AddConfigPath("config")

	// Set the file name of the configurations file
	viper.SetConfigName("")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	// Convert _ underscore in env to . dot notation in viper
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(c)
	if err != nil {
		log.Fatalf("Unable to load config to struct, %v", err)
	}
}
